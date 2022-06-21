package main

type Js8callEvent struct {
	Type   string                 `json:"type"`
	Value  string                 `json:"value"`
	Params map[string]interface{} `json:"params"`
}

type Event interface {
}

type EventRxActivity struct {
	Type    string
	Dial    uint32
	Freq    uint32
	Grid    string
	From    string
	To      string
	Offset  uint16
	Snr     int16
	Speed   int
	Command string
	Text    string
	Extra   string
}
