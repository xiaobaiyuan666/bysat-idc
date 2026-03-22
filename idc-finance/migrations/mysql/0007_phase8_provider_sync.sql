ALTER TABLE services
  ADD COLUMN provider_resource_id VARCHAR(128) NULL AFTER provider_type,
  ADD COLUMN region_name VARCHAR(128) NOT NULL DEFAULT '' AFTER provider_resource_id,
  ADD COLUMN ip_address VARCHAR(128) NOT NULL DEFAULT '' AFTER region_name,
  ADD COLUMN sync_status VARCHAR(32) NOT NULL DEFAULT 'IDLE' AFTER ip_address,
  ADD COLUMN sync_message VARCHAR(255) NOT NULL DEFAULT '' AFTER sync_status,
  ADD COLUMN last_sync_at DATETIME NULL AFTER sync_message,
  ADD UNIQUE KEY uk_services_provider_resource (provider_type, provider_resource_id),
  ADD INDEX idx_services_sync_status (sync_status, last_sync_at);

CREATE TABLE IF NOT EXISTS provider_sync_logs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  provider_type VARCHAR(64) NOT NULL,
  action VARCHAR(64) NOT NULL,
  resource_type VARCHAR(64) NOT NULL DEFAULT 'instance',
  resource_id VARCHAR(128) NOT NULL DEFAULT '',
  service_id BIGINT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'SUCCESS',
  message VARCHAR(255) NOT NULL DEFAULT '',
  request_body JSON NULL,
  response_body JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX idx_provider_sync_logs_resource (provider_type, resource_type, resource_id, created_at),
  INDEX idx_provider_sync_logs_service (service_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_network_interfaces (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  provider_interface_id VARCHAR(128) NULL,
  name VARCHAR(128) NOT NULL DEFAULT '',
  mac_address VARCHAR(64) NOT NULL DEFAULT '',
  bridge VARCHAR(128) NOT NULL DEFAULT '',
  network_name VARCHAR(128) NOT NULL DEFAULT '',
  interface_name VARCHAR(128) NOT NULL DEFAULT '',
  nic_model VARCHAR(64) NOT NULL DEFAULT '',
  inbound_mbps INT NOT NULL DEFAULT 0,
  outbound_mbps INT NOT NULL DEFAULT 0,
  vpc_name VARCHAR(128) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_network_interfaces_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_ip_addresses (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  network_interface_id BIGINT NULL,
  provider_ip_id VARCHAR(128) NULL,
  address VARCHAR(128) NOT NULL,
  version VARCHAR(16) NOT NULL DEFAULT 'IPv4',
  ip_role VARCHAR(32) NOT NULL DEFAULT 'public',
  subnet_mask VARCHAR(64) NOT NULL DEFAULT '',
  gateway VARCHAR(128) NOT NULL DEFAULT '',
  bandwidth_mbps INT NOT NULL DEFAULT 0,
  is_primary TINYINT(1) NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_ip_addresses_service_id (service_id),
  INDEX idx_service_ip_addresses_address (address)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_disks (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  provider_disk_id VARCHAR(128) NULL,
  name VARCHAR(128) NOT NULL,
  disk_type VARCHAR(32) NOT NULL DEFAULT 'data',
  size_gb INT NOT NULL DEFAULT 0,
  device_name VARCHAR(64) NOT NULL DEFAULT '',
  driver_name VARCHAR(64) NOT NULL DEFAULT '',
  fs_type VARCHAR(64) NOT NULL DEFAULT '',
  mount_point VARCHAR(128) NOT NULL DEFAULT '',
  is_system TINYINT(1) NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_disks_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_snapshots (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  disk_id BIGINT NULL,
  provider_snapshot_id VARCHAR(128) NULL,
  name VARCHAR(128) NOT NULL,
  size_gb INT NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'READY',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_snapshots_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_backups (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  provider_backup_id VARCHAR(128) NULL,
  name VARCHAR(128) NOT NULL,
  size_gb INT NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'READY',
  expires_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_backups_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_vpc_networks (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  provider_vpc_id VARCHAR(128) NULL,
  name VARCHAR(128) NOT NULL,
  cidr VARCHAR(64) NOT NULL DEFAULT '',
  gateway VARCHAR(128) NOT NULL DEFAULT '',
  interface_name VARCHAR(128) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_vpc_networks_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_security_groups (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  provider_security_group_id VARCHAR(128) NULL,
  name VARCHAR(128) NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_security_groups_service_id (service_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS service_security_group_rules (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  security_group_id BIGINT NOT NULL,
  direction VARCHAR(16) NOT NULL DEFAULT 'ingress',
  protocol VARCHAR(32) NOT NULL DEFAULT 'tcp',
  port_range VARCHAR(64) NOT NULL DEFAULT 'all',
  source_cidr VARCHAR(128) NOT NULL DEFAULT '0.0.0.0/0',
  action VARCHAR(16) NOT NULL DEFAULT 'ALLOW',
  description VARCHAR(255) NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_service_security_group_rules_group_id (security_group_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
