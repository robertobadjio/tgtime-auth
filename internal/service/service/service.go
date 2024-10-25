package service

import (
	"context"
)

// Service Интерфейс сервиса
type Service interface {
	ServiceStatus(ctx context.Context) int
}
