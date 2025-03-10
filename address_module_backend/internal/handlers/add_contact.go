package handlers

import (
	"address_module/api"
	"address_module/internal/tools"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrMissingContactFields = errors.New("missing required contact fields")
var ErrMissingFirmID = errors.New("missing required firm ID")

// AddContact handles contact registration requests
func AddContact(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var params api.ContactParams
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error("Failed to parse request body: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Validate required fields
	if params.Name == "" || params.Email == "" || params.Kontotyp == "" {
		log.Warn(ErrMissingContactFields)
		api.RequestErrorHandler(w, ErrMissingContactFields)
		return
	}

	// Validate firm ID
	if params.FirmaID <= 0 {
		log.Warn(ErrMissingFirmID)
		api.RequestErrorHandler(w, ErrMissingFirmID)
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

	// Convert api.ContactParams to tools.ContactParams
	contact := tools.ContactParams{
		Anrede:     params.Anrede,
		Name:       params.Name,
		Position:   params.Position,
		Telefon:    params.Telefon,
		Mobil:      params.Mobil,
		Email:      params.Email,
		Abteilung:  params.Abteilung,
		Geburtstag: params.Geburtstag,
		Bemerkung:  params.Bemerkung,
		Kontotyp:   params.Kontotyp,
		FirmaID:    params.FirmaID,
	}

	// Insert contact and create relationship
	contactID, err := db.InsertContact(contact)
	if err != nil {
		log.Error("Failed to insert contact data: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Log received data
	log.WithFields(log.Fields{
		"contact_id": contactID,
		"name":       params.Name,
		"email":      params.Email,
		"firma_id":   params.FirmaID,
		"kontotyp":   params.Kontotyp,
	}).Info("Contact successfully registered in database")

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Contact successfully registered",
		"contact_id": contactID,
	})
	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}