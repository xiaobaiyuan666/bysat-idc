import { requestJson } from "./http";

export interface ApiEnvelope<T> {
  code: string;
  message: string;
  requestId?: string;
  data: T;
}

export interface ListResponse<T> {
  items: T[];
  total: number;
}

export interface PriceOption {
  cycleCode: string;
  cycleName: string;
  price: number;
  setupFee: number;
}

export interface PortalProductConfigChoice {
  value: string;
  label: string;
  priceDelta: number;
}

export interface PortalProductConfigOption {
  code: string;
  name: string;
  inputType: string;
  required: boolean;
  defaultValue: string;
  description: string;
  choices: PortalProductConfigChoice[];
}

export interface PortalProduct {
  id: number;
  productNo: string;
  groupName: string;
  name: string;
  description: string;
  productType: string;
  status: string;
  pricing: PriceOption[];
  configOptions: PortalProductConfigOption[];
}

export interface PortalServiceConfigSelection {
  code: string;
  name: string;
  value: string;
  valueLabel: string;
  priceDelta: number;
}

export type PortalServiceConfig = PortalServiceConfigSelection;

export interface PortalService {
  id: number;
  serviceNo: string;
  customerId: number;
  productName: string;
  status: string;
  nextDueAt: string;
  createdAt: string;
  configuration?: PortalServiceConfigSelection[];
}

export interface PortalOrder {
  id: number;
  orderNo: string;
  customerId: number;
  customerName: string;
  productId: number;
  productName: string;
  productType: string;
  billingCycle: string;
  amount: number;
  status: string;
  createdAt: string;
  configuration?: PortalServiceConfigSelection[];
}

export interface PortalInvoice {
  id: number;
  invoiceNo: string;
  orderId: number;
  orderNo: string;
  customerId: number;
  productName: string;
  totalAmount: number;
  status: string;
  dueAt: string;
  paidAt?: string;
  billingCycle: string;
}

export interface PortalPayment {
  id: number;
  paymentNo: string;
  invoiceId: number;
  orderId: number;
  customerId: number;
  channel: string;
  tradeNo: string;
  amount: number;
  source: string;
  status: string;
  operator: string;
  paidAt: string;
}

export interface PortalTicket {
  no: string;
  title: string;
  status: string;
  updatedAt: string;
}

export interface PortalWallet {
  balance: string;
  creditLimit: string;
}

export interface PortalStat {
  label: string;
  value: string;
}

export interface PortalDashboard {
  stats: PortalStat[];
  services: PortalService[];
  orders: PortalOrder[];
  invoices: PortalInvoice[];
  tickets: PortalTicket[];
  wallet: PortalWallet;
}

export interface PortalContact {
  id: number;
  name: string;
  email: string;
  mobile: string;
  roleName: string;
  isPrimary: boolean;
}

export interface PortalIdentity {
  id: number;
  identityType: string;
  verifyStatus: string;
  subjectName: string;
  certNo: string;
  countryCode: string;
  reviewRemark: string;
  reviewedAt: string;
}

export interface PortalCustomer {
  id: number;
  customerNo: string;
  name: string;
  email: string;
  mobile: string;
  type: string;
  status: string;
  groupName: string;
  levelName: string;
  salesOwner: string;
  remarks: string;
  contacts: PortalContact[];
  identity?: PortalIdentity;
}

export interface PortalAccount {
  customer: PortalCustomer;
  primaryContact?: PortalContact;
  wallet: PortalWallet;
}

export interface CheckoutSelection {
  code: string;
  value: string;
}

export interface CheckoutRequest {
  productId: number;
  cycleCode: string;
  selectedOptions: CheckoutSelection[];
}

export interface CheckoutResult {
  order: PortalOrder;
  invoice: PortalInvoice;
}

export interface PayInvoiceResult {
  invoice: PortalInvoice;
  service?: PortalService;
  payment?: PortalPayment;
}

async function unwrap<T>(promise: Promise<ApiEnvelope<T>>): Promise<T> {
  const result = await promise;
  return result.data;
}

export function loadDashboard() {
  return unwrap(requestJson<ApiEnvelope<PortalDashboard>>("/dashboard"));
}

export async function loadStoreProducts() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalProduct>>>("/products");
  return result.data.items;
}

export async function loadServices() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalService>>>("/services");
  return result.data.items;
}

export async function loadOrders() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalOrder>>>("/orders");
  return result.data.items;
}

export async function loadInvoices() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalInvoice>>>("/invoices");
  return result.data.items;
}

export async function loadTickets() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalTicket>>>("/tickets");
  return result.data.items;
}

export function loadAccount() {
  return unwrap(requestJson<ApiEnvelope<PortalAccount>>("/account"));
}

export function checkoutProduct(payload: CheckoutRequest) {
  return unwrap(
    requestJson<ApiEnvelope<CheckoutResult>>("/orders/checkout", {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export function payInvoice(invoiceId: number) {
  return unwrap(
    requestJson<ApiEnvelope<PayInvoiceResult>>(`/invoices/${invoiceId}/pay`, {
      method: "POST"
    })
  );
}

export function formatCurrency(value: number | string) {
  const numeric = typeof value === "number" ? value : Number.parseFloat(String(value));
  if (Number.isNaN(numeric)) {
    return "￥0.00";
  }
  return `￥${numeric.toFixed(2)}`;
}

export function formatProductType(type: string) {
  const mapping: Record<string, string> = {
    CLOUD: "云主机",
    BANDWIDTH: "带宽产品",
    COLOCATION: "机柜托管",
    HOST: "物理服务器",
    NETWORK: "网络产品"
  };
  return mapping[type] ?? type;
}

export function formatBillingCycle(cycle: string) {
  const mapping: Record<string, string> = {
    monthly: "月付",
    quarterly: "季付",
    semiannual: "半年付",
    annual: "年付",
    annually: "年付"
  };
  return mapping[cycle] ?? cycle;
}

export function formatOrderStatus(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "待支付",
    ACTIVE: "已开通"
  };
  return mapping[status] ?? status;
}

export function formatInvoiceStatus(status: string) {
  const mapping: Record<string, string> = {
    UNPAID: "未支付",
    PAID: "已支付",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

export function formatServiceStatus(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "待开通",
    ACTIVE: "运行中",
    SUSPENDED: "已暂停",
    TERMINATED: "已终止"
  };
  return mapping[status] ?? status;
}

export function formatIdentityStatus(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "待审核",
    APPROVED: "已通过",
    REJECTED: "已驳回"
  };
  return mapping[status] ?? status;
}

export function tagTypeByStatus(status: string): "success" | "warning" | "danger" | "info" {
  switch (status) {
    case "ACTIVE":
    case "PAID":
    case "APPROVED":
    case "CLOSED":
    case "COMPLETED":
      return "success";
    case "PENDING":
    case "UNPAID":
    case "PROCESSING":
      return "warning";
    case "REJECTED":
    case "DISABLED":
    case "REFUNDED":
    case "TERMINATED":
      return "danger";
    default:
      return "info";
  }
}
