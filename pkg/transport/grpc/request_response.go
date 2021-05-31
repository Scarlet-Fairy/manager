package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pb"
	"github.com/Scarlet-Fairy/manager/pkg/endpoint"
)

func decodeDeployRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.DeployRequest)
	if req.Envs == nil {
		req.Envs = map[string]string{}
	}

	return &endpoint.DeployRequest{
		GitRepo: req.GitRepo,
		Name:    req.Name,
		Envs:    req.Envs,
	}, nil
}

func encodeDeployResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*endpoint.DeployResponse)

	return &pb.DeployResponse{
		DeployId: res.DeployId,
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

func decodeGetDeployRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetDeployRequest)

	return &endpoint.GetDeployRequest{
		Id: req.DeployId,
	}, nil
}

func encodeGetDeployResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*endpoint.GetDeployResponse)

	return &pb.GetDeployResponse{
		Deploy: coreDeployToTransportDeploy(res.Deploy),
	}, nil
}

func decodeListDeploysRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	return &endpoint.ListDeploysRequest{}, nil
}

func encodeListDeploysResponse(_ context.Context, resp interface{}) (interface{}, error) {
	res := resp.(*endpoint.ListDeploysResponse)

	var deploys []*pb.Deploy
	for _, deploy := range res.Deploys {
		deploys = append(deploys, coreDeployToTransportDeploy(deploy))
	}

	return &pb.ListDeploysResponse{
		Deploys: deploys,
	}, nil
}
