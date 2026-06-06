package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
)

func authRequest[T any](t *testing.T, msg *T, token string) *connect.Request[T] {
	t.Helper()
	req := connect.NewRequest(msg)
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req
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

func TestCreateJoinGetRoom(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-host-device")
	guest := mustValidateGuest(t, app, "room-guest-device")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	_, err = app.JoinRoom(context.Background(), authRequest(t, &roomv1.JoinRoomRequest{Code: code}, guest.GetToken()))
	if err != nil {
		t.Fatalf("JoinRoom: %v", err)
	}

	getResp, err := app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: code}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetRoom: %v", err)
	}
	if len(getResp.Msg.GetPlayers()) != 2 {
		t.Fatalf("expected 2 players, got %d", len(getResp.Msg.GetPlayers()))
	}
}

func TestJoinRoomFull(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-full-host")
	guest1 := mustValidateGuest(t, app, "room-full-guest1")
	guest2 := mustValidateGuest(t, app, "room-full-guest2")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	_, err = app.JoinRoom(context.Background(), authRequest(t, &roomv1.JoinRoomRequest{Code: code}, guest1.GetToken()))
	if err != nil {
		t.Fatalf("JoinRoom guest1: %v", err)
	}

	_, err = app.JoinRoom(context.Background(), authRequest(t, &roomv1.JoinRoomRequest{Code: code}, guest2.GetToken()))
	assertConnectCode(t, err, connect.CodeFailedPrecondition)
}

func TestJoinRoomExpired(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-expired-host")
	guest := mustValidateGuest(t, app, "room-expired-guest")

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

	_, err = app.JoinRoom(context.Background(), authRequest(t, &roomv1.JoinRoomRequest{Code: code}, guest.GetToken()))
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestHostLeaveDeletesRoom(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-leave-host")
	guest := mustValidateGuest(t, app, "room-leave-guest")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	_, err = app.JoinRoom(context.Background(), authRequest(t, &roomv1.JoinRoomRequest{Code: code}, guest.GetToken()))
	if err != nil {
		t.Fatalf("JoinRoom: %v", err)
	}

	_, err = app.LeaveRoom(context.Background(), authRequest(t, &roomv1.LeaveRoomRequest{Code: code}, host.GetToken()))
	if err != nil {
		t.Fatalf("LeaveRoom: %v", err)
	}

	_, err = app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: code}, guest.GetToken()))
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestJoinRoomIdempotent(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-idem-host")
	guest := mustValidateGuest(t, app, "room-idem-guest")

	createResp, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom: %v", err)
	}
	code := createResp.Msg.GetCode()

	joinReq := authRequest(t, &roomv1.JoinRoomRequest{Code: code}, guest.GetToken())
	_, err = app.JoinRoom(context.Background(), joinReq)
	if err != nil {
		t.Fatalf("JoinRoom first: %v", err)
	}

	_, err = app.JoinRoom(context.Background(), joinReq)
	if err != nil {
		t.Fatalf("JoinRoom idempotent: %v", err)
	}
}

func TestGuestCanStayInRoomAfterHostCreatesAnother(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "room-stay-host")
	guest := mustValidateGuest(t, app, "room-stay-guest")

	first, err := app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom first: %v", err)
	}
	code1 := first.Msg.GetCode()

	_, err = app.JoinRoom(context.Background(), authRequest(t, &roomv1.JoinRoomRequest{Code: code1}, guest.GetToken()))
	if err != nil {
		t.Fatalf("JoinRoom: %v", err)
	}

	_, err = app.CreateRoom(context.Background(), authRequest(t, &roomv1.CreateRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateRoom second: %v", err)
	}

	getResp, err := app.GetRoom(context.Background(), authRequest(t, &roomv1.GetRoomRequest{Code: code1}, guest.GetToken()))
	if err != nil {
		t.Fatalf("GetRoom: %v", err)
	}
	if len(getResp.Msg.GetPlayers()) != 2 {
		t.Fatalf("expected guest to remain in first room, got %d players", len(getResp.Msg.GetPlayers()))
	}
}
