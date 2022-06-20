package main

import (
	"database/sql"
	"errors"
	"os"

	"github.com/PiotrTopa/js8web/model"
	_ "github.com/mattn/go-sqlite3"
)

// Initializes DB for the first time
func initDb(db *sql.DB) {
	logger.Sugar().Infow(
		"Initializing empty database",
		"file", DB_FILE_PATH,
	)

	_, err := db.Exec(RESOURCE_INIT_DB_SQL)
	if err != nil {
		logger.Sugar().Fatalw(
			"Could not initialize database",
			"file", DB_FILE_PATH,
			"error", err,
		)
	}

	err = model.DefaultAdminUser.Insert(db)
	if err != nil {
		logger.Sugar().Fatal(
			"Could not setup default admin user",
			"error", err,
		)
	}
	logger.Sugar().Info("Empty database initialized")
}

func initDbConnection() *sql.DB {
	var recreate bool = false

	if _, err := os.Stat(DB_FILE_PATH); errors.Is(err, os.ErrNotExist) {
		logger.Sugar().Warnw(
			"Database file not found",
			"file", DB_FILE_PATH,
		)
		recreate = true
	}

	db, err := sql.Open("sqlite3", DB_FILE_PATH)
	if err != nil {
		logger.Sugar().Fatalw("Could not open or create database file",
			"file", DB_FILE_PATH,
			"error", err,
		)
	}
	defer db.Close()

	if recreate == true {
		initDb(db)
	}

	return db
}
