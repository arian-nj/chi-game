package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"sync"

	"github.com/arian-nj/chigame/backend/api"
	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/db"
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

	app.logger.Printf("Serving...")
	apiApp := api.NewApiApplication(app.config, app.Queries, app.logger, app.AllRooms, app.MatchMaking)
	go apiApp.MakeMatches()
	apiApp.RunApi(parentCtx, app.Wg)

}
