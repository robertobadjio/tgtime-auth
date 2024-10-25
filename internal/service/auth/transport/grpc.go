package transport

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"

	"github.com/robertobadjio/tgtime-auth/internal/service/auth/endpoints"
	"github.com/robertobadjio/tgtime-auth/pkg/api/auth_v1"
)

type grpcServer struct {
	login           grpctransport.Handler
	getRefreshToken grpctransport.Handler
	getAccessToken  grpctransport.Handler
	auth_v1.UnimplementedAuthV1Server
}

// NewGRPCServer ???
func NewGRPCServer(ep endpoints.Set) auth_v1.AuthV1Server {
	return &grpcServer{
		login: grpctransport.NewServer(
			ep.LoginEndpoint,
			decodeGRPCLoginRequest,
			encodeGRPCLoginResponse,
		),
		getRefreshToken: grpctransport.NewServer(
			ep.GetRefreshTokenEndpoint,
			decodeGRPCGetRefreshTokenRequest,
			encodeGRPCGetRefreshTokenResponse,
		),
		getAccessToken: grpctransport.NewServer(
			ep.GetAccessTokenEndpoint,
			decodeGRPCGetAccessTokenRequest,
			encodeGRPCGetAccessTokenResponse,
		),
	}
}

func (g *grpcServer) Login(
	ctx context.Context,
	r *auth_v1.LoginRequest,
) (*auth_v1.LoginResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth_v1.LoginResponse), nil
}

func decodeGRPCLoginRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth_v1.LoginRequest)

	return endpoints.LoginRequest{Login: req.Username}, nil
}

func encodeGRPCLoginResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(endpoints.LoginResponse)

	return &auth_v1.LoginResponse{RefreshToken: res.RefreshToken}, nil
}

func (g *grpcServer) GetRefreshToken(
	ctx context.Context,
	r *auth_v1.GetRefreshTokenRequest,
) (*auth_v1.GetRefreshTokenResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth_v1.GetRefreshTokenResponse), nil
}

func decodeGRPCGetRefreshTokenRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth_v1.GetRefreshTokenRequest)

	return endpoints.GetRefreshTokenRequest{RefreshToken: req.RefreshToken}, nil
}

func encodeGRPCGetRefreshTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(endpoints.GetRefreshTokenResponse)

	return &auth_v1.GetRefreshTokenResponse{RefreshToken: res.RefreshToken}, nil
}

func (g *grpcServer) GetAccessToken(
	ctx context.Context,
	r *auth_v1.GetAccessTokenRequest,
) (*auth_v1.GetAccessTokenResponse, error) {
	_, resp, err := g.login.ServeGRPC(ctx, r)
	if err != nil {
		return nil, err
	}
	return resp.(*auth_v1.GetAccessTokenResponse), nil
}

func decodeGRPCGetAccessTokenRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*auth_v1.GetAccessTokenRequest)

	return endpoints.GetAccessTokenRequest{RefreshToken: req.RefreshToken}, nil
}

func encodeGRPCGetAccessTokenResponse(_ context.Context, response interface{}) (interface{}, error) {
	res := response.(endpoints.GetRefreshTokenResponse)

	return &auth_v1.GetAccessTokenResponse{AccessToken: res.RefreshToken}, nil
}
