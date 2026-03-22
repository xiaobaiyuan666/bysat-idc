CREATE TABLE IF NOT EXISTS service_change_orders (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  service_id BIGINT NOT NULL,
  order_id BIGINT NOT NULL UNIQUE,
  invoice_id BIGINT NOT NULL UNIQUE,
  action_name VARCHAR(64) NOT NULL,
  title VARCHAR(191) NOT NULL,
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  status VARCHAR(32) NOT NULL DEFAULT 'PENDING',
  reason TEXT NOT NULL,
  payload_json JSON NULL,
  paid_at DATETIME NULL,
  refunded_at DATETIME NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY idx_service_change_orders_service (service_id, created_at),
  KEY idx_service_change_orders_status (status, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
