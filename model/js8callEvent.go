package model

var (
	MODE_JS8                    = "js8"
	EVENT_TYPE_RX_ACTIVITY      = "RX.ACTIVITY"
	EVENT_TYPE_RX_DIRECTED      = "RX.DIRECTED"
	EVENT_TYPE_RX_DIRECTED_ME   = "RX.DIRECTED.ME"
	EVENT_TYPE_RX_SPOT          = "RX.SPOT"
	EVENT_TYPE_RIG_PTT          = "RIG.PTT"
	EVENT_TYPE_TX_FRAME         = "TX.FRAME"
	EVENT_TYPE_RIG_STATUS       = "RIG.STATUS"
	EVENT_TYPE_STATION_STATUS   = "STATION.STATUS"
	EVENT_TYPE_STATION_INFO     = "STATION.INFO"
	EVENT_TYPE_STATION_CALLSIGN = "STATION.CALLSIGN"
	EVENT_TYPE_STATION_GRID     = "STATION.GRID"

	// event types as seen in Websocket communication
	WS_EVENT_TYPE_RIG_PTT      = "RIG.PTT"
	WS_EVENT_TYPE_RIG_STATUS   = "RIG.STATUS"
	WS_EVENT_TYPE_STATION_INFO = "STATION.INFO"
	WS_OBJ_TYPE_RX_PACKET      = "RX_PACKET"
	WS_OBJ_TYPE_RX_SPOT        = "RX_SPOT"
	WS_OBJ_TYPE_TX_FRAME       = "TX_FRAME"
	WS_OBJ_TYPE_OTHER          = "OTHER"
)

type Js8callEvent struct {
	Type   string             `json:"type"`
	Value  string             `json:"value"`
	Params Js8callEventParams `json:"params"`

	DataType string
	Data     interface{}
}

type Js8callEventParams struct {
	Id        interface{} `json:"_ID"`
	Dial      uint32      `json:"DIAL"`
	Freq      uint32      `json:"FREQ"`
	Offset    uint16      `json:"OFFSET"`
	Snr       int16       `json:"SNR"`
	Speed     int         `json:"SPEED"`
	TimeDrift float32     `json:"TDRIFT"`
	Grid      string      `json:"GRID"`
	From      string      `json:"FROM"`
	Call      string      `json:"CALL"`
	To        string      `json:"TO"`
	Text      string      `json:"TEXT"`
	Command   string      `json:"CMD"`
	Extra     string      `json:"EXTRA"`
	PTT       bool        `json:"PTT"`
	Tones     []int       `json:"TONES"`
	UTC       int64       `json:"UTC"`
	Selected  string      `json:"SELECTED"`
	Band      string      `json:"BAND"`
	Mode      string      `json:"MODE"`
	Submode   string      `json:"SUBMODE"`
	RptSent   string      `json:"RPT.SENT"`
	RptRecv   string      `json:"RPT.RECV"`
}

func calcCahnnelFromOffset(offset uint16) uint16 {
	return uint16((offset - 25) / 50)
}

func speedName(speed int) string {
	switch speed {
	case 0:
		return "normal"
	case 1:
		return "fast"
	case 2:
		return "turbo"
	case 4:
		return "slow"
	case 8:
		return "ultra"
	default:
		return "unknown"
	}
}
