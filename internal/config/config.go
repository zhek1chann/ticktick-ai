package config

import (
	"github.com/joho/godotenv"
)

type Config struct {
	pg       PGConfig
	app      AppConfig
	http     HTTPConfig
	s3       S3Config
	jwt      JWTConfig
	tg       TgConfig
	gemini   GeminiConfig
	ticktick TickTickConfig
}

func (c *Config) PG() PGConfig {
	return c.pg
}

func (c *Config) App() AppConfig {
	return c.app
}

func (c *Config) HTTP() HTTPConfig {
	return c.http
}

func (c *Config) S3() S3Config {
	return c.s3
}

func (c *Config) JWT() JWTConfig {
	return c.jwt
}

func (c *Config) Tg() TgConfig {
	return c.tg
}

func (c *Config) Gemini() GeminiConfig {
	return c.gemini
}

func (c *Config) TickTick() TickTickConfig {
	return c.ticktick
}

func LoadConfig(path string) (*Config, error) {
	_ = godotenv.Load(path)
	return newEnvProvider()
}

func newEnvProvider() (*Config, error) {
	appConfig, err := newAppConfigEnv()
	if err != nil {
		return nil, err
	}

	tgConfig, err := newTgConfigEnv()
	if err != nil {
		return nil, err
	}

	geminiConfig, err := newGeminiConfigEnv()
	if err != nil {
		return nil, err
	}

	ticktickConfig, err := newTickTickConfigEnv()
	if err != nil {
		return nil, err
	}

	return &Config{
		app:      appConfig,
		tg:       tgConfig,
		gemini:   geminiConfig,
		ticktick: ticktickConfig,
	}, nil
}
