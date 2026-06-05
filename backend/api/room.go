package api

import (
	"context"
	"sync"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/socket"
)

// AllRooms
type RoomsStore struct {
	sync.RWMutex
	rooms map[int64]*Room
}

func NewRoomsStore() *RoomsStore {
	return &RoomsStore{
		rooms:   make(map[int64]*Room),
		RWMutex: sync.RWMutex{},
	}
}

// Room
type Room struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time

	MsgChnl chan *RoomEvent
}

func (app *APIApplication) NewRoom(id int64) *Room {
	return &Room{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		MsgChnl:   make(chan *RoomEvent, 10),
	}
}

func (app *APIApplication) CreateRoom(ctx context.Context, personID int64) (*Room, error) {
	roomRow, err := app.Queries.InsertRoom(ctx)
	if err != nil {
		return nil, err
	}
	room := app.NewRoom(roomRow.ID)
	return room, nil
}

func (app *APIApplication) JoinRoom(ctx context.Context, roomID, personID int64) error {
	_, err := app.Queries.InsertRoomPlayer(ctx, database.InsertRoomPlayerParams{
		RoomID:   roomID,
		PersonID: personID,
	})
	return err
}

func (app *APIApplication) LeaveRoom(ctx context.Context, roomID, personID int64) error {
	err := app.Queries.DeleteRoomPlayer(ctx, database.DeleteRoomPlayerParams{
		RoomID:   roomID,
		PersonID: personID,
	})
	return err
}

// Room Member
type RoomMember struct {
	PersonID int64
	JoinedAt time.Time
	Socket   *socket.Socket
}

func NewRoomMember(personID int64) *RoomMember {
	return &RoomMember{
		PersonID: personID,
		JoinedAt: time.Now(),
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
