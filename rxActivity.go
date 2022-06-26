package main

import (
	"errors"

	"github.com/PiotrTopa/js8web/model"
)

func rxActivityNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent, databaseObjects chan<- model.DbObj) error {
	obj, err := model.CreateRxPacketObj(event)
	if err != nil {
		return errors.New("can not convert RxActivity event to db object")
	}
	databaseObjects <- obj
	return nil
}

func rxSpotNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent, databaseObjects chan<- model.DbObj) error {
	obj, err := model.CreateRxSpotObj(event)
	if err != nil {
		return errors.New("can not convert RxSpot event to db object")
	}
	databaseObjects <- obj
	return nil
}
