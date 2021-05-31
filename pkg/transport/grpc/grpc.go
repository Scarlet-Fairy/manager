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
	deploy      grpctransport.Handler
	destroy     grpctransport.Handler
	getDeploy   grpctransport.Handler
	listDeploys grpctransport.Handler
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
		getDeploy: grpctransport.NewServer(
			endpoints.GetDeployEndpoint,
			decodeGetDeployRequest,
			encodeGetDeployResponse,
			options...,
		),
		listDeploys: grpctransport.NewServer(
			endpoints.ListDeployEndpoint,
			decodeListDeploysRequest,
			encodeListDeploysResponse,
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

func (g grpcServer) GetDeploy(ctx context.Context, request *pb.GetDeployRequest) (*pb.GetDeployResponse, error) {
	_, resp, err := g.getDeploy.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.GetDeployResponse), nil
}

func (g grpcServer) ListDeploys(ctx context.Context, request *pb.ListDeploysRequest) (*pb.ListDeploysResponse, error) {
	_, resp, err := g.listDeploys.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}

	return resp.(*pb.ListDeploysResponse), nil
}
