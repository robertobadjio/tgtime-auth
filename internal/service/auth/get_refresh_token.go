package auth

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/model"
)

// GetRefreshToken ???
func (s *service) GetRefreshToken(ctx context.Context, rt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "auth.GetRefreshToken")
	defer span.Finish()

	claims, err := helper.VerifyToken(rt, []byte(s.token.RefreshTokenSecretKey()))
	if err != nil {
		return "", fmt.Errorf("invalid refresh token: %w", err)
	}

	user, err := s.userRepo.GetUser(ctx, claims.Email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	refreshToken, err := helper.GenerateToken(model.UserInfo{
		Email: user.Email,
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
