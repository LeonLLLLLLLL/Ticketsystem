package handlers

import (
	"address_module/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// AssignUserRole assigns a role to a user
func AssignUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	var ur model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.InsertUserRole(ur); err != nil {
		log.Errorf("Failed to assign user-role: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not assign role to user")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Role assigned to user successfully",
	})
}

// RemoveUserRole deletes a user-role mapping
func RemoveUserRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	roleIDStr := r.URL.Query().Get("role_id")

	userID, err1 := strconv.ParseInt(userIDStr, 10, 64)
	roleID, err2 := strconv.ParseInt(roleIDStr, 10, 64)

	if err1 != nil || err2 != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Invalid user_id or role_id")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.DeleteUserRole(userID, roleID); err != nil {
		log.Errorf("Failed to remove user-role: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not remove user-role mapping")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User-role mapping removed successfully",
	})
}
