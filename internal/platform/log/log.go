package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/LevantateLabs/basaltrouter/internal/platform/config"
)

type ctxKey struct{}

var defaultLogger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))

func Init(cfg config.LoggingConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	var w io.Writer
	switch strings.ToLower(cfg.Output) {
	case "stderr":
		w = os.Stderr
	default:
		w = os.Stdout
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{Level: level}
	switch strings.ToLower(cfg.Format) {
	case "text":
		handler = slog.NewTextHandler(w, opts)
	default:
		handler = slog.NewJSONHandler(w, opts)
	}

	logger := slog.New(handler)
	defaultLogger = logger
	slog.SetDefault(logger)
	return logger
}

func Default() *slog.Logger {
	return defaultLogger
}

func WithContext(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, logger)
}

func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok && logger != nil {
		return logger
	}
	return defaultLogger
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	logger := FromContext(ctx).With("request_id", requestID)
	return WithContext(ctx, logger)
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
