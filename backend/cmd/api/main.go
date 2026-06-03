package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/arian-nj/chigame/backend/api"
	"github.com/arian-nj/chigame/backend/database"
	"github.com/arian-nj/chigame/backend/db"
	"github.com/arian-nj/chigame/backend/internals/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

const version = "0.1.0"

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})))

	cfg, err := config.ParseConfig()
	if err != nil {
		slog.Error("config error", "err", err)
		os.Exit(1)
	}

	if err := db.Migrate(cfg.DatabaseURL); err != nil {
		slog.Error("migration failed", "err", err)
		os.Exit(1)
	}

	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		slog.Error("database connection failed", "err", err)
		os.Exit(1)
	}
	defer pool.Close()

	queries := database.New(pool)
	apiApp := api.NewAPIApplication(cfg, queries)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup
	slog.Info("chi game api", "version", version)
	apiApp.RunAPI(ctx, &wg)
	wg.Wait()
}
