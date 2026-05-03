package config

import (
	"errors"
	"os"
	"strconv"
	"time"
)

type JWTConfig interface {
	SecretKey() string
	AccessTokenTTL() time.Duration
	RefreshTokenTTL() time.Duration
}

const (
	jwtSecretKeyEnvName       = "JWT_SECRET_KEY"
	jwtAccessTokenTTLEnvName  = "JWT_ACCESS_TOKEN_TTL"
	jwtRefreshTokenTTLEnvName = "JWT_REFRESH_TOKEN_TTL"
)

type jwtConfig struct {
	secretKey       string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func newJWTConfigEnv() (JWTConfig, error) {
	secretKey := os.Getenv(jwtSecretKeyEnvName)
	if secretKey == "" {
		return nil, errors.New(jwtSecretKeyEnvName + " not found")
	}

	accessTokenTTLStr := os.Getenv(jwtAccessTokenTTLEnvName)
	if accessTokenTTLStr == "" {
		return nil, errors.New(jwtAccessTokenTTLEnvName + " not found")
	}

	accessTokenTTLMinutes, err := strconv.Atoi(accessTokenTTLStr)
	if err != nil {
		return nil, errors.New("invalid value for " + jwtAccessTokenTTLEnvName)
	}

	refreshTokenTTLStr := os.Getenv(jwtRefreshTokenTTLEnvName)
	if refreshTokenTTLStr == "" {
		return nil, errors.New(jwtRefreshTokenTTLEnvName + " not found")
	}

	refreshTokenTTLHours, err := strconv.Atoi(refreshTokenTTLStr)
	if err != nil {
		return nil, errors.New("invalid value for " + jwtRefreshTokenTTLEnvName)
	}

	return &jwtConfig{
		secretKey:       secretKey,
		accessTokenTTL:  time.Duration(accessTokenTTLMinutes) * time.Minute,
		refreshTokenTTL: time.Duration(refreshTokenTTLHours) * time.Hour,
	}, nil
}

func (c *jwtConfig) SecretKey() string {
	return c.secretKey
}

func (c *jwtConfig) AccessTokenTTL() time.Duration {
	return c.accessTokenTTL
}

func (c *jwtConfig) RefreshTokenTTL() time.Duration {
	return c.refreshTokenTTL
}
