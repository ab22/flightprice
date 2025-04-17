package models

type SearchFlightsOut struct {
	Cheapest *Flight `json:"cheapest_flights"`
	Fastest  *Flight `json:"fastest_flights"`
}
