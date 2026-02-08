package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/PiotrTopa/js8web/model"
	"github.com/google/uuid"
)

const SESSION_COOKIE_NAME = "js8web_session"
const SESSION_MAX_AGE = 24 * time.Hour

type session struct {
	username  string
	role      string
	createdAt time.Time
}

var (
	sessionsMu sync.RWMutex
	sessions   = make(map[string]session)
)

func createSession(user *model.User) string {
	token := uuid.New().String()
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	sessions[token] = session{
		username:  user.Name,
		role:      user.Role,
		createdAt: time.Now(),
	}
	return token
}

func getSession(token string) (session, bool) {
	sessionsMu.RLock()
	s, ok := sessions[token]
	sessionsMu.RUnlock()
	if !ok {
		return s, false
	}
	if time.Since(s.createdAt) > SESSION_MAX_AGE {
		sessionsMu.Lock()
		delete(sessions, token)
		sessionsMu.Unlock()
		return s, false
	}
	return s, true
}

func deleteSession(token string) {
	sessionsMu.Lock()
	defer sessionsMu.Unlock()
	delete(sessions, token)
}

func getSessionFromRequest(r *http.Request) (session, bool) {
	cookie, err := r.Cookie(SESSION_COOKIE_NAME)
	if err != nil {
		return session{}, false
	}
	return getSession(cookie.Value)
}

// authRequired wraps an http.HandlerFunc and rejects unauthenticated requests.
func authRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, ok := getSessionFromRequest(r)
		if !ok {
			http.Error(w, `{"error":"authentication required"}`, http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

type loginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type authResponse struct {
	Ok       bool   `json:"ok"`
	Username string `json:"username,omitempty"`
	Role     string `json:"role,omitempty"`
	Error    string `json:"error,omitempty"`
}

func apiAuthLoginPost(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var lr loginRequest
	if err := json.NewDecoder(req.Body).Decode(&lr); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(authResponse{Ok: false, Error: "invalid request body"})
		return
	}

	user, err := model.FetchUserByName(db, lr.Username)
	if err != nil || !user.CheckPassword(lr.Password) {
		logger.Sugar().Warnw("Failed login attempt", "username", lr.Username)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(authResponse{Ok: false, Error: "invalid username or password"})
		return
	}

	token := createSession(user)
	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    token,
		Path:     "/",
		MaxAge:   int(SESSION_MAX_AGE.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	logger.Sugar().Infow("User logged in", "username", user.Name, "role", user.Role)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{Ok: true, Username: user.Name, Role: user.Role})
}

func apiAuthLogoutPost(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	cookie, err := req.Cookie(SESSION_COOKIE_NAME)
	if err == nil {
		deleteSession(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     SESSION_COOKIE_NAME,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(authResponse{Ok: true})
}

func apiAuthCheckGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	s, ok := getSessionFromRequest(req)
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		json.NewEncoder(w).Encode(authResponse{Ok: false})
		return
	}
	json.NewEncoder(w).Encode(authResponse{Ok: true, Username: s.username, Role: s.role})
}
