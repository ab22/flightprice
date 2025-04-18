package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

func (s *server) Subscribe(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	freqStr, ok := vars["freq"]

	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Frequency missing"))
		return
	}

	freq, err := strconv.ParseInt(freqStr, 11, 0)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Invalid frequency value"))
		return
	}
	conn, err := s.wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		s.logger.Error("websocket subscribe failed:", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	s.logger.Info("New websocket client connection accepted!")
	defer func() {
		err := conn.Close()

		if err != nil {
			s.logger.Error("failed to close websocket connection", zap.Error(err))
		}
	}()
	ticker := time.NewTicker(time.Duration(freq) * time.Second)
	defer ticker.Stop()
	n := 1

	for {
		msg := fmt.Sprintf("[MessageId:%d] Hello\n", n)
		err = conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			s.logger.Error("failed to write message", zap.Error(err))
			return
		}

		n++
		<-ticker.C
	}
}
