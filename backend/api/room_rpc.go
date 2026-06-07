package api

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	"github.com/arian-nj/chigame/backend/database"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errRoomCodeInvalid = connect.NewError(connect.CodeNotFound, errors.New("invalid or expired room code"))
	errInvalidGameKey  = connect.NewError(connect.CodeInvalidArgument, errors.New("invalid game key"))
)

func (app *APIApplication) CreateRoom(
	ctx context.Context,
	req *connect.Request[roomv1.CreateRoomRequest],
) (*connect.Response[roomv1.CreateRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	gameKey := strings.TrimSpace(req.Msg.GetGameKey())
	if gameKey != "" {
		if _, ok := AllowedRoomGameKeys[gameKey]; !ok {
			return nil, errInvalidGameKey
		}
	}

	newRoom, err := app.RoomsStore.CreateRoom(person.ID, gameKey, app.Queries)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	roomRow, err := app.Queries.InsertRoom(ctx, database.InsertRoomParams{
		Code:         newRoom.Code,
		HostPersonID: newRoom.HostPersonID,
		ExpiresAt:    pgtype.Timestamp{Time: newRoom.ExpiresAt, Valid: true},
	})
	if err != nil {
		slog.Error("failed to insert room", "error", err)
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("failed to create room"))
	}
	newRoom.ID = roomRow.ID
	app.RoomsStore.AddRoom(newRoom)
	app.RunRoom(newRoom)

	return connect.NewResponse(&roomv1.CreateRoomResponse{
		Code: newRoom.Code,
	}), nil
}

func (app *APIApplication) GetRoom(
	ctx context.Context,
	req *connect.Request[roomv1.GetRoomRequest],
) (*connect.Response[roomv1.GetRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	code := normalizeRoomCode(req.Msg.GetCode())
	if code == "" {
		return nil, errRoomCodeInvalid
	}

	room, ok := app.RoomsStore.GetByCode(code)
	if !ok {
		return nil, errRoomCodeInvalid
	}

	if !room.HasMember(person.ID) && person.ID != room.HostPersonID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not in this room"))
	}

	hostRow, err := app.Queries.GetPersonByID(ctx, room.HostPersonID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errRoomCodeInvalid
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	hostAccount := personToAccount(&hostRow)
	players := roomConnectedPlayers(room)
	if !playersIncludeID(players, room.HostPersonID) {
		players = append(players, hostAccount)
	}

	return connect.NewResponse(&roomv1.GetRoomResponse{
		Code:       room.Code,
		GameKey:    room.GameKey,
		HostPlayer: hostAccount,
		Players:    players,
	}), nil
}

func roomConnectedPlayers(room *Room) []*accountv1.Account {
	room.TaskSync.Lock()
	defer room.TaskSync.Unlock()

	players := make([]*accountv1.Account, 0, len(room.Members))
	for _, member := range room.Members {
		players = append(players, personToAccount(member.Person))
	}
	return players
}

func playersIncludeID(players []*accountv1.Account, id int64) bool {
	for _, player := range players {
		if player.GetId() == id {
			return true
		}
	}
	return false
}

func normalizeRoomCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
