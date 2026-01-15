package redisconn

import (
	"context"
	"log/slog"
	"time"

	"github.com/kiselevos/new_tax/internal/config"
	"github.com/redis/go-redis/v9"
)

func Connect(cfg *config.RedisConfig, logger *slog.Logger) (*redis.Client, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:         cfg.Addr,
		DialTimeout:  200 * time.Millisecond,
		ReadTimeout:  200 * time.Millisecond,
		WriteTimeout: 200 * time.Millisecond,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		logger.Warn("redis_unavailable", "err", err)
		if cerr := rdb.Close(); cerr != nil {
			logger.Warn("redis_close_failed", "err", cerr)
		}
		return nil, err
	}

	logger.Info("redis_connected", "addr", cfg.Addr)
	return rdb, nil
}
