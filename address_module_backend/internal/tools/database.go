package tools

import (
	"address_module/internal/model"
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
	InsertFirmWithContacts(firm FirmParams, contactIDs []int64) (int64, error)
	InsertContactWithFirms(contact ContactParams, firmIDs []int64) (int64, error)
	CreateContactFirmRelationship(contactID int64, firmID int64) error
	GetUserLoginDetails(username string) *LoginDetails
	GetUserCoins(username string) *CoinDetails
	GetFirmsByContactID(contactID int64) ([]FirmParams, error)
	GetContactsByFirmID(firmID int64) ([]ContactParams, error)
	GetAllFirms() ([]FirmParams, error)
	GetAllContacts() ([]ContactParams, error)
}

// FirmParams struct matching the MySQL table
type FirmParams struct {
	ID        int64 // Added ID field
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
		port = "3307"
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

// InsertFirmWithContacts inserts a firm and creates relationships with multiple contacts
func (db *MySQLDB) InsertFirmWithContacts(firm FirmParams, contactIDs []int64) (int64, error) {
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

	// Step 2: Create relationships with all contacts
	relationQuery := `
	INSERT INTO firms_contacts (firma_id, contact_id) 
	VALUES (?, ?)`

	for _, contactID := range contactIDs {
		_, err = tx.Exec(relationQuery, firmID, contactID)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to create relationship with contact ID ", contactID, ": ", err)
			return 0, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: ", err)
		return 0, err
	}

	log.WithFields(log.Fields{
		"firm_id":     firmID,
		"contact_ids": contactIDs,
	}).Info("Firm inserted and relationships created successfully")

	return firmID, nil
}

// CreateContactFirmRelationship creates a relationship between a contact and a firm
func (db *MySQLDB) CreateContactFirmRelationship(contactID int64, firmID int64) error {
	query := `
	INSERT INTO firms_contacts (firma_id, contact_id) 
	VALUES (?, ?)`

	_, err := db.DB.Exec(query, firmID, contactID)
	if err != nil {
		log.Error("Failed to create relationship: ", err)
		return err
	}

	log.WithFields(log.Fields{
		"contact_id": contactID,
		"firm_id":    firmID,
	}).Info("Relationship created successfully")

	return nil
}

// ContactParams struct matching the MySQL table
type ContactParams struct {
	ID         int64 // Added ID field
	Anrede     string
	Vorname    string
	Nachname   string
	Position   string
	Telefon    string
	Mobil      string
	Email      string
	Abteilung  string
	Geburtstag string // Store as string in YYYY-MM-DD format
	Bemerkung  string
	Kontotyp   string
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
	INSERT INTO contacts (anrede, vorname, nachname, position, telefon, mobil, email, abteilung, geburtstag, bemerkung, kontotyp)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	contactResult, err := tx.Exec(contactQuery,
		contact.Anrede,
		contact.Vorname,
		contact.Nachname,
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

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: ", err)
		return 0, err
	}

	log.WithFields(log.Fields{
		"contact_id": contactID,
	}).Info("Contact inserted successfully")

	return contactID, nil
}

// InsertContactWithFirms inserts a contact and creates relationships with multiple firms
func (db *MySQLDB) InsertContactWithFirms(contact ContactParams, firmIDs []int64) (int64, error) {
	// Begin transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Error("Failed to begin transaction: ", err)
		return 0, err
	}

	// Step 1: Insert the contact
	contactQuery := `
	INSERT INTO contacts (anrede, vorname, nachname, position, telefon, mobil, email, abteilung, geburtstag, bemerkung, kontotyp) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	contactResult, err := tx.Exec(contactQuery,
		contact.Anrede,
		contact.Vorname,
		contact.Nachname,
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

	// Step 2: Create relationships with all firms
	relationQuery := `
	INSERT INTO firms_contacts (firma_id, contact_id) 
	VALUES (?, ?)`

	for _, firmID := range firmIDs {
		_, err = tx.Exec(relationQuery, firmID, contactID)
		if err != nil {
			tx.Rollback()
			log.Error("Failed to create relationship with firm ID ", firmID, ": ", err)
			return 0, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		log.Error("Failed to commit transaction: ", err)
		return 0, err
	}

	log.WithFields(log.Fields{
		"contact_id": contactID,
		"firm_ids":   firmIDs,
	}).Info("Contact inserted and relationships created successfully")

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
	SELECT id, anrede, vorname, nachname, position, telefon, mobil, 
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
			&contact.ID, &contact.Anrede, &contact.Vorname, &contact.Nachname, &contact.Position,
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

func (db *MySQLDB) InsertUser(user model.User) (int64, error) {
	query := `
	INSERT INTO users (username, email, hashed_password, created_by, last_login)
	VALUES (?, ?, ?, ?, ?)`

	result, err := db.DB.Exec(query, user.Username, user.Email, user.HashedPassword, user.CreatedBy, user.LastLogin)
	if err != nil {
		log.Error("Failed to insert user: ", err)
		return 0, err
	}

	userID, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get user ID: ", err)
		return 0, err
	}

	log.Info("User inserted successfully with ID: ", userID)
	return userID, nil
}

// GetUserByID fetches a user by ID
func (db *MySQLDB) GetUserByID(id int64) (*model.User, error) {
	query := `SELECT id, username, email, hashed_password, created_at, created_by, last_login FROM users WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.CreatedBy, &user.LastLogin)
	if err != nil {
		log.Error("Failed to get user: ", err)
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user data
func (db *MySQLDB) UpdateUser(user model.User) error {
	query := `
	UPDATE users
	SET username = ?, email = ?, hashed_password = ?, created_by = ?, last_login = ?
	WHERE id = ?`

	_, err := db.DB.Exec(query, user.Username, user.Email, user.HashedPassword, user.CreatedBy, user.LastLogin, user.ID)
	if err != nil {
		log.Error("Failed to update user: ", err)
		return err
	}

	log.Info("User updated successfully")
	return nil
}

// DeleteUser deletes a user by ID
func (db *MySQLDB) DeleteUser(id int64) error {
	query := `DELETE FROM users WHERE id = ?`

	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Error("Failed to delete user: ", err)
		return err
	}

	log.Info("User deleted successfully")
	return nil
}

func (db *MySQLDB) InsertRole(role model.Role) (int64, error) {
	query := `
	INSERT INTO roles (name, description)
	VALUES (?, ?)`

	result, err := db.DB.Exec(query, role.Name, role.Description)
	if err != nil {
		log.Error("Failed to insert role: ", err)
		return 0, err
	}

	roleID, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get role ID: ", err)
		return 0, err
	}

	log.Info("Role inserted successfully with ID: ", roleID)
	return roleID, nil
}

// GetRoleByID fetches a role by ID
func (db *MySQLDB) GetRoleByID(id int64) (*model.Role, error) {
	query := `SELECT id, name, description FROM roles WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var role model.Role
	err := row.Scan(&role.ID, &role.Name, &role.Description)
	if err != nil {
		log.Error("Failed to get role: ", err)
		return nil, err
	}
	return &role, nil
}

// UpdateRole updates a role
func (db *MySQLDB) UpdateRole(role model.Role) error {
	query := `UPDATE roles SET name = ?, description = ? WHERE id = ?`

	_, err := db.DB.Exec(query, role.Name, role.Description, role.ID)
	if err != nil {
		log.Error("Failed to update role: ", err)
		return err
	}
	log.Info("Role updated successfully")
	return nil
}

// DeleteRole deletes a role by ID
func (db *MySQLDB) DeleteRole(id int64) error {
	query := `DELETE FROM roles WHERE id = ?`

	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Error("Failed to delete role: ", err)
		return err
	}
	log.Info("Role deleted successfully")
	return nil
}

func (db *MySQLDB) InsertPermission(permission model.Permission) (int64, error) {
	query := `
	INSERT INTO permissions (name, description)
	VALUES (?, ?)`

	result, err := db.DB.Exec(query, permission.Name, permission.Description)
	if err != nil {
		log.Error("Failed to insert permission: ", err)
		return 0, err
	}

	permID, err := result.LastInsertId()
	if err != nil {
		log.Error("Failed to get permission ID: ", err)
		return 0, err
	}

	log.Info("Permission inserted successfully with ID: ", permID)
	return permID, nil
}

// GetPermissionByID fetches a permission by ID
func (db *MySQLDB) GetPermissionByID(id int64) (*model.Permission, error) {
	query := `SELECT id, name, description FROM permissions WHERE id = ?`
	row := db.DB.QueryRow(query, id)

	var perm model.Permission
	err := row.Scan(&perm.ID, &perm.Name, &perm.Description)
	if err != nil {
		log.Error("Failed to get permission: ", err)
		return nil, err
	}
	return &perm, nil
}

// UpdatePermission updates a permission
func (db *MySQLDB) UpdatePermission(permission model.Permission) error {
	query := `UPDATE permissions SET name = ?, description = ? WHERE id = ?`

	_, err := db.DB.Exec(query, permission.Name, permission.Description, permission.ID)
	if err != nil {
		log.Error("Failed to update permission: ", err)
		return err
	}
	log.Info("Permission updated successfully")
	return nil
}

// DeletePermission deletes a permission by ID
func (db *MySQLDB) DeletePermission(id int64) error {
	query := `DELETE FROM permissions WHERE id = ?`

	_, err := db.DB.Exec(query, id)
	if err != nil {
		log.Error("Failed to delete permission: ", err)
		return err
	}
	log.Info("Permission deleted successfully")
	return nil
}

// GetUserRoles returns all roles assigned to a user
func (db *MySQLDB) GetUserRoles(userID int64) ([]model.Role, error) {
	query := `
	SELECT r.id, r.name, r.description
	FROM roles r
	JOIN user_roles ur ON ur.role_id = r.id
	WHERE ur.user_id = ?`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		log.Error("Failed to get user roles: ", err)
		return nil, err
	}
	defer rows.Close()

	var roles []model.Role
	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Description); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}
	return roles, nil
}

// DeleteUserRole deletes a user-role mapping
func (db *MySQLDB) DeleteUserRole(userID, roleID int64) error {
	query := `DELETE FROM user_roles WHERE user_id = ? AND role_id = ?`

	_, err := db.DB.Exec(query, userID, roleID)
	if err != nil {
		log.Error("Failed to delete user-role mapping: ", err)
		return err
	}
	log.Info("User-Role mapping deleted")
	return nil
}

// GetRolePermissions returns all permissions assigned to a role
func (db *MySQLDB) GetRolePermissions(roleID int64) ([]model.Permission, error) {
	query := `
	SELECT p.id, p.name, p.description
	FROM permissions p
	JOIN role_permissions rp ON rp.permission_id = p.id
	WHERE rp.role_id = ?`

	rows, err := db.DB.Query(query, roleID)
	if err != nil {
		log.Error("Failed to get role permissions: ", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var perm model.Permission
		if err := rows.Scan(&perm.ID, &perm.Name, &perm.Description); err != nil {
			return nil, err
		}
		permissions = append(permissions, perm)
	}
	return permissions, nil
}

// DeleteRolePermission deletes a role-permission mapping
func (db *MySQLDB) DeleteRolePermission(roleID, permissionID int64) error {
	query := `DELETE FROM role_permissions WHERE role_id = ? AND permission_id = ?`

	_, err := db.DB.Exec(query, roleID, permissionID)
	if err != nil {
		log.Error("Failed to delete role-permission mapping: ", err)
		return err
	}
	log.Info("Role-Permission mapping deleted")
	return nil
}

func (db *MySQLDB) InsertRolePermission(rp model.RolePermission) error {
	query := `
	INSERT INTO role_permissions (role_id, permission_id)
	VALUES (?, ?)`

	_, err := db.DB.Exec(query, rp.RoleID, rp.PermissionID)
	if err != nil {
		log.Error("Failed to insert role-permission mapping: ", err)
		return err
	}

	log.Info("Role-Permission mapping inserted successfully")
	return nil
}

func (db *MySQLDB) InsertUserRole(ur model.UserRole) error {
	query := `
	INSERT INTO user_roles (user_id, role_id)
	VALUES (?, ?)`

	_, err := db.DB.Exec(query, ur.UserID, ur.RoleID)
	if err != nil {
		log.Error("Failed to insert user-role mapping: ", err)
		return err
	}

	log.Info("User-Role mapping inserted successfully")
	return nil
}

func (db *MySQLDB) GetUserByEmail(email string) (*model.User, error) {
	query := `SELECT id, username, email, hashed_password, created_at, created_by, last_login FROM users WHERE email = ?`

	row := db.DB.QueryRow(query, email)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.CreatedBy, &user.LastLogin)
	if err != nil {
		log.Error("Failed to fetch user by email: ", err)
		return nil, err
	}
	return &user, nil
}

func (db *MySQLDB) GetUserByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, email, hashed_password, created_at, created_by, last_login FROM users WHERE username = ?`

	row := db.DB.QueryRow(query, username)
	var user model.User
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.HashedPassword, &user.CreatedAt, &user.CreatedBy, &user.LastLogin)
	if err != nil {
		log.Error("Failed to fetch user by username: ", err)
		return nil, err
	}
	return &user, nil
}

func (db *MySQLDB) GetUserPermissions(userID int64) ([]model.Permission, error) {
	query := `
		SELECT p.id, p.name, p.description
		FROM permissions p
		JOIN role_permissions rp ON rp.permission_id = p.id
		JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ?`

	rows, err := db.DB.Query(query, userID)
	if err != nil {
		log.Error("Failed to get user permissions: ", err)
		return nil, err
	}
	defer rows.Close()

	var permissions []model.Permission
	for rows.Next() {
		var p model.Permission
		err := rows.Scan(&p.ID, &p.Name, &p.Description)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, p)
	}
	return permissions, nil
}

func (db *MySQLDB) GetPermissionByName(name string) (*model.Permission, error) {
	query := `SELECT id, name, description FROM permissions WHERE name = ?`
	row := db.DB.QueryRow(query, name)

	var p model.Permission
	err := row.Scan(&p.ID, &p.Name, &p.Description)
	return &p, err
}

func (db *MySQLDB) GetRoleByName(name string) (*model.Role, error) {
	query := `SELECT id, name, description FROM roles WHERE name = ?`
	row := db.DB.QueryRow(query, name)

	var r model.Role
	err := row.Scan(&r.ID, &r.Name, &r.Description)
	return &r, err
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
		INSERT INTO contacts (anrede, vorname, nachname, position, telefon, mobil, email, abteilung, geburtstag, bemerkung, kontotyp)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	contacts := [][]interface{}{
		{"Herr", "Max", "Mustermann", "Vertriebsleiter", "+49 30 123457", "+49 170 1234567", "max.mustermann@musterfirma.de", "Vertrieb", "1985-04-12", "Hauptansprechpartner für Vertrieb", "K"},
		{"Frau", "Lisa", "Meier", "IT-Projektmanagerin", "+49 89 987655", "+49 172 9876543", "lisa.meier@techsolutions.de", "Projektmanagement", "1990-09-25", "Leitet Softwareprojekte", "L"},
		{"Herr", "Johann", "Schmidt", "Einkaufsleiter", "+49 221 456788", "", "johann.schmidt@baustoffe-schmidt.de", "Einkauf", "1978-06-30", "Verantwortlich für Materialbeschaffung", "L"},
		{"Herr", "Michael", "König", "Key Account Manager", "+49 30 7654321", "+49 172 7654321", "michael.koenig@musterfirma.de", "Vertrieb", "1982-11-05", "Betreut Großkunden", "K"},
		{"Frau", "Anna", "Weber", "Kundensupport", "+49 30 2345678", "+49 173 2345678", "anna.weber@musterfirma.de", "Support", "1995-07-20", "Support für Geschäftskunden", "K"},
		{"Frau", "Sarah", "Bauer", "Finanzleiterin", "+49 89 1234599", "+49 175 1234599", "sarah.bauer@finanzexperten.de", "Finanzen", "1983-03-15", "Buchhaltung und Finanzen", "L"},
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

func (db *MySQLDB) InsertUserRolesTestData() error {
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

	// Insert roles data
	roleQuery := `
		INSERT INTO roles (name, description)
		VALUES (?, ?);
	`
	roles := [][]interface{}{
		{"admin", "Full system access and control"},
		{"manager", "Department level control and reporting access"},
		{"user", "Standard user access to application"},
		{"guest", "Limited read-only access to public resources"},
		{"customer_support", "Access to customer facing tools and information"},
	}

	// Prepare the role insert statement
	roleStmt, err := tx.Prepare(roleQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare roles statement: %v", err)
	}
	defer roleStmt.Close()

	// Insert roles with prepared statement
	var roleIDs []int64
	for _, role := range roles {
		result, err := roleStmt.Exec(role...)
		if err != nil {
			return fmt.Errorf("error inserting role %v: %v", role[0], err)
		}
		roleID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID for role %v: %v", role[0], err)
		}
		roleIDs = append(roleIDs, roleID)
	}

	// Insert permissions data
	permQuery := `
		INSERT INTO permissions (name, description)
		VALUES (?, ?);
	`
	permissions := [][]interface{}{
		{"users_create", "Create new user accounts"},
		{"users_read", "View user account details"},
		{"users_update", "Modify user account data"},
		{"users_delete", "Delete user accounts"},
		{"roles_manage", "Create, edit, or delete user roles"},
		{"reports_view", "Access system reports"},
		{"reports_export", "Export system reports"},
		{"settings_edit", "Modify system configuration"},
		{"customers_view", "View customer records"},
		{"customers_edit", "Modify customer records"},
	}

	// Prepare the permission insert statement
	permStmt, err := tx.Prepare(permQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare permissions statement: %v", err)
	}
	defer permStmt.Close()

	// Insert permissions with prepared statement
	var permIDs []int64
	for _, perm := range permissions {
		result, err := permStmt.Exec(perm...)
		if err != nil {
			return fmt.Errorf("error inserting permission %v: %v", perm[0], err)
		}
		permID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID for permission %v: %v", perm[0], err)
		}
		permIDs = append(permIDs, permID)
	}

	// Insert users data
	// Note: In a real application, you would use proper password hashing
	// using packages like bcrypt, argon2, etc.
	hashedPassword := "password123" // hash for "password123"
	userQuery := `
		INSERT INTO users (username, email, hashed_password, created_at)
		VALUES (?, ?, ?, NOW());
	`
	users := [][]interface{}{
		{"admin", "admin@example.com", hashedPassword},
		{"jdoe", "john.doe@example.com", hashedPassword},
		{"asmith", "anna.smith@example.com", hashedPassword},
		{"mjones", "mark.jones@example.com", hashedPassword},
		{"lchen", "lily.chen@example.com", hashedPassword},
	}

	// Prepare the user insert statement
	userStmt, err := tx.Prepare(userQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare users statement: %v", err)
	}
	defer userStmt.Close()

	// Insert users with prepared statement
	var userIDs []int64
	for _, user := range users {
		result, err := userStmt.Exec(user...)
		if err != nil {
			return fmt.Errorf("error inserting user %v: %v", user[0], err)
		}
		userID, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("error getting last insert ID for user %v: %v", user[0], err)
		}
		userIDs = append(userIDs, userID)

		// Update the first user to have created_by set to itself (bootstrap admin)
		if len(userIDs) == 1 {
			_, err = tx.Exec("UPDATE users SET created_by = ? WHERE id = ?", userID, userID)
			if err != nil {
				return fmt.Errorf("error updating admin user's created_by: %v", err)
			}
		} else {
			// Set all other users to be created by the admin
			_, err = tx.Exec("UPDATE users SET created_by = ? WHERE id = ?", userIDs[0], userID)
			if err != nil {
				return fmt.Errorf("error updating user's created_by: %v", err)
			}
		}
	}

	// Insert role-permission relationships
	rolePermQuery := `
		INSERT INTO role_permissions (role_id, permission_id)
		VALUES (?, ?);
	`
	rolePermStmt, err := tx.Prepare(rolePermQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare role_permissions statement: %v", err)
	}
	defer rolePermStmt.Close()

	// Define role-permission mappings
	rolePermMappings := []struct {
		roleIndex int
		permIndex int
	}{
		// Admin has all permissions
		{0, 0}, {0, 1}, {0, 2}, {0, 3}, {0, 4}, {0, 5}, {0, 6}, {0, 7}, {0, 8}, {0, 9},
		// Manager has most permissions except delete users and manage roles
		{1, 0}, {1, 1}, {1, 2}, {1, 5}, {1, 6}, {1, 7}, {1, 8}, {1, 9},
		// Regular user has basic permissions
		{2, 1}, {2, 5}, {2, 8},
		// Guest has minimal permissions
		{3, 1}, {3, 8},
		// Customer support has customer-focused permissions
		{4, 1}, {4, 8}, {4, 9},
	}

	// Insert role-permission relationships
	for _, mapping := range rolePermMappings {
		_, err := rolePermStmt.Exec(roleIDs[mapping.roleIndex], permIDs[mapping.permIndex])
		if err != nil {
			return fmt.Errorf("error inserting role-permission mapping role_id=%d, perm_id=%d: %v",
				roleIDs[mapping.roleIndex], permIDs[mapping.permIndex], err)
		}
	}

	// Insert user-role relationships
	userRoleQuery := `
		INSERT INTO user_roles (user_id, role_id)
		VALUES (?, ?);
	`
	userRoleStmt, err := tx.Prepare(userRoleQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare user_roles statement: %v", err)
	}
	defer userRoleStmt.Close()

	// Define user-role mappings
	userRoleMappings := []struct {
		userIndex int
		roleIndex int
	}{
		{0, 0}, // admin has admin role
		{1, 1}, // John Doe has manager role
		{2, 2}, // Anna Smith has user role
		{3, 2}, // Mark Jones has user role
		{3, 4}, // Mark Jones also has customer_support role
		{4, 3}, // Lily Chen has guest role
	}

	// Insert user-role relationships
	for _, mapping := range userRoleMappings {
		_, err := userRoleStmt.Exec(userIDs[mapping.userIndex], roleIDs[mapping.roleIndex])
		if err != nil {
			return fmt.Errorf("error inserting user-role mapping user_id=%d, role_id=%d: %v",
				userIDs[mapping.userIndex], roleIDs[mapping.roleIndex], err)
		}
	}

	log.Info("User, role, and permission test data successfully inserted into database.")
	return nil
}

func (db *MySQLDB) RunCRUDTests() error {
	log.Info("---- Running Full CRUD Tests ----")

	// Insert a test user
	testUser := model.User{
		Username:       "testuser1",
		Email:          "testuser1@example.com",
		HashedPassword: "hashed_pass_123",
	}
	userID, err := db.InsertUser(testUser)
	if err != nil {
		return err
	}
	log.Infof("Inserted User ID: %d", userID)

	// Insert a test role
	testRole := model.Role{
		Name:        "admin",
		Description: "Administrator role",
	}
	roleID, err := db.InsertRole(testRole)
	if err != nil {
		return err
	}
	log.Infof("Inserted Role ID: %d", roleID)

	// Insert a test permission
	testPerm := model.Permission{
		Name:        "read_all",
		Description: "Permission to read everything",
	}
	permID, err := db.InsertPermission(testPerm)
	if err != nil {
		return err
	}
	log.Infof("Inserted Permission ID: %d", permID)

	// Assign role to user
	err = db.InsertUserRole(model.UserRole{UserID: userID, RoleID: roleID})
	if err != nil {
		return err
	}
	log.Infof("Assigned Role ID %d to User ID %d", roleID, userID)

	// Assign permission to role
	err = db.InsertRolePermission(model.RolePermission{RoleID: roleID, PermissionID: permID})
	if err != nil {
		return err
	}
	log.Infof("Assigned Permission ID %d to Role ID %d", permID, roleID)

	// Read user and their roles
	user, err := db.GetUserByID(userID)
	if err != nil {
		return err
	}
	log.Infof("Fetched User: %+v", user)

	roles, err := db.GetUserRoles(userID)
	if err != nil {
		return err
	}
	log.Infof("User Roles: %+v", roles)

	// Read role and its permissions
	perms, err := db.GetRolePermissions(roleID)
	if err != nil {
		return err
	}
	log.Infof("Role Permissions: %+v", perms)

	// Update the user
	user.Username = "updateduser"
	user.Email = "updateduser@example.com"
	err = db.UpdateUser(*user)
	if err != nil {
		return err
	}
	log.Info("Updated User")

	// Update role
	testRole.ID = roleID
	testRole.Description = "Updated admin role"
	err = db.UpdateRole(testRole)
	if err != nil {
		return err
	}
	log.Info("Updated Role")

	// Update permission
	testPerm.ID = permID
	testPerm.Description = "Updated permission to read everything"
	err = db.UpdatePermission(testPerm)
	if err != nil {
		return err
	}
	log.Info("Updated Permission")

	// Cleanup: Delete mappings first
	err = db.DeleteUserRole(userID, roleID)
	if err != nil {
		return err
	}
	log.Infof("Deleted user-role mapping")

	err = db.DeleteRolePermission(roleID, permID)
	if err != nil {
		return err
	}
	log.Infof("Deleted role-permission mapping")

	// Delete entities
	err = db.DeletePermission(permID)
	if err != nil {
		return err
	}
	log.Infof("Deleted Permission")

	err = db.DeleteRole(roleID)
	if err != nil {
		return err
	}
	log.Infof("Deleted Role")

	err = db.DeleteUser(userID)
	if err != nil {
		return err
	}
	log.Infof("Deleted User")

	log.Info("---- CRUD Tests Completed Successfully ✅ ----")
	return nil
}
