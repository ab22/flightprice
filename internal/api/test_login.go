package api

import (
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

func (s *server) TestLogin(w http.ResponseWriter, r *http.Request) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"username": "test-username",
			// Expiration can be set here:
			// "exp": time.Now().Add(time.Hour * 24).Unix()
		})

	tokenString, err := token.SignedString([]byte(s.cfg.JWTSecretKey))

	if err != nil {
		s.logger.Error("failed to sign token", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	w.Write([]byte(tokenString))
}
