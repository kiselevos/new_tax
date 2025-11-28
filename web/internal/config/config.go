package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	WebPort  string
	Backend  string
	APIRPS   int
	APIBurst int
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

	rps, err := envInt("API_RPS")
	if err != nil {
		return nil, fmt.Errorf("invalid API_RPS: %w", err)
	}

	burst, err := envInt("API_BURST")
	if err != nil {
		return nil, fmt.Errorf("invalid API_BURST: %w", err)
	}

	return &Config{
		WebPort:  webPort,
		Backend:  backAddr,
		APIRPS:   rps,
		APIBurst: burst,
	}, nil
}

func envInt(key string) (int, error) {
	val := os.Getenv(key)
	if val == "" {
		return 0, fmt.Errorf("%s is not set", key)
	}

	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return 0, fmt.Errorf("env %s must be > 0, got `%s`", key, val)
	}

	return n, nil
}
