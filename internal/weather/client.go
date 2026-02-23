package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	"weather-api/internal/config"
)

type Client struct {
	httpClient  *http.Client
	redisClient *redis.Client
}

type coords struct {
	Lat float64
	Lon float64
}

var cityCoordinates = map[string]coords{
	"London":   {51.5074, -0.1278},
	"Paris":    {48.8566, 2.3522},
	"Tokyo":    {35.6895, 139.6917},
	"New York": {40.7128, -74.0060},
	"Sydney":   {-33.8688, 151.2093},
	"Rome":     {41.9028, 12.4964},
	"Madrid":   {40.4168, -3.7038},
	"Berlin":   {52.5200, 13.4050},
}

type CurrentWeatherResponse struct {
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		Windspeed   float64 `json:"windspeed"`
		Weathercode int     `json:"weathercode"`
		Time        string  `json:"time"`
	} `json:"current_weather"`
}

const defaultHTTPTimeout = 10 * time.Second
const cacheTTL = 5 * time.Minute

func NewClient(cfg *config.Config) *Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: cfg.RedisURL,
	})
	return &Client{
		httpClient:  &http.Client{Timeout: defaultHTTPTimeout},
		redisClient: rdb,
	}
}

func (c *Client) GetCurrentWeather(ctx context.Context, city string) (CurrentWeatherResponse, error) {
	cityKey := fmt.Sprintf("weather:current:%s", city)

	result, err := c.redisClient.Get(ctx, cityKey).Result()

	if err == nil {
		var cachedResponse CurrentWeatherResponse
		if err := json.Unmarshal([]byte(result), &cachedResponse); err == nil {
			return cachedResponse, nil
		}
	}

	loc, ok := cityCoordinates[city]
	if !ok {
		return CurrentWeatherResponse{}, fmt.Errorf("city not found")
	}
	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current_weather=true",
		loc.Lat,
		loc.Lon,
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return CurrentWeatherResponse{}, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return CurrentWeatherResponse{}, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return CurrentWeatherResponse{}, err
	}

	var currentWeatherResponse CurrentWeatherResponse
	err = json.Unmarshal(body, &currentWeatherResponse)
	if err != nil {
		return CurrentWeatherResponse{}, err
	}

	if c.redisClient != nil {
		if data, err := json.Marshal(currentWeatherResponse); err == nil {
			_ = c.redisClient.Set(ctx, cityKey, data, cacheTTL).Err()
		}

	}
	return currentWeatherResponse, nil
}
