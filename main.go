package main

import (
	"fmt"
	"time"

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

	db := initDbConnection()
	defer db.Close()

	initJs8callConnection(incomingEvents, outgoingEvents)

	outgoingWebsocketEvents, newObjects := dispatchStateChangeEvents(incomingEvents)

	go func() {
		for object := range newObjects {
			fmt.Print("DATABASE: ", object, "\n")
			err := object.Save(db)
			if err != nil {
				logger.Sugar().Errorw(
					"Error when saving object to DB",
					"object", object,
					"error", err,
				)
			}
		}
	}()

	go func() {
		for event := range outgoingWebsocketEvents {
			fmt.Print("WEBSOCKET: ", event, "\n")
		}
	}()

	for {
		time.Sleep(time.Second)
	}

}
