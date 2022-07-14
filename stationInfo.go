package main

import (
	"database/sql"

	"github.com/PiotrTopa/js8web/model"
)

var stationInfoCache model.StationInfoWsEvent = model.StationInfoWsEvent{}

func stationInfoNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent, databaseObjects chan<- model.DbObj) error {
	newStationInfo := stationInfoCache
	err := newStationInfo.UpdateFromEvent(event)
	if err != nil {
		logger.Sugar().Errorw(
			"Can not undertand StationInfo event",
			"event", event,
			"error", err,
		)
		return nil
	}

	if newStationInfo != stationInfoCache {
		stationInfoCache = newStationInfo
		websocketEvents <- &stationInfoCache

		stationInfoObj := model.CreateStationInfoObj(stationInfoCache)
		databaseObjects <- stationInfoObj
	}
	return nil
}

func initStationInfoCache(db *sql.DB) {
	stationInfo, err := model.FetchLatestStationInfo(db)
	if err != nil {
		logger.Sugar().Warnw(
			"Can not initialize StationInfo",
			"error", err,
		)
		return
	}
	stationInfoCache = stationInfo.StationInfoWsEvent
}
