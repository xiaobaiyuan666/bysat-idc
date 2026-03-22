import { addMonths, format } from "date-fns";
import { zhCN } from "date-fns/locale";

import {
  BILLING_CYCLE_LABELS,
  CUSTOMER_STATUS_LABELS,
  CYCLE_MONTHS,
  INVOICE_STATUS_LABELS,
  NOTIFICATION_CHANNEL_LABELS,
  NOTIFICATION_PRIORITY_LABELS,
  NOTIFICATION_STATUS_LABELS,
  ORDER_STATUS_LABELS,
  PAYMENT_METHOD_LABELS,
  PAYMENT_STATUS_LABELS,
  PRODUCT_CATEGORY_LABELS,
  REFUND_STATUS_LABELS,
  SERVICE_STATUS_LABELS,
  TICKET_PRIORITY_LABELS,
  TICKET_STATUS_LABELS,
} from "@/lib/constants";

export function formatCurrency(value: number) {
  return new Intl.NumberFormat("zh-CN", {
    style: "currency",
    currency: "CNY",
    minimumFractionDigits: 2,
  }).format(value / 100);
}

export function formatDate(value?: Date | string | null) {
  if (!value) {
    return "未设置";
  }

  return format(new Date(value), "yyyy-MM-dd", { locale: zhCN });
}

export function formatDateTime(value?: Date | string | null) {
  if (!value) {
    return "未设置";
  }

  return format(new Date(value), "yyyy-MM-dd HH:mm", { locale: zhCN });
}

export function formatPercent(value: number) {
  return `${(value / 100).toFixed(2)}%`;
}

export function parseMoneyToCent(value: FormDataEntryValue | null) {
  const raw = typeof value === "string" ? value.trim() : "";
  if (!raw) {
    return 0;
  }

  const parsed = Number(raw);
  if (Number.isNaN(parsed)) {
    return 0;
  }

  return Math.round(parsed * 100);
}

export function addCycle(date: Date, cycle: string) {
  const months = CYCLE_MONTHS[cycle] ?? 1;

  if (!months) {
    return date;
  }

  return addMonths(date, months);
}

export function labelByMap(value: string, map: Record<string, string>) {
  return map[value] ?? value;
}

export function cycleLabel(value: string) {
  return labelByMap(value, BILLING_CYCLE_LABELS);
}

export function serviceStatusLabel(value: string) {
  return labelByMap(value, SERVICE_STATUS_LABELS);
}

export function orderStatusLabel(value: string) {
  return labelByMap(value, ORDER_STATUS_LABELS);
}

export function invoiceStatusLabel(value: string) {
  return labelByMap(value, INVOICE_STATUS_LABELS);
}

export function customerStatusLabel(value: string) {
  return labelByMap(value, CUSTOMER_STATUS_LABELS);
}

export function productCategoryLabel(value: string) {
  return labelByMap(value, PRODUCT_CATEGORY_LABELS);
}

export function paymentMethodLabel(value: string) {
  return labelByMap(value, PAYMENT_METHOD_LABELS);
}

export function paymentStatusLabel(value: string) {
  return labelByMap(value, PAYMENT_STATUS_LABELS);
}

export function refundStatusLabel(value: string) {
  return labelByMap(value, REFUND_STATUS_LABELS);
}

export function notificationChannelLabel(value: string) {
  return labelByMap(value, NOTIFICATION_CHANNEL_LABELS);
}

export function notificationPriorityLabel(value: string) {
  return labelByMap(value, NOTIFICATION_PRIORITY_LABELS);
}

export function notificationStatusLabel(value: string) {
  return labelByMap(value, NOTIFICATION_STATUS_LABELS);
}

export function ticketStatusLabel(value: string) {
  return labelByMap(value, TICKET_STATUS_LABELS);
}

export function ticketPriorityLabel(value: string) {
  return labelByMap(value, TICKET_PRIORITY_LABELS);
}
