package config

import (
	"errors"
	"os"
)

type TickTickConfig interface {
	BaseURL() string
	AccessToken() string
	DefaultProjectID() string
}

const (
	ticktickBaseURLEnvName          = "TICKTICK_BASE_URL"
	ticktickAccessTokenEnvName      = "TICKTICK_ACCESS_TOKEN"
	ticktickDefaultProjectIDEnvName = "TICKTICK_DEFAULT_PROJECT_ID"
)

type ticktickConfig struct {
	baseURL          string
	accessToken      string
	defaultProjectID string
}

func newTickTickConfigEnv() (TickTickConfig, error) {
	accessToken := os.Getenv(ticktickAccessTokenEnvName)
	if accessToken == "" {
		return nil, errors.New(ticktickAccessTokenEnvName + " is not set")
	}

	baseURL := os.Getenv(ticktickBaseURLEnvName)
	if baseURL == "" {
		baseURL = "https://api.ticktick.com/open/v1"
	}

	return &ticktickConfig{
		baseURL:          baseURL,
		accessToken:      accessToken,
		defaultProjectID: os.Getenv(ticktickDefaultProjectIDEnvName),
	}, nil
}

func (c *ticktickConfig) BaseURL() string {
	return c.baseURL
}

func (c *ticktickConfig) AccessToken() string {
	return c.accessToken
}

func (c *ticktickConfig) DefaultProjectID() string {
	return c.defaultProjectID
}
