package main

import (
	_ "embed"
	"flag"
	"fmt"
	"os"
	"strconv"
)

var (
	JS8CALL_TCP_CONNECTION_STRING    string
	JS8CALL_TCP_CONNECTION_RETRY_SEC int
	DB_FILE_PATH                     string
	WEBAPP_PORT                      int
)

//resource files

//go:embed res/initDb.sql
var RESOURCE_INIT_DB_SQL string

func envOrDefault(key string, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func envOrDefaultInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if n, err := strconv.Atoi(val); err == nil {
			return n
		}
	}
	return fallback
}

func initConfig() {
	flag.StringVar(&JS8CALL_TCP_CONNECTION_STRING, "js8call-addr", envOrDefault("JS8WEB_JS8CALL_ADDR", "localhost:2442"), "JS8Call TCP API address (host:port)")
	flag.IntVar(&JS8CALL_TCP_CONNECTION_RETRY_SEC, "reconnect-interval", envOrDefaultInt("JS8WEB_RECONNECT_SEC", 5), "Seconds between JS8Call reconnection attempts")
	flag.StringVar(&DB_FILE_PATH, "db", envOrDefault("JS8WEB_DB_PATH", "./js8web.db"), "Path to SQLite database file")
	flag.IntVar(&WEBAPP_PORT, "port", envOrDefaultInt("JS8WEB_PORT", 8080), "HTTP server port")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "js8web — Web monitor and control interface for JS8Call\n\n")
		fmt.Fprintf(os.Stderr, "Usage: js8web [options]\n\nOptions:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nEnvironment variables:\n")
		fmt.Fprintf(os.Stderr, "  JS8WEB_JS8CALL_ADDR   JS8Call TCP API address (default: localhost:2442)\n")
		fmt.Fprintf(os.Stderr, "  JS8WEB_RECONNECT_SEC  Reconnect interval in seconds (default: 5)\n")
		fmt.Fprintf(os.Stderr, "  JS8WEB_DB_PATH        SQLite database file path (default: ./js8web.db)\n")
		fmt.Fprintf(os.Stderr, "  JS8WEB_PORT           HTTP server port (default: 8080)\n")
	}

	flag.Parse()
}
