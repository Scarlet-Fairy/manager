package service

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type Service interface {
	Deploy(ctx context.Context, gitRepo string, name string, envs map[string]string) (*Deploy, error)
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

func (s *basicService) Deploy(ctx context.Context, gitRepoUrl string, name string, envs map[string]string) (*Deploy, error) {
	id, err := s.repository.CreateDeploy(ctx, &Deploy{
		Name:     name,
		GitRepo:  gitRepoUrl,
		Build:    &Build{},
		Workload: &Workload{},
	})
	if err != nil {
		return nil, err
	}

	buildJobName, imageName, err := s.scheduler.ScheduleImageBuild(ctx, id, gitRepoUrl)
	if err != nil {
		return nil, err
	}

	if err := s.repository.InitBuild(ctx, id, buildJobName, buildJobName, imageName); err != nil {
		return nil, err
	}

	events, clear, err := s.message.ConsumeBuildEvents(id)
	if err != nil {
		return nil, err
	}

	for event := range events {
		if !event.Step.IsValid() {
			if err := s.repository.SetBuildStatus(ctx, id, StatusError); err != nil {
				return nil, err
			}

			return nil, errors.Wrap(
				errors.New("Internal error while parsing"),
				"Build failed",
			)
		}

		if event.Error != "" {
			if err := s.repository.RecordBuildStep(ctx, id, *event); err != nil {
				return nil, err
			}

			if err := s.repository.SetBuildStatus(ctx, id, StatusError); err != nil {
				return nil, err
			}

			return nil, errors.Wrap(
				errors.New(fmt.Sprintf("[%s] %s", event.Step.ToString(), event.Error)),
				"Build failed",
			)
		}

		if err := s.repository.RecordBuildStep(ctx, id, *event); err != nil {
			return nil, err
		}

		if event.Step == StepPush {
			if err := s.repository.SetBuildStatus(ctx, id, StatusCompleted); err != nil {
				return nil, err
			}

			break
		}
	}
	if err := clear(); err != nil {
		return nil, err
	}

	jobName, err := s.scheduler.ScheduleWorkload(ctx, envs, id)
	if err != nil {
		return nil, err
	}

	if err := s.repository.InitWorkload(ctx, id, jobName, jobName, envs); err != nil {
		return nil, err
	}

	deploy, err := s.repository.GetDeploy(ctx, id)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}

func (s *basicService) Destroy(ctx context.Context, deployId string) error {
	if err := s.repository.DeleteDeploy(ctx, deployId); err != nil {
		return err
	}

	if err := s.scheduler.UnScheduleJob(ctx, JobId(deployId).NameWorkload()); err != nil {
		return err
	}

	return nil
}
