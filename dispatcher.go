package main

import (
	"errors"

	"github.com/PiotrTopa/js8web/model"
)

func separateStateChangesAndObjects(in <-chan model.Js8callEvent) (<-chan model.Js8callEvent, <-chan model.DbObj) {
	outEvents := make(chan model.Js8callEvent, 1)
	outObjects := make(chan model.DbObj, 1)

	go func() {
		defer close(outEvents)
		defer close(outObjects)

		for event := range in {

			dbObj, err := createDbObject(&event)
			if err == nil {
				outObjects <- dbObj
			} else {
				outEvents <- event
			}
		}
	}()

	return outEvents, outObjects
}

func createDbObject(event *model.Js8callEvent) (model.DbObj, error) {
	if event.Type == model.EVENT_TYPE_RX_ACTIVITY || event.Type == model.EVENT_TYPE_RX_DIRECTED || event.Type == model.EVENT_TYPE_RX_DIRECTED_ME {
		return model.CreateRxPacketObj(event)
	}
	return nil, errors.New("event is not DB object")
}

// JS8Call uses the same event name for two different kinds of events.
// STATION.STATUS is used for both current RIG status (freq, selected callsign)
// and as an answer for "STATION.GET_STATUS" command with value of STATUS message.
// Here we set the former to "RIG.STATUS" and leave the latter untached.
func fixSameNameForDifferentStationStatusEvents(events <-chan model.Js8callEvent) <-chan model.Js8callEvent {
	fixedEvents := make(chan model.Js8callEvent, 1)
	go func() {
		defer close(fixedEvents)
		for event := range events {
			if event.Type == model.EVENT_TYPE_STATION_STATUS {
				if event.Params.Freq > 0 {
					event.Type = model.EVENT_TYPE_RIG_STATUS
				}
			}
			fixedEvents <- event
		}
	}()
	return fixedEvents
}

func dispatchStateChangeEvents(events <-chan model.Js8callEvent) <-chan model.WebsocketEvent {
	websocketEvents := make(chan model.WebsocketEvent, 1)
	for event := range events {
		switch event.Type {
		case model.EVENT_TYPE_RIG_STATUS:
			rigStatusNotifier(&event, websocketEvents)
		}
	}
	return websocketEvents
}
