package service

import "context"

type Repository interface {
	CreateDeploy(ctx context.Context, deploy *Deploy) (string, error)
	GetDeploy(ctx context.Context, id string) (*Deploy, error)
	ListDeploy(ctx context.Context) ([]*Deploy, error)
	UpdateDeploy(ctx context.Context, deploy *Deploy) error
	DeleteDeploy(ctx context.Context, id string) error

	InitBuild(ctx context.Context, id string, jobName, jobId, imageName string) error
	InitWorkload(ctx context.Context, id string, jobName, jobId string, envs map[string]string) error
	SetBuildStatus(ctx context.Context, id string, status Status) error
	RecordBuildStep(ctx context.Context, id string, buildStep BuildStep) error
}
