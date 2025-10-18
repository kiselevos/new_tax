package server

import (
	"crypto/rand"
	"encoding/hex"
	"log/slog"
	"net/http"
	"new_tax/pkg/logx"
	"time"
)

// WithLogging — middleware для логирования всех HTTP/Connect запросов
func WithLogging(base *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rid := r.Header.Get("x-request-id")
		if rid == "" {
			b := make([]byte, 8)
			_, _ = rand.Read(b)
			rid = hex.EncodeToString(b)
		}

		logger := base.With(
			"rid", rid,
			"method", r.Method,
			"path", r.URL.Path,
			"remote", r.RemoteAddr,
		)

		ctx := logx.Into(r.Context(), logger)
		r = r.WithContext(ctx)

		rw := &responseWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(rw, r)

		dur := time.Since(start)
		if rw.status >= 400 {
			logger.Error("http",
				"status", rw.status,
				"duration_ms", dur.Milliseconds(),
			)
		} else {
			logger.Info("http",
				"status", rw.status,
				"duration_ms", dur.Milliseconds(),
			)
		}
	})
}

// responseWriter нужен, чтобы перехватывать код ответа
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}
