package access

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/metadata"

	"github.com/robertobadjio/tgtime-auth/internal/helper"
	"github.com/robertobadjio/tgtime-auth/internal/service/auth/model"
)

const authPrefix = "Bearer "

var accessibleRoles map[string]string

func (s *service) Check(ctx context.Context, endpointAddress string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return errors.New("metadata is not provided")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return errors.New("authorization header is not provided")
	}

	if !strings.HasPrefix(authHeader[0], authPrefix) {
		return errors.New("invalid authorization header format")
	}

	accessToken := strings.TrimPrefix(authHeader[0], authPrefix)

	claims, err := helper.VerifyToken(accessToken, []byte(s.token.AccessTokenSecretKey()))
	if err != nil {
		return errors.New("access token is invalid")
	}

	accessibleMap, err := s.accessibleRoles(ctx)
	if err != nil {
		return errors.New("failed to get accessible roles")
	}

	role, ok := accessibleMap[endpointAddress]
	if !ok {
		return errors.New("route not found")
	}

	if role == claims.Role {
		return nil
	}

	return errors.New("access denied")
}

// Возвращает мапу с адресом эндпоинта и ролью, которая имеет доступ к нему
func (s *service) accessibleRoles(_ context.Context) (map[string]string, error) {
	if accessibleRoles == nil {
		accessibleRoles = make(map[string]string)

		// TODO: Лезем в базу за данными о доступных ролях для каждого эндпоинта
		// Можно кэшировать данные, чтобы не лезть в базу каждый раз

		// Например, для эндпоинта /note_v1.NoteV1/Get доступна только роль admin
		accessibleRoles[model.ExamplePath] = "admin"
	}

	return accessibleRoles, nil
}
