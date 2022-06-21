package main

import (
	"github.com/PiotrTopa/js8web/model"
)

// This file contains all generic modifications to raw events
// as they are coming from JS8call applied before any other
// dispatcher takes care
var num int = 0

func applyJs8callEventParser(in <-chan model.Js8callEvent) <-chan model.Js8callEvent {
	out := make(chan model.Js8callEvent, 1)

	go func() {
		defer close(out)
		for event := range in {
			testParser(&event)
			out <- event
		}
	}()

	return out
}

func testParser(event *model.Js8callEvent) {
	if event.Type == EVENT_TYPE_RX_ACTIVITY || event.Type == EVENT_TYPE_RX_DIRECTED {
		event.DataType = "RX"
		data := model.Rx{}
		data.Parse(event)
		event.Data = data
	}
}
