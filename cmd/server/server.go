package main

import (
	"net/http"

	"weather-api/internal/handler"
	"weather-api/internal/weather"
)

func NewMuxServer(weatherClient *weather.Client) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/weather", handler.CurrentWeather(weatherClient))
	return mux
}
