package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/ab22/flightprice/internal/client"
	"github.com/ab22/flightprice/internal/config"
	"github.com/ab22/flightprice/internal/service/models"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type FlightsService interface {
	SearchFlights(ctx context.Context) (*models.SearchFlightsOut, error)
}

type flightsService struct {
	amadeusClient       client.ThirdPartyAPI
	googleflightsClient client.ThirdPartyAPI
	skyscannerClient    client.ThirdPartyAPI
	redisClient         *redis.Client
	logger              *zap.Logger
	cfg                 *config.Config
}

func NewFlightsService(
	amadeusClient client.ThirdPartyAPI,
	googleflightsClient client.ThirdPartyAPI,
	skyscannerClient client.ThirdPartyAPI,
	redisClient *redis.Client,
	logger *zap.Logger,
	cfg *config.Config) FlightsService {
	return &flightsService{
		amadeusClient,
		googleflightsClient,
		skyscannerClient,
		redisClient,
		logger,
		cfg,
	}
}

func (s *flightsService) launchRequest(c client.ThirdPartyAPI, service string, wg *sync.WaitGroup, mu *sync.Mutex, res *[]models.Flight) {
	defer wg.Done()
	flights, err := c.FetchFlights()
	if err != nil {
		s.logger.Error("failed to get flights from "+service, zap.Error(err))
	}

	mu.Lock()
	defer mu.Unlock()
	for _, f := range flights {
		*res = append(*res, models.Flight{
			Service:         service,
			Price:           f.Price,
			DurationMinutes: f.DurationMinutes,
		})
	}
}

func (s *flightsService) findCheapestFlight(c chan *models.Flight, flights []models.Flight) {
	var cheapest *models.Flight

	for _, f := range flights {
		if cheapest == nil {
			cheapest = &f
		} else if cheapest.Price > f.Price {
			cheapest = &f
		}
	}

	c <- cheapest
}

func (s *flightsService) findFastestFlight(c chan *models.Flight, flights []models.Flight) {
	var fastest *models.Flight

	for _, f := range flights {
		if fastest == nil {
			fastest = &f
		} else if fastest.DurationMinutes > f.DurationMinutes {
			fastest = &f
		}
	}

	c <- fastest
}

func (s *flightsService) searchCached(ctx context.Context) (*models.SearchFlightsOut, error) {
	data, err := s.redisClient.Get(ctx, "flights").Bytes()

	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}

		return nil, fmt.Errorf("FlightsService.searchCached: redis get failed: %w", err)
	}

	var flights *models.SearchFlightsOut
	err = json.Unmarshal(data, &flights)

	if err != nil {
		return nil, fmt.Errorf("FlightsService.searchCached: json unmarshalling failed: %w", err)
	}

	return flights, err
}

func (s *flightsService) fetchFlightsFromClients(_ context.Context) *models.SearchFlightsOut {
	var wg sync.WaitGroup
	var mu sync.Mutex
	var flights = make([]models.Flight, 0, 20)

	wg.Add(3)
	go s.launchRequest(s.amadeusClient, "Amadeus", &wg, &mu, &flights)
	go s.launchRequest(s.googleflightsClient, "GoogleFlights", &wg, &mu, &flights)
	go s.launchRequest(s.skyscannerClient, "SkyScanner", &wg, &mu, &flights)
	wg.Wait()

	cheapestFlight := make(chan *models.Flight, 1)
	fastestFlight := make(chan *models.Flight, 1)

	go s.findCheapestFlight(cheapestFlight, flights)
	go s.findFastestFlight(fastestFlight, flights)

	return &models.SearchFlightsOut{
		Cheapest: <-cheapestFlight,
		Fastest:  <-fastestFlight,
	}
}

func (s *flightsService) SearchFlights(ctx context.Context) (*models.SearchFlightsOut, error) {
	cachedFlights, err := s.searchCached(ctx)

	if err != nil {
		return nil, fmt.Errorf("FlightsService.SearchFlights: %w", err)
	} else if cachedFlights != nil {
		s.logger.Debug("CACHE HIT")
		return cachedFlights, nil
	}

	s.logger.Debug("CACHE MISS")
	flights := s.fetchFlightsFromClients(ctx)
	data, err := json.Marshal(flights)

	if err != nil {
		return nil, fmt.Errorf("FlightsService.SearchFlights: json unmarshalling failed: %w", err)
	}

	s.redisClient.Set(ctx, "flights", data, 30*time.Second)
	return s.fetchFlightsFromClients(ctx), nil
}
