UPDATE account_transactions
SET
  summary = '线下充值演示资金',
  remark = '初始演示充值流水',
  operator_name = '系统管理员'
WHERE id = 1;

UPDATE account_transactions
SET
  summary = '演示余额手工扣款',
  remark = '用于展示资金台账明细',
  operator_name = '系统管理员'
WHERE id = 2;

UPDATE account_transactions
SET
  summary = '开通客户授信额度',
  remark = '演示授信额度调整',
  operator_name = '系统管理员'
WHERE id = 3;

UPDATE account_transactions
SET
  summary = '线上充值演示资金',
  remark = '第二个客户的钱包初始余额',
  operator_name = '云脉互联'
WHERE id = 4;
