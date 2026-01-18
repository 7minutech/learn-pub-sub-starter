package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"

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
		log.Fatalf("could not create channel on connection: %v", err)
	}

	userName, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not welcome user: %v", err)
	}

	queueName := routing.PauseKey + "." + userName
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatalf("could not declare and bind client queue: %v", err)
	}

	gameState := gamelogic.NewGameState(userName)

	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient, handlerPause(gameState))
	if err != nil {
		log.Fatalf("could not call SubscribeJSON for pause: %v", err)
	}

	moveKey := routing.ArmyMoveKey + ".*"
	moveQueueName := routing.ArmyMoveKey + "." + userName
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, moveQueueName, moveKey, pubsub.Transient, handlerMove(gameState, ch))
	if err != nil {
		log.Fatalf("could not call SubscribeJSON for move: %v", err)
	}

	warKey := routing.WarRecognitionsPrefix + ".*"
	warQueue := "war"
	err = pubsub.SubscribeJSON(conn, routing.ExchangePerilTopic, warQueue, warKey, pubsub.Durable, handlerWar(gameState, ch))
	if err != nil {
		log.Fatalf("could not call SubscribeJSON for move: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		if words[0] == "spawn" {
			if err := gameState.CommandSpawn(words); err != nil {
				fmt.Printf("could not spawn unit: %v", err)
				continue
			}
			log.Print("unit spawned")
		} else if words[0] == "move" {
			mv, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Printf("could not move unit: %v", err)
				continue
			}
			log.Print("unit moved")
			pubsub.PublishJSON(ch, routing.ExchangePerilTopic, moveQueueName, mv)
		} else if words[0] == "status" {
			gameState.CommandStatus()
		} else if words[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if words[0] == "spam" && len(words) == 2 {
			n, err := strconv.Atoi(words[1])
			if err != nil {
				fmt.Println("must provide a number with spam")
				continue
			}
			for range n {
				msg := gamelogic.GetMaliciousLog()
				err := pubsub.PublishGob(ch, routing.ExchangePerilTopic, routing.GameLogSlug+"."+userName, routing.GameLog{Message: msg})
				if err != nil {
					log.Println("failed to publish malicious log")
					continue
				}
			}
		} else if words[0] == "quit" {
			gamelogic.PrintQuit()
			break
		} else {
			fmt.Println("unknown command")
		}
	}

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	fmt.Println("RabbitMQ connection closed.")
}
