package auth

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/robertobadjio/platform-common/pkg/sys"
	"github.com/robertobadjio/platform-common/pkg/sys/codes"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/model"
)

const spanGetRefreshTokenOperationName = "auth.GetRefreshToken"

// GetRefreshToken ???
func (s *service) GetRefreshToken(ctx context.Context, rt string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, spanGetRefreshTokenOperationName)
	defer span.Finish()

	claims, err := helper.VerifyToken(rt, []byte(s.token.RefreshTokenSecretKey()))
	if err != nil {
		return "", sys.NewCommonError("invalid refresh token", codes.Aborted)
	}

	user, err := s.userRepo.GetUser(ctx, claims.Email)
	if err != nil {
		return "", sys.NewCommonError("failed to get user", codes.Internal)
	}

	refreshToken, err := helper.GenerateToken(model.UserInfo{
		Email: user.Email,
		Role:  user.Role,
	},
		[]byte(s.token.RefreshTokenSecretKey()),
		s.token.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", sys.NewCommonError("failed to generate refresh token", codes.Aborted)
	}

	return refreshToken, nil
}
