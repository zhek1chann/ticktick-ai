package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"ticktick-ai/internal/config"
	ai "ticktick-ai/internal/modules/ai"
	"ticktick-ai/internal/modules/ai/gemini"
	tgbot "ticktick-ai/internal/modules/tg-bot/service"
	"ticktick-ai/internal/modules/ticktick"
	ticktickclient "ticktick-ai/internal/modules/ticktick/clients/ticktick"
	"ticktick-ai/pkg/logger"
	"time"

	tele "gopkg.in/telebot.v3"
)

const configEnvPath = ".env"

type serviceProvider struct {
	ctx context.Context
	cfg *config.Config

	httpClient *http.Client
	tgBot      *tele.Bot

	geminiClient   *gemini.Client
	aiService      *ai.Service
	ticktickClient *ticktickclient.Client
	ticktickSvc    *ticktick.Service
	tgBotService   *tgbot.BotService
}

func newServiceProvider(ctx context.Context) *serviceProvider {
	return &serviceProvider{ctx: ctx}
}

func (s *serviceProvider) Config() *config.Config {
	if s.cfg == nil {
		var err error
		s.cfg, err = config.LoadConfig(configEnvPath)
		if err != nil {
			slog.ErrorContext(logger.ErrorCtx(s.ctx, err), "config error: "+err.Error())
			os.Exit(1)
		}
	}
	return s.cfg
}

func (s *serviceProvider) HTTPClient() *http.Client {
	if s.httpClient == nil {
		s.httpClient = &http.Client{Timeout: 45 * time.Second}
	}
	return s.httpClient
}

func (s *serviceProvider) TgBot(ctx context.Context) *tele.Bot {
	if s.tgBot == nil {
		cfg := s.Config().Tg()
		bot, err := tele.NewBot(tele.Settings{
			Token:     cfg.Token(),
			Poller:    &tele.LongPoller{Timeout: cfg.Timeout()},
			ParseMode: tele.ModeHTML,
		})
		if err != nil {
			log.Fatalf("failed to create telegram bot: %s", err.Error())
		}
		s.tgBot = bot
	}
	return s.tgBot
}

func (s *serviceProvider) GeminiClient(ctx context.Context) *gemini.Client {
	if s.geminiClient == nil {
		cfg := s.Config().Gemini()
		s.geminiClient = gemini.NewClient(s.HTTPClient(), cfg.APIKey(), cfg.Model())
	}
	return s.geminiClient
}

func (s *serviceProvider) AIService(ctx context.Context) *ai.Service {
	if s.aiService == nil {
		s.aiService = ai.NewService(s.GeminiClient(ctx))
	}
	return s.aiService
}

func (s *serviceProvider) TickTickClient(ctx context.Context) *ticktickclient.Client {
	if s.ticktickClient == nil {
		cfg := s.Config().TickTick()
		s.ticktickClient = ticktickclient.NewClient(s.HTTPClient(), cfg.BaseURL(), cfg.AccessToken())
	}
	return s.ticktickClient
}

func (s *serviceProvider) TickTickService(ctx context.Context) *ticktick.Service {
	if s.ticktickSvc == nil {
		s.ticktickSvc = ticktick.NewService(s.TickTickClient(ctx), s.Config().TickTick().DefaultProjectID())
	}
	return s.ticktickSvc
}

func (s *serviceProvider) TgBotService(ctx context.Context) *tgbot.BotService {
	if s.tgBotService == nil {
		s.tgBotService = tgbot.NewBotService(
			s.AIService(ctx),
			s.TickTickService(ctx),
			s.Config().App().Timezone(),
		)
	}
	return s.tgBotService
}
