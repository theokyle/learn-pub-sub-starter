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
	rabbitURL := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("Successfully connected to RabbitMQ")

	gamelogic.PrintServerHelp()
	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, "game_logs", "game_logs.*", pubsub.DurableQueue)
	if err != nil {
		log.Fatalf("error binding queue: %v", err)
	}

	quitloop := false
	for !quitloop {
		input := gamelogic.GetInput()
		switch input[0] {
		case "pause":
			err = pubsub.PublishJSON(channel, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Fatalf("could not create publish time: %v", err)
			}
			fmt.Println("Pause message sent.")
		case "resume":
			err = pubsub.PublishJSON(channel, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Fatalf("could not create publish time: %v", err)
			}
			fmt.Println("Resume message sent.")
		case "quit":
			fmt.Println("Exiting...")
			quitloop = true
		default:
			fmt.Println("Invalid command.")
		}
	}
}
