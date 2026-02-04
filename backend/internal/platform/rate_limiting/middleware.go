package ratelimiting

import (
	"fmt"
	"net"
	"net/http"
)

func Middleware(strategy RateLimitingStrategy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, port, _ := net.SplitHostPort(r.RemoteAddr)
			is_allowed, err := strategy.Allow(r.Context(), port)
			info, err := strategy.GetLimit(r.Context(), port)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d requests/per minute", info.Remaining))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprint(info.Remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprint(info.ResetTime.Local().String()))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			if is_allowed {
				next.ServeHTTP(w, r)
			} else {
				http.Error(w, "REQUEST_LIMITED", http.StatusTooManyRequests)
			}

		})
	}
}
