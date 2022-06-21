package model

import "database/sql"

type RxPackets struct {
	Id        int
	Type      string
	Dial      uint32
	Channel   uint16
	Freq      uint32
	Offset    uint16
	Snr       int16
	Mode      string
	TimeDrift uint16
	Grid      string
	From      string
	To        string
	Text      string
	Command   string
	Extra     string
	MessageId int
}

func (o *Rx) Parse(event *Js8callEvent) error {
	o.Type = event.Type
	o.Dial = event.Params.Dial
	o.Channel = uint16(event.Params.Offset / 50)
	o.Freq = event.Params.Freq
	return nil
}

func (obj *Rx) Insert(db *sql.DB) error {
	stmt, err := db.Prepare("INSERT INTO `RX`(`NAME`, `PASSWORD`, `ROLE`, `BIO`) values(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	//_, err = stmt.Exec(&obj.Name, &obj.Password, &obj.Role, &obj.Bio)
	if err != nil {
		return err
	}
	return nil
}
