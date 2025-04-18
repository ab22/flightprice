package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ab22/flightprice/internal/client/models"
)

var (
	ErrInvalidStatusCode = errors.New("an invalid status code was recieved")
)

type ThirdPartyAPI interface {
	FetchFlights(ctx context.Context) ([]models.Flight, error)
}

type thirdPartyAPI struct {
	httpClient *http.Client
	address    string
}

func NewThirdPartyAPI(httpClient *http.Client, address string) ThirdPartyAPI {
	return &thirdPartyAPI{
		httpClient,
		address,
	}
}

func (a *thirdPartyAPI) FetchFlights(ctx context.Context) ([]models.Flight, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.address, nil)
	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: create request failed: %w", err)
	}

	res, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: http request failed: %w", err)
	}
	defer res.Body.Close()

	var flights []models.Flight
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: io.ReadAll failed: %w", err)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code [%d] - data received [%s]", ErrInvalidStatusCode, res.StatusCode, string(data))
	}

	err = json.Unmarshal(data, &flights)

	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: json.Unmarshal failed: %w", err)
	}

	return flights, nil
}
