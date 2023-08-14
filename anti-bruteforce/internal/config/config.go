package config

import (
	"fmt"
	"os"
)

type Config struct {
	Logger struct {
		Level string `json:"level"`
		File  string `json:"file"`
	} `json:"logger"`
	HTTP struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"http"`
	CLI bool `json:"cli"`
	DB  struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		SSLMode  string `json:"sslMode"`
	} `json:"db"`
	Limitations struct {
		Login    int `json:"login"`
		Password int `json:"password"`
		IP       int `json:"ip"`
	} `json:"limitations"`
}

func Parse(filePath string) (*Config, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading the configuration file: %w", err)
	}
	cfg := new(Config)
	err = cfg.UnmarshalJSON(bytes)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshalling the configuration file:: %w", err)
	}
	return cfg, err
}
