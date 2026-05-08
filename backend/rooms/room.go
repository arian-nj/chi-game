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
)

type RoomType string

const (
	PrivateRoom RoomType = "private"
	RandomRoom  RoomType = "random"
)

type Room struct {
	ID int

	Queries *database.Queries

	IsGameEnded bool
	ChatIsOn    bool
	GameState   games.Game

	MsgChnl chan *RoomEvent

	Players []*RoomMember

	CancelRoom context.CancelFunc
	RoomCtx    context.Context

	CreatedAt       time.Time
	ExpireDuaration time.Duration

	allRooms *AllRooms

	*commander.Commander
}

func NewRoom(Queries *database.Queries, roomId int, allRoom *AllRooms) *Room {
	ctx, cancel := context.WithCancel(context.Background())
	gs := &Room{
		ID:              roomId,
		CreatedAt:       time.Now(),
		ChatIsOn:        true,
		Players:         []*RoomMember{},
		ExpireDuaration: 2*time.Minute*2 + 30,
		MsgChnl:         make(chan *RoomEvent, 10),

		Queries: Queries,

		CancelRoom: cancel,
		RoomCtx:    ctx,

		Commander: commander.NewCommander(),
		allRooms:  allRoom,
	}
	return gs
}

type RoomMember struct {
	ID int
	// TgID   int
	Name   string
	Socket *socket.Socket
}

func NewRoomPlayer(ID int, name string) *RoomMember {
	return &RoomMember{
		ID:   ID,
		Name: keybul.EscapeReserved(name),
	}
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

		room.GameState.OnEnd(room.EndGame)
		gameErr := room.GameState.StartGame()
		if gameErr != nil {
			slog.Error("gamer error", "error", gameErr)
			return
		}
	})
}

func (room *Room) AddPlayerToRoom(player *RoomMember) {
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

	if room.ChatIsOn == false {
		return
	}

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
		}
	}
}
