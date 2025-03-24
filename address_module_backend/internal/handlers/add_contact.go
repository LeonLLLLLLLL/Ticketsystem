package handlers

import (
	"address_module/api"
	"address_module/internal/model"
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
	var params model.ContactParams
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error("Failed to parse request body: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Validate required fields
	if params.Vorname == "" || params.Email == "" /*|| params.Kontotyp == "" */ {
		log.Warn(ErrMissingContactFields)
		api.RequestErrorHandler(w, ErrMissingContactFields)
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
		Vorname:    params.Vorname,
		Nachname:   params.Nachname,
		Position:   params.Position,
		Telefon:    params.Telefon,
		Mobil:      params.Mobil,
		Email:      params.Email,
		Abteilung:  params.Abteilung,
		Geburtstag: params.Geburtstag,
		Bemerkung:  params.Bemerkung,
		Kontotyp:   params.Kontotyp,
	}

	var contactID int64

	// Insert contact with relationships
	if len(params.Firms) > 0 {
		// Insert contact and create relationships with multiple firms
		contactID, err = db.InsertContactWithFirms(contact, params.Firms)
		if err != nil {
			log.Error("Failed to insert contact data with firm relationships: ", err)
			api.InternalErrorHandler(w)
			return
		}
		log.WithField("firm_ids", params.Firms).Info("Contact associated with multiple firms")
	} else {
		//Insert a contact without any firm relationships
		log.Info("Inserting contact without firm relationships")
		contactID, err = db.InsertContact(contact)
		if err != nil {
			log.Error("Failed to insert contact data: ", err)
			api.InternalErrorHandler(w)
			return
		}
	}

	// Log received data
	logFields := log.Fields{
		"contact_id": contactID,
		"vorname":    params.Vorname,
		"email":      params.Email,
		//"kontotyp":   params.Kontotyp,
	}

	// Add firm relationship info to logs
	if len(params.Firms) > 0 {
		logFields["firma_ids"] = params.Firms
	} else {
		logFields["firma_id"] = params.Firms
	}

	log.WithFields(logFields).Info("Contact successfully registered in database")

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"message":    "Contact successfully registered",
		"contact_id": contactID,
	}

	// Add relationship info to response
	if len(params.Firms) > 0 {
		response["firm_ids"] = params.Firms
	} else {
		response["firm_id"] = params.Firms
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}
