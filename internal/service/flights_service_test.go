package service

import (
	"testing"
	"time"

	"github.com/ab22/flightprice/internal/service/models"
	"github.com/stretchr/testify/require"
)

func NewTestService() flightsService {
	return flightsService{}
}

func TestFindCheapestFlight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		flights        []models.Flight
		expectedFlight *models.Flight
	}{
		{
			name:           "when_flights_array_is_nil_it_should_send_nil_on_channel",
			flights:        nil,
			expectedFlight: nil,
		},
		{
			name:           "when_flights_array_is_empty_it_should_return_nil_on_channel",
			flights:        []models.Flight{},
			expectedFlight: nil,
		},
		{
			name: "when_flights_array_has_entries_with_same_price_it_should_return_first_found",
			flights: []models.Flight{
				models.Flight{
					Service: "service1",
					Price:   100,
				},
				models.Flight{
					Service: "service2",
					Price:   100,
				},
			},
			expectedFlight: &models.Flight{
				Service: "service1",
				Price:   100,
			},
		},
		{
			name: "when_multiple_flights_are_passed_it_finds_the_cheapest_flight",
			flights: []models.Flight{
				models.Flight{
					Service: "service1",
					Price:   100,
				},
				models.Flight{
					Service: "service2",
					Price:   101,
				},
				models.Flight{
					Service: "service3",
					Price:   102,
				},
				models.Flight{
					Service: "service4",
					Price:   10,
				},
			},
			expectedFlight: &models.Flight{
				Service: "service4",
				Price:   10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s = NewTestService()
				c = make(chan *models.Flight, 1)
			)

			s.findCheapestFlight(c, tt.flights)

			select {
			case <-time.After(1 * time.Second):
				t.Fatalf("test timeout")
			case actual := <-c:
				require.Equal(t, tt.expectedFlight, actual)
				break
			}
		})
	}
}

func TestFindFastestFlight(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		flights        []models.Flight
		expectedFlight *models.Flight
	}{
		{
			name:           "when_flights_array_is_nil_it_should_send_nil_on_channel",
			flights:        nil,
			expectedFlight: nil,
		},
		{
			name:           "when_flights_array_is_empty_it_should_return_nil_on_channel",
			flights:        []models.Flight{},
			expectedFlight: nil,
		},
		{
			name: "when_flights_array_has_entries_with_same_duration_it_should_return_first_found",
			flights: []models.Flight{
				models.Flight{
					Service:         "service1",
					DurationMinutes: 100,
				},
				models.Flight{
					Service:         "service2",
					DurationMinutes: 100,
				},
			},
			expectedFlight: &models.Flight{
				Service:         "service1",
				DurationMinutes: 100,
			},
		},
		{
			name: "when_multiple_flights_are_passed_it_finds_the_fastest_flight",
			flights: []models.Flight{
				models.Flight{
					Service:         "service1",
					DurationMinutes: 100,
				},
				models.Flight{
					Service:         "service2",
					DurationMinutes: 101,
				},
				models.Flight{
					Service:         "service3",
					DurationMinutes: 102,
				},
				models.Flight{
					Service:         "service4",
					DurationMinutes: 10,
				},
			},
			expectedFlight: &models.Flight{
				Service:         "service4",
				DurationMinutes: 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var (
				s = NewTestService()
				c = make(chan *models.Flight, 1)
			)

			s.findFastestFlight(c, tt.flights)

			select {
			case <-time.After(1 * time.Second):
				t.Fatalf("test timeout")
			case actual := <-c:
				require.Equal(t, tt.expectedFlight, actual)
				break
			}
		})
	}
}
