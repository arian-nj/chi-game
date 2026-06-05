package api

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/arian-nj/chigame/backend/gen/account/v1/accountv1connect"
	"github.com/arian-nj/chigame/backend/gen/admin/v1/adminv1connect"
	"github.com/arian-nj/chigame/backend/gen/auth/v1/authv1connect"
	"github.com/arian-nj/chigame/backend/gen/healthcheck/v1/healthcheckv1connect"
	"github.com/arian-nj/chigame/backend/gen/invite/v1/invitev1connect"
	"github.com/rs/cors"
)

var CorsPatterns = []string{
	"http://localhost:8080",
	"https://localhost:8080",
	"http://localhost:5173",
	"https://localhost:5173",
	"http://localhost:3000",
	"https://localhost:3000",
	"https://chigame.site",
}

func (app *APIApplication) createRouter() http.Handler {
	mux := http.NewServeMux()

	healthPath, healthHandler := healthcheckv1connect.NewHealthcheckServiceHandler(app)
	mux.Handle(healthPath, healthHandler)

	authPath, authHandler := authv1connect.NewAuthServiceHandler(app)
	mux.Handle(authPath, authHandler)

	accountPath, accountHandler := accountv1connect.NewAccountServiceHandler(app)
	mux.Handle(accountPath, accountHandler)

	adminPath, adminHandler := adminv1connect.NewAdminServiceHandler(app)
	mux.Handle(adminPath, adminHandler)

	invitePath, inviteHandler := invitev1connect.NewInviteServiceHandler(app)
	mux.Handle(invitePath, inviteHandler)

	mux.Handle("/room/websocket", http.HandlerFunc(app.roomWebsocket))

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	allowedHeaders := []string{"Accept", "Authorization", "Content-Type"}
	allowedHeaders = append(allowedHeaders, connectcors.AllowedHeaders()...)

	c := cors.New(cors.Options{
		AllowedOrigins:   CorsPatterns,
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   connectcors.ExposedHeaders(),
		MaxAge:           7200,
		AllowCredentials: true,
	})

	return c.Handler(next)
}
