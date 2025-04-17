package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

func (s *server) FlightsSearch(w http.ResponseWriter, r *http.Request) {
	flights, err := s.flightsService.SearchFlights()

	if err != nil {
		s.logger.Error("failed to get flights", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(flights)

	if err != nil {
		s.logger.Error("failed to marshal flights", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
