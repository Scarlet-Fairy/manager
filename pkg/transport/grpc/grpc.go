package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pb"
	"github.com/Scarlet-Fairy/manager/pkg/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	grpctransport "github.com/go-kit/kit/transport/grpc"
)

type grpcServer struct {
	pb.UnimplementedManagerServer
	deploy  grpctransport.Handler
	destroy grpctransport.Handler
}

func NewGRPCServer(endpoints endpoint.ManagerEndpoint, logger log.Logger) pb.ManagerServer {
	options := []grpctransport.ServerOption{
		grpctransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
	}

	return &grpcServer{
		deploy: grpctransport.NewServer(
			endpoints.DeployEndpoint,
			decodeDeployRequest,
			encodeDeployResponse,
			options...,
		),
		destroy: grpctransport.NewServer(
			endpoints.DestroyEndpoint,
			decodeDestroyRequest,
			encodeDestroyResponse,
			options...,
		),
	}
}

func (g grpcServer) Deploy(ctx context.Context, request *pb.DeployRequest) (*pb.DeployResponse, error) {
	_, resp, err := g.deploy.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.DeployResponse), nil
}

func (g grpcServer) Destroy(ctx context.Context, request *pb.DestroyRequest) (*pb.DestroyResponse, error) {
	_, resp, err := g.destroy.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.DestroyResponse), nil
}
