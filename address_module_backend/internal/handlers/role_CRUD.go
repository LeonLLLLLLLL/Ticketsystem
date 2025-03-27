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

func AddRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	var role model.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	roleID, err := db.InsertRole(role)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ErrorResponse(w, http.StatusConflict, "duplicate_entry", "Role name must be unique")
			return
		}

		log.Errorf("Insert role failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not insert role")
		return
	}

	role.ID = roleID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func GetRoleByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing role ID")
		return
	}

	roleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Role ID must be a number")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	role, err := db.GetRoleByID(roleID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "Role not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(role)
}

func UpdateRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only PUT allowed")
		return
	}

	var role model.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.UpdateRole(role); err != nil {
		log.Errorf("Update role failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "update_failed", "Could not update role")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Role updated successfully",
	})
}

func DeleteRole(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing role ID")
		return
	}

	roleID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Role ID must be a number")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.DeleteRole(roleID); err != nil {
		log.Errorf("Delete role failed: %v", err)
		ErrorResponse(w, http.StatusNotFound, "not_found", "Role not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Role deleted successfully",
	})
}
