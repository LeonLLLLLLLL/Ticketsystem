package handlers

import (
	"address_module/internal/model"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// AddUser creates a new user
func AddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	userID, err := db.InsertUser(user)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ErrorResponse(w, http.StatusConflict, "duplicate_entry", "Username or email already exists")
			return
		}

		log.Errorf("Failed to insert user: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not insert user")
		return
	}

	user.ID = userID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetUserByID fetches a user by ID
func GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing user ID")
		return
	}

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "User ID must be a number")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	user, err := db.GetUserByID(userID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// UpdateUser updates an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only PUT allowed")
		return
	}

	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.UpdateUser(user); err != nil {
		log.Errorf("Update failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "update_failed", "User could not be updated")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User updated successfully",
	})
}

// DeleteUser removes a user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing user ID")
		return
	}

	userID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "User ID must be a number")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.DeleteUser(userID); err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "User not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User deleted successfully",
	})
}
