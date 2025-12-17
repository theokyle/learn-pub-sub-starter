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
	fmt.Println("Connecting to Rabbit...")
	connUrl := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connUrl)
	if err != nil {
		log.Fatalf("Error connecting to Rabbit: %v", err)
	}

	defer conn.Close()

	fmt.Println("Connection succesful!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Error getting username: %v", err)
	}

	ch, _, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("pause.%s", username),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("Error declaring queue: %v", err)
	}

	defer ch.Close()

	fmt.Println("Client running, press Ctrl+C to exit")

	game := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilDirect,
		fmt.Sprintf("%s.%s", routing.PauseKey, username),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(game))
	if err != nil {
		log.Fatalf("Error subscribing to connection: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username),
		fmt.Sprintf("%s.*", routing.ArmyMovesPrefix),
		pubsub.SimpleQueueTransient,
		handlerMove(game, ch))
	if err != nil {
		log.Fatalf("Error subscribing to connection: %v", err)
	}

	err = pubsub.SubscribeJSON(
		conn,
		string(routing.ExchangePerilTopic),
		string(routing.WarRecognitionsPrefix),
		fmt.Sprintf("%s.*", routing.WarRecognitionsPrefix),
		pubsub.SimpleQueueDurable,
		handlerWar(game),
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
		case "spawn":
			err = game.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
			}
		case "move":
			mv, err := game.CommandMove(words)
			if err != nil {
				fmt.Printf("Error moving piece: %v", err)
				continue
			}
			err = pubsub.PublishJSON(ch, string(routing.ExchangePerilTopic), fmt.Sprintf("army_moves.%s", username), mv)
			if err != nil {
				fmt.Printf("Error publishing move: %v", err)
				continue
			}
			fmt.Println("Move published succesfully")
		case "status":
			game.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			isRunning = false
		default:
			log.Println("Command not recognized")
		}
	}
}
