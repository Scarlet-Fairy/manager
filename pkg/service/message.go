package service

type Message interface {
	Init() error
	ConsumeBuildEvents(id string) (<-chan *BuildStep, func() error, error)
}
