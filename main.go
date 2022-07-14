package main

import (
	"database/sql"
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

	initStationInfoCache(db)

	initJs8callConnection(incomingEvents, outgoingEvents)
	outgoingWebsocketEvents, newObjects := dispatchStateChangeEvents(incomingEvents)

	websocketMessages := make(chan model.WebsocketMessage, 1)
	defer close(websocketMessages)
	go mainDispatcher(db, websocketMessages, outgoingWebsocketEvents, newObjects)

	wsSessionContainer := new(websocketSessionContainer)
	wsSessionContainer.init()
	go wsSessionContainer.process(websocketMessages)

	go startWebappServer(db, wsSessionContainer)

	for {
		time.Sleep(time.Second)
	}

}

func mainDispatcher(db *sql.DB, websocketMessages chan<- model.WebsocketMessage, outgoingWebsocketEvents <-chan model.WebsocketEvent, newObjects <-chan model.DbObj) {
	go func() {
		for object := range newObjects {
			err := object.Save(db)
			if err != nil {
				logger.Sugar().Errorw(
					"Error when saving object to DB",
					"object", object,
					"error", err,
				)
			}

			websocketMessages <- model.WebsocketMessage{
				EventType: "object",
				Event:     object,
			}
		}
	}()

	go func() {
		for event := range outgoingWebsocketEvents {
			websocketMessages <- model.WebsocketMessage{
				EventType: "event",
				Event:     event,
			}
		}
	}()
}
