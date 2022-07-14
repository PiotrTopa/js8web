package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/PiotrTopa/js8web/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type websocketSession struct {
	user   string
	events chan<- model.WebsocketMessage
}

type websocketSessionContainer struct {
	sessions map[string]websocketSession
}

func (o *websocketSessionContainer) init() {
	o.sessions = make(map[string]websocketSession)
}

func (o *websocketSessionContainer) register(id string, session websocketSession) error {
	if len(id) == 0 {
		return errors.New("no session id specified")
	}

	if _, exists := o.sessions[id]; exists {
		return errors.New("session already registered")
	}

	o.sessions[id] = session
	return nil
}

func (o *websocketSessionContainer) deregister(id string) error {
	if len(id) == 0 {
		return errors.New("no session id specified")
	}

	if session, exists := o.sessions[id]; exists {
		close(session.events)
		delete(o.sessions, id)
		return nil
	}

	return errors.New("session with given id has not been rgistered")
}

func (o *websocketSessionContainer) process(messages <-chan model.WebsocketMessage) {
	for message := range messages {
		for _, session := range o.sessions {
			session.events <- message
		}
	}
}

func websocketHandler(sessionContainer *websocketSessionContainer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Sugar().Errorw(
				"Can not upgrade WebSocket connection",
				"request", r,
				"error", err,
			)
			http.Error(w, "cannot upgrade connection", http.StatusUpgradeRequired)
			return
		}
		defer ws.Close()

		messages := make(chan model.WebsocketMessage)
		sessionId := uuid.New().String()
		session := websocketSession{
			user:   "test",
			events: messages,
		}

		err = sessionContainer.register(
			sessionId,
			session,
		)
		if err != nil {
			logger.Sugar().Errorw(
				"Can not register WebSocket session",
				"request", r,
				"error", err,
			)
			http.Error(w, "cannot register connection session", http.StatusUpgradeRequired)
		}
		defer sessionContainer.deregister(sessionId)

		for message := range messages {
			messageJson, err := json.Marshal(message)
			if err != nil {
				logger.Sugar().Errorw(
					"Cannot marshall WebSocket message",
					"message", message,
					"session", session,
				)
			}
			if err := ws.WriteMessage(websocket.TextMessage, messageJson); err != nil {
				log.Println(err)
				return
			}
		}
	}
}
