package api

import (
	"context"
	"log/slog"
	"net/http"
	"strings"

	"github.com/arian-nj/chigame/backend/database"
)

func (app *ApiApplication) AuthenticateHeader(ctx context.Context, header http.Header) *database.Person {
	header.Add("Vary", "Authorization")

	authorizationHeader := header.Get("Authorization")
	if authorizationHeader == "" {
		return nil
	}

	headerParts := strings.Split(authorizationHeader, " ")

	if (len(headerParts) == 2 && headerParts[0] == "Bearer") == false {
		return nil
	}
	token := headerParts[1]
	userID, err := app.ValidateToken(token)
	if err != nil {
		slog.Error("authorize header ", "error", err)
		return nil
	}

	person, err := app.Queries.GetPersonByID(ctx, int32(userID))
	if err != nil {
		slog.Error("no user found auth header middleware", "err", err)
		return nil
	}
	return &person
}
