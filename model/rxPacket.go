package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"
)

var (
	SQL_RX_PACKET_INSERT      = "INSERT INTO `RX_PACKET` (`TIMESTAMP`, `TYPE`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `SNR`, `MODE`, `SPEED`, `TIME_DRIFT`, `GRID`, `FROM`, `TO`, `TEXT`, `COMMAND`, `EXTRA`) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	SQL_RX_PACKET_LIST_AFTER  = "SELECT `ID`, `TIMESTAMP`, `TYPE`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `SNR`, `MODE`, `SPEED`, `TIME_DRIFT`, `GRID`, `FROM`, `TO`, `TEXT`, `COMMAND`, `EXTRA` FROM `RX_PACKET` WHERE `TIMESTAMP` >= ? ORDER BY `ID` ASC LIMIT 100"
	SQL_RX_PACKET_LIST_BEFORE = "SELECT * FROM (SELECT `ID`, `TIMESTAMP`, `TYPE`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `SNR`, `MODE`, `SPEED`, `TIME_DRIFT`, `GRID`, `FROM`, `TO`, `TEXT`, `COMMAND`, `EXTRA` FROM `RX_PACKET` WHERE `TIMESTAMP` <= ? ORDER BY `ID` DESC LIMIT 100) ORDER BY `ID` ASC"
)

type RxPacketObj struct {
	Id        int64
	Timestamp time.Time
	Type      string
	Dial      uint32
	Channel   uint16
	Freq      uint32
	Offset    uint16
	Snr       int16
	Mode      string
	Speed     string
	TimeDrift int16
	Grid      string
	From      string
	To        string
	Text      string
	Command   string
	Extra     string
}

func CreateRxPacketObj(event *Js8callEvent) (*RxPacketObj, error) {
	if event.Type != EVENT_TYPE_RX_ACTIVITY && event.Type != EVENT_TYPE_RX_DIRECTED && event.Type != EVENT_TYPE_RX_DIRECTED_ME {
		return nil, errors.New("wrong event type, cannot parse params")
	}

	o := new(RxPacketObj)
	o.Timestamp = fromJs8Timestamp(event.Params.UTC)
	o.Type = event.Type
	o.Dial = event.Params.Dial
	o.Freq = event.Params.Freq
	o.Offset = event.Params.Offset
	o.Channel = calcCahnnelFromOffset(o.Offset)
	o.Snr = event.Params.Snr
	o.Mode = MODE_JS8
	o.Speed = speedName(event.Params.Speed)
	o.TimeDrift = int16(1000 * event.Params.TimeDrift)
	o.Grid = event.Params.Grid
	o.From = event.Params.From
	o.To = event.Params.To
	o.Command = event.Params.Command
	o.Extra = event.Params.Extra

	if event.Params.Text != "" {
		o.Text = event.Params.Text
	} else {
		o.Text = event.Value
	}

	return o, nil
}

func (obj *RxPacketObj) Insert(db *sql.DB) error {
	stmt, err := db.Prepare(SQL_RX_PACKET_INSERT)
	if err != nil {
		return fmt.Errorf("error preparing SQL query fo inserting new RxPacket record, caused by %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		toSqlTime(obj.Timestamp),
		&obj.Type,
		&obj.Channel,
		&obj.Dial,
		&obj.Freq,
		&obj.Offset,
		&obj.Snr,
		&obj.Mode,
		&obj.Speed,
		&obj.TimeDrift,
		&obj.Grid,
		&obj.From,
		&obj.To,
		&obj.Text,
		&obj.Command,
		&obj.Extra,
	)
	if err != nil {
		return fmt.Errorf("error executing SQL query inserting new RxPacket record, becouse of %w", err)
	}

	obj.Id, _ = res.LastInsertId()
	return nil
}

func (obj *RxPacketObj) Save(db *sql.DB) error {
	return obj.Insert(db)
}

func (obj *RxPacketObj) Scan(rows *sql.Rows) error {
	var timestamp string
	err := rows.Scan(
		&obj.Id,
		&timestamp,
		&obj.Type,
		&obj.Channel,
		&obj.Dial,
		&obj.Freq,
		&obj.Offset,
		&obj.Snr,
		&obj.Mode,
		&obj.Speed,
		&obj.TimeDrift,
		&obj.Grid,
		&obj.From,
		&obj.To,
		&obj.Text,
		&obj.Command,
		&obj.Extra,
	)
	if err != nil {
		return err
	}

	obj.Timestamp, err = fromSqlTime(timestamp)
	return err
}

func fetchRxPackets(db *sql.DB, query string, args ...any) ([]RxPacketObj, error) {
	l := make([]RxPacketObj, 0)

	stmt, err := db.Prepare(query)
	if err != nil {
		return l, fmt.Errorf("error preparing SQL, caused by %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	err = nil
	for rows.Next() {
		obj := RxPacketObj{}
		err = obj.Scan(rows)
		if err != nil {
			return l, err
		}
		l = append(l, obj)
	}
	err = rows.Err()
	if err != nil {
		return l, err
	}

	return l, nil
}

func FetchRxPacketList(db *sql.DB, from time.Time, direction string) ([]RxPacketObj, error) {
	if direction == "before" {
		return fetchRxPackets(db, SQL_RX_PACKET_LIST_BEFORE, toSqlTime(from))
	}
	return fetchRxPackets(db, SQL_RX_PACKET_LIST_AFTER, toSqlTime(from))
}
