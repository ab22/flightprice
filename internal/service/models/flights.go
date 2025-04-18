package models

type Flight struct {
	Service         string `json:"service"`
	Price           int    `json:"price"`
	DurationMinutes int    `json:"duration_minutes"`
}
