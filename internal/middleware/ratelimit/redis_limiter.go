package ratelimit

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisLimiter struct {
	rdb    *redis.Client
	script *redis.Script
	ttl    time.Duration
}

func NewRedisLimiter(rdb *redis.Client, ttl time.Duration) *redisLimiter {
	return &redisLimiter{
		rdb:    rdb,
		script: redis.NewScript(redisTokenBucketLua),
		ttl:    ttl,
	}
}

func (rl *redisLimiter) Allow(ctx context.Context, key string, rps float64, burst int) (bool, error) {

	if rl.rdb == nil {
		return false, fmt.Errorf("redis is nil")
	}

	if rps <= 0 || burst <= 0 {
		return false, nil
	}

	redisKey := "rl:{" + key + "}"

	res, err := rl.script.Run(ctx, rl.rdb, []string{redisKey}, rps, burst, rl.ttl.Microseconds()).Result()
	if err != nil {
		return false, err
	}

	arr, ok := res.([]interface{})
	if !ok || len(arr) < 1 {
		return false, fmt.Errorf("unexpected redis script response: %T %v", res, res)
	}

	allowed, ok := arr[0].(int64)
	if !ok {
		return false, fmt.Errorf("unexpected allowed type: %T", arr[0])
	}

	return allowed == 1, nil
}
