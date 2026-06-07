package api

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	errRoomFullRPC     = connect.NewError(connect.CodeFailedPrecondition, errors.New("room is full"))
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

func (app *APIApplication) JoinRoom(
	ctx context.Context,
	req *connect.Request[roomv1.JoinRoomRequest],
) (*connect.Response[roomv1.JoinRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	code := normalizeRoomCode(req.Msg.GetCode())
	if code == "" {
		return nil, errRoomCodeInvalid
	}

	if _, ok := app.RoomsStore.GetByCode(code); !ok {
		return nil, errRoomCodeInvalid
	}

	if err := app.RoomsStore.AddPlayer(code, person.ID); err != nil {
		if errors.Is(err, errRoomFull) {
			return nil, errRoomFullRPC
		}
		if errors.Is(err, errRoomNotFound) {
			return nil, errRoomCodeInvalid
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&roomv1.JoinRoomResponse{}), nil
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

	if !app.RoomsStore.HasPlayer(code, person.ID) {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not in this room"))
	}

	var hostPerson *accountv1.Account
	playersAccounts := make([]*accountv1.Account, 0, len(room.PlayerIDs))
	for _, playerID := range room.PlayerIDs {
		p, err := app.Queries.GetPersonByID(ctx, playerID)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				continue
			}
			return nil, connect.NewError(connect.CodeInternal, err)
		}
		personToAccount := personToAccount(&p)
		if playerID == room.HostPersonID {
			hostPerson = personToAccount
		}
		playersAccounts = append(playersAccounts, personToAccount)
	}

	return connect.NewResponse(&roomv1.GetRoomResponse{
		Code:       room.Code,
		GameKey:    room.GameKey,
		HostPlayer: hostPerson,
	}), nil
}

func (app *APIApplication) LeaveRoom(
	ctx context.Context,
	req *connect.Request[roomv1.LeaveRoomRequest],
) (*connect.Response[roomv1.LeaveRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	code := normalizeRoomCode(req.Msg.GetCode())
	if code == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("room code required"))
	}

	if room, ok := app.RoomsStore.GetByCode(code); ok {
		if member, exists := room.Members[person.ID]; exists {
			room.RemoveMember(member)
		}
	}

	if err := app.RoomsStore.RemovePlayer(code, person.ID); err != nil {
		if errors.Is(err, errNotInRoom) || errors.Is(err, errRoomNotFound) {
			return connect.NewResponse(&roomv1.LeaveRoomResponse{Ok: true}), nil
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&roomv1.LeaveRoomResponse{Ok: true}), nil
}

func normalizeRoomCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
