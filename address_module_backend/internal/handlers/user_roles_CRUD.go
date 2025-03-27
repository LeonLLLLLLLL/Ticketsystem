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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var ur model.UserRole
	if err := json.NewDecoder(r.Body).Decode(&ur); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.InsertUserRole(ur); err != nil {
		log.Errorf("Failed to assign user-role: %v", err)
		http.Error(w, "Assignment failed", http.StatusInternalServerError)
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
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("user_id")
	roleIDStr := r.URL.Query().Get("role_id")

	userID, err1 := strconv.ParseInt(userIDStr, 10, 64)
	roleID, err2 := strconv.ParseInt(roleIDStr, 10, 64)

	if err1 != nil || err2 != nil {
		http.Error(w, "Invalid user_id or role_id", http.StatusBadRequest)
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.DeleteUserRole(userID, roleID); err != nil {
		log.Errorf("Failed to remove user-role: %v", err)
		http.Error(w, "Removal failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
