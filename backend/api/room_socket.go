package api

import (
	"errors"
	"log/slog"
	"net/http"

	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/socket"
	"github.com/arian-nj/chigame/backend/internals/utils"
	"github.com/coder/websocket"
	"google.golang.org/protobuf/proto"
)

func (app *APIApplication) roomWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: CorsPatterns,
	})
	if err != nil {
		slog.Error("failed to accept websocket", "error", err)
		return
	}
	defer func() {
		conn.Close(websocket.StatusNormalClosure, "bye websocket")
		slog.Info("websocket closed", "remote addr", r.RemoteAddr)
	}()

	socketClient := socket.NewSocketClient(conn)

	person := app.AuthenticateQuery(r.Context(), *r)
	if person == nil {
		sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_AUTH)
		return
	}
	// check room code
	code := r.URL.Query().Get("code")
	if code == "" {
		sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_INVALID_CODE)
		return
	}
	code = normalizeRoomCode(code)

	currentRoom, ok := app.RoomsStore.GetByCode(code)
	if !ok {
		sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_INVALID_CODE)
		return
	}

	roomMember := NewRoomMember(person, socketClient)
	if err := currentRoom.AddMember(roomMember); err != nil {
		if errors.Is(err, errRoomFull) {
			sendRoomSocketError(socketClient, roomv1.RoomErrorType_ROOM_ERROR_TYPE_INVALID_CODE)
		}
		return
	}

	// listen
	utils.RunBackgroundTask(func() {
		socketClient.Listen(r)
	})

	for {
		select {
		case <-socketClient.Ctx.Done():
			slog.Info("socket context cancelled", "addr", r.RemoteAddr)
			currentRoom.RemoveMember(roomMember)
			app.RoomsStore.MaybeDeleteRoom(code, person.ID)
			return
		case newMsgBytes := <-socketClient.EventChan:
			newRoomMsg := &roomv1.RoomMessage{}
			err := proto.Unmarshal(newMsgBytes, newRoomMsg)
			if err != nil {
				slog.Error("failed to unmarshal message", "error", err)
				continue
			}
			currentRoom.MsgChnl <- NewRoomEvent(roomMember, newRoomMsg)
		}
	}
}

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
