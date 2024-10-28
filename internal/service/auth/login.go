package auth

import (
	"context"
	"fmt"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/model"
)

// Login ???
func (s *service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.userRepo.GetUser(ctx, email)
	if err != nil {
		return "", fmt.Errorf("failed to get user: %w", err)
	}

	if !helper.VerifyPassword(user.Password, password) {
		return "", fmt.Errorf("failed to verify password: %w", err)
	}

	refreshToken, err := helper.GenerateToken(model.UserInfo{
		Email: email,
		Role:  user.Role,
	},
		[]byte(s.token.RefreshTokenSecretKey()),
		s.token.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return refreshToken, nil
}
