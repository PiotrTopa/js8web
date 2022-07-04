package main

import (
	_ "embed"
)

var (
	JS8CALL_TCP_CONNECTION_STRING    = "localhost:2442"
	JS8CALL_TCP_CONNECTION_RETRY_SEC = 5
	DB_FILE_PATH                     = "./js8web.db"
	WEBAPP_PORT                      = 8080
)

//resource files

//go:embed res/initDb.sql
var RESOURCE_INIT_DB_SQL string
