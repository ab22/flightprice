package api

import (
	"net/http"

	"github.com/ab22/flightprice/internal/api/middleware"
	"github.com/ab22/flightprice/internal/config"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
	Middlewares []mux.MiddlewareFunc
}

func RegisterRoutes(s *server, logger *zap.Logger, cfg *config.Config) *mux.Router {
	var (
		router           = mux.NewRouter()
		authMiddleware   = middleware.NewAuthMiddleware(logger, cfg.JWTSecretKey)
		loggerMiddleware = middleware.NewRequestLoggerMiddleware(logger)
	)

	// Global middlewares
	router.Use(loggerMiddleware)

	// Health endpoint
	router.HandleFunc("/ping", s.Ping)
	// Test Login
	router.HandleFunc("/login", s.TestLogin)

	// Flight endpoints
	flightsAPI := router.PathPrefix("/flights").Subrouter()
	flightsAPI.HandleFunc("/search", s.FlightsSearch)
	flightsAPI.Use(authMiddleware)

	return router
}
