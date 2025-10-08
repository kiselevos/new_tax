package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"new_tax/pkg/logx"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const requestIDKey = "x-request-id"

func UnaryLogger(base *slog.Logger) grpc.UnaryServerInterceptor {

	return func(
		ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
	) (interface{}, error) {
		var rid string
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if vals := md.Get(requestIDKey); len(vals) > 0 {
				rid = vals[0]
			}
		}
		if rid == "" {
			rid = newRID()
		}

		logger := base.With("method", info.FullMethod, "rid", rid)
		if p, ok := peer.FromContext(ctx); ok && p != nil && p.Addr != nil {
			logger = logger.With("peer", p.Addr.String())
		}
		ctx = logx.Into(ctx, logger)

		start := time.Now()
		resp, err := handler(ctx, req)
		dur := time.Since(start)

		st := status.Convert(err)
		logger.Info("grpc",
			"code", st.Code().String(),
			"duration_ms", dur.Milliseconds(),
		)
		return resp, err
	}
}

func newRID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
