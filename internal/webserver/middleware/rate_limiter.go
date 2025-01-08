package middleware

import (
	"context"
	"net"
	"net/http"

	"github.com/vs0uz4/rate-limit/internal/contract"
)

func extractIP(remoteAddr string) string {
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return remoteAddr
	}

	if host == "::1" {
		return "127.0.0.1"
	}

	return host
}

func RateLimiterMiddleware(limiter contract.RateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := extractIP(r.RemoteAddr)
			token := r.Header.Get("API_KEY")

			allowed, err := limiter.Allow(context.Background(), ip, token)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "You have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
