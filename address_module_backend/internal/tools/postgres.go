package tools

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type PostgresDB struct {
	DB *sql.DB
}

func NewPostgresDatabase(maxRetries int, delay time.Duration) (*PostgresDB, error) {
	host := os.Getenv("DEVICE_DB_HOST")
	user := os.Getenv("DEVICE_DB_USER")
	password := os.Getenv("DEVICE_DB_PASSWORD")
	dbName := os.Getenv("DEVICE_DB_NAME")
	port := os.Getenv("DEVICE_DB_PORT")

	if dbName == "" {
		dbName = "device_management_database"
	}
	if port == "" {
		port = "5432"
	}

	log.Infof("🚧 DEVICE_DB_HOST: %s", host)
	log.Infof("👤 DEVICE_DB_USER: %s", user)
	log.Infof("🗃️  DEVICE_DB_NAME: %s", dbName)
	log.Infof("📦 DEVICE_DB_PORT: %s", port)

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			log.Errorf("Postgres open error (attempt %d): %v", attempt, err)
			time.Sleep(delay)
			continue
		}

		if err := db.Ping(); err != nil {
			log.Warnf("Postgres ping failed (attempt %d): %v", attempt, err)
			db.Close()
			time.Sleep(delay)
			continue
		}

		log.Infof("✅ Connected to PostgreSQL DB on attempt %d", attempt)
		return &PostgresDB{DB: db}, nil
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts", maxRetries)
}

func (p *PostgresDB) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// Simplified: JSONB fields now stored as TEXT
type DeviceParams struct {
	ID                    int64      `json:"id"`
	Name                  string     `json:"name"`
	Hostname              string     `json:"hostname"`
	IP                    string     `json:"ip"`
	Domain                string     `json:"domain"`
	Manufacturer          string     `json:"manufacturer"`
	ModelType             string     `json:"model_type"`
	SerialNumbers         string     `json:"serial_numbers"`
	MAC                   string     `json:"mac"`
	Description           string     `json:"description"`
	Equipment             string     `json:"equipment"`
	Function              string     `json:"function"`
	Settings              string     `json:"settings"`
	DeviceLink            string     `json:"device_link"`
	CommissioningDate     *time.Time `json:"commissioning_date"`
	Origin                string     `json:"origin"`
	WarrantyServiceNumber string     `json:"warranty_service_number"`
	WarrantyUntil         *time.Time `json:"warranty_until"`
	Licenses              string     `json:"licenses"`
	LocationText          string     `json:"location_text"`
	Department            string     `json:"department"`
	InternalContact       string     `json:"internal_contact"`
	ExternalContact       string     `json:"external_contact"`
	MapLink               string     `json:"map_link"`
	SoftwareInterfaces    string     `json:"software_interfaces"`
	BackupMethod          string     `json:"backup_method"`
	BackupFileLink        string     `json:"backup_file_link"`
	SoftwareAsset         string     `json:"software_asset"`
	PasswordLink          string     `json:"password_link"`
	InternalAccess        string     `json:"internal_access"`
	ExternalAccess        string     `json:"external_access"`
	MiscLinks             string     `json:"misc_links"`
	ExternallyAccessible  bool       `json:"externally_accessible"`
	RestartHow            string     `json:"restart_how"`
	RestartNotes          string     `json:"restart_notes"`
	RestartCoordination   string     `json:"restart_coordination"`
	NetworkConnection     string     `json:"network_connection"`
	PatchLocation         string     `json:"patch_location"`
	Documents             string     `json:"documents"`
	CreatedAt             *time.Time `json:"created_at,omitempty"`
}

func (p *PostgresDB) GetAllDevices() ([]DeviceParams, error) {
	rows, err := p.DB.Query(`SELECT * FROM devices ORDER BY id DESC;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var devices []DeviceParams
	for rows.Next() {
		var d DeviceParams
		err := rows.Scan(
			&d.ID, &d.Name, &d.Hostname, &d.IP, &d.Domain, &d.Manufacturer, &d.ModelType,
			&d.SerialNumbers, &d.MAC, &d.Description, &d.Equipment, &d.Function, &d.Settings,
			&d.DeviceLink, &d.CommissioningDate, &d.Origin, &d.WarrantyServiceNumber,
			&d.WarrantyUntil, &d.Licenses, &d.LocationText, &d.Department, &d.InternalContact,
			&d.ExternalContact, &d.MapLink, &d.SoftwareInterfaces, &d.BackupMethod,
			&d.BackupFileLink, &d.SoftwareAsset, &d.PasswordLink, &d.InternalAccess,
			&d.ExternalAccess, &d.MiscLinks, &d.ExternallyAccessible, &d.RestartHow,
			&d.RestartNotes, &d.RestartCoordination, &d.NetworkConnection, &d.PatchLocation,
			&d.Documents, &d.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		devices = append(devices, d)
	}

	return devices, nil
}

func (p *PostgresDB) GetDeviceByID(id int64) (*DeviceParams, error) {
	var d DeviceParams
	err := p.DB.QueryRow(`SELECT * FROM devices WHERE id = $1`, id).Scan(
		&d.ID, &d.Name, &d.Hostname, &d.IP, &d.Domain, &d.Manufacturer, &d.ModelType,
		&d.SerialNumbers, &d.MAC, &d.Description, &d.Equipment, &d.Function, &d.Settings,
		&d.DeviceLink, &d.CommissioningDate, &d.Origin, &d.WarrantyServiceNumber,
		&d.WarrantyUntil, &d.Licenses, &d.LocationText, &d.Department, &d.InternalContact,
		&d.ExternalContact, &d.MapLink, &d.SoftwareInterfaces, &d.BackupMethod,
		&d.BackupFileLink, &d.SoftwareAsset, &d.PasswordLink, &d.InternalAccess,
		&d.ExternalAccess, &d.MiscLinks, &d.ExternallyAccessible, &d.RestartHow,
		&d.RestartNotes, &d.RestartCoordination, &d.NetworkConnection, &d.PatchLocation,
		&d.Documents, &d.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (p *PostgresDB) InsertDevice(device *DeviceParams) (int64, error) {
	query := `
	INSERT INTO devices (
		name, hostname, ip, domain, manufacturer, model_type, serial_numbers,
		mac, description, equipment, function, settings, device_link,
		commissioning_date, origin, warranty_service_number, warranty_until,
		licenses, location_text, department, internal_contact, external_contact,
		map_link, software_interfaces, backup_method, backup_file_link,
		software_asset, password_link, internal_access, external_access,
		misc_links, externally_accessible, restart_how, restart_notes,
		restart_coordination, network_connection, patch_location, documents
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7,
		$8, $9, $10, $11, $12, $13,
		$14, $15, $16, $17,
		$18, $19, $20, $21, $22,
		$23, $24, $25, $26,
		$27, $28, $29, $30,
		$31, $32, $33, $34,
		$35, $36, $37, $38
	) RETURNING id;
	`

	var id int64
	err := p.DB.QueryRow(query,
		device.Name, device.Hostname, device.IP, device.Domain, device.Manufacturer, device.ModelType, device.SerialNumbers,
		device.MAC, device.Description, device.Equipment, device.Function, device.Settings, device.DeviceLink,
		device.CommissioningDate, device.Origin, device.WarrantyServiceNumber, device.WarrantyUntil,
		device.Licenses, device.LocationText, device.Department, device.InternalContact, device.ExternalContact,
		device.MapLink, device.SoftwareInterfaces, device.BackupMethod, device.BackupFileLink,
		device.SoftwareAsset, device.PasswordLink, device.InternalAccess, device.ExternalAccess,
		device.MiscLinks, device.ExternallyAccessible, device.RestartHow, device.RestartNotes,
		device.RestartCoordination, device.NetworkConnection, device.PatchLocation, device.Documents,
	).Scan(&id)
	return id, err
}

func (p *PostgresDB) UpdateDevice(device *DeviceParams) error {
	query := `
	UPDATE devices SET
		name = $1, hostname = $2, ip = $3, domain = $4, manufacturer = $5, model_type = $6, serial_numbers = $7,
		mac = $8, description = $9, equipment = $10, function = $11, settings = $12, device_link = $13,
		commissioning_date = $14, origin = $15, warranty_service_number = $16, warranty_until = $17,
		licenses = $18, location_text = $19, department = $20, internal_contact = $21, external_contact = $22,
		map_link = $23, software_interfaces = $24, backup_method = $25, backup_file_link = $26,
		software_asset = $27, password_link = $28, internal_access = $29, external_access = $30,
		misc_links = $31, externally_accessible = $32, restart_how = $33, restart_notes = $34,
		restart_coordination = $35, network_connection = $36, patch_location = $37, documents = $38
	WHERE id = $39;
	`

	_, err := p.DB.Exec(query,
		device.Name, device.Hostname, device.IP, device.Domain, device.Manufacturer, device.ModelType, device.SerialNumbers,
		device.MAC, device.Description, device.Equipment, device.Function, device.Settings, device.DeviceLink,
		device.CommissioningDate, device.Origin, device.WarrantyServiceNumber, device.WarrantyUntil,
		device.Licenses, device.LocationText, device.Department, device.InternalContact, device.ExternalContact,
		device.MapLink, device.SoftwareInterfaces, device.BackupMethod, device.BackupFileLink,
		device.SoftwareAsset, device.PasswordLink, device.InternalAccess, device.ExternalAccess,
		device.MiscLinks, device.ExternallyAccessible, device.RestartHow, device.RestartNotes,
		device.RestartCoordination, device.NetworkConnection, device.PatchLocation, device.Documents,
		device.ID,
	)
	return err
}

func (p *PostgresDB) DeleteDevice(id int64) error {
	_, err := p.DB.Exec(`DELETE FROM devices WHERE id = $1`, id)
	return err
}

// ---- Device Links ----

type DeviceLink struct {
	ID           int64 `json:"id"`
	FromDeviceID int64 `json:"from_device_id"`
	ToDeviceID   int64 `json:"to_device_id"`
}

func (p *PostgresDB) GetAllDeviceLinks() ([]DeviceLink, error) {
	rows, err := p.DB.Query(`SELECT id, from_device_id, to_device_id FROM device_links ORDER BY id;`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []DeviceLink
	for rows.Next() {
		var l DeviceLink
		if err := rows.Scan(&l.ID, &l.FromDeviceID, &l.ToDeviceID); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

func (p *PostgresDB) InsertDeviceLink(link *DeviceLink) (int64, error) {
	var id int64
	err := p.DB.QueryRow(`INSERT INTO device_links (from_device_id, to_device_id) VALUES ($1, $2) RETURNING id;`,
		link.FromDeviceID, link.ToDeviceID).Scan(&id)
	return id, err
}

func (p *PostgresDB) GetLinksForDevice(deviceID int64) ([]DeviceLink, error) {
	rows, err := p.DB.Query(`SELECT id, from_device_id, to_device_id FROM device_links WHERE from_device_id = $1;`, deviceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []DeviceLink
	for rows.Next() {
		var l DeviceLink
		if err := rows.Scan(&l.ID, &l.FromDeviceID, &l.ToDeviceID); err != nil {
			return nil, err
		}
		links = append(links, l)
	}
	return links, nil
}

func (p *PostgresDB) UpdateDeviceLink(link *DeviceLink) error {
	query := `
		UPDATE device_links
		SET from_device_id = $1, to_device_id = $2
		WHERE id = $3;
	`
	_, err := p.DB.Exec(query, link.FromDeviceID, link.ToDeviceID, link.ID)
	return err
}

func (p *PostgresDB) DeleteDeviceLink(id int64) error {
	_, err := p.DB.Exec(`DELETE FROM device_links WHERE id = $1;`, id)
	return err
}

func RunPostgresDeviceCRUDTests(pgDB *PostgresDB) error {
	fmt.Println("🚀 Running Postgres Device CRUD Test...")

	now := time.Now()

	newDevice := &DeviceParams{
		Name:                  "Test Router",
		Hostname:              "testrouter01",
		IP:                    "192.168.100.1",
		Domain:                "test.local",
		Manufacturer:          "NetGear",
		ModelType:             "R7000",
		SerialNumbers:         `[{"element": "Board", "sn": "SN123456"}]`,
		MAC:                   "DE:AD:BE:EF:00:01",
		Description:           "Temporary test router for validation",
		Equipment:             "Standard rack mount",
		Function:              "Routing",
		Settings:              "Default configuration",
		DeviceLink:            "http://test.local/router01",
		CommissioningDate:     &now,
		Origin:                "Lab",
		WarrantyServiceNumber: "WSN-TEST-001",
		WarrantyUntil:         &now,
		Licenses:              `[{"element": "Firmware", "key": "LIC-TEST-001"}]`,
		LocationText:          "Test Lab",
		Department:            "QA",
		InternalContact:       "Test Admin",
		ExternalContact:       "NetGear Support",
		MapLink:               "http://maps.local/test01",
		SoftwareInterfaces:    `[{"name": "API"}, {"name": "SNMP"}]`,
		BackupMethod:          "Manual",
		BackupFileLink:        "http://backup.local/test.zip",
		SoftwareAsset:         "v1.0-test",
		PasswordLink:          "http://secrets.local/test",
		InternalAccess:        "SSH",
		ExternalAccess:        "VPN",
		MiscLinks:             `[{"label": "Manual", "url": "http://docs.local/testrouter01.pdf"}]`,
		ExternallyAccessible:  true,
		RestartHow:            "Power cycle",
		RestartNotes:          "Coordinate with lab",
		RestartCoordination:   "QA lead",
		NetworkConnection:     "Test VLAN",
		PatchLocation:         "Patch Panel X",
		Documents:             "http://docs.local/testrouter01.pdf",
	}

	// CREATE
	id, err := pgDB.InsertDevice(newDevice)
	if err != nil {
		return fmt.Errorf("❌ InsertDevice failed: %v", err)
	}
	fmt.Printf("✅ Created device with ID: %d\n", id)

	// READ
	device, err := pgDB.GetDeviceByID(id)
	if err != nil {
		return fmt.Errorf("❌ GetDeviceByID failed: %v", err)
	}
	fmt.Printf("📦 Retrieved device: %s [%s]\n", device.Name, device.IP)

	// UPDATE
	device.Description = "UPDATED: Test device after refactor"
	if err := pgDB.UpdateDevice(device); err != nil {
		return fmt.Errorf("❌ UpdateDevice failed: %v", err)
	}
	fmt.Println("🔄 Updated device description successfully.")

	// LIST
	devices, err := pgDB.GetAllDevices()
	if err != nil {
		return fmt.Errorf("❌ GetAllDevices failed: %v", err)
	}
	fmt.Printf("📋 Total devices in DB: %d\n", len(devices))

	// CREATE LINK
	link := &DeviceLink{
		FromDeviceID: id,
		ToDeviceID:   1, // Ensure device ID 1 exists
	}
	linkID, err := pgDB.InsertDeviceLink(link)
	if err != nil {
		return fmt.Errorf("❌ InsertDeviceLink failed: %v", err)
	}
	fmt.Printf("🔗 Created device link with ID: %d\n", linkID)

	// READ LINK
	links, err := pgDB.GetLinksForDevice(id)
	if err != nil {
		return fmt.Errorf("❌ GetLinksForDevice failed: %v", err)
	}
	fmt.Printf("🔎 Device %d has %d link(s)\n", id, len(links))

	// DELETE LINK
	if err := pgDB.DeleteDeviceLink(linkID); err != nil {
		return fmt.Errorf("❌ DeleteDeviceLink failed: %v", err)
	}
	fmt.Printf("🧹 Deleted device link ID %d\n", linkID)

	// DELETE DEVICE
	if err := pgDB.DeleteDevice(id); err != nil {
		return fmt.Errorf("❌ DeleteDevice failed: %v", err)
	}
	fmt.Printf("🗑️ Deleted device with ID: %d\n", id)

	fmt.Println("✅ CRUD test completed successfully.")
	return nil
}
