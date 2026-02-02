package adapter

import (
	"context"
	"dukkanim-api/internal/platform/observability/tracing"
	"net/http"
	"time"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	"go.opentelemetry.io/otel/trace"
)

type otelTracer struct {
	tracer trace.Tracer
}

type otelSpan struct {
	span trace.Span
}

func (s otelSpan) End() {
	s.span.End(trace.WithTimestamp(time.Now()))
}
func (s otelSpan) SetAttributes(attributes ...tracing.SpanAttributePair) {
	for _, attr := range attributes {
		s.span.SetAttributes(attributeToOtelAttribute(attr))
	}
}
func (s otelSpan) SetStatus(status tracing.OperationStatus, description string) {
	switch status {
	case tracing.Error:
		s.span.SetStatus(codes.Error, description)
	case tracing.Success:
		s.span.SetStatus(codes.Ok, description)
	default:
		s.span.SetStatus(codes.Unset, description)
	}
}
func (s otelSpan) RecordError(err error) {
	s.span.RecordError(err)
}
func (s otelSpan) SpanContext() tracing.SpanContext {
	span_context := s.span.SpanContext()
	return tracing.SpanContext{
		SpanID:  span_context.SpanID().String(),
		TraceID: span_context.TraceID().String(),
	}
}

func NewOtelTracer() tracing.Tracer {

	exp, _ := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)

	r, _ := resource.Merge(resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("dukkanim-api"),
			semconv.ServiceVersion("1.0.0"),
		),
	)

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(r),
	)

	return &otelTracer{
		tracer: provider.Tracer("dukkanim-api"),
	}

}

func (o otelTracer) Start(ctx context.Context, name string) (context.Context, tracing.Span) {
	ctx, span := o.tracer.Start(ctx, name)
	return ctx, &otelSpan{
		span: span,
	}
}
func (o otelTracer) Extract(ctx context.Context, header http.Header) context.Context {
	requestCtx := otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(header))

	if _, exists := requestCtx.Value("trace-id").(string); exists {
		return requestCtx
	}

	trace_id := uuid.New().String()
	return context.WithValue(ctx, "trace-id", trace_id)

}

func (o otelTracer) Inject(ctx context.Context, header http.Header) {
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(header))
}

func attributeToOtelAttribute(t tracing.SpanAttributePair) attribute.KeyValue {
	switch v := t.Value.(type) {
	case int:
		return attribute.Int(t.Key, v)
	case int64:
		return attribute.Int64(t.Key, v)
	case string:
		return attribute.String(t.Key, v)
	case bool:
		return attribute.Bool(t.Key, v)
	case float64:
		return attribute.Float64(t.Key, v)
	default:
		return attribute.String(t.Key, "un-processable")
	}
}
