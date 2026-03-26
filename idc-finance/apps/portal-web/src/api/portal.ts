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
  orderId: number;
  invoiceId: number;
  productName: string;
  status: string;
  nextDueAt: string;
  createdAt: string;
  providerType?: string;
  configuration?: PortalServiceConfigSelection[];
}

export interface PortalServiceResourceSnapshot {
  regionName: string;
  zoneName: string;
  hostname: string;
  operatingSystem: string;
  loginUsername: string;
  passwordHint: string;
  securityGroup: string;
  cpuCores: number;
  memoryGB: number;
  systemDiskGB: number;
  dataDiskGB: number;
  bandwidthMbps: number;
  publicIpv4: string;
  publicIpv6: string;
}

export interface PortalServiceDetailRecord extends PortalService {
  providerType: string;
  providerAccountId: number;
  providerResourceId: string;
  regionName: string;
  ipAddress: string;
  syncStatus: string;
  syncMessage: string;
  lastAction: string;
  lastSyncAt: string;
  updatedAt: string;
  resourceSnapshot: PortalServiceResourceSnapshot;
}

export interface PortalServiceChangeOrder {
  id: number;
  serviceId: number;
  orderId: number;
  orderNo: string;
  invoiceId: number;
  invoiceNo: string;
  actionName: string;
  title: string;
  amount: number;
  status: string;
  executionStatus: string;
  executionMessage: string;
  reason: string;
  billingCycle: string;
  payload?: Record<string, unknown>;
  paidAt: string;
  refundedAt: string;
  createdAt: string;
  updatedAt: string;
}

export interface PortalServiceDetailResponse {
  service: PortalServiceDetailRecord;
  order?: PortalOrder;
  invoice?: PortalInvoice;
  changeOrders: PortalServiceChangeOrder[];
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

export interface PortalTicketRecord {
  id: number;
  ticketNo: string;
  customerId: number;
  customerNo: string;
  customerName: string;
  serviceId: number;
  serviceNo: string;
  productName: string;
  title: string;
  content: string;
  status: string;
  priority: string;
  source: string;
  departmentName: string;
  assignedAdminId: number;
  assignedAdminName: string;
  latestReplyExcerpt: string;
  lastReplyAt: string;
  closedAt: string;
  createdAt: string;
  updatedAt: string;
  slaStatus: string;
  slaDeadlineAt: string;
  slaRemainingMins: number;
  slaPaused: boolean;
  autoCloseAt: string;
  autoCloseMins: number;
}

export interface PortalTicketReply {
  id: number;
  ticketId: number;
  authorType: string;
  authorId: number;
  authorName: string;
  content: string;
  isInternal: boolean;
  createdAt: string;
  updatedAt: string;
}

export interface PortalTicketDetailResponse {
  ticket: PortalTicketRecord;
  replies: PortalTicketReply[];
}

export interface PortalTicketDepartment {
  key: string;
  name: string;
  description: string;
  enabled: boolean;
  isDefault: boolean;
  sort: number;
}

export interface CreatePortalTicketRequest {
  serviceId?: number;
  title: string;
  content: string;
  priority?: string;
  departmentName?: string;
}

export interface ReplyPortalTicketRequest {
  content: string;
  status?: string;
}

export interface PortalWallet {
  balance: string;
  creditLimit: string;
  creditUsed: string;
  availableCredit: string;
}

export interface PortalWalletTransaction {
  id: number;
  transactionNo: string;
  customerId: number;
  customerNo: string;
  customerName: string;
  orderId: number;
  orderNo: string;
  invoiceId: number;
  invoiceNo: string;
  paymentId: number;
  paymentNo: string;
  refundId: number;
  refundNo: string;
  transactionType: string;
  direction: string;
  amount: number;
  balanceBefore: number;
  balanceAfter: number;
  creditBefore: number;
  creditAfter: number;
  channel: string;
  summary: string;
  remark: string;
  operatorType: string;
  operatorId: number;
  operatorName: string;
  occurredAt: string;
}

export interface PortalWalletOverview {
  wallet: PortalWallet;
  transactions: PortalWalletTransaction[];
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

export function loadWallet() {
  return unwrap(requestJson<ApiEnvelope<PortalWalletOverview>>("/wallet"));
}

export async function loadStoreProducts() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalProduct>>>("/products");
  return result.data.items;
}

export async function loadServices() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalService>>>("/services");
  return result.data.items;
}

export function loadServiceDetail(id: number | string) {
  return unwrap(requestJson<ApiEnvelope<PortalServiceDetailResponse>>(`/services/${id}`));
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
  const result = await requestJson<ApiEnvelope<ListResponse<PortalTicketRecord>>>("/tickets");
  return result.data.items;
}

export async function loadTicketDepartments() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalTicketDepartment>>>("/tickets/departments");
  return result.data.items;
}

export function loadTicketDetail(ticketId: number | string) {
  return unwrap(requestJson<ApiEnvelope<PortalTicketDetailResponse>>(`/tickets/${ticketId}`));
}

export function createTicket(payload: CreatePortalTicketRequest) {
  return unwrap(
    requestJson<ApiEnvelope<PortalTicketRecord>>("/tickets", {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export function replyTicket(ticketId: number | string, payload: ReplyPortalTicketRequest) {
  return unwrap(
    requestJson<ApiEnvelope<PortalTicketDetailResponse>>(`/tickets/${ticketId}/replies`, {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export function closeTicket(ticketId: number | string) {
  return unwrap(
    requestJson<ApiEnvelope<PortalTicketRecord>>(`/tickets/${ticketId}/close`, {
      method: "POST"
    })
  );
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
    return "¥0.00";
  }
  return `¥${numeric.toFixed(2)}`;
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
