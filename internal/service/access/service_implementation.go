package access

import (
	"github.com/robertobadjio/tgtime-auth/internal/config"
	"github.com/robertobadjio/tgtime-auth/internal/repository/access"
)

type service struct {
	token      config.Token
	accessRepo access.Repository
}

// NewService ???
func NewService(token config.Token, accessRepo access.Repository) Service {
	return &service{token: token, accessRepo: accessRepo}
}
