package api

import (
	"log"
	"net/http"

	connectcors "connectrpc.com/cors"
	"github.com/arian-nj/chigame/backend/gen/account/v1/accountv1connect"
	"github.com/arian-nj/chigame/backend/gen/auth/v1/authv1connect"
	"github.com/arian-nj/chigame/backend/gen/healthcheck/v1/healthcheckv1connect"
	"github.com/arian-nj/chigame/backend/gen/session/v1/sessionv1connect"
	"github.com/rs/cors"
)

var CORS_PATTERNS = []string{
	"*",
	"http://localhost:8080", "https://localhost:8080", "localhost:8080",
	"http://localhost:5173", "https://localhost:5173", "localhost:5173",
	"http://localhost:3000", "https://localhost:3000", "localhost:3000",
	"https://chigame.ir", "chigame.ir",
}

func (app *ApiApplication) createRouter() http.Handler {
	mux := http.NewServeMux()

	// if app.Config.ReleaseMode == config.Develop {
	// 	mux.Use(func(next http.Handler) http.Handler {
	// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 			log.Printf("Request Received: Method=%s, Path=%s, RemoteAddr=%s", r.Method, r.URL.Path, r.RemoteAddr)
	// 			next.ServeHTTP(w, r)
	// 		})
	// 	})
	// }

	sessionPath, sessionHandler := sessionv1connect.NewSessionServiceHandler(app)
	mux.Handle(sessionPath, sessionHandler)

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

	mux.Handle("/api/session/", app.AuthenticateQuery(http.HandlerFunc(app.sessionWebsocket)))
	mux.Handle("/api/match_making/ticket/", app.AuthenticateQuery(http.HandlerFunc(app.makeMatchMakingTicketWS)))

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

// cannot use app (variable of type *ApiApplication) as sessionv1connect.SessionServiceHandler value in argument to sessionv1connect.NewSessionServiceHandler: *ApiApplication does not implement sessionv1connect.SessionServiceHandler (wrong type for method GetChatHistory)
// 		have GetChatHistory(context.Context, *"connectrpc.com/connect".Request[sessionv1.GetChatHistoryRequest]) (*"connectrpc.com/connect".Response[sessionv1.GetChatHistoryResponse], error)
// 		want GetChatHistory(context.Context, *"github.com/bufbuild/connect-go".Request[sessionv1.GetChatHistoryRequest]) (*"github.com/bufbuild/connect-go".Response[sessionv1.GetChatHistoryResponse], error)
