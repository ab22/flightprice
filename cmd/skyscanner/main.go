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
			Price:           102,
			DurationMinutes: 122,
		},
		{
			Price:           152,
			DurationMinutes: 127,
		},
		{
			Price:           177,
			DurationMinutes: 117,
		},
	}
	data, err := json.Marshal(flights)

	if err != nil {
		log.Fatal(err)
	}
	r.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		log.Println("SkyScanner request hit")

		w.WriteHeader(200)
		w.Header().Add("Content-Type", "application/json")

		if _, err := w.Write(data); err != nil {
			log.Println("SkyScanner: failed to write response:", err)
		}
	})

	log.Fatalln(http.ListenAndServe(":8080", r))
}
