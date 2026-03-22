import http from "./http";

export interface Customer {
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
  identity?: IdentityRecord | null;
  contacts?: Contact[];
}

export interface Contact {
  id: number;
  name: string;
  email: string;
  mobile: string;
  roleName: string;
  isPrimary: boolean;
}

export interface IdentityRecord {
  id: number;
  identityType: string;
  verifyStatus: "PENDING" | "APPROVED" | "REJECTED";
  subjectName: string;
  certNo: string;
  countryCode: string;
  reviewRemark?: string;
  reviewedAt?: string;
  submittedAt?: string;
}

export interface CustomerGroupRecord {
  id: number;
  name: string;
  description: string;
  customerCount: number;
}

export interface CustomerLevelRecord {
  id: number;
  name: string;
  priority: number;
  description: string;
  customerCount: number;
}

export interface CustomerIdentityOverview extends IdentityRecord {
  customerId: number;
  customerNo: string;
  customerName: string;
  customerType: string;
}

export interface RelatedItem {
  id?: number;
  serviceId?: number;
  invoiceId?: number;
  ticketId?: number;
  no?: string;
  name?: string;
  status?: string;
  amount?: string;
  dueAt?: string;
  updatedAt?: string;
  billingCycle?: string;
  providerType?: string;
  providerResourceId?: string;
  regionName?: string;
  ipAddress?: string;
  description?: string;
  serviceNo?: string;
  productName?: string;
  nextDueAt?: string;
  invoiceNo?: string;
  totalAmount?: string;
  ticketNo?: string;
  title?: string;
}

export interface AuditLog {
  id: number;
  actor: string;
  action: string;
  description: string;
  createdAt: string;
  target?: string;
  payload?: Record<string, unknown>;
}

export interface ProductPriceOption {
  cycleCode: string;
  cycleName: string;
  price: number;
  setupFee: number;
}

export interface ProductConfigChoice {
  value: string;
  label: string;
  priceDelta: number;
}

export interface ProductConfigOption {
  code: string;
  name: string;
  inputType: string;
  required: boolean;
  defaultValue: string;
  description: string;
  choices: ProductConfigChoice[];
}

export interface ProductResourceTemplate {
  regionName: string;
  zoneName: string;
  operatingSystem: string;
  loginUsername: string;
  securityGroup: string;
  cpuCores: number;
  memoryGB: number;
  systemDiskGB: number;
  dataDiskGB: number;
  bandwidthMbps: number;
  publicIpCount: number;
}

export interface ProductAutomationConfig {
  providerAccountId: number;
  channel: string;
  moduleType: string;
  provisionStage: string;
  autoProvision: boolean;
  serverGroup: string;
  providerNode: string;
}

export interface ProductUpstreamMapping {
  providerAccountId: number;
  providerType: string;
  sourceName: string;
  remoteProductCode: string;
  remoteProductName: string;
  pricePolicy: string;
  autoSyncPricing: boolean;
  autoSyncConfig: boolean;
  autoSyncTemplate: boolean;
  syncStatus: string;
  syncMessage: string;
  lastSyncedAt: string;
}

export interface Product {
  id: number;
  productNo: string;
  groupName: string;
  name: string;
  description: string;
  productType: string;
  status: string;
  pricing: ProductPriceOption[];
  configOptions: ProductConfigOption[];
  resourceTemplate: ProductResourceTemplate;
  automationConfig: ProductAutomationConfig;
  upstreamMapping: ProductUpstreamMapping;
}

export interface ProductDetailResponse {
  product: Product;
  auditLogs: AuditLog[];
}

export interface OrderRecord {
  id: number;
  orderNo: string;
  customerId: number;
  customerName: string;
  productId: number;
  productName: string;
  productType: string;
  automationType: string;
  providerAccountId: number;
  billingCycle: string;
  amount: number;
  status: string;
  payment: string;
  payStatus: string;
  createdAt: string;
}

export interface OrderQuery {
  page?: number;
  limit?: number;
  order?: string;
  sort?: string;
  status?: string;
  ordernum?: string;
  start_time?: string;
  end_time?: string;
  amount?: number | string;
  uid?: number | string;
  payment?: string;
  pay_status?: string;
  product_name?: string;
}

export interface InvoiceQuery {
  page?: number;
  limit?: number;
  sort?: string;
  order?: string;
  status?: string;
  invoice_no?: string;
  order_no?: string;
  product_name?: string;
  billing_cycle?: string;
  uid?: number | string;
}

export interface ServiceQuery {
  page?: number;
  limit?: number;
  sort?: string;
  order?: string;
  status?: string;
  provider_type?: string;
  provider_account_id?: number | string;
  sync_status?: string;
  keyword?: string;
}

export interface ServiceRecord {
  id: number;
  serviceNo: string;
  orderId: number;
  invoiceId: number;
  customerId: number;
  productName: string;
  providerType: string;
  providerAccountId: number;
  providerResourceId: string;
  regionName: string;
  ipAddress: string;
  status: string;
  syncStatus: string;
  syncMessage: string;
  nextDueAt: string;
  lastAction: string;
  lastSyncAt: string;
  updatedAt: string;
  configuration: ServiceConfigSelection[];
  resourceSnapshot: ServiceResourceSnapshot;
  createdAt: string;
}

export interface ServiceConfigSelection {
  code: string;
  name: string;
  value: string;
  valueLabel: string;
  priceDelta: number;
}

export interface ServiceResourceSnapshot {
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

export interface RefundRecord {
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

export interface PaymentRecord {
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

export interface InvoiceRecord {
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

export interface OrderDetailResponse {
  order: OrderRecord;
  invoices: InvoiceRecord[];
  services: ServiceRecord[];
  changeOrder?: ServiceChangeOrderRecord;
  auditLogs: AuditLog[];
}

export interface InvoiceDetailResponse {
  invoice: InvoiceRecord;
  order?: OrderRecord;
  services: ServiceRecord[];
  changeOrder?: ServiceChangeOrderRecord;
  payments: PaymentRecord[];
  refunds: RefundRecord[];
  auditLogs: AuditLog[];
}

export interface ServiceDetailResponse {
  service: ServiceRecord;
  order?: OrderRecord;
  invoice?: InvoiceRecord;
  changeOrders: ServiceChangeOrderRecord[];
  auditLogs: AuditLog[];
}

export interface ServiceChangeOrderRecord {
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
  payload: Record<string, unknown>;
  paidAt: string;
  refundedAt: string;
  createdAt: string;
  updatedAt: string;
  latestTask?: AutomationTask;
}

export interface MofangHealthResponse {
  enabled: boolean;
  connected: boolean;
  baseUrl: string;
  activeUrl?: string;
  authMode: string;
  message: string;
  lastAuthAt?: string;
}

export interface ProviderAccount {
  id: number;
  providerType: string;
  name: string;
  baseUrl: string;
  username: string;
  password: string;
  sourceName: string;
  contactWay: string;
  description: string;
  accountMode: string;
  lang: string;
  listPath: string;
  detailPath: string;
  insecureSkipVerify: boolean;
  autoUpdate: boolean;
  productCount: number;
  status: string;
  extraConfig: string;
  createdAt: string;
  updatedAt: string;
}

export interface MofangInstanceSummary {
  remoteId: string;
  name: string;
  status: string;
  region?: string;
  ipAddress?: string;
  consoleUrl?: string;
  expiresAt?: string;
  raw?: Record<string, unknown>;
}

export interface MofangInstanceDetail extends MofangInstanceSummary {}

export interface MofangInstanceActionResponse {
  ok: boolean;
  action: string;
  remoteId: string;
  status: string;
  message: string;
  response?: Record<string, unknown>;
}

export interface MofangSyncItem {
  remoteId: string;
  serviceId?: number;
  serviceNo?: string;
  customerId?: number;
  customerName?: string;
  operation: "created" | "updated" | "failed";
  status?: string;
  message: string;
}

export interface MofangSyncSummary {
  remoteCount: number;
  processedCount: number;
  createdCustomers: number;
  createdServices: number;
  updatedServices: number;
  failedServices: number;
  syncedNetworkInterfaces: number;
  syncedIps: number;
  syncedDisks: number;
  syncedSnapshots: number;
  syncedBackups: number;
  syncedVpcs: number;
  syncedSecurityGroups: number;
  syncedSecurityRules: number;
  message: string;
}

export interface MofangSyncResponse {
  summary: MofangSyncSummary;
  items: MofangSyncItem[];
}

export interface MofangSyncLogItem {
  id: number;
  providerType: string;
  action: string;
  resourceType: string;
  resourceId: string;
  serviceId?: number;
  status: string;
  message: string;
  createdAt: string;
}

export interface MofangNetworkInterface {
  id: number;
  providerInterfaceId?: string;
  name: string;
  macAddress: string;
  bridge: string;
  networkName: string;
  interfaceName: string;
  nicModel: string;
  inboundMbps: number;
  outboundMbps: number;
  vpcName: string;
  status: string;
}

export interface MofangIpAddress {
  id: number;
  networkInterfaceId?: number;
  providerIpId?: string;
  address: string;
  version: string;
  ipRole: string;
  subnetMask: string;
  gateway: string;
  bandwidthMbps: number;
  isPrimary: boolean;
  status: string;
}

export interface MofangDisk {
  id: number;
  providerDiskId?: string;
  name: string;
  diskType: string;
  sizeGb: number;
  deviceName: string;
  driverName: string;
  fsType: string;
  mountPoint: string;
  isSystem: boolean;
  status: string;
}

export interface MofangSnapshot {
  id: number;
  diskId?: number;
  providerSnapshotId?: string;
  name: string;
  sizeGb: number;
  status: string;
}

export interface MofangBackup {
  id: number;
  providerBackupId?: string;
  name: string;
  sizeGb: number;
  status: string;
  expiresAt?: string;
}

export interface AutomationTaskSummary {
  total: number;
  running: number;
  success: number;
  failed: number;
  blocked: number;
}

export interface AutomationSettings {
  autoProvisionEnabled: boolean;
  autoSuspendEnabled: boolean;
  suspendAfterDays: number;
  suspendNoticeDays: number;
  autoTerminateEnabled: boolean;
  terminateAfterDays: number;
  terminateNoticeDays: number;
  invoiceReminderOn: boolean;
  invoiceReminderDays: string;
  creditReminderOn: boolean;
  creditReminderDays: string;
  ticketAutoCloseOn: boolean;
  ticketAutoCloseHours: number;
  logRetentionDays: number;
}

export interface AutomationTask {
  id: number;
  taskNo: string;
  taskType: string;
  title: string;
  channel: string;
  stage: string;
  status: "PENDING" | "RUNNING" | "SUCCESS" | "FAILED" | "BLOCKED" | string;
  sourceType: string;
  sourceId: number;
  customerId: number;
  customerName: string;
  productName: string;
  orderId: number;
  invoiceId: number;
  serviceId: number;
  serviceNo: string;
  providerType: string;
  providerResourceId: string;
  actionName: string;
  operatorType: string;
  operatorName: string;
  requestPayload: string;
  resultPayload: string;
  message: string;
  createdAt: string;
  startedAt: string;
  finishedAt: string;
  updatedAt: string;
}

export interface AutomationTaskListResponse {
  items: AutomationTask[];
  summary: AutomationTaskSummary;
  total: number;
}

export interface AutomationRetryResult {
  sourceTask: AutomationTask;
  triggeredTask?: AutomationTask;
  message: string;
}

export interface MofangVpcNetwork {
  id: number;
  providerVpcId?: string;
  name: string;
  cidr: string;
  gateway: string;
  interfaceName: string;
  status: string;
}

export interface MofangSecurityGroupRule {
  id: number;
  direction: string;
  protocol: string;
  portRange: string;
  sourceCidr: string;
  action: string;
  description: string;
}

export interface MofangSecurityGroup {
  id: number;
  providerSecurityGroupId?: string;
  name: string;
  status: string;
  rules: MofangSecurityGroupRule[];
}

export interface MofangServiceResourcesResponse {
  serviceId: number;
  serviceNo: string;
  providerType: string;
  providerResourceId: string;
  status: string;
  regionName: string;
  ipAddress: string;
  syncStatus: string;
  syncMessage: string;
  lastSyncAt: string;
  networkInterfaces: MofangNetworkInterface[];
  ipAddresses: MofangIpAddress[];
  disks: MofangDisk[];
  snapshots: MofangSnapshot[];
  backups: MofangBackup[];
  vpcNetworks: MofangVpcNetwork[];
  securityGroups: MofangSecurityGroup[];
  syncLogs: MofangSyncLogItem[];
}

export interface MofangResourceActionRequest {
  count?: number;
  sizeGb?: number;
  storeId?: number;
  ipGroup?: number;
  driver?: string;
  name?: string;
  diskId?: string;
  snapshotId?: string;
  backupId?: string;
}

export interface MofangResourceActionResponse {
  ok: boolean;
  action: string;
  serviceId: number;
  serviceNo: string;
  remoteId: string;
  resourceId?: string;
  message: string;
  poweredOff: boolean;
  poweredOn: boolean;
  syncItem?: MofangSyncItem;
}

export interface UpstreamCatalogProduct {
  remoteProductCode: string;
  groupName: string;
  name: string;
  productType: string;
  description: string;
}

export interface UpstreamCatalogGroup {
  groupId: string;
  groupName: string;
  products: UpstreamCatalogProduct[];
}

export interface ImportUpstreamProductsRequest {
  providerAccountId: number;
  providerType: string;
  importAll: boolean;
  remoteProductCodes: string[];
  autoSyncPricing: boolean;
  autoSyncConfig: boolean;
  autoSyncTemplate: boolean;
}

export interface ImportUpstreamProductResult {
  remoteProductCode: string;
  remoteProductName: string;
  groupName?: string;
  productId?: number;
  productNo?: string;
  status?: string;
  operation: "create" | "update" | "skip" | "failed" | string;
  message: string;
}

export interface ImportUpstreamProductsResponse {
  providerAccountId: number;
  items: ImportUpstreamProductResult[];
  importedCount?: number;
  updatedCount?: number;
  failedCount?: number;
  total: number;
  created: number;
  updated: number;
  skipped: number;
  failed: number;
  message: string;
}

export interface UpstreamCyclePrice {
  cycleCode: string;
  cycleName: string;
  price: number;
  setupFee: number;
}

export interface UpstreamProductTemplate {
  remoteProductCode: string;
  groupName: string;
  name: string;
  description: string;
  productType: string;
  currency: string;
  pricing: UpstreamCyclePrice[];
  configOptions: ProductConfigOption[];
}

export interface ReceivePaymentResponse {
  invoice: InvoiceRecord;
  service?: ServiceRecord;
  payment: PaymentRecord;
}

export interface RefundInvoiceResponse {
  invoice: InvoiceRecord;
  service?: ServiceRecord;
  refund: RefundRecord;
}

export interface UpdatePendingOrderRequest {
  productName?: string;
  billingCycle?: string;
  amount?: number;
  dueAt?: string;
  status?: string;
  reason?: string;
}

export interface UpdateUnpaidInvoiceRequest {
  productName?: string;
  billingCycle?: string;
  amount?: number;
  dueAt?: string;
  status?: string;
  reason?: string;
}

export interface UpdateServiceRequest {
  providerType?: string;
  providerAccountId?: number | null;
  providerResourceId?: string;
  regionName?: string;
  ipAddress?: string;
  nextDueAt?: string;
  status?: string;
  syncStatus?: string;
  syncMessage?: string;
  reason?: string;
}

export interface CreateServiceChangeOrderRequest {
  actionName: string;
  title: string;
  billingCycle: string;
  amount: number;
  reason: string;
  payload?: Record<string, unknown>;
}

export interface CreateServiceChangeOrderResponse {
  service: ServiceRecord;
  order: OrderRecord;
  invoice: InvoiceRecord;
}

export interface MetricCard {
  key: string;
  label: string;
  value: string;
  hint: string;
  tone: string;
}

export interface StatusBucket {
  name: string;
  count: number;
  amount: number;
}

export interface TrendPoint {
  label: string;
  amount: number;
  count: number;
}

export interface PendingIdentityTodo {
  customerId: number;
  customerNo: string;
  customerName: string;
  subjectName: string;
  submittedAt: string;
}

export interface OverdueInvoiceTodo {
  invoiceId: number;
  invoiceNo: string;
  customerName: string;
  amount: number;
  dueAt: string;
  daysOverdue: number;
}

export interface ExpiringServiceTodo {
  serviceId: number;
  serviceNo: string;
  customerName: string;
  productName: string;
  status: string;
  nextDueAt: string;
  daysRemaining: number;
}

export interface TicketTodo {
  ticketId: number;
  ticketNo: string;
  customerName: string;
  title: string;
  status: string;
  updatedAt: string;
}

export interface AuditTimelineItem {
  id: number;
  actor: string;
  action: string;
  target: string;
  description: string;
  createdAt: string;
}

export interface ReportRankItem {
  name: string;
  count: number;
  amount: number;
  extra: string;
}

export interface WorkbenchResponse {
  summaryCards: MetricCard[];
  riskCards: MetricCard[];
  revenueTrend: TrendPoint[];
  serviceStatus: StatusBucket[];
  pendingIdentities: PendingIdentityTodo[];
  overdueInvoices: OverdueInvoiceTodo[];
  expiringServices: ExpiringServiceTodo[];
  openTickets: TicketTodo[];
  recentAudits: AuditTimelineItem[];
}

export interface ReportOverviewResponse {
  headline: MetricCard[];
  revenueTrend: TrendPoint[];
  refundTrend: TrendPoint[];
  invoiceStatus: StatusBucket[];
  serviceStatus: StatusBucket[];
  paymentChannels: StatusBucket[];
  billingCycles: StatusBucket[];
  customerGroups: StatusBucket[];
  topProducts: ReportRankItem[];
  topReceivables: ReportRankItem[];
}

export async function fetchCustomers() {
  const { data } = await http.get("/customers");
  return data.data as { items: Customer[]; total: number };
}

export async function fetchCustomerDetail(id: number | string) {
  const { data } = await http.get(`/customers/${id}`);
  return data.data as Customer;
}

export async function updateCustomer(id: number | string, payload: Partial<Customer>) {
  const { data } = await http.patch(`/customers/${id}`, payload);
  return data.data as Customer;
}

export async function fetchCustomerContacts(id: number | string) {
  const { data } = await http.get(`/customers/${id}/contacts`);
  return data.data as Contact[];
}

export async function createCustomerContact(id: number | string, payload: Record<string, unknown>) {
  const { data } = await http.post(`/customers/${id}/contacts`, payload);
  return data.data as Contact;
}

export async function updateCustomerContact(customerId: number | string, contactId: number | string, payload: Record<string, unknown>) {
  const { data } = await http.put(`/customers/${customerId}/contacts/${contactId}`, payload);
  return data.data as Contact;
}

export async function deleteCustomerContact(customerId: number | string, contactId: number | string) {
  const { data } = await http.delete(`/customers/${customerId}/contacts/${contactId}`);
  return data.data as { deleted: boolean };
}

export async function fetchCustomerIdentitiesDetail(id: number | string) {
  const { data } = await http.get(`/customers/${id}/identities`);
  return data.data as IdentityRecord[];
}

export async function reviewCustomerIdentity(
  customerId: number | string,
  identityId: number | string,
  payload: { status: "APPROVED" | "REJECTED"; remark?: string }
) {
  const { data } = await http.post(`/customers/${customerId}/identities/${identityId}/review`, payload);
  return data.data as IdentityRecord;
}

export async function fetchCustomerServices(id: number | string) {
  const { data } = await http.get(`/customers/${id}/services`);
  return data.data as RelatedItem[];
}

export async function fetchCustomerInvoices(id: number | string) {
  const { data } = await http.get(`/customers/${id}/invoices`);
  return data.data as RelatedItem[];
}

export async function fetchCustomerTickets(id: number | string) {
  const { data } = await http.get(`/customers/${id}/tickets`);
  return data.data as RelatedItem[];
}

export async function fetchCustomerAuditLogs(id: number | string) {
  const { data } = await http.get(`/customers/${id}/audit-logs`);
  return data.data as AuditLog[];
}

export async function createCustomer(payload: Record<string, unknown>) {
  const { data } = await http.post("/customers", payload);
  return data.data as Customer;
}

export async function fetchCustomerGroups() {
  const { data } = await http.get("/customer-groups");
  return data.data as CustomerGroupRecord[];
}

export async function createCustomerGroup(payload: { name: string; description: string }) {
  const { data } = await http.post("/customer-groups", payload);
  return data.data as CustomerGroupRecord;
}

export async function updateCustomerGroup(id: number | string, payload: { name: string; description: string }) {
  const { data } = await http.patch(`/customer-groups/${id}`, payload);
  return data.data as CustomerGroupRecord;
}

export async function deleteCustomerGroup(id: number | string) {
  const { data } = await http.delete(`/customer-groups/${id}`);
  return data.data as { deleted: boolean };
}

export async function fetchCustomerLevels() {
  const { data } = await http.get("/customer-levels");
  return data.data as CustomerLevelRecord[];
}

export async function createCustomerLevel(payload: { name: string; priority: number; description: string }) {
  const { data } = await http.post("/customer-levels", payload);
  return data.data as CustomerLevelRecord;
}

export async function updateCustomerLevel(
  id: number | string,
  payload: { name: string; priority: number; description: string }
) {
  const { data } = await http.patch(`/customer-levels/${id}`, payload);
  return data.data as CustomerLevelRecord;
}

export async function deleteCustomerLevel(id: number | string) {
  const { data } = await http.delete(`/customer-levels/${id}`);
  return data.data as { deleted: boolean };
}

export async function fetchCustomerIdentities() {
  const { data } = await http.get("/customer-identities");
  return data.data as CustomerIdentityOverview[];
}

export async function fetchProductGroups() {
  const { data } = await http.get("/product-groups");
  return data.data as string[];
}

export async function fetchProducts() {
  const { data } = await http.get("/products");
  return data.data as { items: Product[]; total: number };
}

export async function fetchProductDetail(id: number | string) {
  const { data } = await http.get(`/products/${id}`);
  return data.data as ProductDetailResponse;
}

export async function createProduct(payload: Record<string, unknown>) {
  const { data } = await http.post("/products", payload);
  return data.data as Product;
}

export async function updateProduct(id: number | string, payload: Record<string, unknown>) {
  const { data } = await http.patch(`/products/${id}`, payload);
  return data.data as Product;
}

export async function syncProductUpstream(
  id: number | string,
  payload: Partial<ProductUpstreamMapping>
) {
  const { data } = await http.post(`/products/${id}/upstream/sync`, payload);
  return data.data as Product;
}

export async function importUpstreamProducts(payload: ImportUpstreamProductsRequest) {
  const { data } = await http.post("/products/import-upstream", payload);
  return data.data as ImportUpstreamProductsResponse;
}

export async function fetchOrders(params?: OrderQuery) {
  const { data } = await http.get("/orders", { params });
  return data.data as { items: OrderRecord[]; total: number };
}

export async function fetchOrderDetail(id: number | string) {
  const { data } = await http.get(`/orders/${id}`);
  return data.data as OrderDetailResponse;
}

export async function updatePendingOrder(id: number | string, payload: UpdatePendingOrderRequest) {
  const { data } = await http.patch(`/orders/${id}`, payload);
  return data.data as OrderDetailResponse;
}

export async function fetchInvoices(params?: InvoiceQuery) {
  const { data } = await http.get("/invoices", { params });
  return data.data as { items: InvoiceRecord[]; total: number };
}

export async function fetchInvoiceDetail(id: number | string) {
  const { data } = await http.get(`/invoices/${id}`);
  return data.data as InvoiceDetailResponse;
}

export async function updateUnpaidInvoice(id: number | string, payload: UpdateUnpaidInvoiceRequest) {
  const { data } = await http.patch(`/invoices/${id}`, payload);
  return data.data as InvoiceDetailResponse;
}

export async function receiveInvoicePayment(
  id: number | string,
  payload?: { channel?: string; tradeNo?: string }
) {
  const { data } = await http.post(`/invoices/${id}/receive-payment`, payload ?? {});
  return data.data as ReceivePaymentResponse;
}

export async function refundInvoice(id: number | string, reason: string) {
  const { data } = await http.post(`/invoices/${id}/refund`, { reason });
  return data.data as RefundInvoiceResponse;
}

export async function fetchServices(params?: ServiceQuery) {
  const { data } = await http.get("/services", { params });
  return data.data as { items: ServiceRecord[]; total: number };
}

export async function fetchServiceDetail(id: number | string) {
  const { data } = await http.get(`/services/${id}`);
  return data.data as ServiceDetailResponse;
}

export async function updateServiceRecord(id: number | string, payload: UpdateServiceRequest) {
  const { data } = await http.patch(`/services/${id}`, payload);
  return data.data as ServiceDetailResponse;
}

export async function createServiceChangeOrder(id: number | string, payload: CreateServiceChangeOrderRequest) {
  const { data } = await http.post(`/services/${id}/change-orders`, payload);
  return data.data as CreateServiceChangeOrderResponse;
}

export async function fetchProviderAccounts(providerType?: string) {
  const query = providerType ? `?providerType=${encodeURIComponent(providerType)}` : "";
  const { data } = await http.get(`/providers/accounts${query}`);
  return data.data as ProviderAccount[];
}

export async function createProviderAccount(payload: Partial<ProviderAccount>) {
  const { data } = await http.post("/providers/accounts", payload);
  return data.data as ProviderAccount;
}

export async function updateProviderAccount(id: number | string, payload: Partial<ProviderAccount>) {
  const { data } = await http.patch(`/providers/accounts/${id}`, payload);
  return data.data as ProviderAccount;
}

export async function deleteProviderAccount(id: number | string) {
  const { data } = await http.delete(`/providers/accounts/${id}`);
  return data.data as { deleted: boolean };
}

export async function fetchProviderAccountHealth(id: number | string) {
  const { data } = await http.get(`/providers/accounts/${id}/health`);
  return data.data as MofangHealthResponse;
}

export async function fetchMofangHealth(accountId?: number | string) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.get(`/providers/mofang/health${query}`);
  return data.data as MofangHealthResponse;
}

export async function fetchUpstreamHealth(accountId?: number | string) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.get(`/providers/upstream/health${query}`);
  return data.data as MofangHealthResponse;
}

export async function fetchUpstreamProducts(accountId?: number | string) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.get(`/providers/upstream/products${query}`);
  return data.data as { groups: UpstreamCatalogGroup[]; total: number };
}

export async function fetchUpstreamProductTemplate(id: string | number, accountId?: number | string) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.get(`/providers/upstream/products/${id}/template${query}`);
  return data.data as UpstreamProductTemplate;
}

export async function fetchMofangInstances(accountId?: number | string) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.get(`/providers/mofang/instances${query}`);
  return data.data as { items: MofangInstanceSummary[]; total: number };
}

export async function fetchMofangInstanceDetail(id: number | string, accountId?: number | string) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.get(`/providers/mofang/instances/${id}${query}`);
  return data.data as MofangInstanceDetail;
}

export async function runMofangInstanceAction(
  id: number | string,
  action:
    | "suspend"
    | "power-on"
    | "power-off"
    | "reboot"
    | "hard-reboot"
    | "hard-power-off"
    | "reset-password"
    | "reinstall"
    | "unsuspend"
    | "get-vnc"
    | "rescue-start"
    | "rescue-stop"
    | "lock"
    | "unlock",
  payload?: { imageName?: string; password?: string },
  accountId?: number | string
) {
  const query = accountId !== undefined ? `?accountId=${encodeURIComponent(String(accountId))}` : "";
  const { data } = await http.post(`/providers/mofang/instances/${id}/actions/${action}${query}`, payload ?? {});
  return data.data as MofangInstanceActionResponse;
}

export async function pullMofangSync(payload?: {
  providerAccountId?: number;
  limit?: number;
  includeResources?: boolean;
  remoteIds?: string[];
}) {
  const { data } = await http.post("/providers/mofang/sync", payload ?? {});
  return data.data as MofangSyncResponse;
}

export async function syncMofangService(id: number | string, includeResources = true) {
  const { data } = await http.post(`/providers/mofang/services/${id}/sync`, {
    includeResources
  });
  return data.data as MofangSyncItem;
}

export async function fetchMofangServiceResources(id: number | string) {
  const { data } = await http.get(`/providers/mofang/services/${id}/resources`);
  return data.data as MofangServiceResourcesResponse;
}

export async function runMofangServiceResourceAction(
  id: number | string,
  action:
    | "add-ipv4"
    | "add-ipv6"
    | "add-disk"
    | "resize-disk"
    | "create-snapshot"
    | "delete-snapshot"
    | "restore-snapshot"
    | "create-backup"
    | "delete-backup"
    | "restore-backup",
  payload?: MofangResourceActionRequest
) {
  const { data } = await http.post(`/providers/mofang/services/${id}/resource-actions/${action}`, payload ?? {});
  return data.data as MofangResourceActionResponse;
}

export async function fetchMofangSyncLogs(params?: { serviceId?: number | string; limit?: number }) {
  const search = new URLSearchParams();
  if (params?.serviceId !== undefined) search.set("serviceId", String(params.serviceId));
  if (params?.limit !== undefined) search.set("limit", String(params.limit));
  const query = search.toString();
  const { data } = await http.get(`/providers/mofang/sync-logs${query ? `?${query}` : ""}`);
  return data.data as { items: MofangSyncLogItem[]; total: number };
}

export async function fetchAutomationTasks(params?: {
  status?: string;
  taskType?: string;
  channel?: string;
  stage?: string;
  sourceType?: string;
  sourceId?: number | string;
  orderId?: number | string;
  invoiceId?: number | string;
  serviceId?: number | string;
  keyword?: string;
  limit?: number;
}) {
  const search = new URLSearchParams();
  if (params?.status) search.set("status", params.status);
  if (params?.taskType) search.set("taskType", params.taskType);
  if (params?.channel) search.set("channel", params.channel);
  if (params?.stage) search.set("stage", params.stage);
  if (params?.sourceType) search.set("sourceType", params.sourceType);
  if (params?.sourceId !== undefined) search.set("sourceId", String(params.sourceId));
  if (params?.orderId !== undefined) search.set("orderId", String(params.orderId));
  if (params?.invoiceId !== undefined) search.set("invoiceId", String(params.invoiceId));
  if (params?.serviceId !== undefined) search.set("serviceId", String(params.serviceId));
  if (params?.keyword) search.set("keyword", params.keyword);
  if (params?.limit !== undefined) search.set("limit", String(params.limit));
  const query = search.toString();
  const { data } = await http.get(`/automation/tasks${query ? `?${query}` : ""}`);
  return data.data as AutomationTaskListResponse;
}

export async function fetchAutomationTaskDetail(id: number | string) {
  const { data } = await http.get(`/automation/tasks/${id}`);
  return data.data as AutomationTask;
}

export async function retryAutomationTask(id: number | string) {
  const { data } = await http.post(`/automation/tasks/${id}/retry`);
  return data.data as AutomationRetryResult;
}

export async function fetchAutomationSettings() {
  const { data } = await http.get("/automation/settings");
  return data.data as AutomationSettings;
}

export async function updateAutomationSettings(payload: AutomationSettings) {
  const { data } = await http.patch("/automation/settings", payload);
  return data.data as AutomationSettings;
}

export async function runServiceAction(
  id: number | string,
  action: "activate" | "suspend" | "terminate" | "reboot" | "reset-password" | "reinstall",
  payload?: { imageName?: string; password?: string }
) {
  const { data } = await http.post(`/services/${id}/actions/${action}`, payload ?? {});
  return data.data as ServiceRecord;
}

export async function fetchAdmins() {
  const { data } = await http.get("/admins");
  return data.data as Array<{
    id: number;
    username: string;
    displayName: string;
    status: string;
    roles: string[];
  }>;
}

export async function fetchRoles() {
  const { data } = await http.get("/roles");
  return data.data as Array<{
    id: number;
    name: string;
    code: string;
    users: number;
  }>;
}

export async function fetchAuditLogs() {
  const { data } = await http.get("/audit-logs");
  return data.data as AuditLog[];
}

export async function fetchWorkbenchOverview() {
  const { data } = await http.get("/workbench");
  return data.data as WorkbenchResponse;
}

export async function fetchReportOverview() {
  const { data } = await http.get("/reports/overview");
  return data.data as ReportOverviewResponse;
}
