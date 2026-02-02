package middleware

import (
	"context"
	"dukkanim-api/internal/platform/observability/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RequestLogging(baseLogger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)

			request_id := uuid.New().String()

			ctx := context.WithValue(r.Context(), "request-id", request_id)

			requestLogger := baseLogger.Using(ctx)
			requestLogger.With(r.Context(), logging.SetRequest(logging.RequestEntry{
				Method:    r.Method,
				Path:      r.URL.Path,
				Query:     r.URL.RawQuery,
				Fragment:  r.URL.Fragment,
				UserAgent: r.UserAgent(),
				Headers:   r.Header,
			}))

			// Request Initiating Logging
			requestLogger.Debug(ctx, "Request Started")
			// Injecting Logger to Request Context
			ctx = logging.InjectContext(ctx, requestLogger)
			// Serving the Next Handler with new Request Context
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Response Logging

			duration_ms := time.Since(start).Milliseconds()

			if duration_ms > 5000 {
				requestLogger.Warning(ctx, "Duration Higher Than Treshold")
			}

			requestLogger.With(
				ctx,
				logging.SetDurationMS(int(duration_ms)),
				logging.SetTimeStamp(time.Now()),
				logging.SetResponse(logging.ResponseEntry{
					StatusCode: wrapped.Status(),
				}),
			)

			requestLogger.Info(
				r.Context(),
				"Request Completed",
			)

		})
	}
}

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func wrapResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{ResponseWriter: w}
}

func (rw *responseWriter) Status() int {
	return rw.status
}

func (rw *responseWriter) WriteHeader(code int) {
	if rw.wroteHeader {
		return
	}

	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
	rw.wroteHeader = true

}
