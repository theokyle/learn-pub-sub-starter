package main

import (
	"fmt"
	"log"

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

	channel, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	err = pubsub.PublishJSON(channel, string(routing.ExchangePerilDirect), string(routing.PauseKey), routing.PlayingState{IsPaused: true})
	if err != nil {
		log.Fatalf("could not create publish time: %v", err)
	}

	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan

	fmt.Println("Pause message sent")
}
