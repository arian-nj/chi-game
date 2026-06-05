package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/arian-nj/chigame/backend/database"
)

var ErrCantAuthenticateUser = errors.New("can't authenticate user")

func (app *APIApplication) AuthenticateHeader(ctx context.Context, header http.Header) *database.Person {
	authorizationHeader := header.Get("Authorization")
	if authorizationHeader == "" {
		return nil
	}

	headerParts := strings.Split(authorizationHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil
	}

	userID, err := app.ValidateToken(headerParts[1])
	if err != nil {
		slog.Error("authorize header", "error", err)
		return nil
	}

	person, err := app.Queries.GetPersonByID(ctx, int64(userID))
	if err != nil {
		slog.Error("no user found for auth header", "err", err)
		return nil
	}

	return &person
}

func (app *APIApplication) AuthenticateQuery(ctx context.Context, header http.Header) *database.Person {
	token := header.Get("auth_token")
	if token == "" {
		return nil
	}

	userID, err := app.ValidateToken(token)
	if err != nil {
		return nil
	}

	person, err := app.Queries.GetPersonByID(ctx, int64(userID))
	if err != nil {
		slog.Error("no user found for auth header", "err", err)
		return nil
	}

	return &person
}
