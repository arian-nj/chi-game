package api

import (
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
	"github.com/arian-nj/chigame/backend/internals/random"
	"github.com/arian-nj/chigame/backend/internals/socket"
	"github.com/arian-nj/chigame/backend/internals/utils"
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
	errRoomFull = errors.New("room is full")
)

type Room struct {
	ID           int64
	Code         string
	HostPersonID int64
	Members      map[int64]*RoomMember

	GameKey string

	CreatedAt time.Time
	ExpiresAt time.Time

	Commander *commander.Commander

	MsgChnl  chan *RoomEvent
	TaskSync sync.Mutex

	Queries *database.Queries
}

func NewRoom(code string, hostPersonID int64, gameKey string, queries *database.Queries) *Room {
	now := time.Now()
	expiresAt := now.Add(roomLifetime)
	return &Room{
		Code:         code,
		GameKey:      gameKey,
		HostPersonID: hostPersonID,
		Members:      make(map[int64]*RoomMember),
		CreatedAt:    now,
		ExpiresAt:    expiresAt,

		MsgChnl: make(chan *RoomEvent, 10),

		Commander: commander.NewCommander(),

		Queries: queries,
	}
}

func (app *APIApplication) RunRoom(r *Room) {
	utils.RunBackgroundTask(func() {
		for {
			select {
			case roomEvent := <-r.MsgChnl:
				func() {
					r.TaskSync.Lock()
					defer r.TaskSync.Unlock()

					switch msg := roomEvent.Event.Content.(type) {
					case *roomv1.RoomMessage_ChatReq:
						r.Commander.PushCommand(NewRoomMessageCommand(r, roomEvent.Player, msg.ChatReq))
					default:
						slog.Error("unhandled room event", "event", roomEvent.Event)
						_ = msg // Avoid unused variable warning if not handling other cases
					}
				}()
			case <-r.Commander.CommandNotifire:
				if len(r.Commander.Commands) > 0 {
					com := r.Commander.PopCommand()
					r.Commander.ApplyCommand(com)
				}
			}
		}
	})
}

func (r *Room) AddMember(member *RoomMember) error {
	r.TaskSync.Lock()
	defer r.TaskSync.Unlock()

	if _, exists := r.Members[member.Person.ID]; exists {
		return nil
	}
	if len(r.Members) >= maxPlayersPerRoom {
		return errRoomFull
	}

	r.Members[member.Person.ID] = member
	r.Commander.PushCommand(NewRoomMemberJoinedCommand(r, member))
	return nil
}

func (r *Room) RemoveMember(member *RoomMember) {
	r.TaskSync.Lock()
	defer r.TaskSync.Unlock()

	if _, ok := r.Members[member.Person.ID]; !ok {
		return
	}
	r.Commander.PushCommand(NewRoomMemberLeftCommand(r, member))
	delete(r.Members, member.Person.ID)
	if member.Socket != nil {
		member.Socket.Cancel()
	}
}

func (r *Room) HasMember(personID int64) bool {
	r.TaskSync.Lock()
	defer r.TaskSync.Unlock()
	_, ok := r.Members[personID]
	return ok
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

func (s *RoomsStore) CreateRoom(hostPersonID int64, gameKey string, queries *database.Queries) (*Room, error) {
	for range 8 {
		code := random.GenerateInviteCode(roomCodeLength)
		if _, exists := s.GetByCode(code); exists {
			continue
		}
		newRoom := NewRoom(code, hostPersonID, gameKey, queries)

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

func (s *RoomsStore) getByCodeLocked(code string) (*Room, bool) {
	room, ok := s.byCode[code]
	if !ok {
		return nil, false
	}
	if time.Now().After(room.ExpiresAt) {
		delete(s.byCode, code)
		return nil, false
	}
	return room, true
}

// Room Member
type RoomMember struct {
	Person   *database.Person
	JoinedAt time.Time
	Socket   *socket.Socket
}

func NewRoomMember(person *database.Person, socket *socket.Socket) *RoomMember {
	return &RoomMember{
		Person:   person,
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
