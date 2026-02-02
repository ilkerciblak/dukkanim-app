package middleware

import (
	"dukkanim-api/internal/platform/observability/metrics"
	"fmt"
	"net/http"
	"time"
)

func HttpMetrics(httpMetrics metrics.HttpMetrics) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start_time := time.Now()
			defer httpMetrics.Gauge.WithLabelValues(
				map[string]string{
					"http_method": r.Method,
					"http_path":   r.URL.Path,
				},
			).Dec()

			httpMetrics.Gauge.WithLabelValues(
				map[string]string{
					"http_method": r.Method,
					"http_path":   r.URL.Path,
				},
			).Inc()

			wrapped := wrapResponseWriter(w)

			next.ServeHTTP(wrapped, r)

			duration_ms := time.Since(start_time).Milliseconds()

			httpMetrics.Histogram.WithLabelValues(
				map[string]string{
					"http.method": r.Method,
					"http.path":   r.URL.Path,
				},
			).Observe(int(duration_ms))

			httpMetrics.Counter.WithLabelValues(
				map[string]string{
					"http.method":               r.Method,
					"http.path":                 r.URL.Path,
					"http.response_status_code": fmt.Sprint(wrapped.Status()),
				},
			).Inc()

		})
	}
}
