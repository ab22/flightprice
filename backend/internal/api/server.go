package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ab22/flightprice/internal/config"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server interface {
	Serve() error
	Stop(ctx context.Context) error
}

type server struct {
	cfg    *config.Config
	router *mux.Router
	srv    *http.Server
	logger *zap.Logger
}

func New(logger *zap.Logger, cfg *config.Config) Server {
	var (
		apiServer = &server{}
		addr      = fmt.Sprintf("0.0.0.0:%s", cfg.APIPort)
		router    = RegisterRoutes(apiServer, logger, cfg)
		srv       = &http.Server{
			Handler:      router,
			Addr:         addr,
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}
	)

	apiServer.router = router
	apiServer.srv = srv
	apiServer.logger = logger
	apiServer.cfg = cfg
	return apiServer
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
