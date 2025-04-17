package api

import "net/http"

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))

	s.logger.Info("Ping handler hit")
}
