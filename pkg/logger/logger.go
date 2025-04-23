package logger

import (
	"io"
	"log/slog"
	"os"
)

func InitLogger(env EnvType, minLevel slog.Level) {
	var handler slog.Handler

	switch env {
	case EnvTypeTesting:
		handler = slog.NewJSONHandler(io.Discard, nil)
	case EnvTypeLocal:
		handler = newPrettyJSONHandler(os.Stdout, minLevel)
	case EnvTypeProd:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: minLevel,
		})
	default:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: minLevel,
		})
	}

	contextHandler := newContextHandler(handler)
	slog.SetDefault(slog.New(contextHandler))
}
