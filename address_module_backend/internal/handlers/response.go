package handlers

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse sends a structured error response
func ErrorResponse(w http.ResponseWriter, statusCode int, errType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   errType,
		"message": message,
	})
}
