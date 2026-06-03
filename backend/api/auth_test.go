package api

import (
	"context"
	"net/http"
	"testing"

	"connectrpc.com/connect"
	authv1 "github.com/arian-nj/chigame/backend/gen/auth/v1"
	"github.com/arian-nj/chigame/backend/database"
)

func TestCreateGuestToken_roundTrip(t *testing.T) {
	app := &APIApplication{Config: testJWTConfig()}
	const (
		userID   int64 = 42
		deviceID       = "device-abc"
	)

	token, err := createGuestToken(userID, deviceID, app.Config.Jwt.SecretKey)
	if err != nil {
		t.Fatalf("createGuestToken: %v", err)
	}

	gotID, err := app.ValidateToken(token)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if gotID != int(userID) {
		t.Fatalf("userID = %d, want %d", gotID, userID)
	}

	claims := parseGuestClaims(t, app.Config.Jwt.SecretKey, token)
	if !claims.IsGuest {
		t.Fatal("expected IsGuest claim")
	}
	if claims.DeviceID != deviceID {
		t.Fatalf("DeviceID = %q, want %q", claims.DeviceID, deviceID)
	}
}

func TestValidateToken_rejectsInvalidTokens(t *testing.T) {
	app := &APIApplication{Config: testJWTConfig()}

	tests := []struct {
		name  string
		token string
	}{
		{name: "garbage", token: "not-a-jwt"},
		{name: "empty", token: ""},
		{name: "wrong secret", token: mustSignWithSecret(t, 1, "dev", []byte("other-secret-key-for-tests"))},
		{name: "expired", token: signedExpiredGuestToken(t, 1, "dev", app.Config.Jwt.SecretKey)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := app.ValidateToken(tc.token); err == nil {
				t.Fatal("expected error")
			}
		})
	}
}

func mustSignWithSecret(t *testing.T, userID int64, deviceID string, secret []byte) string {
	t.Helper()
	signed, err := createGuestToken(userID, deviceID, secret)
	if err != nil {
		t.Fatalf("createGuestToken: %v", err)
	}
	return signed
}

func TestGenerateSecureDeviceID_format(t *testing.T) {
	id := generateSecureDeviceID()
	if len(id) != 64 {
		t.Fatalf("length = %d, want 64 hex chars", len(id))
	}
	for _, c := range id {
		if (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			t.Fatalf("non-hex character %q in device ID", c)
		}
	}
}

func TestGetOrCreateGuestUser_createsAndReuses(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()
	deviceID := "test-device-001"

	first, err := app.GetOrCreateGuestUser(ctx, deviceID)
	if err != nil {
		t.Fatalf("first create: %v", err)
	}
	if !first.IsGuest {
		t.Fatal("expected guest user")
	}
	if first.Username == "" || first.Username != first.DisplayName {
		t.Fatalf("unexpected username/display_name: %q / %q", first.Username, first.DisplayName)
	}

	stored, err := app.Queries.GetPersonByAuthMethod(ctx, database.GetPersonByAuthMethodParams{
		AuthType:  AuthTypeGuestDevice,
		AuthValue: deviceID,
	})
	if err != nil {
		t.Fatalf("GetPersonByAuthMethod: %v", err)
	}
	if stored.ID != first.ID {
		t.Fatalf("auth method user id = %d, want %d", stored.ID, first.ID)
	}

	second, err := app.GetOrCreateGuestUser(ctx, deviceID)
	if err != nil {
		t.Fatalf("second lookup: %v", err)
	}
	if first.ID != second.ID {
		t.Fatalf("user IDs differ: %d vs %d", first.ID, second.ID)
	}
	if first.Username != second.Username {
		t.Fatalf("username changed: %q -> %q", first.Username, second.Username)
	}
}

func TestValidateGuest_createsGuestWithToken(t *testing.T) {
	app := setupTestApp(t)
	ctx := context.Background()

	resp, err := app.ValidateGuest(ctx, connect.NewRequest(&authv1.ValidateGuestRequest{}))
	if err != nil {
		t.Fatalf("ValidateGuest: %v", err)
	}

	person, err := app.Queries.GetPersonByAuthMethod(ctx, database.GetPersonByAuthMethodParams{
		AuthType:  AuthTypeGuestDevice,
		AuthValue: resp.Msg.GetDeviceId(),
	})
	if err != nil {
		t.Fatalf("GetPersonByAuthMethod: %v", err)
	}

	userID, err := app.ValidateToken(resp.Msg.GetToken())
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if int64(userID) != person.ID {
		t.Fatalf("token user id = %d, db person id = %d", userID, person.ID)
	}

	claims := parseGuestClaims(t, app.Config.Jwt.SecretKey, resp.Msg.GetToken())
	if claims.DeviceID != resp.Msg.GetDeviceId() {
		t.Fatalf("claim device id = %q, want %q", claims.DeviceID, resp.Msg.GetDeviceId())
	}
}

func TestValidateGuest_returnsSameUserForDevice(t *testing.T) {
	app := setupTestApp(t)

	first := mustValidateGuest(t, app, "")
	second := mustValidateGuest(t, app, first.GetDeviceId())

	userID1, err := app.ValidateToken(first.GetToken())
	if err != nil {
		t.Fatalf("token 1: %v", err)
	}
	userID2, err := app.ValidateToken(second.GetToken())
	if err != nil {
		t.Fatalf("token 2: %v", err)
	}
	if userID1 != userID2 {
		t.Fatalf("user IDs differ: %d vs %d", userID1, userID2)
	}
	if first.GetDeviceId() != second.GetDeviceId() {
		t.Fatal("device ID should stay the same")
	}
}

func TestAuthenticateHeader_rejectsBadAuth(t *testing.T) {
	app := setupTestApp(t)
	guest := mustValidateGuest(t, app, "")

	tests := []struct {
		name       string
		authHeader string
	}{
		{name: "missing"},
		{name: "invalid jwt", authHeader: "Bearer not-a-jwt"},
		{name: "wrong scheme", authHeader: "Basic " + guest.GetToken()},
		{name: "bearer only", authHeader: "Bearer"},
		{name: "wrong secret", authHeader: "Bearer " + mustSignWithSecret(t, 1, guest.GetDeviceId(), []byte("wrong-secret"))},
		{name: "expired", authHeader: "Bearer " + signedExpiredGuestToken(t, 1, guest.GetDeviceId(), app.Config.Jwt.SecretKey)},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			header := http.Header{}
			if tc.authHeader != "" {
				header.Set("Authorization", tc.authHeader)
			}
			if got := app.AuthenticateHeader(context.Background(), header); got != nil {
				t.Fatal("expected nil person")
			}
		})
	}
}
