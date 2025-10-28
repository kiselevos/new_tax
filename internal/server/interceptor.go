package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net"
	"runtime/debug"
	"time"

	"github.com/kiselevos/new_tax/pkg/logx"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const requestIDKey = "x-request-id"

func UnaryLogger(base *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// request id
		rid := newRID(ctx)

		// базовый логгер
		logger := base.With("method", info.FullMethod, "rid", rid)

		// peer ip (без порта)
		if p, ok := peer.FromContext(ctx); ok && p != nil && p.Addr != nil {
			host, _, _ := net.SplitHostPort(p.Addr.String())
			if host == "" {
				host = p.Addr.String()
			}
			logger = logger.With("peer_ip", host)
		}

		// прокидываем rid обратно клиенту
		_ = grpc.SetHeader(ctx, metadata.Pairs(requestIDKey, rid))

		// кладём request-scoped логгер в контекст
		ctx = logx.Into(ctx, logger)

		start := time.Now()
		resp, err := handler(ctx, req)

		code := status.Code(err)
		level := slog.LevelInfo
		switch code {
		case codes.OK, codes.Canceled, codes.InvalidArgument, codes.NotFound,
			codes.AlreadyExists, codes.PermissionDenied, codes.Unauthenticated,
			codes.FailedPrecondition, codes.OutOfRange:
			level = slog.LevelInfo
		case codes.ResourceExhausted, codes.Aborted, codes.Unavailable, codes.DeadlineExceeded:
			level = slog.LevelWarn
		default: // Unknown, Internal, DataLoss, Unimplemented
			level = slog.LevelError
		}

		attrs := []any{
			"code", code.String(),
			"duration", time.Since(start).String(),
		}
		if err != nil {
			attrs = append(attrs, "err", err)
		}

		logger.Log(ctx, level, "grpc_request_done", attrs...)
		return resp, err
	}
}

func UnaryRecovery(base *slog.Logger) grpc.UnaryServerInterceptor {

	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				stack := string(debug.Stack())
				logx.From(ctx).Error("panic_recovered",
					"method", info.FullMethod,
					"recover", r,
					"stack", stack,
				)
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(ctx, req)
	}
}

func newRID(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if vals := md.Get(requestIDKey); len(vals) > 0 && vals[0] != "" {
			return vals[0]
		}
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
