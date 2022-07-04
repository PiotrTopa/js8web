package main

import (
	"errors"

	"github.com/PiotrTopa/js8web/model"
)

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

func dispatchStateChangeEvents(events <-chan model.Js8callEvent) (<-chan model.WebsocketEvent, <-chan model.DbObj) {
	events = fixSameNameForDifferentStationStatusEvents(events)

	websocketEvents := make(chan model.WebsocketEvent, 1)
	dbObjects := make(chan model.DbObj, 1)
	go func() {
		defer close(websocketEvents)
		defer close(dbObjects)

		//var f func(*model.Js8callEvent, chan<- model.WebsocketEvent, chan<- model.DbObj) error
		f := defaultNotifier
		for event := range events {
			switch event.Type {
			case model.EVENT_TYPE_RX_ACTIVITY, model.EVENT_TYPE_RX_DIRECTED, model.EVENT_TYPE_RX_DIRECTED_ME:
				f = rxActivityNotifier
			case model.EVENT_TYPE_RX_SPOT:
				f = rxSpotNotifier
			case model.EVENT_TYPE_TX_FRAME:
				f = txFrameNotifier
			case model.EVENT_TYPE_RIG_STATUS:
				f = rigStatusNotifier
			case model.EVENT_TYPE_STATION_CALLSIGN, model.EVENT_TYPE_STATION_GRID, model.EVENT_TYPE_STATION_INFO, model.EVENT_TYPE_STATION_STATUS:
				f = stationInfoNotifier
			default:
				f = defaultNotifier
			}

			err := f(&event, websocketEvents, dbObjects)
			if err != nil {
				logger.Sugar().Errorw(
					"Error dispatching event",
					"event", event,
					"error", err,
				)
			}
		}
	}()

	return websocketEvents, dbObjects
}

func defaultNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent, databaseObjects chan<- model.DbObj) error {
	return errors.New("unrecognized event type")
}
