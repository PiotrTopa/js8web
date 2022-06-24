package model

import (
	"errors"
)

type RigStatusWsEvent struct {
	Dial     uint32
	Freq     uint32
	Offset   uint16
	Speed    int
	Selected string
}

func (o *RigStatusWsEvent) Type() string {
	return EVENT_TYPE_RIG_STATUS
}

func CreateRigStatusWsEvent(event *Js8callEvent) (*RigStatusWsEvent, error) {
	if event.Type != EVENT_TYPE_RIG_STATUS {
		return nil, errors.New("wrong event type, cannot parse params")
	}
	o := new(RigStatusWsEvent)
	o.Dial = event.Params.Dial
	o.Freq = event.Params.Freq
	o.Offset = event.Params.Offset
	o.Selected = event.Params.Selected
	o.Speed = event.Params.Speed
	return o, nil
}
