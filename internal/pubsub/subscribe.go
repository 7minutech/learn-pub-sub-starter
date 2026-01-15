package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type AckType int

const (
	ack         AckType = 0
	nackRequeue AckType = 1
	nackDiscard AckType = 2
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T) AckType,
) error {
	ch, _, err := DeclareAndBind(conn, exchange, queueName, key, queueType)
	if err != nil {
		return err
	}

	deliveryCh, err := ch.Consume(queueName, "", false, false, false, false, amqp.Table{})
	if err != nil {
		return err
	}

	recieveMessages := func() {
		for delivery := range deliveryCh {
			var obj T
			if err := json.Unmarshal(delivery.Body, &obj); err != nil {
				log.Print("error: unmarshling deliveries", err)
				return
			}
			ackType := handler(obj)
			switch ackType {
			case ack:
				delivery.Ack(false)
				log.Print("Message was ack")
			case nackRequeue:
				delivery.Nack(false, true)
				log.Print("Message was nack and requeued")
			case nackDiscard:
				delivery.Nack(false, false)
				log.Print("Message was nack and discarded")
			}
		}
	}

	go recieveMessages()

	return nil
}
