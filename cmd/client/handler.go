package main

import (
	"fmt"
	"time"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func handlerPause(gs *gamelogic.GameState) func(routing.PlayingState) pubsub.AckType {
	return func(ps routing.PlayingState) pubsub.AckType {
		defer fmt.Print("> ")
		gs.HandlePause(ps)
		return pubsub.Ack
	}
}

func handlerMove(gs *gamelogic.GameState, publishCh *amqp.Channel) func(gamelogic.ArmyMove) pubsub.AckType {
	return func(move gamelogic.ArmyMove) pubsub.AckType {
		defer fmt.Print("> ")
		outcome := gs.HandleMove(move)
		if outcome == gamelogic.MoveOutcomeMakeWar {
			err := pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.WarRecognitionsPrefix, gs.GetUsername()),
				gamelogic.RecognitionOfWar{
					Attacker: move.Player,
					Defender: gs.GetPlayerSnap(),
				},
			)
			if err != nil {
				fmt.Printf("error: %v\n", err)
				return pubsub.NackRequeue
			}
			return pubsub.Ack
		}
		if outcome == gamelogic.MoveOutComeSafe {
			return pubsub.Ack
		}
		return pubsub.NackDiscard
	}
}

func handlerWar(gs *gamelogic.GameState, publishCh *amqp.Channel) func(dw gamelogic.RecognitionOfWar) pubsub.AckType {
	return func(dw gamelogic.RecognitionOfWar) pubsub.AckType {
		defer fmt.Print("> ")
		warOutcome, winner, loser := gs.HandleWar(dw)
		switch warOutcome {
		case gamelogic.WarOutcomeNoUnits:
			return pubsub.NackDiscard

		case gamelogic.WarOutcomeOpponentWon:
			logmsg := fmt.Sprintf("%s won a war against %s", winner, loser)
			err := publishGameLog(publishCh, gs, logmsg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack

		case gamelogic.WarOutcomeYouWon:
			logmsg := fmt.Sprintf("%s won a war against %s", winner, loser)
			err := publishGameLog(publishCh, gs, logmsg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack

		case gamelogic.WarOutcomeDraw:
			logmsg := fmt.Sprintf("A war between %s and %s resulted in a draw", winner, loser)
			err := publishGameLog(publishCh, gs, logmsg)
			if err != nil {
				return pubsub.NackRequeue
			}
			return pubsub.Ack

		case gamelogic.WarOutcomeNotInvolved:
			return pubsub.NackRequeue

		default:
			fmt.Println("error: unknown war outcome")
			return pubsub.NackDiscard
		}
	}
}

func publishGameLog(publishCh *amqp.Channel, gs *gamelogic.GameState, message string) error {
	err := pubsub.PublishGob(
		publishCh,
		string(routing.ExchangePerilTopic),
		fmt.Sprintf("%s.%s", routing.GameLogSlug, gs.GetUsername()),
		routing.GameLog{
			CurrentTime: time.Now(),
			Message:     message,
			Username:    gs.GetUsername(),
		},
	)
	if err != nil {
		return err
	}
	return nil
}
