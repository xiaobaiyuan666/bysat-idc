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

INSERT INTO customers (id, customer_no, customer_type, name, email, mobile, status, group_id, level_id, sales_owner, remarks)
VALUES
  (1, 'CUST-00000001', 'COMPANY', '星际网络', 'ops@starring.cn', '13800000001', 'ACTIVE', 1, 1, '华东销售一组', '重点云业务客户'),
  (2, 'CUST-00000002', 'COMPANY', '云脉互联', 'admin@cloudpulse.cn', '13800000003', 'ACTIVE', 2, 2, '华南销售二组', '以弹性云和带宽业务为主')
ON DUPLICATE KEY UPDATE
  customer_type = VALUES(customer_type),
  name = VALUES(name),
  email = VALUES(email),
  mobile = VALUES(mobile),
  status = VALUES(status),
  group_id = VALUES(group_id),
  level_id = VALUES(level_id),
  sales_owner = VALUES(sales_owner),
  remarks = VALUES(remarks);

INSERT INTO customer_contacts (id, customer_id, name, email, mobile, role_name, is_primary)
VALUES
  (1, 1, '林淼', 'ops@starring.cn', '13800000001', '主联系人', 1),
  (2, 1, '陈会', 'finance@starring.cn', '13800000002', '财务联系人', 0),
  (3, 2, '周越', 'admin@cloudpulse.cn', '13800000003', '主联系人', 1)
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  email = VALUES(email),
  mobile = VALUES(mobile),
  role_name = VALUES(role_name),
  is_primary = VALUES(is_primary);

INSERT INTO customer_identities (id, customer_id, identity_type, verify_status, subject_name, cert_no, country_code, review_remark, submitted_at, reviewed_at)
VALUES
  (1, 1, 'COMPANY', 'APPROVED', '上海星际网络科技有限公司', '91310000MA1K123456', 'CN', '历史实名认证数据已通过导入校验', '2026-03-20 17:50:00', '2026-03-20 18:00:00'),
  (2, 2, 'COMPANY', 'PENDING', '深圳云脉互联科技有限公司', '91440300MA5P654321', 'CN', '', '2026-03-20 18:05:00', NULL)
ON DUPLICATE KEY UPDATE
  identity_type = VALUES(identity_type),
  verify_status = VALUES(verify_status),
  subject_name = VALUES(subject_name),
  cert_no = VALUES(cert_no),
  country_code = VALUES(country_code),
  review_remark = VALUES(review_remark),
  submitted_at = VALUES(submitted_at),
  reviewed_at = VALUES(reviewed_at);

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

INSERT INTO tickets (
  id,
  customer_id,
  service_id,
  ticket_no,
  title,
  content,
  priority,
  source,
  department_name,
  assigned_admin_id,
  assigned_admin_name,
  status,
  latest_reply_excerpt,
  last_reply_at
)
VALUES
  (
    1,
    1,
    1,
    'TIC-00000001',
    '实例网络异常排查',
    '客户反馈云主机在晚高峰时段出现网络抖动和丢包，要求排查宿主机网络与带宽配置。',
    'HIGH',
    'PORTAL',
    '技术支持',
    1,
    '系统管理员',
    'PROCESSING',
    '已安排运维检查宿主机网络与安全组配置。',
    '2026-03-20 18:20:00'
  ),
  (
    2,
    2,
    NULL,
    'TIC-00000002',
    '续费发票申请',
    '客户需要为季度续费订单补开发票，并确认开票抬头与税号信息。',
    'NORMAL',
    'PORTAL',
    '财务支持',
    NULL,
    '',
    'OPEN',
    '等待财务确认开票资料。',
    '2026-03-20 18:25:00'
  )
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  title = VALUES(title),
  content = VALUES(content),
  priority = VALUES(priority),
  source = VALUES(source),
  department_name = VALUES(department_name),
  assigned_admin_id = VALUES(assigned_admin_id),
  assigned_admin_name = VALUES(assigned_admin_name),
  status = VALUES(status);

INSERT INTO ticket_replies (id, ticket_id, author_type, author_id, author_name, content, is_internal, created_at)
VALUES
  (1, 1, 'CUSTOMER', 1, '星际网络', '昨晚 20:00 到 20:30 之间业务访问明显抖动，请协助确认。', 0, '2026-03-20 18:05:00'),
  (2, 1, 'ADMIN', 1, '系统管理员', '已安排运维检查宿主机网络与安全组配置，稍后同步检查结果。', 0, '2026-03-20 18:20:00'),
  (3, 2, 'CUSTOMER', 2, '云脉互联', '请为本季度续费账单补开发票，抬头信息和税号已在账户资料中维护。', 0, '2026-03-20 18:25:00')
ON DUPLICATE KEY UPDATE
  author_type = VALUES(author_type),
  author_id = VALUES(author_id),
  author_name = VALUES(author_name),
  content = VALUES(content),
  is_internal = VALUES(is_internal),
  created_at = VALUES(created_at);

INSERT INTO orders (id, order_no, customer_id, customer_name, product_id, product_name, product_type, billing_cycle, amount, status, configuration_snapshot, resource_snapshot)
VALUES
  (
    1,
    'ORD-00000001',
    1,
    '星际网络',
    1,
    '弹性云主机标准型',
    'CLOUD_HOST',
    'monthly',
    629.00,
    'ACTIVE',
    JSON_ARRAY(
      JSON_OBJECT('code', 'cpu', 'name', 'CPU 规格', 'value', '8', 'valueLabel', '8 核', 'priceDelta', 120.00),
      JSON_OBJECT('code', 'memory', 'name', '内存规格', 'value', '16', 'valueLabel', '16 GB', 'priceDelta', 180.00),
      JSON_OBJECT('code', 'backup', 'name', '云备份', 'value', 'enabled', 'valueLabel', '启用', 'priceDelta', 30.00)
    ),
    JSON_OBJECT(
      'regionName', '华东 1',
      'zoneName', '上海可用区 A',
      'hostname', 'srv-cloud-0001',
      'operatingSystem', 'Ubuntu 24.04 LTS',
      'loginUsername', 'root',
      'passwordHint', '最近一次密码重置：2026-03-20 21:40:00',
      'securityGroup', 'default-web',
      'cpuCores', 8,
      'memoryGB', 16,
      'systemDiskGB', 60,
      'dataDiskGB', 100,
      'bandwidthMbps', 20,
      'publicIpv4', '203.0.113.10',
      'publicIpv6', '240e:1234::10'
    )
  )
ON DUPLICATE KEY UPDATE
  customer_name = VALUES(customer_name),
  product_id = VALUES(product_id),
  product_name = VALUES(product_name),
  product_type = VALUES(product_type),
  billing_cycle = VALUES(billing_cycle),
  amount = VALUES(amount),
  status = VALUES(status),
  configuration_snapshot = VALUES(configuration_snapshot),
  resource_snapshot = VALUES(resource_snapshot);

INSERT INTO order_items (id, order_id, product_id, product_name, billing_cycle, quantity, amount)
VALUES
  (1, 1, 1, '弹性云主机标准型', 'monthly', 1, 629.00)
ON DUPLICATE KEY UPDATE
  product_name = VALUES(product_name),
  billing_cycle = VALUES(billing_cycle),
  quantity = VALUES(quantity),
  amount = VALUES(amount);

INSERT INTO invoices (id, customer_id, order_id, order_no, invoice_no, product_name, billing_cycle, status, total_amount, due_at, paid_at)
VALUES
  (1, 1, 1, 'ORD-00000001', 'INV-00000001', '弹性云主机标准型', 'monthly', 'PAID', 629.00, '2026-03-23 21:00:00', '2026-03-20 21:10:00')
ON DUPLICATE KEY UPDATE
  customer_id = VALUES(customer_id),
  order_id = VALUES(order_id),
  order_no = VALUES(order_no),
  product_name = VALUES(product_name),
  billing_cycle = VALUES(billing_cycle),
  status = VALUES(status),
  total_amount = VALUES(total_amount),
  due_at = VALUES(due_at),
  paid_at = VALUES(paid_at);

INSERT INTO invoice_items (id, invoice_id, order_id, product_id, product_name, billing_cycle, amount)
VALUES
  (1, 1, 1, 1, '弹性云主机标准型', 'monthly', 629.00)
ON DUPLICATE KEY UPDATE
  product_name = VALUES(product_name),
  billing_cycle = VALUES(billing_cycle),
  amount = VALUES(amount);

INSERT INTO services (id, customer_id, order_id, invoice_id, service_no, product_name, provider_type, status, next_due_at, last_action, configuration_snapshot, resource_snapshot)
VALUES
  (
    1,
    1,
    1,
    1,
    'SRV-00000001',
    '弹性云主机标准型',
    'MYSQL_PROVIDER',
    'ACTIVE',
    '2026-04-20 21:10:00',
    'activate',
    JSON_ARRAY(
      JSON_OBJECT('code', 'cpu', 'name', 'CPU 规格', 'value', '8', 'valueLabel', '8 核', 'priceDelta', 120.00),
      JSON_OBJECT('code', 'memory', 'name', '内存规格', 'value', '16', 'valueLabel', '16 GB', 'priceDelta', 180.00),
      JSON_OBJECT('code', 'backup', 'name', '云备份', 'value', 'enabled', 'valueLabel', '启用', 'priceDelta', 30.00)
    ),
    JSON_OBJECT(
      'regionName', '华东 1',
      'zoneName', '上海可用区 A',
      'hostname', 'srv-cloud-0001',
      'operatingSystem', 'Ubuntu 24.04 LTS',
      'loginUsername', 'root',
      'passwordHint', '最近一次密码重置：2026-03-20 21:40:00',
      'securityGroup', 'default-web',
      'cpuCores', 8,
      'memoryGB', 16,
      'systemDiskGB', 60,
      'dataDiskGB', 100,
      'bandwidthMbps', 20,
      'publicIpv4', '203.0.113.10',
      'publicIpv6', '240e:1234::10'
    )
  )
ON DUPLICATE KEY UPDATE
  customer_id = VALUES(customer_id),
  order_id = VALUES(order_id),
  invoice_id = VALUES(invoice_id),
  product_name = VALUES(product_name),
  provider_type = VALUES(provider_type),
  status = VALUES(status),
  next_due_at = VALUES(next_due_at),
  last_action = VALUES(last_action),
  configuration_snapshot = VALUES(configuration_snapshot),
  resource_snapshot = VALUES(resource_snapshot);

INSERT INTO payments (id, payment_no, invoice_id, order_id, customer_id, channel, trade_no, amount, source, status, operator_name, paid_at)
VALUES
  (1, 'PAY-00000001', 1, 1, 1, 'OFFLINE', 'BANK-20260320-0001', 629.00, 'ADMIN', 'COMPLETED', '系统管理员', '2026-03-20 21:10:00')
ON DUPLICATE KEY UPDATE
  invoice_id = VALUES(invoice_id),
  order_id = VALUES(order_id),
  customer_id = VALUES(customer_id),
  channel = VALUES(channel),
  trade_no = VALUES(trade_no),
  amount = VALUES(amount),
  source = VALUES(source),
  status = VALUES(status),
  operator_name = VALUES(operator_name),
  paid_at = VALUES(paid_at);

INSERT INTO order_requests (
  id, request_no, order_id, order_no, customer_id, customer_name, product_name, type, status, summary, reason,
  current_amount, requested_amount, current_billing_cycle, requested_billing_cycle,
  source_type, source_id, source_name, processor_type, processor_id, processor_name, process_note, payload_json, processed_at, created_at, updated_at
)
VALUES
  (
    1,
    'REQ-00000001',
    1,
    'ORD-00000001',
    1,
    '演示客户',
    '弹性云主机标准型',
    'RENEW',
    'PENDING',
    '弹性云主机续费请求',
    '客户希望提前确认续费并保持当前配置',
    629.00,
    629.00,
    'monthly',
    'monthly',
    'PORTAL',
    1,
    '演示客户',
    '',
    NULL,
    '',
    '',
    JSON_OBJECT('channel', 'portal', 'priority', 'normal'),
    NULL,
    '2026-03-21 09:00:00',
    '2026-03-21 09:00:00'
  ),
  (
    2,
    'REQ-00000002',
    1,
    'ORD-00000001',
    1,
    '演示客户',
    '弹性云主机标准型',
    'PRICE_ADJUST',
    'APPROVED',
    '云主机改价申请',
    '销售承诺首月优惠价，待财务复核',
    629.00,
    599.00,
    'monthly',
    'monthly',
    'ADMIN',
    1,
    '系统管理员',
    'ADMIN',
    1,
    '系统管理员',
    '同意首月优惠，续费恢复标准价',
    JSON_OBJECT('discountType', 'first_month'),
    '2026-03-21 11:00:00',
    '2026-03-21 10:30:00',
    '2026-03-21 11:00:00'
  )
ON DUPLICATE KEY UPDATE
  request_no = VALUES(request_no),
  order_id = VALUES(order_id),
  order_no = VALUES(order_no),
  customer_id = VALUES(customer_id),
  customer_name = VALUES(customer_name),
  product_name = VALUES(product_name),
  type = VALUES(type),
  status = VALUES(status),
  summary = VALUES(summary),
  reason = VALUES(reason),
  current_amount = VALUES(current_amount),
  requested_amount = VALUES(requested_amount),
  current_billing_cycle = VALUES(current_billing_cycle),
  requested_billing_cycle = VALUES(requested_billing_cycle),
  source_type = VALUES(source_type),
  source_id = VALUES(source_id),
  source_name = VALUES(source_name),
  processor_type = VALUES(processor_type),
  processor_id = VALUES(processor_id),
  processor_name = VALUES(processor_name),
  process_note = VALUES(process_note),
  payload_json = VALUES(payload_json),
  processed_at = VALUES(processed_at),
  updated_at = VALUES(updated_at);

INSERT INTO audit_logs (id, actor_type, actor_id, actor, action, target_type, target_id, target, request_id, description, payload, created_at)
VALUES
  (
    1,
    'ADMIN',
    1,
    '系统管理员',
    'customer.bootstrap',
    'customer',
    1,
    'CUST-00000001',
    'seed-demo-0001',
    '初始化演示客户与联系人数据',
    JSON_OBJECT('customerNo', 'CUST-00000001'),
    '2026-03-20 18:00:00'
  ),
  (
    2,
    'ADMIN',
    1,
    '系统管理员',
    'order.activate',
    'service',
    1,
    'SRV-00000001',
    'seed-demo-0002',
    '演示订单支付后自动激活服务',
    JSON_OBJECT('invoiceNo', 'INV-00000001', 'serviceNo', 'SRV-00000001'),
    '2026-03-20 21:10:00'
  )
ON DUPLICATE KEY UPDATE
  actor = VALUES(actor),
  action = VALUES(action),
  target = VALUES(target),
  request_id = VALUES(request_id),
  description = VALUES(description),
  payload = VALUES(payload),
  created_at = VALUES(created_at);
