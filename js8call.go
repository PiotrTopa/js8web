package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"time"
)

var (
	JS8CALL_TCP_CONNECTION_STRING    = "localhost:2442"
	JS8CALL_TCP_CONNECTION_RETRY_SEC = 5
	JS8CALL_TCP_CONNECTION_TIMEOUT   = 10
)

func readEventsFromJs8call(events chan<- Js8callEvent, disconnected chan<- int, reader *bufio.Reader) {
	for {
		var event Js8callEvent
		jsonData, err := reader.ReadBytes('\n')
		if err != nil {
			logger.Sugar().Warnw("Error reading from Js8Call",
				"error", err,
			)
			disconnected <- 1
			return
		}

		errJson := json.Unmarshal(jsonData, &event)

		if errJson != nil {
			logger.Sugar().Warnw("Cannot unmarshal JSON",
				"json", jsonData,
				"error", errJson,
			)
		} else {
			events <- event
		}
	}
}

func writeEventsToJs8call(events <-chan Js8callEvent, disconnected chan<- int, writer *bufio.Writer) {
	fmt.Printf("Writer started")
	for event := range events {
		fmt.Printf("Writer object:", event)
	}
	fmt.Printf("Writer stopped")
}

func attachEventsStreamsToJs8callConnection(incomingEvents chan<- Js8callEvent, outgoingEvents <-chan Js8callEvent, conn net.Conn) {
	disconnected := make(chan int)
	incomingJs8callEvents := make(chan Js8callEvent, 1)
	outgoingJs8callEvents := make(chan Js8callEvent, 1)

	defer close(incomingJs8callEvents)
	defer close(outgoingJs8callEvents)
	defer close(disconnected)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	go readEventsFromJs8call(incomingJs8callEvents, disconnected, reader)
	go writeEventsToJs8call(outgoingJs8callEvents, disconnected, writer)

	for {
		select {
		case <-disconnected:
			return
		case event := <-incomingJs8callEvents:
			incomingEvents <- event
		case event := <-outgoingEvents:
			outgoingJs8callEvents <- event
		}
	}
}

func keepConnectedToJs8call(incomingEvents chan<- Js8callEvent, outgoingEvents <-chan Js8callEvent) {
	for {
		conn, err := net.Dial("tcp", JS8CALL_TCP_CONNECTION_STRING)
		if err != nil {
			logger.Sugar().Warnw("JS8Call connection error",
				"address", JS8CALL_TCP_CONNECTION_STRING,
				"error", err,
			)
			time.Sleep(time.Second * time.Duration(JS8CALL_TCP_CONNECTION_RETRY_SEC))
			continue
		}
		logger.Sugar().Info("Connected to JS8call")
		attachEventsStreamsToJs8callConnection(incomingEvents, outgoingEvents, conn)
		logger.Sugar().Warn("JS8call disconnected")
	}
}

func initJs8callConnection(incomingEvents chan<- Js8callEvent, outgoingEvents <-chan Js8callEvent) {
	go keepConnectedToJs8call(incomingEvents, outgoingEvents)
}
