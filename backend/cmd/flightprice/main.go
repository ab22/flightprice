package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ab22/flightprice/internal/api"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Use a colorized logger for local environments. Ideally, we
// should detect which env we are running on and build a different one
// for non-local environments.
func NewLocalLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	return config.Build()
}

func main() {
	logger, err := NewLocalLogger()
	if err != nil {
		log.Fatalln("Failed to create logger:", err)
	}

	logger.Info("Starting API on port 8080")
	port := 8080
	api := api.New(logger, port)

	// Serve on a separate goroutine.
	go func() {
		if err := api.Serve(); err != nil {
			log.Fatalln("serve error:", err)
		}
	}()

	// Listen for termination signals
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)
	s := <-sigc

	logger.Info("Received termination signal", zap.String("signal", s.String()))
	logger.Info("Stopping HTTP server with a hard timeout of 30 seconds")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := api.Stop(ctx); err != nil {
		logger.Error("Failed to stop server", zap.Error(err))
	} else {
		logger.Info("HTTP server stopped successfully")
	}
}
