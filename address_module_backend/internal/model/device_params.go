package model

type DeviceParams struct {
	ID                    int64
	Name                  string
	Hostname              string
	IP                    string
	Domain                string
	Manufacturer          string
	ModelType             string
	SerialNumbers         string // store as JSON string
	MAC                   string
	Description           string
	Equipment             string
	Function              string
	Settings              string
	DeviceLink            string
	CommissioningDate     string // YYYY-MM-DD
	Origin                string
	WarrantyServiceNumber string
	WarrantyUntil         string // YYYY-MM-DD
	Licenses              string // store as JSON string
	LocationText          string
	Department            string
	InternalContact       string
	ExternalContact       string
	MapLink               string
	SoftwareInterfaces    string // JSON string
	BackupMethod          string
	BackupFileLink        string
	SoftwareAsset         string
	PasswordLink          string
	InternalAccess        string
	ExternalAccess        string
	MiscLinks             string // JSON string
	ExternallyAccessible  bool
	RestartHow            string
	RestartNotes          string
	RestartCoordination   string
	NetworkConnection     string
	PatchLocation         string
	Documents             string
}
