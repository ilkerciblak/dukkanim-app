package tracing

import (
	"context"
	"net/http"
)

type Tracer interface {
	Start(ctx context.Context, name string) (context.Context, Span)
	Extract(ctx context.Context, header http.Header) context.Context
	Inject(ctx context.Context, header http.Header)
}

type Span interface {
	End()
	SetAttributes(attributes ...SpanAttributePair)
	SetStatus(OperationStatus, string)
	RecordError(err error)
	SpanContext() SpanContext
}

type OperationStatus int

const (
	Success OperationStatus = iota
	Error
	Unset
)

type SpanContext struct {
	SpanID  string `json:"span_id"`
	TraceID string `json:"trace_id"`
}

type SpanAttributePair struct {
	Key   string
	Value any
}
