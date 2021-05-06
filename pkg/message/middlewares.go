package message

import (
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/log"
)

type Middleware func(message service.Message) service.Message

type messageLogger struct {
	logger log.Logger
	next   service.Message
}

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(message service.Message) service.Message {
		return &messageLogger{
			logger: logger,
			next:   message,
		}
	}
}

func (m messageLogger) Init() (err error) {
	defer func() {
		m.logger.Log(
			"method", "Init",
			"err", err,
		)
	}()

	return m.next.Init()
}

func (m messageLogger) ConsumeBuildEvents(id string) (events <-chan *service.BuildStep, clear func() error, err error) {
	defer func() {
		m.logger.Log(
			"method", "ConsumeBuildEvents",
			"id", id,
			"err", err,
		)
	}()

	return m.next.ConsumeBuildEvents(id)
}
