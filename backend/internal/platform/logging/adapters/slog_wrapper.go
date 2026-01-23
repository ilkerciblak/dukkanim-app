package logging

import (
	"context"
	"dukkanim-api/internal/platform/logging"
	"log/slog"
	"os"
)

type slogAdapter struct {
	handler *slog.Logger
}

func NewSlogLogger() logging.Logger {
	handler := slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{

				Level: slog.LevelDebug,
			},
		),
	)

	return &slogAdapter{
		handler: handler,
	}
}

func (l slogAdapter) DEBUG(ctx context.Context, msg string, fields ...any) {
	l.handler.LogAttrs(ctx, slog.LevelDebug, msg, fieldToAnyAttr(fields)...)
}
func (l slogAdapter) INFO(ctx context.Context, msg string, fields ...any) {
	l.handler.LogAttrs(ctx, slog.LevelInfo, msg, fieldToAnyAttr(fields)...)
}
func (l slogAdapter) WARN(ctx context.Context, msg string, fields ...any) {
	l.handler.LogAttrs(ctx, slog.LevelWarn, msg, fieldToAnyAttr(fields)...)
}
func (l slogAdapter) ERROR(ctx context.Context, msg string, fields ...any) {
	l.handler.LogAttrs(ctx, slog.LevelError, msg, fieldToAnyAttr(fields)...)
}
func (l slogAdapter) FATAL(ctx context.Context, msg string, fields ...any) {
	l.handler.LogAttrs(ctx, slog.LevelError, msg, fieldToAnyAttr(fields)...)
}

func (l *slogAdapter) With(ctx context.Context, fields ...any) logging.Logger {

	return &slogAdapter{
		handler: l.handler.With(fields...),
	}

}

func fieldToAnyAttr(fields []any) []slog.Attr {
	var attrs []slog.Attr
	for i := 0; i <= len(fields)-2; i += 2 {
		if abs, k := fields[i].(string); k {
			attrs = append(attrs, slog.Attr{Key: abs, Value: slog.AnyValue(fields[i+1])})
		} else {
			attrs = append(attrs, slog.Attr{Key: "WierdKey", Value: slog.AnyValue(fields[i+1])})
		}

	}

	return attrs
}
