<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import { formatPaymentChannel } from "@/utils/business";
import {
  type AccountTransactionRecord,
  createCustomerContact,
  deleteCustomerContact,
  fetchCustomerAuditLogs,
  fetchCustomerContacts,
  fetchCustomerDetail,
  fetchCustomerIdentitiesDetail,
  fetchCustomerInvoices,
  fetchPayments,
  fetchRefunds,
  fetchCustomerServices,
  fetchCustomerTickets,
  fetchCustomerWallet,
  reviewCustomerIdentity,
  updateCustomer,
  updateCustomerContact,
  type AuditLog,
  type Contact,
  type Customer,
  type CustomerWalletResponse,
  type IdentityRecord,
  type PaymentRecord,
  type RelatedItem,
  type RefundRecord
} from "@/api/admin";

type WorkbenchRelatedItem = RelatedItem & { entityId?: number };

type FinanceTimelineItem = {
  key: string;
  occurredAt: string;
  title: string;
  description: string;
  referenceNo: string;
  amount: number | undefined;
  categoryLabel: string;
  type: "primary" | "success" | "warning" | "danger" | "info";
  invoiceId?: number;
  linkKeyword?: string;
  transactionType?: string;
  target: "invoice" | "payment" | "refund" | "account";
};

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const contactDialogVisible = ref(false);
const contactSubmitting = ref(false);
const savingCustomer = ref(false);
const editingContactId = ref<number | null>(null);
const walletDetail = ref<CustomerWalletResponse | null>(null);

const customer = reactive<Customer>({
  id: 0,
  customerNo: "",
  name: "",
  email: "",
  mobile: "",
  type: "COMPANY",
  status: "ACTIVE",
  groupName: "",
  levelName: "",
  salesOwner: "",
  remarks: ""
});

const contacts = ref<Contact[]>([]);
const identities = ref<IdentityRecord[]>([]);
const services = ref<WorkbenchRelatedItem[]>([]);
const invoices = ref<WorkbenchRelatedItem[]>([]);
const payments = ref<PaymentRecord[]>([]);
const refunds = ref<RefundRecord[]>([]);
const tickets = ref<WorkbenchRelatedItem[]>([]);
const auditLogs = ref<AuditLog[]>([]);

const contactForm = reactive({
  name: "",
  email: "",
  mobile: "",
  roleName: "",
  isPrimary: false
});

const customerId = computed(() => route.params.id as string);
const approvedIdentityCount = computed(
  () => identities.value.filter(item => item.verifyStatus === "APPROVED").length
);
const walletTransactionCount = computed(() => walletDetail.value?.transactions.length ?? 0);
const walletBalance = computed(() => walletDetail.value?.wallet.balance ?? 0);
const walletAvailableCredit = computed(() => walletDetail.value?.wallet.availableCredit ?? 0);
const walletUpdatedAt = computed(() => walletDetail.value?.wallet.updatedAt || "-");
const financeEventCount = computed(() => payments.value.length + refunds.value.length);
const financeTimeline = computed<FinanceTimelineItem[]>(() => {
  const invoiceItems = invoices.value
    .map(item => {
      const paidAt = String((item as Record<string, unknown>).paidAt ?? "");
      const occurredAt = paidAt || item.dueAt || "";
      if (!occurredAt) return null;
      const amount = Number(item.totalAmount);
      return {
        key: `invoice-${item.entityId ?? item.invoiceNo}`,
        occurredAt,
        title:
          item.status === "PAID" ? "账单已支付" : item.status === "REFUNDED" ? "账单已退款" : "账单待支付",
        description: `${item.productName ?? "关联产品"} · ${invoiceStatusLabel(item.status ?? "")}`,
        referenceNo: item.invoiceNo ?? "-",
        amount: Number.isFinite(amount) ? amount : undefined,
        categoryLabel: "账单",
        type: invoiceStatusTag(item.status ?? ""),
        invoiceId: item.entityId,
        target: "invoice" as const
      };
    })
    .filter(Boolean) as FinanceTimelineItem[];

  const paymentItems = payments.value.map(item => ({
    key: `payment-${item.id}`,
    occurredAt: item.paidAt,
    title: "支付入账",
    description: `${invoiceNoById(item.invoiceId)} · ${formatPaymentChannel(localeStore.locale, item.channel || "SYSTEM")}`,
    referenceNo: item.paymentNo,
    linkKeyword: item.paymentNo,
    amount: item.amount,
    categoryLabel: "支付",
    type: "success" as const,
    invoiceId: item.invoiceId,
    target: "payment" as const
  }));

  const refundItems = refunds.value.map(item => ({
    key: `refund-${item.id}`,
    occurredAt: item.createdAt,
    title: "退款完成",
    description: `${invoiceNoById(item.invoiceId)} · ${item.reason || "后台人工退款"}`,
    referenceNo: item.refundNo,
    linkKeyword: item.refundNo,
    amount: item.amount,
    categoryLabel: "退款",
    type: "warning" as const,
    invoiceId: item.invoiceId,
    target: "refund" as const
  }));

  const accountItems = (walletDetail.value?.transactions ?? []).map(item => ({
    key: `account-${item.id}`,
    occurredAt: item.occurredAt,
    title: `资金台账 · ${transactionTypeLabel(item.transactionType)}`,
    description: item.summary || item.remark || walletActionLabel(item),
    referenceNo: item.transactionNo,
    linkKeyword: item.transactionNo,
    amount: item.amount,
    categoryLabel: "台账",
    type: directionTag(item.direction) as FinanceTimelineItem["type"],
    invoiceId: item.invoiceId > 0 ? item.invoiceId : undefined,
    transactionType: item.transactionType,
    target: "account" as const
  }));

  return ([...paymentItems, ...refundItems, ...accountItems, ...invoiceItems] as FinanceTimelineItem[])
    .filter(item => Boolean(item.occurredAt))
    .sort((left, right) => parseTimelineTime(right.occurredAt) - parseTimelineTime(left.occurredAt))
    .slice(0, 24);
});

const contextItems = computed(() => [
  {
    key: "customer",
    label: "客户工作台",
    active: true,
    badge: customer.customerNo || "-"
  },
  {
    key: "identity",
    label: "实名审核",
    to: "/customer/identities",
    badge: identities.value.length
  },
  {
    key: "service",
    label: "业务列表",
    to: "/services/list",
    badge: services.value.length
  },
  {
    key: "invoice",
    label: "账单管理",
    to: customer.id ? `/billing/invoices?customerId=${customer.id}` : "/billing/invoices",
    badge: invoices.value.length
  },
  {
    key: "finance",
    label: "资金台账",
    to: customer.id ? `/billing/accounts?customerId=${customer.id}` : "/billing/accounts",
    badge: financeEventCount.value || walletTransactionCount.value
  },
  {
    key: "order",
    label: "产品订单",
    to: "/orders/list"
  }
]);

function pickEntityId(item: RelatedItem, keys: string[]) {
  for (const key of keys) {
    const value = Number((item as Record<string, unknown>)[key]);
    if (Number.isFinite(value) && value > 0) return value;
  }
  return undefined;
}

function resetContactForm() {
  editingContactId.value = null;
  contactForm.name = "";
  contactForm.email = "";
  contactForm.mobile = "";
  contactForm.roleName = "";
  contactForm.isPrimary = false;
}

function normalizeServices(items: RelatedItem[]) {
  return items.map(item => ({
    ...item,
    entityId: pickEntityId(item, ["serviceId", "id"]),
    serviceNo: item.serviceNo ?? item.no ?? "-",
    productName: item.productName ?? item.name ?? "-",
    nextDueAt: item.nextDueAt ?? item.dueAt ?? "-"
  }));
}

function normalizeInvoices(items: RelatedItem[]) {
  return items.map(item => ({
    ...item,
    entityId: pickEntityId(item, ["invoiceId", "id"]),
    invoiceNo: item.invoiceNo ?? item.no ?? "-",
    totalAmount: item.totalAmount ?? item.amount ?? "-",
    dueAt: item.dueAt ?? "-"
  }));
}

function normalizeTickets(items: RelatedItem[]) {
  return items.map(item => ({
    ...item,
    entityId: pickEntityId(item, ["ticketId", "id"]),
    ticketNo: item.ticketNo ?? item.no ?? "-",
    title: item.title ?? item.name ?? "-",
    updatedAt: item.updatedAt ?? "-"
  }));
}

function formatCustomerType(type: string) {
  const mapping: Record<string, string> = {
    COMPANY: "企业客户",
    PERSONAL: "个人客户"
  };
  return mapping[type] ?? type;
}

function formatCustomerStatus(status: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "正常",
    SUSPENDED: "暂停",
    DISABLED: "禁用"
  };
  return mapping[status] ?? status;
}

function customerStatusTag(status: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "success",
    SUSPENDED: "warning",
    DISABLED: "danger"
  };
  return mapping[status] ?? "info";
}

function identityStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "待审核",
    APPROVED: "已实名",
    REJECTED: "已驳回"
  };
  return mapping[status] ?? status;
}

function identityTagType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    APPROVED: "success",
    REJECTED: "danger"
  };
  return mapping[status] ?? "info";
}

function formatAmount(value: string | number | undefined) {
  if (value === undefined || value === null || value === "") return "-";
  const amount = typeof value === "number" ? value : Number(value);
  if (!Number.isFinite(amount)) return String(value);
  return `¥${amount.toFixed(2)}`;
}

function formatWalletMoney(value: string | number | undefined) {
  if (value === undefined || value === null || value === "") return "¥0.00";
  const amount = typeof value === "number" ? value : Number(value);
  return Number.isFinite(amount) ? `¥${amount.toFixed(2)}` : String(value);
}

function transactionTypeLabel(value: string) {
  const mapping: Record<string, string> = {
    RECHARGE: "充值",
    CONSUME: "余额支付",
    REFUND: "退款回退",
    ADJUSTMENT: "手工调整",
    CREDIT_LIMIT: "授信调整"
  };
  return mapping[value] ?? value;
}

function directionTag(value: string) {
  const mapping: Record<string, string> = {
    IN: "success",
    OUT: "danger",
    FLAT: "info"
  };
  return mapping[value] ?? "info";
}

function directionLabel(value: string) {
  const mapping: Record<string, string> = {
    IN: "收入",
    OUT: "支出",
    FLAT: "平移"
  };
  return mapping[value] ?? value;
}

function openServiceDetail(row: WorkbenchRelatedItem) {
  if (!row.entityId) return;
  void router.push(`/services/detail/${row.entityId}`);
}

function openInvoiceDetail(row: WorkbenchRelatedItem) {
  if (!row.entityId) return;
  void router.push(`/billing/invoices/${row.entityId}`);
}

function openInvoiceById(invoiceId: number) {
  const invoice = invoices.value.find(item => item.entityId === invoiceId);
  if (invoice?.entityId) {
    openInvoiceDetail(invoice);
    return;
  }
  void router.push(`/billing/invoices/${invoiceId}`);
}

function openInvoiceWorkbench(status?: "UNPAID" | "PAID" | "REFUNDED") {
  const query: Record<string, string> = {};
  if (customer.id) {
    query.customerId = String(customer.id);
  }
  if (status) {
    query.status = status;
  }
  void router.push({ path: "/billing/invoices", query });
}

function openPaymentWorkbench(keyword?: string) {
  const query: Record<string, string> = {};
  if (customer.id) {
    query.customerId = String(customer.id);
  }
  if (keyword) {
    query.keyword = keyword;
  }
  void router.push({ path: "/billing/payments", query });
}

function openRefundWorkbench(keyword?: string) {
  const query: Record<string, string> = {};
  if (customer.id) {
    query.customerId = String(customer.id);
  }
  if (keyword) {
    query.keyword = keyword;
  }
  void router.push({ path: "/billing/refunds", query });
}

function openFinanceWorkbench(action?: "recharge" | "deduct" | "credit") {
  const query: Record<string, string> = {};
  if (customer.id) {
    query.customerId = String(customer.id);
  }
  if (action) {
    query.action = action;
  }
  void router.push({ path: "/billing/accounts", query });
}

function openRechargeWorkbench(keyword?: string) {
  const query: Record<string, string> = {};
  if (customer.id) {
    query.customerId = String(customer.id);
  }
  if (keyword) {
    query.keyword = keyword;
  }
  void router.push({ path: "/billing/recharges", query });
}

function openAdjustmentWorkbench(keyword?: string, transactionType?: string) {
  const query: Record<string, string> = {};
  if (customer.id) {
    query.customerId = String(customer.id);
  }
  if (keyword) {
    query.keyword = keyword;
  }
  if (transactionType) {
    query.transactionType = transactionType;
  }
  void router.push({ path: "/billing/adjustments", query });
}

function walletActionLabel(row: AccountTransactionRecord) {
  return `${directionLabel(row.direction)} · ${transactionTypeLabel(row.transactionType)}`;
}

function invoiceNoById(invoiceId: number) {
  return invoices.value.find(item => item.entityId === invoiceId)?.invoiceNo ?? `#${invoiceId}`;
}

function paymentStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "已收款"
  };
  return mapping[status] ?? status;
}

function paymentStatusTag(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "success"
  };
  return mapping[status] ?? "info";
}

function refundStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "已退款"
  };
  return mapping[status] ?? status;
}

function refundStatusTag(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "warning"
  };
  return mapping[status] ?? "info";
}

function invoiceStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    UNPAID: "待支付",
    PAID: "已支付",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

function invoiceStatusTag(status: string) {
  const mapping: Record<string, FinanceTimelineItem["type"]> = {
    UNPAID: "warning",
    PAID: "success",
    REFUNDED: "info"
  };
  return mapping[status] ?? "primary";
}

function parseTimelineTime(value: string) {
  const normalized = String(value || "").trim().replace(" ", "T");
  const parsed = Date.parse(normalized);
  return Number.isFinite(parsed) ? parsed : 0;
}

function openFinanceTimelineTarget(item: FinanceTimelineItem) {
  if (item.target === "account") {
    if (item.transactionType === "RECHARGE") {
      openRechargeWorkbench(item.linkKeyword);
      return;
    }
    if (item.transactionType === "ADJUSTMENT" || item.transactionType === "CREDIT_LIMIT") {
      openAdjustmentWorkbench(item.linkKeyword, item.transactionType);
      return;
    }
    openFinanceWorkbench();
    return;
  }
  if (item.target === "payment") {
    openPaymentWorkbench(item.linkKeyword);
    return;
  }
  if (item.target === "refund") {
    openRefundWorkbench(item.linkKeyword);
    return;
  }
  if (item.invoiceId) {
    openInvoiceById(item.invoiceId);
    return;
  }
  openInvoiceWorkbench();
}

async function loadDetail() {
  loading.value = true;
  try {
    const id = Number(customerId.value);
    const [
      customerData,
      contactData,
      identityData,
      serviceData,
      invoiceData,
      paymentData,
      refundData,
      ticketData,
      auditData,
      walletData
    ] =
      await Promise.all([
        fetchCustomerDetail(id),
        fetchCustomerContacts(id),
        fetchCustomerIdentitiesDetail(id),
        fetchCustomerServices(id),
        fetchCustomerInvoices(id),
        fetchPayments({
          customerId: id,
          limit: 20,
          sort: "paid_at",
          order: "desc"
        }),
        fetchRefunds({
          customerId: id,
          limit: 20,
          sort: "created_at",
          order: "desc"
        }),
        fetchCustomerTickets(id),
        fetchCustomerAuditLogs(id),
        fetchCustomerWallet(id).catch(() => null)
      ]);

    Object.assign(customer, customerData);
    contacts.value = contactData;
    identities.value = identityData;
    services.value = normalizeServices(serviceData);
    invoices.value = normalizeInvoices(invoiceData);
    payments.value = paymentData.items;
    refunds.value = refundData.items;
    tickets.value = normalizeTickets(ticketData);
    auditLogs.value = auditData;
    walletDetail.value = walletData;
  } finally {
    loading.value = false;
  }
}

async function saveCustomer() {
  savingCustomer.value = true;
  try {
    const updated = await updateCustomer(customer.id, customer);
    Object.assign(customer, updated);
    ElMessage.success("客户资料已更新");
    await loadDetail();
  } finally {
    savingCustomer.value = false;
  }
}

function openCreateContact() {
  resetContactForm();
  contactDialogVisible.value = true;
}

function openEditContact(contact: Contact) {
  editingContactId.value = contact.id;
  contactForm.name = contact.name;
  contactForm.email = contact.email;
  contactForm.mobile = contact.mobile;
  contactForm.roleName = contact.roleName;
  contactForm.isPrimary = contact.isPrimary;
  contactDialogVisible.value = true;
}

async function submitContact() {
  contactSubmitting.value = true;
  try {
    if (editingContactId.value) {
      await updateCustomerContact(customer.id, editingContactId.value, contactForm);
      ElMessage.success("联系人已更新");
    } else {
      await createCustomerContact(customer.id, contactForm);
      ElMessage.success("联系人已新增");
    }
    contactDialogVisible.value = false;
    resetContactForm();
    await loadDetail();
  } finally {
    contactSubmitting.value = false;
  }
}

async function removeContact(contact: Contact) {
  await ElMessageBox.confirm(`确认删除联系人 ${contact.name} 吗？`, "删除联系人", {
    type: "warning"
  });
  await deleteCustomerContact(customer.id, contact.id);
  ElMessage.success("联系人已删除");
  await loadDetail();
}

async function handleReview(identity: IdentityRecord, status: "APPROVED" | "REJECTED") {
  const remark =
    status === "APPROVED" ? "资料齐全，审核通过" : "资料不完整，请补充后重新提交";
  await reviewCustomerIdentity(customer.id, identity.id, { status, remark });
  ElMessage.success(status === "APPROVED" ? "已通过实名认证" : "已驳回实名认证");
  await loadDetail();
}

watch(
  () => route.params.id,
  () => {
    void loadDetail();
  }
);

onMounted(() => {
  void loadDetail();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="客户"
      title="客户工作台"
      subtitle="把客户资料、联系人、实名、服务、账单、工单和审计记录集中到同一页处理，作为后台客户主工作台。"
    >
      <template #actions>
        <el-button @click="router.push('/customer/list')">返回列表</el-button>
        <el-button @click="loadDetail">刷新详情</el-button>
        <el-button type="primary" plain @click="openFinanceWorkbench()">资金台账</el-button>
        <el-button plain @click="openRechargeWorkbench()">充值记录</el-button>
        <el-button plain @click="openAdjustmentWorkbench()">资金调整</el-button>
        <el-button type="primary" :loading="savingCustomer" @click="saveCustomer">
          保存客户资料
        </el-button>
      </template>

      <template #context>
        <ContextTabs :items="contextItems" />
      </template>

      <template #metrics>
        <div class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户编号</span>
            <strong class="detail-kpi-card__value">{{ customer.customerNo || "N/A" }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户状态</span>
            <el-tag :type="customerStatusTag(customer.status)" effect="light">
              {{ formatCustomerStatus(customer.status) }}
            </el-tag>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户分组</span>
            <strong class="detail-kpi-card__value">{{ customer.groupName || "未设置" }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户等级</span>
            <strong class="detail-kpi-card__value">{{ customer.levelName || "未设置" }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">已实名档案</span>
            <strong class="detail-kpi-card__value">{{ approvedIdentityCount }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">关联服务</span>
            <strong class="detail-kpi-card__value">{{ services.length }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">账户余额</span>
            <strong class="detail-kpi-card__value">{{ formatWalletMoney(walletBalance) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">可用授信</span>
            <strong class="detail-kpi-card__value">{{ formatWalletMoney(walletAvailableCredit) }}</strong>
          </div>
        </div>
      </template>

      <el-tabs>
        <el-tab-pane label="客户概览">
          <div class="customer-grid">
            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>基础资料</strong>
                <span class="section-card__meta">{{ formatCustomerType(customer.type) }}</span>
              </div>
              <el-form label-position="top" class="customer-form-grid">
                <el-form-item label="客户名称">
                  <el-input v-model="customer.name" />
                </el-form-item>
                <el-form-item label="邮箱">
                  <el-input v-model="customer.email" />
                </el-form-item>
                <el-form-item label="手机号">
                  <el-input v-model="customer.mobile" />
                </el-form-item>
                <el-form-item label="销售归属">
                  <el-input v-model="customer.salesOwner" />
                </el-form-item>
                <el-form-item label="客户分组">
                  <el-input v-model="customer.groupName" />
                </el-form-item>
                <el-form-item label="客户等级">
                  <el-input v-model="customer.levelName" />
                </el-form-item>
                <el-form-item label="备注" class="customer-form-grid__full">
                  <el-input v-model="customer.remarks" type="textarea" :rows="4" />
                </el-form-item>
              </el-form>
            </div>

            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>客户摘要</strong>
                <span class="section-card__meta">对象关系总览</span>
              </div>
              <el-descriptions :column="1" border>
                <el-descriptions-item label="客户类型">
                  {{ formatCustomerType(customer.type) }}
                </el-descriptions-item>
                <el-descriptions-item label="联系人数量">
                  {{ contacts.length }}
                </el-descriptions-item>
                <el-descriptions-item label="实名档案">
                  {{ identities.length }}
                </el-descriptions-item>
                <el-descriptions-item label="关联账单">
                  {{ invoices.length }}
                </el-descriptions-item>
                <el-descriptions-item label="支付与退款">
                  {{ financeEventCount }}
                </el-descriptions-item>
                <el-descriptions-item label="关联工单">
                  {{ tickets.length }}
                </el-descriptions-item>
                <el-descriptions-item label="审计记录">
                  {{ auditLogs.length }}
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </div>

          <div class="customer-grid customer-grid--finance" style="margin-top: 16px">
            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>钱包与授信</strong>
                <span class="section-card__meta">{{ walletUpdatedAt }}</span>
              </div>
              <div v-if="walletDetail" class="summary-strip" style="margin-bottom: 16px">
                <div class="summary-pill">
                  <span>账户余额</span>
                  <strong>{{ formatWalletMoney(walletDetail.wallet.balance) }}</strong>
                </div>
                <div class="summary-pill">
                  <span>授信额度</span>
                  <strong>{{ formatWalletMoney(walletDetail.wallet.creditLimit) }}</strong>
                </div>
                <div class="summary-pill">
                  <span>已用授信</span>
                  <strong>{{ formatWalletMoney(walletDetail.wallet.creditUsed) }}</strong>
                </div>
                <div class="summary-pill">
                  <span>可用授信</span>
                  <strong>{{ formatWalletMoney(walletDetail.wallet.availableCredit) }}</strong>
                </div>
              </div>
              <el-empty v-else description="暂无钱包与授信记录" :image-size="72" />
              <div class="inline-actions">
                <el-button type="primary" plain @click="openFinanceWorkbench('recharge')">线下充值</el-button>
                <el-button type="warning" plain @click="openFinanceWorkbench('deduct')">扣减余额</el-button>
                <el-button type="primary" @click="openFinanceWorkbench('credit')">调整授信</el-button>
                <el-button @click="openFinanceWorkbench()">打开资金台账</el-button>
                <el-button plain @click="openRechargeWorkbench()">查看充值记录</el-button>
                <el-button plain @click="openAdjustmentWorkbench()">查看调整记录</el-button>
              </div>
            </div>

            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>最近资金流水</strong>
                <span class="section-card__meta">最近 5 条</span>
              </div>
              <el-table
                v-if="walletDetail?.transactions.length"
                :data="walletDetail.transactions.slice(0, 5)"
                border
                stripe
                size="small"
                empty-text="暂无流水"
              >
                <el-table-column prop="transactionNo" label="流水号" min-width="150" />
                <el-table-column label="动作" min-width="130">
                  <template #default="{ row }">
                    <el-tag :type="directionTag(row.direction)" effect="light">{{ walletActionLabel(row) }}</el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="渠道" min-width="110">
                  <template #default="{ row }">{{ formatPaymentChannel(localeStore.locale, row.channel || "SYSTEM") }}</template>
                </el-table-column>
                <el-table-column label="金额" min-width="110">
                  <template #default="{ row }">{{ formatWalletMoney(row.amount) }}</template>
                </el-table-column>
                <el-table-column prop="summary" label="摘要" min-width="180" show-overflow-tooltip />
                <el-table-column prop="occurredAt" label="时间" min-width="160" />
              </el-table>
              <el-empty v-else description="暂无资金流水" :image-size="72" />
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="联系人">
          <div class="page-header">
            <div>
              <h2 class="section-title">联系人管理</h2>
              <p class="page-subtitle">支持新增、编辑、删除和设置主联系人。</p>
            </div>
            <el-button type="primary" @click="openCreateContact">新增联系人</el-button>
          </div>
          <el-table :data="contacts" border stripe empty-text="暂无联系人">
            <el-table-column prop="name" label="姓名" min-width="140" />
            <el-table-column prop="email" label="邮箱" min-width="180" />
            <el-table-column prop="mobile" label="手机号" min-width="140" />
            <el-table-column prop="roleName" label="角色" min-width="160" />
            <el-table-column label="主联系人" min-width="120">
              <template #default="{ row }">
                <el-tag :type="row.isPrimary ? 'primary' : 'info'" effect="plain">
                  {{ row.isPrimary ? "是" : "否" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" min-width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openEditContact(row)">编辑</el-button>
                <el-button link type="danger" @click="removeContact(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="实名信息">
          <el-table :data="identities" border stripe empty-text="暂无实名档案">
            <el-table-column prop="subjectName" label="主体名称" min-width="220" />
            <el-table-column prop="certNo" label="证件号" min-width="220" />
            <el-table-column prop="identityType" label="证件类型" min-width="120" />
            <el-table-column prop="countryCode" label="国家/地区" min-width="120" />
            <el-table-column prop="submittedAt" label="提交时间" min-width="180" />
            <el-table-column label="状态" min-width="120">
              <template #default="{ row }">
                <el-tag :type="identityTagType(row.verifyStatus)" effect="light">
                  {{ identityStatusLabel(row.verifyStatus) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="reviewRemark" label="审核备注" min-width="220" />
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <el-button
                  link
                  type="primary"
                  :disabled="row.verifyStatus === 'APPROVED'"
                  @click="handleReview(row, 'APPROVED')"
                >
                  通过
                </el-button>
                <el-button
                  link
                  type="danger"
                  :disabled="row.verifyStatus === 'REJECTED'"
                  @click="handleReview(row, 'REJECTED')"
                >
                  驳回
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="关联服务">
          <el-table :data="services" border stripe empty-text="暂无关联服务">
            <el-table-column prop="serviceNo" label="服务编号" min-width="180" />
            <el-table-column prop="productName" label="产品名称" min-width="220" />
            <el-table-column prop="status" label="状态" min-width="120" />
            <el-table-column prop="nextDueAt" label="到期时间" min-width="180" />
            <el-table-column label="操作" min-width="140" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link :disabled="!row.entityId" @click="openServiceDetail(row)">
                  打开服务
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="关联账单">
          <el-table :data="invoices" border stripe empty-text="暂无关联账单">
            <el-table-column prop="invoiceNo" label="账单编号" min-width="180" />
            <el-table-column label="金额" min-width="120">
              <template #default="{ row }">{{ formatAmount(row.totalAmount) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" min-width="120" />
            <el-table-column prop="dueAt" label="到期时间" min-width="180" />
            <el-table-column label="操作" min-width="140" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link :disabled="!row.entityId" @click="openInvoiceDetail(row)">
                  打开账单
                </el-button>
              </template>
            </el-table-column>
          </el-table>
          <div class="inline-actions" style="margin-top: 16px">
            <el-button type="primary" plain @click="openInvoiceWorkbench()">客户账单列表</el-button>
            <el-button type="warning" plain @click="openInvoiceWorkbench('UNPAID')">待支付账单</el-button>
          </div>
        </el-tab-pane>

        <el-tab-pane label="支付与退款">
          <div class="finance-stack">
            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>支付记录</strong>
                <span class="section-card__meta">最近 {{ payments.length }} 笔</span>
              </div>
              <el-table :data="payments" border stripe empty-text="暂无支付记录">
                <el-table-column prop="paymentNo" label="支付单号" min-width="180" />
                <el-table-column label="关联账单" min-width="170">
                  <template #default="{ row }">
                    <el-button link type="primary" @click="openInvoiceById(row.invoiceId)">
                      {{ invoiceNoById(row.invoiceId) }}
                    </el-button>
                  </template>
                </el-table-column>
                <el-table-column label="渠道" min-width="120">
                  <template #default="{ row }">
                    {{ formatPaymentChannel(localeStore.locale, row.channel || "SYSTEM") }}
                  </template>
                </el-table-column>
                <el-table-column label="金额" min-width="120">
                  <template #default="{ row }">{{ formatWalletMoney(row.amount) }}</template>
                </el-table-column>
                <el-table-column label="状态" min-width="120">
                  <template #default="{ row }">
                    <el-tag :type="paymentStatusTag(row.status)" effect="light">
                      {{ paymentStatusLabel(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="tradeNo" label="交易流水" min-width="180" show-overflow-tooltip />
                <el-table-column prop="operator" label="登记人" min-width="120" />
                <el-table-column prop="paidAt" label="支付时间" min-width="170" />
              </el-table>
            </div>

            <div class="page-card inner-card">
              <div class="section-card__head">
                <strong>退款记录</strong>
                <span class="section-card__meta">最近 {{ refunds.length }} 笔</span>
              </div>
              <el-table :data="refunds" border stripe empty-text="暂无退款记录">
                <el-table-column prop="refundNo" label="退款单号" min-width="180" />
                <el-table-column label="关联账单" min-width="170">
                  <template #default="{ row }">
                    <el-button link type="primary" @click="openInvoiceById(row.invoiceId)">
                      {{ invoiceNoById(row.invoiceId) }}
                    </el-button>
                  </template>
                </el-table-column>
                <el-table-column label="金额" min-width="120">
                  <template #default="{ row }">{{ formatWalletMoney(row.amount) }}</template>
                </el-table-column>
                <el-table-column label="状态" min-width="120">
                  <template #default="{ row }">
                    <el-tag :type="refundStatusTag(row.status)" effect="light">
                      {{ refundStatusLabel(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="reason" label="退款原因" min-width="240" show-overflow-tooltip />
                <el-table-column prop="createdAt" label="退款时间" min-width="170" />
              </el-table>
            </div>

            <div class="inline-actions">
              <el-button type="primary" plain @click="openPaymentWorkbench()">查看全部支付</el-button>
              <el-button type="warning" plain @click="openRefundWorkbench()">查看全部退款</el-button>
              <el-button plain @click="openRechargeWorkbench()">查看充值记录</el-button>
              <el-button plain @click="openAdjustmentWorkbench()">查看调整记录</el-button>
              <el-button @click="openInvoiceWorkbench()">联查客户账单</el-button>
            </div>
          </div>
        </el-tab-pane>

        <el-tab-pane label="财务时间线">
          <div class="page-card inner-card">
            <div class="section-card__head">
              <strong>客户财务时间线</strong>
              <span class="section-card__meta">最近 {{ financeTimeline.length }} 条</span>
            </div>
            <el-timeline v-if="financeTimeline.length" class="finance-timeline">
              <el-timeline-item
                v-for="item in financeTimeline"
                :key="item.key"
                :timestamp="item.occurredAt"
                :type="item.type"
                placement="top"
              >
                <div class="timeline-card">
                  <div class="timeline-card__head">
                    <strong>{{ item.title }}</strong>
                    <div class="timeline-card__meta">
                      <el-tag size="small" effect="plain">{{ item.categoryLabel }}</el-tag>
                      <el-tag size="small" effect="plain">{{ item.referenceNo }}</el-tag>
                      <el-tag v-if="item.amount !== undefined" :type="item.type" effect="light">
                        {{ formatWalletMoney(item.amount) }}
                      </el-tag>
                    </div>
                  </div>
                  <div class="timeline-card__description">{{ item.description }}</div>
                  <div class="timeline-card__actions">
                    <el-button link type="primary" @click="openFinanceTimelineTarget(item)">
                      打开关联记录
                    </el-button>
                  </div>
                </div>
              </el-timeline-item>
            </el-timeline>
            <el-empty v-else description="暂无客户财务时间线" :image-size="72" />
          </div>
        </el-tab-pane>

        <el-tab-pane label="关联工单">
          <el-table :data="tickets" border stripe empty-text="暂无关联工单">
            <el-table-column prop="ticketNo" label="工单编号" min-width="180" />
            <el-table-column prop="title" label="工单标题" min-width="260" />
            <el-table-column prop="status" label="状态" min-width="120" />
            <el-table-column prop="updatedAt" label="更新时间" min-width="180" />
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="审计日志">
          <el-table :data="auditLogs" border stripe empty-text="暂无审计记录">
            <el-table-column prop="createdAt" label="时间" min-width="180" />
            <el-table-column prop="actor" label="操作人" min-width="140" />
            <el-table-column prop="action" label="动作" min-width="180" />
            <el-table-column prop="description" label="说明" min-width="280" />
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <el-dialog
        v-model="contactDialogVisible"
        :title="editingContactId ? '编辑联系人' : '新增联系人'"
        width="520px"
      >
        <el-form label-position="top">
          <el-form-item label="姓名">
            <el-input v-model="contactForm.name" />
          </el-form-item>
          <el-form-item label="邮箱">
            <el-input v-model="contactForm.email" />
          </el-form-item>
          <el-form-item label="手机号">
            <el-input v-model="contactForm.mobile" />
          </el-form-item>
          <el-form-item label="角色">
            <el-input v-model="contactForm.roleName" placeholder="例如：技术联系人 / 财务联系人" />
          </el-form-item>
          <el-form-item label="主联系人">
            <el-switch v-model="contactForm.isPrimary" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="contactDialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="contactSubmitting" @click="submitContact">
            保存联系人
          </el-button>
        </template>
      </el-dialog>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.customer-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.35fr) minmax(320px, 0.65fr);
  gap: 16px;
}

.customer-grid--finance {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

.inner-card {
  padding: 18px;
}

.customer-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 12px 16px;
}

.customer-form-grid__full {
  grid-column: 1 / -1;
}

.section-card__head {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  align-items: center;
  margin-bottom: 16px;
}

.section-card__meta {
  color: #5e7093;
  font-size: 12px;
}

.finance-stack {
  display: grid;
  gap: 16px;
}

.finance-timeline {
  padding-top: 8px;
}

.timeline-card {
  display: grid;
  gap: 10px;
  padding: 14px 16px;
  border: 1px solid #dbe5f1;
  border-radius: 14px;
  background: linear-gradient(180deg, #fff, #f8fbff);
}

.timeline-card__head,
.timeline-card__meta,
.timeline-card__actions {
  display: flex;
  align-items: center;
}

.timeline-card__head,
.timeline-card__actions {
  justify-content: space-between;
  gap: 12px;
}

.timeline-card__meta {
  flex-wrap: wrap;
  gap: 8px;
}

.timeline-card__description {
  color: #4f617d;
  line-height: 1.6;
}

@media (max-width: 1280px) {
  .customer-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .customer-form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
