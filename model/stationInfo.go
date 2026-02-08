package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	SQL_STATION_INFO_INSERT        = "INSERT INTO `STATION_INFO` (`TIMESTAMP`, `LATEST`, `CALL`, `GRID`, `INFO`, `STATUS`) values(?, ?, ?, ?, ?, ?)"
	SQL_STATION_INFO_UPDATE_LATEST = "UPDATE `STATION_INFO` SET `LATEST` = 0 WHERE `LATEST` = 1 AND `ID` != ?"
	SQL_STATION_INFO_FETCH_LATEST  = "SELECT * FROM `STATION_INFO` WHERE `LATEST` = 1"
)

type StationInfoWsEvent struct {
	Call   string
	Grid   string
	Info   string
	Status string
}

type StationInfoObj struct {
	Id        int64
	Timestamp string
	Latest    bool
	StationInfoWsEvent
}

func (o *StationInfoWsEvent) WsType() string {
	return WS_EVENT_TYPE_STATION_INFO
}

func (o *StationInfoWsEvent) UpdateFromEvent(event *Js8callEvent) error {
	switch event.Type {
	case EVENT_TYPE_STATION_CALLSIGN:
		o.Call = event.Value
	case EVENT_TYPE_STATION_GRID:
		o.Grid = event.Value
	case EVENT_TYPE_STATION_STATUS:
		o.Status = event.Value
	case EVENT_TYPE_STATION_INFO:
		o.Info = event.Value
	default:
		return errors.New("event type does not match stationInfo type")
	}
	return nil
}

func CreateStationInfoObj(stationInfo StationInfoWsEvent) *StationInfoObj {
	return &StationInfoObj{
		StationInfoWsEvent: stationInfo,
		Timestamp:          toSqlTime(time.Now()),
		Latest:             true,
	}
}

func (obj *StationInfoObj) Insert(db *sql.DB) error {
	stmt, err := db.Prepare(SQL_STATION_INFO_INSERT)
	if err != nil {
		return fmt.Errorf("error preparing SQL, caused by %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		&obj.Timestamp,
		&obj.Latest,
		&obj.Call,
		&obj.Grid,
		&obj.Info,
		&obj.Status,
	)
	if err != nil {
		return fmt.Errorf("error executing SQL Insert, caused by %w", err)
	}

	obj.Id, _ = res.LastInsertId()
	return nil
}

func (obj *StationInfoObj) updateLatest(db *sql.DB) error {
	if obj.Id == 0 {
		return errors.New("cannot update latest flag without ID set")
	}

	stmt, err := db.Prepare(SQL_STATION_INFO_UPDATE_LATEST)
	if err != nil {
		return fmt.Errorf("error preparing SQL, caused by %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(&obj.Id)
	if err != nil {
		return fmt.Errorf("error executing SQL updateLatest, caused by %w", err)
	}

	return nil
}

func (obj *StationInfoObj) Save(db *sql.DB) error {
	err := obj.Insert(db)
	if err != nil {
		return err
	}

	err = obj.updateLatest(db)
	return err
}

func (obj *StationInfoObj) Scan(rows *sql.Row) error {
	err := rows.Scan(
		&obj.Id,
		&obj.Timestamp,
		&obj.Latest,
		&obj.Call,
		&obj.Grid,
		&obj.Info,
		&obj.Status,
	)
	return err
}

func FetchLatestStationInfo(db *sql.DB) (StationInfoObj, error) {
	row := db.QueryRow(SQL_STATION_INFO_FETCH_LATEST)
	obj := StationInfoObj{}
	err := obj.Scan(row)
	if err != nil {
		return obj, fmt.Errorf("error fetching latest StationInfo, caused by %w", err)
	}
	return obj, nil
}
