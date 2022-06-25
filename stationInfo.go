package main

import "github.com/PiotrTopa/js8web/model"

var stationInfoCache model.StationInfoWsEvent = model.StationInfoWsEvent{}

func stationInfoNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent) {
	newStationInfo := stationInfoCache
	err := newStationInfo.UpdateFromEvent(event)
	if err != nil {
		logger.Sugar().Errorw(
			"Can not undertand StationInfo event",
			"event", event,
			"error", err,
		)
		return
	}

	if newStationInfo != stationInfoCache {
		stationInfoCache = newStationInfo
		websocketEvents <- stationInfoCache
		updateStationInfoInDb(stationInfoCache)
	}
}

func updateStationInfoInDb(stationInfo model.StationInfoWsEvent) {
	go func() {
		stationInfoObj := model.CreateStationInfoObj(stationInfo)
		stationInfoObj.Save()
	}
}
