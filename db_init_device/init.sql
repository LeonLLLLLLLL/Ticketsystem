-- Create the main devices table (✅ all JSONB → TEXT)
CREATE TABLE devices (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    hostname VARCHAR(255),
    ip VARCHAR(50),
    domain VARCHAR(255),
    manufacturer VARCHAR(255),
    model_type VARCHAR(255),
    serial_numbers TEXT,                -- was JSONB
    mac VARCHAR(50),
    description TEXT,
    equipment TEXT,
    function TEXT,
    settings TEXT,
    device_link TEXT,
    commissioning_date DATE,
    origin TEXT,
    warranty_service_number TEXT,
    warranty_until DATE,
    licenses TEXT,                      -- was JSONB
    location_text TEXT,
    department TEXT,
    internal_contact TEXT,
    external_contact TEXT,
    map_link TEXT,
    software_interfaces TEXT,           -- was JSONB
    backup_method TEXT,
    backup_file_link TEXT,
    software_asset TEXT,
    password_link TEXT,
    internal_access TEXT,
    external_access TEXT,
    misc_links TEXT,                    -- was JSONB
    externally_accessible BOOLEAN,
    restart_how TEXT,
    restart_notes TEXT,
    restart_coordination TEXT,
    network_connection TEXT,
    patch_location TEXT,
    documents TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create many-to-many relationship table for linking device assets
CREATE TABLE device_links (
    id SERIAL PRIMARY KEY,
    from_device_id INT NOT NULL,
    to_device_id INT NOT NULL,
    CONSTRAINT fk_from_device FOREIGN KEY (from_device_id) REFERENCES devices (id) ON DELETE CASCADE,
    CONSTRAINT fk_to_device FOREIGN KEY (to_device_id) REFERENCES devices (id) ON DELETE CASCADE
);

-- Insert sample device 1: Main Router
INSERT INTO devices (
    name, hostname, ip, domain, manufacturer, model_type,
    serial_numbers, mac, description, equipment, function, settings,
    device_link, commissioning_date, origin, warranty_service_number,
    warranty_until, licenses, location_text, department, internal_contact,
    external_contact, map_link, software_interfaces, backup_method,
    backup_file_link, software_asset, password_link, internal_access,
    external_access, misc_links, externally_accessible, restart_how,
    restart_notes, restart_coordination, network_connection, patch_location,
    documents
)
VALUES (
    'Main Router', 'router01', '192.168.0.1', 'corp.local', 'Cisco', 'RV340',
    '[{"element": "Board", "sn": "SN-MAIN-001"}]',
    'AA:BB:CC:DD:EE:01', 'Main internet router for corporate network',
    'Rack mount, power cable', 'Routing & Firewall', 'Default Config',
    'http://example.com/router01', '2022-05-15', 'Warehouse',
    'WSN-0001', '2025-05-15',
    '[{"element": "Firmware", "key": "LIC-12345"}]',
    'Server Room A', 'IT', 'Alice Admin', 'Cisco Support',
    'http://map.local/router01',
    '[{"name": "API v1"}, {"name": "SNMP"}]',
    'Cloud backup', 'http://backup.local/router01.zip', 'IOS 1.0.3',
    'http://secrets.local/router01', 'SSH', 'VPN',
    '[{"label": "Manual", "url": "http://docs.local/router01.pdf"}]',
    true, 'SSH reboot command', 'Inform network team', 'After business hours',
    'LAN', 'Patch Panel 3A', 'http://docs.local/router01_docs.pdf'
);

-- Insert sample device 2: Backup Router
INSERT INTO devices (
    name, hostname, ip, domain, manufacturer, model_type,
    serial_numbers, mac, description, equipment, function, settings,
    device_link, commissioning_date, origin, warranty_service_number,
    warranty_until, licenses, location_text, department, internal_contact,
    external_contact, map_link, software_interfaces, backup_method,
    backup_file_link, software_asset, password_link, internal_access,
    external_access, misc_links, externally_accessible, restart_how,
    restart_notes, restart_coordination, network_connection, patch_location,
    documents
)
VALUES (
    'Backup Router', 'router02', '192.168.0.2', 'corp.local', 'Cisco', 'RV340',
    '[{"element": "Board", "sn": "SN-BACK-002"}]',
    'AA:BB:CC:DD:EE:02', 'Backup router for redundancy',
    'Rack mount, power cable', 'Failover Routing', 'Default Backup Config',
    'http://example.com/router02', '2022-06-01', 'Warehouse',
    'WSN-0002', '2025-06-01',
    '[{"element": "Firmware", "key": "LIC-67890"}]',
    'Server Room B', 'IT', 'Bob Backup', 'Cisco Support',
    'http://map.local/router02',
    '[{"name": "API v1"}, {"name": "SNMP"}]',
    'Cloud backup', 'http://backup.local/router02.zip', 'IOS 1.0.3',
    'http://secrets.local/router02', 'SSH', 'VPN',
    '[{"label": "Manual", "url": "http://docs.local/router02.pdf"}]',
    true, 'SSH reboot command', 'Inform IT manager', 'Only after failover',
    'LAN', 'Patch Panel 3B', 'http://docs.local/router02_docs.pdf'
);

-- Link Main Router to Backup Router and vice versa
INSERT INTO device_links (from_device_id, to_device_id)
VALUES 
    (1, 2),
    (2, 1);
