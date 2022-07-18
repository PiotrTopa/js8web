package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/PiotrTopa/js8web/model"
)

func parseTimestamp(t string) (time.Time, error) {
	return time.Parse(time.RFC3339, t)
}

func apiStationInfoGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	stationInfoJson, err := json.Marshal(stationInfoCache)
	if err != nil {
		logger.Sugar().Errorw(
			"Cannot marshal stationInfo",
			"stationInfo", stationInfoCache,
			"error", err,
		)
		http.Error(w, "cannot marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(stationInfoJson)
}

func apiRigStatusGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	rigStatusJson, err := json.Marshal(rigStatusCache)
	if err != nil {
		logger.Sugar().Errorw(
			"Cannot marshal rigStatus",
			"rigStatus", rigStatusCache,
			"error", err,
		)
		http.Error(w, "cannot marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(rigStatusJson)
}

func apiRxPacketsGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	q := req.URL.Query()
	if !q.Has("startTime") {
		http.Error(w, "'startTime' parameter is required", http.StatusBadRequest)
		return
	}

	startTime, err := parseTimestamp(q.Get("startTime"))
	if err != nil {
		logger.Sugar().Warnw(
			"Cannot parse timestamp",
			"time", startTime,
			"error", err,
		)
		http.Error(w, "cannot parse timestamp in 'startTime' parameter", http.StatusBadRequest)
		return
	}

	if !q.Has("direction") {
		http.Error(w, "'direction' parameter is required", http.StatusBadRequest)
		return
	}

	direction := q.Get("direction")
	if direction != "after" && direction != "before" {
		http.Error(w, "'direction' parameter has to be 'before' or 'after'", http.StatusBadRequest)
		return
	}

	filter := &model.RxPacketFilter{}
	if q.Has("filter") {
		err := json.Unmarshal([]byte(q.Get("filter")), filter)
		if err != nil {
			http.Error(w, "unable to parse filter", http.StatusInternalServerError)
			return
		}
	}

	list, err := model.FetchRxPacketList(db, filter, startTime, direction)
	if err != nil {
		logger.Sugar().Errorw(
			"Cannot fetch RxPacket records from DB",
			"error", err,
		)
		http.Error(w, "cannot fetch RxPacket records", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(list)
	if err != nil {
		logger.Sugar().Errorw(
			"Cannot marshal RxPacket records json",
			"error", err,
		)
		http.Error(w, "cannot marshal RxPacket records", http.StatusInternalServerError)
		return
	}

	w.Write(response)
}
