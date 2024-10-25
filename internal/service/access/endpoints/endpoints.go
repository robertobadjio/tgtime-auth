package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/robertobadjio/tgtime-auth/internal/service/access"
)

// Set ???
type Set struct {
	CheckEndpoint endpoint.Endpoint
}

// NewEndpointSet ???
func NewEndpointSet(svc access.Service) Set {
	return Set{
		CheckEndpoint: MakeCheckEndpoint(svc),
	}
}

// MakeCheckEndpoint ???
func MakeCheckEndpoint(svc access.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CheckRequest)
		err := svc.Check(ctx, req.EndpointAddress)

		return nil, err
	}
}
