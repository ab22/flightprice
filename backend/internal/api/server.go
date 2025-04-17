package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Server interface {
	Serve() error
	Stop(ctx context.Context) error
}

type server struct {
	router *mux.Router
	srv    *http.Server
	logger *zap.Logger
}

func New(logger *zap.Logger, port int) Server {
	var (
		apiServer = &server{}
		addr      = fmt.Sprintf("0.0.0.0:%d", port)
		router    = RegisterRoutes(apiServer, logger)
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
	return apiServer
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

func (s *server) Stop(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}
