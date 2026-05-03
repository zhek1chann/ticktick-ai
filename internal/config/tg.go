package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type TgConfig interface {
	Token() string
	Timeout() time.Duration
	AdminIDs() []int64
}

const (
	tgTokenEnvName    = "TG_TOKEN"
	tgTimeoutEnvName  = "TG_TIMEOUT"
	tgAdminIDsEnvName = "TG_ADMIN_IDS"
)

type tgConfig struct {
	token    string
	timeout  time.Duration
	adminIDs []int64
}

func newTgConfigEnv() (TgConfig, error) {
	token := os.Getenv(tgTokenEnvName)
	if token == "" {
		return nil, errors.New(tgTokenEnvName + " is not set")
	}

	timeoutValue := os.Getenv(tgTimeoutEnvName)
	if timeoutValue == "" {
		timeoutValue = "10s"
	}

	timeout, err := time.ParseDuration(timeoutValue)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", tgTimeoutEnvName, err)
	}

	var adminIDs []int64
	for _, s := range strings.Split(os.Getenv(tgAdminIDsEnvName), ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid %s value %q: %w", tgAdminIDsEnvName, s, err)
		}
		adminIDs = append(adminIDs, id)
	}

	return &tgConfig{
		token:    token,
		timeout:  timeout,
		adminIDs: adminIDs,
	}, nil
}

func (c *tgConfig) Token() string {
	return c.token
}

func (c *tgConfig) Timeout() time.Duration {
	return c.timeout
}

func (c *tgConfig) AdminIDs() []int64 {
	return c.adminIDs
}
