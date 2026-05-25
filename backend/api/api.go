// Package api contains all api rpcs
package api

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/internals/config"
	matchmaking "github.com/arian-nj/chigame/backend/match_making"
	"github.com/arian-nj/chigame/backend/rooms"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var ErrCantAuthenticateUser = errors.New("can't get user")

type APIApplication struct {
	Config      *config.Config
	Queries     *database.Queries
	AllRooms    *rooms.AllRooms
	MatchMaking *matchmaking.MatchMaking
	// InboxMap    *AllInbox
}

func NewAPIApplication(config *config.Config,
	queries *database.Queries,
	logger *log.Logger,
	AllRoom *rooms.AllRooms,
	matchMaking *matchmaking.MatchMaking,
) *APIApplication {
	return &APIApplication{
		Config:      config,
		Queries:     queries,
		AllRooms:    AllRoom,
		MatchMaking: matchMaking,
		// InboxMap:    NewAllInbox(),
	}
}

func (app *APIApplication) RunAPI(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()

	router := app.createRouter()
	srv := &http.Server{
		Addr:         ":8383",
		Handler:      h2c.NewHandler(router, &http2.Server{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		slog.Info("Starting server on " + srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen error", "err", err)
			return
		}
	}()

	<-ctx.Done()

	slog.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("server forces to shutdown ", "err", err.Error())
	}
	slog.Info("Server exited")
}
