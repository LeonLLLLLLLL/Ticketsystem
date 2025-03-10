package tools

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// DatabaseInterface defines database operations
type DatabaseInterface interface {
	SetupDatabase() error
	InsertFirm(firm FirmParams) (int64, error)
	InsertFirmWithContact(firm FirmParams, contactID int64) (int64, error)
	GetUserLoginDetails(username string) *LoginDetails
	GetUserCoins(username string) *CoinDetails
	GetFirmsByContactID(contactID int64) ([]FirmParams, error)
	GetContactsByFirmID(firmID int64) ([]ContactParams, error)
	GetAllFirms() ([]FirmParams, error)
	GetAllContacts() ([]ContactParams, error)
}

// FirmParams struct matching the MySQL table
type FirmParams struct {
	ID        int64  // Added ID field
	Anrede    string
	Name1     string
	Name2     string
	Name3     string
	Straße    string
	Land      string
	PLZ       string
	Ort       string
	Telefon   string
	Email     string
	Website   string
	Kunde     bool
	Lieferant bool
	Gesperrt  bool
	Bemerkung string
	FirmaTyp  string
}

// LoginDetails struct
type LoginDetails struct {
	AuthToken string
	Username  string
}

// CoinDetails struct
type CoinDetails struct {
	Coins    int64
	Username string
}

type MySQLDB struct {
	DB *sql.DB
}

func debugMySQLEnvironmentVariables() {
	// List of environment variables to check
	envVars := []string{
		"MYSQL_HOST",
		"MYSQL_USER",
		"MYSQL_PASSWORD",
		"MYSQL_DATABASE",
	}

	// Print out each environment variable
	fmt.Println("MySQL Environment Variables Debug:")
	fmt.Println("-----------------------------------")

	for _, varName := range envVars {
		value := os.Getenv(varName)

		// Check if the variable is set
		if value == "" {
			fmt.Printf("%s: <NOT SET>\n", varName)
		} else {
			// For security, partially mask the password
			if varName == "MYSQL_PASSWORD" {
				maskedValue := "****" + value[len(value)-4:]
				fmt.Printf("%s: %s\n", varName, maskedValue)
			} else {
				fmt.Printf("%s: %s\n", varName, value)
			}
		}
	}

	// Print all environment variables (optional, can be very verbose)
	fmt.Println("\nFull Environment Variables:")
	fmt.Println("-----------------------------------")
	for _, env := range os.Environ() {
		fmt.Println(env)
	}
}

// NewDatabase initializes a MySQL connection
func NewDatabase(maxRetries int, delay time.Duration) (*MySQLDB, error) {
	host := os.Getenv("MYSQL_HOST")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")
	dbName := os.Getenv("MYSQL_DATABASE")
	port := os.Getenv("MYSQL_PORT")

	if port == "" {
		port = "3306"
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			user,
			password,
			host,
			port,
			dbName,
		)

		// Open database connection
		sqlDB, err := sql.Open("mysql", dsn)
		if err != nil {
			log.Errorf("Failed to open database (Attempt %d/%d): %v", attempt, maxRetries, err)
			return nil, err
		}

		// Test the connection
		err = sqlDB.Ping()
		if err == nil {
			log.Infof("Successfully connected to MySQL database on attempt %d", attempt)
			return &MySQLDB{DB: sqlDB}, nil
		}

		// Close the connection if ping fails
		sqlDB.Close()

		log.Errorf("Database connection test failed (Attempt %d/%d): %v", attempt, maxRetries, err)

		// Sleep before next attempt
		time.Sleep(delay)
	}

	return nil, fmt.Errorf("failed to connect to database after %d attempts with host %s", maxRetries, host)
}

func (m *MySQLDB) Close() error {
	if m.DB != nil {
		err := m.DB.Close()
		if err != nil {
			log.Printf("Error closing database connection: %v", err)
			return err
		}
		log.Println("Database connection closed successfully")
	}
	return nil
}

// InsertFirm inserts firm data into MySQL and returns the firm ID
func (db *MySQLDB) InsertFirm(firm FirmParams) (int64, error) {
	query := `
	INSERT INTO firms (anrede, name_1, name_2, name_3, straße, land, plz, ort, telefon, email, website, kunde, lieferant, gesperrt, bemerkung, firma_typ) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.DB.Exec(query, firm.Anrede, firm.Name1, firm.Name2, firm.Name3, firm.Straße, firm.Land, firm.PLZ, firm.Ort, firm.Telefon, firm.Email, firm.Website, firm.Kunde, firm.Lieferant, firm.Gesperrt, firm.Bemerkung, firm.FirmaTyp)
	if err != nil {
		log.Error("Failed to insert firm: ", err)
		return 0, err
	}

	// Get the firm ID
	firmID, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get firm ID: ", err)
		return 0, err
	}

	log.Info("Firm inserted successfully with ID: ", firmID)
	return firmID, nil
}

// InsertFirmWithContact inserts a firm and creates a relationship with a contact
func (db *MySQLDB) InsertFirmWithContact(firm FirmParams, contactID int64) (int64, error) {
	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Error("Failed to begin transaction: ", err)
		return 0, err
	}

	// Step 1: Insert the firm
	firmQuery := `
	INSERT INTO firms (anrede, name_1, name_2, name_3, straße, land, plz, ort, telefon, email, website, kunde, lieferant, gesperrt, bemerkung, firma_typ) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	firmResult, err := tx.Exec(firmQuery, firm.Anrede, firm.Name1, firm.Name2, firm.Name3, firm.Straße, firm.Land, firm.PLZ, firm.Ort, firm.Telefon, firm.Email, firm.Website, firm.Kunde, firm.Lieferant, firm.Gesperrt, firm.Bemerkung, firm.FirmaTyp)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to insert firm: ", err)
		return 0, err
	}

	// Get the firm ID
	firmID, err := firmResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Error("Failed to get firm ID: ", err)
		return 0, err
	}

	// Step 2: Create relationship with contact
	relationQuery := `
	INSERT INTO firms_contacts (firma_id, contact_id) 
	VALUES (?, ?)`

	_, err = tx.Exec(relationQuery, firmID, contactID)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to create relationship: ", err)
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: ", err)
		return 0, err
	}

	log.WithFields(log.Fields{
		"firm_id":    firmID,
		"contact_id": contactID,
	}).Info("Firm inserted and relationship created successfully")

	return firmID, nil
}

// ContactParams struct matching the MySQL table
type ContactParams struct {
	ID         int64  // Added ID field
	Anrede     string
	Name       string
	Position   string
	Telefon    string
	Mobil      string
	Email      string
	Abteilung  string
	Geburtstag string // Store as string in YYYY-MM-DD format
	Bemerkung  string
	Kontotyp   string
	FirmaID    int64 // The ID of the firm this contact belongs to
}

// InsertContact inserts contact data into MySQL and creates relationship with firm
func (db *MySQLDB) InsertContact(contact ContactParams) (int64, error) {
	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Error("Failed to begin transaction: ", err)
		return 0, err
	}

	// Step 1: Insert the contact
	contactQuery := `
	INSERT INTO contacts (anrede, name, position, telefon, mobil, email, abteilung, geburtstag, bemerkung, kontotyp) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	contactResult, err := tx.Exec(contactQuery,
		contact.Anrede,
		contact.Name,
		contact.Position,
		contact.Telefon,
		contact.Mobil,
		contact.Email,
		contact.Abteilung,
		contact.Geburtstag,
		contact.Bemerkung,
		contact.Kontotyp)

	if err != nil {
		tx.Rollback()
		log.Error("Failed to insert contact: ", err)
		return 0, err
	}

	// Get the contact ID
	contactID, err := contactResult.LastInsertId()
	if err != nil {
		tx.Rollback()
		log.Error("Failed to get contact ID: ", err)
		return 0, err
	}

	// Step 2: Create relationship with firm
	relationQuery := `
	INSERT INTO firms_contacts (firma_id, contact_id) 
	VALUES (?, ?)`

	_, err = tx.Exec(relationQuery, contact.FirmaID, contactID)
	if err != nil {
		tx.Rollback()
		log.Error("Failed to create relationship: ", err)
		return 0, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: ", err)
		return 0, err
	}

	return contactID, nil
}

// GetFirmsByContactID retrieves all firms associated with a specific contact
func (db *MySQLDB) GetFirmsByContactID(contactID int64) ([]FirmParams, error) {
	query := `
	SELECT f.id, f.anrede, f.name_1, f.name_2, f.name_3, f.straße, f.land, 
	       f.plz, f.ort, f.telefon, f.email, f.website, f.kunde, 
	       f.lieferant, f.gesperrt, f.bemerkung, f.firma_typ
	FROM firms f
	JOIN firms_contacts fc ON f.id = fc.firma_id
	WHERE fc.contact_id = ?`

	rows, err := db.DB.Query(query, contactID)
	if err != nil {
		log.Error("Failed to query firms by contact ID: ", err)
		return nil, err
	}
	defer rows.Close()

	var firms []FirmParams
	for rows.Next() {
		var firm FirmParams
		err := rows.Scan(
			&firm.ID, &firm.Anrede, &firm.Name1, &firm.Name2, &firm.Name3,
			&firm.Straße, &firm.Land, &firm.PLZ, &firm.Ort, &firm.Telefon,
			&firm.Email, &firm.Website, &firm.Kunde, &firm.Lieferant,
			&firm.Gesperrt, &firm.Bemerkung, &firm.FirmaTyp,
		)
		if err != nil {
			log.Error("Failed to scan firm row: ", err)
			return nil, err
		}
		firms = append(firms, firm)
	}

	if err = rows.Err(); err != nil {
		log.Error("Error during rows iteration: ", err)
		return nil, err
	}

	return firms, nil
}

// GetContactsByFirmID retrieves all contacts associated with a specific firm
func (db *MySQLDB) GetContactsByFirmID(firmID int64) ([]ContactParams, error) {
	query := `
	SELECT c.id, c.anrede, c.name, c.position, c.telefon, c.mobil, 
	       c.email, c.abteilung, c.geburtstag, c.bemerkung, c.kontotyp
	FROM contacts c
	JOIN firms_contacts fc ON c.id = fc.contact_id
	WHERE fc.firma_id = ?`

	rows, err := db.DB.Query(query, firmID)
	if err != nil {
		log.Error("Failed to query contacts by firm ID: ", err)
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactParams
	for rows.Next() {
		var contact ContactParams
		err := rows.Scan(
			&contact.ID, &contact.Anrede, &contact.Name, &contact.Position,
			&contact.Telefon, &contact.Mobil, &contact.Email, &contact.Abteilung,
			&contact.Geburtstag, &contact.Bemerkung, &contact.Kontotyp,
		)
		if err != nil {
			log.Error("Failed to scan contact row: ", err)
			return nil, err
		}
		contact.FirmaID = firmID // Set the firm ID in the contact
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		log.Error("Error during rows iteration: ", err)
		return nil, err
	}

	return contacts, nil
}

// GetAllFirms retrieves all firms from the database
func (db *MySQLDB) GetAllFirms() ([]FirmParams, error) {
	query := `
	SELECT id, anrede, name_1, name_2, name_3, straße, land, 
	       plz, ort, telefon, email, website, kunde, 
	       lieferant, gesperrt, bemerkung, firma_typ
	FROM firms
	ORDER BY id DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Error("Failed to query all firms: ", err)
		return nil, err
	}
	defer rows.Close()

	var firms []FirmParams
	for rows.Next() {
		var firm FirmParams
		err := rows.Scan(
			&firm.ID, &firm.Anrede, &firm.Name1, &firm.Name2, &firm.Name3,
			&firm.Straße, &firm.Land, &firm.PLZ, &firm.Ort, &firm.Telefon,
			&firm.Email, &firm.Website, &firm.Kunde, &firm.Lieferant,
			&firm.Gesperrt, &firm.Bemerkung, &firm.FirmaTyp,
		)
		if err != nil {
			log.Error("Failed to scan firm row: ", err)
			return nil, err
		}
		firms = append(firms, firm)
	}

	if err = rows.Err(); err != nil {
		log.Error("Error during rows iteration: ", err)
		return nil, err
	}

	return firms, nil
}

// GetAllContacts retrieves all contacts from the database
func (db *MySQLDB) GetAllContacts() ([]ContactParams, error) {
	query := `
	SELECT id, anrede, name, position, telefon, mobil, 
	       email, abteilung, geburtstag, bemerkung, kontotyp
	FROM contacts
	ORDER BY id DESC`

	rows, err := db.DB.Query(query)
	if err != nil {
		log.Error("Failed to query all contacts: ", err)
		return nil, err
	}
	defer rows.Close()

	var contacts []ContactParams
	for rows.Next() {
		var contact ContactParams
		err := rows.Scan(
			&contact.ID, &contact.Anrede, &contact.Name, &contact.Position, 
			&contact.Telefon, &contact.Mobil, &contact.Email, &contact.Abteilung, 
			&contact.Geburtstag, &contact.Bemerkung, &contact.Kontotyp,
		)
		if err != nil {
			log.Error("Failed to scan contact row: ", err)
			return nil, err
		}
		contacts = append(contacts, contact)
	}

	if err = rows.Err(); err != nil {
		log.Error("Error during rows iteration: ", err)
		return nil, err
	}

	return contacts, nil
}

/*
TEST FUNCTION REMOVE LATER!!!
*/

// InsertTestData inserts predefined test data into the database
func (db *MySQLDB) InsertTestData() error {
	// Start a transaction to ensure data integrity
	tx, err := db.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Firms data (unchanged)
	firmQuery := `
		INSERT INTO firms (anrede, name_1, name_2, name_3, straße, land, plz, ort, telefon, email, website, kunde, lieferant, gesperrt, bemerkung, firma_typ)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	firms := [][]interface{}{
		{"Herr", "Musterfirma GmbH", "Abteilung Vertrieb", "", "Musterstraße 1", "Deutschland", "10115", "Berlin", "+49 30 123456", "info@musterfirma.de", "https://www.musterfirma.de", true, false, false, "Großkunde", "GmbH"},
		{"Frau", "Tech Solutions AG", "", "", "Innovationsweg 2", "Deutschland", "80331", "München", "+49 89 987654", "contact@techsolutions.de", "https://www.techsolutions.de", true, true, false, "IT-Dienstleister", "AG"},
		{"Herr", "Baustoffe Schmidt KG", "", "", "Bauallee 5", "Deutschland", "50667", "Köln", "+49 221 456789", "info@baustoffe-schmidt.de", "", true, true, false, "Baumateriallieferant", "KG"},
	}

	// Prepare the firm insert statement
	firmStmt, err := tx.Prepare(firmQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare firms statement: %v", err)
	}
	defer firmStmt.Close()

	// Insert firms with prepared statement
	var firmIDs []int64
	for _, firm := range firms {
		result, err := firmStmt.Exec(firm...)
		if err != nil {
			return fmt.Errorf("error inserting firm %v: %v", firm[1], err)
		}
		firmID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID for firm %v: %v", firm[1], err)
		}
		firmIDs = append(firmIDs, firmID)
	}

	// Contacts data with truncated kontotyp
	contactQuery := `
		INSERT INTO contacts (anrede, name, position, telefon, mobil, email, abteilung, geburtstag, bemerkung, kontotyp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	contacts := [][]interface{}{
		{"Herr", "Max Mustermann", "Vertriebsleiter", "+49 30 123457", "+49 170 1234567", "max.mustermann@musterfirma.de", "Vertrieb", "1985-04-12", "Hauptansprechpartner für Vertrieb", "K"}, // Changed from "Kunde"
		{"Frau", "Lisa Meier", "IT-Projektmanagerin", "+49 89 987655", "+49 172 9876543", "lisa.meier@techsolutions.de", "Projektmanagement", "1990-09-25", "Leitet Softwareprojekte", "L"},    // Changed from "Lieferant"
		{"Herr", "Johann Schmidt", "Einkaufsleiter", "+49 221 456788", "", "johann.schmidt@baustoffe-schmidt.de", "Einkauf", "1978-06-30", "Verantwortlich für Materialbeschaffung", "L"},
		{"Herr", "Michael König", "Key Account Manager", "+49 30 7654321", "+49 172 7654321", "michael.koenig@musterfirma.de", "Vertrieb", "1982-11-05", "Betreut Großkunden", "K"},
		{"Frau", "Anna Weber", "Kundensupport", "+49 30 2345678", "+49 173 2345678", "anna.weber@musterfirma.de", "Support", "1995-07-20", "Support für Geschäftskunden", "K"},
		{"Frau", "Sarah Bauer", "Finanzleiterin", "+49 89 1234599", "+49 175 1234599", "sarah.bauer@finanzexperten.de", "Finanzen", "1983-03-15", "Buchhaltung und Finanzen", "L"},
	}

	// Prepare the contact insert statement
	contactStmt, err := tx.Prepare(contactQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare contacts statement: %v", err)
	}
	defer contactStmt.Close()

	// Insert contacts with prepared statement
	var contactIDs []int64
	for _, contact := range contacts {
		result, err := contactStmt.Exec(contact...)
		if err != nil {
			return fmt.Errorf("error inserting contact %v: %v", contact[1], err)
		}
		contactID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID for contact %v: %v", contact[1], err)
		}
		contactIDs = append(contactIDs, contactID)
	}

	// Firms-Contacts relationships data (unchanged)
	firmContactQuery := `
		INSERT INTO firms_contacts (firma_id, contact_id, beziehung, hauptansprechpartner)
		VALUES (?, ?, ?, ?);
	`
	firmsContacts := [][]interface{}{
		{firmIDs[0], contactIDs[0], "Hauptkontakt für Bestellungen", true},            // Max Mustermann -> Musterfirma GmbH
		{firmIDs[1], contactIDs[1], "Projektleitung IT-Systeme", true},                // Lisa Meier -> Tech Solutions
		{firmIDs[2], contactIDs[2], "Ansprechpartner für Materiallieferungen", false}, // Johann Schmidt -> Baustoffe Schmidt KG
		{firmIDs[0], contactIDs[3], "Verantwortlich für Großkundenbetreuung", false},  // Michael König -> Musterfirma GmbH
		{firmIDs[0], contactIDs[4], "Support für Bestandskunden", false},              // Anna Weber -> Musterfirma GmbH
		{firmIDs[1], contactIDs[5], "Finanzierung und Rechnungswesen", false},         // Sarah Bauer -> Tech Solutions
		{firmIDs[2], contactIDs[5], "Beratung zur Kostenoptimierung", false},          // Sarah Bauer -> Baustoffe Schmidt KG
	}

	// Prepare the firm-contact relationship insert statement
	firmContactStmt, err := tx.Prepare(firmContactQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare firms_contacts statement: %v", err)
	}
	defer firmContactStmt.Close()

	// Insert firm-contact relationships
	for _, relation := range firmsContacts {
		_, err := firmContactStmt.Exec(relation...)
		if err != nil {
			return fmt.Errorf("error inserting firm-contact relation %v: %v", relation, err)
		}
	}

	log.Println("Test data successfully inserted into database.")
	return nil
}
