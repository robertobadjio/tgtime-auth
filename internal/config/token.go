package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

const (
	// #nosec G101
	refreshTokenSecretKeyEnvName = "REFRESH_TOKEN_SECRET_KEY"
	// #nosec G101
	accessTokenSecretKeyEnvName = "ACCESS_TOKEN_SECRET_KEY"
	// #nosec G101
	refreshTokenExpirationEnvName = "REFRESH_TOKEN_EXPIRATION_IN_MINUTES"
	// #nosec G101
	accessTokenExpirationEnvName = "ACCESS_TOKEN_EXPIRATION_IN_MINUTES"
)

// Token ???
type Token interface {
	AccessTokenSecretKey() string
	RefreshTokenSecretKey() string
	RefreshTokenExpiration() time.Duration
	AccessTokenExpiration() time.Duration
}

type token struct {
	refreshTokenSecretKey  string
	accessTokenSecretKey   string
	refreshTokenExpiration time.Duration
	accessTokenExpiration  time.Duration
}

// NewToken ???
func NewToken() (Token, error) {
	refreshTokenSecretKey := os.Getenv(refreshTokenSecretKeyEnvName)
	if len(refreshTokenSecretKey) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", refreshTokenSecretKeyEnvName)
	}

	accessTokenSecretKey := os.Getenv(accessTokenSecretKeyEnvName)
	if len(accessTokenSecretKey) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", accessTokenSecretKeyEnvName)
	}

	refreshTokenExpiration := os.Getenv(refreshTokenExpirationEnvName)
	if len(refreshTokenExpiration) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", refreshTokenExpirationEnvName)
	}

	accessTokenExpiration := os.Getenv(accessTokenExpirationEnvName)
	if len(accessTokenExpiration) == 0 {
		return nil, fmt.Errorf("environment variable %s must be set", accessTokenExpirationEnvName)
	}

	refreshTokenExpirationInt, err := strconv.Atoi(refreshTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s must be an integer", refreshTokenExpirationEnvName)
	}

	accessTokenExpirationInt, err := strconv.Atoi(accessTokenExpiration)
	if err != nil {
		return nil, fmt.Errorf("environment variable %s must be an integer", accessTokenExpirationEnvName)
	}

	return &token{
		refreshTokenSecretKey:  refreshTokenSecretKey,
		accessTokenSecretKey:   accessTokenSecretKey,
		refreshTokenExpiration: time.Duration(refreshTokenExpirationInt) * time.Minute,
		accessTokenExpiration:  time.Duration(accessTokenExpirationInt) * time.Minute,
	}, nil
}

// AccessTokenSecretKey ???
func (t *token) AccessTokenSecretKey() string {
	return t.refreshTokenSecretKey
}

// RefreshTokenSecretKey ???
func (t *token) RefreshTokenSecretKey() string {
	return t.accessTokenSecretKey
}

// RefreshTokenExpiration ???
func (t *token) RefreshTokenExpiration() time.Duration {
	return t.refreshTokenExpiration
}

// AccessTokenExpiration ???
func (t *token) AccessTokenExpiration() time.Duration {
	return t.accessTokenExpiration
}
