package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authInfo struct {
	Type      string
	KeyPrefix string
	KeyValid  bool
}

type authKey struct{}

func WithAuthInfo(ctx context.Context, info authInfo) context.Context {
	return context.WithValue(ctx, authKey{}, info)
}

func GetAuthInfo(ctx context.Context) (authInfo, bool) {
	v := ctx.Value(authKey{})
	if v == nil {
		return authInfo{}, false
	}
	return v.(authInfo), true
}

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

		// Пропускаем Public метод без apiKey
		if strings.HasSuffix(info.FullMethod, "/CalculatePublic") {
			ctx = WithAuthInfo(ctx, authInfo{
				Type:     "public",
				KeyValid: true,
			})
			return handler(ctx, req)
		}

		md, _ := metadata.FromIncomingContext(ctx)

		if vals := md.Get("x-internal"); len(vals) > 0 && vals[0] == "true" {
			ctx = WithAuthInfo(ctx, authInfo{
				Type:     "internal",
				KeyValid: true,
			})
			return handler(ctx, req)
		}

		if vals := md.Get("x-api-key"); len(vals) > 0 {
			key := vals[0]
			prefix := maskKeyPrefix(key)

			if key == privateAPIKey {
				ctx = WithAuthInfo(ctx, authInfo{
					Type:      "api_key",
					KeyPrefix: prefix,
					KeyValid:  true,
				})
				return handler(ctx, req)
			}

			return nil, status.Error(codes.PermissionDenied, "invalid api key")
		}

		return nil, status.Error(codes.PermissionDenied, "missing api key")
	}
}

func maskKeyPrefix(key string) string {
	if len(key) < 8 {
		return "****"
	}

	return key[:8] + "..."
}
