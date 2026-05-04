package rooms

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chigame/backend/games/games"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
)

type RoomSocketListener struct {
	Player *RoomPlayer
}

func NewRoomSocketListener(player *RoomPlayer) *RoomSocketListener {
	return &RoomSocketListener{
		Player: player,
	}
}
func (sl *RoomSocketListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *MessageCommand:
		if c.Reciever.ID == sl.Player.ID {
			err := SendChatMessageInWeb(c.Reciever, c.Sender, c.Text)
			if err != nil {
				slog.Error("Can not send chat message in room socket", "error", err)
			}
		}
	case *GameStartCommand:
		SendGametypeOverSocket(c.Room, sl.Player)
	}
}

func (room *Room) SocketRequestSendMsg(roomPlayer *RoomPlayer, chatMsgReq *roomv1.ChatMessageRequest) {
	if !room.Chat.IsOn {
		return
	}

	messageText := chatMsgReq.Text
	if len(messageText) > 256 {
		slog.Error("message is to long")
		return
	}

	senderID := roomPlayer.ID

	var senderPlayer *RoomPlayer
	var recieverPlayer *RoomPlayer

	for _, p := range room.Players {
		if p.ID == senderID {
			senderPlayer = p
		} else {
			recieverPlayer = p
		}
	}

	room.PushCommand(NewMessageCommand(room, messageText, senderPlayer, recieverPlayer))
}
func SendChatMessageInWeb(recieverPlayer *RoomPlayer, senderPlayer *RoomPlayer, messageText string) error {
	if recieverPlayer.Socket != nil {
		newChatMsg := &roomv1.RoomMessage{
			Content: &roomv1.RoomMessage_Chat{
				Chat: &roomv1.ChatMessage{
					PlayerId: int64(senderPlayer.ID),
					Text:     messageText,
				},
			},
		}
		return recieverPlayer.Socket.SendMessage(newChatMsg)
	}
	return fmt.Errorf("no socket found")
}

func SendGametypeOverSocket(room *Room, player *RoomPlayer) {
	if room.GameState == nil {
		slog.Error("can not send game type message game state is nil")
		return
	}
	gameData := room.GameState.GetGameData()
	currentGameType := gameData.GameType

	var protoGameType roomv1.GameType
	switch currentGameType {
	case games.XOGameType3X3:
		protoGameType = roomv1.GameType_GAME_TYPE_XO3X3
	case games.Conn4GameType:
		protoGameType = roomv1.GameType_GAME_TYPE_CONN4
	default:
		slog.Error("unknown game type to send", "game_type", currentGameType)
		return
	}

	message := roomv1.RoomMessage{Content: &roomv1.RoomMessage_GameType{
		GameType: &roomv1.ChangeGameTypeMessage{
			GameType: protoGameType,
		},
	}}

	player.Socket.SendMessage(&message)
}
