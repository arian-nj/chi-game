package api

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
)

// GetMe implements [accountv1connect.AccountServiceHandler].
func (app *APIApplication) GetMe(ctx context.Context, req *connect.Request[accountv1.GetMeRequest]) (*connect.Response[accountv1.GetMeResponse], error) {
	personRow := app.AuthenticateHeader(ctx, req.Header())
	if personRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}
	return connect.NewResponse(&accountv1.GetMeResponse{
		Account: &accountv1.Account{
			Id:          int64(personRow.ID),
			Username:    personRow.Username,
			DisplayName: personRow.DisplayName,
		},
	}), nil
}

func (app *APIApplication) GetPerson(ctx context.Context, req *connect.Request[accountv1.GetPersonRequest]) (*connect.Response[accountv1.GetPersonResponse], error) {
	personRow := app.AuthenticateHeader(ctx, req.Header())
	if personRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	person, err := app.Queries.GetPersonByID(ctx, int(req.Msg.Id))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("can't get person"))
	}

	return connect.NewResponse(&accountv1.GetPersonResponse{
		Account: &accountv1.Account{
			Id:          int64(person.ID),
			Username:    person.Username,
			DisplayName: person.DisplayName,
		},
	}), nil
}
