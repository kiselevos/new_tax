package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type memoryLimiter struct {
	mu    sync.Mutex
	items map[string]*memoryItem

	ttl          time.Duration
	cleanupEvery int
	calls        int
}

type memoryItem struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func NewMemoryLimiter(ttl time.Duration, cleanupEvery int) *memoryLimiter {
	if cleanupEvery <= 0 {
		cleanupEvery = 1000
	}
	return &memoryLimiter{
		items:        make(map[string]*memoryItem),
		ttl:          ttl,
		cleanupEvery: cleanupEvery,
	}
}

func (m *memoryLimiter) Allow(ctx context.Context, key string, rps float64, burst int) (bool, error) {

	now := time.Now()

	m.mu.Lock()

	it := m.items[key]

	if it != nil && m.ttl > 0 && now.Sub(it.lastSeen) > m.ttl {
		delete(m.items, key)
		it = nil
	}

	if it == nil {
		it = &memoryItem{
			limiter:  rate.NewLimiter(rate.Limit(rps), burst),
			lastSeen: now,
		}
		m.items[key] = it
	} else {
		it.lastSeen = now
	}

	m.calls++
	doCleanup := m.calls%m.cleanupEvery == 0

	lim := it.limiter

	m.mu.Unlock()

	allowed := lim.Allow()

	if doCleanup {
		m.cleanup(now)
	}

	return allowed, nil
}

func (m *memoryLimiter) cleanup(now time.Time) {

	cutoff := now.Add(-m.ttl)

	m.mu.Lock()
	for k, it := range m.items {
		if it.lastSeen.Before(cutoff) {
			delete(m.items, k)
		}
	}
	m.mu.Unlock()
}
