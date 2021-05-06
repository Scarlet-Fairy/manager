package repository

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/log"
)

type Middleware func(repository service.Repository) service.Repository

type repositoryLogger struct {
	next   service.Repository
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(repository service.Repository) service.Repository {
		return &repositoryLogger{
			next:   repository,
			logger: logger,
		}
	}
}

func (r repositoryLogger) CreateDeploy(ctx context.Context, deploy *service.Deploy) (id string, err error) {
	defer func() {
		r.logger.Log(
			"method", "CreateDeploy",
			"deploy", deploy,
			"id", id,
			"err", err,
		)
	}()

	return r.next.CreateDeploy(ctx, deploy)
}

func (r repositoryLogger) GetDeploy(ctx context.Context, id string) (deploy *service.Deploy, err error) {
	defer func() {
		r.logger.Log(
			"method", "GetDeploy",
			"id", id,
			"deploy", deploy,
			"err", err,
		)
	}()

	return r.next.GetDeploy(ctx, id)
}

func (r repositoryLogger) ListDeploy(ctx context.Context) (deploys []*service.Deploy, err error) {
	defer func() {
		r.logger.Log(
			"method", "ListDeploy",
			"deploys", deploys,
			"err", err,
		)
	}()

	return r.next.ListDeploy(ctx)
}

func (r repositoryLogger) UpdateDeploy(ctx context.Context, deploy *service.Deploy) (err error) {
	defer func() {
		r.logger.Log(
			"method", "UpdateDeploy",
			"deploy", deploy,
			"err", err,
		)
	}()

	return r.next.UpdateDeploy(ctx, deploy)
}

func (r repositoryLogger) DeleteDeploy(ctx context.Context, id string) (err error) {
	defer func() {
		r.logger.Log(
			"method", "DeleteDeploy",
			"id", id,
			"err", err,
		)
	}()

	return r.next.DeleteDeploy(ctx, id)
}

func (r repositoryLogger) InitBuild(ctx context.Context, id string, jobName, jobId, imageName string) (err error) {
	defer func() {
		r.logger.Log(
			"method", "InitBuild",
			"id", id,
			"jobName", jobName,
			"jobId", jobId,
			"imageName", imageName,
			"err", err,
		)
	}()

	return r.next.InitBuild(ctx, id, jobName, jobId, imageName)
}

func (r repositoryLogger) InitWorkload(ctx context.Context, id string, jobName, jobId string, envs map[string]string) (err error) {
	defer func() {
		r.logger.Log(
			"method", "InitWorkload",
			"id", id,
			"jobName", jobName,
			"jobId", jobId,
			"envs", envs,
			"err", err,
		)
	}()

	return r.next.InitWorkload(ctx, id, jobName, jobId, envs)
}

func (r repositoryLogger) SetBuildStatus(ctx context.Context, id string, status service.Status) (err error) {
	defer func() {
		r.logger.Log(
			"method", "SetBuildStatus",
			"id", id,
			"status", status,
			"err", err,
		)
	}()

	return r.next.SetBuildStatus(ctx, id, status)
}

func (r repositoryLogger) RecordBuildStep(ctx context.Context, id string, buildStep service.BuildStep) (err error) {
	defer func() {
		r.logger.Log(
			"method", "RecordBuildStep",
			"id", id,
			"buildStep", buildStep,
			"err", err,
		)
	}()

	return r.next.RecordBuildStep(ctx, id, buildStep)
}
