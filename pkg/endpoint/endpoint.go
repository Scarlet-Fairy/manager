package endpoint

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

type ManagerEndpoint struct {
	DeployEndpoint  endpoint.Endpoint
	DestroyEndpoint endpoint.Endpoint
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

	return ManagerEndpoint{
		DeployEndpoint:  deployEndpoint,
		DestroyEndpoint: destroyEndpoint,
	}
}

// compile time assertions for our response types implementing endpoint.Failer.
var (
	_ endpoint.Failer = DeployResponse{}
	_ endpoint.Failer = DestroyResponse{}
)

type DeployRequest struct {
	GitRepo string
	Name    string
	Envs    map[string]string
}

type DeployResponse struct {
	Deploy *service.Deploy
	Err    error `json:"-"`
}

func (r DeployResponse) Failed() error {
	return r.Err
}

func makeDeployEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(*DeployRequest)
		deploy, err := s.Deploy(ctx, req.GitRepo, req.Name, req.Envs)

		return &DeployResponse{
			Deploy: deploy,
			Err:    err,
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
