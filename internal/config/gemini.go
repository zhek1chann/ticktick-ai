package config

import (
	"errors"
	"os"
)

type GeminiConfig interface {
	APIKey() string
	Model() string
}

const (
	geminiAPIKeyEnvName = "GEMINI_API_KEY"
	geminiModelEnvName  = "GEMINI_MODEL"
)

type geminiConfig struct {
	apiKey string
	model  string
}

func newGeminiConfigEnv() (GeminiConfig, error) {
	apiKey := os.Getenv(geminiAPIKeyEnvName)
	if apiKey == "" {
		return nil, errors.New(geminiAPIKeyEnvName + " is not set")
	}

	model := os.Getenv(geminiModelEnvName)
	if model == "" {
		model = "gemini-1.5-flash"
	}

	return &geminiConfig{
		apiKey: apiKey,
		model:  model,
	}, nil
}

func (c *geminiConfig) APIKey() string {
	return c.apiKey
}

func (c *geminiConfig) Model() string {
	return c.model
}
