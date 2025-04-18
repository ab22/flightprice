package api

import (
	"net/http"

	"go.uber.org/zap"
)

func (s *server) Ping(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("Ping handler hit")
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write([]byte("pong")); err != nil {
		s.logger.Error("failed to write response", zap.Error(err))
	}
}
