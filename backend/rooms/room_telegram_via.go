package rooms

import (
	"fmt"
	"log/slog"

	"github.com/arian-nj/chigame/backend/internals/commander"
	"github.com/arian-nj/chigame/backend/internals/keybul"
	"gopkg.in/telebot.v4"
)

type RoomTelegramViaListener struct {
	Bot          *telebot.Bot
	ViaMessageId string // Via Bots
}

func NewRoomTelegramViaListener(bot *telebot.Bot, viaMessageID string) *RoomTelegramViaListener {
	return &RoomTelegramViaListener{
		Bot:          bot,
		ViaMessageId: viaMessageID,
	}
}

func (tg *RoomTelegramViaListener) MessageSig() (string, int64) {
	return tg.ViaMessageId, 0
}

func (tg *RoomTelegramViaListener) Update(command commander.Command) {
	switch c := command.(type) {
	case *WaitForPlayerCommand:
		err := tg.SendWaitPanel(c)
		if err != nil {
			slog.Error("can not send wait panel", "error", err)
		}
	}
}

func (tg *RoomTelegramViaListener) SendWaitPanel(wait *WaitForPlayerCommand) error {
	room := wait.Room
	gameData := room.GameState.GetGameData()
	creator := wait.Creator
	inlineKeyboard := keybul.CreateInlineKeyboard(
		keybul.JoinGameInlineButtons,
	)
	text := gameData.StartText + "\n\n" + gameData.RulesText + "\n\n🕹 بازیکن " + fmt.Sprintf("%s", creator.Name) + " منتظر حریفه"

	return keybul.EditMessage(tg.Bot, tg, text, inlineKeyboard)
}
