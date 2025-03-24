package model

// ContactParams represents the expected structure for a contact request
type ContactParams struct {
	Anrede     string  `json:"anrede"`
	Vorname    string  `json:"vorname"`
	Nachname   string  `json:"nachname"`
	Position   string  `json:"position"`
	Telefon    string  `json:"telefon"`
	Mobil      string  `json:"mobil"`
	Email      string  `json:"email"`
	Abteilung  string  `json:"abteilung"`
	Geburtstag string  `json:"geburtstag"` // Format: YYYY-MM-DD
	Bemerkung  string  `json:"bemerkung"`
	Kontotyp   string  `json:"kontotyp"`
	Firms      []int64 `json:"firms"` // IDs of firms this contact belongs to
}

// ContactResponse represents a contact returned in API responses
type ContactResponse struct {
	ID         int64  `json:"id"`
	Anrede     string `json:"anrede"`
	Vorname    string `json:"vorname"`
	Nachname   string `json:"nachname"`
	Position   string `json:"position"`
	Telefon    string `json:"telefon"`
	Mobil      string `json:"mobil"`
	Email      string `json:"email"`
	Abteilung  string `json:"abteilung"`
	Geburtstag string `json:"geburtstag"`
	Bemerkung  string `json:"bemerkung"`
	Kontotyp   string `json:"kontotyp"`
}
