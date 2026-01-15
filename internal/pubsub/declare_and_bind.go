package pubsub

import (
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SimpleQueueType int

const (
	Durable   SimpleQueueType = 0
	Transient SimpleQueueType = 1
)

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {

	var queue amqp.Queue
	ch, err := conn.Channel()
	if err != nil {
		return ch, queue, err
	}

	queue, err = ch.QueueDeclare(
		queueName,
		queueType == Durable,
		queueType != Durable,
		queueType != Durable,
		false,
		amqp.Table{"x-dead-letter-exchange": routing.Exchange_DLX},
	)
	if err != nil {
		return ch, queue, err
	}

	if err := ch.QueueBind(queueName, key, exchange, false, nil); err != nil {
		return ch, queue, err
	}

	return ch, queue, nil

}
