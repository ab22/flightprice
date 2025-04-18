package api

import (
	"net/http"

	"github.com/ab22/flightprice/internal/api/middleware"
	"github.com/gorilla/mux"
)

type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
	Middlewares []mux.MiddlewareFunc
}

func (s *server) RegisterRoutes() *mux.Router {
	var (
		router           = mux.NewRouter()
		authMiddleware   = middleware.NewAuthMiddleware(s.logger, s.cfg.JWTSecretKey)
		loggerMiddleware = middleware.NewRequestLoggerMiddleware(s.logger)
	)

	// Global middlewares
	router.Use(loggerMiddleware)

	// Health endpoint
	router.HandleFunc("/ping", s.Ping)
	// Test Login
	router.HandleFunc("/login", s.TestLogin)

	// Websocket endpoints.
	router.HandleFunc("/subscribe/{freq}", s.Subscribe)

	// Flight endpoints
	flightsAPI := router.PathPrefix("/flights").Subrouter()
	flightsAPI.HandleFunc("/search", s.FlightsSearch)
	flightsAPI.Use(authMiddleware)

	return router
}
