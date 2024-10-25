package auth

import (
	"context"
)

// Service Интерфейс сервиса
type Service interface {
	Login(ctx context.Context, login, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}
