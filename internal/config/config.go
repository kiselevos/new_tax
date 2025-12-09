package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	BackPort string
	ApiKey   string
	LogMode  string
	LogLevel string
}

func Load() (*Config, error) {
	backPort := os.Getenv("BACKEND_PORT")
	if backPort == "" {
		return nil, fmt.Errorf("BACKEND_PORT is required")
	}
	if strings.Contains(backPort, ":") {
		return nil, fmt.Errorf("BACKEND_PORT must contain only port (example: 50051)")
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("API_KEY is required")
	}

	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "json"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	conf := &Config{
		BackPort: ":" + backPort,
		ApiKey:   apiKey,
		LogMode:  logMode,
		LogLevel: logLevel,
	}

	return conf, nil
}
