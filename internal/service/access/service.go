package access

import (
	"context"
)

// Service Интерфейс сервиса
type Service interface {
	Check(ctx context.Context, endpointAddress string) error
}
