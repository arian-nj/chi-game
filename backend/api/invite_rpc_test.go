package api

import (
	"context"
	"fmt"
	"testing"
	"time"

	"connectrpc.com/connect"
	invitev1 "github.com/arian-nj/chigame/backend/gen/invite/v1"
)

func authRequest[T any](t *testing.T, msg *T, token string) *connect.Request[T] {
	t.Helper()
	req := connect.NewRequest(msg)
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return req
}

func TestCreateInviteRoomWhenNotInRoom(t *testing.T) {
	app := setupTestApp(t)
	guest := mustValidateGuest(t, app, "invite-create-device")

	resp, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, guest.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	if resp.Msg.GetInviteCode() == "" {
		t.Fatal("expected invite code")
	}
}

func TestCreateMultipleInviteRoomsSameHost(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-multi-host")

	resp1, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom first: %v", err)
	}
	resp2, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom second: %v", err)
	}
	if resp1.Msg.GetInviteCode() == resp2.Msg.GetInviteCode() {
		t.Fatal("expected different invite codes")
	}

	_, err = app.GetInviteRoom(context.Background(), authRequest(t, &invitev1.GetInviteRoomRequest{InviteCode: resp1.Msg.GetInviteCode()}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetInviteRoom first: %v", err)
	}
	_, err = app.GetInviteRoom(context.Background(), authRequest(t, &invitev1.GetInviteRoomRequest{InviteCode: resp2.Msg.GetInviteCode()}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetInviteRoom second: %v", err)
	}
}

func TestCreateJoinGetInviteRoom(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-host-device")
	guest := mustValidateGuest(t, app, "invite-guest-device")

	createResp, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	code := createResp.Msg.GetInviteCode()

	_, err = app.JoinInviteRoom(context.Background(), authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code}, guest.GetToken()))
	if err != nil {
		t.Fatalf("JoinInviteRoom: %v", err)
	}

	getResp, err := app.GetInviteRoom(context.Background(), authRequest(t, &invitev1.GetInviteRoomRequest{InviteCode: code}, host.GetToken()))
	if err != nil {
		t.Fatalf("GetInviteRoom: %v", err)
	}
	if len(getResp.Msg.GetPlayers()) != 2 {
		t.Fatalf("expected 2 players, got %d", len(getResp.Msg.GetPlayers()))
	}
}

func TestJoinInviteRoomFull(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-full-host")
	guest1 := mustValidateGuest(t, app, "invite-full-guest1")
	guest2 := mustValidateGuest(t, app, "invite-full-guest2")

	createResp, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	code := createResp.Msg.GetInviteCode()

	_, err = app.JoinInviteRoom(context.Background(), authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code}, guest1.GetToken()))
	if err != nil {
		t.Fatalf("JoinInviteRoom guest1: %v", err)
	}

	_, err = app.JoinInviteRoom(context.Background(), authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code}, guest2.GetToken()))
	assertConnectCode(t, err, connect.CodeFailedPrecondition)
}

func TestJoinInviteRoomExpired(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-expired-host")
	guest := mustValidateGuest(t, app, "invite-expired-guest")

	createResp, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	code := createResp.Msg.GetInviteCode()

	room, ok := app.InviteStore.GetByCode(code)
	if !ok {
		t.Fatal("expected room in store")
	}
	room.ExpiresAt = time.Now().Add(-time.Minute)

	_, err = app.JoinInviteRoom(context.Background(), authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code}, guest.GetToken()))
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestInviteHostLeaveDeletesRoom(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-leave-host")
	guest := mustValidateGuest(t, app, "invite-leave-guest")

	createResp, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	code := createResp.Msg.GetInviteCode()

	_, err = app.JoinInviteRoom(context.Background(), authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code}, guest.GetToken()))
	if err != nil {
		t.Fatalf("JoinInviteRoom: %v", err)
	}

	_, err = app.LeaveInviteRoom(context.Background(), authRequest(t, &invitev1.LeaveInviteRoomRequest{InviteCode: code}, host.GetToken()))
	if err != nil {
		t.Fatalf("LeaveInviteRoom: %v", err)
	}

	_, err = app.GetInviteRoom(context.Background(), authRequest(t, &invitev1.GetInviteRoomRequest{InviteCode: code}, guest.GetToken()))
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestJoinInviteRoomIdempotent(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-idem-host")
	guest := mustValidateGuest(t, app, "invite-idem-guest")

	createResp, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	code := createResp.Msg.GetInviteCode()

	joinReq := authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code}, guest.GetToken())
	_, err = app.JoinInviteRoom(context.Background(), joinReq)
	if err != nil {
		t.Fatalf("JoinInviteRoom first: %v", err)
	}

	_, err = app.JoinInviteRoom(context.Background(), joinReq)
	if err != nil {
		t.Fatalf("JoinInviteRoom idempotent: %v", err)
	}
}

func TestGuestCanStayInRoomAfterHostCreatesAnother(t *testing.T) {
	app := setupTestApp(t)
	host := mustValidateGuest(t, app, "invite-stay-host")
	guest := mustValidateGuest(t, app, "invite-stay-guest")

	first, err := app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom first: %v", err)
	}
	code1 := first.Msg.GetInviteCode()

	_, err = app.JoinInviteRoom(context.Background(), authRequest(t, &invitev1.JoinInviteRoomRequest{InviteCode: code1}, guest.GetToken()))
	if err != nil {
		t.Fatalf("JoinInviteRoom: %v", err)
	}

	_, err = app.CreateInviteRoom(context.Background(), authRequest(t, &invitev1.CreateInviteRoomRequest{GameKey: ""}, host.GetToken()))
	if err != nil {
		t.Fatalf("CreateInviteRoom second: %v", err)
	}

	getResp, err := app.GetInviteRoom(context.Background(), authRequest(t, &invitev1.GetInviteRoomRequest{InviteCode: code1}, guest.GetToken()))
	if err != nil {
		t.Fatalf("GetInviteRoom: %v", err)
	}
	if len(getResp.Msg.GetPlayers()) != 2 {
		t.Fatalf("expected guest to remain in first room, got %d players", len(getResp.Msg.GetPlayers()))
	}
}
