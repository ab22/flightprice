package config

import "github.com/ab22/env"

type Config struct {
	JWTSecretKey string `env:"JWT_SECRET_KEY"`
	APIPort      string `env:"API_PORT"`
}

func New() (Config, error) {
	c := Config{}
	err := env.Parse(&c)

	return c, err
}
