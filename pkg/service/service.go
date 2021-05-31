package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type Service interface {
	Deploy(ctx context.Context, gitRepo string, name string, envs map[string]string) (string, error)
	HandleEvent(ctx context.Context, event *BuildStep, buildId string, envs map[string]string) (bool, error)
	Destroy(ctx context.Context, deployId string) error
	GetDeploy(ctx context.Context, name string) (*Deploy, error)
}

type basicService struct {
	repository Repository
	message    Message
	scheduler  Scheduler
}

func NewService(repository Repository, message Message, scheduler Scheduler, logger log.Logger) Service {
	var service Service
	{
		service = &basicService{
			repository: repository,
			message:    message,
			scheduler:  scheduler,
		}
		service = LoggingMiddleware(logger)(service)
	}

	return service
}

func (s *basicService) Deploy(ctx context.Context, gitRepoUrl string, name string, envs map[string]string) (string, error) {
	id, err := s.repository.CreateDeploy(ctx, &Deploy{
		Name:    name,
		GitRepo: gitRepoUrl,
		Build:   &Build{},
		Workload: &Workload{
			Envs: envs,
		},
	})
	if err != nil {
		return "", errors.Wrap(err, "Deploy creation")
	}

	buildJobName, imageName, err := s.scheduler.ScheduleImageBuild(ctx, id, gitRepoUrl)
	if err != nil {
		return "", errors.Wrap(err, "Image Build Schedulation")
	}

	if err := s.repository.InitBuild(ctx, id, buildJobName, buildJobName, imageName); err != nil {
		return "", errors.Wrap(err, "Storing build infos")
	}

	if err := s.repository.SetBuildStatus(ctx, id, StatusLoading); err != nil {
		return "", errors.Wrap(err, "Failed to set build status")
	}

	events, clear, err := s.message.ConsumeBuildEvents(id)
	if err != nil {
		if err := s.repository.SetBuildStatus(ctx, id, StatusError); err != nil {
			return "", errors.Wrap(err, "Failed to set build status")
		}

		return "", errors.Wrap(err, "Failed to consume events")
	}

	go func() {
		for event := range events {
			done, err := s.HandleEvent(context.Background(), event, id, envs)
			if err != nil {
				return
			}

			if done {
				break
			}
		}
		_ = clear()
	}()

	return id, nil
}

func (s *basicService) HandleEvent(ctx context.Context, event *BuildStep, buildId string, envs map[string]string) (bool, error) {
	if !event.Step.IsValid() {
		if err := s.repository.SetBuildStatus(ctx, buildId, StatusError); err != nil {
			return false, errors.Wrap(err, "Settings Build status on Error")
		}

		return false, errors.Wrap(
			errors.New("Internal error while parsing"),
			"Build failed",
		)
	}

	if event.Error != "" {
		if err := s.repository.RecordBuildStep(ctx, buildId, *event); err != nil {
			return false, errors.Wrap(err, "Recording Build Step with error")
		}

		if err := s.repository.SetBuildStatus(ctx, buildId, StatusError); err != nil {
			return false, errors.Wrap(err, "Settings Build Status on Error")
		}

		return true, nil
	}
	if err := s.repository.RecordBuildStep(ctx, buildId, *event); err != nil {
		return false, errors.Wrap(err, "Recoding Build Step")
	}

	if event.Step == StepPush {
		if err := s.repository.SetBuildStatus(ctx, buildId, StatusCompleted); err != nil {
			return false, errors.Wrap(err, "Settings Build Status on Completed")
		}

		jobName, url, err := s.scheduler.ScheduleWorkload(ctx, envs, buildId)
		if err != nil {
			if err := s.repository.SetBuildStatus(ctx, buildId, StatusError); err != nil {
				return false, errors.Wrap(err, "Settings Build Status on Error")
			}

			return false, errors.Wrap(err, "Scheduling Workload")
		}

		if err := s.repository.InitWorkload(ctx, buildId, jobName, jobName, envs, url); err != nil {
			return false, errors.Wrap(err, "Storing workload infos")
		}

		return true, nil
	}

	return false, nil
}

func (s *basicService) Destroy(ctx context.Context, deployId string) error {
	deploy, err := s.repository.GetDeploy(ctx, deployId)
	if err != nil {
		return errors.Wrap(err, "Retrieving Deploy")
	}

	if err := s.repository.DeleteDeploy(ctx, deployId); err != nil {
		return errors.Wrap(err, "Deleting Deploy")
	}

	if err := s.scheduler.UnScheduleJob(ctx, deploy.Workload.JobId); err != nil {
		return errors.Wrap(err, "UnScheduling Workload")
	}

	return nil
}

func (s *basicService) GetDeploy(ctx context.Context, id string) (*Deploy, error) {
	deploy, err := s.repository.GetDeploy(ctx, id)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}
