package logger

import (
	"context"
	"log/slog"
	"os"
)

// Available logging levels
const (
	levelInfo = "info"
	levelWarn = "warning"
	levelErr  = "error"
)

// Create new text logger with level. Debug by default
func NewTextLogger(level string) *slog.Logger {
	switch level {
	case levelInfo:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case levelWarn:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	case levelErr:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	default:
		return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	}
}

// Set logger into context
func NewCtxWithLog(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, slog.Logger{}, log)
}

// Use context logger value to log data
// Returns logger with debug level by default
func LogUse(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(slog.Logger{}).(*slog.Logger); ok {
		return log
	}

	return NewTextLogger("")
}
