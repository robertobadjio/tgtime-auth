package access

import (
	"context"
	"errors"
	"strings"

	"github.com/opentracing/opentracing-go"
	"github.com/robertobadjio/platform-common/pkg/sys"
	"github.com/robertobadjio/platform-common/pkg/sys/codes"
	"google.golang.org/grpc/metadata"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
)

const authPrefix = "Bearer "

func (s *service) Check(ctx context.Context, endpointAddress string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "access.Check")
	defer span.Finish()

	span.SetTag("endpointAddress", endpointAddress)

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return sys.NewCommonError("metadata is not provided", codes.Internal)
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return sys.NewCommonError("authorization header is not provided", codes.InvalidArgument)
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := helper.VerifyToken(accessToken, []byte(s.token.AccessTokenSecretKey()))
	if err != nil {
		return sys.NewCommonError("access denied", codes.PermissionDenied)
	}

	roles, err := s.accessRepo.GetAccessibleRolesByEndpoint(ctx, endpointAddress)
	for _, role := range roles {
		if claims.Role == role {
			return nil
		}
	}

	return sys.NewCommonError("access denied", codes.PermissionDenied)
}
