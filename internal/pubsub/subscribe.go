package pubsub

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func SubscribeJSON[T any](
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType SimpleQueueType, // an enum to represent "durable" or "transient"
	handler func(T),
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
			handler(obj)
			delivery.Ack(false)
		}
	}

	go recieveMessages()

	return nil
}
