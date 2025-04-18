package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ab22/flightprice/internal/api"
	"github.com/ab22/flightprice/internal/client"
	"github.com/ab22/flightprice/internal/config"
	"github.com/ab22/flightprice/internal/service"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Use a colorized logger for local environments. Ideally, we
// should detect which env we are running on and build a different one
// for non-local environments.
func NewLocalLogger() (*zap.Logger, error) {
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.Level.SetLevel(zapcore.DebugLevel)

	return config.Build()
}

func main() {
	logger, err := NewLocalLogger()
	if err != nil {
		log.Fatalln("Failed to create logger:", err)
	}

	defer func() {
		err := logger.Sync()

		// This will probably not be flushed.
		if err != nil {
			log.Println("Failed to sync logger:", err)
		}
	}()
	logger.Info("Starting API on port 8080")
	cfg, err := config.New()

	if err != nil {
		logger.Fatal("Failed to parse configuration", zap.Error(err))
	}

	var (
		httpClient          = &http.Client{}
		amadeusClient       = client.NewThirdPartyAPI(httpClient, "http://amadeus:8080/search", logger)
		googleflightsClient = client.NewThirdPartyAPI(httpClient, "http://googleflights:8080/search", logger)
		skyscannerClient    = client.NewThirdPartyAPI(httpClient, "http://skyscanner:8080/search", logger)
		redisClient         = redis.NewClient(&redis.Options{
			Addr:     "redis:6379",
			Password: "", // No password set
			DB:       0,  // use default DB
		})
		flightsService = service.NewFlightsService(
			amadeusClient,
			googleflightsClient,
			skyscannerClient,
			redisClient,
			logger,
			&cfg)
		api = api.New(logger, &cfg, flightsService)
	)

	// Serve on a separate goroutine.
	go func() {
		if err := api.Serve(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Fatalln("serve error:", err)
			}
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
		logger.Fatal("Failed to stop server", zap.Error(err))
	}

	logger.Info("HTTP server stopped successfully")
}
