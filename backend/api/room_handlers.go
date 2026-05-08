package api

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/socket"
	"github.com/arian-nj/chigame/backend/rooms"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

var (
	ErrorUnauthenticated = connect.NewError(connect.CodeUnauthenticated, errors.New("unknown user"))
)

func sendRoomSocketError(socketClient *socket.Socket, errType roomv1.RoomErrorType) {
	addEvent := roomv1.RoomMessage{
		Content: &roomv1.RoomMessage_Error{
			Error: errType,
		},
	}
	err := socketClient.SendMessage(&addEvent)
	if err != nil {
		slog.Error("", "error", err)
	}
}

func (app *ApiApplication) roomWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: CORS_PATTERNS,
	})
	if err != nil {
		slog.Error("error accepting new connection", "err", err)
		return
	}
	defer conn.Close(websocket.StatusNormalClosure, "websocket ended")

	socketClient := socket.NewSocketClient(conn)

	personRow, err := ContextGetAuthenticatedUser(app.Queries, r)
	if err != nil {
		sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_AUTH)
		return
	}

	currentRoom, found := app.AllRooms.Get(strconv.Itoa(personRow.ID))
	if found == false {
		sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_NOROOM)
		return
	}

	socketClient.ListenInBackground(r)

	var roomPlayer *rooms.RoomMember

	for _, RPlayer := range currentRoom.Players {
		if RPlayer.ID == personRow.ID {
			roomPlayer = RPlayer
			break
		}
	}
	if roomPlayer == nil {
		slog.Error("no player found in room", "id", personRow.ID)
		return
	}
	roomPlayer.Socket = socketClient
	rooms.SendGametypeOverSocket(currentRoom, roomPlayer)

	socketSubber := rooms.NewRoomSocketListener(roomPlayer)
	currentRoom.Subscribe(socketSubber)
	defer currentRoom.Unsubscribe(socketSubber)

	if currentRoom.GameState != nil {
		cancel := currentRoom.GameState.SubToSocket(personRow.ID, socketClient)
		if cancel == nil {
			sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_UNSPECIFIED)
			return
		}
		defer cancel()
	}

	for {
		select {
		case <-socketClient.Ctx.Done():
			slog.Info("socket context cancelled", "addr", r.RemoteAddr)
			return
		case newMsgBytes := <-roomPlayer.Socket.EventChan:
			newRoomMsg := &roomv1.RoomMessage{}
			err := proto.Unmarshal(newMsgBytes, newRoomMsg)
			if err != nil {
				slog.Error("can't unmarshal room msg", "error", err)
				continue
			}
			currentRoom.MsgChnl <- rooms.NewRoomEvent(roomPlayer, newRoomMsg)
		}
	}
}

func (app *ApiApplication) GetChatHistory(
	ctx context.Context,
	req *connect.Request[roomv1.GetChatHistoryRequest],
) (*connect.Response[roomv1.GetChatHistoryResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}
	gameRoom, ok := app.AllRooms.Get(strconv.Itoa(person.ID))
	if !ok {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("no game found"))
	}
	// NOTE: in case of dynamic Limit and offset cap limit by 50
	allMessages, err := app.Queries.GetRoomMessages(context.Background(), database.GetRoomMessagesParams{
		RoomID: gameRoom.ID,
		Limit:  50,
	})
	if err != nil {
		slog.Error("can't get message history", "error", err)
		return nil, connect.NewError(connect.CodeUnknown, errors.New("internal"))
	}

	messageLen := len(allMessages)
	// mhOut := chatHistoryOut{
	// 	Messages: make([]messageOut, messageLen),
	// }
	response := &roomv1.GetChatHistoryResponse{
		Messages: make([]*roomv1.ChatMessage, messageLen),
	}
	for index, message := range allMessages {
		response.Messages[messageLen-1-index] = &roomv1.ChatMessage{
			Id:       int64(message.ID),
			Text:     message.Content,
			PlayerId: int64(message.PersonID),
		}
	}

	return connect.NewResponse(response), nil
}

func (app *ApiApplication) GetRoomOpponent(
	ctx context.Context,
	req *connect.Request[roomv1.GetRoomOpponentRequest],
) (*connect.Response[roomv1.GetRoomOpponentResponse], error) {
	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}

	gs, gsExist := app.AllRooms.Get(strconv.Itoa(person.ID))
	if !gsExist {
		slog.Error("room not found", "person", person.ID)
		return nil, connect.NewError(connect.CodeNotFound, errors.New("room not found"))
	}

	var opponents *rooms.RoomMember
	for _, p := range gs.Players {
		if p.ID != person.ID {
			opponents = p
			break
		}
	}

	return connect.NewResponse(&roomv1.GetRoomOpponentResponse{
		Opponent: &accountv1.Account{
			Id:   int64(opponents.ID),
			Name: opponents.Name,
		},
	}), nil
}

func (app *ApiApplication) HasRoom(
	ctx context.Context,
	req *connect.Request[roomv1.HasRoomRequest],
) (*connect.Response[roomv1.HasRoomResponse], error) {

	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}
	_, hasRoom := app.AllRooms.Get(strconv.Itoa(person.ID))

	return connect.NewResponse(&roomv1.HasRoomResponse{
		HasRoom: hasRoom,
	}), nil
}

func (app *ApiApplication) CloseRoom(
	ctx context.Context,
	req *connect.Request[roomv1.CloseRoomRequest],
) (*connect.Response[roomv1.CloseRoomResponse], error) {

	person := app.AuthenticateHeader(ctx, req.Header())
	if person == nil {
		return nil, ErrorUnauthenticated
	}

	isOk := app.AllRooms.Delete(strconv.Itoa(person.ID))

	return connect.NewResponse(&roomv1.CloseRoomResponse{
		IsOk: isOk,
	}), nil
}
