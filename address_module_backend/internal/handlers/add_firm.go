package handlers

import (
	"address_module/api"
	"address_module/internal/tools" // Import the database package
	"encoding/json"
	"errors"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

var ErrMissingFields = errors.New("missing required fields")

// AddFirm handles firm registration requests
func AddFirm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Decode request body
	var params api.FirmParams
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&params)
	if err != nil {
		log.Error("Failed to parse request body: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Validate required fields using a single if statement with OR conditions
	if params.Anrede == "" || params.Name1 == "" || params.PLZ == "" || params.Ort == "" || params.Telefon == "" || params.Email == "" {
		log.Warn(ErrMissingFields)
		api.RequestErrorHandler(w, ErrMissingFields)
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

	// Convert `api.FirmParams` to `tools.FirmParams`
	firm := tools.FirmParams{
		Anrede:    params.Anrede,
		Name1:     params.Name1,
		Name2:     params.Name2,
		Name3:     params.Name3,
		Straße:    params.Straße,
		Land:      params.Land,
		PLZ:       params.PLZ,
		Ort:       params.Ort,
		Telefon:   params.Telefon,
		Email:     params.Email,
		Website:   params.Website,
		Kunde:     params.Kunde,
		Lieferant: params.Lieferant,
		Gesperrt:  params.Gesperrt,
		Bemerkung: params.Bemerkung,
		FirmaTyp:  params.FirmaTyp,
	}

	var firmID int64
	var insertErr error

	// Check if a contact ID was provided to create a relationship
	if params.ContactID > 0 {
		// Insert firm data and create relationship with contact
		firmID, insertErr = db.InsertFirmWithContact(firm, params.ContactID)
		if insertErr != nil {
			log.Error("Failed to insert firm data with contact relationship: ", insertErr)
			api.InternalErrorHandler(w)
			return
		}
		log.WithField("contact_id", params.ContactID).Info("Firm associated with contact")
	} else {
		// Insert firm data only
		firmID, insertErr = db.InsertFirm(firm)
		if insertErr != nil {
			log.Error("Failed to insert firm data: ", insertErr)
			api.InternalErrorHandler(w)
			return
		}
	}

	// Log received data
	log.WithFields(log.Fields{
		"anrede":    params.Anrede,
		"name_1":    params.Name1,
		"name_2":    params.Name2,
		"name_3":    params.Name3,
		"straße":    params.Straße,
		"land":      params.Land,
		"plz":       params.PLZ,
		"ort":       params.Ort,
		"telefon":   params.Telefon,
		"email":     params.Email,
		"website":   params.Website,
		"kunde":     params.Kunde,
		"lieferant": params.Lieferant,
		"gesperrt":  params.Gesperrt,
		"bemerkung": params.Bemerkung,
		"firma_typ": params.FirmaTyp,
	}).Info("Firm successfully registered in database")

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Firm successfully registered",
		"firm_id": firmID,
	})
	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}
