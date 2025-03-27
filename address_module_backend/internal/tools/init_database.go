package tools

import (
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

// SetupFirmsTable creates the firms table
func (db *MySQLDB) SetupFirmsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS firms (
		id INT AUTO_INCREMENT PRIMARY KEY,
		anrede VARCHAR(50) NOT NULL,
		name_1 VARCHAR(255) NOT NULL,
		name_2 VARCHAR(255),
		name_3 VARCHAR(255),
		straße VARCHAR(255),
		land VARCHAR(100),
		plz VARCHAR(20) NOT NULL,
		ort VARCHAR(255) NOT NULL,
		telefon VARCHAR(50) NOT NULL,
		email VARCHAR(255) NOT NULL,
		website VARCHAR(255),
		kunde BOOLEAN DEFAULT FALSE,
		lieferant BOOLEAN DEFAULT FALSE,
		gesperrt BOOLEAN DEFAULT FALSE,
		bemerkung TEXT,
		firma_typ VARCHAR(100),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create firms table: ", err)
		return err
	}
	log.Info("Firms table setup completed")
	return nil
}

// SetupContactsTable creates the contacts table
func (db *MySQLDB) SetupContactsTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS contacts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		anrede VARCHAR(50),
		vorname VARCHAR(255) NOT NULL,
		nachname VARCHAR(255) NOT NULL,
		position VARCHAR(255),
		telefon VARCHAR(50),
		mobil VARCHAR(50),
		email VARCHAR(255) NOT NULL,
		abteilung VARCHAR(255),
		geburtstag DATE,
		bemerkung TEXT,
		kontotyp VARCHAR(50) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE KEY idx_contact_email (email)
	);`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create contacts table: ", err)
		return err
	}
	log.Info("Contacts table setup completed")
	return nil
}

// SetupFirmsContactsRelationTable creates the junction table
func (db *MySQLDB) SetupFirmsContactsRelationTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS firms_contacts (
		id INT AUTO_INCREMENT PRIMARY KEY,
		firma_id INT NOT NULL,
		contact_id INT NOT NULL,
		beziehung VARCHAR(255),
		hauptansprechpartner BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (firma_id) REFERENCES firms(id) ON DELETE CASCADE,
		FOREIGN KEY (contact_id) REFERENCES contacts(id) ON DELETE CASCADE,
		UNIQUE KEY unique_firm_contact (firma_id, contact_id)
	);`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create firms_contacts junction table: ", err)
		return err
	}
	log.Info("Firms-Contacts junction table setup completed")
	return nil
}

// SetupPerformanceIndexes creates additional indexes
func (db *MySQLDB) SetupPerformanceIndexes() error {
	queries := []string{
		"CREATE INDEX IF NOT EXISTS idx_firms_name ON firms(name_1);",
		"CREATE INDEX IF NOT EXISTS idx_contacts_name ON contacts(name);",
		"CREATE INDEX IF NOT EXISTS idx_firms_contacts_firma ON firms_contacts(firma_id);",
		"CREATE INDEX IF NOT EXISTS idx_firms_contacts_contact ON firms_contacts(contact_id);",
	}

	for _, query := range queries {
		_, err := db.DB.Exec(query)
		if err != nil {
			log.Error("Failed to create index: ", err)
			return err
		}
	}
	log.Info("Performance indexes setup completed")
	return nil
}

// SetupUsersTable creates the users table
func (db *MySQLDB) SetupUsersTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(50) NOT NULL,
        email VARCHAR(100) NOT NULL UNIQUE,
        hashed_password VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        created_by INT,
        last_login TIMESTAMP NULL,
        FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE SET NULL
    );`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create users table: ", err)
		return err
	}
	log.Info("Users table setup completed")
	return nil
}

// SetupRolesTable creates the roles table
func (db *MySQLDB) SetupRolesTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS roles (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(50) NOT NULL UNIQUE,
        description VARCHAR(255)
    );`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create roles table: ", err)
		return err
	}
	log.Info("Roles table setup completed")
	return nil
}

// SetupPermissionsTable creates the permissions table
func (db *MySQLDB) SetupPermissionsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS permissions (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(50) NOT NULL UNIQUE,
        description VARCHAR(255)
    );`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create permissions table: ", err)
		return err
	}
	log.Info("Permissions table setup completed")
	return nil
}

// SetupRolePermissionsTable creates the junction table for roles and permissions
func (db *MySQLDB) SetupRolePermissionsTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS role_permissions (
        role_id INT NOT NULL,
        permission_id INT NOT NULL,
        PRIMARY KEY (role_id, permission_id),
        FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
        FOREIGN KEY (permission_id) REFERENCES permissions(id) ON DELETE CASCADE
    );`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create role_permissions junction table: ", err)
		return err
	}
	log.Info("Role-Permissions junction table setup completed")
	return nil
}

// SetupUserRolesTable creates the junction table for users and roles
func (db *MySQLDB) SetupUserRolesTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS user_roles (
        user_id INT NOT NULL,
        role_id INT NOT NULL,
        PRIMARY KEY (user_id, role_id),
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE
    );`

	_, err := db.DB.Exec(query)
	if err != nil {
		log.Error("Failed to create user_roles junction table: ", err)
		return err
	}
	log.Info("User-Roles junction table setup completed")
	return nil
}

// SetupDatabase sets up all tables and indexes
func (db *MySQLDB) SetupDatabase() error {
	// Order matters due to foreign key constraints
	setupFuncs := []func() error{
		db.SetupFirmsTable,
		db.SetupContactsTable,
		db.SetupFirmsContactsRelationTable,
		db.SetupUsersTable,
		db.SetupRolesTable,
		db.SetupPermissionsTable,
		db.SetupRolePermissionsTable,
		db.SetupUserRolesTable,
		//db.SetupPerformanceIndexes,
	}

	for _, setupFunc := range setupFuncs {
		err := setupFunc()
		if err != nil {
			return err
		}
	}

	log.Info("Complete database setup finished successfully")
	return nil
}
