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
	defer close(incomingEvents)
	defer close(outgoingEvents)

	initDbConnection()
	initJs8callConnection(incomingEvents, outgoingEvents)

	parsedIncomingEvents := applyJs8callEventParser(incomingEvents)
	dbIncomingEvents, webappIncomingEvents := copyChannel(parsedIncomingEvents)

	go func() {
		for event := range dbIncomingEvents {
			fmt.Print("Processed DB incoming: ", event)
		}
	}()

	go func() {
		for event := range webappIncomingEvents {
			fmt.Print("Processed webapp incoming: ", event)
		}
	}()

	for {
	}

}
