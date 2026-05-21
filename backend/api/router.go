package api

import (
	"log"
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/arian-nj/chigame/backend/gen/account/v1/accountv1connect"
	"github.com/arian-nj/chigame/backend/gen/auth/v1/authv1connect"
	"github.com/arian-nj/chigame/backend/gen/friends/v1/friendsv1connect"
	"github.com/arian-nj/chigame/backend/gen/healthcheck/v1/healthcheckv1connect"
	"github.com/arian-nj/chigame/backend/gen/room/v1/roomv1connect"
	"github.com/rs/cors"
)

var CORS_PATTERNS = []string{
	"*",
	"http://localhost:8080", "https://localhost:8080", "localhost:8080",
	"http://localhost:5173", "https://localhost:5173", "localhost:5173",
	"http://localhost:3000", "https://localhost:3000", "localhost:3000",
	"https://chigame.ir", "chigame.ir",
}

func (app *APIApplication) createRouter() http.Handler {
	mux := http.NewServeMux()

	// if app.Config.ReleaseMode == config.Develop {
	// 	mux.Use(func(next http.Handler) http.Handler {
	// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			log.Printf("Request Received: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)
	// 			next.ServeHTTP(w, r)
	// 		})
	// 	})
	// }

	roomPath, roomHandler := roomv1connect.NewRoomServiceHandler(app)
	mux.Handle(roomPath, roomHandler)

	// if app.Config.ReleaseMode == config.Develop {
	// 	dummyAuthPath, dummyAuthHandler := dummy_authv1connect.NewDummyAuthServiceHandler(app)
	// 	mux.Handle(dummyAuthPath, dummyAuthHandler)
	// }

	accountPath, accountHandler := accountv1connect.NewAccountServiceHandler(app)
	mux.Handle(accountPath, accountHandler)

	healthPath, healthHandler := healthcheckv1connect.NewHealthcheckServiceHandler(app)
	mux.Handle(healthPath, healthHandler)

	authPath, authHandler := authv1connect.NewAuthServiceHandler(app)
	mux.Handle(authPath, authHandler)

	friendsPath, friendsHandler := friendsv1connect.NewFriendsServiceHandler(app)
	mux.Handle(friendsPath, friendsHandler)

	mux.Handle("/api/room/", app.AuthenticateQuery(http.HandlerFunc(app.roomWebsocket)))
	mux.Handle("/api/match_making/ticket/", app.AuthenticateQuery(http.HandlerFunc(app.finderWS)))

	return withCORS(mux)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Request URI:", r.RequestURI)

		next.ServeHTTP(w, r)
	})
}
func withCORS(next http.Handler) http.Handler {

	allowdHeaders := []string{"Accept", "Authorization", "Content-Type"}
	allowdHeaders = append(allowdHeaders, connectcors.AllowedHeaders()...)

	c := cors.New(cors.Options{
		AllowedOrigins: CORS_PATTERNS,
		// AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", , connectcors.AllowedMethods(),
		// AllowedHeaders:   []string{},
		// ExposedHeaders:   []string{"Link"},
		AllowedMethods:   connectcors.AllowedMethods(),
		AllowedHeaders:   allowdHeaders,
		ExposedHeaders:   connectcors.ExposedHeaders(),
		MaxAge:           7200,
		AllowCredentials: true,
	})

	return c.Handler(next)

}
