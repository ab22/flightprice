package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ab22/flightprice/internal/client/models"
)

var (
	InvalidStatusCodeErr = errors.New("an invalid status code was recieved")
)

type ThirdPartyAPI interface {
	FetchFlights() ([]models.Flight, error)
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

func (a *thirdPartyAPI) FetchFlights() ([]models.Flight, error) {
	res, err := a.httpClient.Get(a.address)

	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: http request failed: %w", err)
	}
	defer res.Body.Close()

	var flights []models.Flight
	data, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: io.ReadAll failed:", err)
	} else if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: status code [%d] - data received [%s]", InvalidStatusCodeErr, res.StatusCode, string(data))
	}

	err = json.Unmarshal(data, &flights)

	if err != nil {
		return nil, fmt.Errorf("ThirdPartyAPI.FetchFlights: json.Unmarshal failed: %w", err)
	}

	return flights, nil
}
