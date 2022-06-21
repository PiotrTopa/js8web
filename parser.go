package main

// This file contains all generic modifications to raw events
// as they are coming from JS8call applied before any other
// dispatcher takes care
var num int = 0

func applyJs8callEventParser(in <-chan Js8callEvent) <-chan Js8callEvent {
	out := make(chan Js8callEvent, 1)

	go func() {
		defer close(out)
		for event := range in {
			testParser(&event)
			out <- event
		}
	}()

	return out
}

func testParser(event *Js8callEvent) {
	num = num + 1
	event.Test = num
}
