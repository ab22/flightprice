package api

import "net/http"

func (s *server) FlightsSearch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
