package service

type Status byte

const (
	StatusError     Status = 1
	StatusLoading   Status = 2
	StatusCompleted Status = 3
)

func (s Status) IsValid() bool {
	return s == StatusError || s == StatusLoading || s == StatusCompleted
}

type Step byte

const (
	StepInit  Step = 0
	StepClone Step = 1
	StepBuild Step = 2
	StepPush  Step = 3
)

func (s Step) IsValid() bool {
	return s == StepInit || s == StepClone || s == StepBuild || s == StepPush
}

type BuildStep struct {
	Step  Step
	Error string
}

type Build struct {
	JobId   string
	JobName string
	Status  Status
	Steps   []BuildStep
}

type Workload struct {
	JobId   string
	JobName string
	Envs    map[string]string
}

type Deploy struct {
	Id         string
	Name       string
	GithubRepo string
	Build      Build
	Workload   Workload
}
