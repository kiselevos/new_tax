package ratelimit

import (
	"context"
	"sync/atomic"
	"time"
)

type HybridLimiter struct {
	primary  Limiter // Redis
	fallback Limiter // Memory

	cooldown  time.Duration
	failUntil atomic.Int64
}

func NewHybridLimiter(primary Limiter, fallback Limiter) *HybridLimiter {
	return &HybridLimiter{
		primary:  primary,
		fallback: fallback,
		cooldown: 2 * time.Second,
	}
}

func (hl *HybridLimiter) Allow(ctx context.Context, key string, rps float64, burst int) (bool, error) {

	until := time.Unix(0, hl.failUntil.Load())
	if time.Now().Before(until) {
		return hl.fallback.Allow(ctx, key, rps, burst)
	}

	ok, err := hl.primary.Allow(ctx, key, rps, burst)
	if err == nil {
		return ok, nil
	}

	hl.failUntil.Store(time.Now().Add(hl.cooldown).UnixNano())
	return hl.fallback.Allow(ctx, key, rps, burst)
}
