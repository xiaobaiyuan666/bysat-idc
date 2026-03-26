ALTER TABLE tickets
  ADD COLUMN service_id BIGINT NULL AFTER customer_id,
  ADD COLUMN content TEXT NULL AFTER title,
  ADD COLUMN priority VARCHAR(32) NOT NULL DEFAULT 'NORMAL' AFTER content,
  ADD COLUMN source VARCHAR(32) NOT NULL DEFAULT 'PORTAL' AFTER priority,
  ADD COLUMN department_name VARCHAR(64) NOT NULL DEFAULT '技术支持' AFTER source,
  ADD COLUMN assigned_admin_id BIGINT NULL AFTER department_name,
  ADD COLUMN assigned_admin_name VARCHAR(128) NOT NULL DEFAULT '' AFTER assigned_admin_id,
  ADD COLUMN latest_reply_excerpt VARCHAR(255) NOT NULL DEFAULT '' AFTER assigned_admin_name,
  ADD COLUMN last_reply_at DATETIME NULL AFTER latest_reply_excerpt,
  ADD COLUMN closed_at DATETIME NULL AFTER last_reply_at;

CREATE TABLE IF NOT EXISTS ticket_replies (
  id BIGINT PRIMARY KEY AUTO_INCREMENT,
  ticket_id BIGINT NOT NULL,
  author_type VARCHAR(32) NOT NULL DEFAULT 'CUSTOMER',
  author_id BIGINT NOT NULL DEFAULT 0,
  author_name VARCHAR(128) NOT NULL DEFAULT '',
  content TEXT NOT NULL,
  is_internal TINYINT(1) NOT NULL DEFAULT 0,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX idx_ticket_replies_ticket_id (ticket_id, created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

UPDATE tickets
SET
  priority = COALESCE(NULLIF(priority, ''), 'NORMAL'),
  source = COALESCE(NULLIF(source, ''), 'PORTAL'),
  department_name = COALESCE(NULLIF(department_name, ''), '技术支持'),
  latest_reply_excerpt = COALESCE(NULLIF(latest_reply_excerpt, ''), title),
  last_reply_at = COALESCE(last_reply_at, updated_at)
WHERE 1 = 1;
