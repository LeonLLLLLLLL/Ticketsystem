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

var ErrMissingContactID = errors.New("missing required contact ID")

// GetFirmsByContactID handles GET requests to retrieve all firms associated with a contact
func GetFirmsByContactID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Get contact ID from query parameters
	contactIDStr := r.URL.Query().Get("id")
	if contactIDStr == "" {
		log.Warn(ErrMissingContactID)
		api.RequestErrorHandler(w, ErrMissingContactID)
		return
	}

	// Convert contact ID to int64
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		log.Error("Invalid contact ID format: ", err)
		api.RequestErrorHandler(w, errors.New("invalid contact ID format"))
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

	// Get firms by contact ID
	firms, err := db.GetFirmsByContactID(contactID)
	if err != nil {
		log.Error("Failed to get firms: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Convert tools.FirmParams to api.FirmResponse objects
	var firmResponses []api.FirmResponse
	for _, firm := range firms {
		firmResponse := api.FirmResponse{
			ID:        firm.ID, // Include the ID
			Anrede:    firm.Anrede,
			Name1:     firm.Name1,
			Name2:     firm.Name2,
			Name3:     firm.Name3,
			Straße:    firm.Straße,
			Land:      firm.Land,
			PLZ:       firm.PLZ,
			Ort:       firm.Ort,
			Telefon:   firm.Telefon,
			Email:     firm.Email,
			Website:   firm.Website,
			Kunde:     firm.Kunde,
			Lieferant: firm.Lieferant,
			Gesperrt:  firm.Gesperrt,
			Bemerkung: firm.Bemerkung,
			FirmaTyp:  firm.FirmaTyp,
		}
		firmResponses = append(firmResponses, firmResponse)
	}

	// Response with the list of firms
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	response := map[string]interface{}{
		"firms": firmResponses,
		"count": len(firmResponses),
	}
	
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error("Failed to encode response: ", err)
		api.InternalErrorHandler(w)
		return
	}
}