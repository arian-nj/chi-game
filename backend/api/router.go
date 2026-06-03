package api

import (
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/arian-nj/chigame/backend/gen/healthcheck/v1/healthcheckv1connect"
	"github.com/rs/cors"
)

var corsOrigins = []string{
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

	return withCORS(mux)
}

func withCORS(next http.Handler) http.Handler {
	allowedHeaders := []string{"Accept", "Authorization", "Content-Type"}
	allowedHeaders = append(allowedHeaders, connectcors.AllowedHeaders()...)

	c := cors.New(cors.Options{
		AllowedOrigins:   corsOrigins,
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   allowedHeaders,
		ExposedHeaders:   connectcors.ExposedHeaders(),
		MaxAge:           7200,
		AllowCredentials: true,
	})

	return c.Handler(next)
}
