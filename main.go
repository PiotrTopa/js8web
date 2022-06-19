package main

import "go.uber.org/zap"

var (
	logger *zap.Logger
)

func main() {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()

	incomingEvents := make(chan Js8callEvent, 1)
	outgoingEvents := make(chan Js8callEvent, 1)

	initJs8callEventStreams(incomingEvents, outgoingEvents)

	defer close(incomingEvents)
	defer close(outgoingEvents)

	for {

	}
}
