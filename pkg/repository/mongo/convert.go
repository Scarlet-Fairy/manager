package mongo

import (
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func businessToData(deploy *service.Deploy) *Deploy {
	var id primitive.ObjectID
	if deploy.Id != "" {
		parsedId, err := primitive.ObjectIDFromHex(deploy.Id)
		if err != nil {
			panic(err)
		}

		id = parsedId
	} else {
		id = primitive.NewObjectID()
	}

	var steps []*BuildStep
	for _, step := range deploy.Build.Steps {
		steps = append(steps, &BuildStep{
			Step:  int(step.Step),
			Error: step.Error,
		})
	}
	envs := mapEnvToArrEnv(deploy.Workload.Envs)

	return &Deploy{
		Id:      id,
		Name:    deploy.Name,
		GitRepo: deploy.GitRepo,
		Build: &Build{
			JobId:   deploy.Build.JobId,
			JobName: deploy.Build.JobName,
			Status:  int(deploy.Build.Status),
			Steps:   steps,
		},
		Workload: &Workload{
			JobId:   deploy.Workload.JobId,
			JobName: deploy.Workload.JobName,
			Envs:    envs,
			Url:     deploy.Workload.Url,
		},
	}
}

func dataToBusiness(deploy *Deploy) *service.Deploy {
	id := ""
	if deploy.Id != primitive.NilObjectID {
		id = deploy.Id.Hex()
	}

	var steps []*service.BuildStep
	for _, step := range deploy.Build.Steps {
		steps = append(steps, &service.BuildStep{
			Step:  service.Step(step.Step),
			Error: step.Error,
		})
	}

	envs := arrEnvToMapEnv(deploy.Workload.Envs)

	return &service.Deploy{
		Id:      id,
		Name:    deploy.Name,
		GitRepo: deploy.GitRepo,
		Build: &service.Build{
			JobId:     deploy.Build.JobId,
			JobName:   deploy.Build.JobName,
			ImageName: deploy.Build.ImageName,
			Status:    service.Status(deploy.Build.Status),
			Steps:     steps,
		},
		Workload: &service.Workload{
			JobId:   deploy.Workload.JobId,
			JobName: deploy.Workload.JobName,
			Envs:    envs,
			Url:     deploy.Workload.Url,
		},
	}
}

func mapEnvToArrEnv(envsToConvert map[string]string) []*Env {
	var envs []*Env
	if envsToConvert != nil {
		for key, value := range envsToConvert {
			envs = append(envs, &Env{
				Key:   key,
				Value: value,
			})
		}
	}

	return envs
}

func arrEnvToMapEnv(envsToConvert []*Env) map[string]string {
	envs := make(map[string]string)

	for _, env := range envsToConvert {
		envs[env.Key] = env.Value
	}

	return envs
}
