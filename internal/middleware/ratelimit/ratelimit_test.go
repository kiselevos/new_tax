package ratelimit

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"
	"time"
)

type fakeLimiter struct {
	allow bool
	err   error
	calls int64
}

func (f *fakeLimiter) Allow(ctx context.Context, key string, rps float64, burst int) (bool, error) {
	atomic.AddInt64(&f.calls, 1)
	return f.allow, f.err
}

func TestHybridLimiter_FallbackOnError(t *testing.T) {
	primary := &fakeLimiter{allow: false, err: fmt.Errorf("redis down")}
	fallback := &fakeLimiter{allow: true, err: nil}

	h := NewHybridLimiter(primary, fallback)

	ok, err := h.Allow(context.Background(), "k", 1, 1)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	if !ok {
		t.Fatalf("expected allowed via fallback")
	}

	if atomic.LoadInt64(&primary.calls) != 1 {
		t.Fatalf("expected primary called once, got %d", atomic.LoadInt64(&primary.calls))
	}
	if atomic.LoadInt64(&fallback.calls) != 1 {
		t.Fatalf("expected fallback called once, got %d", atomic.LoadInt64(&fallback.calls))
	}
}

func TestHybridLimiter_CooldownSkipsPrimary(t *testing.T) {
	primary := &fakeLimiter{allow: true, err: fmt.Errorf("redis down")}
	fallback := &fakeLimiter{allow: true, err: nil}

	h := NewHybridLimiter(primary, fallback)
	h.cooldown = 200 * time.Millisecond

	_, _ = h.Allow(context.Background(), "k", 1, 1)

	_, _ = h.Allow(context.Background(), "k", 1, 1)

	if atomic.LoadInt64(&primary.calls) != 1 {
		t.Fatalf("expected primary called only once during cooldown, got %d", atomic.LoadInt64(&primary.calls))
	}
	if atomic.LoadInt64(&fallback.calls) != 2 {
		t.Fatalf("expected fallback called twice, got %d", atomic.LoadInt64(&fallback.calls))
	}

	time.Sleep(250 * time.Millisecond)
	_, _ = h.Allow(context.Background(), "k", 1, 1)

	if atomic.LoadInt64(&primary.calls) != 2 {
		t.Fatalf("expected primary called again after cooldown, got %d", atomic.LoadInt64(&primary.calls))
	}
}
