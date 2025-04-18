package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ab22/flightprice/internal/config"
	"github.com/ab22/flightprice/internal/service"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Server interface {
	Serve() error
	Stop(ctx context.Context) error
}

type server struct {
	cfg            *config.Config
	router         *mux.Router
	srv            *http.Server
	wsUpgrader     websocket.Upgrader
	logger         *zap.Logger
	flightsService service.FlightsService
}

func New(logger *zap.Logger, cfg *config.Config, flightsService service.FlightsService) Server {
	server := &server{
		logger:         logger,
		cfg:            cfg,
		wsUpgrader:     websocket.Upgrader{},
		flightsService: flightsService,
	}
	server.router = server.RegisterRoutes()
	server.srv = &http.Server{
		Handler:      server.router,
		Addr:         fmt.Sprintf("0.0.0.0:%s", cfg.APIPort),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	return server
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
