package endpoint

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type ManagerEndpoint struct {
	DeployEndpoint    endpoint.Endpoint
	DestroyEndpoint   endpoint.Endpoint
	GetDeployEndpoint endpoint.Endpoint
}

func NewEndpoint(s service.Service, logger log.Logger) ManagerEndpoint {
	var deployEndpoint endpoint.Endpoint
	{
		deployEndpoint = makeDeployEndpoint(s)
		deployEndpoint = LoggingMiddleware(log.With(logger, "method", "Deploy"))(deployEndpoint)
		deployEndpoint = UnwrapErrorMiddleware()(deployEndpoint)
	}

	var destroyEndpoint endpoint.Endpoint
	{
		destroyEndpoint = makeDestroyEndpoint(s)
		destroyEndpoint = LoggingMiddleware(log.With(logger, "method", "Destroy"))(destroyEndpoint)
		destroyEndpoint = UnwrapErrorMiddleware()(destroyEndpoint)
	}

	var getDeployEndpoint endpoint.Endpoint
	{
		getDeployEndpoint = makeGetDeployEndpoint(s)
		getDeployEndpoint = LoggingMiddleware(log.With(logger, "method", "GetDeploy"))(getDeployEndpoint)
		getDeployEndpoint = UnwrapErrorMiddleware()(getDeployEndpoint)
	}

	return ManagerEndpoint{
		DeployEndpoint:    deployEndpoint,
		DestroyEndpoint:   destroyEndpoint,
		GetDeployEndpoint: getDeployEndpoint,
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = DeployResponse{}
	_ endpoint.Failer = DestroyResponse{}
	_ endpoint.Failer = GetDeployResponse{}
)

type DeployRequest struct {
	GitRepo string
	Name    string
	Envs    map[string]string
}

type DeployResponse struct {
	DeployId string
	Err      error `json:"-"`
}

func (r DeployResponse) Failed() error {
	return r.Err
}

func makeDeployEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*DeployRequest)
		id, err := s.Deploy(ctx, req.GitRepo, req.Name, req.Envs)

		return &DeployResponse{
			DeployId: id,
			Err:      err,
		}, nil
	}
}

type DestroyRequest struct {
	Id string
}

type DestroyResponse struct {
	Err error `json:"-"`
}

func (r DestroyResponse) Failed() error {
	return r.Err
}

func makeDestroyEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*DestroyRequest)
		err := s.Destroy(ctx, req.Id)

		return &DestroyResponse{
			Err: err,
		}, nil
	}
}

type GetDeployRequest struct {
	Name string
}

type GetDeployResponse struct {
	Deploy *service.Deploy
	Err    error `json:"-"`
}

func (r GetDeployResponse) Failed() error {
	return r.Err
}

func makeGetDeployEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*GetDeployRequest)
		deploy, err := s.GetDeploy(ctx, req.Name)

		return &GetDeployResponse{
			Deploy: deploy,
			Err:    err,
		}, nil
	}
}
