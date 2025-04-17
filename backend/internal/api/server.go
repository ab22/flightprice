package api

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server interface {
	Serve() error
	Stop()
}

type server struct {
	router *mux.Router
	srv    *http.Server
}

func New() Server {
	var (
		router = BuildMuxRouter()
		srv    = &http.Server{
			Handler:      router,
			Addr:         "0.0.0.0:8080",
			WriteTimeout: 10 * time.Second,
			ReadTimeout:  10 * time.Second,
		}
	)

	return &server{
		router,
		srv,
	}
}

func (s *server) Serve() error {
	return s.srv.ListenAndServe()
}

func (s *server) Stop() {
}
