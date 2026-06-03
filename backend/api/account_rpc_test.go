package api

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
)

func TestGetMe_returnsAuthenticatedGuest(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()
	guest := mustValidateGuest(t, app, "")

	userID, err := app.ValidateToken(guest.GetToken())
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}

	req := connect.NewRequest(&accountv1.GetMeRequest{})
	req.Header().Set("Authorization", "Bearer "+guest.GetToken())

	resp, err := app.GetMe(ctx, req)
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}

	account := resp.Msg.GetAccount()
	if account.GetUsername() == "" {
		t.Fatal("expected username")
	}
	if account.GetId() != int64(userID) {
		t.Fatalf("account id = %d, token user id = %d", account.GetId(), userID)
	}
	if account.GetDisplayName() != account.GetUsername() {
		t.Fatalf("display_name = %q, username = %q", account.GetDisplayName(), account.GetUsername())
	}
}

func TestGetMe_unauthenticated(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		authHeader string
		want       connect.Code
	}{
		{name: "missing", want: connect.CodeUnauthenticated},
		{name: "invalid jwt", authHeader: "Bearer not-a-jwt", want: connect.CodeUnauthenticated},
		{name: "malformed", authHeader: "Bearer", want: connect.CodeUnauthenticated},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := connect.NewRequest(&accountv1.GetMeRequest{})
			if tc.authHeader != "" {
				req.Header().Set("Authorization", tc.authHeader)
			}
			_, err := app.GetMe(ctx, req)
			assertConnectCode(t, err, tc.want)
		})
	}
}

func TestGetPerson_returnsGuestByID(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()
	guest := mustValidateGuest(t, app, "")

	meReq := connect.NewRequest(&accountv1.GetMeRequest{})
	meReq.Header().Set("Authorization", "Bearer "+guest.GetToken())
	me, err := app.GetMe(ctx, meReq)
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}

	personReq := connect.NewRequest(&accountv1.GetPersonRequest{Id: me.Msg.GetAccount().GetId()})
	personReq.Header().Set("Authorization", "Bearer "+guest.GetToken())

	resp, err := app.GetPerson(ctx, personReq)
	if err != nil {
		t.Fatalf("GetPerson: %v", err)
	}
	if resp.Msg.GetAccount().GetId() != me.Msg.GetAccount().GetId() {
		t.Fatalf("GetPerson id = %d, want %d", resp.Msg.GetAccount().GetId(), me.Msg.GetAccount().GetId())
	}
	if resp.Msg.GetAccount().GetUsername() != me.Msg.GetAccount().GetUsername() {
		t.Fatal("username mismatch between GetMe and GetPerson")
	}
}

func TestGetPerson_notFound(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()
	guest := mustValidateGuest(t, app, "")

	req := connect.NewRequest(&accountv1.GetPersonRequest{Id: 999_999})
	req.Header().Set("Authorization", "Bearer "+guest.GetToken())

	_, err := app.GetPerson(ctx, req)
	assertConnectCode(t, err, connect.CodeNotFound)
}

func TestGetPerson_unauthenticated(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()

	_, err := app.GetPerson(ctx, connect.NewRequest(&accountv1.GetPersonRequest{Id: 1}))
	assertConnectCode(t, err, connect.CodeUnauthenticated)
}

func TestAuthenticateHeader_acceptsValidGuest(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()
	guest := mustValidateGuest(t, app, "")

	header := http.Header{}
	header.Set("Authorization", "Bearer "+guest.GetToken())

	person := app.AuthenticateHeader(ctx, header)
	if person == nil {
		t.Fatal("expected authenticated person")
	}

	userID, err := app.ValidateToken(guest.GetToken())
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if person.ID != int64(userID) {
		t.Fatalf("person id = %d, token user id = %d", person.ID, userID)
	}
}
