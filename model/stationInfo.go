package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	SQL_STATION_INFO_INSERT = "INSERT INTO `STATION_INFO` (`TIMESTAMP`, `LATEST`, `CALL`, `GRID`, `INFO`, `STATUS`) values(?, ?, ?, ?, ?, ?)"
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

func (o *StationInfoWsEvent) Type() string {
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
		return fmt.Errorf("error preparing SQL inserting new StationInfo record, caused by %w", err)
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
		return fmt.Errorf("error executing SQL inserting new StationInfo record, becouse of %w", err)
	}

	obj.Id, _ = res.LastInsertId()
	return nil
}

func (obj *StationInfoObj) Save(db *sql.DB) error {
	return nil
}
