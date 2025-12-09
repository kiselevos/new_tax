package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// Аутентификация для private рассчета
//
// Правила:
// 1. Public рассчет разрешен всегда.
// 2. Private метод (CalculatePrivate) разрешен в 2 случаях:
//   - Внутренний трафик (запрос с UI): x-internal: true
//   - У запроса имеется: x-api-key: <PRIVATE_API_KEY>
func Auth(privateAPIKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		// Ограничиваем только приватный метод
		if strings.HasSuffix(info.FullMethod, "/CalculatePrivate") {

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "metadata missing")
			}

			// 1) Internal UI call
			if vals := md.Get("x-internal"); len(vals) > 0 && vals[0] == "true" {
				return handler(ctx, req)
			}

			// 2) External API call — requires API-key
			if vals := md.Get("x-api-key"); len(vals) > 0 {
				apiKey := vals[0]
				if apiKey == privateAPIKey {
					return handler(ctx, req)
				}
				return nil, status.Error(codes.PermissionDenied, "invalid api key")
			}

			return nil, status.Error(codes.PermissionDenied, "private api requires a valid api key")
		}

		return handler(ctx, req)
	}
}
