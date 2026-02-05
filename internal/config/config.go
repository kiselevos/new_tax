package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	BackPort string
	ApiKey   string
	LogMode  string
	LogLevel string

	RateLimitCfg *RateLimitConfig

	RedisCfg *RedisConfig
}

type RateLimitConfig struct {
	PublicRPS     float64
	PublicBurst   int
	PrivateRPS    float64
	PrivateBurst  int
	TTL           time.Duration
	Cleanup_Every int
}

type RedisConfig struct {
	Enabled bool
	Addr    string
	Ttl     time.Duration
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
		BackPort:     ":" + backPort,
		ApiKey:       apiKey,
		LogMode:      logMode,
		LogLevel:     logLevel,
		RateLimitCfg: LoadRateLimitConf(),
		RedisCfg:     LoadRedisConfig(),
	}

	return conf, nil
}

func LoadRateLimitConf() *RateLimitConfig {

	publicRPS, err := strconv.ParseFloat(os.Getenv("RATE_LIMIT_PUBLIC_RPS"), 64)
	if err != nil {
		publicRPS = 1
	}

	publicBurst, err := strconv.Atoi(os.Getenv("RATE_LIMIT_PUBLIC_BURST"))
	if err != nil {
		publicBurst = 10
	}

	privateRPS, err := strconv.ParseFloat(os.Getenv("RATE_LIMIT_PRIVATE_RPS"), 64)
	if err != nil {
		privateRPS = 2
	}

	privateBurst, err := strconv.Atoi(os.Getenv("RATE_LIMIT_PRIVATE_BURST"))
	if err != nil {
		privateBurst = 20
	}

	ttl := os.Getenv("TTL")
	rlTTL := 10 * time.Minute

	if ttl != "" {
		if parse, err := time.ParseDuration(ttl); err == nil {
			rlTTL = parse
		}
	}

	cleanup, err := strconv.Atoi(os.Getenv("CLEANUP_EVERY"))
	if err != nil {
		cleanup = 1000
	}

	return &RateLimitConfig{
		PublicRPS:     publicRPS,
		PublicBurst:   publicBurst,
		PrivateRPS:    privateRPS,
		PrivateBurst:  privateBurst,
		Cleanup_Every: cleanup,
		TTL:           rlTTL,
	}
}

func LoadRedisConfig() *RedisConfig {
	enabledStr := strings.TrimSpace(strings.ToLower(os.Getenv("REDIS_ENABLED")))
	enabled := enabledStr == "1" || enabledStr == "true" || enabledStr == "yes"

	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "redis:6379"
	}

	ttl := os.Getenv("TTL")
	cacheTTL := 10 * time.Minute

	if ttl != "" {
		if parse, err := time.ParseDuration(ttl); err == nil {
			cacheTTL = parse
		}
	}

	return &RedisConfig{
		Enabled: enabled,
		Addr:    addr,
		Ttl:     cacheTTL,
	}
}
