CREATE TABLE IF NOT EXISTS provider_accounts (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  provider_type VARCHAR(64) NOT NULL,
  account_name VARCHAR(128) NOT NULL,
  base_url VARCHAR(255) NOT NULL DEFAULT '',
  username VARCHAR(128) NOT NULL DEFAULT '',
  password VARCHAR(255) NOT NULL DEFAULT '',
  source_name VARCHAR(128) NOT NULL DEFAULT '',
  contact_way VARCHAR(128) NOT NULL DEFAULT '',
  description VARCHAR(255) NOT NULL DEFAULT '',
  account_mode VARCHAR(64) NOT NULL DEFAULT '',
  lang VARCHAR(32) NOT NULL DEFAULT 'zh-cn',
  list_path VARCHAR(255) NOT NULL DEFAULT '',
  detail_path VARCHAR(255) NOT NULL DEFAULT '',
  insecure_skip_verify TINYINT(1) NOT NULL DEFAULT 1,
  auto_update TINYINT(1) NOT NULL DEFAULT 1,
  product_count INT NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'ACTIVE',
  extra_config JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_provider_accounts_name (provider_type, account_name),
  INDEX idx_provider_accounts_type_status (provider_type, status)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE services
  ADD COLUMN provider_account_id BIGINT NOT NULL DEFAULT 0 AFTER provider_type,
  ADD INDEX idx_services_provider_account_id (provider_account_id);

ALTER TABLE orders
  ADD COLUMN provider_account_id BIGINT NOT NULL DEFAULT 0 AFTER automation_type,
  ADD INDEX idx_orders_provider_account_id (provider_account_id);
