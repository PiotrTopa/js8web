package model

type WebsocketEvent interface {
	WsType() string
}
