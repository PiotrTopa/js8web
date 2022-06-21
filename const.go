package main

import _ "embed"

var (
	JS8CALL_TCP_CONNECTION_STRING    = "localhost:2442"
	JS8CALL_TCP_CONNECTION_RETRY_SEC = 5
	DB_FILE_PATH                     = "./js8web.db"

	EVENT_TYPE_RX_ACTIVITY = "RX.ACTIVITY"
	EVENT_TYPE_RX_SPOT     = "RX.SPOT"
	EVENT_TYPE_RIG_PTT     = "RIG.PTT"
	EVENT_TYPE_TX_FRAME    = "TX.FRAME"
)

//resource files

//go:embed res/initDb.sql
var RESOURCE_INIT_DB_SQL string
