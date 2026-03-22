ALTER TABLE orders
  ADD COLUMN automation_type VARCHAR(64) NOT NULL DEFAULT 'LOCAL' AFTER product_type;

UPDATE orders
SET automation_type = 'LOCAL'
WHERE automation_type IS NULL OR automation_type = '';
