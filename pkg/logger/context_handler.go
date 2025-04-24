package logger

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
)

type contextHandler struct {
	handler slog.Handler
}

func newContextHandler(handler slog.Handler) *contextHandler {
	return &contextHandler{
		handler: handler,
	}
}

func (h *contextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if reqID, ok := ctx.Value(contextKeyRequestID).(uuid.UUID); ok {
		r.AddAttrs(slog.String(string(contextKeyRequestID), reqID.String()))
	}

	if envType, ok := ctx.Value(contextKeyEnvType).(EnvType); ok {
		r.AddAttrs(slog.String(string(contextKeyEnvType), string(envType)))
	}

	return h.handler.Handle(ctx, r)
}

func (h *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return newContextHandler(h.handler.WithAttrs(attrs))
}

func (h *contextHandler) WithGroup(name string) slog.Handler {
	return newContextHandler(h.handler.WithGroup(name))
}
