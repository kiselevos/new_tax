package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type authInfo struct {
	Type      string // public | private | internal | health
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

func maskKeyPrefix(key string) string {
	if len(key) < 8 {
		return "****"
	}

	return key[:8] + "..."
}

// Аутентификация для private рассчета
//
// Правила:
// 1. Public рассчет разрешен всегда.
// 2. Private метод (CalculatePrivate) разрешен в 2 случаях:
//   - Внутренний трафик (запрос с UI): x-internal: true
//   - У запроса имеется: x-api-key: <PRIVATE_API_KEY>
func Auth(privateAPIKey string) grpc.UnaryServerInterceptor {

	privateMethods := map[string]bool{
		"/tax.TaxService/CalculatePrivate": true,
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		method := info.FullMethod

		if !privateMethods[method] {
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

			if key == privateAPIKey {
				ctx = WithAuthInfo(ctx, authInfo{
					Type:      "private",
					KeyPrefix: maskKeyPrefix(key),
					KeyValid:  true,
				})
				return handler(ctx, req)
			}

			return nil, status.Error(codes.PermissionDenied, "invalid api key")
		}

		return nil, status.Error(codes.PermissionDenied, "missing api key")
	}
}
