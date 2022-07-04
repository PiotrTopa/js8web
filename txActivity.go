package main

import (
	"errors"

	"github.com/PiotrTopa/js8web/model"
)

func txFrameNotifier(event *model.Js8callEvent, websocketEvents chan<- model.WebsocketEvent, databaseObjects chan<- model.DbObj) error {
	obj, err := model.CreateTxFrameObj(event)
	if err != nil {
		return errors.New("can not convert TxFrame event to db object")
	}
	obj.ApplyRigStatus(&rigStatusCache)
	databaseObjects <- obj
	return nil
}
