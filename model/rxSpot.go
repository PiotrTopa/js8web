package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

var (
	SQL_RX_SPOT_INSERT = "INSERT INTO `RX_SPOTS` (`TIMESTAMP`, `ID`, `CALL`, `GRID`, `SNR`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`) values(?, ?, ?, ?, ?, ?, ?, ?)"
)

type RxSpotObj struct {
	Id        int64
	Timestamp time.Time
	Call      string
	Grid      string
	Snr       int16
	Channel   uint16
	Dial      uint32
	Freq      uint32
	Offset    uint16
}

func CreateRxSpotObj(event *Js8callEvent) (*RxSpotObj, error) {
	if event.Type != EVENT_TYPE_RX_SPOT {
		return nil, errors.New("wrong event type, cannot parse params")
	}

	o := new(RxSpotObj)
	o.Timestamp = fromJs8Timestamp(event.Params.UTC)
	o.Call = event.Params.Call
	o.Grid = event.Params.Grid
	o.Snr = event.Params.Snr
	o.Dial = event.Params.Dial
	o.Channel = uint16(event.Params.Offset / 50)
	o.Freq = event.Params.Freq
	o.Offset = event.Params.Offset

	return o, nil
}

func (obj *RxSpotObj) Insert(db *sql.DB) error {
	stmt, err := db.Prepare(SQL_RX_SPOT_INSERT)
	if err != nil {
		return fmt.Errorf("error inserting new RxSpot record, caused by %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		toSqlTime(obj.Timestamp),
		&obj.Call,
		&obj.Grid,
		&obj.Snr,
		&obj.Channel,
		&obj.Dial,
		&obj.Freq,
		&obj.Offset,
	)
	if err != nil {
		return fmt.Errorf("error inserting new RxSpot record, becouse of %w", err)
	}

	obj.Id, _ = res.LastInsertId()
	return nil
}

func (obj *RxSpotObj) Save(db *sql.DB) error {
	return obj.Insert(db)
}
