package config

import (
	"fmt"
	"os"
)

const (
	dbNameEnvName     = "DATABASE_PG_NAME"
	dbHostEnvName     = "DATABASE_PG_HOST"
	dbPortEnvName     = "DATABASE_PG_PORT"
	dbUserEnvName     = "DATABASE_PG_USER"
	dbPasswordEnvName = "DATABASE_PG_PASSWORD"
	dbSSLEnvName      = "DATABASE_PG_SSL_MODE"
)

// PGConfig Конфиг для подключения к DB PostgresQl
type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

// NewPGConfig Конструктор конфига для подключения к DB PostgresQl
func NewPGConfig() (PGConfig, error) {
	host := os.Getenv(dbHostEnvName)
	if len(host) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dbHostEnvName)
	}

	port := os.Getenv(dbPortEnvName)
	if len(port) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dbPortEnvName)
	}

	dbName := os.Getenv(dbNameEnvName)
	if len(dbName) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dbNameEnvName)
	}

	user := os.Getenv(dbUserEnvName)
	if len(user) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dbUserEnvName)
	}

	password := os.Getenv(dbPasswordEnvName)
	if len(password) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dbPasswordEnvName)
	}

	sslMode := os.Getenv(dbSSLEnvName)
	if len(sslMode) == 0 {
		return nil, fmt.Errorf("environment variable %s is not set", dbSSLEnvName)
	}

	dsn := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		host, port, dbName, user, password, sslMode,
	)

	return &pgConfig{
		dsn: dsn,
	}, nil
}

// DSN Возвращает DSN для подключения
func (cfg *pgConfig) DSN() string {
	return cfg.dsn
}
