package main

import "github.com/PiotrTopa/js8web/model"

var rigStatusCache model.RigStatusWsEvent = model.RigStatusWsEvent{}

func rigStatusNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent, databaseObjects chan<- model.DbObj) error {
	newRigStatus, err := model.CreateRigStatusWsEvent(event)
	if err != nil {
		logger.Sugar().Errorw(
			"Can not undertand RIG.STATUS event",
			"event", event,
			"error", err,
		)
		return nil
	}

	if *newRigStatus != rigStatusCache {
		rigStatusCache = *newRigStatus
		websocketEvents <- rigStatusCache
	}
	return nil
}
