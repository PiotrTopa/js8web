package model

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	SQL_TX_FRAME_INSERT = "INSERT INTO `TX_FRAME` (`TIMESTAMP`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `MODE`, `SPEED`, `SELECTED`, `TONES`) values(?, ?, ?, ?, ?, ?, ?, ?, ?)"
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
