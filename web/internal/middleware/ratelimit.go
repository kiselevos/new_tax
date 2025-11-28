package middleware

import (
	"net/http"

	"golang.org/x/time/rate"
)

// RateLimiterMiddleware ограничивает частоту запросов.
func RateLimiterMiddleware(rps int, burst int) func(http.Handler) http.Handler {

	limiter := rate.NewLimiter(rate.Limit(rps), burst)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if !limiter.Allow() {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				if _, err := w.Write([]byte(`{"error":"too many requests"}`)); err != nil {
					return
				}
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
