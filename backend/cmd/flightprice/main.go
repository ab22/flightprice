package main

import (
	"log"

	"github.com/ab22/flightprice/internal/api"
)

func main() {
	api := api.New()

	if err := api.Serve(); err != nil {
		log.Fatalln("serve error:", err)
	}
}
