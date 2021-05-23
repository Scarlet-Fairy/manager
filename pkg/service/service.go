package service

import (
	"context"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

type Service interface {
	Deploy(ctx context.Context, gitRepo string, name string, envs map[string]string) (string, error)
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

	go func() {
		events, clear, err := s.message.ConsumeBuildEvents(id)
		if err != nil {
			return
		}

		for event := range events {

			if !event.Step.IsValid() {
				_ = s.repository.SetBuildStatus(ctx, id, StatusError)

				return
			}

			if event.Error != "" {
				if err := s.repository.RecordBuildStep(ctx, id, *event); err != nil {
					return
				}

				if err := s.repository.SetBuildStatus(ctx, id, StatusError); err != nil {
					return
				}

				break
			}
			if err := s.repository.RecordBuildStep(ctx, id, *event); err != nil {
				return
			}

			if event.Step == StepPush {
				if err := s.repository.SetBuildStatus(ctx, id, StatusCompleted); err != nil {
					return
				}

				jobName, err := s.scheduler.ScheduleWorkload(ctx, envs, id)
				if err != nil {
					_ = s.repository.SetBuildStatus(ctx, id, StatusError)

					return
				}

				if err := s.repository.InitWorkload(ctx, id, jobName, jobName, envs); err != nil {
					return
				}

				break
			}
		}
		_ = clear()
	}()

	return id, nil
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

func (s *basicService) GetDeploy(ctx context.Context, name string) (*Deploy, error) {
	deploy, err := s.repository.GetDeployByName(ctx, name)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}
