CREATE TABLE IF NOT EXISTS plugin_configs (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  code VARCHAR(64) NOT NULL UNIQUE,
  name VARCHAR(128) NOT NULL,
  enabled TINYINT(1) NOT NULL DEFAULT 0,
  config_json JSON NULL,
  updated_by VARCHAR(128) NOT NULL DEFAULT '',
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_plugin_configs_enabled (enabled)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS plugin_kyc_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  plugin_code VARCHAR(64) NOT NULL,
  customer_id BIGINT NOT NULL,
  real_name VARCHAR(128) NOT NULL,
  id_card_no VARCHAR(64) NOT NULL,
  mobile VARCHAR(32) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  message VARCHAR(255) NOT NULL DEFAULT '',
  raw_json JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_plugin_kyc_customer_created (customer_id, created_at),
  KEY idx_plugin_kyc_plugin_created (plugin_code, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS plugin_sms_records (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  plugin_code VARCHAR(64) NOT NULL,
  customer_id BIGINT NOT NULL,
  mobile VARCHAR(32) NOT NULL,
  template_code VARCHAR(64) NOT NULL DEFAULT '',
  content TEXT NOT NULL,
  status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  message VARCHAR(255) NOT NULL DEFAULT '',
  biz_no VARCHAR(128) NOT NULL DEFAULT '',
  raw_json JSON NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_plugin_sms_customer_created (customer_id, created_at),
  KEY idx_plugin_sms_plugin_created (plugin_code, created_at),
  KEY idx_plugin_sms_biz_no (biz_no)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS plugin_payment_orders (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  plugin_code VARCHAR(64) NOT NULL,
  order_no VARCHAR(64) NOT NULL UNIQUE,
  customer_id BIGINT NOT NULL,
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  currency VARCHAR(16) NOT NULL DEFAULT 'CNY',
  subject VARCHAR(191) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'UNPAID',
  pay_url TEXT NOT NULL,
  notify_url VARCHAR(255) NOT NULL DEFAULT '',
  return_url VARCHAR(255) NOT NULL DEFAULT '',
  trade_no VARCHAR(128) NOT NULL DEFAULT '',
  raw_json JSON NULL,
  paid_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_plugin_payment_customer_created (customer_id, created_at),
  KEY idx_plugin_payment_status_created (status, created_at),
  KEY idx_plugin_payment_plugin_created (plugin_code, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO plugin_configs(code, name, enabled, config_json, updated_by)
VALUES
  ('zhima_kyc', '芝麻实名认证', 1, JSON_OBJECT(), 'system'),
  ('smsbao_sms', '短信宝短信', 1, JSON_OBJECT(), 'system'),
  ('epay', '易支付', 1, JSON_OBJECT(), 'system')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  enabled = VALUES(enabled),
  config_json = VALUES(config_json),
  updated_by = VALUES(updated_by);
