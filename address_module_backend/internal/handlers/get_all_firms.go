package handlers

import (
	"address_module/api"
	"address_module/internal/model"
	"address_module/internal/tools"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// GetAllFirms handles GET requests to retrieve all firms
func GetAllFirms(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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

	// Get all firms
	firms, err := db.GetAllFirms()
	if err != nil {
		log.Error("Failed to get firms: ", err)
		api.InternalErrorHandler(w)
		return
	}

	// Convert tools.FirmParams to api.FirmResponse objects
	var firmResponses []model.FirmResponse
	for _, firm := range firms {
		firmResponse := model.FirmResponse{
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
