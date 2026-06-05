package api

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	"github.com/jackc/pgx/v5"
)

func (app *APIApplication) GetMe(
	ctx context.Context,
	req *connect.Request[accountv1.GetMeRequest],
) (*connect.Response[accountv1.GetMeResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	return connect.NewResponse(&accountv1.GetMeResponse{
		Account: personToAccount(person),
	}), nil
}

func (app *APIApplication) GetPerson(
	ctx context.Context,
	req *connect.Request[accountv1.GetPersonRequest],
) (*connect.Response[accountv1.GetPersonResponse], error) {
	if app.AuthenticateHeader(ctx, req.Header()) == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	person, err := app.Queries.GetPersonByID(ctx, req.Msg.GetId())
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, connect.NewError(connect.CodeNotFound, errors.New("person not found"))
		}
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&accountv1.GetPersonResponse{
		Account: personToAccount(&person),
	}), nil
}

func personToAccount(person *database.Person) *accountv1.Account {
	return &accountv1.Account{
		Id:          person.ID,
		Username:    person.Username,
		DisplayName: person.DisplayName,
	}
}
