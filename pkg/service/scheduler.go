package service

import "context"

type Scheduler interface {
	ScheduleImageBuild(ctx context.Context, workloadId string, gitRepoUrl string) (jobName string, imageName string, err error)
	ScheduleWorkload(ctx context.Context, envs map[string]string, workloadId string) (jobName string, url string, err error)
	UnScheduleJob(ctx context.Context, jobId string) error
}
