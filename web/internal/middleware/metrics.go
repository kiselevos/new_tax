package middleware

import (
	"net/http"

	"github.com/kiselevos/new_tax/web/internal/metrics"
)

func MetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// исключаем Prometheus
		if r.URL.Path != "/metrics" {
			metrics.M.System.HTTPRequests.Inc()
		}

		next.ServeHTTP(w, r)
	})
}
