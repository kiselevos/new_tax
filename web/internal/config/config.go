package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	WebPort    string
	Backend    string
	APIVersion string
	LogMode    string
	LogLevel   string
}

func Load() (*Config, error) {

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		return nil, fmt.Errorf("WEB_PORT is required")
	}
	if strings.Contains(webPort, ":") {
		return nil, fmt.Errorf("WEB_PORT must contain only port (example: 8080)")
	}

	backAddr := os.Getenv("BACKEND_ADDR")
	if backAddr == "" {
		return nil, fmt.Errorf("BACKEND_ADDR is required")
	}
	if !strings.Contains(backAddr, ":") {
		return nil, fmt.Errorf("BACKEND_ADDR must contain port (example: tax-backend:50051)")
	}

	apiVers := os.Getenv("API_VERSION")
	if webPort == "" {
		apiVers = "v1"
	}

	logMode := os.Getenv("LOG_MODE")
	if logMode == "" {
		logMode = "json"
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}

	return &Config{
		WebPort:    webPort,
		Backend:    backAddr,
		APIVersion: apiVers,
		LogMode:    logMode,
		LogLevel:   logLevel,
	}, nil
}
