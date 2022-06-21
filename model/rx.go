package model

import "database/sql"

var MODE_JS8 = "js8"

type RxPacket struct {
	Id        int
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

func (o *RxPacket) Parse(event *Js8callEvent) error {
	o.Type = event.Type
	o.Dial = event.Params.Dial
	o.Channel = uint16(event.Params.Offset / 50)
	o.Freq = event.Params.Freq
	o.Offset = event.Params.Offset
	o.Snr = event.Params.Snr
	o.Mode = MODE_JS8
	o.Speed = speedName(event.Params.Speed)
	return nil
}

func (obj *RxPacket) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO `RX`(`NAME`, `PASSWORD`, `ROLE`, `BIO`) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err = stmt.Exec(&obj.Name, &obj.Password, &obj.Role, &obj.Bio)
	if err != nil {
		return err
	}

	return nil
}
