package user

import (
	"context"

	"github.com/robertobadjio/tgtime-auth/internal/repository/user/model"
)

// Repository ???
type Repository interface {
	GetUser(ctx context.Context, email string) (*model.User, error)
}
