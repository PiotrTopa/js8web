package model

type WebsocketMessage struct {
	EventType string
	WsType    string
	Event     interface{}
}
