package rooms

import (
	"context"
	"log/slog"
	"strconv"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/internals/utils"
)

type MessageCommand struct {
	Text     string
	Sender   *RoomPlayer
	Reciever *RoomPlayer
	Room     *Room
}

func NewMessageCommand(room *Room, text string, senderPlayer *RoomPlayer, recPlayer *RoomPlayer) *MessageCommand {
	return &MessageCommand{
		Text:     text,
		Room:     room,
		Sender:   senderPlayer,
		Reciever: recPlayer,
	}
}

func (message *MessageCommand) Execute() {
	room := message.Room

	utils.RunBackgroundTask(func() {
		_, err := room.Queries.CreateRoomMessage(context.Background(), database.CreateRoomMessageParams{
			RoomID:   room.ID,
			PlayerID: message.Sender.ID,
			Message:  message.Text,
		})
		if err != nil {
			slog.Error("error creating new message in db")
		}
	})
}

type WaitForPlayerCommand struct {
	Room    *Room
	Creator *RoomPlayer
}

func NewWaitForPlayerCommand(room *Room, creatorUser *RoomPlayer) *WaitForPlayerCommand {
	return &WaitForPlayerCommand{
		Room:    room,
		Creator: creatorUser,
	}
}

func (wait *WaitForPlayerCommand) Execute() {
}

type JoinRoomCommand struct {
	Room         *Room
	JoinedPlayer *RoomPlayer
	AllRoom      *AllRooms
}

func NewJoinRoomCommand(room *Room, JoinedUser *RoomPlayer, allRoom *AllRooms) *JoinRoomCommand {
	return &JoinRoomCommand{
		Room:         room,
		JoinedPlayer: JoinedUser,
		AllRoom:      allRoom,
	}
}

func (join *JoinRoomCommand) Execute() {
	join.Room.AddPlayerToRoom(join.JoinedPlayer)
	for _, player := range join.Room.Players {
		join.AllRoom.Add(strconv.Itoa(player.ID), join.Room)
	}
	join.Room.StartGame()
}

type GameEndedCommand struct {
	Room *Room
}

func NewGameEndedRoomCommand(room *Room) *GameEndedCommand {
	return &GameEndedCommand{
		Room: room,
	}
}

func (EndGame *GameEndedCommand) Execute() {
}

type GameStartCommand struct {
	Room *Room
}

func NewGameStartRoomCommand(room *Room) *GameEndedCommand {
	return &GameEndedCommand{
		Room: room,
	}
}

func (EndGame *GameStartCommand) Execute() {
}

type RequestEndRoomCommand struct {
	Room   *Room
	Player *RoomPlayer
}

func NewRequestEndRoomCommand(room *Room, player *RoomPlayer) *RequestEndRoomCommand {
	return &RequestEndRoomCommand{
		Room:   room,
		Player: player,
	}
}

func (ReqEnd *RequestEndRoomCommand) Execute() {
}
