package service

import (
	"context"
	"errors"
	"log/slog"
	"ticktick-ai/internal/domain"
)

type IAIService interface {
	ParseText(ctx context.Context, text string, timezone string) (domain.ParsedIntent, error)
	ParseAudio(ctx context.Context, audio []byte, mimeType string, timezone string) (domain.ParsedIntent, error)
}

type ITaskService interface {
	ExecuteIntent(ctx context.Context, intent domain.ParsedIntent, timezone string) (domain.TaskResult, error)
}

type BotService struct {
	ai       IAIService
	tasks    ITaskService
	timezone string
}

func NewBotService(ai IAIService, tasks ITaskService, timezone string) *BotService {
	return &BotService{
		ai:       ai,
		tasks:    tasks,
		timezone: timezone,
	}
}

func (s *BotService) ProcessText(ctx context.Context, text string) string {
	if text == "" {
		return "Напиши задачу текстом или отправь голосовое сообщение."
	}

	intent, err := s.ai.ParseText(ctx, text, s.timezone)
	if err != nil {
		slog.ErrorContext(ctx, "parse text failed", "err", err, "input", text)
		return "Не понял запрос. Попробуй переформулировать."
	}

	return s.execute(ctx, intent)
}

func (s *BotService) ProcessAudio(ctx context.Context, audio []byte, mimeType string) string {
	if len(audio) == 0 {
		return "Не смог скачать голосовое сообщение. Попробуй еще раз."
	}

	intent, err := s.ai.ParseAudio(ctx, audio, mimeType, s.timezone)
	if err != nil {
		slog.ErrorContext(ctx, "parse audio failed", "err", err)
		return "Не понял голосовое сообщение. Попробуй сказать короче или отправь текстом."
	}


	return s.execute(ctx, intent)
}

func (s *BotService) execute(ctx context.Context, intent domain.ParsedIntent) string {
	result, err := s.tasks.ExecuteIntent(ctx, intent, s.timezone)
	if err != nil {
		slog.ErrorContext(ctx, "execute intent failed", "err", err, "intent_type", intent.Type, "task_title", intent.TaskTitle)
		return userMessageFromError(err)
	}

	if result.Message != "" {
		return result.Message
	}
	return "Готово."
}

func userMessageFromError(err error) string {
	if errors.Is(err, context.DeadlineExceeded) {
		return "Запрос занял слишком много времени. Попробуй еще раз."
	}
	return "TickTick сейчас не принял запрос. Попробуй еще раз."
}
