<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import AutomationTaskPanel from "@/components/workbench/AutomationTaskPanel.vue";
import AuditTrailTable from "@/components/workbench/AuditTrailTable.vue";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import {
  fetchOrderDetail,
  fetchProviderAccounts,
  updatePendingOrder,
  type OrderDetailResponse,
  type ProviderAccount
} from "@/api/admin";
import {
  billingCycleOptions as createBillingCycleOptions,
  formatAuditAction,
  formatAuditDescription,
  formatAuditReason,
  formatBillingCycle,
  formatChangeOrderAction,
  formatChangeOrderExecution,
  formatFieldValue,
  formatInvoiceStatus,
  formatMoney,
  formatOrderStatus,
  formatProductType,
  formatProviderType,
  formatServiceStatus,
  fieldLabel,
  orderStatusOptions as createOrderStatusOptions
} from "@/utils/business";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();
const billingCycles = computed(() => createBillingCycleOptions(localeStore.locale));
const orderStatusOptions = computed(() => createOrderStatusOptions(localeStore.locale));
const _legacyBillingCycles = [
  { label: "月付", value: "monthly" },
  { label: "季付", value: "quarterly" },
  { label: "半年付", value: "semiannual" },
  { label: "年付", value: "annual" },
  { label: "两年付", value: "biennially" },
  { label: "三年付", value: "triennially" },
  { label: "一次性", value: "onetime" }
];
const _legacyOrderStatusOptions = [
  { label: "待支付", value: "PENDING" },
  { label: "待交付", value: "ACTIVE" },
  { label: "已完成", value: "COMPLETED" },
  { label: "已取消", value: "CANCELLED" }
];

const loading = ref(false);
const savingEdit = ref(false);
const savingLifecycle = ref(false);
const editDialogVisible = ref(false);
const lifecycleDialogVisible = ref(false);
const detail = ref<OrderDetailResponse | null>(null);
const providerAccounts = ref<ProviderAccount[]>([]);
const editForm = reactive({
  productName: "",
  billingCycle: "monthly",
  amount: 0,
  dueAt: "",
  status: "PENDING",
  reason: ""
});
const lifecycleForm = reactive({
  status: "ACTIVE",
  reason: ""
});

const customerId = computed(() => detail.value?.order.customerId ?? 0);
const primaryInvoice = computed(() => detail.value?.invoices[0] ?? null);
const orderAccount = computed(() => {
  const accountId = detail.value?.order.providerAccountId || detail.value?.services[0]?.providerAccountId || 0;
  if (!accountId) return null;
  return providerAccounts.value.find(item => item.id === accountId) ?? null;
});
const orderEditImpact = computed(() => {
  switch (editForm.status) {
    case "ACTIVE":
      return {
        type: "success" as const,
        title: "将订单改为待交付",
        description: "系统会把关联服务拉回运行中，并保持账单链路为可交付状态。"
      };
    case "COMPLETED":
      return {
        type: "success" as const,
        title: "将订单改为已完成",
        description: "适合人工确认业务已完成的场景，已关联服务会保持生效状态。"
      };
    case "CANCELLED":
      return {
        type: "warning" as const,
        title: "将订单改为已取消",
        description: "系统会把关联服务终止，适合订单关闭、作废或撤销业务的场景。"
      };
    default:
      return {
        type: "info" as const,
        title: "将订单改为待支付",
        description: "系统会把关联服务回退为挂起状态，等待后续收款或再次人工确认。"
      };
  }
});
const lifecycleActionOptions = computed(() =>
  ["PENDING", "ACTIVE", "COMPLETED", "CANCELLED"]
    .filter(status => status !== (detail.value?.order.status || ""))
    .map(status => ({
      value: status,
      label: formatStatus(status)
    }))
);
const lifecycleNotice = computed(() => {
  const payStatus = primaryInvoice.value ? formatInvoiceStatus(localeStore.locale, primaryInvoice.value.status) : "无账单";
  return `当前订单状态 ${formatStatus(detail.value?.order.status || "PENDING")}，主账单 ${payStatus}，关联服务 ${detail.value?.services.length ?? 0} 个。`;
});

const contextTabs = computed(() => [
  { key: "customer", label: "客户", to: customerId.value ? `/customer/detail/${customerId.value}` : undefined },
  { key: "order", label: "订单工作台", active: true, badge: detail.value?.order.orderNo },
  {
    key: "invoice",
    label: "账单",
    to: detail.value?.invoices[0] ? `/billing/invoices/${detail.value.invoices[0].id}` : undefined,
    badge: detail.value?.invoices.length ?? 0
  },
  {
    key: "service",
    label: "服务",
    to: detail.value?.services[0] ? `/services/detail/${detail.value.services[0].id}` : undefined,
    badge: detail.value?.services.length ?? 0
  }
]);

interface OrderTimelineItem {
  key: string;
  time: string;
  title: string;
  description: string;
  amount?: number;
  tag?: string;
  tagType?: "primary" | "success" | "warning" | "info" | "danger";
  routePath?: string;
  routeQuery?: Record<string, string | undefined>;
  linkLabel?: string;
}

interface OrderAuditChangeRow {
  key: string;
  label: string;
  before: string;
  after: string;
}

interface OrderAuditRow {
  id: number;
  createdAt: string;
  actor: string;
  action: string;
  actionLabel: string;
  description: string;
  reason: string;
  changes: OrderAuditChangeRow[];
}

function formatCurrency(value: number) {
  return formatMoney(localeStore.locale, value);
  return `¥${value.toFixed(2)}`;
}

function formatTimelineAmount(value: number) {
  if (value < 0) return `-${formatCurrency(Math.abs(value))}`;
  return formatCurrency(value);
}

function asRecord(value: unknown): Record<string, unknown> {
  if (value && typeof value === "object" && !Array.isArray(value)) {
    return value as Record<string, unknown>;
  }
  return {};
}

function formatCycle(cycle: string) {
  return formatBillingCycle(localeStore.locale, cycle);
  const mapping: Record<string, string> = {
    monthly: "月付",
    quarterly: "季付",
    semiannual: "半年付",
    semiannually: "半年付",
    annual: "年付",
    annually: "年付",
    biennially: "两年付",
    triennially: "三年付",
    onetime: "一次性"
  };
  return mapping[cycle] ?? cycle;
}

function formatStatus(status: string) {
  return formatOrderStatus(localeStore.locale, status);
  const mapping: Record<string, string> = {
    PENDING: "待支付",
    ACTIVE: "待交付",
    PAID: "已支付",
    COMPLETED: "已完成",
    CANCELLED: "已取消"
  };
  return mapping[status] ?? status;
}

function formatAutomationType(type: string) {
  return formatProviderType(localeStore.locale, type);
  const mapping: Record<string, string> = {
    LOCAL: "本地模块",
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上下游财务",
    RESOURCE: "资源池",
    MANUAL: "手动资源"
  };
  return mapping[type] ?? type;
}

function formatProductTypeLabel(type: string) {
  return formatProductType(localeStore.locale, type);
}

function statusTagType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    ACTIVE: "primary",
    PAID: "success",
    COMPLETED: "success",
    CANCELLED: "info"
  };
  return mapping[status] ?? "info";
}

function changeOrderActionLabel(action: string) {
  return formatChangeOrderAction(localeStore.locale, action);
  const mapping: Record<string, string> = {
    "add-ipv4": "新增 IPv4",
    "add-ipv6": "新增 IPv6",
    "add-disk": "新增数据盘",
    "resize-disk": "扩容数据盘",
    "create-snapshot": "创建快照",
    "create-backup": "创建备份"
  };
  return mapping[action] ?? action;
}

function changeOrderExecutionLabel(status: string) {
  return formatChangeOrderExecution(localeStore.locale, status);
  const mapping: Record<string, string> = {
    WAITING_PAYMENT: "待支付",
    PAID: "待回写",
    EXECUTING: "执行中",
    EXECUTED: "已执行",
    EXECUTE_FAILED: "执行失败",
    EXECUTE_BLOCKED: "执行阻塞",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

function changeOrderExecutionType(status: string) {
  const mapping: Record<string, string> = {
    WAITING_PAYMENT: "warning",
    PAID: "info",
    EXECUTING: "primary",
    EXECUTED: "success",
    EXECUTE_FAILED: "danger",
    EXECUTE_BLOCKED: "warning",
    REFUNDED: "info"
  };
  return mapping[status] ?? "info";
}

async function loadProviderAccounts() {
  providerAccounts.value = await fetchProviderAccounts();
}

async function loadDetail() {
  loading.value = true;
  try {
    detail.value = await fetchOrderDetail(route.params.id as string);
  } finally {
    loading.value = false;
  }
}

function openEditDialog() {
  if (!detail.value) return;
  editForm.productName = detail.value.order.productName;
  editForm.billingCycle = detail.value.order.billingCycle;
  editForm.amount = detail.value.order.amount;
  editForm.dueAt = detail.value.invoices[0]?.dueAt ?? "";
  editForm.status = detail.value.order.status;
  editForm.reason = "";
  editDialogVisible.value = true;
}

async function submitEditOrder() {
  if (!detail.value) return;
  savingEdit.value = true;
  try {
    const result = await updatePendingOrder(detail.value.order.id, {
      productName: editForm.productName,
      billingCycle: editForm.billingCycle,
      amount: editForm.amount,
      dueAt: editForm.dueAt,
      status: editForm.status,
      reason: editForm.reason
    });
    detail.value = result;
    editDialogVisible.value = false;
    ElMessage.success("订单和关联账单已更新，变更记录已写入审计");
  } finally {
    savingEdit.value = false;
  }
}

function defaultLifecycleReason(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "后台回退为待支付，等待重新收款或人工复核",
    ACTIVE: "后台确认订单进入待开通流程",
    COMPLETED: "后台确认订单已完成交付",
    CANCELLED: "后台关闭订单并终止后续交付"
  };
  return mapping[status] ?? "后台调整订单生命周期状态";
}

function openLifecycleDialog(status: string) {
  lifecycleForm.status = status;
  lifecycleForm.reason = defaultLifecycleReason(status);
  lifecycleDialogVisible.value = true;
}

function closeLifecycleDialog() {
  lifecycleDialogVisible.value = false;
  lifecycleForm.reason = "";
}

async function submitLifecycleAction() {
  if (!detail.value) return;
  savingLifecycle.value = true;
  try {
    const result = await updatePendingOrder(detail.value.order.id, {
      status: lifecycleForm.status,
      reason: lifecycleForm.reason.trim() || defaultLifecycleReason(lifecycleForm.status)
    });
    detail.value = result;
    lifecycleDialogVisible.value = false;
    ElMessage.success(`订单已更新为${formatStatus(lifecycleForm.status)}`);
  } finally {
    savingLifecycle.value = false;
  }
}

function parseTimelineTime(value: string) {
  if (!value) return 0;
  const parsed = Date.parse(value.includes("T") ? value : value.replace(" ", "T"));
  return Number.isFinite(parsed) ? parsed : 0;
}

function providerAccountLabel(accountId?: number) {
  if (!accountId) return "未绑定";
  const account = providerAccounts.value.find(item => item.id === accountId);
  if (!account) return "未绑定";
  return `${account.name} / ${account.baseUrl}`;
}

function openCustomerWorkbench() {
  if (!detail.value) return;
  void router.push(`/customer/detail/${detail.value.order.customerId}`);
}

function openProductWorkbench() {
  if (!detail.value) return;
  void router.push(`/catalog/products/${detail.value.order.productId}`);
}

function openProviderAccountWorkbench() {
  if (!detail.value) return;
  const providerType = detail.value.services[0]?.providerType || detail.value.order.automationType;
  const accountId = detail.value.services[0]?.providerAccountId || detail.value.order.providerAccountId;
  void router.push({
    path: "/providers/accounts",
    query: {
      providerType: providerType || undefined,
      accountId: accountId ? String(accountId) : undefined
    }
  });
}

function openInvoiceWorkbench(invoiceId?: number) {
  const id = invoiceId || detail.value?.invoices[0]?.id;
  if (!id) return;
  void router.push(`/billing/invoices/${id}`);
}

function openServiceWorkbench(serviceId?: number) {
  const id = serviceId || detail.value?.services[0]?.id;
  if (!id) return;
  void router.push(`/services/detail/${id}`);
}

function openProviderAutomationContext(options?: {
  invoiceId?: number;
  serviceId?: number;
  channel?: string;
}) {
  if (!detail.value) return;
  const matchedService = options?.serviceId
    ? detail.value.services.find(item => item.id === options.serviceId)
    : detail.value.services[0];
  void router.push({
    path: "/providers/automation",
    query: {
      orderId: String(detail.value.order.id),
      invoiceId: options?.invoiceId ? String(options.invoiceId) : undefined,
      serviceId: options?.serviceId ? String(options.serviceId) : undefined,
      channel: options?.channel || matchedService?.providerType || detail.value.order.automationType || undefined
    }
  });
}

function openProviderAutomationWorkbench() {
  openProviderAutomationContext();
}

function openProviderResourcesContext(options?: {
  serviceId?: number;
  channel?: string;
  accountId?: number;
}) {
  if (!detail.value) return;
  const matchedService = options?.serviceId
    ? detail.value.services.find(item => item.id === options.serviceId)
    : detail.value.services[0];
  void router.push({
    path: "/providers/resources",
    query: {
      providerType: options?.channel || matchedService?.providerType || detail.value.order.automationType || undefined,
      accountId: options?.accountId
        ? String(options.accountId)
        : matchedService?.providerAccountId
          ? String(matchedService.providerAccountId)
          : undefined,
      serviceId: options?.serviceId ? String(options.serviceId) : matchedService?.id ? String(matchedService.id) : undefined
    }
  });
}

function openProviderResourcesWorkbench() {
  openProviderResourcesContext();
}

function openOrderTimelineTarget(item: OrderTimelineItem) {
  if (!item.routePath) return;
  void router.push({
    path: item.routePath,
    query: item.routeQuery
  });
}

const orderTimeline = computed<OrderTimelineItem[]>(() => {
  if (!detail.value) return [];

  const items: OrderTimelineItem[] = [
    {
      key: `order-${detail.value.order.id}`,
      time: detail.value.order.createdAt,
      title: `订单创建 ${detail.value.order.orderNo}`,
      description: `商品 ${detail.value.order.productName}，状态 ${formatStatus(detail.value.order.status)}，自动化渠道 ${formatAutomationType(detail.value.order.automationType)}`,
      amount: detail.value.order.amount,
      tag: detail.value.order.orderNo,
      tagType: statusTagType(detail.value.order.status) as OrderTimelineItem["tagType"],
      routePath: `/customer/detail/${detail.value.order.customerId}`,
      linkLabel: "打开客户工作台"
    }
  ];

  detail.value.invoices.forEach(invoice => {
    items.push({
      key: `invoice-${invoice.id}`,
      time: invoice.paidAt || invoice.dueAt,
      title: `账单联动 ${invoice.invoiceNo}`,
      description: `账单状态 ${formatInvoiceStatus(localeStore.locale, invoice.status)}，金额 ${formatCurrency(invoice.totalAmount)}，到期 ${invoice.dueAt || "-"}`,
      amount: invoice.totalAmount,
      tag: formatInvoiceStatus(localeStore.locale, invoice.status),
      tagType: invoice.status === "PAID" ? "success" : invoice.status === "UNPAID" ? "warning" : "info",
      routePath: `/billing/invoices/${invoice.id}`,
      linkLabel: "打开账单工作台"
    });
  });

  detail.value.services.forEach(service => {
    items.push({
      key: `service-${service.id}`,
      time: service.updatedAt || service.lastSyncAt || service.createdAt,
      title: `服务交付 ${service.serviceNo}`,
      description: `当前状态 ${formatServiceStatus(localeStore.locale, service.status)}，接口账户 ${providerAccountLabel(service.providerAccountId)}，下次到期 ${service.nextDueAt || "-"}`,
      tag: formatAutomationType(service.providerType),
      tagType: "primary",
      routePath: `/services/detail/${service.id}`,
      linkLabel: "打开服务工作台"
    });
  });

  if (detail.value.changeOrder) {
    items.push({
      key: `change-order-${detail.value.changeOrder.id}`,
      time: detail.value.changeOrder.paidAt || detail.value.changeOrder.createdAt,
      title: `改配执行 ${detail.value.changeOrder.invoiceNo}`,
      description: `动作 ${changeOrderActionLabel(detail.value.changeOrder.actionName)}，支付状态 ${formatInvoiceStatus(localeStore.locale, detail.value.changeOrder.status)}，执行状态 ${changeOrderExecutionLabel(detail.value.changeOrder.executionStatus)}`,
      amount: detail.value.changeOrder.amount,
      tag: detail.value.changeOrder.executionStatus,
      tagType: changeOrderExecutionType(detail.value.changeOrder.executionStatus) as OrderTimelineItem["tagType"],
      routePath: "/providers/automation",
      routeQuery: {
        orderId: String(detail.value.order.id),
        invoiceId: detail.value.changeOrder.invoiceId ? String(detail.value.changeOrder.invoiceId) : undefined,
        serviceId: detail.value.changeOrder.serviceId ? String(detail.value.changeOrder.serviceId) : undefined,
        channel: detail.value.services[0]?.providerType || detail.value.order.automationType || undefined
      },
      linkLabel: "查看改配执行任务"
    });
  }

  detail.value.auditLogs.slice(0, 6).forEach(log => {
    items.push({
      key: `audit-${log.id}`,
      time: log.createdAt,
      title: `审计记录 ${formatAuditAction(localeStore.locale, log.action)}`,
      description: `${formatAuditDescription(localeStore.locale, log.action, log.description)}｜操作人 ${log.actor || "-"}${log.target ? `｜对象 ${log.target}` : ""}`,
      tag: log.actor || "系统",
      tagType: "info"
    });
  });

  return items.sort((left, right) => parseTimelineTime(right.time) - parseTimelineTime(left.time));
});

function buildAuditChanges(item: OrderDetailResponse["auditLogs"][number]): OrderAuditChangeRow[] {
  const payload = asRecord(item.payload);
  const before = asRecord(payload.before);
  const after = asRecord(payload.after);
  const keys = Array.from(new Set([...Object.keys(before), ...Object.keys(after)]));
  return keys
    .filter(
      key =>
        formatFieldValue(localeStore.locale, key, before[key], item.action) !==
        formatFieldValue(localeStore.locale, key, after[key], item.action)
    )
    .map(key => ({
      key,
      label: fieldLabel(localeStore.locale, key),
      before: formatFieldValue(localeStore.locale, key, before[key], item.action),
      after: formatFieldValue(localeStore.locale, key, after[key], item.action)
    }));
}

const orderAuditRows = computed<OrderAuditRow[]>(() =>
  (detail.value?.auditLogs ?? []).map(item => ({
    id: item.id,
    createdAt: item.createdAt,
    actor: item.actor || "系统",
    action: item.action,
    actionLabel: formatAuditAction(localeStore.locale, item.action),
    description: formatAuditDescription(localeStore.locale, item.action, item.description),
    reason: formatAuditReason(asRecord(item.payload).reason),
    changes: buildAuditChanges(item)
  }))
);

const priceRelatedKeys = new Set(["amount", "billingCycle", "productName", "dueAt"]);
const lifecycleRelatedKeys = new Set(["status"]);

const pricingAuditRows = computed(() =>
  orderAuditRows.value.filter(row => row.changes.some(change => priceRelatedKeys.has(change.key)))
);

const lifecycleAuditRows = computed(() =>
  orderAuditRows.value.filter(
    row =>
      row.action === "order.checkout" ||
      row.changes.some(change => lifecycleRelatedKeys.has(change.key))
  )
);

const auditSummary = computed(() => ({
  total: orderAuditRows.value.length,
  manualAdjustments: orderAuditRows.value.filter(row => row.action === "order.manual_adjust").length,
  priceAdjustments: pricingAuditRows.value.length,
  lifecycleChanges: lifecycleAuditRows.value.length,
  latestOperator: orderAuditRows.value[0]?.actor || "系统"
}));

function summarizeAuditChanges(row: OrderAuditRow) {
  if (row.changes.length === 0) return "-";
  return row.changes
    .slice(0, 3)
    .map(change => `${change.label} ${change.before} -> ${change.after}`)
    .join("；");
}

watch(() => route.params.id, () => void loadDetail());

onMounted(async () => {
  await Promise.all([loadProviderAccounts(), loadDetail()]);
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 订单"
      title="订单工作台"
      subtitle="围绕订单对象集中处理客户、账单、服务交付、接口账户与自动化任务。"
    >
      <template #actions>
        <el-button @click="router.push('/orders/list')">返回列表</el-button>
        <el-button @click="loadDetail">刷新详情</el-button>
        <el-button v-if="detail" type="primary" plain @click="openEditDialog">
          人工调整订单
        </el-button>
      </template>

      <template #context>
        <ContextTabs v-if="detail" :items="contextTabs" />
      </template>

      <template #metrics>
        <div v-if="detail" class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">订单编号</span>
            <strong class="detail-kpi-card__value">{{ detail.order.orderNo }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">客户</span>
            <strong class="detail-kpi-card__value">{{ detail.order.customerName }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">订单金额</span>
            <strong class="detail-kpi-card__value">{{ formatCurrency(detail.order.amount) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">订单状态</span>
            <el-tag :type="statusTagType(detail.order.status)" effect="light">
              {{ formatStatus(detail.order.status) }}
            </el-tag>
          </div>
        </div>
      </template>

      <template v-if="detail">
        <el-tabs>
          <el-tab-pane label="订单概览">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>基础信息</strong>
                  <span class="section-card__meta">{{ formatProductTypeLabel(detail.order.productType) }}</span>
                </div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="商品名称">{{ detail.order.productName }}</el-descriptions-item>
                  <el-descriptions-item label="计费周期">{{ formatCycle(detail.order.billingCycle) }}</el-descriptions-item>
                  <el-descriptions-item label="自动化渠道">
                    {{ formatAutomationType(detail.order.automationType) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="接口账户">
                    {{ orderAccount ? `${orderAccount.name} / ${orderAccount.baseUrl}` : "未绑定" }}
                  </el-descriptions-item>
                  <el-descriptions-item label="创建时间">{{ detail.order.createdAt }}</el-descriptions-item>
                  <el-descriptions-item label="订单编号">{{ detail.order.orderNo }}</el-descriptions-item>
                </el-descriptions>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>交付链路</strong>
                  <span class="section-card__meta">订单 -> 账单 -> 服务 -> 自动化任务</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill">
                    <span>关联账单</span>
                    <strong>{{ detail.invoices.length }}</strong>
                  </div>
                  <div class="summary-pill">
                    <span>已开通服务</span>
                    <strong>{{ detail.services.length }}</strong>
                  </div>
                  <div class="summary-pill">
                    <span>审计记录</span>
                    <strong>{{ detail.auditLogs.length }}</strong>
                  </div>
                </div>
                <el-alert
                  title="订单工作台负责串联支付、交付与同步状态。当前接口账户会直接决定后续自动化执行落到哪套魔方或上下游环境。"
                  type="info"
                  :closable="false"
                  show-icon
                />
                <div class="inline-actions" style="margin-top: 16px">
                  <el-button plain @click="openCustomerWorkbench">客户</el-button>
                  <el-button plain @click="openProductWorkbench">商品</el-button>
                  <el-button plain @click="openProviderAccountWorkbench">接口账户</el-button>
                  <el-button plain @click="openProviderAutomationWorkbench">自动化任务</el-button>
                  <el-button plain :disabled="detail.services.length === 0" @click="openProviderResourcesWorkbench">渠道资源</el-button>
                </div>
                <div style="margin-top: 20px">
                  <div class="section-card__head">
                    <strong>生命周期处置</strong>
                    <span class="section-card__meta">当前状态 {{ formatStatus(detail.order.status) }}</span>
                  </div>
                  <el-alert
                    :title="lifecycleNotice"
                    description="常用处置包括回退待支付、转待开通、确认完成和订单取消，动作会写入审计并回写关联账单/服务状态。"
                    type="warning"
                    :closable="false"
                    show-icon
                  />
                  <div class="inline-actions" style="margin-top: 16px">
                    <el-button
                      v-for="item in lifecycleActionOptions"
                      :key="item.value"
                      plain
                      @click="openLifecycleDialog(item.value)"
                    >
                      {{ item.label }}
                    </el-button>
                  </div>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="业务时间线">
            <div class="table-toolbar">
              <div class="table-toolbar__meta">
                <strong>订单业务时间线</strong>
                <span>共 {{ orderTimeline.length }} 条对象事件</span>
              </div>
              <div class="inline-actions">
                <el-button type="primary" link @click="openProviderAutomationWorkbench">任务中心</el-button>
                <el-button type="primary" link @click="openProviderResourcesWorkbench">渠道资源</el-button>
              </div>
            </div>
            <el-timeline v-if="orderTimeline.length">
              <el-timeline-item
                v-for="item in orderTimeline"
                :key="item.key"
                :timestamp="item.time || '未记录时间'"
                placement="top"
                :type="item.tagType || 'info'"
              >
                <div class="panel-card">
                  <div class="section-card__head">
                    <strong>{{ item.title }}</strong>
                    <div style="display: flex; gap: 8px; align-items: center; flex-wrap: wrap">
                      <el-tag v-if="item.tag" :type="item.tagType || 'info'" effect="light">{{ item.tag }}</el-tag>
                      <strong v-if="item.amount !== undefined">{{ formatTimelineAmount(item.amount) }}</strong>
                    </div>
                  </div>
                  <div style="color: var(--el-text-color-regular); line-height: 1.7">
                    {{ item.description }}
                  </div>
                  <div v-if="item.linkLabel" class="inline-actions" style="margin-top: 12px">
                    <el-button type="primary" link @click="openOrderTimelineTarget(item)">{{ item.linkLabel }}</el-button>
                  </div>
                </div>
              </el-timeline-item>
            </el-timeline>
            <el-empty v-else description="当前订单暂无可展示的业务时间线" :image-size="80" />
          </el-tab-pane>

          <el-tab-pane label="关联账单">
            <el-table :data="detail.invoices" border stripe empty-text="暂无关联账单">
              <el-table-column prop="invoiceNo" label="账单编号" min-width="180" />
              <el-table-column label="金额" min-width="120">
                <template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template>
              </el-table-column>
              <el-table-column label="状态" min-width="120">
                <template #default="{ row }">{{ formatInvoiceStatus(localeStore.locale, row.status) }}</template>
              </el-table-column>
              <el-table-column prop="dueAt" label="到期时间" min-width="180" />
              <el-table-column label="操作" min-width="220" fixed="right">
                <template #default="{ row }">
                  <div class="inline-actions">
                    <el-button type="primary" link @click="openInvoiceWorkbench(row.id)">打开账单</el-button>
                    <el-button type="primary" link @click="openProviderAutomationContext({ invoiceId: row.id })">任务中心</el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="已开通服务">
            <el-table :data="detail.services" border stripe empty-text="暂无关联服务">
              <el-table-column prop="serviceNo" label="服务编号" min-width="170" />
              <el-table-column prop="productName" label="产品名称" min-width="220" />
              <el-table-column label="渠道" min-width="140">
                <template #default="{ row }">{{ formatAutomationType(row.providerType) }}</template>
              </el-table-column>
              <el-table-column label="接口账户" min-width="220">
                <template #default="{ row }">{{ providerAccountLabel(row.providerAccountId) }}</template>
              </el-table-column>
              <el-table-column label="状态" min-width="120">
                <template #default="{ row }">{{ formatServiceStatus(localeStore.locale, row.status) }}</template>
              </el-table-column>
              <el-table-column prop="createdAt" label="开通时间" min-width="180" />
              <el-table-column prop="nextDueAt" label="下次到期" min-width="160" />
              <el-table-column label="操作" min-width="260" fixed="right">
                <template #default="{ row }">
                  <div class="inline-actions">
                    <el-button type="primary" link @click="openServiceWorkbench(row.id)">打开服务</el-button>
                    <el-button type="primary" link @click="openProviderAutomationContext({ serviceId: row.id, channel: row.providerType })">任务中心</el-button>
                    <el-button
                      type="primary"
                      link
                      @click="openProviderResourcesContext({ serviceId: row.id, channel: row.providerType, accountId: row.providerAccountId })"
                    >
                      渠道资源
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane v-if="detail.changeOrder" label="改配记录">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="改配订单">{{ detail.changeOrder.orderNo }}</el-descriptions-item>
              <el-descriptions-item label="改配账单">{{ detail.changeOrder.invoiceNo }}</el-descriptions-item>
              <el-descriptions-item label="资源动作">
                {{ changeOrderActionLabel(detail.changeOrder.actionName) }}
              </el-descriptions-item>
              <el-descriptions-item label="支付状态">{{ formatInvoiceStatus(localeStore.locale, detail.changeOrder.status) }}</el-descriptions-item>
              <el-descriptions-item label="执行状态">
                <el-tag :type="changeOrderExecutionType(detail.changeOrder.executionStatus)" effect="light">
                  {{ changeOrderExecutionLabel(detail.changeOrder.executionStatus) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="改配金额">{{ formatCurrency(detail.changeOrder.amount) }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detail.changeOrder.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="已支付时间">{{ detail.changeOrder.paidAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="执行结果" :span="2">
                {{ detail.changeOrder.executionMessage || "-" }}
              </el-descriptions-item>
            </el-descriptions>
            <div class="inline-actions" style="margin-top: 16px">
              <el-button type="primary" link @click="router.push(`/orders/change-orders?orderId=${detail.changeOrder.orderId}`)">改配单工作台</el-button>
              <el-button type="primary" link @click="openInvoiceWorkbench(detail.changeOrder.invoiceId)">打开改配账单</el-button>
              <el-button type="primary" link @click="openServiceWorkbench(detail.changeOrder.serviceId)">打开关联服务</el-button>
              <el-button
                type="primary"
                link
                @click="openProviderAutomationContext({ serviceId: detail.changeOrder.serviceId, invoiceId: detail.changeOrder.invoiceId })"
              >
                查看执行任务
              </el-button>
              <el-button type="primary" link @click="openProviderResourcesContext({ serviceId: detail.changeOrder.serviceId })">
                渠道资源
              </el-button>
            </div>
          </el-tab-pane>

          <el-tab-pane label="自动化任务">
            <AutomationTaskPanel title="订单自动化任务" :order-id="detail.order.id" />
          </el-tab-pane>

          <el-tab-pane label="处置记录">
            <div class="summary-strip" style="margin-bottom: 16px">
              <div class="summary-pill">
                <span>审计总数</span>
                <strong>{{ auditSummary.total }}</strong>
              </div>
              <div class="summary-pill">
                <span>人工调整</span>
                <strong>{{ auditSummary.manualAdjustments }}</strong>
              </div>
              <div class="summary-pill">
                <span>改价记录</span>
                <strong>{{ auditSummary.priceAdjustments }}</strong>
              </div>
              <div class="summary-pill">
                <span>状态变更</span>
                <strong>{{ auditSummary.lifecycleChanges }}</strong>
              </div>
              <div class="summary-pill">
                <span>最近操作人</span>
                <strong>{{ auditSummary.latestOperator }}</strong>
              </div>
            </div>

            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>价格与配置变更</strong>
                  <span class="section-card__meta">记录金额、周期、商品名、到期时间的调整</span>
                </div>
                <el-table :data="pricingAuditRows" border stripe empty-text="暂无价格或配置调整记录">
                  <el-table-column prop="createdAt" label="时间" min-width="170" />
                  <el-table-column prop="actor" label="操作人" min-width="120" />
                  <el-table-column prop="actionLabel" label="动作" min-width="160" />
                  <el-table-column label="变更摘要" min-width="320" show-overflow-tooltip>
                    <template #default="{ row }">{{ summarizeAuditChanges(row) }}</template>
                  </el-table-column>
                  <el-table-column prop="reason" label="原因" min-width="220" show-overflow-tooltip />
                </el-table>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>生命周期记录</strong>
                  <span class="section-card__meta">记录订单状态流转和下单动作</span>
                </div>
                <el-table :data="lifecycleAuditRows" border stripe empty-text="暂无生命周期处置记录">
                  <el-table-column prop="createdAt" label="时间" min-width="170" />
                  <el-table-column prop="actor" label="操作人" min-width="120" />
                  <el-table-column prop="actionLabel" label="动作" min-width="160" />
                  <el-table-column prop="description" label="说明" min-width="240" show-overflow-tooltip />
                  <el-table-column label="状态摘要" min-width="320" show-overflow-tooltip>
                    <template #default="{ row }">{{ summarizeAuditChanges(row) }}</template>
                  </el-table-column>
                  <el-table-column prop="reason" label="原因" min-width="220" show-overflow-tooltip />
                </el-table>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="审计日志">
            <AuditTrailTable :items="detail.auditLogs" empty-text="暂无审计记录" />
          </el-tab-pane>
        </el-tabs>
      </template>
    </PageWorkbench>

    <el-dialog v-model="lifecycleDialogVisible" title="生命周期处置" width="520px" @closed="closeLifecycleDialog">
      <el-form label-position="top">
        <el-alert
          :title="`将订单更新为${formatStatus(lifecycleForm.status)}`"
          description="这个动作会走后台人工调整链路，并同步把变更写入审计日志。若订单已经关联服务，服务交付状态也会一并回写。"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="目标状态">
          <el-select v-model="lifecycleForm.status" style="width: 100%">
            <el-option v-for="item in orderStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="变更原因">
          <el-input
            v-model="lifecycleForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入后台处置原因，例如人工收款后转待开通、客户取消订单、人工确认已交付等"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeLifecycleDialog">取消</el-button>
        <el-button type="primary" :loading="savingLifecycle" @click="submitLifecycleAction">确认处置</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="editDialogVisible" title="人工调整订单" width="560px">
      <el-form label-position="top">
        <el-alert
          :title="orderEditImpact.title"
          :description="orderEditImpact.description"
          :type="orderEditImpact.type"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="商品名称">
          <el-input v-model="editForm.productName" placeholder="请输入商品名称" />
        </el-form-item>
        <el-form-item label="计费周期">
          <el-select v-model="editForm.billingCycle" style="width: 100%">
            <el-option v-for="item in billingCycles" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="待支付金额">
          <el-input-number v-model="editForm.amount" :min="0" :precision="2" :step="10" style="width: 100%" />
        </el-form-item>
        <el-form-item label="订单状态">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option v-for="item in orderStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="账单到期时间">
          <el-date-picker
            v-model="editForm.dueAt"
            type="datetime"
            value-format="YYYY-MM-DD HH:mm:ss"
            placeholder="请选择到期时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="变更原因">
          <el-input
            v-model="editForm.reason"
            type="textarea"
            :rows="3"
            placeholder="例如：客户申请改周期、财务复核改价、销售补贴调整"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingEdit" @click="submitEditOrder">保存调整</el-button>
      </template>
    </el-dialog>
  </div>
</template>
