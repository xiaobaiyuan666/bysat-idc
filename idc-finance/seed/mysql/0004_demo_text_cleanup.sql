UPDATE product_groups
SET
  name = '弹性云主机标准型',
  description = '弹性云主机标准型 自动维护'
WHERE id = 2
  AND (name LIKE '%?%' OR description LIKE '%?%');

UPDATE products
SET description = '适用于中小型业务站点与通用应用服务'
WHERE id = 1
  AND description LIKE '%?%';
