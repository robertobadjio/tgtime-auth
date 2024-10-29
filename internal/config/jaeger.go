package config

import (
	"fmt"
	"net"
	"os"
)

const (
	jaegerHostEnvName = "JAEGER_CLIENT_HOST"
	jaegerPortEnvName = "JAEGER_CLIENT_PORT"
)

// JaegerConfig ???
type JaegerConfig interface {
	Address() string
}

type jaegerConfig struct {
	host string
	port string
}

// NewJaegerConfig Конструктор конфига
func NewJaegerConfig() (JaegerConfig, error) {
	host := os.Getenv(jaegerHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", jaegerHostEnvName)
	}

	port := os.Getenv(jaegerPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", jaegerPortEnvName)
	}

	return &jaegerConfig{
		host: host,
		port: port,
	}, nil
}

// Address Возвращает адрес для подключения
func (cfg *jaegerConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
