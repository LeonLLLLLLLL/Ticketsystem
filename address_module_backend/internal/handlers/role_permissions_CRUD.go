package handlers

import (
	"address_module/internal/model"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// AssignRolePermission assigns a permission to a role
func AssignRolePermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	var rp model.RolePermission
	if err := json.NewDecoder(r.Body).Decode(&rp); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.InsertRolePermission(rp); err != nil {
		log.Errorf("Failed to assign role-permission: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not assign permission to role")
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Permission assigned to role successfully",
	})
}

// RemoveRolePermission deletes a role-permission mapping
func RemoveRolePermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	roleIDStr := r.URL.Query().Get("role_id")
	permIDStr := r.URL.Query().Get("permission_id")

	roleID, err1 := strconv.ParseInt(roleIDStr, 10, 64)
	permID, err2 := strconv.ParseInt(permIDStr, 10, 64)

	if err1 != nil || err2 != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Invalid role_id or permission_id")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.DeleteRolePermission(roleID, permID); err != nil {
		log.Errorf("Failed to remove role-permission: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not remove role-permission mapping")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Role-permission mapping removed successfully",
	})
}
