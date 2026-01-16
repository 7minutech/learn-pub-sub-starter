package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(recogWar gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(recogWar)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			gl := routing.GameLog{Username: gs.GetUsername(), Message: msg}
			err := publishGameLog(gl, publishCh)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			msg := fmt.Sprintf("%s won a war against %s", winner, loser)
			gl := routing.GameLog{Username: gs.GetUsername(), Message: msg}
			err := publishGameLog(gl, publishCh)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			msg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			gl := routing.GameLog{Username: gs.GetUsername(), Message: msg}
			err := publishGameLog(gl, publishCh)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		default:
			log.Print("error: handling war outcome")
			return pubsub.NackDiscard
		}
	}
}

func publishGameLog(gl routing.GameLog, publishCh *amqp.Channel) error {
	err := pubsub.PublishGob(publishCh, routing.ExchangePerilTopic, routing.GameLogSlug+gl.Username, gl)
	if err != nil {
		return err
	}
	return nil
}
