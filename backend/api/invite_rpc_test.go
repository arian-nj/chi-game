package api

import (
	"context"
	"fmt"
	"testing"

	"connectrpc.com/connect"
	invitev1 "github.com/arian-nj/chigame/backend/gen/invite/v1"
)

func TestCreateInviteRoomWhenNotInRoom(t *testing.T) {
	app := setupTestApp(t)
	guest := mustValidateGuest(t, app, "invite-create-device")

	req := connect.NewRequest(&invitev1.CreateInviteRoomRequest{GameKey: ""})
	req.Header().Set("Authorization", fmt.Sprintf("Bearer %s", guest.GetToken()))

	resp, err := app.CreateInviteRoom(context.Background(), req)
	if err != nil {
		t.Fatalf("CreateInviteRoom: %v", err)
	}
	if resp.Msg.GetInviteCode() == "" {
		t.Fatal("expected invite code")
	}
}

func TestClearPersonFromInviteRoomNoMembership(t *testing.T) {
	app := setupTestApp(t)
	guest := mustValidateGuest(t, app, "invite-clear-device")

	userID, err := app.ValidateToken(guest.GetToken())
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}

	if err := app.clearPersonFromInviteRoom(context.Background(), int64(userID)); err != nil {
		t.Fatalf("expected nil when not in room, got %v", err)
	}
}
