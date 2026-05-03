package ai

import (
	"context"
	"ticktick-ai/internal/domain"
)

type Parser interface {
	ParseText(ctx context.Context, text string, timezone string) (domain.ParsedIntent, error)
	ParseAudio(ctx context.Context, audio []byte, mimeType string, timezone string) (domain.ParsedIntent, error)
}

type Service struct {
	parser Parser
}

func NewService(parser Parser) *Service {
	return &Service{parser: parser}
}

func (s *Service) ParseText(ctx context.Context, text string, timezone string) (domain.ParsedIntent, error) {
	return s.parser.ParseText(ctx, text, timezone)
}

func (s *Service) ParseAudio(ctx context.Context, audio []byte, mimeType string, timezone string) (domain.ParsedIntent, error) {
	return s.parser.ParseAudio(ctx, audio, mimeType, timezone)
}
