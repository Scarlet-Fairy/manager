package amqp

import "github.com/Scarlet-Fairy/manager/pkg/service"

type message struct {
	Topic string `json:"topic"`
	Error string `json:"error"`
}

func (m message) ParseTopic() service.Step {
	switch m.Topic {
	case "clone":
		return service.StepClone
	case "build":
		return service.StepBuild
	case "push":
		return service.StepPush
	default:
		return service.StepUnknown
	}
}
