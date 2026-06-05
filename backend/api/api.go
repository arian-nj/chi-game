// Package api contains Connect RPC handlers and HTTP routing.
package api

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/internals/config"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type APIApplication struct {
	Config      *config.Config
	Queries     *database.Queries
	InviteStore *InviteStore
	RoomsStore  *RoomsStore
}

func NewAPIApplication(cfg *config.Config, queries *database.Queries) *APIApplication {
	return &APIApplication{
		Config:      cfg,
		Queries:     queries,
		InviteStore: NewInviteStore(),
		RoomsStore:  NewRoomsStore(),
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
		slog.Info("starting server", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen error", "err", err)
		}
	}()

	<-ctx.Done()

	slog.Info("shutting down server")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("server shutdown error", "err", err)
	}
	slog.Info("server exited")
}
