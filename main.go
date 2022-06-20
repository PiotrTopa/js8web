package main

import (
	"fmt"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()

	incomingEvents := make(chan Js8callEvent, 1)
	outgoingEvents := make(chan Js8callEvent, 1)

	initJs8callConnection(incomingEvents, outgoingEvents)

	defer close(incomingEvents)
	defer close(outgoingEvents)

	for event := range incomingEvents {
		fmt.Print("Processed incoming: ", event)
	}
}
