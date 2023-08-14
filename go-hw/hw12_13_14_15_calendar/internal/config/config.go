package config

import "os"

type Config struct {
	Logger struct {
		Level string `json:"level"`
		File  string `json:"file"`
	} `json:"logger"`
	HTTP struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"http"`
	GRPC struct {
		Host string `json:"host"`
		Port string `json:"port"`
	} `json:"grpc"`
	SQL bool `json:"sql"`
	DB  struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		Database string `json:"database"`
		SSLMode  string `json:"sslMode"`
	} `json:"db"`
	AMQP struct {
		URI   string `json:"uri"`
		Queue string `json:"queue"`
	} `json:"amqp"`
}

func Parse(filePath string) (*Config, error) {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = cfg.UnmarshalJSON(bytes)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
