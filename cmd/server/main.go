package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril server...")
	connUrl := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connUrl)
	if err != nil {
		log.Fatalf("Error connecting to Rabbit: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connection succesful!")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error establishing channel: %v", err)
	}

	gamelogic.PrintServerHelp()

	err = pubsub.SubscribeGob(
		conn,
		string(routing.ExchangePerilTopic),
		string(routing.GameLogSlug),
		fmt.Sprintf("%s.*", routing.GameLogSlug),
		pubsub.SimpleQueueDurable,
		handlerLog,
	)
	if err != nil {
		log.Fatalf("Error subscribing to connection: %v", err)
	}

	isRunning := true

	for isRunning {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			log.Println("Sending pause message...")
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("Error publishing data: %v", err)
			}

		case "resume":
			log.Println("Sending resume message...")
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("Error publishing data: %v", err)
			}

		case "quit":
			log.Println("Shutting down...")
			isRunning = false

		default:
			log.Println("Command not recognized")
		}
	}
}
