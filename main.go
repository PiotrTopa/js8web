package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PiotrTopa/js8web/model"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
)

func main() {
	logger, _ = zap.NewDevelopment()
	defer logger.Sync()

	initConfig()

	logger.Sugar().Infow("js8web starting",
		"js8callAddr", JS8CALL_TCP_CONNECTION_STRING,
		"port", WEBAPP_PORT,
		"db", DB_FILE_PATH,
	)

	incomingEvents := make(chan model.Js8callEvent, 1)
	outgoingEvents := make(chan model.Js8callEvent, 1)

	db := initDbConnection()
	defer db.Close()

	initStationInfoCache(db)

	initJs8callConnection(incomingEvents, outgoingEvents)
	outgoingWebsocketEvents, newObjects := dispatchStateChangeEvents(incomingEvents)

	websocketMessages := make(chan model.WebsocketMessage, 1)
	go mainDispatcher(db, websocketMessages, outgoingWebsocketEvents, newObjects)

	wsSessionContainer := new(websocketSessionContainer)
	wsSessionContainer.init()
	go wsSessionContainer.process(websocketMessages)

	go startWebappServer(db, wsSessionContainer)

	logger.Sugar().Infof("js8web ready — http://localhost:%d", WEBAPP_PORT)

	// Wait for termination signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit

	fmt.Println() // newline after ^C
	logger.Sugar().Infow("Shutting down", "signal", sig.String())
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
				WsType:    object.WsType(),
				Event:     object,
			}
		}
	}()

	for event := range outgoingWebsocketEvents {
		websocketMessages <- model.WebsocketMessage{
			EventType: "event",
			WsType:    event.WsType(),
			Event:     event,
		}
	}
}
