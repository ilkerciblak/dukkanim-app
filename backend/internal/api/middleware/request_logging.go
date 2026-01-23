package middleware

import (
	"context"
	response "dukkanim-api/internal/platform/http_response"
	"dukkanim-api/internal/platform/logging"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func RequestLogging(baseLogger logging.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			wrapped := wrapResponseWriter(w)

			request_id := uuid.New()

			ctx := context.WithValue(r.Context(), response.RequestIdKey, request_id)

			request := map[string]any{
				"host":       r.Host,
				"path":       r.URL.Path,
				"query":      r.URL.RawQuery,
				"fragment":   r.URL.Fragment,
				"method":     r.Method,
				"user-agent": r.UserAgent(),
				"target":     r.RequestURI,
				"headers":    r.Header,
			}

			requestLogger := baseLogger.With(
				ctx,
				"request_id", request_id,
				"request", request,
			)
			// Request Initiating Logging
			requestLogger.DEBUG(ctx, "Request Started")
			// Injecting Logger to Request Context
			ctx = logging.InjectLogger(ctx, requestLogger)
			// Serving the Next Handler with new Request Context
			next.ServeHTTP(wrapped, r.WithContext(ctx))

			// Response Logging
			responseStatus := wrapped.status
			duration_ms := time.Since(start).Milliseconds()

			requestLogger.DEBUG(
				r.Context(),
				"Request Completed",
				"status_code", responseStatus,
				"duration_ms", duration_ms,
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
