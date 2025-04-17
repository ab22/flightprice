package middleware

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

func NewAuthMiddleware(logger *zap.Logger, secretKey string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("X-API-Token")
			if len(tokenString) == 0 {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
				return []byte(secretKey), nil
			})

			if err != nil {
				logger.Error("failed to jwt parse", zap.Error(err))
				w.WriteHeader(http.StatusNotFound)
				return
			} else if !token.Valid {
				logger.Error("invalid token", zap.Error(err))
				w.WriteHeader(http.StatusNotFound)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
