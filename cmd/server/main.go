package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const rabbitConnString = "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(rabbitConnString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Peril game server connected to RabbitMQ!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}

	gamelogic.PrintServerHelp()

	exchange := routing.ExchangePerilDirect
	key := routing.PauseKey

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		if words[0] == "pause" {
			log.Print("sending pause message")
			pubsub.PublishJSON(ch, exchange, key, routing.PlayingState{IsPaused: true})
		} else if words[0] == "resume" {
			log.Print("sending resume message")
			pubsub.PublishJSON(ch, exchange, key, routing.PlayingState{IsPaused: false})
		} else if words[0] == "quit" {
			log.Print("sending quit message")
			break
		} else {
			log.Printf("unknown command: %s", words[0])
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
