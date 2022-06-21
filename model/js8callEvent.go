package model

type Js8callEvent struct {
	Type   string             `json:"type"`
	Value  string             `json:"value"`
	Params Js8callEventParams `json:"params"`

	DataType string
	Data     interface{}
}

type Js8callEventParams struct {
	Id        int     `json:"_ID"`
	Dial      uint32  `json:"DIAL"`
	Freq      uint32  `json:"FREQ"`
	Offset    uint16  `json:"OFFSET"`
	Snr       int16   `json:"SNR"`
	Speed     int     `json:"SPEED"`
	TimeDrift float32 `json:"TDRIFT"`
	Grid      string  `json:"GRID"`
	From      string  `json:"FROM"`
	To        string  `json:"TO"`
	Text      string  `json:"TEXT"`
	Command   string  `json:"CMD"`
	Extra     string  `json:"EXTRA"`
	PTT       bool    `json:"PTT"`
	Tones     []uint8 `json:"TONES"`
	UTC       uint64  `json:"UTC"`
	Selected  string  `json:"SELECTED"`
	Band      string  `json:"BAND"`
	Mode      string  `json:"MODE"`
	Submode   string  `json:"SUBMODE"`
	RptSent   string  `json:"RPT.SENT"`
	RptRecv   string  `json:"RPT.RECV"`
}
