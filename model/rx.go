package model

import (
	"database/sql"
	"errors"
	"fmt"
)

var (
	MODE_JS8                  = "js8"
	EVENT_TYPE_RX_ACTIVITY    = "RX.ACTIVITY"
	EVENT_TYPE_RX_DIRECTED    = "RX.DIRECTED"
	EVENT_TYPE_RX_DIRECTED_ME = "RX.DIRECTED.ME:"
	EVENT_TYPE_RX_SPOT        = "RX.SPOT"
	EVENT_TYPE_RIG_PTT        = "RIG.PTT"
	EVENT_TYPE_TX_FRAME       = "TX.FRAME"

	SQL_RX_PACKETS_INSERT = "INSERT INTO `RX_MESSAGES` (`TYPE`, `CHANNEL`, `DIAL`, `FREQ`, `OFFSET`, `SNR`, `MODE`, `TIME_DRIFT`, `GRID`, `FROM`, `TO`, `TEXT`, `COMMAND`, `EXTRA`) values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
)

type RxPacket struct {
	Id        int64
	Type      string
	Dial      uint32
	Channel   uint16
	Freq      uint32
	Offset    uint16
	Snr       int16
	Mode      string
	Speed     string
	TimeDrift uint16
	Grid      string
	From      string
	To        string
	Text      string
	Command   string
	Extra     string
}

func speedName(speed int) string {
	switch speed {
	case 0:
		return "normal"
	case 1:
		return "fast"
	default:
		return "normal"
	}
}

func CreateRxPacket(event *Js8callEvent) (*RxPacket, error) {
	if event.Type != EVENT_TYPE_RX_ACTIVITY && event.Type != EVENT_TYPE_RX_DIRECTED && event.Type != EVENT_TYPE_RX_DIRECTED_ME {
		return nil, errors.New("Wrong event type, cannot parse params")
	}

	o := new(RxPacket)
	o.Type = event.Type
	o.Dial = event.Params.Dial
	o.Channel = uint16(event.Params.Offset / 50)
	o.Freq = event.Params.Freq
	o.Offset = event.Params.Offset
	o.Snr = event.Params.Snr
	o.Mode = MODE_JS8
	o.Speed = speedName(event.Params.Speed)
	o.TimeDrift = uint16(1000 * event.Params.TimeDrift)
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

func (obj *RxPacket) Insert(db *sql.DB) error {
	stmt, err := db.Prepare(SQL_RX_PACKETS_INSERT)
	if err != nil {
		return fmt.Errorf("Error inserting new RxPacket record, caused by %w", err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(
		&obj.Type,
		&obj.Channel,
		&obj.Dial,
		&obj.Freq,
		&obj.Offset,
		&obj.Snr,
		&obj.Mode,
		&obj.TimeDrift,
		&obj.Grid,
		&obj.From,
		&obj.To,
		&obj.Text,
		&obj.Command,
		&obj.Extra,
	)
	if err != nil {
		return fmt.Errorf("Error inserting new RxPacket record, becouse of %w", err)
	}
	obj.Id, _ = res.LastInsertId()

	return nil
}

func (obj *RxPacket) Save(db *sql.DB) error {
	return obj.Insert(db)
}
