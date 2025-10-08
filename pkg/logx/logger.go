package logx

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey struct{}

func Into(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, l)
}

// Достаем логгер из контекста или дефолтный
func From(ctx context.Context) *slog.Logger {
	if v := ctx.Value(ctxKey{}); v != nil {
		if l, ok := v.(*slog.Logger); ok && l != nil {
			return l
		}
	}
	return slog.Default()
}

func New() *slog.Logger {

	mode := os.Getenv("LOG_MODE")
	level := os.Getenv("LOG_LEVEL")

	var h slog.Handler
	if mode == "text" {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLevel(level),
		})
	} else {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: parseLevel(level),
		})
	}

	l := slog.New(h)
	slog.SetDefault(l)
	return l
}

func parseLevel(lvl string) slog.Level {
	switch lvl {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
