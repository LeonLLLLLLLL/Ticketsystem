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

func AddPermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	var perm model.Permission
	if err := json.NewDecoder(r.Body).Decode(&perm); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	permID, err := db.InsertPermission(perm)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			ErrorResponse(w, http.StatusConflict, "duplicate_entry", "Permission name must be unique")
			return
		}

		log.Errorf("Insert permission failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not insert permission")
		return
	}

	perm.ID = permID
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perm)
}

func GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing permission ID")
		return
	}

	permID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Permission ID must be a number")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	perm, err := db.GetPermissionByID(permID)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "Permission not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(perm)
}

func UpdatePermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only PUT allowed")
		return
	}

	var perm model.Permission
	if err := json.NewDecoder(r.Body).Decode(&perm); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.UpdatePermission(perm); err != nil {
		log.Errorf("Update permission failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "update_failed", "Could not update permission")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Permission updated successfully",
	})
}

func DeletePermission(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing permission ID")
		return
	}

	permID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Permission ID must be a number")
		return
	}

	db, ok := getDBInstance(w)
	if !ok {
		return
	}
	defer db.Close()

	if err := db.DeletePermission(permID); err != nil {
		log.Errorf("Delete permission failed: %v", err)
		ErrorResponse(w, http.StatusNotFound, "not_found", "Permission not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Permission deleted successfully",
	})
}
