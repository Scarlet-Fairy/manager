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
	StepInit    Step = 0
	StepClone   Step = 1
	StepBuild   Step = 2
	StepPush    Step = 3
	StepUnknown Step = 4
)

func (s Step) IsValid() bool {
	return s == StepInit || s == StepClone || s == StepBuild || s == StepPush
}

func (s Step) ToString() string {
	switch s {
	case StepInit:
		return "init"
	case StepClone:
		return "clone"
	case StepBuild:
		return "build"
	case StepPush:
		return "push"
	default:
		return "unknown"
	}
}

type BuildStep struct {
	Step  Step
	Error string
}

type Build struct {
	JobId     string
	JobName   string
	ImageName string
	Status    Status
	Steps     []*BuildStep
}

type Workload struct {
	JobId   string
	JobName string
	Envs    map[string]string
	Url     string
}

type Deploy struct {
	Id       string
	Name     string
	GitRepo  string
	Build    *Build
	Workload *Workload
}
