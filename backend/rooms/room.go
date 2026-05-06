package rooms

import (
	"context"
	"log/slog"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/games/games"
	roomv1 "github.com/arian-nj/chigame/backend/gen/room/v1"
	"github.com/arian-nj/chigame/backend/internals/commander"
	"github.com/arian-nj/chigame/backend/internals/keybul"
	"github.com/arian-nj/chigame/backend/internals/socket"
	"github.com/arian-nj/chigame/backend/internals/utils"
	"gopkg.in/telebot.v4"
)

type RoomType string

const (
	PrivateRoom RoomType = "private"
	RandomRoom  RoomType = "random"
)

const ExpirationDur = 30 * time.Second

type RoomPlayer struct {
	ID int
	// TgID   int
	Name   string
	Socket *socket.Socket
}

func NewRoomPlayer(ID int, name string) *RoomPlayer {
	return &RoomPlayer{
		ID: ID,
		// TgID: tgID,
		Name: keybul.EscapeReserved(name),
	}
}

type Chat struct {
	IsOn bool
}

type Room struct {
	ID int

	Bot     *telebot.Bot
	Queries *database.Queries

	IsGameEnded bool
	Chat        Chat
	GameState   games.Game

	MsgChnl chan *RoomEvent

	Players []*RoomPlayer

	CancelRoom context.CancelFunc
	RoomCtx    context.Context

	CreatedAt       time.Time
	ExpireDuaration time.Duration

	allRooms *AllRooms

	*commander.Commander

	ShutdownTimer <-chan time.Time
}

func NewRoom(Queries *database.Queries, roomId int, allRoom *AllRooms) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	gs := &Room{
		ID: roomId,
		// Bot:       bot,
		CreatedAt: time.Now(),
		Chat: Chat{
			IsOn: true,
		},
		Players:         []*RoomPlayer{},
		ExpireDuaration: 2*time.Minute*2 + 30,
		MsgChnl:         make(chan *RoomEvent, 10),

		Queries: Queries,

		CancelRoom: cancel,
		RoomCtx:    ctx,

		Commander: commander.NewCommander(),
		allRooms:  allRoom,
		// ShutdownTimer: make(<-chan time.Time),
	}
	return gs
}

func (room *Room) RunBgMonitor() {
	utils.RunBackgroundTask(func() {
		room.MonitorGameRoom()
	})
}

func (room *Room) StartGame() {
	utils.RunBackgroundTask(func() {
		if room.GameState == nil {
			slog.Error("failed to start game field GameState is nil")
			return
		}
		room.PushCommand(
			NewGameStartRoomCommand(room),
		)

		for _, player := range room.Players {
			room.GameState.AddPlayer(player.ID, player.Name, player.Socket)
		}

		for _, suber := range room.Subscribers {
			// if listener, ok := suber.(*RoomTelegramBotListener); ok {
			// 	room.GameState.SubToTelegram(listener.UserID, listener.Bot, "")
			// }
			if listener, ok := suber.(*RoomTelegramViaListener); ok {
				room.GameState.SubToTelegram(0, listener.Bot, listener.ViaMessageId)
			}
		}

		room.GameState.OnEnd(room.EndGame)
		gameErr := room.GameState.StartGame()
		if gameErr != nil {
			slog.Error("gamer error", "error", gameErr)
			return
		}
	})
}

func (room *Room) AddPlayerToRoom(player *RoomPlayer) {
	room.Players = append(room.Players, player)
	utils.RunBackgroundTask(func() {
		_, err := room.Queries.CreateRoomPlayer(context.Background(), database.CreateRoomPlayerParams{
			RoomID:   room.ID,
			PersonID: player.ID,
		})
		if err != nil {
			slog.Error("can't add room player", "error", err)
		}
	})
}

func (room *Room) EndGame() {
	room.IsGameEnded = true

	gameEnded := NewGameEndedRoomCommand(room)
	room.PushCommand(gameEnded)

	if room.Chat.IsOn == false {
		return
	}
	// room.ShutdownTimer = time.After(30 * time.Second)

}

func (room *Room) MonitorGameRoom() {
	for {
		select {
		case newRoomEvent := <-room.MsgChnl:
			switch newRoomEvent.Event.Content.(type) {
			case *roomv1.RoomMessage_ChatReq:
				room.SocketRequestSendMsg(newRoomEvent.Player, newRoomEvent.Event.GetChatReq())
			case *roomv1.RoomMessage_Game:
				room.GameState.SocketRouter(newRoomEvent.Event.GetGame(), newRoomEvent.Player.ID)
			}
		case <-room.CommandNotifire:
			if len(room.Commands) > 0 {
				com := room.PopCommand()
				room.ApplyCommand(com)
			}
			// case <-room.ShutdownTimer:
			// room.CleanAndDisconnect()
			// return
		}
	}
}

// func (room *Room) CleanAndDisconnect() {
// 	room.allRooms.Mutex.Lock()
// 	defer room.allRooms.Mutex.Unlock()

// 	for _, player := range room.Players {
// 		delete(room.allRooms.Rooms, strconv.Itoa(player.ID))
// 	}
// }

func (room *Room) HandleCallback(c telebot.Context, queries *database.Queries, allRoom *AllRooms) error {
	// callbackData := c.Callback().Data
	// if callbackData == "join" {
	// 	personRow, err := queries.GetPersonByID(context.Background(), int(c.Sender().ID))
	// 	if err != nil {
	// 		slog.Error("can not get user at room handle callback", "err", err)
	// 		return c.RespondText("خطا")
	// 	}
	// 	if c.Sender().ID == int64(room.Players[0].TgID) {
	// 		text := "خودت بازیو ساختی تو بازی هستی"
	// 		return c.RespondText(text)
	// 	}

	// 	newJoinCommand := NewJoinRoomCommand(room, NewRoomPlayer(personRow.ID, personRow.TgID, personRow.Name), allRoom)
	// 	room.PushCommand(newJoinCommand)

	// 	text := "اضافه شدی بازی شروع شد"
	// 	err = c.RespondText(text)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// }

	// err := room.GameState.CallBackRouter(c)
	// if err != nil {
	// 	slog.Error("error in call back router", "error", err)
	// 	return c.RespondText("خطا")
	// }
	return nil
}
