package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"weather-api/internal/weather"
)

type errorResponse struct {
	Error string `json:"error"`
}

func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func CurrentWeather(client *weather.Client) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		city := request.URL.Query().Get("city")
		if city == "" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(writer).Encode(errorResponse{Error: "city is required"})
			return
		}

		current, err := client.GetCurrentWeather(request.Context(), city)
		if err != nil {
			log.Printf("Failed to get weather: %v", err)
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusInternalServerError)
			_ = json.NewEncoder(writer).Encode(errorResponse{Error: "could not fetch weather right now"})
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(writer).Encode(current)
	}
}
