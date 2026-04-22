package ratelimit

import (
	"context"
	"sync"
	"testing"
	"time"
)


func TestMemoryLimiter_Burst(t *testing.T) {
	lim := NewMemoryLimiter(1*time.Minute, 100)
	key := "ip1"

	rps := 5.0
	burst := 3

	for i := 0; i < burst; i++ {
		ok, err := lim.Allow(context.Background(), key, rps, burst)
		if err != nil || !ok {
			t.Fatalf("expected allow at %d, ok=%v err=%v", i, ok, err)
		}
	}

	ok, _ := lim.Allow(context.Background(), key, rps, burst)
	if ok {
		t.Fatalf("expected deny after burst exhausted")
	}
}

func TestMemoryLimiter_Refill(t *testing.T) {
	lim := NewMemoryLimiter(1*time.Minute, 100)
	key := "ip2"

	rps := 5.0
	burst := 1

	ok, _ := lim.Allow(context.Background(), key, rps, burst)
	if !ok {
		t.Fatalf("first request must pass")
	}

	ok, _ = lim.Allow(context.Background(), key, rps, burst)
	if ok {
		t.Fatalf("should be denied immediately after")
	}

	time.Sleep(220 * time.Millisecond)

	ok, _ = lim.Allow(context.Background(), key, rps, burst)
	if !ok {
		t.Fatalf("expected allow after refill")
	}
}

func TestMemoryLimiter_Concurrent(t *testing.T) {
	lim := NewMemoryLimiter(1*time.Minute, 100)
	key := "ip-concurrent"

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = lim.Allow(context.Background(), key, 10, 5)
		}()
	}
	wg.Wait()
}

func TestMemoryLimiter_TTL(t *testing.T) {
	ttl := 100 * time.Millisecond
	lim := NewMemoryLimiter(ttl, 1)
	key := "ip-ttl"

	ok, _ := lim.Allow(context.Background(), key, 1, 1)
	if !ok {
		t.Fatalf("first allow failed")
	}

	time.Sleep(ttl + 50*time.Millisecond)

	ok, _ = lim.Allow(context.Background(), key, 1, 1)
	if !ok {
		t.Fatalf("expected allow after TTL cleanup")
	}
}
