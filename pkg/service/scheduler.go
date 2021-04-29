package service

import "context"

type Scheduler interface {
	ScheduleImageBuild(ctx context.Context, workloadId string, gitRepoUrl string) (string, string, error)
	ScheduleWorkload(ctx context.Context, envs map[string]string, workloadId string) error
	UnScheduleJob(ctx context.Context, jobId string) error
}
