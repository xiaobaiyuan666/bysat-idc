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

export interface PortalOrderRequestRecord {
  id: number;
  requestNo: string;
  orderId: number;
  orderNo: string;
  customerId: number;
  customerName: string;
  productName: string;
  type: string;
  status: string;
  summary: string;
  reason: string;
  currentAmount: number;
  requestedAmount: number;
  currentBillingCycle: string;
  requestedBillingCycle: string;
  sourceType: string;
  sourceId: number;
  sourceName: string;
  processorType: string;
  processorId: number;
  processorName: string;
  processNote: string;
  payload?: Record<string, unknown>;
  createdAt: string;
  updatedAt: string;
  processedAt: string;
}

export interface PortalServiceChangeOrder {
  id: number;
  serviceId: number;
  serviceNo: string;
  orderId: number;
  orderNo: string;
  invoiceId: number;
  invoiceNo: string;
  productName: string;
  providerType: string;
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

export interface PortalRefund {
  id: number;
  refundNo: string;
  invoiceId: number;
  orderId: number;
  customerId: number;
  amount: number;
  reason: string;
  status: string;
  createdAt: string;
}

export interface PortalServiceDetailResponse {
  service: PortalServiceDetailRecord;
  order?: PortalOrder;
  invoice?: PortalInvoice;
  changeOrders: PortalServiceChangeOrder[];
}

export interface PortalOrderDetailResponse {
  order: PortalOrder;
  invoices: PortalInvoice[];
  services: PortalServiceDetailRecord[];
  requests: PortalOrderRequestRecord[];
  changeOrder?: PortalServiceChangeOrder;
  auditLogs: Array<Record<string, unknown>>;
}

export interface PortalInvoiceDetailResponse {
  invoice: PortalInvoice;
  order?: PortalOrder;
  services: PortalServiceDetailRecord[];
  changeOrder?: PortalServiceChangeOrder;
  payments: PortalPayment[];
  refunds: PortalRefund[];
  auditLogs: Array<Record<string, unknown>>;
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
  submittedAt: string;
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

export interface CreatePortalOrderRequestPayload {
  type: string;
  summary?: string;
  reason: string;
  requestedAmount?: number;
  requestedBillingCycle?: string;
  payload?: Record<string, unknown>;
}

export interface UpdatePortalProfilePayload {
  name: string;
  email: string;
  mobile: string;
  remarks?: string;
}

export interface SavePortalContactPayload {
  name: string;
  email?: string;
  mobile?: string;
  roleName?: string;
  isPrimary?: boolean;
}

export interface SubmitPortalIdentityPayload {
  identityType: string;
  subjectName: string;
  certNo: string;
  countryCode?: string;
}

export interface ChangePortalPasswordPayload {
  currentPassword: string;
  newPassword: string;
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

export interface PortalKYCRecord {
  id: number;
  pluginCode: string;
  customerId: number;
  realName: string;
  idCardNo: string;
  mobile: string;
  status: string;
  message: string;
  createdAt: string;
}

export interface SubmitZhimaKYCRequest {
  realName: string;
  idCardNo: string;
  mobile?: string;
}

export interface CreateRechargeOrderRequest {
  amount: number;
  subject?: string;
}

export interface RechargeOrderResult {
  id: number;
  pluginCode: string;
  orderNo: string;
  amount: number;
  currency: string;
  subject: string;
  status: string;
  payUrl: string;
  createdAt: string;
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

export function loadWalletTransactions(params?: {
  transactionType?: string;
  keyword?: string;
  channel?: string;
  limit?: number;
}) {
  const search = new URLSearchParams();
  if (params?.transactionType) search.set("transactionType", params.transactionType);
  if (params?.keyword) search.set("keyword", params.keyword);
  if (params?.channel) search.set("channel", params.channel);
  if (params?.limit) search.set("limit", String(params.limit));
  const suffix = search.size ? `?${search.toString()}` : "";
  return unwrap(requestJson<ApiEnvelope<ListResponse<PortalWalletTransaction>>>(`/wallet/transactions${suffix}`));
}

export function loadWalletRecharges(params?: {
  keyword?: string;
  channel?: string;
  limit?: number;
}) {
  const search = new URLSearchParams();
  if (params?.keyword) search.set("keyword", params.keyword);
  if (params?.channel) search.set("channel", params.channel);
  if (params?.limit) search.set("limit", String(params.limit));
  const suffix = search.size ? `?${search.toString()}` : "";
  return unwrap(requestJson<ApiEnvelope<ListResponse<PortalWalletTransaction>>>(`/wallet/recharges${suffix}`));
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

export function loadOrderDetail(id: number | string) {
  return unwrap(requestJson<ApiEnvelope<PortalOrderDetailResponse>>(`/orders/${id}`));
}

export function createOrderRequest(orderId: number | string, payload: CreatePortalOrderRequestPayload) {
  return unwrap(
    requestJson<ApiEnvelope<PortalOrderDetailResponse>>(`/orders/${orderId}/requests`, {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export async function loadInvoices() {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalInvoice>>>("/invoices");
  return result.data.items;
}

export function loadInvoiceDetail(id: number | string) {
  return unwrap(requestJson<ApiEnvelope<PortalInvoiceDetailResponse>>(`/invoices/${id}`));
}

export async function loadPayments(params?: {
  invoiceId?: number;
  keyword?: string;
  channel?: string;
  status?: string;
}) {
  const search = new URLSearchParams();
  if (params?.invoiceId) search.set("invoiceId", String(params.invoiceId));
  if (params?.keyword) search.set("keyword", params.keyword);
  if (params?.channel) search.set("channel", params.channel);
  if (params?.status) search.set("status", params.status);
  const suffix = search.size ? `?${search.toString()}` : "";
  const result = await requestJson<ApiEnvelope<ListResponse<PortalPayment>>>(`/payments${suffix}`);
  return result.data.items;
}

export async function loadRefunds(params?: {
  invoiceId?: number;
  keyword?: string;
  status?: string;
}) {
  const search = new URLSearchParams();
  if (params?.invoiceId) search.set("invoiceId", String(params.invoiceId));
  if (params?.keyword) search.set("keyword", params.keyword);
  if (params?.status) search.set("status", params.status);
  const suffix = search.size ? `?${search.toString()}` : "";
  const result = await requestJson<ApiEnvelope<ListResponse<PortalRefund>>>(`/refunds${suffix}`);
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

export function updateProfile(payload: UpdatePortalProfilePayload) {
  return unwrap(
    requestJson<ApiEnvelope<PortalAccount>>("/account/profile", {
      method: "PUT",
      body: JSON.stringify(payload)
    })
  );
}

export function loadContacts() {
  return unwrap(requestJson<ApiEnvelope<ListResponse<PortalContact>>>("/account/contacts"));
}

export function createContact(payload: SavePortalContactPayload) {
  return unwrap(
    requestJson<ApiEnvelope<PortalContact>>("/account/contacts", {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export function updateContact(contactId: number | string, payload: SavePortalContactPayload) {
  return unwrap(
    requestJson<ApiEnvelope<PortalContact>>(`/account/contacts/${contactId}`, {
      method: "PATCH",
      body: JSON.stringify(payload)
    })
  );
}

export function deleteContact(contactId: number | string) {
  return unwrap(
    requestJson<ApiEnvelope<{ deleted: boolean }>>(`/account/contacts/${contactId}`, {
      method: "DELETE"
    })
  );
}

export function loadIdentity() {
  return unwrap(requestJson<ApiEnvelope<{ identity?: PortalIdentity | null }>>("/account/identity"));
}

export function submitIdentity(payload: SubmitPortalIdentityPayload) {
  return unwrap(
    requestJson<ApiEnvelope<PortalIdentity>>("/account/identity", {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export function changePortalPassword(payload: ChangePortalPasswordPayload) {
  return unwrap(
    requestJson<ApiEnvelope<{ changed: boolean }>>("/account/security/password", {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
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

export function submitZhimaKYC(payload: SubmitZhimaKYCRequest) {
  return unwrap(
    requestJson<ApiEnvelope<PortalKYCRecord>>("/plugins/kyc/zhima_kyc/verify", {
      method: "POST",
      body: JSON.stringify(payload)
    })
  );
}

export async function loadZhimaKYCRecords(limit = 20) {
  const result = await requestJson<ApiEnvelope<ListResponse<PortalKYCRecord>>>(`/plugins/kyc/records?limit=${limit}`);
  return result.data.items;
}

export function createRechargeOrder(payload: CreateRechargeOrderRequest) {
  return unwrap(
    requestJson<ApiEnvelope<RechargeOrderResult>>("/plugins/payment/epay/orders", {
      method: "POST",
      body: JSON.stringify({
        orderNo: `RECHARGE-${Date.now()}`,
        amount: payload.amount,
        currency: "CNY",
        subject: payload.subject || "钱包在线充值"
      })
    })
  );
}
