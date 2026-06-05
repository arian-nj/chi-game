package api

import (
	"errors"
	"sync"
	"time"

	"github.com/arian-nj/chigame/backend/internals/random"
)

const (
	inviteCodeLength   = 6
	inviteRoomLifetime = 12 * time.Hour
	maxPlayersPerRoom  = 2
)

var AllowedInviteGameKeys = map[string]struct{}{
	"tic-tac-toe": {},
	"connect-4":   {},
}

var (
	errNotInInviteRoom    = errors.New("not in invite room")
	errInviteStoreFull    = errors.New("invite room is full")
	errInviteRoomNotFound = errors.New("invite room not found")
)

type InviteRoom struct {
	Code         string
	GameKey      string
	HostPersonID int64
	PlayerIDs    []int64
	CreatedAt    time.Time
	ExpiresAt    time.Time
}

type InviteStore struct {
	mu     sync.Mutex
	byCode map[string]*InviteRoom
}

func NewInviteStore() *InviteStore {
	return &InviteStore{
		byCode: make(map[string]*InviteRoom),
	}
}

func (s *InviteStore) Create(hostPersonID int64, gameKey string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	for range 8 {
		code := random.GenerateInviteCode(inviteCodeLength)
		if _, exists := s.byCode[code]; exists {
			continue
		}

		s.byCode[code] = &InviteRoom{
			Code:         code,
			GameKey:      gameKey,
			HostPersonID: hostPersonID,
			PlayerIDs:    []int64{hostPersonID},
			CreatedAt:    now,
			ExpiresAt:    now.Add(inviteRoomLifetime),
		}
		return code, nil
	}
	return "", errors.New("could not allocate invite code")
}

func (s *InviteStore) GetByCode(code string) (*InviteRoom, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.getByCodeLocked(code)
}

func (s *InviteStore) HasPlayer(code string, personID int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.getByCodeLocked(code)
	if !ok {
		return false
	}
	return hasPlayerID(room.PlayerIDs, personID)
}

func (s *InviteStore) AddPlayer(code string, personID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.getByCodeLocked(code)
	if !ok {
		return errInviteRoomNotFound
	}

	if hasPlayerID(room.PlayerIDs, personID) {
		return nil
	}
	if len(room.PlayerIDs) >= maxPlayersPerRoom {
		return errInviteStoreFull
	}

	room.PlayerIDs = append(room.PlayerIDs, personID)
	return nil
}

func (s *InviteStore) RemovePlayer(code string, personID int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	room, ok := s.getByCodeLocked(code)
	if !ok {
		return errInviteRoomNotFound
	}
	if !hasPlayerID(room.PlayerIDs, personID) {
		return errNotInInviteRoom
	}

	room.PlayerIDs = removePlayerID(room.PlayerIDs, personID)
	if len(room.PlayerIDs) == 0 || room.HostPersonID == personID {
		s.deleteLocked(code)
	}
	return nil
}

func (s *InviteStore) deleteLocked(code string) {
	delete(s.byCode, code)
}

func (s *InviteStore) getByCodeLocked(code string) (*InviteRoom, bool) {
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

func hasPlayerID(ids []int64, personID int64) bool {
	for _, id := range ids {
		if id == personID {
			return true
		}
	}
	return false
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
