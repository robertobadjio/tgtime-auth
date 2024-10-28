package config

import (
	"fmt"
	"net"
	"os"
)

const (
	promAppPortEnvName = "PROMETHEUS_APP_PORT"
)

// PromConfig ???
type PromConfig interface {
	Address() string
}

type promConfig struct {
	host string
	port string
}

// NewPromConfig ???
func NewPromConfig() (PromConfig, error) {
	port := os.Getenv(promAppPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", promAppPortEnvName)
	}

	return &promConfig{
		host: "",
		port: port,
	}, nil
}

// Address Возвращает адрес для подключения
func (cfg *promConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
