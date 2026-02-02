package middleware

import (
	"dukkanim-api/internal/platform/observability/tracing"
	"fmt"
	"net/http"
)

func DistributedTracing(tracer tracing.Tracer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestCtx := tracer.Extract(r.Context(), r.Header)

			ctx, span := tracer.Start(requestCtx, fmt.Sprintf("%s %s.middleware", r.Method, r.URL.Path))

			span.SetAttributes(
				tracing.SpanAttributePair{
					Key:   "http.method",
					Value: r.Method,
				},
				tracing.SpanAttributePair{
					Key:   "Route",
					Value: r.URL.Path,
				},

				tracing.SpanAttributePair{
					Key:   "Client",
					Value: r.RemoteAddr,
				},
			)

			defer span.End()

			r = r.WithContext(ctx)

			tracer.Inject(ctx, r.Header)

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			span.SetAttributes(

				tracing.SpanAttributePair{
					Key: "http.status-code", Value: wrapped.Status(),
				})

		})
	}
}

// import (
// 	"net/http"

// 	"dukkanim-api/internal/platform/observability/tracing"
// )

// func DistributedTracing(tracer tracing.Tracer) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

// 			// Extract or generate trrace ID from request

// 			// Start a new span for the HTTP request
// 			ctx, span := tracer.StartSpan(r.Context(), r.Method+" "+r.URL.Path)
// 			span.SetAttribute("http.method", r.Method)
// 			span.SetAttribute("http.route", r.URL.Path)
// 			span.SetAttribute("http.client-ip", r.RemoteAddr)

// 			defer span.End()
// 			trace_id := span.TraceID()

// 			// Inject span context into request ctx
// 			r = r.WithContext(ctx)
// 			w.Header().Set("X-Trace-ID", trace_id)

// 			// Propagate trace context in response header
// 			wrapped := wrapResponseWriter(w)
// 			next.ServeHTTP(wrapped, r)
// 			span.SetAttribute("http.status-code", wrapped.Status())

// 			// End span after rerquesst compeleted with status code

// 		})
// 	}
// }
