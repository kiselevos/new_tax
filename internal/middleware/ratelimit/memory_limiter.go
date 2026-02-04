package ratelimit

import (
	"context"
	"sync"

	"golang.org/x/time/rate"
)

type memoryLimiter struct {
	mu       sync.Mutex
	limiters map[string]*rate.Limiter
}

func NewMemoryLimiter() *memoryLimiter {
	return &memoryLimiter{
		limiters: make(map[string]*rate.Limiter),
	}
}

func (m *memoryLimiter) Allow(ctx context.Context, key string, rps float64, burst int) (bool, error) {

	m.mu.Lock()
	defer m.mu.Unlock()

	limiter, ok := m.limiters[key]
	if !ok {
		limiter = rate.NewLimiter(rate.Limit(rps), burst)
		m.limiters[key] = limiter
	}

	return limiter.Allow(), nil
}
