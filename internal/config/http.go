package config

import (
	"fmt"
	"net"
	"os"
)

const (
	httpPortEnvName = "HTTP_PORT"
)

// HTTPConfig Конфиг HTTP-интерфейса API
type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

// NewHTTPConfig Конструктор конфига HTTP-интерфейса API
func NewHTTPConfig() (HTTPConfig, error) {
	port := os.Getenv(httpPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", httpPortEnvName)
	}

	return &httpConfig{
		host: "",
		port: port,
	}, nil
}

// Address Возвращает адрес для подключения
func (cfg *httpConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
