<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import {
  fetchServiceChangeOrders,
  type ServiceChangeOrderQuery,
  type ServiceChangeOrderRecord
} from "@/api/admin";
import {
  formatChangeOrderAction,
  formatChangeOrderExecution,
  formatInvoiceStatus,
  formatMoney,
  formatProviderType
} from "@/utils/business";
import { useLocaleStore } from "@/store";

type TabKey = "ALL" | "UNPAID" | "PAID" | "REFUNDED";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const activeTab = ref<TabKey>("ALL");
const selectedRows = ref<ServiceChangeOrderRecord[]>([]);
const items = ref<ServiceChangeOrderRecord[]>([]);
const total = ref(0);
const detailVisible = ref(false);
const detailItem = ref<ServiceChangeOrderRecord | null>(null);

const pagination = reactive({
  page: 1,
  limit: 20
});

const filters = reactive({
  keyword: "",
  action: "",
  serviceId: "",
  orderId: "",
  invoiceId: "",
  executionStatus: ""
});

const actionOptions = [
  { label: "全部动作", value: "" },
  { label: "新增 IPv4", value: "add-ipv4" },
  { label: "新增 IPv6", value: "add-ipv6" },
  { label: "新增磁盘", value: "add-disk" },
  { label: "扩容磁盘", value: "resize-disk" },
  { label: "创建快照", value: "create-snapshot" },
  { label: "创建备份", value: "create-backup" }
];

const executionOptions = [
  { label: "全部执行状态", value: "" },
  { label: "待支付", value: "WAITING_PAYMENT" },
  { label: "已支付待执行", value: "PAID" },
  { label: "执行中", value: "EXECUTING" },
  { label: "执行完成", value: "EXECUTED" },
  { label: "执行失败", value: "EXECUTE_FAILED" },
  { label: "执行阻塞", value: "EXECUTE_BLOCKED" },
  { label: "已退款", value: "REFUNDED" }
];

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: total.value },
  { key: "UNPAID", label: "待支付", count: items.value.filter(item => item.status === "UNPAID").length },
  { key: "PAID", label: "已支付", count: items.value.filter(item => item.status === "PAID").length },
  { key: "REFUNDED", label: "已退款", count: items.value.filter(item => item.status === "REFUNDED").length }
]);

const totalAmount = computed(() => items.value.reduce((sum, item) => sum + item.amount, 0));
const executedCount = computed(() => items.value.filter(item => item.executionStatus === "EXECUTED").length);
const failedCount = computed(() =>
  items.value.filter(item => ["EXECUTE_FAILED", "EXECUTE_BLOCKED"].includes(item.executionStatus)).length
);
const waitingCount = computed(() =>
  items.value.filter(item => ["WAITING_PAYMENT", "PAID", "EXECUTING"].includes(item.executionStatus)).length
);
const filteredItems = computed(() =>
  items.value.filter(item => !filters.executionStatus || item.executionStatus === filters.executionStatus)
);
const contextEntries = computed(() => {
  const entries: Array<{ label: string; value: string }> = [];
  if (filters.serviceId) entries.push({ label: "服务", value: filters.serviceId });
  if (filters.orderId) entries.push({ label: "订单", value: filters.orderId });
  if (filters.invoiceId) entries.push({ label: "账单", value: filters.invoiceId });
  if (filters.action) entries.push({ label: "动作", value: filters.action });
  if (filters.executionStatus) entries.push({ label: "执行状态", value: filters.executionStatus });
  return entries;
});

function formatCurrency(value: number) {
  return formatMoney(localeStore.locale, value);
}

function paymentStatusType(value: string) {
  if (value === "PAID") return "success";
  if (value === "REFUNDED") return "info";
  return "warning";
}

function executionStatusType(value: string) {
  const mapping: Record<string, string> = {
    WAITING_PAYMENT: "warning",
    PAID: "info",
    EXECUTING: "primary",
    EXECUTED: "success",
    EXECUTE_FAILED: "danger",
    EXECUTE_BLOCKED: "warning",
    REFUNDED: "info"
  };
  return mapping[value] ?? "info";
}

function readRouteQueryValue(value: unknown) {
  if (Array.isArray(value)) return String(value[0] ?? "");
  if (value === undefined || value === null) return "";
  return String(value);
}

function syncFiltersFromRoute() {
  activeTab.value = (readRouteQueryValue(route.query.status).toUpperCase() as TabKey) || "ALL";
  filters.keyword = readRouteQueryValue(route.query.keyword);
  filters.action = readRouteQueryValue(route.query.action);
  filters.serviceId = readRouteQueryValue(route.query.serviceId);
  filters.orderId = readRouteQueryValue(route.query.orderId);
  filters.invoiceId = readRouteQueryValue(route.query.invoiceId);
  filters.executionStatus = readRouteQueryValue(route.query.executionStatus).toUpperCase();
  pagination.page = 1;
}

function buildQuery(): ServiceChangeOrderQuery {
  const query: ServiceChangeOrderQuery = {
    page: pagination.page,
    limit: pagination.limit,
    sort: "created_at",
    order: "desc"
  };
  if (activeTab.value !== "ALL") query.status = activeTab.value;
  if (filters.keyword.trim()) query.keyword = filters.keyword.trim();
  if (filters.action) query.action = filters.action;
  if (filters.serviceId.trim()) query.serviceId = filters.serviceId.trim();
  if (filters.orderId.trim()) query.orderId = filters.orderId.trim();
  if (filters.invoiceId.trim()) query.invoiceId = filters.invoiceId.trim();
  return query;
}

async function loadItems() {
  loading.value = true;
  try {
    const data = await fetchServiceChangeOrders(buildQuery());
    items.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.action = "";
  filters.serviceId = "";
  filters.orderId = "";
  filters.invoiceId = "";
  filters.executionStatus = "";
  pagination.page = 1;
  void router.push({ path: "/orders/change-orders", query: {} });
}

function applyFilters() {
  pagination.page = 1;
  void loadItems();
}

function handlePageChange(page: number) {
  pagination.page = page;
  void loadItems();
}

function handleLimitChange(limit: number) {
  pagination.limit = limit;
  pagination.page = 1;
  void loadItems();
}

function handleSelectionChange(rows: ServiceChangeOrderRecord[]) {
  selectedRows.value = rows;
}

function openServiceDetail(row: ServiceChangeOrderRecord) {
  void router.push(`/services/detail/${row.serviceId}`);
}

function openOrderDetail(row: ServiceChangeOrderRecord) {
  void router.push(`/orders/detail/${row.orderId}`);
}

function openInvoiceDetail(row: ServiceChangeOrderRecord) {
  void router.push(`/billing/invoices/${row.invoiceId}`);
}

function openAutomationWorkbench(row: ServiceChangeOrderRecord) {
  void router.push({
    path: "/providers/automation",
    query: {
      orderId: String(row.orderId),
      invoiceId: String(row.invoiceId),
      serviceId: String(row.serviceId),
      channel: row.providerType || undefined
    }
  });
}

function openResourcesWorkbench(row: ServiceChangeOrderRecord) {
  void router.push({
    path: "/providers/resources",
    query: {
      serviceId: String(row.serviceId),
      providerType: row.providerType || undefined
    }
  });
}

function openDetailDrawer(row: ServiceChangeOrderRecord) {
  detailItem.value = row;
  detailVisible.value = true;
}

function closeDetailDrawer() {
  detailVisible.value = false;
  detailItem.value = null;
}

function openExceptionWorkbench() {
  void router.push({
    path: "/providers/automation",
    query: {
      serviceId: filters.serviceId || undefined,
      orderId: filters.orderId || undefined,
      invoiceId: filters.invoiceId || undefined,
      channel: detailItem.value?.providerType || undefined,
      status: "FAILED"
    }
  });
}

function safeJson(value: unknown) {
  if (value === undefined || value === null || value === "") return "-";
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

async function copySelectedNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择改配单");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.invoiceNo || `改配单#${item.id}`).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个改配账单号`);
}

function exportRows(rows: ServiceChangeOrderRecord[], filename: string) {
  downloadCsv(
    filename,
    ["服务编号", "商品名称", "改配账单", "改配订单", "动作", "金额", "支付状态", "执行状态", "渠道", "创建时间"],
    rows.map(item => [
      item.serviceNo || String(item.serviceId),
      item.productName || "-",
      item.invoiceNo || String(item.invoiceId),
      item.orderNo || String(item.orderId),
      formatChangeOrderAction(localeStore.locale, item.actionName),
      formatCurrency(item.amount),
      formatInvoiceStatus(localeStore.locale, item.status),
      formatChangeOrderExecution(localeStore.locale, item.executionStatus),
      formatProviderType(localeStore.locale, item.providerType || "LOCAL"),
      item.createdAt
    ])
  );
}

function exportCurrent() {
  exportRows(filteredItems.value, `change-orders-page-${pagination.page}.csv`);
  ElMessage.success("当前页改配单已导出");
}

onMounted(() => {
  syncFiltersFromRoute();
  void loadItems();
});

watch(activeTab, () => {
  pagination.page = 1;
  void loadItems();
});

watch(
  () => route.fullPath,
  () => {
    syncFiltersFromRoute();
    void loadItems();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 订单"
      title="改配单工作台"
      subtitle="集中管理服务升级、扩容、附加资源等改配订单，联动订单、账单、服务和自动化执行链。"
    >
      <template #actions>
        <el-button @click="loadItems">刷新列表</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>改配单总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatCurrency(totalAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>执行完成</span>
            <strong>{{ executedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>待跟进</span>
            <strong>{{ waitingCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>执行异常</span>
            <strong>{{ failedCount }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="账单号 / 订单号 / 服务号 / 标题" clearable />
          <el-select v-model="filters.action" placeholder="改配动作" clearable>
            <el-option v-for="item in actionOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="filters.executionStatus" placeholder="执行状态" clearable>
            <el-option v-for="item in executionOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-input v-model="filters.serviceId" placeholder="服务 ID" clearable />
          <el-input v-model="filters.orderId" placeholder="订单 ID" clearable />
          <el-input v-model="filters.invoiceId" placeholder="账单 ID" clearable />
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div v-if="contextEntries.length > 0" class="table-toolbar">
          <div class="table-toolbar__meta">
            <strong>当前上下文</strong>
            <span v-for="item in contextEntries" :key="`${item.label}-${item.value}`">{{ item.label }}：{{ item.value }}</span>
          </div>
          <div class="action-group">
            <el-button v-if="filters.serviceId" plain @click="router.push(`/services/detail/${filters.serviceId}`)">打开服务</el-button>
            <el-button v-if="filters.orderId" plain @click="router.push(`/orders/detail/${filters.orderId}`)">打开订单</el-button>
            <el-button v-if="filters.invoiceId" plain @click="router.push(`/billing/invoices/${filters.invoiceId}`)">打开账单</el-button>
            <el-button
              v-if="filters.executionStatus === 'EXECUTE_FAILED' || filters.executionStatus === 'EXECUTE_BLOCKED'"
              type="warning"
              plain
              @click="openExceptionWorkbench"
            >
              联查异常任务
            </el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>改配运营列表</strong>
          <span>当前页 {{ filteredItems.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain @click="copySelectedNos">复制改配单号</el-button>
          <el-button plain @click="exportCurrent">导出当前页</el-button>
        </div>
      </div>

      <el-table :data="filteredItems" border stripe empty-text="暂无改配单" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="serviceNo" label="服务编号" min-width="160" />
        <el-table-column prop="productName" label="商品名称" min-width="220" show-overflow-tooltip />
        <el-table-column prop="title" label="改配标题" min-width="220" show-overflow-tooltip />
        <el-table-column label="动作" min-width="140">
          <template #default="{ row }">{{ formatChangeOrderAction(localeStore.locale, row.actionName) }}</template>
        </el-table-column>
        <el-table-column label="金额" min-width="120">
          <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="支付状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="paymentStatusType(row.status)" effect="light">
              {{ formatInvoiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="执行状态" min-width="140">
          <template #default="{ row }">
            <el-tag :type="executionStatusType(row.executionStatus)" effect="light">
              {{ formatChangeOrderExecution(localeStore.locale, row.executionStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="渠道" min-width="140">
          <template #default="{ row }">{{ formatProviderType(localeStore.locale, row.providerType || "LOCAL") }}</template>
        </el-table-column>
        <el-table-column prop="invoiceNo" label="改配账单" min-width="170" />
        <el-table-column prop="orderNo" label="改配订单" min-width="170" />
        <el-table-column label="最新任务" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">
            {{
              row.latestTask
                ? `${row.latestTask.taskNo} / ${row.latestTask.status}`
                : "未生成执行任务"
            }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="180" />
        <el-table-column label="操作" min-width="340" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetailDrawer(row)">详情</el-button>
              <el-button type="primary" link @click="openServiceDetail(row)">服务</el-button>
              <el-button type="primary" link @click="openOrderDetail(row)">订单</el-button>
              <el-button type="primary" link @click="openInvoiceDetail(row)">账单</el-button>
              <el-button type="primary" link @click="openAutomationWorkbench(row)">任务</el-button>
              <el-button type="primary" link @click="openResourcesWorkbench(row)">资源</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-pagination">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next"
          :total="total"
          :current-page="pagination.page"
          :page-size="pagination.limit"
          :page-sizes="[20, 50, 100]"
          @current-change="handlePageChange"
          @size-change="handleLimitChange"
        />
      </div>
    </PageWorkbench>

    <el-drawer v-model="detailVisible" size="720px" :with-header="false" @closed="closeDetailDrawer">
      <div v-if="detailItem" class="drawer-panel">
        <div class="page-header">
          <div>
            <h3 class="page-header__title">改配单详情</h3>
            <p class="page-header__desc">查看改配原因、执行状态、最新任务和请求内容。</p>
          </div>
          <div class="inline-actions">
            <el-button plain @click="openServiceDetail(detailItem)">服务</el-button>
            <el-button plain @click="openOrderDetail(detailItem)">订单</el-button>
            <el-button plain @click="openInvoiceDetail(detailItem)">账单</el-button>
            <el-button plain @click="openAutomationWorkbench(detailItem)">任务</el-button>
            <el-button plain @click="openResourcesWorkbench(detailItem)">资源</el-button>
          </div>
        </div>

        <div class="portal-grid portal-grid--two">
          <div class="panel-card">
            <div class="section-card__head">
              <strong>基础信息</strong>
              <span class="section-card__meta">{{ detailItem.invoiceNo || `改配单 #${detailItem.id}` }}</span>
            </div>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="服务编号">{{ detailItem.serviceNo || "-" }}</el-descriptions-item>
              <el-descriptions-item label="商品名称">{{ detailItem.productName || "-" }}</el-descriptions-item>
              <el-descriptions-item label="改配动作">
                {{ formatChangeOrderAction(localeStore.locale, detailItem.actionName) }}
              </el-descriptions-item>
              <el-descriptions-item label="渠道">
                {{ formatProviderType(localeStore.locale, detailItem.providerType || "LOCAL") }}
              </el-descriptions-item>
              <el-descriptions-item label="支付状态">
                <el-tag :type="paymentStatusType(detailItem.status)" effect="light">
                  {{ formatInvoiceStatus(localeStore.locale, detailItem.status) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="执行状态">
                <el-tag :type="executionStatusType(detailItem.executionStatus)" effect="light">
                  {{ formatChangeOrderExecution(localeStore.locale, detailItem.executionStatus) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="金额">{{ formatCurrency(detailItem.amount) }}</el-descriptions-item>
              <el-descriptions-item label="计费周期">{{ detailItem.billingCycle || "-" }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detailItem.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="已支付时间">{{ detailItem.paidAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="已退款时间">{{ detailItem.refundedAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="最近更新">{{ detailItem.updatedAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="执行结果" :span="2">{{ detailItem.executionMessage || "-" }}</el-descriptions-item>
              <el-descriptions-item label="改配原因" :span="2">{{ detailItem.reason || "-" }}</el-descriptions-item>
            </el-descriptions>
          </div>

          <div class="panel-card">
            <div class="section-card__head">
              <strong>最新执行任务</strong>
              <span class="section-card__meta">
                {{ detailItem.latestTask ? detailItem.latestTask.taskNo : "未生成任务" }}
              </span>
            </div>
            <el-descriptions v-if="detailItem.latestTask" :column="2" border>
              <el-descriptions-item label="任务编号">{{ detailItem.latestTask.taskNo }}</el-descriptions-item>
              <el-descriptions-item label="状态">{{ detailItem.latestTask.status }}</el-descriptions-item>
              <el-descriptions-item label="阶段">{{ detailItem.latestTask.stage || "-" }}</el-descriptions-item>
              <el-descriptions-item label="渠道">{{ detailItem.latestTask.channel || "-" }}</el-descriptions-item>
              <el-descriptions-item label="执行方">{{ detailItem.latestTask.operatorName || "-" }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ detailItem.latestTask.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="消息" :span="2">{{ detailItem.latestTask.message || "-" }}</el-descriptions-item>
            </el-descriptions>
            <el-empty v-else description="当前改配单还没有下发执行任务" :image-size="72" />
          </div>
        </div>

        <div class="panel-card" style="margin-top: 16px">
          <div class="section-card__head">
            <strong>改配 Payload</strong>
            <span class="section-card__meta">查看本次改配请求内容和参数</span>
          </div>
          <pre class="json-block">{{ safeJson(detailItem.payload) }}</pre>
        </div>
      </div>
    </el-drawer>
  </div>
</template>
