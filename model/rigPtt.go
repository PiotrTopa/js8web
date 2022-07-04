package model

import (
	"errors"
)

type RigPttWsEvent struct {
	Enabled bool
}

func (o *RigPttWsEvent) Type() string {
	return WS_EVENT_TYPE_RIG_PTT
}

func CreateRigPttWsEvent(event *Js8callEvent) (*RigPttWsEvent, error) {
	if event.Type != EVENT_TYPE_RIG_PTT {
		return nil, errors.New("wrong event type, cannot parse params")
	}
	o := new(RigPttWsEvent)
	o.Enabled = event.Params.PTT
	return o, nil
}
