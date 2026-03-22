export const roleMap: Record<string, string> = {
  SUPER_ADMIN: "超级管理员",
  FINANCE: "财务",
  SUPPORT: "客服",
  OPERATIONS: "运营",
};

export const portalRoleMap: Record<string, string> = {
  OWNER: "主账号",
  TECH: "技术联系人",
  BILLING: "财务联系人",
  READONLY: "只读成员",
};

export const customerTypeMap: Record<string, string> = {
  PERSONAL: "个人",
  COMPANY: "企业",
  RESELLER: "代理商",
};

export const customerStatusMap: Record<string, string> = {
  ACTIVE: "正常",
  SUSPENDED: "暂停",
  OVERDUE: "逾期",
  ARCHIVED: "归档",
};

export const orderStatusMap: Record<string, string> = {
  PENDING: "待支付",
  PAID: "已付款",
  PROVISIONING: "开通中",
  ACTIVE: "已生效",
  CANCELLED: "已取消",
  REFUNDED: "已退款",
};

export const orderSourceMap: Record<string, string> = {
  portal: "前台门户",
  sales: "销售录入",
  admin: "后台录入",
  "billing-engine": "计费引擎",
  "mofang-sync": "魔方同步",
};

export const productStatusMap: Record<string, string> = {
  ACTIVE: "销售中",
  DRAFT: "草稿",
  DISABLED: "已下架",
};

export const serviceStatusMap: Record<string, string> = {
  PENDING: "待开通",
  PROVISIONING: "开通中",
  ACTIVE: "运行中",
  SUSPENDED: "已暂停",
  OVERDUE: "已逾期",
  TERMINATED: "已终止",
  EXPIRED: "已到期",
  FAILED: "异常",
};

export const invoiceTypeMap: Record<string, string> = {
  ORDER: "订单账单",
  RENEWAL: "续费账单",
  MANUAL: "手工账单",
  CREDIT: "余额账单",
};

export const invoiceStatusMap: Record<string, string> = {
  DRAFT: "草稿",
  ISSUED: "待支付",
  PARTIAL: "部分支付",
  PAID: "已支付",
  OVERDUE: "已逾期",
  VOID: "已作废",
};

export const paymentMethodMap: Record<string, string> = {
  BALANCE: "余额",
  ALIPAY: "支付宝",
  WECHAT: "微信支付",
  BANK_TRANSFER: "银行转账",
  OFFLINE: "线下收款",
  PAYPAL: "PayPal",
  STRIPE: "Stripe",
  OTHER: "其他",
};

export const paymentStatusMap: Record<string, string> = {
  PENDING: "待确认",
  SUCCESS: "成功",
  FAILED: "失败",
  REFUNDED: "已退款",
};

export const refundStatusMap: Record<string, string> = {
  PENDING: "待处理",
  SUCCESS: "已退款",
  FAILED: "失败",
  CANCELED: "已取消",
};

export const notificationChannelMap: Record<string, string> = {
  SYSTEM: "站内通知",
  EMAIL: "邮件",
  SMS: "短信",
  WEBHOOK: "Webhook",
};

export const notificationPriorityMap: Record<string, string> = {
  LOW: "低",
  NORMAL: "普通",
  HIGH: "高",
  URGENT: "紧急",
};

export const notificationStatusMap: Record<string, string> = {
  PENDING: "待发送",
  SENT: "已送达",
  FAILED: "发送失败",
  CANCELED: "已取消",
};

export const callbackStatusMap: Record<string, string> = {
  PENDING: "待处理",
  SUCCESS: "成功",
  FAILED: "失败",
  REFUNDED: "已退款",
};

export const ticketStatusMap: Record<string, string> = {
  OPEN: "待处理",
  PROCESSING: "处理中",
  WAITING_CUSTOMER: "待客户回复",
  RESOLVED: "已解决",
  CLOSED: "已关闭",
};

export const ticketPriorityMap: Record<string, string> = {
  LOW: "低",
  NORMAL: "普通",
  HIGH: "高",
  URGENT: "紧急",
};

export const cycleMap: Record<string, string> = {
  MONTHLY: "月付",
  QUARTERLY: "季付",
  SEMI_ANNUALLY: "半年付",
  ANNUALLY: "年付",
  BIENNIALLY: "两年付",
  TRIENNIALLY: "三年付",
  ONETIME: "一次性",
};

export const productCategoryMap: Record<string, string> = {
  CLOUD_SERVER: "云服务器",
  BARE_METAL: "裸金属",
  STORAGE: "存储",
  CDN: "CDN",
  SECURITY: "安全",
  NETWORK: "网络",
  LICENSE: "许可",
  DOMAIN: "域名",
  OTHER: "其他",
};

export const providerTypeMap: Record<string, string> = {
  MOFANG_CLOUD: "魔方云",
  MANUAL: "人工开通",
};

export const creditTransactionTypeMap: Record<string, string> = {
  RECHARGE: "充值",
  CONSUME: "消费",
  REFUND: "退款",
  ADJUSTMENT: "调整",
  AUTO_RENEW: "自动续费",
};

export const providerSyncStatusMap: Record<string, string> = {
  SUCCESS: "成功",
  FAILED: "失败",
  PENDING: "排队中",
};

export const resourceStatusMap: Record<string, string> = {
  ACTIVE: "正常",
  ATTACHED: "已挂载",
  DETACHED: "未挂载",
  SUSPENDED: "暂停",
  READY: "可用",
  ARCHIVED: "归档",
  EXPIRED: "已过期",
  FAILED: "失败",
  PENDING: "处理中",
};

export const providerActionMap: Record<string, string> = {
  sync: "同步",
  provision: "开通",
  activate: "开机",
  powerOn: "开机",
  powerOff: "关机",
  reboot: "重启",
  hardReboot: "强制重启",
  hardPowerOff: "强制关机",
  reinstall: "重装系统",
  resetPassword: "重置密码",
  suspend: "暂停",
  unsuspend: "解除暂停",
  renew: "续费",
  terminate: "终止",
  backup: "备份",
  createSnapshot: "创建快照",
  createBackup: "创建备份",
  restore: "恢复",
  deleteSnapshot: "删除快照",
  deleteBackup: "删除备份",
  detach: "卸载磁盘",
  attach: "挂载磁盘",
  setBoot: "设为系统盘",
  addRule: "新增规则",
  deleteRule: "删除规则",
  deleteGroup: "删除安全组",
  getVnc: "VNC 控制台",
  rescueStart: "进入救援",
  rescueStop: "退出救援",
  lock: "锁定",
  unlock: "解锁",
};

export const billingJobStatusMap: Record<string, string> = {
  PENDING: "待执行",
  SUCCESS: "成功",
  FAILED: "失败",
  SKIPPED: "跳过",
};

export const billingJobTypeMap: Record<string, string> = {
  CREATE_RENEWAL_INVOICE: "生成续费账单",
  AUTO_RENEW_BALANCE: "余额自动续费",
  MARK_OVERDUE: "标记逾期",
  AUTO_SUSPEND: "自动暂停",
  AUTO_TERMINATE: "自动终止",
  RESOURCE_SYNC: "资源同步",
};

export const auditModuleMap: Record<string, string> = {
  auth: "登录认证",
  architecture: "云架构中心",
  billing: "计费中心",
  customer: "客户管理",
  invoice: "账单管理",
  notification: "通知中心",
  order: "订单管理",
  product: "产品管理",
  resource: "资源中心",
  payment: "收款中心",
  service: "服务管理",
  ticket: "工单中心",
};

export const auditActionMap: Record<string, string> = {
  login: "登录",
  logout: "退出登录",
  create: "创建",
  update: "更新",
  delete: "删除",
  reply: "回复",
  refund: "退款",
  "provider-sync": "魔方云同步",
  sync: "同步",
  activate: "开通",
  suspend: "暂停",
  unsuspend: "解除暂停",
  renew: "续费",
  terminate: "终止",
  powerOn: "开机",
  powerOff: "关机",
  reboot: "重启",
  hardReboot: "强制重启",
  hardPowerOff: "强制关机",
  reinstall: "重装系统",
  resetPassword: "重置密码",
  getVnc: "获取 VNC",
  rescueStart: "进入救援",
  rescueStop: "退出救援",
  lock: "锁定",
  unlock: "解除锁定",
};

export const auditTargetTypeMap: Record<string, string> = {
  auth: "认证",
  architecture: "云架构",
  billing: "计费任务",
  customer: "客户",
  invoice: "账单",
  notification: "通知",
  order: "订单",
  payment: "支付",
  product: "产品",
  resource: "资源",
  service: "服务",
  ticket: "工单",
};

export function getLabel(map: Record<string, string>, value?: string | null) {
  if (!value) {
    return "-";
  }

  return map[value] ?? value;
}

export function getStatusTagType(status?: string | null) {
  switch (status) {
    case "ACTIVE":
    case "SUCCESS":
    case "PAID":
    case "SENT":
    case "RESOLVED":
    case "READY":
      return "success";
    case "PENDING":
    case "PROVISIONING":
    case "PARTIAL":
    case "PROCESSING":
      return "warning";
    case "SUSPENDED":
    case "OVERDUE":
    case "EXPIRED":
    case "FAILED":
    case "VOID":
    case "CANCELLED":
    case "REFUNDED":
    case "CLOSED":
    case "TERMINATED":
      return "danger";
    default:
      return "info";
  }
}
