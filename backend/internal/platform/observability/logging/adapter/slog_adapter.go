package adapter

import (
	"context"
	logging "dukkanim-api/internal/platform/observability/logging"
	"log/slog"
	"os"
)

type slogLogger struct {
	logger *slog.Logger
	entry  logging.Entry
}

func NewslogLogger(logLevel string) logging.Logger {
	var logLevelInt slog.Level
	switch logLevel {
	case "DEBUG":
		logLevelInt = slog.LevelDebug
	case "INFO":
		logLevelInt = slog.LevelInfo
	case "WARN":
		logLevelInt = slog.LevelWarn
	case "ERROR":
		logLevelInt = slog.LevelError
	}

	handler := slog.NewJSONHandler(
		os.Stdout,
		&slog.HandlerOptions{
			Level: logLevelInt,
		},
	)

	logger := slog.New(handler)

	return &slogLogger{
		logger: logger,
	}

}

func (s *slogLogger) Debug(ctx context.Context, msg string, attrs ...any) {

	s.logger.DebugContext(ctx, msg, s.appendFields(attrs...)...)

}
func (s *slogLogger) Info(ctx context.Context, msg string, attrs ...any) {
	s.logger.InfoContext(ctx, msg, s.appendFields(attrs...)...)

}

func (s *slogLogger) Warning(ctx context.Context, msg string, attrs ...any) {
	s.logger.WarnContext(ctx, msg, s.appendFields(attrs...)...)

}

func (s *slogLogger) Error(ctx context.Context, msg string, attrs ...any) {
	s.logger.ErrorContext(
		ctx,
		msg,
		s.appendFields(attrs...)...,
	)
}

func (s *slogLogger) Fatal(ctx context.Context, msg string, attrs ...any) {
	s.logger.ErrorContext(
		ctx,
		msg,
		s.appendFields(attrs...)...,
	)
}

func (s *slogLogger) With(ctx context.Context, withAttr ...logging.WithEntryField) {
	for _, f := range withAttr {
		s.entry = f(s.entry)
	}
}

func (s *slogLogger) Using(ctx context.Context) logging.Logger {
	transactionId := ctx.Value("request-id").(string)
	// requestId := ctx.Value("request-id").(string)
	logger := slog.New(s.logger.Handler())

	return &slogLogger{
		logger: logger,
		entry: logging.Entry{
			TransactionId: transactionId,
		},
	}
}

func (s slogLogger) appendFields(attrs ...any) []any {
	attrs = append(attrs, s.entry.ToFieldSlice()...)
	return attrs
}
