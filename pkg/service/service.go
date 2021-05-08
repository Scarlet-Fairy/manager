package service

import (
	"context"
	"github.com/pkg/errors"
)

type Service interface {
	Deploy(ctx context.Context, gitRepo string, name string, envs map[string]string) (*Deploy, error)
	Destroy(ctx context.Context, deployId string) error
}

type basicService struct {
	repository Repository
	message    Message
	scheduler  Scheduler
}

func NewService(repository Repository, message Message, scheduler Scheduler) Service {
	return &basicService{
		repository: repository,
		message:    message,
		scheduler:  scheduler,
	}
}

func (s *basicService) Deploy(ctx context.Context, gitRepoUrl string, name string, envs map[string]string) (*Deploy, error) {
	id, err := s.repository.CreateDeploy(ctx, &Deploy{
		Name:    name,
		GitRepo: gitRepoUrl,
		Build:   &Build{},
		Workload: &Workload{
			Envs: envs,
		},
	})
	if err != nil {
		return nil, errors.Wrap(err, "Deploy creation")
	}

	buildJobName, imageName, err := s.scheduler.ScheduleImageBuild(ctx, id, gitRepoUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Image Build Schedulation")
	}

	if err := s.repository.InitBuild(ctx, id, buildJobName, buildJobName, imageName); err != nil {
		return nil, errors.Wrap(err, "Storing build infos")
	}

	events, clear, err := s.message.ConsumeBuildEvents(id)
	if err != nil {
		return nil, errors.Wrap(err, "Queue consuming")
	}

	for event := range events {

		if !event.Step.IsValid() {
			if err := s.repository.SetBuildStatus(ctx, id, StatusError); err != nil {
				return nil, errors.Wrap(err, "Settings Build status on Error")
			}

			return nil, errors.Wrap(
				errors.New("Internal error while parsing"),
				"Build failed",
			)
		}

		if event.Error != "" {
			if err := s.repository.RecordBuildStep(ctx, id, *event); err != nil {
				return nil, errors.Wrap(err, "Recording Build Step with error")
			}

			if err := s.repository.SetBuildStatus(ctx, id, StatusError); err != nil {
				return nil, errors.Wrap(err, "Settings Build Status on Error")
			}

			break
		}
		if err := s.repository.RecordBuildStep(ctx, id, *event); err != nil {
			return nil, errors.Wrap(err, "Recoding Build Step")
		}

		if event.Step == StepPush {
			if err := s.repository.SetBuildStatus(ctx, id, StatusCompleted); err != nil {
				return nil, errors.Wrap(err, "Settings Build Status on Completed")
			}

			jobName, err := s.scheduler.ScheduleWorkload(ctx, envs, id)
			if err != nil {
				return nil, errors.Wrap(err, "Scheduling Workload")
			}

			if err := s.repository.InitWorkload(ctx, id, jobName, jobName, envs); err != nil {
				return nil, errors.Wrap(err, "Storing workload infos")
			}

			break
		}
	}
	if err := clear(); err != nil {
		return nil, err
	}

	deploy, err := s.repository.GetDeploy(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "Retrieving final Deploy infos")
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
