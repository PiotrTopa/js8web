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
		return model.CreateRxPacket(event)
	}
	return nil, errors.New("event is not DB object")
}
