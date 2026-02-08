package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/PiotrTopa/js8web/model"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Bio      string `json:"bio"`
}

type updateUserRequest struct {
	Role string `json:"role"`
	Bio  string `json:"bio"`
}

type changePasswordRequest struct {
	Password string `json:"password"`
}

type userResponse struct {
	Ok    bool        `json:"ok"`
	Error string      `json:"error,omitempty"`
	User  interface{} `json:"user,omitempty"`
	Users interface{} `json:"users,omitempty"`
}

func apiUsersGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	users, err := model.FetchAllUsers(db)
	if err != nil {
		logger.Sugar().Errorw("Cannot fetch users", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "cannot fetch users"})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Ok: true, Users: users})
}

func apiUsersPost(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	var body createUserRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid request body"})
		return
	}

	body.Username = strings.TrimSpace(body.Username)
	body.Password = strings.TrimSpace(body.Password)
	body.Role = strings.TrimSpace(body.Role)

	if body.Username == "" || body.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "username and password are required"})
		return
	}

	if !model.IsValidRole(body.Role) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid role; must be admin, operator, or monitor"})
		return
	}

	// Check for duplicate username
	existing, _ := model.FetchUserByName(db, body.Username)
	if existing != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "username already exists"})
		return
	}

	user := &model.User{
		Name: body.Username,
		Role: body.Role,
		Bio:  body.Bio,
	}
	user.SetPassword(body.Password)

	if err := user.Insert(db); err != nil {
		logger.Sugar().Errorw("Cannot create user", "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "cannot create user"})
		return
	}

	logger.Sugar().Infow("User created", "username", user.Name, "role", user.Role)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userResponse{Ok: true, User: user.Public()})
}

// extractUserID parses the user ID from the URL path.
// Expected paths: /api/users/123 or /api/users/123/password
func extractUserID(path string) (int64, error) {
	// Strip prefix /api/users/
	rest := strings.TrimPrefix(path, "/api/users/")
	// Take only the numeric part (before any further /)
	parts := strings.SplitN(rest, "/", 2)
	return strconv.ParseInt(parts[0], 10, 64)
}

func apiUserGet(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	id, err := extractUserID(req.URL.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid user ID"})
		return
	}

	user, err := model.FetchUserById(db, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "user not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Ok: true, User: user.Public()})
}

func apiUserPut(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	id, err := extractUserID(req.URL.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid user ID"})
		return
	}

	var body updateUserRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid request body"})
		return
	}

	body.Role = strings.TrimSpace(body.Role)
	if !model.IsValidRole(body.Role) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid role; must be admin, operator, or monitor"})
		return
	}

	if err := model.UpdateUser(db, id, body.Role, body.Bio); err != nil {
		logger.Sugar().Errorw("Cannot update user", "id", id, "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "cannot update user"})
		return
	}

	logger.Sugar().Infow("User updated", "id", id, "role", body.Role)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Ok: true})
}

func apiUserDelete(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	id, err := extractUserID(req.URL.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid user ID"})
		return
	}

	// Prevent deleting yourself
	s, _ := getSessionFromRequest(req)
	user, err := model.FetchUserById(db, id)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "user not found"})
		return
	}
	if user.Name == s.username {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "cannot delete your own account"})
		return
	}

	if err := model.DeleteUser(db, id); err != nil {
		logger.Sugar().Errorw("Cannot delete user", "id", id, "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "cannot delete user"})
		return
	}

	logger.Sugar().Infow("User deleted", "id", id, "username", user.Name)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Ok: true})
}

func apiUserPasswordPut(w http.ResponseWriter, req *http.Request, db *sql.DB) {
	id, err := extractUserID(req.URL.Path)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid user ID"})
		return
	}

	var body changePasswordRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "invalid request body"})
		return
	}

	body.Password = strings.TrimSpace(body.Password)
	if body.Password == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "password is required"})
		return
	}

	if err := model.UpdateUserPassword(db, id, body.Password); err != nil {
		logger.Sugar().Errorw("Cannot update password", "id", id, "error", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(userResponse{Ok: false, Error: "cannot update password"})
		return
	}

	logger.Sugar().Infow("User password updated", "id", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse{Ok: true})
}
