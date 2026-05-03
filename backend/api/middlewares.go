package api

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/pkg/response"
	"github.com/arian-nj/chigame/backend/pkg/validator"
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

	person, err := app.Queries.GetPersonByID(ctx, userID)
	if err != nil {
		slog.Error("no user found auth header middleware", "err", err)
		return nil
	}
	return &person
}

func (app *ApiApplication) AuthenticateQuery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("auth_token")
		if token == "" {
			app.invalidAuthenticationCreds(w, r)
			return
		}
		userID, err := app.ValidateToken(token)
		if err != nil {
			app.invalidAuthenticationCreds(w, r)
			return
		}

		reqUser := &ReqContextUser{UserID: userID}
		newRequest := ContextSetAuthenticatedUser(r, reqUser)

		next.ServeHTTP(w, newRequest)
	})
}

type contextKey string

const (
	authenticatedUserContextKey = contextKey("authenticatedUser")
)

type ReqContextUser struct {
	UserID int
}

func ContextSetAuthenticatedUser(r *http.Request, user *ReqContextUser) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func ContextGetAuthenticatedUser(queries *database.Queries, r *http.Request) (*database.Person, error) {
	val := r.Context().Value(authenticatedUserContextKey)
	reqConUser, ok := val.(*ReqContextUser)
	if !ok || reqConUser == nil {
		return nil, errors.New("authenticated user missing or invalid type in context")
	}

	var user database.Person
	var err error
	user, err = queries.GetPersonByID(r.Context(), reqConUser.UserID)

	if err != nil {
		slog.Error("no user found in context get authticated user", "err", err)
		return nil, err
	}

	return &user, nil
}

func (app *ApiApplication) ReportServerError(r *http.Request, err error) {
	var (
		message = err.Error()
		// method  = r.Method
		// url     = r.URL.String()
		trace = string(debug.Stack())
	)

	// requestAttrs := slog.Group("request", "method", method, "url", url)
	// app.Logger.Error(message, requestAttrs, "trace", trace)
	// log.Println(message, method, url)
	slog.Error("report server error ", "err", message)
	log.Println(trace)
}

func (app *ApiApplication) errorMessage(w http.ResponseWriter, r *http.Request, status int, message string, headers http.Header) {
	message = strings.ToUpper(message[:1]) + message[1:]

	err := response.JSONWithHeaders(w, status, map[string]string{"Error": message}, headers)
	if err != nil {
		app.ReportServerError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (app *ApiApplication) InvalidAuthenticationToken(w http.ResponseWriter, r *http.Request) {
	headers := make(http.Header)
	headers.Set("WWW-Authenticate", "Bearer")

	app.errorMessage(w, r, http.StatusUnauthorized, "Invalid authentication token", headers)
}

func (app *ApiApplication) invalidAuthenticationCreds(w http.ResponseWriter, r *http.Request) {
	app.errorMessage(w, r, http.StatusUnauthorized, "Invalid authentication credentials", nil)
}

func (app *ApiApplication) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.ReportServerError(r, err)

	message := "The server encountered a problem and could not process your request"
	app.errorMessage(w, r, http.StatusInternalServerError, message, nil)
}

func (app *ApiApplication) NotFound(w http.ResponseWriter, r *http.Request) {
	message := "The requested resource could not be found"
	app.errorMessage(w, r, http.StatusNotFound, message, nil)
}

func (app *ApiApplication) MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("The %s method is not supported for this resource", r.Method)
	app.errorMessage(w, r, http.StatusMethodNotAllowed, message, nil)
}

func (app *ApiApplication) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	app.errorMessage(w, r, http.StatusBadRequest, err.Error(), nil)
}

func (app *ApiApplication) FailedValidation(w http.ResponseWriter, r *http.Request, v validator.Validator) {
	err := response.JSON(w, http.StatusUnprocessableEntity, v)
	if err != nil {
		app.ServerError(w, r, err)
	}
}
