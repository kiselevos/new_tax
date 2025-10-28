package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"github.com/kiselevos/new_tax/pkg/logx"
)

type ctxKeyRID struct{}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rid := getRequestID(r)
		logger := logx.New().With("rid", rid, "path", r.URL.Path, "method", r.Method)
		ctx := logx.Into(r.Context(), logger)
		ctx = context.WithValue(ctx, ctxKeyRID{}, rid)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getRequestID(r *http.Request) string {
	if v := r.Header.Get("X-Request-ID"); v != "" {
		return v
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func GetRID(ctx context.Context) string {
	if v := ctx.Value(ctxKeyRID{}); v != nil {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}
