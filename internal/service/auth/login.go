package auth

import (
	"context"
	"net/mail"

	"github.com/opentracing/opentracing-go"
	"github.com/robertobadjio/platform-common/pkg/sys"
	"github.com/robertobadjio/platform-common/pkg/sys/codes"
	"github.com/robertobadjio/platform-common/pkg/sys/validate"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/model"
)

const spanLoginOperationName = "auth.Login"

// Login ???
func (s *service) Login(ctx context.Context, email string, password string) (string, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, spanLoginOperationName)
	defer span.Finish()

	span.SetTag("email", email)

	err := validate.Validate(
		ctx,
		validateEmail(email),
		validatePassword(password),
	)
	if err != nil {
		return "", err
	}

	user, err := s.userRepo.GetUser(ctx, email)
	if err != nil {
		return "", sys.NewCommonError("failed to get user", codes.Internal)
	}

	if !helper.VerifyPassword(user.Password, password) {
		return "", sys.NewCommonError("failed to verify password", codes.PermissionDenied)
	}

	refreshToken, err := helper.GenerateToken(model.UserInfo{
		Email: email,
		Role:  user.Role,
	},
		[]byte(s.token.RefreshTokenSecretKey()),
		s.token.RefreshTokenExpiration(),
	)
	if err != nil {
		return "", sys.NewCommonError("failed to generate refresh token", codes.Internal)
	}

	return refreshToken, nil
}

func validateEmail(email string) validate.Condition {
	return func(ctx context.Context) error {
		_, err := mail.ParseAddress(email)
		if err != nil {
			return validate.NewValidationErrors("email not valid")
		}

		return nil
	}
}

func validatePassword(password string) validate.Condition {
	return func(ctx context.Context) error {
		if len(password) == 0 {
			return validate.NewValidationErrors("password not valid")
		}

		return nil
	}
}
