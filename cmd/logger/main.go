package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/meteogo/logger/pkg/logger"
)

func main() {
	logger.InitLogger(logger.EnvTypeLocal, slog.LevelDebug)

	ctx := logger.WithRequestID(context.Background(), uuid.New())
	logger.Info(ctx, "some test message", slog.String("key", "value"))
}
