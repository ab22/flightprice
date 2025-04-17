package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const APIToken = "mysecrettoken"

func NewAuthMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestAPIToken := r.Header.Get("X-API-Token")

			if requestAPIToken != APIToken {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
