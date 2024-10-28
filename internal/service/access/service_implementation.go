package access

import (
	"github.com/robertobadjio/tgtime-auth/internal/config"
)

type service struct {
	token config.Token
}

// NewService ???
func NewService(token config.Token) Service {
	return &service{token: token}
}
