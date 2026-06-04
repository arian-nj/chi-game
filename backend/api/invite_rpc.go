package api

import (
	"context"
	"errors"
	"strings"
	"time"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	invitev1 "github.com/arian-nj/chigame/backend/gen/invite/v1"
	"github.com/arian-nj/chigame/backend/internals/random"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

const (
	inviteCodeLength   = 6
	inviteRoomLifetime = 24 * time.Hour
	maxPlayersPerRoom  = 2
)

var (
	errInviteRoomFull    = connect.NewError(connect.CodeFailedPrecondition, errors.New("room is full"))
	errInviteCodeInvalid = connect.NewError(connect.CodeNotFound, errors.New("invalid or expired invite code"))
	errInvalidGameKey    = connect.NewError(connect.CodeInvalidArgument, errors.New("invalid game key"))
)

var allowedInviteGameKeys = map[string]struct{}{
	"tic-tac-toe": {},
	"connect-4":   {},
}

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
		if _, ok := allowedInviteGameKeys[gameKey]; !ok {
			return nil, errInvalidGameKey
		}
	}

	if err := app.clearPersonFromInviteRoom(ctx, person.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	expiresAt := pgtype.Timestamp{Time: time.Now().Add(inviteRoomLifetime), Valid: true}
	var room database.InviteRoom
	var err error
	for range 8 {
		code := random.GenerateInviteCode(inviteCodeLength)
		room, err = app.Queries.InsertInviteRoom(ctx, database.InsertInviteRoomParams{
			InviteCode:   code,
			GameKey:      gameKey,
			HostPersonID: person.ID,
			ExpiresAt:    expiresAt,
		})
		if err == nil {
			break
		}
		if !isUniqueViolation(err) {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("could not allocate invite code"))
	}

	if err := app.Queries.InsertInviteRoomPlayer(ctx, database.InsertInviteRoomPlayerParams{
		RoomID:   room.ID,
		PersonID: person.ID,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&invitev1.CreateInviteRoomResponse{
		InviteCode: room.InviteCode,
		RoomId:     room.ID,
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

	// Check if user is already in a room
	if existingRoom, err := app.Queries.GetInviteRoomForPerson(ctx, person.ID); err == nil {
		// If already in the requested room, treat as idempotent join
		if existingRoom.InviteCode == code {
			return connect.NewResponse(&invitev1.JoinInviteRoomResponse{RoomId: existingRoom.ID}), nil
		}
		// Otherwise, leave current room (swap rooms)
		if err := app.clearPersonFromInviteRoom(ctx, person.ID); err != nil {
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else if !errors.Is(err, pgx.ErrNoRows) {
		// Real error (not just "not in any room")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Lookup and validate the room
	room, err := app.Queries.GetInviteRoomByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errInviteCodeInvalid
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	// Prevent overfilling the room
	if count, err := app.Queries.CountInviteRoomPlayers(ctx, room.ID); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	} else if count >= maxPlayersPerRoom {
		return nil, errInviteRoomFull
	}

	// Add the user to the room (idempotent in DB via ON CONFLICT DO NOTHING)
	if err := app.Queries.InsertInviteRoomPlayer(ctx, database.InsertInviteRoomPlayerParams{
		RoomID:   room.ID,
		PersonID: person.ID,
	}); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&invitev1.JoinInviteRoomResponse{RoomId: room.ID}), nil
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

	room, err := app.Queries.GetInviteRoomByCode(ctx, code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errInviteCodeInvalid
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	memberRoom, err := app.Queries.GetInviteRoomForPerson(ctx, person.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not in this room"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	if memberRoom.ID != room.ID {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("not in this room"))
	}

	players, err := app.Queries.ListInviteRoomPlayers(ctx, room.ID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	accounts := make([]*accountv1.Account, len(players))
	for i, p := range players {
		accounts[i] = personToAccount(p)
	}

	return connect.NewResponse(&invitev1.GetInviteRoomResponse{
		InviteCode: room.InviteCode,
		GameKey:    room.GameKey,
		RoomId:     room.ID,
		Players:    accounts,
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

	if err := app.clearPersonFromInviteRoom(ctx, person.ID); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return connect.NewResponse(&invitev1.LeaveInviteRoomResponse{Ok: true}), nil
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&invitev1.LeaveInviteRoomResponse{Ok: true}), nil
}

// clearPersonFromInviteRoom removes the person from their current invite room and deletes the room if empty or they were the host.
func (app *APIApplication) clearPersonFromInviteRoom(ctx context.Context, personID int64) error {
	room, err := app.Queries.GetInviteRoomForPerson(ctx, personID)
	if err != nil {
		return err
	}

	if err := app.Queries.DeleteInviteRoomPlayer(ctx, database.DeleteInviteRoomPlayerParams{
		RoomID:   room.ID,
		PersonID: personID,
	}); err != nil {
		return err
	}

	count, err := app.Queries.CountInviteRoomPlayers(ctx, room.ID)
	if err != nil {
		return err
	}
	if count == 0 || room.HostPersonID == personID {
		return app.Queries.DeleteInviteRoom(ctx, room.ID)
	}
	return nil
}

func normalizeInviteCode(code string) string {
	return strings.ToUpper(strings.TrimSpace(code))
}

func isUniqueViolation(err error) bool {
	return err != nil && strings.Contains(err.Error(), "duplicate key")
}
