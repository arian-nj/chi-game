package api

import (
	"context"
	"time"

	"github.com/arian-nj/chigame/backend/database"
)

type Room struct {
	ID        int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (app *APIApplication) NewRoom(id int64) *Room {
	return &Room{
		ID:        id,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
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
