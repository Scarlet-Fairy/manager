package grpc

import (
	"context"
	"github.com/Scarlet-Fairy/manager/pb"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcScheduler struct {
	client pb.SchedulerClient
}

func New(client pb.SchedulerClient) service.Scheduler {
	return &grpcScheduler{
		client: client,
	}
}

func (g grpcScheduler) ScheduleImageBuild(ctx context.Context, workloadId string, gitRepoUrl string) (string, string, error) {
	res, err := g.client.ScheduleImageBuild(ctx, &pb.ScheduleImageBuildRequest{
		WorkloadId: workloadId,
		GitRepoUrl: gitRepoUrl,
	})
	if err != nil {
		return "", "", g.handleGrpcError(err)
	}

	return res.JobName, res.ImageName, nil
}

func (g grpcScheduler) ScheduleWorkload(ctx context.Context, envs map[string]string, workloadId string) error {
	_, err := g.client.ScheduleWorkload(ctx, &pb.ScheduleWorkloadRequest{
		Envs:       envs,
		WorkloadId: workloadId,
	})
	if err != nil {
		return g.handleGrpcError(err)
	}

	return nil
}

func (g grpcScheduler) UnScheduleJob(ctx context.Context, jobId string) error {
	_, err := g.client.UnScheduleJob(ctx, &pb.UnScheduleJobRequest{
		JobId: jobId,
	})
	if err != nil {
		return g.handleGrpcError(err)
	}

	return nil
}

func (g grpcScheduler) handleGrpcError(err error) error {
	if e, ok := status.FromError(err); ok {
		switch e.Code() {
		case codes.InvalidArgument:
			return errors.Wrap(errors.New(e.Message()), "Scheduler")
		default:
			return errors.Errorf("%s: %s", e.Code(), e.Message())
		}
	} else {
		return errors.New("Unable to parse gRPC Scheduler error")
	}
}
