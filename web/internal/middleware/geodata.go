package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/kiselevos/new_tax/web/data"
	"github.com/kiselevos/new_tax/web/internal/geoip"
)

type ctxKeyRegion struct{}

func getClientIP(r *http.Request) string {
	// 1. X-Forwarded-For: client, proxy1, proxy2
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		ip := strings.TrimSpace(parts[0])
		if ip != "" {
			return ip
		}
	}

	// 2. X-Real-IP
	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		return strings.TrimSpace(xrip)
	}

	// 3. RemoteAddr
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}

	return r.RemoteAddr
}

func isIPv4(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	return ip != nil && ip.To4() != nil
}

func RegionMiddleware(db *geoip.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ip := getClientIP(r)

			region := data.Region{
				Name:  "Unknown",
				Label: "other",
			}
			if isIPv4(ip) {
				regionName := db.LookupRegion(ip)
				region = data.NormalizeRegion(regionName)
			}

			ctx := context.WithValue(r.Context(), ctxKeyRegion{}, region)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetRegion(ctx context.Context) data.Region {
	if v := ctx.Value(ctxKeyRegion{}); v != nil {
		if region, ok := v.(data.Region); ok {
			return region
		}
	}
	return data.Region{
		Name:  "Unknown",
		Label: "other",
	}
}
