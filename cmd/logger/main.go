package main

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/meteogo/logger/pkg/logger"
)

func main() {
	logger.InitLogger(logger.EnvTypeLocal, slog.LevelDebug)

	ctx := context.Background()
	ctx = logger.WithRequestID(ctx, uuid.New())
	ctx = logger.WithEnvType(ctx, logger.EnvTypeLocal)

	logger.Error(ctx, "user is authenticated", slog.String("userID", uuid.NewString()))
}
