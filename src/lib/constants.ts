export const APP_NAME = "IDC 云业务管理系统";

export const NAV_ITEMS = [
  { href: "/dashboard", label: "财务总览" },
  { href: "/customers", label: "客户管理" },
  { href: "/products", label: "产品目录" },
  { href: "/orders", label: "订单中心" },
  { href: "/services", label: "服务实例" },
  { href: "/invoices", label: "账单发票" },
  { href: "/payments", label: "收款记录" },
  { href: "/tickets", label: "工单支持" },
] as const;

export const PORTAL_NAV_ITEMS = [
  { href: "/portal", label: "控制台" },
  { href: "/portal/store", label: "云产品" },
  { href: "/portal/orders", label: "订单" },
  { href: "/portal/services", label: "服务" },
  { href: "/portal/invoices", label: "账单" },
  { href: "/portal/wallet", label: "余额" },
  { href: "/portal/notifications", label: "通知" },
  { href: "/portal/tickets", label: "工单" },
] as const;

export const CUSTOMER_TYPE_LABELS: Record<string, string> = {
  PERSONAL: "个人",
  COMPANY: "企业",
  RESELLER: "代理商",
};

export const CUSTOMER_STATUS_LABELS: Record<string, string> = {
  ACTIVE: "正常",
  SUSPENDED: "暂停",
  OVERDUE: "逾期",
  ARCHIVED: "归档",
};

export const PRODUCT_CATEGORY_LABELS: Record<string, string> = {
  CLOUD_SERVER: "云服务器",
  BARE_METAL: "裸金属",
  STORAGE: "存储",
  CDN: "CDN",
  SECURITY: "安全产品",
  NETWORK: "网络资源",
  LICENSE: "软件许可",
  DOMAIN: "域名 / SSL",
  OTHER: "其他",
};

export const BILLING_CYCLE_LABELS: Record<string, string> = {
  MONTHLY: "月付",
  QUARTERLY: "季付",
  SEMI_ANNUALLY: "半年付",
  ANNUALLY: "年付",
  BIENNIALLY: "两年付",
  TRIENNIALLY: "三年付",
  ONETIME: "一次性",
};

export const ORDER_STATUS_LABELS: Record<string, string> = {
  PENDING: "待支付",
  PAID: "已付款",
  PROVISIONING: "开通中",
  ACTIVE: "已生效",
  CANCELLED: "已取消",
  REFUNDED: "已退款",
};

export const SERVICE_STATUS_LABELS: Record<string, string> = {
  PENDING: "待开通",
  PROVISIONING: "开通中",
  ACTIVE: "运行中",
  SUSPENDED: "已暂停",
  OVERDUE: "已逾期",
  TERMINATED: "已终止",
  EXPIRED: "已到期",
  FAILED: "开通失败",
};

export const INVOICE_STATUS_LABELS: Record<string, string> = {
  DRAFT: "草稿",
  ISSUED: "待支付",
  PARTIAL: "部分支付",
  PAID: "已支付",
  OVERDUE: "已逾期",
  VOID: "已作废",
};

export const PAYMENT_METHOD_LABELS: Record<string, string> = {
  BALANCE: "余额",
  ALIPAY: "支付宝",
  WECHAT: "微信支付",
  BANK_TRANSFER: "银行转账",
  OFFLINE: "线下打款",
  PAYPAL: "PayPal",
  STRIPE: "Stripe",
  OTHER: "其他",
};

export const PAYMENT_STATUS_LABELS: Record<string, string> = {
  PENDING: "待确认",
  SUCCESS: "成功",
  FAILED: "失败",
  REFUNDED: "已退款",
};

export const REFUND_STATUS_LABELS: Record<string, string> = {
  PENDING: "待处理",
  SUCCESS: "已退款",
  FAILED: "失败",
  CANCELED: "已取消",
};

export const NOTIFICATION_CHANNEL_LABELS: Record<string, string> = {
  SYSTEM: "站内通知",
  EMAIL: "邮件",
  SMS: "短信",
  WEBHOOK: "Webhook",
};

export const NOTIFICATION_PRIORITY_LABELS: Record<string, string> = {
  LOW: "低",
  NORMAL: "普通",
  HIGH: "高",
  URGENT: "紧急",
};

export const NOTIFICATION_STATUS_LABELS: Record<string, string> = {
  PENDING: "待投递",
  SENT: "已送达",
  FAILED: "发送失败",
  CANCELED: "已取消",
};

export const TICKET_STATUS_LABELS: Record<string, string> = {
  OPEN: "待处理",
  PROCESSING: "处理中",
  WAITING_CUSTOMER: "待客户回复",
  RESOLVED: "已解决",
  CLOSED: "已关闭",
};

export const TICKET_PRIORITY_LABELS: Record<string, string> = {
  LOW: "低",
  NORMAL: "中",
  HIGH: "高",
  URGENT: "紧急",
};

export const ROLE_LABELS: Record<string, string> = {
  SUPER_ADMIN: "超级管理员",
  FINANCE: "财务",
  SUPPORT: "客服",
  OPERATIONS: "运维",
};

export const PORTAL_ROLE_LABELS: Record<string, string> = {
  OWNER: "主账户",
  BILLING: "财务联系人",
  TECH: "技术联系人",
  FINANCE: "财务联系人",
  READONLY: "只读成员",
};

export const CYCLE_MONTHS: Record<string, number> = {
  MONTHLY: 1,
  QUARTERLY: 3,
  SEMI_ANNUALLY: 6,
  ANNUALLY: 12,
  BIENNIALLY: 24,
  TRIENNIALLY: 36,
  ONETIME: 0,
};
