package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/arian-nj/chigame/backend/api"
	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/db"
	"github.com/arian-nj/chigame/backend/games/conn4"
	"github.com/arian-nj/chigame/backend/games/games"
	"github.com/arian-nj/chigame/backend/games/xo"
	"github.com/arian-nj/chigame/backend/internals/config"
	matchmaking "github.com/arian-nj/chigame/backend/match_making"
	"github.com/arian-nj/chigame/backend/rooms"
	"github.com/jackc/pgx/v5/pgxpool"
)

const Version = "1.0.0"

type Config struct {
	port        int
	ReleaseMode string
	DatabseUrl  string
}

type Application struct {
	config      *config.Config
	logger      *log.Logger
	Wg          *sync.WaitGroup
	Queries     *database.Queries
	MatchMaking *matchmaking.MatchMaking
	AllRooms    *rooms.AllRooms
}

func main() {
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	cfg, err := config.ParseConfig()
	if err != nil {
		logger.Fatal(err)
	}

	app := Application{
		config: cfg,
		logger: logger,
		Wg:     &sync.WaitGroup{},
	}

	err = db.Migrate(cfg.DatabseUrl)
	if err != nil {
		app.logger.Fatalln("can not migrate", err)
	}

	conn, err := pgxpool.New(context.Background(), app.config.DatabseUrl)
	if err != nil {
		slog.Error("can not make a new connection ", "err", err)
		return
	}
	defer conn.Close()

	app.Queries = database.New(conn)

	app.AllRooms = rooms.NewAllRooms()
	app.MatchMaking = matchmaking.NewMatchMaking(app.AllRooms, app.Queries)

	parentCtx, pCancel := context.WithCancel(context.Background())
	defer pCancel()

	go app.MakeMatches()

	app.logger.Printf("Serving...")
	apiApp := api.NewApiApplication(app.config, app.Queries, app.logger, app.AllRooms, app.MatchMaking)
	apiApp.RunApi(parentCtx, app.Wg)

}

func (app *Application) MakeMatches() {
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

func (app *Application) createRandomGame(gameType games.GameType, ticketOne *matchmaking.Ticket, ticketTwo *matchmaking.Ticket) {
	newRoomRow, err := app.Queries.CreateRoom(context.Background(), string(rooms.RandomRoom))
	if err != nil {
		slog.Error("can't create new random game", "error", err)
	}
	newGameRoom := rooms.NewGameRoom(app.Queries, newRoomRow.ID, app.AllRooms)
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
