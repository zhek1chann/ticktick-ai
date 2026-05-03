package config

import (
	"errors"
	"log/slog"
	"os"
	"strconv"
)

type Mode string

const (
	Development Mode = "dev"
	Production  Mode = "prod"
	Test        Mode = "test"
)

type AppConfig interface {
	Mode() Mode
	Debug() bool
	LogLevel() slog.Level
	Timezone() string
}

const (
	modeEnvName     = "APP_MODE"
	debugEnvName    = "APP_DEBUG"
	logLevelEnvName = "APP_LOG_LEVEL"
	timezoneEnvName = "USER_TIMEZONE"
)

type appConfig struct {
	mode     Mode
	debug    bool
	logLevel slog.Level
	timezone string
}

func newAppConfigEnv() (AppConfig, error) {
	modeValue := os.Getenv(modeEnvName)
	if modeValue == "" {
		modeValue = string(Development)
	}

	var mode Mode
	switch modeValue {
	case string(Development):
		mode = Development
	case string(Production):
		mode = Production
	case string(Test):
		mode = Test
	default:
		return nil, errors.New("invalid value for " + modeEnvName)
	}

	debugValue := os.Getenv(debugEnvName)
	var debug bool
	if debugValue != "" {
		var err error
		debug, err = strconv.ParseBool(debugValue)
		if err != nil {
			return nil, errors.New("invalid boolean value for " + debugEnvName)
		}
	} else {
		// Default to false if not provided
		debug = false
	}

	logLevelStr := os.Getenv(logLevelEnvName)

	if logLevelStr == "" {
		logLevelStr = "info"
	}
	var logLevel slog.Level
	switch logLevelStr {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		return nil, errors.New("invalid value for " + logLevelEnvName)
	}

	timezone := os.Getenv(timezoneEnvName)
	if timezone == "" {
		timezone = "Asia/Almaty"
	}

	return &appConfig{
		mode:     mode,
		debug:    debug,
		logLevel: logLevel,
		timezone: timezone,
	}, nil
}

func (c *appConfig) Mode() Mode {
	return c.mode
}

func (c *appConfig) Debug() bool {
	return c.debug
}

func (c *appConfig) LogLevel() slog.Level {
	return c.logLevel
}

func (c *appConfig) Timezone() string {
	return c.timezone
}
