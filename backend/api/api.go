package api

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"connectrpc.com/connect"
	"github.com/arian-nj/chigame/backend/database"
	gamesessions "github.com/arian-nj/chigame/backend/game_sessions"
	accountv1 "github.com/arian-nj/chigame/backend/gen/account/v1"
	"github.com/arian-nj/chigame/backend/internals/config"
	matchmaking "github.com/arian-nj/chigame/backend/match_making"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type ApiApplication struct {
	Config      *config.Config
	Queries     *database.Queries
	AllSessions *gamesessions.AllSession
	MatchMaking *matchmaking.MatchMaking
}

// GetMe implements [accountv1connect.AccountServiceHandler].
func (app *ApiApplication) GetMe(ctx context.Context, req *connect.Request[accountv1.GetMeRequest]) (*connect.Response[accountv1.GetMeResponse], error) {
	personRow := app.AuthenticateHeader(ctx, req.Header())
	if personRow == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, errors.New("can't get user"))
	}
	return connect.NewResponse(&accountv1.GetMeResponse{
		Account: &accountv1.Account{
			Id:   int64(personRow.ID),
			Name: personRow.Username,
		},
	}), nil
}

func NewApiApplication(config *config.Config,
	queries *database.Queries,
	logger *log.Logger,
	AllSession *gamesessions.AllSession,
	matchMaking *matchmaking.MatchMaking,
) *ApiApplication {
	return &ApiApplication{
		Config:      config,
		Queries:     queries,
		AllSessions: AllSession,
		MatchMaking: matchMaking,
	}
}

func (app *ApiApplication) RunApi(ctx context.Context, wg *sync.WaitGroup) {
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
