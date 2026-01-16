package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
)

func handlerWar(gs *gamelogic.GameState) func(gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(recogWar gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(recogWar)
		switch warOutcome {
		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard
		case gamelogic.WarOutcomeOpponentWon:
			log.Printf("%s won a war against %s", winner, loser)
			return pubsub.Ack
		case gamelogic.WarOutcomeYouWon:
			log.Printf("%s won a war against %s", winner, loser)
			return pubsub.Ack
		case gamelogic.WarOutcomeDraw:
			log.Printf("A war between %s and %s resulted in a draw", winner, loser)
			return pubsub.Ack
		default:
			log.Print("error: handling war outcome")
			return pubsub.NackDiscard
		}
	}
}
