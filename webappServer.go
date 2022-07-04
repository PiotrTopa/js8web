package main

import (
	"embed"
	"fmt"
	"io/fs"
	"net/http"
)

//go:embed webapp/*
var WEBAPP_FS embed.FS
var WEBAPP_SUBDIR = "webapp"

func startWebappServer() {
	serverRoot, err := fs.Sub(WEBAPP_FS, WEBAPP_SUBDIR)
	if err != nil {
		logger.Sugar().Fatalf(
			"Cannot access WebApp subdirectory",
			"subdir", WEBAPP_SUBDIR,
			"error", err,
		)
	}

	webappFs := http.FileServer(http.FS(serverRoot))
	mux := http.NewServeMux()
	mux.Handle("/", webappFs)

	err = http.ListenAndServe(fmt.Sprintf(":%d", WEBAPP_PORT), mux)
	if err != nil {
		logger.Sugar().Fatalf(
			"Cannot start WebApp HTTP server",
			"port", WEBAPP_PORT,
			"error", err,
		)
	}
}
