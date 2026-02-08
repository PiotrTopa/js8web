package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/PiotrTopa/js8web/model"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

type websocketSession struct {
	user   string
	events chan<- model.WebsocketMessage
}

type websocketSessionContainer struct {
	mu       sync.RWMutex
	sessions map[string]websocketSession
}

func (o *websocketSessionContainer) init() {
	o.sessions = make(map[string]websocketSession)
}

func (o *websocketSessionContainer) register(id string, session websocketSession) error {
	if len(id) == 0 {
		return errors.New("no session id specified")
	}

	o.mu.Lock()
	defer o.mu.Unlock()

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

	o.mu.Lock()
	defer o.mu.Unlock()

	if session, exists := o.sessions[id]; exists {
		close(session.events)
		delete(o.sessions, id)
		return nil
	}

	return errors.New("session with given id has not been registered")
}

func (o *websocketSessionContainer) process(messages <-chan model.WebsocketMessage) {
	for message := range messages {
		o.mu.RLock()
		for _, session := range o.sessions {
			// Non-blocking send to avoid one slow client blocking all others
			select {
			case session.events <- message:
			default:
				logger.Sugar().Warnw("Dropping WebSocket message for slow client")
			}
		}
		o.mu.RUnlock()
	}
}

func websocketHandler(sessionContainer *websocketSessionContainer) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Sugar().Errorw(
				"Cannot upgrade WebSocket connection",
				"error", err,
			)
			return
		}
		defer ws.Close()

		messages := make(chan model.WebsocketMessage, 16)
		sessionId := uuid.New().String()
		session := websocketSession{
			user:   "anonymous",
			events: messages,
		}

		err = sessionContainer.register(sessionId, session)
		if err != nil {
			logger.Sugar().Errorw(
				"Cannot register WebSocket session",
				"error", err,
			)
			return
		}
		defer sessionContainer.deregister(sessionId)

		logger.Sugar().Infow("WebSocket client connected", "sessionId", sessionId)

		for message := range messages {
			messageJson, err := json.Marshal(message)
			if err != nil {
				logger.Sugar().Errorw(
					"Cannot marshal WebSocket message",
					"error", err,
				)
				continue
			}
			if err := ws.WriteMessage(websocket.TextMessage, messageJson); err != nil {
				logger.Sugar().Infow("WebSocket client disconnected", "sessionId", sessionId)
				return
			}
		}
	}
}
