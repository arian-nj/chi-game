package api

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/arian-nj/chigame/backend/games/conn4"
	"github.com/arian-nj/chigame/backend/games/games"
	"github.com/arian-nj/chigame/backend/games/xo"
	matchmaking "github.com/arian-nj/chigame/backend/match_making"
	"github.com/arian-nj/chigame/backend/rooms"
)

func (app *ApiApplication) MakeMatches() {
	defer app.MatchMaking.Mutex.Unlock()
	var doFlag = false
	for {
		doFlag = false
		for gameTypeKey, ticketsList := range app.MatchMaking.WaitingPlayers {
			app.MatchMaking.Mutex.Lock()
			if len(ticketsList) >= 2 {
				doFlag = true
				ticketOne := ticketsList[0]
				ticketTwo := ticketsList[1]
				app.MatchMaking.WaitingPlayers[gameTypeKey] = ticketsList[2:]
				app.createRandomGame(gameTypeKey, ticketOne, ticketTwo)
			}
			app.MatchMaking.Mutex.Unlock()
		}
		if !doFlag {
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func (app *ApiApplication) createRandomGame(gameType games.GameType, ticketOne *matchmaking.Ticket, ticketTwo *matchmaking.Ticket) {
	newRoomRow, err := app.Queries.CreateRoom(context.Background(), string(rooms.RandomRoom))
	if err != nil {
		slog.Error("can't create new random game", "error", err)
	}
	newGameRoom := rooms.NewRoom(app.Queries, newRoomRow.ID, app.AllRooms)
	// newGameRoom.Subscribe(rooms.NewRoomTelegramBotListener(ticketOne.UserID, ticketOne.TgID, app.Bot, ""))
	// newGameRoom.Subscribe(rooms.NewRoomTelegramBotListener(ticketTwo.UserID, ticketTwo.TgID, app.Bot, ""))

	var newGame games.Game

	switch gameType {
	case games.XOGameType3X3, games.XOGameType5X5:
		newXoGame := xo.NewXOGame(newGameRoom.RoomCtx, games.XOGameType3X3, app.Queries)
		newGame = newXoGame
	case games.Conn4GameType:
		newConn4Game := conn4.NewConn4State(newGameRoom.RoomCtx, games.Conn4GameType, app.Queries)
		newGame = newConn4Game

	default:
		slog.Error("not possible random game")
		return
	}
	newGameRoom.GameState = newGame
	newGameRoom.RunBgMonitor()

	playerOne := rooms.NewRoomPlayer(ticketOne.UserID, ticketOne.Name)
	playerTwo := rooms.NewRoomPlayer(ticketTwo.UserID, ticketTwo.Name)

	newGameRoom.AddPlayerToRoom(playerOne)
	newGameRoom.AddPlayerToRoom(playerTwo)

	app.AllRooms.Add(strconv.Itoa(playerOne.ID), newGameRoom)
	app.AllRooms.Add(strconv.Itoa(playerTwo.ID), newGameRoom)

	for _, ticket := range []*matchmaking.Ticket{ticketOne, ticketTwo} {
		ticket.MatchFoundChan <- newGameRoom
	}

	newGameRoom.StartGame()
}
