package main

import "github.com/streadway/amqp"

func main() {

}

func newRabbitMQClient(url string) (*amqp.Connection, error) {
	return amqp.Dial(url)
}
