package transport

import (
	"context"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/robertobadjio/tgtime-auth/internal/service/access/endpoints"
	"github.com/robertobadjio/tgtime-auth/pkg/api/access_v1"
)

type grpcServer struct {
	check grpctransport.Handler
	access_v1.UnimplementedAccessV1Server
}

// NewGRPCServer ???
func NewGRPCServer(ep endpoints.Set) access_v1.AccessV1Server {
	return &grpcServer{
		check: grpctransport.NewServer(
			ep.CheckEndpoint,
			decodeGRPCCheckRequest,
			encodeGRPCCheckResponse,
		),
	}
}

func (g *grpcServer) Check(
	ctx context.Context,
	r *access_v1.CheckRequest,
) (*emptypb.Empty, error) {
	_, _, err := g.check.ServeGRPC(ctx, r)
	return &emptypb.Empty{}, err
}

func decodeGRPCCheckRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*access_v1.CheckRequest)

	return endpoints.CheckRequest{EndpointAddress: req.EndpointAddress}, nil
}

func encodeGRPCCheckResponse(_ context.Context, _ interface{}) (interface{}, error) {
	return &emptypb.Empty{}, nil
}
