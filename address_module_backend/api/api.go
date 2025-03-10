package api

import (
	"encoding/json"
	"net/http"
)

// Coin Balance Params
type CoinBalanceParams struct {
	Username string
}

type CoinBalanceRespond struct {
	// Succec Code, usually 200
	Code    int
	Balance int64
}

// FirmParams represents the expected structure of the request body
type FirmParams struct {
	Anrede     string `json:"anrede"`
	Name1      string `json:"name_1"`
	Name2      string `json:"name_2"`
	Name3      string `json:"name_3"`
	Straße     string `json:"straße"`
	Land       string `json:"land"`
	PLZ        string `json:"plz"`
	Ort        string `json:"ort"`
	Telefon    string `json:"telefon"`
	Email      string `json:"email"`
	Website    string `json:"website"`
	Kunde      bool   `json:"kunde"`
	Lieferant  bool   `json:"lieferant"`
	Gesperrt   bool   `json:"gesperrt"`
	Bemerkung  string `json:"bemerkung"`
	FirmaTyp   string `json:"firma_typ"`
	ContactID  int64  `json:"contact_id"` // Optional: ID of a contact to associate with this firm
}

// FirmResponse represents a firm returned in API responses
type FirmResponse struct {
	ID        int64  `json:"id"`
	Anrede    string `json:"anrede"`
	Name1     string `json:"name_1"`
	Name2     string `json:"name_2"`
	Name3     string `json:"name_3"`
	Straße    string `json:"straße"`
	Land      string `json:"land"`
	PLZ       string `json:"plz"`
	Ort       string `json:"ort"`
	Telefon   string `json:"telefon"`
	Email     string `json:"email"`
	Website   string `json:"website"`
	Kunde     bool   `json:"kunde"`
	Lieferant bool   `json:"lieferant"`
	Gesperrt  bool   `json:"gesperrt"`
	Bemerkung string `json:"bemerkung"`
	FirmaTyp  string `json:"firma_typ"`
}

// ContactParams represents the expected structure for a contact request
type ContactParams struct {
	Anrede     string `json:"anrede"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	Telefon    string `json:"telefon"`
	Mobil      string `json:"mobil"`
	Email      string `json:"email"`
	Abteilung  string `json:"abteilung"`
	Geburtstag string `json:"geburtstag"` // Format: YYYY-MM-DD
	Bemerkung  string `json:"bemerkung"`
	Kontotyp   string `json:"kontotyp"`
	FirmaID    int64  `json:"firma_id"` // The ID of the firm this contact belongs to
}

// ContactResponse represents a contact returned in API responses
type ContactResponse struct {
	ID         int64  `json:"id"`
	Anrede     string `json:"anrede"`
	Name       string `json:"name"`
	Position   string `json:"position"`
	Telefon    string `json:"telefon"`
	Mobil      string `json:"mobil"`
	Email      string `json:"email"`
	Abteilung  string `json:"abteilung"`
	Geburtstag string `json:"geburtstag"`
	Bemerkung  string `json:"bemerkung"`
	Kontotyp   string `json:"kontotyp"`
}

type Error struct {
	// Error Code
	Code int
	// Error Message
	Message string
}

func writeError(w http.ResponseWriter, message string, code int) {
	resp := Error{
		Code:    code,
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(resp)
}

var (
	RequestErrorHandler = func(w http.ResponseWriter, err error) {
		writeError(w, err.Error(), http.StatusBadRequest)
	}
	InternalErrorHandler = func(w http.ResponseWriter) {
		writeError(w, "An Unexpected Error has Occured.", http.StatusInternalServerError)
	}
)
