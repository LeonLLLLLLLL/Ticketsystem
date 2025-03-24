package model

// FirmParams represents the expected structure of the request body
type FirmParams struct {
	Anrede     string  `json:"anrede"`
	Name1      string  `json:"name_1"`
	Name2      string  `json:"name_2"`
	Name3      string  `json:"name_3"`
	Straße     string  `json:"straße"`
	Land       string  `json:"land"`
	PLZ        string  `json:"plz"`
	Ort        string  `json:"ort"`
	Telefon    string  `json:"telefon"`
	Email      string  `json:"email"`
	Website    string  `json:"website"`
	Kunde      bool    `json:"kunde"`
	Lieferant  bool    `json:"lieferant"`
	Gesperrt   bool    `json:"gesperrt"`
	Bemerkung  string  `json:"bemerkung"`
	FirmaTyp   string  `json:"firma_typ"`
	ContactID  int64   `json:"contact_id"`  // Optional: ID of a contact to associate with this firm (for backward compatibility)
	ContactIDs []int64 `json:"contact_ids"` // Optional: IDs of contacts to associate with this firm
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
