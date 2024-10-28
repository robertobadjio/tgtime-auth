package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/robertobadjio/tgtime-auth/internal/service/auth"
)

// Set ???
type Set struct {
	LoginEndpoint           endpoint.Endpoint
	GetRefreshTokenEndpoint endpoint.Endpoint
	GetAccessTokenEndpoint  endpoint.Endpoint
}

// NewEndpointSet ???
func NewEndpointSet(svc auth.Service) Set {
	return Set{
		LoginEndpoint:           MakeLoginEndpoint(svc),
		GetRefreshTokenEndpoint: MakeGetRefreshTokenEndpoint(svc),
		GetAccessTokenEndpoint:  MakeGetAccessTokenEndpoint(svc),
	}
}

// MakeLoginEndpoint ???
func MakeLoginEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(LoginRequest)
		refreshToken, err := svc.Login(ctx, req.Email, req.Password)
		if err != nil {
			return nil, err
		}

		return LoginResponse{RefreshToken: refreshToken}, nil
	}
}

// MakeGetRefreshTokenEndpoint ???
func MakeGetRefreshTokenEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetRefreshTokenRequest)
		refreshToken, err := svc.GetRefreshToken(ctx, req.RefreshToken)
		if err != nil {
			return nil, err
		}
		return GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
	}
}

// MakeGetAccessTokenEndpoint ???
func MakeGetAccessTokenEndpoint(svc auth.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetAccessTokenRequest)
		accessToken, err := svc.GetAccessToken(ctx, req.RefreshToken)
		if err != nil {
			return nil, err
		}
		return GetAccessTokenResponse{AccessToken: accessToken}, nil
	}
}
