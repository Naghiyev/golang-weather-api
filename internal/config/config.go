package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	HTTPPort string `json:"server_port"`
	RedisURL string `json:"redis_url"`
}

func LoadConfig(path string) (*Config, error) {
	res, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config Config
	if err = json.Unmarshal(res, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func MustLoad(path string) *Config {
	cfg, err := LoadConfig(path)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}
