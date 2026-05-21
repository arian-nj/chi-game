package api

import (
	"context"
	"errors"
	"log/slog"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	friendsv1 "github.com/arian-nj/chigame/backend/gen/friends/v1"
	"github.com/jackc/pgx/v5/pgtype"
)

// FindUsername implements [friendsv1connect.FriendsServiceHandler].
func (app *APIApplication) SearchForUsername(ctx context.Context, req *connect.Request[friendsv1.SearchForUsernameRequest]) (
	*connect.Response[friendsv1.SearchForUsernameResponse], error) {
	personRow := app.AuthenticateHeader(ctx, req.Header())
	if personRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrCantAuthenticateUser)
	}

	lookFor := req.Msg.LookFor
	personsRow, err := app.Queries.SearchPersonByUsername(context.Background(), pgtype.Text{String: lookFor, Valid: true})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, errors.New("can't search for username"))
	}

	foundUsers := []*accountv1.Account{}

	for _, person := range personsRow {
		foundUsers = append(foundUsers, &accountv1.Account{
			Id:          int64(person.ID),
			Username:    person.Username,
			DisplayName: person.DisplayName,
		})
	}

	return connect.NewResponse(&friendsv1.SearchForUsernameResponse{
		FoundUsers: foundUsers,
	}), nil
}


func (app *APIApplication) GetFriendshipStatus(ctx context.Context, req *connect.Request[friendsv1.GetFriendshipStatusRequest]) (
	*connect.Response[friendsv1.GetFriendshipStatusResponse], error) {
	personRow := app.AuthenticateHeader(ctx,req.Header())
	if personRow == nil {
		return nil ,connect.NewError(connect.CodeUnauthenticated,ErrCantAuthenticateUser)
	}
	lookFor := req.Msg.ToPerson

	friendshipStatus ,err := app.Queries.GetFriendshipStatus(ctx,database.GetFriendshipStatusParams{
		UserID: personRow.ID,
		FriendID: int(lookFor),
	})
	if err != nil {
		slog.Error("can't get friendship status","err",err,"user_id",personRow.ID,"look_for",lookFor)
		return nil,connect.NewError(connect.CodeInternal,errors.New("can't get friendship status"))
	}

	rpcState := friendsv1.FriendshipStatus_FRIENDSHIP_STATUS_UNSPECIFIED
	switch friendshipStatus {
	case "friends":
		rpcState = friendsv1.FriendshipStatus_FRIENDSHIP_STATUS_FRIEND
	case "request_sent":
		rpcState = friendsv1.FriendshipStatus_FRIENDSHIP_STATUS_REQUESTED
	case "request_received":
		rpcState = friendsv1.FriendshipStatus_FRIENDSHIP_STATUS_RECEIVED_REQUEST
	case "not_connected":
		rpcState = friendsv1.FriendshipStatus_FRIENDSHIP_STATUS_NOTHING
}

	return connect.NewResponse(&friendsv1.GetFriendshipStatusResponse{
		Fstatus: rpcState,
	}),nil
}


func (app *APIApplication) SendFriendReq(ctx context.Context, req *connect.Request[friendsv1.SendFriendReqRequest]) (
	*connect.Response[friendsv1.SendFriendReqResponse], error) {
	personRow := app.AuthenticateHeader(ctx,req.Header())
	if personRow == nil {
		return nil ,connect.NewError(connect.CodeUnauthenticated,ErrCantAuthenticateUser)
	}

	err := app.Queries.InsertFriendRequest(ctx,database.InsertFriendRequestParams{
		SenderID: personRow.ID,
		ReceiverID: int(req.Msg.ToPerson),
	})
	if err != nil {
		slog.Error("can't insert friend status","err",err,"sender_id",personRow.ID,"reciever_id",req.Msg.ToPerson)
		return nil,connect.NewError(connect.CodeInternal,errors.New("can't request friendship"))
	}

	return connect.NewResponse(&friendsv1.SendFriendReqResponse{
		Success: true,
	}),nil
}


func (app *APIApplication) CancelFriendReq(ctx context.Context, req *connect.Request[friendsv1.CancelFriendReqRequest]) (
	*connect.Response[friendsv1.CancelFriendReqResponse], error) {
	personRow := app.AuthenticateHeader(ctx,req.Header())
	if personRow == nil {
		return nil ,connect.NewError(connect.CodeUnauthenticated,ErrCantAuthenticateUser)
	}
	
	receiverId := req.Msg.ToPerson
	_ ,err := app.Queries.CancelFriendRequest(ctx,database.CancelFriendRequestParams{
		SenderID: personRow.ID,
		ReceiverID: int(receiverId),
	})
	if err != nil {
		slog.Error("can't cancel friend request","err",err,"sender_id",personRow.ID,"reciever_id",req.Msg.ToPerson)
		return nil,connect.NewError(connect.CodeInternal,errors.New("can't get friendship status"))
	}

	return connect.NewResponse(&friendsv1.CancelFriendReqResponse{
		Success: true,
	}),nil
}


func (app *APIApplication) AcceptFriendReq(ctx context.Context, req *connect.Request[friendsv1.AcceptFriendReqRequest]) (
	*connect.Response[friendsv1.AcceptFriendReqResponse], error) {

	recPersonRow := app.AuthenticateHeader(ctx,req.Header())
	if recPersonRow == nil {
		return nil ,connect.NewError(connect.CodeUnauthenticated,ErrCantAuthenticateUser)
	}

	senderId := req.Msg.FromPerson

	err := app.Queries.AcceptFriendRequest(ctx,database.AcceptFriendRequestParams{
		SenderID: int(senderId),
		ReceiverID: recPersonRow.ID,
	})
	if err != nil {
		slog.Error("can't accept friend request","err",err,"sender_id",recPersonRow.ID,"reciever_id",req.Msg.FromPerson)
		return nil,connect.NewError(connect.CodeInternal,errors.New("can't accept friendship"))
	}

	err = app.Queries.InsertFriend(ctx,database.InsertFriendParams{
		UserID: int(senderId),
		FriendID: recPersonRow.ID,
	})
	if err != nil {
		slog.Error("can't add friend","err",err,"sender_id",recPersonRow.ID,"reciever_id",req.Msg.FromPerson)
		return nil,connect.NewError(connect.CodeInternal,errors.New("can't add friend"))
	}

	return connect.NewResponse(&friendsv1.AcceptFriendReqResponse{
		Success: true,
	}),nil
}
