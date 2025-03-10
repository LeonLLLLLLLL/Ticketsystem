package handlers

import (
	"address_module/api"
	"address_module/internal/tools"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrMissingFirmIDParam = errors.New("missing required firm ID")

// GetContactsByFirmID handles GET requests to retrieve all contacts associated with a firm
func GetContactsByFirmID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get firm ID from query parameters
	firmIDStr := r.URL.Query().Get("id")
	if firmIDStr == "" {
		log.Warn(ErrMissingFirmIDParam)
		api.RequestErrorHandler(w, ErrMissingFirmIDParam)
		return
	}

	// Convert firm ID to int64
	firmID, err := strconv.ParseInt(firmIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid firm ID format: ", err)
		api.RequestErrorHandler(w, errors.New("invalid firm ID format"))
		return
	}

	// Connect to MySQL database
	db, err := tools.NewDatabase(5, 3*time.Second)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		api.InternalErrorHandler(w)
		return
	}
	defer db.Close()

	// Get contacts by firm ID
	contacts, err := db.GetContactsByFirmID(firmID)
	if err != nil {
		log.Error("Failed to get contacts: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Convert tools.ContactParams to api.ContactResponse objects
	var contactResponses []api.ContactResponse
	for _, contact := range contacts {
		contactResponse := api.ContactResponse{
			ID:         contact.ID, // Include the ID
			Anrede:     contact.Anrede,
			Name:       contact.Name,
			Position:   contact.Position,
			Telefon:    contact.Telefon,
			Mobil:      contact.Mobil,
			Email:      contact.Email,
			Abteilung:  contact.Abteilung,
			Geburtstag: contact.Geburtstag,
			Bemerkung:  contact.Bemerkung,
			Kontotyp:   contact.Kontotyp,
		}
		contactResponses = append(contactResponses, contactResponse)
	}

	// Response with the list of contacts
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"contacts": contactResponses,
		"count":    len(contactResponses),
	}
	
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}