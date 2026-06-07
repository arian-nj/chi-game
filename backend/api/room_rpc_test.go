package api

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	authv1 "github.com/arian-nj/chigame/backend/gen/auth/v1"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
)

func authRequest[T any](t *testing.T, msg *T, token string) *connect.Request[T] {
	t.Helper()
	req := connect.NewRequest(msg)
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req
}

func guestPerson(t *testing.T, app *APIApplication, guest *authv1.ValidateGuestResponse) *database.Person {
	t.Helper()
	claims := parseGuestClaims(t, []byte(testJWTSecret), guest.GetToken())
	personID, err := strconv.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		t.Fatalf("parse person id: %v", err)
	}
	person, err := app.Queries.GetPersonByID(context.Background(), personID)
	if err != nil {
		t.Fatalf("GetPersonByID: %v", err)
	}
	return &person
}

func TestCreateRoomWhenNotInRoom(t *testing.T) {
	app := setupTestApp(t)
	guest := mustValidateGuest(t, app, "room-create-device")

	resp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, guest.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	if resp.Msg.GetCode() == "" {
		t.Fatal("expected room code")
	}
}

func TestCreateMultipleRoomsSameHost(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-multi-host")

	resp1, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom first: %v", err)
	}
	resp2, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom second: %v", err)
	}
	if resp1.Msg.GetCode() == resp2.Msg.GetCode() {
		t.Fatal("expected different room codes")
	}

	_, err = app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: resp1.Msg.GetCode()}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetRoom first: %v", err)
	}
	_, err = app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: resp2.Msg.GetCode()}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetRoom second: %v", err)
	}
}

func TestGetRoomExpired(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-expired-host")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	room, ok := app.RoomsStore.GetByCode(code)
	if !ok {
		t.Fatal("expected room in store")
	}
	room.ExpiresAt = time.Now().Add(-time.Minute)

	_, err = app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: code}, host.GetToken()))
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestAddMemberFull(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-full-host")
	guest1 := mustValidateGuest(t, app, "room-full-guest1")
	guest2 := mustValidateGuest(t, app, "room-full-guest2")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	room, ok := app.RoomsStore.GetByCode(code)
	if !ok {
		t.Fatal("expected room in store")
	}

	if err := room.AddMember(NewRoomMember(guestPerson(t, app, host), nil)); err != nil {
		t.Fatalf("AddMember host: %v", err)
	}
	if err := room.AddMember(NewRoomMember(guestPerson(t, app, guest1), nil)); err != nil {
		t.Fatalf("AddMember guest1: %v", err)
	}
	err = room.AddMember(NewRoomMember(guestPerson(t, app, guest2), nil))
	if err != errRoomFull {
		t.Fatalf("expected room full, got: %v", err)
	}
}

func TestGetRoomReturnsConnectedPlayers(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-get-players-host")
	guest := mustValidateGuest(t, app, "room-get-players-guest")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	room, ok := app.RoomsStore.GetByCode(code)
	if !ok {
		t.Fatal("expected room in store")
	}

	hostPerson := guestPerson(t, app, host)
	guestPersonRow := guestPerson(t, app, guest)
	if err := room.AddMember(NewRoomMember(hostPerson, nil)); err != nil {
		t.Fatalf("AddMember host: %v", err)
	}
	if err := room.AddMember(NewRoomMember(guestPersonRow, nil)); err != nil {
		t.Fatalf("AddMember guest: %v", err)
	}

	resp, err := app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: code}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetRoom: %v", err)
	}
	if len(resp.Msg.GetPlayers()) != 2 {
		t.Fatalf("expected 2 players, got %d", len(resp.Msg.GetPlayers()))
	}
}

func TestAddMemberIdempotent(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-idem-host")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	room, ok := app.RoomsStore.GetByCode(code)
	if !ok {
		t.Fatal("expected room in store")
	}

	member := NewRoomMember(guestPerson(t, app, host), nil)
	if err := room.AddMember(member); err != nil {
		t.Fatalf("AddMember first: %v", err)
	}
	if err := room.AddMember(member); err != nil {
		t.Fatalf("AddMember idempotent: %v", err)
	}
}
