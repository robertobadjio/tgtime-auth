package access

import (
	"github.com/robertobadjio/tgtime-auth/internal/config"
	"github.com/robertobadjio/tgtime-auth/internal/repository/user"
)

type service struct {
	userRepo user.Repository
	token    config.Token
}

// NewService ???
func NewService(userRepo user.Repository, token config.Token) Service {
	return &service{userRepo: userRepo, token: token}
}
