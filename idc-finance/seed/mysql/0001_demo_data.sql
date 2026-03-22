INSERT INTO customer_groups (id, name, description)
VALUES
  (1, '旗舰客户', '高价值 IDC 客户'),
  (2, '云业务客户', '以云主机和网络业务为主')
ON DUPLICATE KEY UPDATE
  description = VALUES(description);

INSERT INTO customer_levels (id, name, priority, description)
VALUES
  (1, 'VIP', 100, '重点维护客户'),
  (2, '标准', 50, '默认标准等级')
ON DUPLICATE KEY UPDATE
  priority = VALUES(priority),
  description = VALUES(description);

INSERT INTO product_groups (id, name, description)
VALUES
  (1, '弹性云主机', '标准公有云产品')
ON DUPLICATE KEY UPDATE
  description = VALUES(description);

INSERT INTO products (id, product_no, group_id, name, description, product_type, status, resource_template)
VALUES
  (
    1,
    'PROD-000001',
    1,
    '弹性云主机标准型',
    '适用于中小型业务站点与通用应用服务',
    'CLOUD_HOST',
    'ACTIVE',
    JSON_OBJECT(
      'regionName', '华东 1',
      'zoneName', '上海可用区 A',
      'operatingSystem', 'Ubuntu 24.04 LTS',
      'loginUsername', 'root',
      'securityGroup', 'default-web',
      'cpuCores', 4,
      'memoryGB', 8,
      'systemDiskGB', 60,
      'dataDiskGB', 100,
      'bandwidthMbps', 20,
      'publicIpCount', 1
    )
  )
ON DUPLICATE KEY UPDATE
  group_id = VALUES(group_id),
  name = VALUES(name),
  description = VALUES(description),
  product_type = VALUES(product_type),
  status = VALUES(status),
  resource_template = VALUES(resource_template);

INSERT INTO product_prices (id, product_id, cycle_code, cycle_name, price, setup_fee)
VALUES
  (1, 1, 'monthly', '月付', 299.00, 30.00),
  (2, 1, 'quarterly', '季付', 849.00, 0.00),
  (3, 1, 'annual', '年付', 3199.00, 0.00)
ON DUPLICATE KEY UPDATE
  cycle_name = VALUES(cycle_name),
  price = VALUES(price),
  setup_fee = VALUES(setup_fee);

INSERT INTO product_config_options (id, product_id, code, name, input_type, is_required, default_value, description, sort_index)
VALUES
  (1, 1, 'cpu', 'CPU 规格', 'select', 1, '4', '选择实例 CPU 配置', 1),
  (2, 1, 'memory', '内存规格', 'select', 1, '8', '选择实例内存配置', 2),
  (3, 1, 'backup', '云备份', 'select', 0, 'disabled', '是否启用云备份服务', 3)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  input_type = VALUES(input_type),
  is_required = VALUES(is_required),
  default_value = VALUES(default_value),
  description = VALUES(description),
  sort_index = VALUES(sort_index);

INSERT INTO product_config_choices (id, option_id, option_value, option_label, price_delta, sort_index)
VALUES
  (1, 1, '4', '4 核', 0.00, 1),
  (2, 1, '8', '8 核', 120.00, 2),
  (3, 1, '16', '16 核', 320.00, 3),
  (4, 2, '8', '8 GB', 0.00, 1),
  (5, 2, '16', '16 GB', 180.00, 2),
  (6, 2, '32', '32 GB', 380.00, 3),
  (7, 3, 'disabled', '关闭', 0.00, 1),
  (8, 3, 'enabled', '启用', 30.00, 2)
ON DUPLICATE KEY UPDATE
  option_label = VALUES(option_label),
  price_delta = VALUES(price_delta),
  sort_index = VALUES(sort_index);
