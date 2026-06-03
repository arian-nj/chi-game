package xo

import (
	"fmt"
	"strconv"

	"github.com/arian-nj/chigame/backend/internals/chronos"
	"github.com/arian-nj/chigame/backend/internals/socket"
)

type XoPlayer struct {
	ID         int
	TelegramID int

	Name string

	MessageID int

	Socket *socket.Socket

	Timer *chronos.Timer
	Move  Cell
}

func NewXoPlayer(id int, name string, socket *socket.Socket) *XoPlayer {
	if len(name) > 20 {
		name = name[:20]
		name += "..."
	}
	return &XoPlayer{
		// Name: keybul.EscapeReserved(name),
		ID:    id,
		Name:  fmt.Sprintf("`%s`", name),
		Timer: chronos.NewTimer(MAX_ALLOWED_TIME),
	}
}

func (p *XoPlayer) MessageSig() (string, int64) {
	return strconv.Itoa(p.MessageID), int64(p.TelegramID)
}

func (p *XoPlayer) Recipient() string {
	return strconv.FormatInt(int64(p.TelegramID), 10)
}
