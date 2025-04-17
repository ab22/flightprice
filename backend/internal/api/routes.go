package api

import (
	"net/http"

	"github.com/ab22/flightprice/internal/api/handler"
	"github.com/ab22/flightprice/internal/api/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Route struct {
	Path        string
	HandlerFunc http.HandlerFunc
	Middlewares []http.HandlerFunc
}

var Routes = []Route{
	{
		Path:        "/ping",
		HandlerFunc: handler.Ping,
		Middlewares: nil,
	},
	{
		Path:        "/flights/search",
		HandlerFunc: nil,
		Middlewares: []http.HandlerFunc{},
	},
}

func BuildMuxRouter(logger *zap.Logger) *mux.Router {
	var (
		router           = mux.NewRouter()
		loggerMiddleware = middleware.NewLoggerMiddleware(logger)
	)

	for _, r := range Routes {
		router.HandleFunc(r.Path, r.HandlerFunc)
		router.Use(loggerMiddleware)
	}

	return router
}
