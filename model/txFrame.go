package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	SQL_TX_FRAME_INSERT      = "INSERT INTO `TX_FRAME` (`TIMESTAMP`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `MODE`, `SPEED`, `SELECTED`, `TONES`) values(?, ?, ?, ?, ?, ?, ?, ?, ?)"
	SQL_TX_FRAME_LIST_AFTER  = "SELECT `ID`, `TIMESTAMP`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `MODE`, `SPEED`, `SELECTED` FROM `TX_FRAME` WHERE `TIMESTAMP` > ?1 ORDER BY `ID` ASC LIMIT 100"
	SQL_TX_FRAME_LIST_BEFORE = "SELECT * FROM (SELECT `ID`, `TIMESTAMP`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `MODE`, `SPEED`, `SELECTED` FROM `TX_FRAME` WHERE `TIMESTAMP` <= ?1 ORDER BY `ID` DESC LIMIT 100) ORDER BY `ID` ASC"
)

type TxFrameObj struct {
	Id        int64
	Timestamp time.Time
	Channel   uint16
	Dial      uint32
	Freq      uint32
	Offset    uint16
	Mode      string
	Speed     string
	Selected  string
	Tones     []int
}

func (o *TxFrameObj) WsType() string {
	return WS_OBJ_TYPE_TX_FRAME
}

func CreateTxFrameObj(event *Js8callEvent) (*TxFrameObj, error) {
	if event.Type != EVENT_TYPE_TX_FRAME {
		return nil, errors.New("wrong event type, cannot parse params")
	}

	o := new(TxFrameObj)
	o.Timestamp = time.Now().UTC()
	o.Tones = event.Params.Tones
	return o, nil
}

func (obj *TxFrameObj) ApplyRigStatus(rig *RigStatusWsEvent) {
	obj.Channel = rig.Channel
	obj.Offset = rig.Offset
	obj.Dial = rig.Dial
	obj.Freq = rig.Freq
	obj.Speed = rig.Speed
	obj.Selected = rig.Selected
}

func (obj *TxFrameObj) Insert(db *sql.DB) error {
	stmt, err := db.Prepare(SQL_TX_FRAME_INSERT)
	if err != nil {
		return fmt.Errorf("error preparing SQL query fo inserting new TxFrame record, caused by %w", err)
	}
	defer stmt.Close()

	marshalledTones, err := json.Marshal(obj.Tones)
	if err != nil {
		return fmt.Errorf("unable to marshall tones %w", err)
	}

	res, err := stmt.Exec(
		toSqlTime(obj.Timestamp),
		&obj.Channel,
		&obj.Dial,
		&obj.Freq,
		&obj.Offset,
		&obj.Mode,
		&obj.Speed,
		&obj.Selected,
		string(marshalledTones),
	)
	if err != nil {
		return fmt.Errorf("error executing SQL query inserting new TxFrame record, caused by %w", err)
	}

	obj.Id, _ = res.LastInsertId()
	return nil
}

func (obj *TxFrameObj) Save(db *sql.DB) error {
	return obj.Insert(db)
}

func (obj *TxFrameObj) Scan(rows *sql.Rows) error {
	var timestamp string
	err := rows.Scan(
		&obj.Id,
		&timestamp,
		&obj.Channel,
		&obj.Dial,
		&obj.Freq,
		&obj.Offset,
		&obj.Mode,
		&obj.Speed,
		&obj.Selected,
	)
	if err != nil {
		return err
	}
	obj.Timestamp, err = fromSqlTime(timestamp)
	return err
}

func fetchTxFrames(db *sql.DB, query string, args ...any) ([]TxFrameObj, error) {
	l := make([]TxFrameObj, 0)

	stmt, err := db.Prepare(query)
	if err != nil {
		return l, fmt.Errorf("error preparing SQL, caused by %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(args...)
	if err != nil {
		return l, fmt.Errorf("error executing SQL query, caused by %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		obj := TxFrameObj{}
		err = obj.Scan(rows)
		if err != nil {
			return l, err
		}
		l = append(l, obj)
	}
	if err = rows.Err(); err != nil {
		return l, err
	}
	return l, nil
}

func FetchTxFrameList(db *sql.DB, startTime time.Time, direction string) ([]TxFrameObj, error) {
	query := SQL_TX_FRAME_LIST_BEFORE
	if direction == "after" {
		query = SQL_TX_FRAME_LIST_AFTER
	}
	return fetchTxFrames(db, query, toSqlTime(startTime))
}
