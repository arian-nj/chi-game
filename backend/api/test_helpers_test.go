package api

import (
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"connectrpc.com/connect"
	authv1 "github.com/arian-nj/chigame/backend/gen/auth/v1"
	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/db"
	"github.com/arian-nj/chigame/backend/internals/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const testJWTSecret = "test-jwt-secret-for-guest-auth-tests"

func requireDatabaseURL(t *testing.T) string {
	t.Helper()
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("DATABASE_URL not set; run ./mash.sh test")
	}
	return connStr
}

func setupTestApp(t *testing.T) *APIApplication {
	t.Helper()
	connStr := requireDatabaseURL(t)

	if err := db.Migrate(connStr); err != nil {
		t.Fatalf("migrate: %v", err)
	}

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		t.Fatalf("db pool: %v", err)
	}

	t.Cleanup(func() {
		_, _ = pool.Exec(context.Background(), "TRUNCATE person_auth_methods, persons RESTART IDENTITY CASCADE")
		pool.Close()
	})

	cfg := &config.Config{
		DatabaseURL: connStr,
		Jwt:         config.JWTConfig{SecretKey: []byte(testJWTSecret)},
	}

	return NewAPIApplication(cfg, database.New(pool))
}

func testJWTConfig() *config.Config {
	return &config.Config{
		Jwt: config.JWTConfig{SecretKey: []byte(testJWTSecret)},
	}
}

func mustValidateGuest(t *testing.T, app *APIApplication, deviceID string) *authv1.ValidateGuestResponse {
	t.Helper()
	ctx := context.Background()
	resp, err := app.ValidateGuest(ctx, connect.NewRequest(&authv1.ValidateGuestRequest{DeviceId: deviceID}))
	if err != nil {
		t.Fatalf("ValidateGuest: %v", err)
	}
	if resp.Msg.GetDeviceId() == "" {
		t.Fatal("expected device ID")
	}
	if resp.Msg.GetToken() == "" {
		t.Fatal("expected token")
	}
	return resp.Msg
}

func parseGuestClaims(t *testing.T, secret []byte, tokenString string) *CustomClaim {
	t.Helper()
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (any, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		t.Fatalf("parse token: %v", err)
	}
	claims, ok := token.Claims.(*CustomClaim)
	if !ok || !token.Valid {
		t.Fatal("invalid token claims")
	}
	return claims
}

func signedExpiredGuestToken(t *testing.T, userID int64, deviceID string, secret []byte) string {
	t.Helper()
	claims := CustomClaim{
		IsGuest:  true,
		DeviceID: deviceID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			NotBefore: jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(secret)
	if err != nil {
		t.Fatalf("sign expired token: %v", err)
	}
	return signed
}

func assertConnectCode(t *testing.T, err error, want connect.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("expected connect error %v", want)
	}
	if got := connect.CodeOf(err); got != want {
		t.Fatalf("connect code = %v, want %v", got, want)
	}
}
