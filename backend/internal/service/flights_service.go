package service

import (
	"sync"

	"github.com/ab22/flightprice/internal/client"
	"github.com/ab22/flightprice/internal/config"
	"github.com/ab22/flightprice/internal/service/models"
	"go.uber.org/zap"
)

type FlightsService interface {
	SearchFlights() (*models.SearchFlightsOut, error)
}

type flightsService struct {
	amadeusClient       client.ThirdPartyAPI
	googleflightsClient client.ThirdPartyAPI
	skyscannerClient    client.ThirdPartyAPI
	logger              *zap.Logger
	cfg                 *config.Config
}

func NewFlightsService(
	amadeusClient client.ThirdPartyAPI,
	googleflightsClient client.ThirdPartyAPI,
	skyscannerClient client.ThirdPartyAPI,
	logger *zap.Logger,
	cfg *config.Config) FlightsService {
	return &flightsService{
		amadeusClient,
		googleflightsClient,
		skyscannerClient,
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

func (s *flightsService) SearchFlights() (*models.SearchFlightsOut, error) {
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
	}, nil
}
