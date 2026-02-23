package main

import (
	"log"
	"net/http"

	"weather-api/internal/config"
	"weather-api/internal/weather"
)

func main() {
	cfg := config.MustLoad("config.json")
	log.Println("starting server on", cfg.HTTPPort)

	weatherClient := weather.NewClient(cfg)
	mux := NewMuxServer(weatherClient)

	if err := http.ListenAndServe(cfg.HTTPPort, mux); err != nil {
		log.Fatal(err)
	}
}
