ALTER TABLE product_prices
  ADD COLUMN setup_fee DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER price;

CREATE TABLE IF NOT EXISTS product_config_options (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  product_id BIGINT NOT NULL,
  code VARCHAR(64) NOT NULL,
  name VARCHAR(128) NOT NULL,
  input_type VARCHAR(32) NOT NULL DEFAULT 'select',
  is_required TINYINT(1) NOT NULL DEFAULT 1,
  default_value VARCHAR(128) NOT NULL DEFAULT '',
  description VARCHAR(255) NOT NULL DEFAULT '',
  sort_index INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_product_config_options_product_code (product_id, code),
  INDEX idx_product_config_options_product_id (product_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS product_config_choices (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  option_id BIGINT NOT NULL,
  option_value VARCHAR(128) NOT NULL,
  option_label VARCHAR(128) NOT NULL,
  price_delta DECIMAL(12,2) NOT NULL DEFAULT 0,
  sort_index INT NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_product_config_choices_option_id (option_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

ALTER TABLE orders
  ADD COLUMN configuration_snapshot JSON NULL AFTER status,
  ADD COLUMN resource_snapshot JSON NULL AFTER configuration_snapshot;

ALTER TABLE invoices
  ADD COLUMN order_id BIGINT NULL AFTER customer_id,
  ADD COLUMN order_no VARCHAR(64) NOT NULL DEFAULT '' AFTER order_id,
  ADD COLUMN product_name VARCHAR(128) NOT NULL DEFAULT '' AFTER order_no,
  ADD COLUMN billing_cycle VARCHAR(32) NOT NULL DEFAULT '' AFTER product_name,
  ADD COLUMN paid_at DATETIME NULL AFTER due_at,
  ADD INDEX idx_invoices_order_id (order_id);

ALTER TABLE services
  ADD COLUMN order_id BIGINT NULL AFTER customer_id,
  ADD COLUMN invoice_id BIGINT NULL AFTER order_id,
  ADD COLUMN provider_type VARCHAR(64) NOT NULL DEFAULT 'LOCAL' AFTER product_name,
  ADD COLUMN last_action VARCHAR(64) NOT NULL DEFAULT '' AFTER status,
  ADD COLUMN configuration_snapshot JSON NULL AFTER last_action,
  ADD COLUMN resource_snapshot JSON NULL AFTER configuration_snapshot,
  ADD INDEX idx_services_order_id (order_id),
  ADD INDEX idx_services_invoice_id (invoice_id);

CREATE TABLE IF NOT EXISTS payments (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  payment_no VARCHAR(64) NOT NULL UNIQUE,
  invoice_id BIGINT NOT NULL,
  order_id BIGINT NOT NULL,
  customer_id BIGINT NOT NULL,
  channel VARCHAR(32) NOT NULL DEFAULT 'ONLINE',
  trade_no VARCHAR(128) NOT NULL DEFAULT '',
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  source VARCHAR(32) NOT NULL DEFAULT 'PORTAL',
  status VARCHAR(32) NOT NULL DEFAULT 'COMPLETED',
  operator_name VARCHAR(128) NOT NULL DEFAULT '',
  paid_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_payments_invoice_id (invoice_id),
  INDEX idx_payments_customer_id (customer_id, paid_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS refunds (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  refund_no VARCHAR(64) NOT NULL UNIQUE,
  invoice_id BIGINT NOT NULL,
  order_id BIGINT NOT NULL,
  customer_id BIGINT NOT NULL,
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  reason VARCHAR(255) NOT NULL DEFAULT '',
  status VARCHAR(32) NOT NULL DEFAULT 'COMPLETED',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_refunds_invoice_id (invoice_id),
  INDEX idx_refunds_customer_id (customer_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
