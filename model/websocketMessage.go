package model

type WebsocketMessage struct {
	EventType string `json:"type"`
	Event     interface{}
}
