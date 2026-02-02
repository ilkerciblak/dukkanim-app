package logging

import (
	"context"
	"time"
)

type Logger interface {
	// Methods to Log Messages
	Debug(ctx context.Context, msg string, attrs ...any)
	Info(ctx context.Context, msg string, attrs ...any)
	Warning(ctx context.Context, msg string, attrs ...any)
	Error(ctx context.Context, msg string, attrs ...any)
	Fatal(ctx context.Context, msg string, attrs ...any)
	// Method to manipulate attributes with single instance
	With(ctx context.Context, withAttr ...WithEntryField)
	// Method to populate new logger instance with new request context
	Using(ctx context.Context) Logger
}

type key int

var loggerKey key

func FromContext(ctx context.Context) Logger {
	return ctx.Value(loggerKey).(Logger)
}

func InjectContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

type Entry struct {
	Timestamp     time.Time
	Level         string
	Service       string
	Message       string
	Context       ContextEntry
	TransactionId string
	DurationMS    int
	Request       RequestEntry
	Response      ResponseEntry
	Fields        map[string]any
}

type ContextEntry struct {
}

type RequestEntry struct {
	Method    string
	Path      string
	Query     string
	Fragment  string
	UserAgent string
	Headers   map[string][]string
}

// Method:    r.Method,
// 	Path:      r.URL.Path,
// 	Query:     r.URL.RawQuery,
// 	Fragment:  r.URL.Fragment,
// 	UserAgent: r.UserAgent(),
// 	Headers:   r.Header,

type ResponseEntry struct {
	StatusCode int
}

type WithEntryField func(Entry) Entry

func SetTimeStamp(timestamp time.Time) WithEntryField {
	return func(e Entry) Entry {
		e.Timestamp = timestamp
		return e
	}
}

func SetLevel(val string) WithEntryField {
	return func(e Entry) Entry {
		e.Level = val
		return e
	}
}
func SetService(val string) WithEntryField {
	return func(e Entry) Entry {
		e.Service = val
		return e
	}
}
func SetMessage(val string) WithEntryField {
	return func(e Entry) Entry {
		e.Message = val
		return e
	}
}
func SetContext(val ContextEntry) WithEntryField {
	return func(e Entry) Entry {
		e.Context = val
		return e
	}
}
func SetTransactionId(val string) WithEntryField {
	return func(e Entry) Entry {
		e.TransactionId = val
		return e
	}
}
func SetDurationMS(val int) WithEntryField {
	return func(e Entry) Entry {
		e.DurationMS = val
		return e
	}
}
func SetRequest(val RequestEntry) WithEntryField {
	return func(e Entry) Entry {
		e.Request = val
		return e
	}
}
func SetResponse(val ResponseEntry) WithEntryField {
	return func(e Entry) Entry {
		e.Response = val
		return e
	}
}
func AppendField(key string, val any) WithEntryField {
	return func(e Entry) Entry {
		if existingVal, exists := e.Fields[key]; exists {
			switch value := existingVal.(type) {
			case []any:
				value = append(value, val)
				e.Fields[key] = value
			case any:
				e.Fields[key] = []any{value, val}
			}
		}

		e.Fields[key] = val

		return e
	}
}

func (e Entry) ToFieldSlice() []any {
	return []any{
		"timestamp", e.Timestamp,
		"level", e.Level,
		"service", e.Service,
		"message", e.Message,
		"context", e.Context,
		"transaction_id", e.TransactionId,
		"durationMS", e.DurationMS,
		"request", e.Request,
		"response", e.Response,
		"fields", e.Fields,
	}
}
