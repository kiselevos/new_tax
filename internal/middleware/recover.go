package middleware

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/kiselevos/new_tax/pkg/logx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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
