package main

import "github.com/PiotrTopa/js8web/model"

var rigStatusCache *model.RigStatusWsEvent

func rigStatusNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent) {
	newRigStatus, err := model.CreateRigStatusWsEvent(event)
	if err != nil {
		logger.Sugar().Errorw(
			"Can not undertand RIG.STATUS event",
			"event", event,
			"error", err,
		)
		return
	}

	if *newRigStatus != *rigStatusCache {
		rigStatusCache = newRigStatus
		websocketEvents <- newRigStatus
	}
}
