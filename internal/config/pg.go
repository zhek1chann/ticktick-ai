package config

import (
	"fmt"
	"os"
	"strings"
)

const (
	pgDSNEnvName = "PG_DSN"
	pgLogEnvName = "PG_LOG"
)

// PGConfig defines the interface for PostgreSQL configuration
type PGConfig interface {
	DSN() string
	IsLog() bool
}

type pgConfig struct {
	dsn    string
	silent bool
	isLog  bool
}

// DSN returns the Data Source Name for connecting to the database
func (c *pgConfig) DSN() string {
	return c.dsn
}

func (c *pgConfig) IsLog() bool {
	return c.isLog
}

// newPGConfigEnv creates a new PostgreSQL config from environment variables
func newPGConfigEnv() (PGConfig, error) {
	dsn := os.Getenv(pgDSNEnvName)
	if dsn == "" {
		return nil, fmt.Errorf("env variable %s is not set", pgDSNEnvName)
	}

	// Strip quotes if present (Docker --env-file includes them from .env file)
	dsn = strings.Trim(dsn, `"'`)

	isLog := os.Getenv(pgLogEnvName) == "true"

	fmt.Println("dsn", dsn)
	fmt.Println("isLog", isLog)

	return &pgConfig{
		dsn:   strings.ReplaceAll(dsn, ";", " "),
		isLog: isLog,
	}, nil
}
