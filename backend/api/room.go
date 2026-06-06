package api

import (
	"errors"
	"slices"
	"sync"
	"time"

	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/random"
	"github.com/arian-nj/chigame/backend/internals/socket"
)

const (
	roomCodeLength    = 6
	roomLifetime      = 12 * time.Hour
	maxPlayersPerRoom = 2
)

var AllowedRoomGameKeys = map[string]struct{}{
	"tic-tac-toe": {},
	"connect-4":   {},
}

var (
	errNotInRoom    = errors.New("not in room")
	errRoomFull     = errors.New("room is full")
	errRoomNotFound = errors.New("room not found")
)

type Room struct {
	Code         string
	GameKey      string
	HostPersonID int64
	PlayerIDs    []int64
	Members      map[int64]*RoomMember
	CreatedAt    time.Time
	ExpiresAt    time.Time

	MsgChnl chan *RoomEvent
}

func NewRoom(code string, hostPersonID int64, gameKey string) *Room {
	now := time.Now()
	expiresAt := now.Add(roomLifetime)
	return &Room{
		Code:         code,
		GameKey:      gameKey,
		HostPersonID: hostPersonID,
		PlayerIDs:    []int64{hostPersonID},
		Members:      make(map[int64]*RoomMember),
		CreatedAt:    now,
		ExpiresAt:    expiresAt,
		MsgChnl:      make(chan *RoomEvent, 10),
	}
}

func (r *Room) Handler() {
	for {
		newRoomEvent := <-r.MsgChnl
		switch newRoomEvent.Event.Content.(type) {
		case *roomv1.RoomMessage_ChatReq:
			r.SocketRequestSendMsg(newRoomEvent.Player, newRoomEvent.Event.GetChatReq())
			// r.handleChatMessageRequest(newEvent.Player, newEvent.Event.ChatMessageRequest)
		}

	}
}

// Room Store

type RoomsStore struct {
	mu     sync.Mutex
	byCode map[string]*Room
}

func NewRoomsStore() *RoomsStore {
	return &RoomsStore{
		byCode: make(map[string]*Room),
	}
}

func (s *RoomsStore) CreateRoom(hostPersonID int64, gameKey string) (*Room, error) {
	for range 8 {
		code := random.GenerateInviteCode(roomCodeLength)
		if _, exists := s.byCode[code]; exists {
			continue
		}
		newRoom := NewRoom(code, hostPersonID, gameKey)

		return newRoom, nil
	}
	return nil, errors.New("could not allocate room code")
}

func (s *RoomsStore) AddRoom(room *Room) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.byCode[room.Code] = room
}

func (s *RoomsStore) GetByCode(code string) (*Room, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getByCodeLocked(code)
}

func (s *RoomsStore) HasPlayer(code string, personID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.getByCodeLocked(code)
	if !ok {
		return false
	}
	return slices.Contains(room.PlayerIDs, personID)
}

func (s *RoomsStore) AddPlayer(code string, personID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.getByCodeLocked(code)
	if !ok {
		return errRoomNotFound
	}

	if slices.Contains(room.PlayerIDs, personID) {
		return nil
	}
	if len(room.PlayerIDs) >= maxPlayersPerRoom {
		return errRoomFull
	}

	room.PlayerIDs = append(room.PlayerIDs, personID)
	return nil
}

func (s *RoomsStore) RemovePlayer(code string, personID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.getByCodeLocked(code)
	if !ok {
		return errRoomNotFound
	}
	if !slices.Contains(room.PlayerIDs, personID) {
		return errNotInRoom
	}

	room.PlayerIDs = removePlayerID(room.PlayerIDs, personID)
	if len(room.PlayerIDs) == 0 || room.HostPersonID == personID {
		s.deleteLocked(code)
	}
	return nil
}

func (s *RoomsStore) deleteLocked(code string) {
	delete(s.byCode, code)
}

func (s *RoomsStore) getByCodeLocked(code string) (*Room, bool) {
	room, ok := s.byCode[code]
	if !ok {
		return nil, false
	}
	if time.Now().After(room.ExpiresAt) {
		s.deleteLocked(code)
		return nil, false
	}
	return room, true
}

func removePlayerID(ids []int64, personID int64) []int64 {
	out := ids[:0]
	for _, id := range ids {
		if id != personID {
			out = append(out, id)
		}
	}
	return out
}

// Room Member
type RoomMember struct {
	PersonID int64
	JoinedAt time.Time
	Socket   *socket.Socket
}

func NewRoomMember(personID int64, socket *socket.Socket) *RoomMember {
	return &RoomMember{
		PersonID: personID,
		JoinedAt: time.Now(),
		Socket:   socket,
	}
}

// Room Event
type RoomEvent struct {
	Player *RoomMember
	Event  *roomv1.RoomMessage
}

func NewRoomEvent(player *RoomMember, event *roomv1.RoomMessage) *RoomEvent {
	return &RoomEvent{
		Player: player,
		Event:  event,
	}
}
