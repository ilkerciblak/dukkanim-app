package middleware

import (
	response "dukkanim-api/internal/platform/http_response"
	ratelimiting "dukkanim-api/internal/platform/rate_limiting"
	"fmt"
	"net"
	"net/http"
)

func RateLimiting(algorithm ratelimiting.RateLimitingStrategy) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_, port, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				response.RespondWithProblemDetails(
					w,
					r.Context(),
					http.StatusInternalServerError,
					err.Error(),
					"SERVER_ERROR",
					nil,
				)
			}
			is_allowed, err := algorithm.Allow(r.Context(), port)
			if err != nil {
				response.RespondWithProblemDetails(
					w,
					r.Context(),
					http.StatusInternalServerError,
					err.Error(),
					"SERVER_ERROR",
					nil,
				)
			}
			info, err := algorithm.GetLimit(r.Context(), port)
			if err != nil {
				response.RespondWithProblemDetails(
					w,
					r.Context(),
					http.StatusInternalServerError,
					err.Error(),
					"SERVER_ERROR",
					nil,
				)
			}

			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d requests/per minute", info.Remaining))
			w.Header().Set("X-RateLimit-Remaining", fmt.Sprint(info.Remaining))
			w.Header().Set("X-RateLimit-Reset", fmt.Sprint(info.ResetTime.Local().String()))

			if is_allowed {
				next.ServeHTTP(w, r)
			} else {
				response.RespondWithProblemDetails(
					w,
					r.Context(),
					http.StatusTooManyRequests,
					"Request is rate limited",
					"RATE_LIMITED",
					nil,
				)
			}

		})

	}
}
