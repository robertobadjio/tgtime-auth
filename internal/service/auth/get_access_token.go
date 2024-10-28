package auth

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/model"
)

func (s *service) GetAccessToken(ctx context.Context, rt string) (string, error) {
	claims, err := helper.VerifyToken(rt, []byte(s.token.RefreshTokenSecretKey()))
	if err != nil {
		return "", status.Errorf(codes.Aborted, "invalid refresh token")
	}

	user, err := s.userRepo.GetUser(ctx, claims.Email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	accessToken, err := helper.GenerateToken(model.UserInfo{
		Email: claims.Email,
		Role:  user.Role,
	},
		[]byte(s.token.AccessTokenSecretKey()),
		s.token.AccessTokenExpiration(),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
