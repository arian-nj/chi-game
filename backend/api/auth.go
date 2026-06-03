package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/internals/random"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5"
)

type AuthType string

const (
	AuthTypeGuestDevice = "guest_device"
)

const GuestJWTExpiryDuration = 365 * 24 * time.Hour

type CustomClaim struct {
	IsGuest  bool   `json:"igu"`
	DeviceID string `json:"did,omitempty"`
	jwt.RegisteredClaims
}

func (app *APIApplication) GetOrCreateGuestUser(ctx context.Context, deviceID string) (*database.Person, error) {
	person, err := app.Queries.GetPersonByAuthMethod(ctx, database.GetPersonByAuthMethodParams{
		AuthType:  AuthTypeGuestDevice,
		AuthValue: deviceID,
	})
	if err == nil {
		return &person, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return nil, fmt.Errorf("query guest user: %w", err)
	}

	randString := random.GenerateRandomUsername(8)
	newPerson, err := app.Queries.InsertPerson(ctx, database.InsertPersonParams{
		IsGuest:     true,
		Username:    randString,
		DisplayName: randString,
	})
	if err != nil {
		return nil, fmt.Errorf("create person: %w", err)
	}

	_, err = app.Queries.InsertAuthMethod(ctx, database.InsertAuthMethodParams{
		UserID:    int64(newPerson.ID),
		AuthType:  AuthTypeGuestDevice,
		AuthValue: deviceID,
	})
	if err != nil {
		return nil, fmt.Errorf("link auth method: %w", err)
	}

	return &newPerson, nil
}

func createGuestToken(userID int64, deviceID string, secret []byte) (string, error) {
	claims := CustomClaim{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(GuestJWTExpiryDuration)),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		IsGuest:  true,
		DeviceID: deviceID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func (app *APIApplication) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaim{}, func(token *jwt.Token) (any, error) {
		return app.Config.Jwt.SecretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		slog.Error("error parsing token", "err", err)
		return 0, err
	}

	claims, ok := token.Claims.(*CustomClaim)
	if !ok || !token.Valid {
		return 0, fmt.Errorf("invalid token claims")
	}

	expireAt, err := claims.GetExpirationTime()
	if err != nil {
		return 0, err
	}
	if expireAt.Time.Before(time.Now()) {
		return 0, fmt.Errorf("token expired")
	}

	notBefore, err := claims.GetNotBefore()
	if err != nil {
		return 0, err
	}
	if notBefore.Time.After(time.Now()) {
		return 0, fmt.Errorf("token not yet valid")
	}

	sub, err := claims.GetSubject()
	if err != nil {
		return 0, err
	}

	userID, err := strconv.Atoi(sub)
	if err != nil {
		return 0, err
	}

	return userID, nil
}
