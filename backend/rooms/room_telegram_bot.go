package rooms

import (
	"gopkg.in/telebot.v4"
)

type RoomTelegramBotListener struct {
	Bot    *telebot.Bot
	UserID int
	TgID   int
}

// func NewRoomTelegramBotListener(playerID int, TgID int, bot *telebot.Bot, viaMessageID string) *RoomTelegramBotListener {
// 	return &RoomTelegramBotListener{
// 		Bot:    bot,
// 		UserID: playerID,
// 		TgID:   TgID,
// 	}
// }

// func (tg *RoomTelegramBotListener) Update(command commander.Command) {
// 	switch c := command.(type) {
// 	case *MessageCommand:
// 		if c.Reciever.ID == tg.UserID {
// 			err := tg.SendChatMessageInBot(tg.TgID, c.Text, c.Sender.Name)
// 			if err != nil {
// 				slog.Error("Can not send chat message in room telegram", "error", err)
// 			}
// 		}

// 	case *GameEndedCommand:
// 		room := c.Room
// 		if room.Chat.IsOn == false {
// 			return
// 		}
// 		text := "چت قطع شد"
// 		_, err := room.Bot.Send(&telebot.User{ID: int64(tg.TgID)}, text, keybul.WelcomeReplyKeyboard)
// 		if err != nil {
// 			slog.Error("can't send chat ended message", "err", err)
// 		}

// 		text = fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(ExpirationDur.Seconds()))
// 		_, err = c.Room.Bot.Send(&telebot.User{ID: int64(tg.TgID)}, text)
// 		if err != nil {
// 			slog.Error("can't send end game chat message", "err", err)
// 		}
// 	case *GameStartCommand:
// 		SendFoundOpponentMessage(c.Room.Players, tg.Bot)
// 	}
// }

// // Handler
// func (room *GameRoom) BotRequestSendMsg(bot telebot.API, senderID int, messageText string) error {
// 	// if !room.Chat.IsOn {
// 	// 	return nil
// 	// }
// 	// if len(messageText) > 256 {
// 	// 	slog.Error("message is to long")
// 	// 	return nil
// 	// }

// 	// var senderPlayer *RoomPlayer
// 	// var recieverPlayer *RoomPlayer

// 	// for _, p := range Room.Players {
// 	// 	if p.TgID == senderID {
// 	// 		senderPlayer = p
// 	// 	} else {
// 	// 		recieverPlayer = p
// 	// 	}
// 	// }
// 	// room.PushCommand(NewMessageCommand(room, messageText, senderPlayer, recieverPlayer))
// 	return nil
// }

// func (tg *RoomTelegramBotListener) SendChatMessageInBot(toId int, text string, senderName string) error {
// 	_, err := tg.Bot.Send(&telebot.User{ID: int64(toId)},
// 		fmt.Sprintf("*_%s:_* %s", senderName, text), telebot.ModeMarkdownV2)
// 	return err
// }
