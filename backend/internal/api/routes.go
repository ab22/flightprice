package api

import (
	"net/http"

	"github.com/ab22/flightprice/internal/api/handler"
	"github.com/gorilla/mux"
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

func BuildMuxRouter() *mux.Router {
	router := mux.NewRouter()

	for _, r := range Routes {
		router.HandleFunc(r.Path, r.HandlerFunc)
	}

	return router
}
