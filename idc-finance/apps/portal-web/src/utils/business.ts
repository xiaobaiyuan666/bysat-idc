import { isLegacyMojibake, pickLabel, type LocaleCode } from "@/locales";
import type { PortalServiceConfig } from "@/api/portal";

function valueOrDash(value: string | undefined) {
  return value && value.trim() ? value : "-";
}

export function formatPortalMoney(locale: LocaleCode, value: number | string | undefined) {
  const numeric = typeof value === "number" ? value : Number.parseFloat(String(value ?? 0));
  const prefix = locale === "en-US" ? "$" : "¥";
  return Number.isNaN(numeric) ? `${prefix}0.00` : `${prefix}${numeric.toFixed(2)}`;
}

export function formatPortalBillingCycle(locale: LocaleCode, cycle: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    monthly: ["月付", "Monthly"],
    quarterly: ["季付", "Quarterly"],
    semiannual: ["半年付", "Semiannual"],
    annual: ["年付", "Annual"],
    annually: ["年付", "Annual"]
  };
  const pair = mapping[String(cycle ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(cycle);
}

export function formatPortalProductType(locale: LocaleCode, type: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    CLOUD: ["云主机", "Cloud Server"],
    CLOUD_HOST: ["云主机", "Cloud Server"],
    BANDWIDTH: ["带宽产品", "Bandwidth"],
    COLOCATION: ["机柜托管", "Colocation"],
    HOST: ["物理服务器", "Dedicated Server"],
    NETWORK: ["网络产品", "Network Product"]
  };
  const pair = mapping[String(type ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(type);
}

export function formatPortalOrderStatus(locale: LocaleCode, status: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    PENDING: ["待支付", "Pending"],
    ACTIVE: ["已生效", "Active"],
    COMPLETED: ["已完成", "Completed"],
    CANCELLED: ["已取消", "Cancelled"]
  };
  const pair = mapping[String(status ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(status);
}

export function formatPortalInvoiceStatus(locale: LocaleCode, status: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    UNPAID: ["未支付", "Unpaid"],
    PAID: ["已支付", "Paid"],
    REFUNDED: ["已退款", "Refunded"],
    VOID: ["已作废", "Voided"]
  };
  const pair = mapping[String(status ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(status);
}

export function formatPortalPaymentChannel(locale: LocaleCode, channel: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    OFFLINE: ["线下汇款", "Offline Transfer"],
    MANUAL: ["线下收款", "Manual Payment"],
    ALIPAY: ["支付宝", "Alipay"],
    WECHAT: ["微信支付", "WeChat Pay"],
    BALANCE: ["余额抵扣", "Balance"],
    ONLINE: ["在线支付", "Online Payment"],
    SYSTEM: ["系统处理", "System"]
  };
  const pair = mapping[String(channel ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(channel);
}

export function formatPortalServiceStatus(locale: LocaleCode, status: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    PENDING: ["待开通", "Pending"],
    ACTIVE: ["运行中", "Active"],
    SUSPENDED: ["已暂停", "Suspended"],
    TERMINATED: ["已终止", "Terminated"]
  };
  const pair = mapping[String(status ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(status);
}

export function formatPortalTicketStatus(locale: LocaleCode, status: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    OPEN: ["处理中", "Open"],
    PROCESSING: ["处理中", "Processing"],
    CLOSED: ["已关闭", "Closed"]
  };
  const pair = mapping[String(status ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(status);
}

export function formatPortalIdentityStatus(locale: LocaleCode, status: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    PENDING: ["待审核", "Pending"],
    APPROVED: ["已通过", "Approved"],
    REJECTED: ["已驳回", "Rejected"]
  };
  const pair = mapping[String(status ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(status);
}

export function formatPortalCustomerStatus(locale: LocaleCode, status: string | undefined) {
  const mapping: Record<string, [string, string]> = {
    ACTIVE: ["正常", "Active"],
    DISABLED: ["停用", "Disabled"]
  };
  const pair = mapping[String(status ?? "")];
  return pair ? pickLabel(locale, pair[0], pair[1]) : valueOrDash(status);
}

export function portalTagTypeByStatus(status: string | undefined): "success" | "warning" | "danger" | "info" {
  switch (status) {
    case "ACTIVE":
    case "PAID":
    case "APPROVED":
    case "CLOSED":
    case "COMPLETED":
      return "success";
    case "PENDING":
    case "UNPAID":
    case "OPEN":
    case "PROCESSING":
      return "warning";
    case "REJECTED":
    case "DISABLED":
    case "REFUNDED":
    case "TERMINATED":
    case "CANCELLED":
      return "danger";
    default:
      return "info";
  }
}

function formatPortalBoolean(locale: LocaleCode, value: string) {
  const normalized = value.trim().toLowerCase();
  if (["true", "yes", "on", "enable", "enabled", "open", "opened", "1"].includes(normalized)) {
    return pickLabel(locale, "启用", "Enabled");
  }
  if (["false", "no", "off", "disable", "disabled", "close", "closed", "0"].includes(normalized)) {
    return pickLabel(locale, "不启用", "Disabled");
  }
  return value;
}

function formatPortalConfigValue(locale: LocaleCode, item: PortalServiceConfig) {
  const rawValue = String(item.value ?? "").trim();
  const rawLabel = String(item.valueLabel ?? "").trim();
  const code = String(item.code ?? "").trim().toLowerCase();
  const label = rawLabel && !isLegacyMojibake(rawLabel) && !rawLabel.includes("?") ? rawLabel : rawValue;

  if (!label) {
    return "-";
  }

  if (["cpu"].includes(code) && /^\d+(\.\d+)?$/.test(rawValue)) {
    return pickLabel(locale, `${rawValue} 核`, `${rawValue} vCPU`);
  }

  if (["memory", "system_disk_size", "data_disk_size"].includes(code) && /^\d+(\.\d+)?$/.test(rawValue)) {
    return `${rawValue} GB`;
  }

  if (["bw", "bandwidth"].includes(code) && /^\d+(\.\d+)?$/.test(rawValue)) {
    return `${rawValue} Mbps`;
  }

  if (["ip_num", "ipv4_num", "backup_num", "snap_num"].includes(code) && /^\d+$/.test(rawValue)) {
    return pickLabel(locale, `${rawValue} 个`, `${rawValue} pcs`);
  }

  if (["backup", "snapshot", "nat_web_limit", "nat_acl_limit"].includes(code)) {
    return formatPortalBoolean(locale, label);
  }

  if (code === "network_type") {
    if (rawValue === "classic" || rawValue === "1") {
      return pickLabel(locale, "经典网络", "Classic Network");
    }
    if (rawValue === "vpc" || rawValue === "2") {
      return "VPC";
    }
  }

  return label;
}

export function formatPortalConfiguration(items: PortalServiceConfig[] | undefined, locale: LocaleCode = "zh-CN") {
  if (!items?.length) {
    return "-";
  }

  return items
    .map(item => formatPortalConfigValue(locale, item))
    .filter(Boolean)
    .join(" / ");
}
