package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"time"

	"github.com/kiselevos/new_tax/pkg/logx"
)

type ctxKeyRID struct{}

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.URL.Path == "/metrics" {
			next.ServeHTTP(w, r)
			return
		}

		rid := getRequestID(r)

		logger := slog.Default().With(
			"rid", rid,
			"path", r.URL.Path,
			"method", r.Method,
		)

		ctx := logx.Into(r.Context(), logger)
		ctx = context.WithValue(ctx, ctxKeyRID{}, rid)

		sr := &statusRecorder{ResponseWriter: w, status: 200}

		start := time.Now()

		next.ServeHTTP(sr, r.WithContext(ctx))

		logger.Info("http_request_completed",
			"status", sr.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)

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

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
