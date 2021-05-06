package amqp

import (
	"encoding/json"
	"fmt"
	middlewares "github.com/Scarlet-Fairy/manager/pkg/message"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	BuildImageExchanger = "build_image"
)

type rabbitMessage struct {
	channel *amqp.Channel
}

func New(ch *amqp.Channel, logger log.Logger) service.Message {
	var instance service.Message
	instance = &rabbitMessage{
		channel: ch,
	}
	instance = middlewares.LoggingMiddleware(logger)(instance)

	return instance
}

func (m *rabbitMessage) Init() error {
	if err := m.channel.ExchangeDeclare(
		BuildImageExchanger,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return errors.Wrap(err, fmt.Sprintf("Failed to declare %s exchange", BuildImageExchanger))
	}

	return nil
}

func (m *rabbitMessage) declareQueues(id string) error {
	_, err := m.channel.QueueDeclare(
		id,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	if err := m.channel.QueueBind(
		id,
		id,
		BuildImageExchanger,
		false,
		nil,
	); err != nil {
		return err
	}

	return nil
}

func (m *rabbitMessage) ConsumeBuildEvents(id string) (<-chan *service.BuildStep, func() error, error) {
	if err := m.declareQueues(id); err != nil {
		return nil, nil, err
	}

	events, err := m.channel.Consume(
		id,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	buildStepEvents := make(chan *service.BuildStep)
	go func() {
		for e := range events {
			var msg message
			if err := json.Unmarshal(e.Body, &msg); err != nil {
				buildStepEvents <- &service.BuildStep{
					Step:  service.StepUnknown,
					Error: "Failed to parse message from cobold",
				}

				return
			}

			step := msg.ParseTopic()
			if !step.IsValid() {
				buildStepEvents <- &service.BuildStep{
					Step:  service.StepUnknown,
					Error: "Failed to parse Topic name",
				}

				return
			}

			if msg.Error != "" {
				buildStepEvents <- &service.BuildStep{
					Step:  step,
					Error: msg.Error,
				}

				return
			}

			buildStepEvents <- &service.BuildStep{
				Step: step,
			}
		}
	}()

	return buildStepEvents, func() error {
		_, err := m.channel.QueueDelete(
			id,
			false,
			false,
			false,
		)
		if err != nil {
			return err
		}

		return nil
	}, nil
}
