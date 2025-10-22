package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"address_module/internal/tools"

	log "github.com/sirupsen/logrus"
)

// AddDevice handles creating a new device
func AddDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only POST allowed")
		return
	}

	// ✅ Read the body into memory
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorResponse(w, http.StatusInternalServerError, "read_error", "Could not read request body")
		return
	}

	// ✅ Log raw JSON for inspection
	log.Printf("📥 Received JSON:\n%s", string(bodyBytes))

	// ✅ Decode with re-created reader
	var device tools.DeviceParams
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&device); err != nil {
		log.Printf("❌ JSON decode error: %v", err)

		// Optional: return detailed error for frontend dev
		ErrorResponse(w, http.StatusBadRequest, "bad_request", fmt.Sprintf("Invalid JSON input: %v", err))
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	id, err := dbi.InsertDevice(&device)
	if err != nil {
		log.Errorf("InsertDevice failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not insert device")
		return
	}

	device.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// GetDeviceByID retrieves a device by ID
func GetDeviceByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing device ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Device ID must be a number")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	device, err := dbi.GetDeviceByID(id)
	if err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "Device not found")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(device)
}

// UpdateDevice modifies an existing device
func UpdateDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only PUT allowed")
		return
	}

	var device tools.DeviceParams
	if err := json.NewDecoder(r.Body).Decode(&device); err != nil {
		ErrorResponse(w, http.StatusBadRequest, "bad_request", "Invalid JSON input")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	if err := dbi.UpdateDevice(&device); err != nil {
		log.Errorf("UpdateDevice failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "update_failed", "Device could not be updated")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Device updated successfully"})
}

// DeleteDevice removes a device
func DeleteDevice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only DELETE allowed")
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		ErrorResponse(w, http.StatusBadRequest, "missing_parameter", "Missing device ID")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		ErrorResponse(w, http.StatusBadRequest, "invalid_id", "Device ID must be a number")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	if err := dbi.DeleteDevice(id); err != nil {
		ErrorResponse(w, http.StatusNotFound, "not_found", "Device not found")
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Device deleted successfully"})
}

// ListDevices returns all devices
func ListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ErrorResponse(w, http.StatusMethodNotAllowed, "invalid_method", "Only GET allowed")
		return
	}

	dbi, ok := getPostgresDBInstance(w)
	if !ok {
		return
	}
	defer dbi.Close()

	devices, err := dbi.GetAllDevices()
	if err != nil {
		log.Errorf("GetAllDevices failed: %v", err)
		ErrorResponse(w, http.StatusInternalServerError, "db_error", "Could not fetch devices")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(devices)
}
