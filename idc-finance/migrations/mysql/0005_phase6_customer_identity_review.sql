ALTER TABLE customer_identities
  ADD COLUMN review_remark VARCHAR(255) NOT NULL DEFAULT '' AFTER country_code;
