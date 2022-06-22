package main

import (
	"fmt"

	"github.com/PiotrTopa/js8web/model"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()

	incomingEvents := make(chan model.Js8callEvent, 1)
	outgoingEvents := make(chan model.Js8callEvent, 1)
	defer close(incomingEvents)
	defer close(outgoingEvents)

	initDbConnection()
	initJs8callConnection(incomingEvents, outgoingEvents)

	stateChangeEvents, newObjects := separateStateChangesAndObjects(incomingEvents)

	go func() {
		for event := range newObjects {
			fmt.Print("OBJECT: ", event, "\n")
		}
	}()

	go func() {
		for event := range stateChangeEvents {
			fmt.Print("STATE CHANGE: ", event, "\n")
		}
	}()

	for {
	}

}
