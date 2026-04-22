package ratelimit

import (
	"context"
	"net"
	"strings"

	"github.com/kiselevos/new_tax/internal/config"
	"github.com/kiselevos/new_tax/internal/middleware"
	"github.com/kiselevos/new_tax/pkg/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

type Limiter interface {
	Allow(ctx context.Context, key string, rps float64, burst int) (bool, error)
}

func RateLimitInterceptor(limiter Limiter, cfg *config.RateLimitConfig) grpc.UnaryServerInterceptor {
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
		if ai, ok := middleware.GetAuthInfo(ctx); ok && ai.Type == "internal" {
			return handler(ctx, req)
		}

		ip := getClientIP(ctx)

		var (
			rps   float64
			burst int
			scope string
		)

		if isPrivate(method) {
			rps = cfg.PrivateRPS
			burst = cfg.PrivateBurst
			scope = "private"
		} else {
			rps = cfg.PublicRPS
			burst = cfg.PublicBurst
			scope = "public"
		}

		key := "ip_" + scope + "_" + ip

		allowed, err := limiter.Allow(ctx, key, rps, burst)
		if err != nil {
			logx.From(ctx).Warn("rate_limiter_failed",
				"err", err,
				"ip", ip,
				"method", method,
				"scope", scope,
			)
			return handler(ctx, req)
		}

		if !allowed {
			logx.From(ctx).Warn("rate_limit_blocked",
				"ip", ip,
				"method", method,
				"scope", scope,
			)
			return nil, status.Error(codes.ResourceExhausted, "too many requests")
		}

		return handler(ctx, req)
	}
}

func getClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok && p.Addr != nil {
		if ta, ok := p.Addr.(*net.TCPAddr); ok {
			return ta.IP.String()
		}
		return p.Addr.String()
	}
	return "unknown"
}

func isPrivate(method string) bool {
	return method == "/tax.TaxService/CalculatePrivate"
}
