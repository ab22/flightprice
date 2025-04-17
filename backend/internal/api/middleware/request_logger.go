package middleware

import (
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewLoggerMiddleware(logger *zap.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)

			logger.Info("request hit",
				zap.String("request_method", r.Method),
				zap.String("request_path", r.RequestURI),
			)
		})
	}
}
