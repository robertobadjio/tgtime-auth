package config

import (
	"fmt"
	"net"
	"os"
)

const (
	grpcPortEnvName = "GRPC_PORT"
)

// GRPCConfig Конфиг GRPC-интерфейса API
type GRPCConfig interface {
	Address() string
}

type grpcConfig struct {
	host string
	port string
}

// NewGRPCConfig Конструктор конфига GRPC-интерфейса API
func NewGRPCConfig() (GRPCConfig, error) {
	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", grpcPortEnvName)
	}

	return &grpcConfig{
		host: "",
		port: port,
	}, nil
}

// Address Возвращает адрес для подключения
func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}
