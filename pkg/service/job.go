package service

import (
	"fmt"
)

type JobId string

const (
	JobTypeImageBuild = "imagebuild"
	JobTypeWorkload   = "workload"
	NameImageBuilder  = "cobold"
)

func (id JobId) NameImageBuild() string {
	return fmt.Sprintf("%s.%s", JobTypeImageBuild, id)
}

func (id JobId) NameWorkload() string {
	return fmt.Sprintf("%s.%s", JobTypeWorkload, id)
}

func (id JobId) ImageName(registry string) string {
	return fmt.Sprintf("%s/%s/%s", registry, NameImageBuilder, id)
}
