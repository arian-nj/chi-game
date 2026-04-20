package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/arian-nj/chigame/backend/api"
	"github.com/arian-nj/chigame/backend/db"
	"github.com/arian-nj/chigame/backend/internals/config"
)

const Version = "1.0.0"

type Config struct {
	port        int
	ReleaseMode string
	DatabseUrl  string
}

type Application struct {
	config *config.Config
	logger *log.Logger
	Wg     *sync.WaitGroup
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

	parentCtx, pCancel := context.WithCancel(context.Background())
	defer pCancel()

	app.logger.Printf("Serving...")
	apiApp := api.NewApiApplication(app.config, nil)
	apiApp.RunApi(parentCtx, app.Wg)

}
