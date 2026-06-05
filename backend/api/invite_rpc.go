package api

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	invitev1 "github.com/arian-nj/chigame/backend/gen/invite/v1"
	"github.com/jackc/pgx/v5"
)

var (
	errInviteRoomFull    = connect.NewError(connect.CodeFailedPrecondition, errors.New("room is full"))
	errInviteCodeInvalid = connect.NewError(connect.CodeNotFound, errors.New("invalid or expired invite code"))
	errInvalidGameKey    = connect.NewError(connect.CodeInvalidArgument, errors.New("invalid game key"))
)

func (app *APIApplication) CreateInviteRoom(
	ctx context.Context,
	req *connect.Request[invitev1.CreateInviteRoomRequest],
) (*connect.Response[invitev1.CreateInviteRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	gameKey := strings.TrimSpace(req.Msg.GetGameKey())
	if gameKey != "" {
		if _, ok := AllowedInviteGameKeys[gameKey]; !ok {
			return nil, errInvalidGameKey
		}
	}

	code, err := app.InviteStore.Create(person.ID, gameKey)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&invitev1.CreateInviteRoomResponse{
		InviteCode: code,
	}), nil
}

func (app *APIApplication) JoinInviteRoom(
	ctx context.Context,
	req *connect.Request[invitev1.JoinInviteRoomRequest],
) (*connect.Response[invitev1.JoinInviteRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	code := normalizeInviteCode(req.Msg.GetInviteCode())
	if code == "" {
		return nil, errInviteCodeInvalid
	}

	if _, ok := app.InviteStore.GetByCode(code); !ok {
		return nil, errInviteCodeInvalid
	}

	if err := app.InviteStore.AddPlayer(code, person.ID); err != nil {
		if errors.Is(err, errInviteStoreFull) {
			return nil, errInviteRoomFull
		}
		if errors.Is(err, errInviteRoomNotFound) {
			return nil, errInviteCodeInvalid
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&invitev1.JoinInviteRoomResponse{}), nil
}

func (app *APIApplication) GetInviteRoom(
	ctx context.Context,
	req *connect.Request[invitev1.GetInviteRoomRequest],
) (*connect.Response[invitev1.GetInviteRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	code := normalizeInviteCode(req.Msg.GetInviteCode())
	if code == "" {
		return nil, errInviteCodeInvalid
	}

	room, ok := app.InviteStore.GetByCode(code)
	if !ok {
		return nil, errInviteCodeInvalid
	}

	if !app.InviteStore.HasPlayer(code, person.ID) {
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

	return connect.NewResponse(&invitev1.GetInviteRoomResponse{
		InviteCode: room.Code,
		GameKey:    room.GameKey,
		Players:    playersAccounts,
		HostPlayer: hostPerson,
	}), nil
}

func (app *APIApplication) LeaveInviteRoom(
	ctx context.Context,
	req *connect.Request[invitev1.LeaveInviteRoomRequest],
) (*connect.Response[invitev1.LeaveInviteRoomResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	code := normalizeInviteCode(req.Msg.GetInviteCode())
	if code == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invite code required"))
	}

	if err := app.InviteStore.RemovePlayer(code, person.ID); err != nil {
		if errors.Is(err, errNotInInviteRoom) || errors.Is(err, errInviteRoomNotFound) {
			return connect.NewResponse(&invitev1.LeaveInviteRoomResponse{Ok: true}), nil
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&invitev1.LeaveInviteRoomResponse{Ok: true}), nil
}

func normalizeInviteCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}
