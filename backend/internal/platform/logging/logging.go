package logging

import (
	"context"
	"time"
)

type Logger interface {
	DEBUG(ctx context.Context, msg string, fields ...any)
	INFO(ctx context.Context, msg string, fields ...any)
	WARN(ctx context.Context, msg string, fields ...any)
	ERROR(ctx context.Context, msg string, fields ...any)
	FATAL(ctx context.Context, msg string, fields ...any)
	With(ctx context.Context, fields ...any) Logger
}

type key int

var loggerKey key

func FromContext(ctx context.Context) Logger {
	logger, k := ctx.Value(loggerKey).(Logger)
	if !k {

		panic("Logger interface not implemented")
	}

	return logger
}

func InjectLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

type Log struct {
	Level         string            `json:"level"`
	Service       string            `json:"service"`
	Timestamp     time.Time         `json:"time_stamp"`
	Message       string            `json:"message"`
	TransactionID string            `json:"transaction_id"`
	Duration      float32           `json:"duration"`
	Request       RequestLog        `json:"request"`
	Response      ResponseLog       `json:"response"`
	Fields        map[string]string `json:"fields"`
}

type RequestLog struct {
	Method    string            `json:"method"`
	Host      string            `json:"host"`
	Path      string            `json:"path"`
	Headers   map[string]string `json:"headers"`
	Query     string            `json:"query,omitempty"`
	Fragment  string            `json:"fragment,omitempty"`
	Body      map[string]string `json:"body,omitempty"`
	UserAgent string            `json:"user-agent"`
}

type ResponseLog struct {
	StatusCode int               `json:"status_code"`
	Body       map[string]string `json:"body"`
}
