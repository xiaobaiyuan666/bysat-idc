ALTER TABLE products
  ADD COLUMN upstream_mapping JSON NULL AFTER resource_template;
