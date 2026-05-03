package config

import (
	"errors"
	"net"
	"os"
)

const (
	httpHostEnvName    = "HTTP_HOST"
	httpPortEnvName    = "HTTP_PORT"
	httpLogEnvName     = "HTTP_LOG"
	swaggerHostEnvName = "SWAGGER_HOST"
)

type HTTPConfig interface {
	Address() string
	IsLog() bool
	SwaggerHost() string
}

type httpConfig struct {
	host        string
	port        string
	log         bool
	swaggerHost string
}

func newHTTPConfigEnv() (HTTPConfig, error) {
	host := os.Getenv(httpHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("http host not found")
	}

	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("http port not found")
	}

	log := os.Getenv(httpLogEnvName) == "true"

	swaggerHost := os.Getenv(swaggerHostEnvName)
	if len(swaggerHost) == 0 {
		// Default to localhost:port if not specified
		swaggerHost = net.JoinHostPort("localhost", port)
	}

	return &httpConfig{
		host:        host,
		port:        port,
		log:         log,
		swaggerHost: swaggerHost,
	}, nil
}

func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *httpConfig) IsLog() bool {
	return cfg.log
}

func (cfg *httpConfig) SwaggerHost() string {
	return cfg.swaggerHost
}
