package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"address_module/internal/tools"

	log "github.com/sirupsen/logrus"
)

// AddDeviceLink creates a new link between two devices
func AddDeviceLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	var link tools.DeviceLink
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	id, err := dbi.InsertDeviceLink(&link)
	if err != nil {
		log.Errorf("InsertDeviceLink failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not insert device link")
		return
	}

	link.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(link)
}

// UpdateDeviceLink modifies an existing device link
func UpdateDeviceLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only PUT allowed")
		return
	}

	var link tools.DeviceLink
	if err := json.NewDecoder(r.Body).Decode(&link); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	if err := dbi.UpdateDeviceLink(&link); err != nil {
		log.Errorf("UpdateDeviceLink failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "update_failed", "Device link could not be updated")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Device link updated successfully"})
}

// DeleteDeviceLink removes a link by ID
func DeleteDeviceLink(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing link ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Link ID must be a number")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	if err := dbi.DeleteDeviceLink(id); err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "Link not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Link deleted successfully"})
}

// ListDeviceLinks returns all links
func ListDeviceLinks(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	links, err := dbi.GetAllDeviceLinks()
	if err != nil {
		log.Errorf("GetAllDeviceLinks failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not fetch links")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}

// GetDeviceLinksForDevice returns links for a given device
func GetDeviceLinksForDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	idStr := r.URL.Query().Get("device_id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing device_id")
		return
	}

	deviceID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "device_id must be numeric")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	links, err := dbi.GetLinksForDevice(deviceID)
	if err != nil {
		log.Errorf("GetLinksForDevice failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not fetch device links")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(links)
}
