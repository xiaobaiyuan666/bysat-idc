ALTER TABLE customers
  ADD COLUMN balance DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER remarks,
  ADD COLUMN credit_limit DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER balance,
  ADD COLUMN credit_used DECIMAL(12,2) NOT NULL DEFAULT 0 AFTER credit_limit;

CREATE TABLE IF NOT EXISTS account_transactions (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  transaction_no VARCHAR(64) NOT NULL UNIQUE,
  customer_id BIGINT NOT NULL,
  order_id BIGINT NULL,
  invoice_id BIGINT NULL,
  payment_id BIGINT NULL,
  refund_id BIGINT NULL,
  transaction_type VARCHAR(32) NOT NULL DEFAULT 'ADJUSTMENT',
  direction VARCHAR(16) NOT NULL DEFAULT 'IN',
  amount DECIMAL(12,2) NOT NULL DEFAULT 0,
  balance_before DECIMAL(12,2) NOT NULL DEFAULT 0,
  balance_after DECIMAL(12,2) NOT NULL DEFAULT 0,
  credit_before DECIMAL(12,2) NOT NULL DEFAULT 0,
  credit_after DECIMAL(12,2) NOT NULL DEFAULT 0,
  channel VARCHAR(32) NOT NULL DEFAULT 'SYSTEM',
  summary VARCHAR(255) NOT NULL DEFAULT '',
  remark TEXT NULL,
  operator_type VARCHAR(32) NOT NULL DEFAULT 'ADMIN',
  operator_id BIGINT NOT NULL DEFAULT 0,
  operator_name VARCHAR(128) NOT NULL DEFAULT '',
  occurred_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_account_transactions_customer_occurred_at (customer_id, occurred_at),
  INDEX idx_account_transactions_invoice_id (invoice_id),
  INDEX idx_account_transactions_payment_id (payment_id),
  INDEX idx_account_transactions_refund_id (refund_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
