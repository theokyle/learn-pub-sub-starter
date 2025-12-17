package main

import (
	"fmt"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
)

func handlerLog(log routing.GameLog) pubsub.AckType {
	defer fmt.Print("> ")
	err := gamelogic.WriteLog(log)
	if err != nil {
		return pubsub.NackRequeue
	}
	return pubsub.Ack
}
