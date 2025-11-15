package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType string

const (
	DurableQueue   = "durable"
	TransientQueue = "transient"
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create channel: %v", err)
	}

	queue, err := channel.QueueDeclare(queueName, queueType == DurableQueue, queueType == TransientQueue, queueType == TransientQueue, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not create queue: %v", err)
	}

	err = channel.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		return nil, amqp.Queue{}, fmt.Errorf("could not bind queue: %v", err)
	}

	return channel, queue, nil
}
