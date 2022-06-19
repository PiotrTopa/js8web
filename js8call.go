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

type readJs8callEventResult struct {
	event Js8callEvent
	err   error
}

func readEventFromJs8call(reader *bufio.Reader) <-chan readJs8callEventResult {
	var (
		result readJs8callEventResult
		event  Js8callEvent
	)

	fmt.Print("Started reading")

	resultChannel := make(chan readJs8callEventResult, 1)
	defer close(resultChannel)

	jsonData, err := reader.ReadBytes('\n')
	if err != nil {
		logger.Sugar().Warnw("Error reading from Js8Call",
			"error", err,
		)
		result.err = err
		resultChannel <- result
		return resultChannel
	}

	errJson := json.Unmarshal(jsonData, &event)

	if errJson != nil {
		logger.Sugar().Warnw("Cannot unmarshal JSON",
			"json", jsonData,
			"error", errJson,
		)
	} else {
		result.event = event
	}

	resultChannel <- result
	fmt.Print("TEST2")
	return resultChannel
}

func keepConnectedToJs8callEventStreams(incomingEvents chan<- Js8callEvent, outgoingEvents <-chan Js8callEvent) {
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

		reader := bufio.NewReader(conn)
		//writer := bufio.NewWriter(conn)

		func() {
			for {
				select {
				case result := <-readEventFromJs8call(reader):
					fmt.Print("Incoming:", result)
					if result.err != nil {
						return
					}
					incomingEvents <- result.event
					break
				case event := <-outgoingEvents:
					fmt.Print("Outgoing:", event)
					break
				}
			}
		}()
		logger.Sugar().Warn("JS8call disconnected")
	}
}

func initJs8callEventStreams(incomingEvents chan<- Js8callEvent, outgoingEvents <-chan Js8callEvent) {
	go keepConnectedToJs8callEventStreams(incomingEvents, outgoingEvents)
}
