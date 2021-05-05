package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pb"
	"github.com/Scarlet-Fairy/manager/pkg/endpoint"
)

func decodeDeployRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeployRequest)

	return &endpoint.DeployRequest{
		GitRepo: req.GitRepo,
		Name:    req.Name,
		Envs:    req.Envs,
	}, nil
}

func encodeDeployResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*endpoint.DeployResponse)

	return &pb.DeployResponse{
		Deploy: coreDeployToTransportDeploy(res.Deploy),
	}, nil
}

func decodeDestroyRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DestroyRequest)

	return &endpoint.DestroyRequest{
		Id: req.DeployId,
	}, nil
}

func encodeDestroyResponse(_ context.Context, resp interface{}) (interface{}, error) {
	return &pb.DestroyResponse{}, nil
}
