package access

import (
	"context"
)

// Repository ???
type Repository interface {
	GetAccessibleRolesByEndpoint(ctx context.Context, endpoint string) ([]string, error)
}
