-- 无穷云 IDC 系统 MySQL 演示数据
-- 目标：让新设备在执行 seed 后，直接得到一套可用于后台、用户端、报表和 Provider 联调的完整演示链路。
-- 说明：
-- 1. 所有数据都使用固定主键，方便重复执行。
-- 2. 所有敏感字段只放演示占位值，不放真实凭据。
-- 3. 订单、账单、支付、服务时间使用相对时间，保证重新 seed 后报表与工作台仍然有可演示的数据。

-- -----------------------------------------------------------------------------
-- 客户基础资料
-- -----------------------------------------------------------------------------

INSERT INTO customer_groups (id, name, description)
VALUES
  (1, '旗舰客户', '重点维护的大客户与关键收入客户'),
  (2, '云业务客户', '以云主机、网络和资源型业务为主'),
  (3, '托管客户', '以机柜托管和物理资源交付为主')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  description = VALUES(description),
  updated_at = NOW();

INSERT INTO customer_levels (id, name, priority, description)
VALUES
  (1, 'VIP', 100, '重点维护客户等级'),
  (2, '标准', 50, '默认标准等级')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  priority = VALUES(priority),
  description = VALUES(description),
  updated_at = NOW();

INSERT INTO customers (
  id, customer_no, customer_type, name, email, mobile, status, group_id, level_id, sales_owner, remarks
)
VALUES
  (
    1,
    'CUST-20260320-0001',
    'COMPANY',
    '星环网络',
    'ops@starring.cn',
    '13800000001',
    'ACTIVE',
    1,
    1,
    '华东销售一组',
    '重点云业务客户，用于演示魔方云实例、服务资源和改配链路。'
  ),
  (
    2,
    'CUST-20260320-0002',
    'COMPANY',
    '云脉互联',
    'admin@cloudpulse.cn',
    '13800000003',
    'ACTIVE',
    2,
    2,
    '华南销售二组',
    '以带宽与托管业务为主，用于演示欠费暂停与即将到期服务场景。'
  )
ON DUPLICATE KEY UPDATE
  customer_no = VALUES(customer_no),
  customer_type = VALUES(customer_type),
  name = VALUES(name),
  email = VALUES(email),
  mobile = VALUES(mobile),
  status = VALUES(status),
  group_id = VALUES(group_id),
  level_id = VALUES(level_id),
  sales_owner = VALUES(sales_owner),
  remarks = VALUES(remarks),
  updated_at = NOW();

INSERT INTO customer_contacts (
  id, customer_id, name, email, mobile, role_name, is_primary
)
VALUES
  (1, 1, '林浩', 'ops@starring.cn', '13800000001', '主联系人', 1),
  (2, 1, '陈会', 'finance@starring.cn', '13800000002', '财务联系人', 0),
  (3, 2, '周越', 'admin@cloudpulse.cn', '13800000003', '主联系人', 1)
ON DUPLICATE KEY UPDATE
  customer_id = VALUES(customer_id),
  name = VALUES(name),
  email = VALUES(email),
  mobile = VALUES(mobile),
  role_name = VALUES(role_name),
  is_primary = VALUES(is_primary),
  updated_at = NOW();

INSERT INTO customer_identities (
  id, customer_id, identity_type, verify_status, subject_name, cert_no, country_code, review_remark, submitted_at, reviewed_at
)
VALUES
  (
    1,
    1,
    'COMPANY',
    'APPROVED',
    '上海星环网络科技有限公司',
    '91310000MA1K123456',
    'CN',
    '历史实名资料已完成核验导入',
    DATE_SUB(NOW(), INTERVAL 8 DAY),
    DATE_SUB(NOW(), INTERVAL 7 DAY)
  ),
  (
    2,
    2,
    'COMPANY',
    'PENDING',
    '深圳云脉互联科技有限公司',
    '91440300MA5P654321',
    'CN',
    '',
    DATE_SUB(NOW(), INTERVAL 1 DAY),
    NULL
  )
ON DUPLICATE KEY UPDATE
  customer_id = VALUES(customer_id),
  identity_type = VALUES(identity_type),
  verify_status = VALUES(verify_status),
  subject_name = VALUES(subject_name),
  cert_no = VALUES(cert_no),
  country_code = VALUES(country_code),
  review_remark = VALUES(review_remark),
  submitted_at = VALUES(submitted_at),
  reviewed_at = VALUES(reviewed_at),
  updated_at = NOW();

INSERT INTO tickets (
  id, customer_id, ticket_no, title, status, created_at, updated_at
)
VALUES
  (
    1,
    1,
    'TIC-20260320-0001',
    '实例网络异常排查',
    'PROCESSING',
    DATE_SUB(NOW(), INTERVAL 10 HOUR),
    DATE_SUB(NOW(), INTERVAL 2 HOUR)
  ),
  (
    2,
    2,
    'TIC-20260320-0002',
    '托管设备上架预约确认',
    'OPEN',
    DATE_SUB(NOW(), INTERVAL 5 HOUR),
    DATE_SUB(NOW(), INTERVAL 1 HOUR)
  )
ON DUPLICATE KEY UPDATE
  customer_id = VALUES(customer_id),
  ticket_no = VALUES(ticket_no),
  title = VALUES(title),
  status = VALUES(status),
  created_at = VALUES(created_at),
  updated_at = VALUES(updated_at);

-- -----------------------------------------------------------------------------
-- 商品、价格矩阵、配置项、资源模板
-- -----------------------------------------------------------------------------

INSERT INTO product_groups (id, name, description)
VALUES
  (1, '云主机', '标准公有云与弹性实例商品'),
  (2, '网络', '带宽与网络增强商品'),
  (3, '机柜托管', '物理资源与托管商品')
ON DUPLICATE KEY UPDATE
  name = VALUES(name),
  description = VALUES(description),
  updated_at = NOW();

INSERT INTO products (
  id, product_no, group_id, name, description, product_type, status, resource_template, automation_config, upstream_mapping
)
VALUES
  (
    1,
    'PROD-000001',
    1,
    '弹性云主机 CN2 标准型',
    '适合中小型 IDC 客户业务部署，含 1 个公网 IPv4，可按月或按年售卖。',
    'CLOUD',
    'ACTIVE',
    '{"regionName":"华南广州","zoneName":"广州三区","operatingSystem":"Rocky Linux 9","loginUsername":"root","securityGroup":"default-cloud","cpuCores":4,"memoryGB":8,"systemDiskGB":60,"dataDiskGB":120,"bandwidthMbps":20,"publicIpCount":1}',
    '{"providerAccountId":1,"channel":"MOFANG_CLOUD","moduleType":"DCIMCLOUD","provisionStage":"AFTER_PAYMENT","autoProvision":true,"serverGroup":"mf-cloud-default","providerNode":"gz-a-01"}',
    '{"providerAccountId":1,"providerType":"MOFANG_CLOUD","sourceName":"魔方云演示环境","remoteProductCode":"mf-cloud-standard","remoteProductName":"弹性云主机标准型","pricePolicy":"CUSTOM","autoSyncPricing":true,"autoSyncConfig":true,"autoSyncTemplate":true,"syncStatus":"SUCCESS","syncMessage":"已预置演示映射","lastSyncedAt":"2026-03-21T10:00:00+08:00"}'
  ),
  (
    2,
    'PROD-000002',
    2,
    '精品带宽 50M',
    '适合高并发 Web 业务，支持后续扩容和带宽包升级。',
    'BANDWIDTH',
    'ACTIVE',
    '{"regionName":"华南广州","zoneName":"骨干网络","operatingSystem":"-","loginUsername":"-","securityGroup":"network-edge","cpuCores":0,"memoryGB":0,"systemDiskGB":0,"dataDiskGB":0,"bandwidthMbps":50,"publicIpCount":1}',
    '{"providerAccountId":0,"channel":"LOCAL","moduleType":"NORMAL","provisionStage":"AFTER_PAYMENT","autoProvision":true,"serverGroup":"","providerNode":""}',
    '{"providerAccountId":0,"providerType":"LOCAL","sourceName":"本地交付","remoteProductCode":"","remoteProductName":"","pricePolicy":"FIXED","autoSyncPricing":false,"autoSyncConfig":false,"autoSyncTemplate":false,"syncStatus":"IDLE","syncMessage":"本地交付商品无需上游映射","lastSyncedAt":""}'
  ),
  (
    3,
    'PROD-000003',
    3,
    '1U 标准托管',
    '适合物理服务器托管业务，含基础运维巡检与上架服务。',
    'COLOCATION',
    'ACTIVE',
    '{"regionName":"华东上海","zoneName":"A1 机房","operatingSystem":"-","loginUsername":"-","securityGroup":"cabinet-standard","cpuCores":0,"memoryGB":0,"systemDiskGB":0,"dataDiskGB":0,"bandwidthMbps":20,"publicIpCount":1}',
    '{"providerAccountId":0,"channel":"LOCAL","moduleType":"DCIM","provisionStage":"AFTER_PAYMENT","autoProvision":true,"serverGroup":"","providerNode":"rack-a1"}',
    '{"providerAccountId":0,"providerType":"LOCAL","sourceName":"本地交付","remoteProductCode":"","remoteProductName":"","pricePolicy":"FIXED","autoSyncPricing":false,"autoSyncConfig":false,"autoSyncTemplate":false,"syncStatus":"IDLE","syncMessage":"托管商品无需上游映射","lastSyncedAt":""}'
  )
ON DUPLICATE KEY UPDATE
  product_no = VALUES(product_no),
  group_id = VALUES(group_id),
  name = VALUES(name),
  description = VALUES(description),
  product_type = VALUES(product_type),
  status = VALUES(status),
  resource_template = VALUES(resource_template),
  automation_config = VALUES(automation_config),
  upstream_mapping = VALUES(upstream_mapping),
  updated_at = NOW();

INSERT INTO product_prices (id, product_id, cycle_code, cycle_name, price, setup_fee)
VALUES
  (1, 1, 'monthly', '月付', 199.00, 0.00),
  (2, 1, 'quarterly', '季付', 569.00, 0.00),
  (3, 1, 'annual', '年付', 1999.00, 0.00),
  (4, 2, 'monthly', '月付', 299.00, 50.00),
  (5, 2, 'annual', '年付', 2990.00, 0.00),
  (6, 3, 'monthly', '月付', 499.00, 200.00),
  (7, 3, 'annual', '年付', 4990.00, 0.00)
ON DUPLICATE KEY UPDATE
  product_id = VALUES(product_id),
  cycle_code = VALUES(cycle_code),
  cycle_name = VALUES(cycle_name),
  price = VALUES(price),
  setup_fee = VALUES(setup_fee),
  updated_at = NOW();

INSERT INTO product_config_options (
  id, product_id, code, name, input_type, is_required, default_value, description, sort_index
)
VALUES
  (1, 1, 'cpu', 'CPU 规格', 'select', 1, '4', '影响实例可分配的 vCPU 核数。', 1),
  (2, 1, 'memory', '内存规格', 'select', 1, '8', '按 GB 计费，可和 CPU 规格组合售卖。', 2),
  (3, 1, 'backup', '云备份', 'radio', 0, 'enabled', '是否附带每周快照备份服务。', 3),
  (4, 2, 'commit', '承诺带宽', 'select', 1, '50', '固定承诺带宽。', 1),
  (5, 3, 'power', '电力额度', 'select', 1, '300', '默认含 A/B 路冗余电力。', 1)
ON DUPLICATE KEY UPDATE
  product_id = VALUES(product_id),
  code = VALUES(code),
  name = VALUES(name),
  input_type = VALUES(input_type),
  is_required = VALUES(is_required),
  default_value = VALUES(default_value),
  description = VALUES(description),
  sort_index = VALUES(sort_index),
  updated_at = NOW();

INSERT INTO product_config_choices (
  id, option_id, option_value, option_label, price_delta, sort_index
)
VALUES
  (1, 1, '2', '2 核', 0.00, 1),
  (2, 1, '4', '4 核', 80.00, 2),
  (3, 1, '8', '8 核', 220.00, 3),
  (4, 2, '4', '4 GB', 0.00, 1),
  (5, 2, '8', '8 GB', 60.00, 2),
  (6, 2, '16', '16 GB', 180.00, 3),
  (7, 3, 'enabled', '启用', 30.00, 1),
  (8, 3, 'disabled', '关闭', 0.00, 2),
  (9, 4, '50', '50 Mbps', 0.00, 1),
  (10, 4, '100', '100 Mbps', 260.00, 2),
  (11, 4, '200', '200 Mbps', 680.00, 3),
  (12, 5, '300', '300W', 0.00, 1),
  (13, 5, '500', '500W', 120.00, 2),
  (14, 5, '800', '800W', 280.00, 3)
ON DUPLICATE KEY UPDATE
  option_id = VALUES(option_id),
  option_value = VALUES(option_value),
  option_label = VALUES(option_label),
  price_delta = VALUES(price_delta),
  sort_index = VALUES(sort_index),
  updated_at = NOW();

-- -----------------------------------------------------------------------------
-- Provider 账号与商品映射
-- -----------------------------------------------------------------------------

INSERT INTO provider_accounts (
  id, provider_type, account_name, base_url, username, password, source_name, contact_way, description,
  account_mode, lang, list_path, detail_path, insecure_skip_verify, auto_update, product_count, status, extra_config
)
VALUES
  (
    1,
    'MOFANG_CLOUD',
    '魔方云演示账号',
    'https://demo-mofang.example.com',
    'demo-mofang-user',
    'demo-mofang-password',
    '魔方云演示环境',
    '仅用于本地联调演示',
    '该账号为演示占位数据，不包含真实生产凭据。',
    'API',
    'zh-cn',
    '/v1/clouds?page=1&per_page=100&sort=desc&orderby=id',
    '/v1/clouds/:id',
    1,
    1,
    1,
    'ACTIVE',
    '{"note":"seed-only demo account","owner":"repository"}'
  )
ON DUPLICATE KEY UPDATE
  provider_type = VALUES(provider_type),
  account_name = VALUES(account_name),
  base_url = VALUES(base_url),
  username = VALUES(username),
  password = VALUES(password),
  source_name = VALUES(source_name),
  contact_way = VALUES(contact_way),
  description = VALUES(description),
  account_mode = VALUES(account_mode),
  lang = VALUES(lang),
  list_path = VALUES(list_path),
  detail_path = VALUES(detail_path),
  insecure_skip_verify = VALUES(insecure_skip_verify),
  auto_update = VALUES(auto_update),
  product_count = VALUES(product_count),
  status = VALUES(status),
  extra_config = VALUES(extra_config),
  updated_at = NOW();

-- -----------------------------------------------------------------------------
-- 订单、账单、支付、服务主链路
-- -----------------------------------------------------------------------------

INSERT INTO orders (
  id, order_no, customer_id, customer_name, product_id, product_name, product_type, automation_type, provider_account_id,
  billing_cycle, amount, status, configuration_snapshot, resource_snapshot, created_at, updated_at
)
VALUES
  (
    1,
    'ORD-20260320-0001',
    1,
    '星环网络',
    1,
    '弹性云主机 CN2 标准型',
    'CLOUD',
    'MOFANG_CLOUD',
    1,
    'annual',
    1999.00,
    'ACTIVE',
    '[{"code":"cpu","name":"CPU 规格","value":"4","valueLabel":"4 核","priceDelta":80},{"code":"memory","name":"内存规格","value":"8","valueLabel":"8 GB","priceDelta":60},{"code":"backup","name":"云备份","value":"enabled","valueLabel":"启用","priceDelta":30}]',
    '{"regionName":"华南广州","zoneName":"广州三区","hostname":"srv-00000001.idc.local","operatingSystem":"Rocky Linux 9","loginUsername":"root","passwordHint":"初始化密码已通过站内信下发","securityGroup":"default-cloud","cpuCores":4,"memoryGB":8,"systemDiskGB":60,"dataDiskGB":120,"bandwidthMbps":20,"publicIpv4":"203.0.113.11","publicIpv6":"2408:4004:10::11"}',
    DATE_SUB(NOW(), INTERVAL 3 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY)
  ),
  (
    2,
    'ORD-20260320-0002',
    1,
    '星环网络',
    2,
    '精品带宽 50M',
    'BANDWIDTH',
    'LOCAL',
    0,
    'monthly',
    299.00,
    'PENDING',
    '[{"code":"commit","name":"承诺带宽","value":"50","valueLabel":"50 Mbps","priceDelta":0}]',
    '{"regionName":"华南广州","zoneName":"骨干网络","hostname":"network-addon-0002","operatingSystem":"-","loginUsername":"-","passwordHint":"-","securityGroup":"network-edge","cpuCores":0,"memoryGB":0,"systemDiskGB":0,"dataDiskGB":0,"bandwidthMbps":50,"publicIpv4":"203.0.113.51","publicIpv6":""}',
    DATE_SUB(NOW(), INTERVAL 12 HOUR),
    DATE_SUB(NOW(), INTERVAL 6 HOUR)
  ),
  (
    3,
    'ORD-20260320-0003',
    2,
    '云脉互联',
    3,
    '1U 标准托管',
    'COLOCATION',
    'LOCAL',
    0,
    'annual',
    4990.00,
    'COMPLETED',
    '[{"code":"power","name":"电力额度","value":"500","valueLabel":"500W","priceDelta":120}]',
    '{"regionName":"华东上海","zoneName":"A1 机房","hostname":"rack-a1-1u-0001","operatingSystem":"-","loginUsername":"-","passwordHint":"由运维现场交付","securityGroup":"cabinet-standard","cpuCores":0,"memoryGB":0,"systemDiskGB":0,"dataDiskGB":0,"bandwidthMbps":20,"publicIpv4":"198.51.100.20","publicIpv6":""}',
    DATE_SUB(NOW(), INTERVAL 10 DAY),
    DATE_SUB(NOW(), INTERVAL 9 DAY)
  ),
  (
    4,
    'ORD-20260320-0004',
    1,
    '星环网络',
    1,
    '弹性云主机 CN2 标准型 升配单',
    'SERVICE_CHANGE',
    'SERVICE_CHANGE',
    1,
    'annual',
    260.00,
    'PENDING',
    '[{"code":"cpu","name":"CPU 规格","value":"8","valueLabel":"8 核","priceDelta":220},{"code":"memory","name":"内存规格","value":"16","valueLabel":"16 GB","priceDelta":180}]',
    '{"regionName":"华南广州","zoneName":"广州三区","hostname":"srv-00000001.idc.local","operatingSystem":"Rocky Linux 9","loginUsername":"root","passwordHint":"支付改配单后自动下发新规格","securityGroup":"default-cloud","cpuCores":8,"memoryGB":16,"systemDiskGB":60,"dataDiskGB":120,"bandwidthMbps":20,"publicIpv4":"203.0.113.11","publicIpv6":"2408:4004:10::11"}',
    DATE_SUB(NOW(), INTERVAL 3 HOUR),
    DATE_SUB(NOW(), INTERVAL 1 HOUR)
  )
ON DUPLICATE KEY UPDATE
  order_no = VALUES(order_no),
  customer_id = VALUES(customer_id),
  customer_name = VALUES(customer_name),
  product_id = VALUES(product_id),
  product_name = VALUES(product_name),
  product_type = VALUES(product_type),
  automation_type = VALUES(automation_type),
  provider_account_id = VALUES(provider_account_id),
  billing_cycle = VALUES(billing_cycle),
  amount = VALUES(amount),
  status = VALUES(status),
  configuration_snapshot = VALUES(configuration_snapshot),
  resource_snapshot = VALUES(resource_snapshot),
  created_at = VALUES(created_at),
  updated_at = VALUES(updated_at);

INSERT INTO order_items (
  id, order_id, product_id, product_name, billing_cycle, quantity, amount, created_at
)
VALUES
  (1, 1, 1, '弹性云主机 CN2 标准型', 'annual', 1, 1999.00, DATE_SUB(NOW(), INTERVAL 3 DAY)),
  (2, 2, 2, '精品带宽 50M', 'monthly', 1, 299.00, DATE_SUB(NOW(), INTERVAL 12 HOUR)),
  (3, 3, 3, '1U 标准托管', 'annual', 1, 4990.00, DATE_SUB(NOW(), INTERVAL 10 DAY)),
  (4, 4, 1, '弹性云主机 CN2 标准型 升配单', 'annual', 1, 260.00, DATE_SUB(NOW(), INTERVAL 3 HOUR))
ON DUPLICATE KEY UPDATE
  order_id = VALUES(order_id),
  product_id = VALUES(product_id),
  product_name = VALUES(product_name),
  billing_cycle = VALUES(billing_cycle),
  quantity = VALUES(quantity),
  amount = VALUES(amount),
  created_at = VALUES(created_at);

INSERT INTO invoices (
  id, customer_id, order_id, order_no, invoice_no, product_name, billing_cycle, status, total_amount, due_at, paid_at, created_at, updated_at
)
VALUES
  (
    1,
    1,
    1,
    'ORD-20260320-0001',
    'INV-20260320-0001',
    '弹性云主机 CN2 标准型',
    'annual',
    'PAID',
    1999.00,
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 3 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY)
  ),
  (
    2,
    1,
    2,
    'ORD-20260320-0002',
    'INV-20260320-0002',
    '精品带宽 50M',
    'monthly',
    'UNPAID',
    299.00,
    DATE_SUB(NOW(), INTERVAL 1 DAY),
    NULL,
    DATE_SUB(NOW(), INTERVAL 12 HOUR),
    DATE_SUB(NOW(), INTERVAL 6 HOUR)
  ),
  (
    3,
    2,
    3,
    'ORD-20260320-0003',
    'INV-20260320-0003',
    '1U 标准托管',
    'annual',
    'PAID',
    4990.00,
    DATE_SUB(NOW(), INTERVAL 9 DAY),
    DATE_SUB(NOW(), INTERVAL 9 DAY),
    DATE_SUB(NOW(), INTERVAL 10 DAY),
    DATE_SUB(NOW(), INTERVAL 9 DAY)
  ),
  (
    4,
    1,
    4,
    'ORD-20260320-0004',
    'INV-20260320-0004',
    '弹性云主机 CN2 标准型 升配单',
    'annual',
    'UNPAID',
    260.00,
    DATE_ADD(NOW(), INTERVAL 2 DAY),
    NULL,
    DATE_SUB(NOW(), INTERVAL 3 HOUR),
    DATE_SUB(NOW(), INTERVAL 1 HOUR)
  )
ON DUPLICATE KEY UPDATE
  customer_id = VALUES(customer_id),
  order_id = VALUES(order_id),
  order_no = VALUES(order_no),
  invoice_no = VALUES(invoice_no),
  product_name = VALUES(product_name),
  billing_cycle = VALUES(billing_cycle),
  status = VALUES(status),
  total_amount = VALUES(total_amount),
  due_at = VALUES(due_at),
  paid_at = VALUES(paid_at),
  created_at = VALUES(created_at),
  updated_at = VALUES(updated_at);

INSERT INTO invoice_items (
  id, invoice_id, order_id, product_id, product_name, billing_cycle, amount, created_at
)
VALUES
  (1, 1, 1, 1, '弹性云主机 CN2 标准型', 'annual', 1999.00, DATE_SUB(NOW(), INTERVAL 3 DAY)),
  (2, 2, 2, 2, '精品带宽 50M', 'monthly', 299.00, DATE_SUB(NOW(), INTERVAL 12 HOUR)),
  (3, 3, 3, 3, '1U 标准托管', 'annual', 4990.00, DATE_SUB(NOW(), INTERVAL 10 DAY)),
  (4, 4, 4, 1, '弹性云主机 CN2 标准型 升配单', 'annual', 260.00, DATE_SUB(NOW(), INTERVAL 3 HOUR))
ON DUPLICATE KEY UPDATE
  invoice_id = VALUES(invoice_id),
  order_id = VALUES(order_id),
  product_id = VALUES(product_id),
  product_name = VALUES(product_name),
  billing_cycle = VALUES(billing_cycle),
  amount = VALUES(amount),
  created_at = VALUES(created_at);

INSERT INTO payments (
  id, payment_no, invoice_id, order_id, customer_id, channel, trade_no, amount, source, status, operator_name, paid_at, created_at, updated_at
)
VALUES
  (
    1,
    'PAY-20260320-0001',
    1,
    1,
    1,
    'ONLINE',
    'TRADE-20260320-0001',
    1999.00,
    'PORTAL',
    'COMPLETED',
    '星环网络',
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY)
  ),
  (
    2,
    'PAY-20260320-0002',
    3,
    3,
    2,
    'OFFLINE',
    'TRADE-20260320-0002',
    4990.00,
    'ADMIN',
    'COMPLETED',
    '财务专员',
    DATE_SUB(NOW(), INTERVAL 9 DAY),
    DATE_SUB(NOW(), INTERVAL 9 DAY),
    DATE_SUB(NOW(), INTERVAL 9 DAY)
  )
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
  paid_at = VALUES(paid_at),
  created_at = VALUES(created_at),
  updated_at = VALUES(updated_at);

INSERT INTO services (
  id, service_no, customer_id, order_id, invoice_id, product_name, provider_type, provider_account_id, provider_resource_id,
  region_name, ip_address, status, sync_status, sync_message, next_due_at, last_action, last_sync_at,
  configuration_snapshot, resource_snapshot, created_at, updated_at
)
VALUES
  (
    1,
    'SRV-20260320-0001',
    1,
    1,
    1,
    '弹性云主机 CN2 标准型',
    'MOFANG_CLOUD',
    1,
    'mf-cloud-1001',
    '华南广州',
    '203.0.113.11',
    'ACTIVE',
    'SUCCESS',
    '魔方云实例已同步到本地资源视图',
    DATE_ADD(CURDATE(), INTERVAL 330 DAY),
    'sync',
    DATE_SUB(NOW(), INTERVAL 4 HOUR),
    '[{"code":"cpu","name":"CPU 规格","value":"4","valueLabel":"4 核","priceDelta":80},{"code":"memory","name":"内存规格","value":"8","valueLabel":"8 GB","priceDelta":60},{"code":"backup","name":"云备份","value":"enabled","valueLabel":"启用","priceDelta":30}]',
    '{"regionName":"华南广州","zoneName":"广州三区","hostname":"srv-00000001.idc.local","operatingSystem":"Rocky Linux 9","loginUsername":"root","passwordHint":"初始化密码已通过站内信下发","securityGroup":"default-cloud","cpuCores":4,"memoryGB":8,"systemDiskGB":60,"dataDiskGB":120,"bandwidthMbps":20,"publicIpv4":"203.0.113.11","publicIpv6":"2408:4004:10::11"}',
    DATE_SUB(NOW(), INTERVAL 3 DAY),
    DATE_SUB(NOW(), INTERVAL 2 HOUR)
  ),
  (
    2,
    'SRV-20260320-0002',
    2,
    3,
    3,
    '1U 标准托管',
    'LOCAL',
    0,
    NULL,
    '华东上海',
    '198.51.100.20',
    'SUSPENDED',
    'SUCCESS',
    '本地交付服务，当前为欠费暂停演示状态',
    DATE_ADD(CURDATE(), INTERVAL 3 DAY),
    'suspend',
    DATE_SUB(NOW(), INTERVAL 1 DAY),
    '[{"code":"power","name":"电力额度","value":"500","valueLabel":"500W","priceDelta":120}]',
    '{"regionName":"华东上海","zoneName":"A1 机房","hostname":"rack-a1-1u-0001","operatingSystem":"-","loginUsername":"-","passwordHint":"由运维现场交付","securityGroup":"cabinet-standard","cpuCores":0,"memoryGB":0,"systemDiskGB":0,"dataDiskGB":0,"bandwidthMbps":20,"publicIpv4":"198.51.100.20","publicIpv6":""}',
    DATE_SUB(NOW(), INTERVAL 10 DAY),
    DATE_SUB(NOW(), INTERVAL 1 DAY)
  )
ON DUPLICATE KEY UPDATE
  service_no = VALUES(service_no),
  customer_id = VALUES(customer_id),
  order_id = VALUES(order_id),
  invoice_id = VALUES(invoice_id),
  product_name = VALUES(product_name),
  provider_type = VALUES(provider_type),
  provider_account_id = VALUES(provider_account_id),
  provider_resource_id = VALUES(provider_resource_id),
  region_name = VALUES(region_name),
  ip_address = VALUES(ip_address),
  status = VALUES(status),
  sync_status = VALUES(sync_status),
  sync_message = VALUES(sync_message),
  next_due_at = VALUES(next_due_at),
  last_action = VALUES(last_action),
  last_sync_at = VALUES(last_sync_at),
  configuration_snapshot = VALUES(configuration_snapshot),
  resource_snapshot = VALUES(resource_snapshot),
  created_at = VALUES(created_at),
  updated_at = VALUES(updated_at);

-- -----------------------------------------------------------------------------
-- 改配单与自动化任务
-- -----------------------------------------------------------------------------

INSERT INTO service_change_orders (
  id, service_id, order_id, invoice_id, action_name, title, amount, status, reason, payload_json, paid_at, refunded_at, created_at, updated_at
)
VALUES
  (
    1,
    1,
    4,
    4,
    'upgrade',
    '弹性云主机 CN2 标准型 升配单',
    260.00,
    'UNPAID',
    '客户申请将现有实例升级到 8 核 16 GB。',
    '{"targetConfiguration":[{"code":"cpu","value":"8"},{"code":"memory","value":"16"}],"trigger":"admin-demo"}',
    NULL,
    NULL,
    DATE_SUB(NOW(), INTERVAL 3 HOUR),
    DATE_SUB(NOW(), INTERVAL 1 HOUR)
  )
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  order_id = VALUES(order_id),
  invoice_id = VALUES(invoice_id),
  action_name = VALUES(action_name),
  title = VALUES(title),
  amount = VALUES(amount),
  status = VALUES(status),
  reason = VALUES(reason),
  payload_json = VALUES(payload_json),
  paid_at = VALUES(paid_at),
  refunded_at = VALUES(refunded_at),
  created_at = VALUES(created_at),
  updated_at = VALUES(updated_at);

INSERT INTO automation_tasks (
  id, task_no, task_type, title, channel, stage_name, status, source_type, source_id, customer_id, customer_name,
  product_name, order_id, invoice_id, service_id, service_no, provider_type, provider_resource_id, action_name,
  operator_type, operator_name, request_payload, result_payload, message, created_at, started_at, finished_at, updated_at
)
VALUES
  (
    1,
    'TASK-20260320-0001',
    'AUTO_PROVISION',
    '商品自动开通',
    'MOFANG_CLOUD',
    'AFTER_PAYMENT',
    'SUCCESS',
    'order',
    1,
    1,
    '星环网络',
    '弹性云主机 CN2 标准型',
    1,
    1,
    1,
    'SRV-20260320-0001',
    'MOFANG_CLOUD',
    'mf-cloud-1001',
    'provision',
    'SYSTEM',
    'seed',
    '{"orderNo":"ORD-20260320-0001","providerAccountId":1}',
    '{"serviceNo":"SRV-20260320-0001","syncStatus":"SUCCESS"}',
    '演示实例已完成自动开通并写入本地服务表。',
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY),
    DATE_SUB(NOW(), INTERVAL 2 DAY)
  ),
  (
    2,
    'TASK-20260320-0002',
    'SERVICE_CHANGE',
    '服务改配待执行',
    'MOFANG_CLOUD',
    'AFTER_PAYMENT',
    'PENDING',
    'order',
    4,
    1,
    '星环网络',
    '弹性云主机 CN2 标准型 升配单',
    4,
    4,
    1,
    'SRV-20260320-0001',
    'MOFANG_CLOUD',
    'mf-cloud-1001',
    'upgrade',
    'ADMIN',
    '系统管理员',
    '{"changeOrderId":1,"targetCpu":8,"targetMemory":16}',
    NULL,
    '等待客户支付改配单后执行规格升级。',
    DATE_SUB(NOW(), INTERVAL 2 HOUR),
    NULL,
    NULL,
    DATE_SUB(NOW(), INTERVAL 1 HOUR)
  )
ON DUPLICATE KEY UPDATE
  task_no = VALUES(task_no),
  task_type = VALUES(task_type),
  title = VALUES(title),
  channel = VALUES(channel),
  stage_name = VALUES(stage_name),
  status = VALUES(status),
  source_type = VALUES(source_type),
  source_id = VALUES(source_id),
  customer_id = VALUES(customer_id),
  customer_name = VALUES(customer_name),
  product_name = VALUES(product_name),
  order_id = VALUES(order_id),
  invoice_id = VALUES(invoice_id),
  service_id = VALUES(service_id),
  service_no = VALUES(service_no),
  provider_type = VALUES(provider_type),
  provider_resource_id = VALUES(provider_resource_id),
  action_name = VALUES(action_name),
  operator_type = VALUES(operator_type),
  operator_name = VALUES(operator_name),
  request_payload = VALUES(request_payload),
  result_payload = VALUES(result_payload),
  message = VALUES(message),
  created_at = VALUES(created_at),
  started_at = VALUES(started_at),
  finished_at = VALUES(finished_at),
  updated_at = VALUES(updated_at);

-- -----------------------------------------------------------------------------
-- Provider 同步日志与服务资源明细
-- -----------------------------------------------------------------------------

INSERT INTO provider_sync_logs (
  id, provider_type, action, resource_type, resource_id, service_id, status, message, request_body, response_body, created_at
)
VALUES
  (
    1,
    'MOFANG_CLOUD',
    'sync-instance',
    'instance',
    'mf-cloud-1001',
    1,
    'SUCCESS',
    '实例详情已同步到本地服务与资源视图。',
    '{"mode":"seed","serviceId":1}',
    '{"instanceId":"mf-cloud-1001","status":"running"}',
    DATE_SUB(NOW(), INTERVAL 4 HOUR)
  ),
  (
    2,
    'MOFANG_CLOUD',
    'open-console',
    'instance',
    'mf-cloud-1001',
    1,
    'FAILED',
    '演示失败日志：控制台令牌过期，需要重新获取。',
    '{"serviceNo":"SRV-20260320-0001"}',
    '{"error":"token expired"}',
    DATE_SUB(NOW(), INTERVAL 90 MINUTE)
  )
ON DUPLICATE KEY UPDATE
  provider_type = VALUES(provider_type),
  action = VALUES(action),
  resource_type = VALUES(resource_type),
  resource_id = VALUES(resource_id),
  service_id = VALUES(service_id),
  status = VALUES(status),
  message = VALUES(message),
  request_body = VALUES(request_body),
  response_body = VALUES(response_body),
  created_at = VALUES(created_at);

INSERT INTO service_network_interfaces (
  id, service_id, provider_interface_id, name, mac_address, bridge, network_name, interface_name, nic_model, inbound_mbps, outbound_mbps, vpc_name, status
)
VALUES
  (
    1,
    1,
    'nic-1001',
    '主网卡',
    '52:54:00:12:34:56',
    'br0',
    '公网网络',
    'eth0',
    'virtio',
    20,
    20,
    'prod-vpc',
    'ACTIVE'
  )
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  provider_interface_id = VALUES(provider_interface_id),
  name = VALUES(name),
  mac_address = VALUES(mac_address),
  bridge = VALUES(bridge),
  network_name = VALUES(network_name),
  interface_name = VALUES(interface_name),
  nic_model = VALUES(nic_model),
  inbound_mbps = VALUES(inbound_mbps),
  outbound_mbps = VALUES(outbound_mbps),
  vpc_name = VALUES(vpc_name),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO service_ip_addresses (
  id, service_id, network_interface_id, provider_ip_id, address, version, ip_role, subnet_mask, gateway, bandwidth_mbps, is_primary, status
)
VALUES
  (1, 1, 1, 'ip-1001', '203.0.113.11', 'IPv4', 'public', '255.255.255.0', '203.0.113.1', 20, 1, 'ACTIVE'),
  (2, 1, 1, 'ip-1002', '2408:4004:10::11', 'IPv6', 'public', '64', '2408:4004:10::1', 20, 0, 'ACTIVE'),
  (3, 1, 1, 'ip-1003', '10.0.1.10', 'IPv4', 'private', '255.255.255.0', '10.0.1.1', 0, 0, 'ACTIVE')
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  network_interface_id = VALUES(network_interface_id),
  provider_ip_id = VALUES(provider_ip_id),
  address = VALUES(address),
  version = VALUES(version),
  ip_role = VALUES(ip_role),
  subnet_mask = VALUES(subnet_mask),
  gateway = VALUES(gateway),
  bandwidth_mbps = VALUES(bandwidth_mbps),
  is_primary = VALUES(is_primary),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO service_disks (
  id, service_id, provider_disk_id, name, disk_type, size_gb, device_name, driver_name, fs_type, mount_point, is_system, status
)
VALUES
  (1, 1, 'disk-1001', 'system-disk', 'system', 60, '/dev/vda', 'virtio', 'xfs', '/', 1, 'ACTIVE'),
  (2, 1, 'disk-1002', 'data-disk-1', 'data', 120, '/dev/vdb', 'virtio', 'xfs', '/data', 0, 'ACTIVE')
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  provider_disk_id = VALUES(provider_disk_id),
  name = VALUES(name),
  disk_type = VALUES(disk_type),
  size_gb = VALUES(size_gb),
  device_name = VALUES(device_name),
  driver_name = VALUES(driver_name),
  fs_type = VALUES(fs_type),
  mount_point = VALUES(mount_point),
  is_system = VALUES(is_system),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO service_snapshots (
  id, service_id, disk_id, provider_snapshot_id, name, size_gb, status
)
VALUES
  (1, 1, 2, 'snapshot-1001', 'daily-backup-001', 120, 'READY')
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  disk_id = VALUES(disk_id),
  provider_snapshot_id = VALUES(provider_snapshot_id),
  name = VALUES(name),
  size_gb = VALUES(size_gb),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO service_backups (
  id, service_id, provider_backup_id, name, size_gb, status, expires_at
)
VALUES
  (1, 1, 'backup-1001', 'weekly-backup-001', 120, 'READY', DATE_ADD(NOW(), INTERVAL 25 DAY))
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  provider_backup_id = VALUES(provider_backup_id),
  name = VALUES(name),
  size_gb = VALUES(size_gb),
  status = VALUES(status),
  expires_at = VALUES(expires_at),
  updated_at = NOW();

INSERT INTO service_vpc_networks (
  id, service_id, provider_vpc_id, name, cidr, gateway, interface_name, status
)
VALUES
  (1, 1, 'vpc-1001', 'prod-vpc', '10.0.1.0/24', '10.0.1.1', 'eth0', 'ACTIVE')
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  provider_vpc_id = VALUES(provider_vpc_id),
  name = VALUES(name),
  cidr = VALUES(cidr),
  gateway = VALUES(gateway),
  interface_name = VALUES(interface_name),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO service_security_groups (
  id, service_id, provider_security_group_id, name, status
)
VALUES
  (1, 1, 'sg-1001', 'default-cloud', 'ACTIVE')
ON DUPLICATE KEY UPDATE
  service_id = VALUES(service_id),
  provider_security_group_id = VALUES(provider_security_group_id),
  name = VALUES(name),
  status = VALUES(status),
  updated_at = NOW();

INSERT INTO service_security_group_rules (
  id, security_group_id, direction, protocol, port_range, source_cidr, action, description
)
VALUES
  (1, 1, 'ingress', 'tcp', '22', '0.0.0.0/0', 'ALLOW', '允许 SSH 登录'),
  (2, 1, 'ingress', 'tcp', '80,443', '0.0.0.0/0', 'ALLOW', '允许 Web 业务访问')
ON DUPLICATE KEY UPDATE
  security_group_id = VALUES(security_group_id),
  direction = VALUES(direction),
  protocol = VALUES(protocol),
  port_range = VALUES(port_range),
  source_cidr = VALUES(source_cidr),
  action = VALUES(action),
  description = VALUES(description),
  updated_at = NOW();
