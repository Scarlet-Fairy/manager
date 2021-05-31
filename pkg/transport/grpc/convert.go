package grpc

import (
	"github.com/Scarlet-Fairy/manager/pb"
	"github.com/Scarlet-Fairy/manager/pkg/service"
)

func transportDeployToCoreDeploy(deploy *pb.Deploy) *service.Deploy {
	buildSteps := make([]*service.BuildStep, len(deploy.Build.Steps))
	for _, step := range deploy.Build.Steps {
		buildSteps = append(buildSteps, &service.BuildStep{
			Step:  service.Step(step.Step),
			Error: step.Error,
		})
	}

	return &service.Deploy{
		Id:      deploy.Id,
		Name:    deploy.Name,
		GitRepo: deploy.GitRepo,
		Build: &service.Build{
			JobId:     deploy.Build.JobId,
			JobName:   deploy.Build.JobName,
			ImageName: deploy.Build.ImageName,
			Status:    service.Status(deploy.Build.Status),
			Steps:     buildSteps,
		},
		Workload: &service.Workload{
			JobId:   deploy.Workload.JobId,
			JobName: deploy.Workload.JobName,
			Envs:    deploy.Workload.Envs,
			Url:     deploy.Workload.Url,
		},
	}
}

func coreDeployToTransportDeploy(deploy *service.Deploy) *pb.Deploy {
	var buildSteps []*pb.Build_BuildStep
	for _, step := range deploy.Build.Steps {
		buildSteps = append(buildSteps, &pb.Build_BuildStep{
			Step:  pb.Build_BuildStep_Step(step.Step),
			Error: step.Error,
		})
	}

	return &pb.Deploy{
		Id:      deploy.Id,
		Name:    deploy.Name,
		GitRepo: deploy.GitRepo,
		Build: &pb.Build{
			JobId:   deploy.Build.JobId,
			JobName: deploy.Build.JobName,
			Status:  pb.Build_Status(deploy.Build.Status),
			Steps:   buildSteps,
		},
		Workload: &pb.Workload{
			JobId:   deploy.Workload.JobId,
			JobName: deploy.Workload.JobName,
			Envs:    deploy.Workload.Envs,
			Url:     deploy.Workload.Url,
		},
	}
}
