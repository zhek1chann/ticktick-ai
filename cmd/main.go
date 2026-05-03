package main

import (
	"context"
	"log/slog"
	"os"
	"ticktick-ai/internal/app"
	"ticktick-ai/pkg/logger"
)

func main() {
	ctx := context.Background()

	app, err := app.New(ctx)
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
		os.Exit(1)
	}

	if err := app.Run(ctx); err != nil {
		slog.ErrorContext(logger.ErrorCtx(ctx, err), "Error: "+err.Error())
		os.Exit(1)
	}
}
