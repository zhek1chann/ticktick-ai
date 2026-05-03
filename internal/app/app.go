package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"ticktick-ai/internal/modules/tg-bot/handler"
	"ticktick-ai/pkg/closer"
	"ticktick-ai/pkg/logger"
)

type App struct {
	serviceProvider *serviceProvider
}

func New(ctx context.Context) (*App, error) {
	app := &App{}

	if err := app.initDeps(ctx); err != nil {
		return nil, err
	}

	slog.Info("dependencies initialized")
	return app, nil
}

func (a *App) Run(ctx context.Context) error {
	runCtx, stop := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer stop()

	bot := a.serviceProvider.TgBot(runCtx)
	tgHandler := handler.New(a.serviceProvider.TgBotService(runCtx))
	tgHandler.RegisterRoutes(bot)

	go func() {
		<-runCtx.Done()
		bot.Stop()
		closer.CloseAll()
	}()

	slog.InfoContext(runCtx, "telegram bot started")
	go bot.Start()

	closer.Wait()
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	for _, f := range []func(context.Context) error{
		a.initServiceProvider,
		a.initLogger,
	} {
		if err := f(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initServiceProvider(ctx context.Context) error {
	a.serviceProvider = newServiceProvider(ctx)
	return nil
}

func (a *App) initLogger(ctx context.Context) error {
	logger.InitLogging(a.serviceProvider.Config().App().LogLevel())
	return nil
}
