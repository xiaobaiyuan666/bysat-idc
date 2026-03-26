<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import {
  fetchOrderRequests,
  processOrderRequest,
  type OrderRequestQuery,
  type OrderRequestRecord
} from "@/api/admin";
import { formatBillingCycle, formatMoney } from "@/utils/business";
import { useLocaleStore } from "@/store";

type TabKey = "ALL" | "PENDING" | "APPROVED" | "REJECTED" | "COMPLETED" | "CANCELLED";
type RequestProcessStatus = "APPROVED" | "REJECTED" | "COMPLETED" | "CANCELLED";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const items = ref<OrderRequestRecord[]>([]);
const total = ref(0);
const selectedRows = ref<OrderRequestRecord[]>([]);
const activeTab = ref<TabKey>("ALL");
const detailVisible = ref(false);
const detailItem = ref<OrderRequestRecord | null>(null);
const processDialogVisible = ref(false);
const batchDialogVisible = ref(false);
const processing = ref(false);
const batchProcessing = ref(false);
const activeRequest = ref<OrderRequestRecord | null>(null);
const batchRows = ref<OrderRequestRecord[]>([]);

const pagination = reactive({
  page: 1,
  limit: 20
});

const filters = reactive({
  keyword: "",
  type: "",
  orderId: "",
  customerId: ""
});

const summary = reactive({
  total: 0,
  pending: 0,
  approved: 0,
  rejected: 0,
  completed: 0,
  cancelled: 0
});

const processForm = reactive({
  status: "APPROVED" as RequestProcessStatus,
  processNote: ""
});

const batchForm = reactive({
  status: "APPROVED" as RequestProcessStatus,
  processNote: ""
});

const requestTypeOptions = [
  { label: "全部类型", value: "" },
  { label: "取消申请", value: "CANCEL" },
  { label: "续费请求", value: "RENEW" },
  { label: "改价申请", value: "PRICE_ADJUST" }
];

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: summary.total },
  { key: "PENDING", label: "待处理", count: summary.pending },
  { key: "APPROVED", label: "已同意", count: summary.approved },
  { key: "REJECTED", label: "已驳回", count: summary.rejected },
  { key: "COMPLETED", label: "已完成", count: summary.completed },
  { key: "CANCELLED", label: "已撤销", count: summary.cancelled }
]);

const totalRequestedAmount = computed(() =>
  items.value.reduce((sum, item) => sum + Number(item.requestedAmount || 0), 0)
);
const currentPagePending = computed(() => items.value.filter(item => item.status === "PENDING").length);
const currentPageProcessed = computed(() => items.value.filter(item => item.status !== "PENDING").length);
const currentPagePriceAdjust = computed(() => items.value.filter(item => item.type === "PRICE_ADJUST").length);
const contextEntries = computed(() => {
  const entries: Array<{ label: string; value: string }> = [];
  if (filters.orderId.trim()) entries.push({ label: "订单", value: filters.orderId.trim() });
  if (filters.customerId.trim()) entries.push({ label: "客户", value: filters.customerId.trim() });
  if (filters.type) entries.push({ label: "类型", value: formatRequestType(filters.type) });
  if (filters.keyword.trim()) entries.push({ label: "关键字", value: filters.keyword.trim() });
  return entries;
});
const processOptions = computed(() => {
  if (!activeRequest.value) {
    return [
      { label: "同意申请", value: "APPROVED" },
      { label: "驳回申请", value: "REJECTED" }
    ];
  }
  if (activeRequest.value.status === "APPROVED") {
    return [
      { label: "处理完成", value: "COMPLETED" },
      { label: "撤销申请", value: "CANCELLED" }
    ];
  }
  return [
    { label: "同意申请", value: "APPROVED" },
    { label: "驳回申请", value: "REJECTED" },
    { label: "撤销申请", value: "CANCELLED" }
  ];
});
const batchTargetCount = computed(() => batchRows.value.length);
const batchTargetLabel = computed(() => formatRequestStatus(batchForm.status));

function readRouteQueryValue(value: unknown) {
  if (Array.isArray(value)) return String(value[0] ?? "");
  if (value === undefined || value === null) return "";
  return String(value);
}

function syncFiltersFromRoute() {
  const status = readRouteQueryValue(route.query.status).toUpperCase();
  activeTab.value = (status as TabKey) || "ALL";
  filters.keyword = readRouteQueryValue(route.query.keyword);
  filters.type = readRouteQueryValue(route.query.type).toUpperCase();
  filters.orderId = readRouteQueryValue(route.query.orderId);
  filters.customerId = readRouteQueryValue(route.query.customerId || route.query.uid);
  pagination.page = 1;
}

function buildQuery(): OrderRequestQuery {
  const query: OrderRequestQuery = {
    page: pagination.page,
    limit: pagination.limit,
    sort: "created_at",
    order: "desc"
  };
  if (activeTab.value !== "ALL") query.status = activeTab.value;
  if (filters.keyword.trim()) query.keyword = filters.keyword.trim();
  if (filters.type) query.type = filters.type;
  if (filters.orderId.trim()) query.orderId = filters.orderId.trim();
  if (filters.customerId.trim()) query.customerId = filters.customerId.trim();
  return query;
}

function buildSummaryQuery(status?: Exclude<TabKey, "ALL">): OrderRequestQuery {
  const query: OrderRequestQuery = {
    limit: 1,
    page: 1
  };
  if (status) query.status = status;
  if (filters.keyword.trim()) query.keyword = filters.keyword.trim();
  if (filters.type) query.type = filters.type;
  if (filters.orderId.trim()) query.orderId = filters.orderId.trim();
  if (filters.customerId.trim()) query.customerId = filters.customerId.trim();
  return query;
}

async function fetchCount(status?: Exclude<TabKey, "ALL">) {
  const response = await fetchOrderRequests(buildSummaryQuery(status));
  return response.total;
}

async function loadItemsAndSummary() {
  loading.value = true;
  try {
    const [listData, totalCount, pendingCount, approvedCount, rejectedCount, completedCount, cancelledCount] =
      await Promise.all([
        fetchOrderRequests(buildQuery()),
        fetchCount(),
        fetchCount("PENDING"),
        fetchCount("APPROVED"),
        fetchCount("REJECTED"),
        fetchCount("COMPLETED"),
        fetchCount("CANCELLED")
      ]);
    items.value = listData.items;
    total.value = listData.total;
    summary.total = totalCount;
    summary.pending = pendingCount;
    summary.approved = approvedCount;
    summary.rejected = rejectedCount;
    summary.completed = completedCount;
    summary.cancelled = cancelledCount;
  } finally {
    loading.value = false;
  }
}

function handlePageChange(page: number) {
  pagination.page = page;
  void loadItemsAndSummary();
}

function handleLimitChange(limit: number) {
  pagination.limit = limit;
  pagination.page = 1;
  void loadItemsAndSummary();
}

function applyFilters() {
  pagination.page = 1;
  void loadItemsAndSummary();
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.type = "";
  filters.orderId = "";
  filters.customerId = "";
  pagination.page = 1;
  detailVisible.value = false;
  void router.push({ path: "/orders/requests", query: {} });
}

function handleSelectionChange(rows: OrderRequestRecord[]) {
  selectedRows.value = rows;
}

function formatCurrency(value: number) {
  return formatMoney(localeStore.locale, value);
}

function formatRequestType(value: string) {
  const mapping: Record<string, string> = {
    CANCEL: "取消申请",
    RENEW: "续费请求",
    PRICE_ADJUST: "改价申请"
  };
  return mapping[value] ?? value;
}

function formatRequestStatus(value: string) {
  const mapping: Record<string, string> = {
    PENDING: "待处理",
    APPROVED: "已同意",
    REJECTED: "已驳回",
    COMPLETED: "已完成",
    CANCELLED: "已撤销"
  };
  return mapping[value] ?? value;
}

function requestStatusTagType(value: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    APPROVED: "primary",
    REJECTED: "danger",
    COMPLETED: "success",
    CANCELLED: "info"
  };
  return mapping[value] ?? "info";
}

function formatRequestTarget(row: OrderRequestRecord) {
  return `${formatBillingCycle(localeStore.locale, row.currentBillingCycle)} / ${formatCurrency(row.currentAmount)} -> ${formatBillingCycle(localeStore.locale, row.requestedBillingCycle)} / ${formatCurrency(row.requestedAmount)}`;
}

function defaultProcessNote(request: OrderRequestRecord, status: RequestProcessStatus) {
  if (status === "APPROVED") return `后台已同意${formatRequestType(request.type)}`;
  if (status === "REJECTED") return `后台驳回${formatRequestType(request.type)}，待客户重新确认`;
  if (status === "COMPLETED") return `${formatRequestType(request.type)}已完成后续处理`;
  if (status === "CANCELLED") return `${formatRequestType(request.type)}已撤销`;
  return "";
}

function defaultBatchProcessNote(status: RequestProcessStatus) {
  const mapping: Record<RequestProcessStatus, string> = {
    APPROVED: "后台批量同意订单申请，已进入后续执行链路",
    REJECTED: "后台批量驳回订单申请，等待重新确认",
    COMPLETED: "后台批量完成订单申请处理",
    CANCELLED: "后台批量撤销订单申请"
  };
  return mapping[status];
}

function openOrderDetail(row: OrderRequestRecord) {
  void router.push(`/orders/detail/${row.orderId}`);
}

function openCustomerDetail(row: OrderRequestRecord) {
  void router.push(`/customer/detail/${row.customerId}`);
}

function openAutomationWorkbench(row: OrderRequestRecord) {
  void router.push({
    path: "/providers/automation",
    query: {
      orderId: String(row.orderId)
    }
  });
}

function openDetailDrawer(row: OrderRequestRecord) {
  detailItem.value = row;
  detailVisible.value = true;
}

function closeDetailDrawer() {
  detailVisible.value = false;
  detailItem.value = null;
}

function closeProcessDialog() {
  processDialogVisible.value = false;
  activeRequest.value = null;
  processForm.processNote = "";
}

function openProcessDialog(row: OrderRequestRecord, status?: RequestProcessStatus) {
  activeRequest.value = row;
  processForm.status = status || (row.status === "APPROVED" ? "COMPLETED" : "APPROVED");
  processForm.processNote = defaultProcessNote(row, processForm.status);
  processDialogVisible.value = true;
}

async function submitProcess() {
  if (!activeRequest.value) return;
  processing.value = true;
  try {
    const result = await processOrderRequest(activeRequest.value.id, {
      status: processForm.status,
      processNote: processForm.processNote.trim()
    });
    const updatedRequest = result.requests.find(item => item.id === activeRequest.value?.id) ?? null;
    if (updatedRequest) {
      detailItem.value = detailItem.value?.id === updatedRequest.id ? updatedRequest : detailItem.value;
    }
    processDialogVisible.value = false;
    ElMessage.success(`订单申请已更新为${formatRequestStatus(processForm.status)}`);
    await loadItemsAndSummary();
  } finally {
    processing.value = false;
  }
}

function isCompatibleForBatch(row: OrderRequestRecord, status: RequestProcessStatus) {
  if (status === "COMPLETED") return row.status === "APPROVED";
  if (status === "APPROVED" || status === "REJECTED") return row.status === "PENDING";
  return row.status === "PENDING" || row.status === "APPROVED";
}

function openBatchDialog(status: RequestProcessStatus) {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择订单申请");
    return;
  }
  const compatible = selectedRows.value.filter(row => isCompatibleForBatch(row, status));
  const skipped = selectedRows.value.length - compatible.length;
  if (compatible.length === 0) {
    ElMessage.warning("所选申请与当前批量处理动作不匹配");
    return;
  }
  if (skipped > 0) {
    ElMessage.info(`已跳过 ${skipped} 条不适配的申请`);
  }
  batchRows.value = compatible;
  batchForm.status = status;
  batchForm.processNote = defaultBatchProcessNote(status);
  batchDialogVisible.value = true;
}

function closeBatchDialog() {
  batchDialogVisible.value = false;
  batchRows.value = [];
  batchForm.processNote = "";
}

async function submitBatchProcess() {
  if (batchRows.value.length === 0) return;
  const targets = batchRows.value.filter(row => isCompatibleForBatch(row, batchForm.status));
  if (targets.length === 0) {
    ElMessage.warning("当前批量处理状态与所选申请不匹配");
    return;
  }
  batchProcessing.value = true;
  let success = 0;
  let failed = 0;
  try {
    for (const row of targets) {
      try {
        await processOrderRequest(row.id, {
          status: batchForm.status,
          processNote: batchForm.processNote.trim() || defaultBatchProcessNote(batchForm.status)
        });
        success += 1;
      } catch {
        failed += 1;
      }
    }
    if (success > 0) {
      ElMessage.success(
        failed > 0
          ? `已处理 ${success} 条申请，失败 ${failed} 条`
          : `已将 ${success} 条申请更新为${batchTargetLabel.value}`
      );
    } else {
      ElMessage.error("批量处理失败");
    }
    closeBatchDialog();
    selectedRows.value = [];
    await loadItemsAndSummary();
  } finally {
    batchProcessing.value = false;
  }
}

async function copySelectedRequestNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择订单申请");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.requestNo).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个申请单号`);
}

function exportRows(rows: OrderRequestRecord[], filename: string) {
  downloadCsv(
    filename,
    [
      "申请单号",
      "订单号",
      "客户",
      "商品",
      "类型",
      "状态",
      "摘要",
      "当前计费",
      "目标计费",
      "来源",
      "处理人",
      "申请时间",
      "处理时间"
    ],
    rows.map(item => [
      item.requestNo,
      item.orderNo,
      item.customerName,
      item.productName,
      formatRequestType(item.type),
      formatRequestStatus(item.status),
      item.summary,
      `${formatBillingCycle(localeStore.locale, item.currentBillingCycle)} / ${formatCurrency(item.currentAmount)}`,
      `${formatBillingCycle(localeStore.locale, item.requestedBillingCycle)} / ${formatCurrency(item.requestedAmount)}`,
      item.sourceName || "-",
      item.processorName || "-",
      item.createdAt || "-",
      item.processedAt || "-"
    ])
  );
}

function exportCurrent() {
  exportRows(items.value, `order-requests-page-${pagination.page}.csv`);
  ElMessage.success("当前页订单申请已导出");
}

function exportSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择订单申请");
    return;
  }
  exportRows(selectedRows.value, "order-requests-selected.csv");
  ElMessage.success("已导出选中订单申请");
}

function safeJson(value: unknown) {
  if (value === undefined || value === null || value === "") return "-";
  try {
    return JSON.stringify(value, null, 2);
  } catch {
    return String(value);
  }
}

onMounted(() => {
  syncFiltersFromRoute();
  void loadItemsAndSummary();
});

watch(activeTab, () => {
  pagination.page = 1;
  void loadItemsAndSummary();
});

watch(
  () => route.fullPath,
  () => {
    syncFiltersFromRoute();
    void loadItemsAndSummary();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 订单"
      title="订单申请工作台"
      subtitle="集中处理取消、续费、改价等订单申请，联动订单工作台、客户工作台和自动化任务中心。"
    >
      <template #actions>
        <el-button @click="loadItemsAndSummary">刷新列表</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>申请总数</span>
            <strong>{{ summary.total }}</strong>
          </div>
          <div class="summary-pill">
            <span>待处理</span>
            <strong>{{ summary.pending }}</strong>
          </div>
          <div class="summary-pill">
            <span>已完成</span>
            <strong>{{ summary.completed }}</strong>
          </div>
          <div class="summary-pill">
            <span>本页申请金额</span>
            <strong>{{ formatCurrency(totalRequestedAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>本页待处理</span>
            <strong>{{ currentPagePending }}</strong>
          </div>
          <div class="summary-pill">
            <span>本页改价申请</span>
            <strong>{{ currentPagePriceAdjust }}</strong>
          </div>
          <div class="summary-pill">
            <span>本页已处理</span>
            <strong>{{ currentPageProcessed }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="申请单号 / 订单号 / 客户 / 商品 / 摘要" clearable />
          <el-select v-model="filters.type" placeholder="申请类型" clearable>
            <el-option v-for="item in requestTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-input v-model="filters.orderId" placeholder="订单 ID" clearable />
          <el-input v-model="filters.customerId" placeholder="客户 ID" clearable />
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
            <el-button v-if="filters.orderId" plain @click="router.push(`/orders/detail/${filters.orderId}`)">打开订单</el-button>
            <el-button v-if="filters.customerId" plain @click="router.push(`/customer/detail/${filters.customerId}`)">
              打开客户
            </el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>申请运营列表</strong>
          <span>当前页 {{ items.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchDialog('APPROVED')">批量同意</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchDialog('REJECTED')">批量驳回</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchDialog('COMPLETED')">批量完成</el-button>
          <el-button plain @click="copySelectedRequestNos">复制申请单号</el-button>
          <el-button plain @click="exportSelected">导出选中</el-button>
          <el-button plain @click="exportCurrent">导出当前页</el-button>
        </div>
      </div>

      <el-table :data="items" border stripe empty-text="暂无订单申请" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="requestNo" label="申请单号" min-width="160" />
        <el-table-column label="类型" min-width="120">
          <template #default="{ row }">{{ formatRequestType(row.type) }}</template>
        </el-table-column>
        <el-table-column prop="orderNo" label="订单号" min-width="150" />
        <el-table-column prop="customerName" label="客户" min-width="160" show-overflow-tooltip />
        <el-table-column prop="productName" label="商品" min-width="200" show-overflow-tooltip />
        <el-table-column prop="summary" label="摘要" min-width="220" show-overflow-tooltip />
        <el-table-column label="申请内容" min-width="320" show-overflow-tooltip>
          <template #default="{ row }">{{ formatRequestTarget(row) }}</template>
        </el-table-column>
        <el-table-column label="状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="requestStatusTagType(row.status)" effect="light">
              {{ formatRequestStatus(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sourceName" label="来源" min-width="140" show-overflow-tooltip />
        <el-table-column prop="processorName" label="处理人" min-width="120" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="申请时间" min-width="170" />
        <el-table-column prop="processedAt" label="处理时间" min-width="170" />
        <el-table-column label="操作" min-width="360" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetailDrawer(row)">详情</el-button>
              <el-button type="primary" link @click="openOrderDetail(row)">订单</el-button>
              <el-button type="primary" link @click="openCustomerDetail(row)">客户</el-button>
              <el-button type="primary" link @click="openAutomationWorkbench(row)">任务</el-button>
              <el-button v-if="row.status === 'PENDING'" type="primary" link @click="openProcessDialog(row, 'APPROVED')">
                同意
              </el-button>
              <el-button v-if="row.status === 'PENDING'" type="danger" link @click="openProcessDialog(row, 'REJECTED')">
                驳回
              </el-button>
              <el-button v-if="row.status === 'APPROVED'" type="success" link @click="openProcessDialog(row, 'COMPLETED')">
                完成
              </el-button>
              <el-button type="primary" link @click="openProcessDialog(row)">处理</el-button>
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
            <h3 class="page-header__title">订单申请详情</h3>
            <p class="page-header__desc">查看申请原因、目标变更、处理结果和关联订单上下文。</p>
          </div>
          <div class="inline-actions">
            <el-button plain @click="openOrderDetail(detailItem)">订单</el-button>
            <el-button plain @click="openCustomerDetail(detailItem)">客户</el-button>
            <el-button plain @click="openAutomationWorkbench(detailItem)">任务</el-button>
            <el-button plain @click="openProcessDialog(detailItem)">处理申请</el-button>
          </div>
        </div>

        <div class="portal-grid portal-grid--two">
          <div class="panel-card">
            <div class="section-card__head">
              <strong>基础信息</strong>
              <span class="section-card__meta">{{ detailItem.requestNo }}</span>
            </div>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="申请单号">{{ detailItem.requestNo }}</el-descriptions-item>
              <el-descriptions-item label="订单号">{{ detailItem.orderNo }}</el-descriptions-item>
              <el-descriptions-item label="客户">{{ detailItem.customerName }}</el-descriptions-item>
              <el-descriptions-item label="商品">{{ detailItem.productName }}</el-descriptions-item>
              <el-descriptions-item label="申请类型">{{ formatRequestType(detailItem.type) }}</el-descriptions-item>
              <el-descriptions-item label="申请状态">
                <el-tag :type="requestStatusTagType(detailItem.status)" effect="light">
                  {{ formatRequestStatus(detailItem.status) }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="申请时间">{{ detailItem.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="处理时间">{{ detailItem.processedAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="申请来源">{{ detailItem.sourceName || "-" }}</el-descriptions-item>
              <el-descriptions-item label="处理人">{{ detailItem.processorName || "-" }}</el-descriptions-item>
              <el-descriptions-item label="摘要" :span="2">{{ detailItem.summary || "-" }}</el-descriptions-item>
              <el-descriptions-item label="申请原因" :span="2">{{ detailItem.reason || "-" }}</el-descriptions-item>
              <el-descriptions-item label="处理备注" :span="2">{{ detailItem.processNote || "-" }}</el-descriptions-item>
            </el-descriptions>
          </div>

          <div class="panel-card">
            <div class="section-card__head">
              <strong>目标变更</strong>
              <span class="section-card__meta">当前值与目标值对比</span>
            </div>
            <el-descriptions :column="2" border>
              <el-descriptions-item label="当前周期">
                {{ formatBillingCycle(localeStore.locale, detailItem.currentBillingCycle) }}
              </el-descriptions-item>
              <el-descriptions-item label="目标周期">
                {{ formatBillingCycle(localeStore.locale, detailItem.requestedBillingCycle) }}
              </el-descriptions-item>
              <el-descriptions-item label="当前金额">{{ formatCurrency(detailItem.currentAmount) }}</el-descriptions-item>
              <el-descriptions-item label="目标金额">{{ formatCurrency(detailItem.requestedAmount) }}</el-descriptions-item>
              <el-descriptions-item label="联查入口" :span="2">
                <div class="inline-actions">
                  <el-button type="primary" link @click="openOrderDetail(detailItem)">订单工作台</el-button>
                  <el-button type="primary" link @click="openCustomerDetail(detailItem)">客户工作台</el-button>
                  <el-button type="primary" link @click="openAutomationWorkbench(detailItem)">任务中心</el-button>
                </div>
              </el-descriptions-item>
            </el-descriptions>
          </div>
        </div>

        <div class="panel-card" style="margin-top: 16px">
          <div class="section-card__head">
            <strong>Payload</strong>
            <span class="section-card__meta">保留本次申请的原始上下文</span>
          </div>
          <pre class="json-block">{{ safeJson(detailItem.payload) }}</pre>
        </div>
      </div>
    </el-drawer>

    <el-dialog v-model="processDialogVisible" title="处理订单申请" width="560px" @closed="closeProcessDialog">
      <el-form v-if="activeRequest" label-position="top">
        <el-alert
          :title="`${activeRequest.requestNo} · ${formatRequestType(activeRequest.type)}`"
          :description="`${activeRequest.summary}｜当前状态 ${formatRequestStatus(activeRequest.status)}`"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="处理结果">
          <el-select v-model="processForm.status" style="width: 100%">
            <el-option v-for="item in processOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理备注">
          <el-input
            v-model="processForm.processNote"
            type="textarea"
            :rows="4"
            placeholder="请输入处理结果、执行说明或需要继续跟进的备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeProcessDialog">取消</el-button>
        <el-button type="primary" :loading="processing" @click="submitProcess">保存处理结果</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="batchDialogVisible" title="批量处理订单申请" width="560px" @closed="closeBatchDialog">
      <el-form label-position="top">
        <el-alert
          :title="`即将把 ${batchTargetCount} 条申请更新为${batchTargetLabel}`"
          description="批量动作会逐条调用订单申请处理链路，并写入对应订单审计记录。"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="目标状态">
          <el-select v-model="batchForm.status" style="width: 100%">
            <el-option label="同意申请" value="APPROVED" />
            <el-option label="驳回申请" value="REJECTED" />
            <el-option label="处理完成" value="COMPLETED" />
            <el-option label="撤销申请" value="CANCELLED" />
          </el-select>
        </el-form-item>
        <el-form-item label="处理备注">
          <el-input
            v-model="batchForm.processNote"
            type="textarea"
            :rows="4"
            placeholder="请输入批量处理说明，便于后续审计和追踪"
          />
        </el-form-item>
        <el-form-item label="本次涉及申请">
          <div class="table-toolbar__meta" style="line-height: 1.8">
            <span v-for="item in batchRows" :key="item.id">{{ item.requestNo }} / {{ item.customerName }}</span>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeBatchDialog">取消</el-button>
        <el-button type="primary" :loading="batchProcessing" @click="submitBatchProcess">确认批量处理</el-button>
      </template>
    </el-dialog>
  </div>
</template>
