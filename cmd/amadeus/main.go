package main

import (
	_ "embed"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Flights struct {
	Price           int `json:"price"`
	DurationMinutes int `json:"duration_minutes"`
}

func main() {
	r := mux.NewRouter()
	flights := []Flights{
		{
			Price:           100,
			DurationMinutes: 120,
		},
		{
			Price:           150,
			DurationMinutes: 125,
		},
		{
			Price:           175,
			DurationMinutes: 115,
		},
	}
	data, err := json.Marshal(flights)

	if err != nil {
		log.Fatal(err)
	}

	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Amadeus request hit")

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")

		if _, err := w.Write(data); err != nil {
			log.Println("Amadeus: failed to write response:", err)
		}
	})

	log.Fatalln(http.ListenAndServe(":8080", r))
}
