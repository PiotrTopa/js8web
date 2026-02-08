package main

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"net/http"

	"github.com/PiotrTopa/js8web/model"
)

//go:embed webapp/*
var WEBAPP_FS embed.FS
var WEBAPP_SUBDIR = "webapp"

func startWebappServer(db *sql.DB, wsEventsSessionContainer *websocketSessionContainer, outgoingEvents chan<- model.Js8callEvent) {
	serverRoot, err := fs.Sub(WEBAPP_FS, WEBAPP_SUBDIR)
	if err != nil {
		logger.Sugar().Fatalw(
			"Cannot access WebApp subdirectory",
			"subdir", WEBAPP_SUBDIR,
			"error", err,
		)
	}

	webappFs := http.FileServer(http.FS(serverRoot))
	mux := http.NewServeMux()

	mux.HandleFunc("/api/rx-packets", methodHandler(methodRouter{
		get: apiRxPacketsGet,
	}, db))
	mux.HandleFunc("/api/station-info", methodHandler(methodRouter{
		get: apiStationInfoGet,
	}, db))
	mux.HandleFunc("/api/rig-status", methodHandler(methodRouter{
		get: apiRigStatusGet,
	}, db))
	mux.HandleFunc("/api/tx-message", authRequired(methodHandler(methodRouter{
		post: apiTxMessagePost(outgoingEvents),
	}, db)))
	mux.HandleFunc("/api/auth/login", methodHandler(methodRouter{
		post: apiAuthLoginPost,
	}, db))
	mux.HandleFunc("/api/auth/logout", methodHandler(methodRouter{
		post: apiAuthLogoutPost,
	}, db))
	mux.HandleFunc("/api/auth/check", methodHandler(methodRouter{
		get: apiAuthCheckGet,
	}, db))
	mux.HandleFunc("/ws/events", websocketHandler(wsEventsSessionContainer))
	mux.Handle("/", webappFs)

	err = http.ListenAndServe(fmt.Sprintf(":%d", WEBAPP_PORT), mux)
	if err != nil {
		logger.Sugar().Fatalw(
			"Cannot start WebApp HTTP server",
			"port", WEBAPP_PORT,
			"error", err,
		)
	}
}

type methodRouter struct {
	get  func(http.ResponseWriter, *http.Request, *sql.DB)
	post func(http.ResponseWriter, *http.Request, *sql.DB)
}

func methodNotSupported(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	logger.Sugar().Errorw(
		"Method is not supported by the API",
		"method", req.Method,
		"url", req.URL,
	)
	http.Error(w, "method not supported", http.StatusNotImplemented)
}

func methodHandler(r methodRouter, db *sql.DB) func(http.ResponseWriter, *http.Request) {
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
		f(w, req, db)
	}
}
