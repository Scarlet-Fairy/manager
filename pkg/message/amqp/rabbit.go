package amqp

import (
	"encoding/json"
	"fmt"
	"github.com/Scarlet-Fairy/manager/pkg/service"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
)

const (
	BuildImageExchanger = "build_image"
)

type rabbitMessage struct {
	channel *amqp.Channel
}

func New(ch *amqp.Channel) service.Message {
	return &rabbitMessage{
		channel: ch,
	}
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

func (m *rabbitMessage) ConsumeBuildEvents(id string) (<-chan *service.BuildStep, error) {
	if err := m.declareQueues(id); err != nil {
		return nil, err
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
		return nil, err
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

	return buildStepEvents, nil
}
