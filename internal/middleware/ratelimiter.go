package middleware

import (
	"context"
	"strings"
	"sync"

	"github.com/kiselevos/new_tax/internal/config"
	"github.com/kiselevos/new_tax/pkg/logx"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var mu sync.Mutex
var limiters = make(map[string]*rate.Limiter)

func getLimiter(key string, r rate.Limit, burst int) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	if l, ok := limiters[key]; ok {
		return l
	}

	l := rate.NewLimiter(r, burst)
	limiters[key] = l
	return l
}

func RateLimitInterceptor(cfg *config.RateLimitConfig) grpc.UnaryServerInterceptor {
	return func(ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		method := info.FullMethod

		if strings.HasSuffix(info.FullMethod, "/Healthz") {
			return handler(ctx, req)
		}

		// Не тротлим UI
		if ai, ok := GetAuthInfo(ctx); ok && ai.Type == "internal" {
			return handler(ctx, req)
		}

		ip := getClientIP(ctx)
		var limiter *rate.Limiter

		if isPrivate(method) {
			limiter = getLimiter("ip_"+ip, rate.Limit(cfg.PrivateRPS), cfg.PrivateBurst)
		} else {
			limiter = getLimiter("ip_"+ip, rate.Limit(cfg.PublicRPS), cfg.PublicBurst)
		}

		if !limiter.Allow() {
			logx.From(ctx).Warn("rate_limit_blocked",
				"ip", ip,
				"method", method,
			)
			return nil, status.Error(codes.ResourceExhausted, "too many requests")
		}

		return handler(ctx, req)
	}
}

func getClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		return p.Addr.String()
	}
	return "unknown"
}

func isPrivate(method string) bool {
	return method == "/tax.TaxService/CalculatePrivate"
}
