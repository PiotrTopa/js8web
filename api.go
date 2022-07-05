package main

import (
	"encoding/json"
	"net/http"
)

type methodRouter struct {
	get  func(http.ResponseWriter, *http.Request)
	post func(http.ResponseWriter, *http.Request)
}

func methodNotSupported(w http.ResponseWriter, req *http.Request) {
	logger.Sugar().Errorf(
		"method is not supported by the API",
		"method", req.Method,
		"url", req.URL,
	)
	http.Error(w, "method not supported", http.StatusNotImplemented)
}

func methodHandler(r methodRouter) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		f := methodNotSupported
		switch req.Method {
		case http.MethodGet:
			f = r.get
		case http.MethodPost:
			f = r.post
		}
		if f == nil {
			f = methodNotSupported
		}
		f(w, req)
	}
}

func apiStationInfoGet(w http.ResponseWriter, req *http.Request) {
	stationInfoJson, err := json.Marshal(stationInfoCache)
	if err != nil {
		logger.Sugar().Errorf(
			"Cannot marshal stationInfo",
			"stationInfo", stationInfoCache,
			"error", err,
		)
		http.Error(w, "cannot marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(stationInfoJson)
}

func apiRigStatusGet(w http.ResponseWriter, req *http.Request) {
	rigStatusJson, err := json.Marshal(rigStatusCache)
	if err != nil {
		logger.Sugar().Errorf(
			"Cannot marshal rigStatus",
			"rigStatus", rigStatusCache,
			"error", err,
		)
		http.Error(w, "cannot marshal json", http.StatusInternalServerError)
		return
	}
	w.Write(rigStatusJson)
}
