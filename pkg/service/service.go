package service

import "context"

type Service interface {
	Deploy(ctx context.Context, gitRepo string, name string) (*Deploy, error)
	Destroy(ctx context.Context, deployId string) error
}

func NewService() Service {
	return &basicService{}
}

type basicService struct {
	repository Repository
	message    Message
	scheduler  Scheduler
}

func (s *basicService) Deploy(ctx context.Context, gitRepoUrl string, name string) (*Deploy, error) {
	panic("implement me")
}

func (s *basicService) Destroy(ctx context.Context, deployId string) error {
	panic("implement me")
}
