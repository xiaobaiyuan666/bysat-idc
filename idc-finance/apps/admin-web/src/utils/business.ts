import type { LocaleCode } from "@/locales";
import { isLegacyMojibake, pickLabel } from "@/locales";

type LabelPair = { zh: string; en: string };

const billingCycleLabels: Record<string, LabelPair> = {
  monthly: { zh: "月付", en: "Monthly" },
  quarterly: { zh: "季付", en: "Quarterly" },
  semiannual: { zh: "半年付", en: "Semiannual" },
  semiannually: { zh: "半年付", en: "Semiannual" },
  annual: { zh: "年付", en: "Annual" },
  annually: { zh: "年付", en: "Annual" },
  biennially: { zh: "两年付", en: "Biennial" },
  triennially: { zh: "三年付", en: "Triennial" },
  onetime: { zh: "一次性", en: "One-time" }
};

const orderStatusLabels: Record<string, LabelPair> = {
  PENDING: { zh: "待支付", en: "Pending" },
  ACTIVE: { zh: "待交付", en: "Ready to Deliver" },
  PAID: { zh: "已支付", en: "Paid" },
  COMPLETED: { zh: "已完成", en: "Completed" },
  CANCELLED: { zh: "已取消", en: "Cancelled" }
};

const invoiceStatusLabels: Record<string, LabelPair> = {
  UNPAID: { zh: "未支付", en: "Unpaid" },
  PAID: { zh: "已支付", en: "Paid" },
  REFUNDED: { zh: "已退款", en: "Refunded" }
};

const serviceStatusLabels: Record<string, LabelPair> = {
  PENDING: { zh: "待开通", en: "Pending" },
  ACTIVE: { zh: "运行中", en: "Active" },
  SUSPENDED: { zh: "已暂停", en: "Suspended" },
  TERMINATED: { zh: "已终止", en: "Terminated" }
};

const syncStatusLabels: Record<string, LabelPair> = {
  PENDING: { zh: "待同步", en: "Pending" },
  RUNNING: { zh: "执行中", en: "Running" },
  SUCCESS: { zh: "同步成功", en: "Success" },
  FAILED: { zh: "同步失败", en: "Failed" },
  PROVISIONING: { zh: "开通中", en: "Provisioning" }
};

const providerTypeLabels: Record<string, LabelPair> = {
  LOCAL: { zh: "本地模块", en: "Local Module" },
  MOFANG_CLOUD: { zh: "魔方云", en: "Mofang Cloud" },
  ZJMF_API: { zh: "上下游财务", en: "Finance Upstream" },
  RESOURCE: { zh: "资源池", en: "Resource Pool" },
  MANUAL: { zh: "手工资源", en: "Manual Resource" },
  SERVICE_CHANGE: { zh: "服务改配", en: "Service Change" }
};

const productTypeLabels: Record<string, LabelPair> = {
  CLOUD_HOST: { zh: "云主机", en: "Cloud Host" },
  DCIMCLOUD: { zh: "云资源商品", en: "Cloud Resource" },
  DCIM: { zh: "物理机商品", en: "Bare Metal" }
};

const paymentChannelLabels: Record<string, LabelPair> = {
  OFFLINE: { zh: "线下汇款", en: "Offline Transfer" },
  MANUAL: { zh: "线下收款", en: "Manual Payment" },
  ALIPAY: { zh: "支付宝", en: "Alipay" },
  WECHAT: { zh: "微信支付", en: "WeChat Pay" },
  BALANCE: { zh: "余额抵扣", en: "Balance" },
  ONLINE: { zh: "在线支付", en: "Online Payment" },
  SYSTEM: { zh: "系统处理", en: "System" }
};

const changeOrderActionLabels: Record<string, LabelPair> = {
  "add-ipv4": { zh: "新增 IPv4", en: "Add IPv4" },
  "add-ipv6": { zh: "新增 IPv6", en: "Add IPv6" },
  "add-disk": { zh: "新增数据盘", en: "Add Disk" },
  "resize-disk": { zh: "扩容数据盘", en: "Resize Disk" },
  "create-snapshot": { zh: "创建快照", en: "Create Snapshot" },
  "create-backup": { zh: "创建备份", en: "Create Backup" }
};

const changeOrderExecutionLabels: Record<string, LabelPair> = {
  WAITING_PAYMENT: { zh: "待支付", en: "Waiting Payment" },
  PAID: { zh: "待回写", en: "Paid, Waiting Sync" },
  EXECUTING: { zh: "执行中", en: "Executing" },
  EXECUTED: { zh: "已执行", en: "Executed" },
  EXECUTE_FAILED: { zh: "执行失败", en: "Execution Failed" },
  EXECUTE_BLOCKED: { zh: "执行阻塞", en: "Execution Blocked" },
  REFUNDED: { zh: "已退款", en: "Refunded" }
};

const automationTaskTypeLabels: Record<string, LabelPair> = {
  AUTO_PROVISION: { zh: "自动开通", en: "Auto Provision" },
  SERVICE_ACTION: { zh: "服务动作", en: "Service Action" },
  PULL_SYNC_SERVICE: { zh: "服务同步", en: "Service Sync" },
  PULL_SYNC_BATCH: { zh: "批量同步", en: "Batch Sync" },
  RESOURCE_ACTION: { zh: "资源动作", en: "Resource Action" },
  INVOICE_ACTION: { zh: "账单动作", en: "Invoice Action" },
  CUSTOM: { zh: "自定义任务", en: "Custom Task" }
};

const resourceTypeLabels: Record<string, LabelPair> = {
  instance: { zh: "实例", en: "Instance" },
  ip: { zh: "IP 地址", en: "IP Address" },
  disk: { zh: "磁盘", en: "Disk" },
  snapshot: { zh: "快照", en: "Snapshot" },
  backup: { zh: "备份", en: "Backup" }
};

const lastActionLabels: Record<string, LabelPair> = {
  "provider-sync": { zh: "同步服务信息", en: "Provider Sync" },
  reboot: { zh: "重启实例", en: "Reboot" },
  reinstall: { zh: "重装系统", en: "Reinstall" },
  "reset-password": { zh: "重置密码", en: "Reset Password" },
  suspend: { zh: "暂停服务", en: "Suspend" },
  activate: { zh: "恢复运行", en: "Activate" },
  terminate: { zh: "终止服务", en: "Terminate" }
};

const auditActionLabels: Record<string, LabelPair> = {
  "invoice.receive_payment": { zh: "登记收款", en: "Receive Payment" },
  "invoice.manual_adjust": { zh: "人工调整账单", en: "Manual Invoice Adjustment" },
  "order.manual_adjust": { zh: "人工调整订单", en: "Manual Order Adjustment" },
  "service.change_order.create": { zh: "生成改配单", en: "Create Change Order" },
  "service.change_order.execute": { zh: "执行改配单", en: "Execute Change Order" },
  "service.manual_adjust": { zh: "人工调整服务", en: "Manual Service Adjustment" },
  "service.pull_sync": { zh: "服务同步", en: "Service Sync" },
  "service.activate": { zh: "恢复服务", en: "Activate Service" },
  "service.terminate": { zh: "终止服务", en: "Terminate Service" },
  "service.reboot": { zh: "重启实例", en: "Reboot Instance" },
  "service.provision.success": { zh: "自动开通完成", en: "Provision Complete" }
};

const auditDescriptionFallbacks: Record<string, LabelPair> = {
  "invoice.receive_payment": { zh: "后台登记账单收款", en: "Payment was recorded from the admin panel." },
  "invoice.manual_adjust": { zh: "管理员人工调整账单内容与状态", en: "Administrator adjusted invoice content and status." },
  "order.manual_adjust": { zh: "管理员人工调整订单内容与状态", en: "Administrator adjusted order content and status." },
  "service.change_order.create": { zh: "管理员为服务生成待支付改配单", en: "Administrator created a pending change order for the service." },
  "service.change_order.execute": { zh: "改配单支付后自动执行资源动作", en: "The paid change order triggered resource execution automatically." },
  "service.manual_adjust": { zh: "管理员人工调整服务状态、接口绑定与同步信息", en: "Administrator adjusted service status, account binding, and sync state." },
  "service.pull_sync": { zh: "从远端平台同步服务信息回本地", en: "Service data was synchronized from the remote platform." },
  "service.activate": { zh: "支付或人工操作后恢复服务", en: "Service was activated after payment or manual handling." },
  "service.terminate": { zh: "服务已被终止", en: "The service was terminated." },
  "service.reboot": { zh: "后台发起实例重启", en: "An instance reboot was triggered from the admin panel." },
  "service.provision.success": { zh: "按自动化渠道完成实例开通回写", en: "Provisioning completed and the service was written back locally." }
};

const fieldLabels: Record<string, LabelPair> = {
  status: { zh: "状态", en: "Status" },
  productName: { zh: "商品名称", en: "Product Name" },
  billingCycle: { zh: "计费周期", en: "Billing Cycle" },
  amount: { zh: "金额", en: "Amount" },
  dueAt: { zh: "到期时间", en: "Due At" },
  orderStatus: { zh: "订单状态", en: "Order Status" },
  orderId: { zh: "关联订单", en: "Order" },
  orderNo: { zh: "订单编号", en: "Order No." },
  invoiceId: { zh: "关联账单", en: "Invoice" },
  invoiceNo: { zh: "账单编号", en: "Invoice No." },
  invoiceStatus: { zh: "账单状态", en: "Invoice Status" },
  invoiceStatuses: { zh: "关联账单", en: "Invoices" },
  serviceStatuses: { zh: "关联服务", en: "Services" },
  paymentCount: { zh: "支付记录", en: "Payments" },
  refundCount: { zh: "退款记录", en: "Refunds" },
  providerType: { zh: "自动化渠道", en: "Automation Channel" },
  providerAccountId: { zh: "接口账户", en: "Provider Account" },
  providerResourceId: { zh: "远端资源 ID", en: "Remote Resource ID" },
  regionName: { zh: "地域", en: "Region" },
  ipAddress: { zh: "IP 地址", en: "IP Address" },
  nextDueAt: { zh: "下次到期", en: "Next Due At" },
  syncStatus: { zh: "同步状态", en: "Sync Status" },
  syncMessage: { zh: "同步说明", en: "Sync Note" },
  lastAction: { zh: "最近动作", en: "Last Action" },
  actionName: { zh: "改配动作", en: "Change Action" },
  serviceStatus: { zh: "服务状态", en: "Service Status" }
};

function resolve(locale: LocaleCode, labels: Record<string, LabelPair>, value: string, fallback?: string) {
  const item = labels[value];
  if (!item) return fallback ?? value;
  return pickLabel(locale, item.zh, item.en);
}

export function formatMoney(locale: LocaleCode, value: number) {
  return new Intl.NumberFormat(locale === "en-US" ? "en-US" : "zh-CN", {
    style: "currency",
    currency: locale === "en-US" ? "USD" : "CNY",
    minimumFractionDigits: 2,
    maximumFractionDigits: 2
  }).format(Number(value || 0));
}

export function billingCycleOptions(locale: LocaleCode) {
  return Object.entries(billingCycleLabels).map(([value, pair]) => ({
    value,
    label: pickLabel(locale, pair.zh, pair.en)
  }));
}

export function orderStatusOptions(locale: LocaleCode) {
  return ["PENDING", "ACTIVE", "COMPLETED", "CANCELLED"].map(value => ({
    value,
    label: formatOrderStatus(locale, value)
  }));
}

export function invoiceStatusOptions(locale: LocaleCode) {
  return ["UNPAID", "PAID", "REFUNDED"].map(value => ({
    value,
    label: formatInvoiceStatus(locale, value)
  }));
}

export function serviceStatusOptions(locale: LocaleCode) {
  return ["PENDING", "ACTIVE", "SUSPENDED", "TERMINATED"].map(value => ({
    value,
    label: formatServiceStatus(locale, value)
  }));
}

export function syncStatusOptions(locale: LocaleCode) {
  return ["PENDING", "RUNNING", "SUCCESS", "FAILED"].map(value => ({
    value,
    label: formatSyncStatus(locale, value)
  }));
}

export function providerTypeOptions(locale: LocaleCode) {
  return ["LOCAL", "MOFANG_CLOUD", "ZJMF_API", "RESOURCE", "MANUAL"].map(value => ({
    value,
    label: formatProviderType(locale, value)
  }));
}

export function formatBillingCycle(locale: LocaleCode, value: string) {
  return resolve(locale, billingCycleLabels, value);
}

export function formatOrderStatus(locale: LocaleCode, value: string) {
  return resolve(locale, orderStatusLabels, value);
}

export function formatInvoiceStatus(locale: LocaleCode, value: string) {
  return resolve(locale, invoiceStatusLabels, value);
}

export function formatServiceStatus(locale: LocaleCode, value: string) {
  return resolve(locale, serviceStatusLabels, value);
}

export function formatSyncStatus(locale: LocaleCode, value: string) {
  return resolve(locale, syncStatusLabels, value);
}

export function formatProviderType(locale: LocaleCode, value: string) {
  return resolve(locale, providerTypeLabels, value);
}

export function formatProductType(locale: LocaleCode, value: string) {
  return resolve(locale, productTypeLabels, value);
}

export function formatPaymentChannel(locale: LocaleCode, value: string) {
  return resolve(locale, paymentChannelLabels, value);
}

export function formatChangeOrderAction(locale: LocaleCode, value: string) {
  return resolve(locale, changeOrderActionLabels, value);
}

export function formatChangeOrderExecution(locale: LocaleCode, value: string) {
  return resolve(locale, changeOrderExecutionLabels, value);
}

export function formatResourceType(locale: LocaleCode, value: string) {
  return resolve(locale, resourceTypeLabels, value);
}

export function formatLastAction(locale: LocaleCode, value: string) {
  return resolve(locale, lastActionLabels, value);
}

export function formatAutomationTaskStatus(locale: LocaleCode, value: string) {
  return resolve(
    locale,
    {
      PENDING: { zh: "待执行", en: "Pending" },
      RUNNING: { zh: "执行中", en: "Running" },
      SUCCESS: { zh: "成功", en: "Success" },
      FAILED: { zh: "失败", en: "Failed" },
      BLOCKED: { zh: "阻塞", en: "Blocked" }
    },
    value
  );
}

export function formatAutomationTaskType(locale: LocaleCode, value: string) {
  return resolve(locale, automationTaskTypeLabels, value);
}

export function formatAutomationTaskTitle(locale: LocaleCode, taskType: string, title: string, actionName = "") {
  if (taskType === "SERVICE_ACTION" && actionName) {
    return pickLabel(locale, `服务动作：${formatLastAction(locale, actionName)}`, `Service Action: ${formatLastAction(locale, actionName)}`);
  }

  if (taskType === "RESOURCE_ACTION" && actionName) {
    return pickLabel(locale, `资源动作：${formatChangeOrderAction(locale, actionName)}`, `Resource Action: ${formatChangeOrderAction(locale, actionName)}`);
  }

  if (taskType === "SERVICE_ACTION") {
    return pickLabel(locale, "服务动作执行", "Service Action");
  }

  if (taskType === "RESOURCE_ACTION") {
    return pickLabel(locale, "资源动作执行", "Resource Action");
  }

  if (taskType === "INVOICE_ACTION") {
    return pickLabel(locale, "账单收款处理", "Invoice Action");
  }

  if (taskType === "PULL_SYNC_SERVICE") {
    return pickLabel(locale, "魔方云单服务同步", "Mofang Service Sync");
  }

  if (taskType === "PULL_SYNC_BATCH") {
    return pickLabel(locale, "魔方云批量同步", "Mofang Batch Sync");
  }

  if (title && !isLegacyMojibake(title)) {
    return title;
  }

  return formatAutomationTaskType(locale, taskType);
}

export function formatSyncLogAction(locale: LocaleCode, action: string) {
  const actionLabels: Record<string, LabelPair> = {
    pull_sync_service: { zh: "单服务同步", en: "Sync Service" },
    pull_sync_batch: { zh: "批量同步", en: "Batch Sync" },
    provision_create: { zh: "自动开通", en: "Provision Create" },
    "provider-sync": { zh: "同步服务信息", en: "Provider Sync" },
    reboot: { zh: "重启实例", en: "Reboot" }
  };

  if (actionLabels[action]) {
    return pickLabel(locale, actionLabels[action].zh, actionLabels[action].en);
  }

  if (action.startsWith("resource_")) {
    return pickLabel(locale, `资源动作：${formatChangeOrderAction(locale, action.replace(/^resource_/, ""))}`, `Resource: ${formatChangeOrderAction(locale, action.replace(/^resource_/, ""))}`);
  }

  return action;
}

export function formatAuditAction(locale: LocaleCode, action: string) {
  return resolve(locale, auditActionLabels, action, action);
}

export function formatAuditDescription(locale: LocaleCode, action: string, description?: string) {
  if (description && !isLegacyMojibake(description)) {
    return description;
  }

  const fallback = auditDescriptionFallbacks[action];
  if (fallback) {
    return pickLabel(locale, fallback.zh, fallback.en);
  }

  return formatAuditAction(locale, action);
}

export function formatAuditReason(reason: unknown) {
  const value = String(reason ?? "").trim();
  if (!value || isLegacyMojibake(value)) return "-";
  return value;
}

export function fieldLabel(locale: LocaleCode, key: string): string {
  return resolve(locale, fieldLabels, key, key);
}

export function formatFieldValue(locale: LocaleCode, key: string, value: unknown, action = ""): string {
  if (value === undefined || value === null || value === "") return "-";
  if (typeof value === "boolean") {
    return pickLabel(locale, value ? "是" : "否", value ? "Yes" : "No");
  }
  if (Array.isArray(value)) {
    return value.length ? value.map(item => formatFieldValue(locale, key, item, action)).join(" / ") : "-";
  }
  if (typeof value === "number") {
    if (key === "amount") return formatMoney(locale, value);
    return String(value);
  }
  if (typeof value === "object") {
    try {
      return JSON.stringify(value);
    } catch {
      return String(value);
    }
  }

  const text = String(value);
  switch (key) {
    case "billingCycle":
      return formatBillingCycle(locale, text);
    case "orderStatus":
      return formatOrderStatus(locale, text);
    case "invoiceStatus":
      return formatInvoiceStatus(locale, text);
    case "serviceStatus":
      return formatServiceStatus(locale, text);
    case "status":
      if (action.startsWith("service.")) return formatServiceStatus(locale, text);
      if (action.startsWith("invoice.")) return formatInvoiceStatus(locale, text);
      if (action.startsWith("order.")) return formatOrderStatus(locale, text);
      return text;
    case "syncStatus":
      return formatSyncStatus(locale, text);
    case "providerType":
      return formatProviderType(locale, text);
    case "lastAction":
      return formatLastAction(locale, text);
    case "actionName":
      return formatChangeOrderAction(locale, text);
    case "amount":
      return formatMoney(locale, Number(text));
    default:
      return isLegacyMojibake(text) ? "-" : text;
  }
}
