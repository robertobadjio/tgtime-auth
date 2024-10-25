package endpoints

import (
	"context"

	"github.com/go-kit/kit/endpoint"

	"github.com/robertobadjio/tgtime-auth/internal/service/service"
)

// Set ???
type Set struct {
	ServiceStatusEndpoint endpoint.Endpoint
}

// NewEndpointSet ???
func NewEndpointSet(svc service.Service) Set {
	return Set{
		ServiceStatusEndpoint: MakeServiceStatusEndpoint(svc),
	}
}

// MakeServiceStatusEndpoint ???
func MakeServiceStatusEndpoint(svc service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_ = request.(ServiceStatusRequest)
		code := svc.ServiceStatus(ctx)
		return ServiceStatusResponse{Code: code}, nil
	}
}
