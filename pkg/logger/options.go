package logger

import (
	"context"

	"github.com/google/uuid"
)

func WithRequestID(ctx context.Context, requestID uuid.UUID) context.Context {
	return context.WithValue(ctx, contextKeyRequestID, requestID)
}
