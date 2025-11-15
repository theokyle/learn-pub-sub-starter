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
	fmt.Println("Starting Peril client...")

	rabbitURL := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(rabbitURL)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("error getting username: %v", err)
	}

	queueName := routing.PauseKey + "." + username
	_, _, err = pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.TransientQueue)
	if err != nil {
		log.Fatalf("error binding queue: %v", err)
	}

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
}
