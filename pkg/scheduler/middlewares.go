package scheduler

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/log"
)

type Middleware func(scheduler service.Scheduler) service.Scheduler

type schedulerLogger struct {
	next   service.Scheduler
	logger log.Logger
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(scheduler service.Scheduler) service.Scheduler {
		return &schedulerLogger{
			next:   scheduler,
			logger: logger,
		}
	}
}

func (s schedulerLogger) ScheduleImageBuild(ctx context.Context, workloadId string, gitRepoUrl string) (jobName string, imageName string, err error) {
	defer func() {
		s.logger.Log(
			"method", "ScheduleImageBuild",
			"workloadId", workloadId,
			"gitRepoUrl", gitRepoUrl,
			"jobName", jobName,
			"imageName", imageName,
			"err", err,
		)
	}()

	return s.next.ScheduleImageBuild(ctx, workloadId, gitRepoUrl)
}

func (s schedulerLogger) ScheduleWorkload(ctx context.Context, envs map[string]string, workloadId string) (jobName string, err error) {
	defer func() {
		s.logger.Log(
			"method", "ScheduleWorkload",
			"envs", envs,
			"workloadId", workloadId,
			"jobName", jobName,
			"err", err,
		)
	}()

	return s.next.ScheduleWorkload(ctx, envs, workloadId)
}

func (s schedulerLogger) UnScheduleJob(ctx context.Context, jobId string) (err error) {
	defer func() {
		s.logger.Log(
			"method", "UnScheduleJob",
			"jobId", jobId,
			"err", err,
		)
	}()

	return s.next.UnScheduleJob(ctx, jobId)
}
