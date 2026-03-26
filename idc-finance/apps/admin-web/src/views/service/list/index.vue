<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import {
  fetchProviderAccounts,
  fetchServices,
  runServiceAction,
  syncMofangService,
  type ProviderAccount,
  type ServiceQuery,
  type ServiceRecord
} from "@/api/admin";

type TabKey = "ALL" | "ACTIVE" | "SUSPENDED" | "TERMINATED" | "MOFANG_CLOUD";
type ServiceBatchAction = "activate" | "suspend" | "reboot" | "terminate";

const router = useRouter();
const route = useRoute();

const loading = ref(false);
const syncLoading = ref(false);
const batchActionLoading = ref(false);
const services = ref<ServiceRecord[]>([]);
const total = ref(0);
const providerAccounts = ref<ProviderAccount[]>([]);
const selectedRows = ref<ServiceRecord[]>([]);
const activeTab = ref<TabKey>("ALL");
const batchActionVisible = ref(false);
const batchActionRows = ref<ServiceRecord[]>([]);

const filters = reactive({
  keyword: "",
  customerId: "",
  orderId: "",
  providerType: "",
  providerAccountId: "",
  syncStatus: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const batchActionForm = reactive({
  action: "activate" as ServiceBatchAction
});

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: total.value },
  { key: "ACTIVE", label: "运行中", count: services.value.filter(item => item.status === "ACTIVE").length },
  { key: "SUSPENDED", label: "已暂停", count: services.value.filter(item => item.status === "SUSPENDED").length },
  { key: "TERMINATED", label: "已终止", count: services.value.filter(item => item.status === "TERMINATED").length },
  { key: "MOFANG_CLOUD", label: "魔方云", count: services.value.filter(item => item.providerType === "MOFANG_CLOUD").length }
]);

const activeCount = computed(() => services.value.filter(item => item.status === "ACTIVE").length);
const suspendedCount = computed(() => services.value.filter(item => item.status === "SUSPENDED").length);
const terminatedCount = computed(() => services.value.filter(item => item.status === "TERMINATED").length);
const syncFailedCount = computed(() => services.value.filter(item => item.syncStatus === "FAILED").length);
const currentContextSummary = computed(() => {
  const items: string[] = [];
  if (filters.customerId.trim()) items.push(`客户 ${filters.customerId.trim()}`);
  if (filters.orderId.trim()) items.push(`订单 ${filters.orderId.trim()}`);
  if (filters.providerType) items.push(providerLabel(filters.providerType));
  if (filters.providerAccountId) items.push(providerAccountLabel(Number(filters.providerAccountId)));
  return items.join(" / ");
});
const selectedMofangRows = computed(() => selectedRows.value.filter(item => item.providerType === "MOFANG_CLOUD"));
const currentPageMofangRows = computed(() => services.value.filter(item => item.providerType === "MOFANG_CLOUD"));
const syncActionLabel = computed(() => {
  if (selectedRows.value.length > 0) {
    return selectedMofangRows.value.length > 0
      ? `同步选中实例服务 (${selectedMofangRows.value.length})`
      : "同步选中实例服务";
  }
  return currentPageMofangRows.value.length > 0
    ? `同步当前页实例服务 (${currentPageMofangRows.value.length})`
    : "同步当前页实例服务";
});
const syncActionDisabled = computed(() => {
  if (syncLoading.value) return true;
  if (selectedRows.value.length > 0) return selectedMofangRows.value.length === 0;
  return currentPageMofangRows.value.length === 0;
});
const batchActionLabel = computed(() => serviceActionLabel(batchActionForm.action));
const batchActionTargetCount = computed(() => batchActionRows.value.length);

function buildQuery(): ServiceQuery {
  const query: ServiceQuery = {
    page: pagination.page,
    limit: pagination.limit,
    sort: "created_at",
    order: "desc"
  };
  if (activeTab.value !== "ALL" && activeTab.value !== "MOFANG_CLOUD") query.status = activeTab.value;
  if (activeTab.value === "MOFANG_CLOUD") query.provider_type = "MOFANG_CLOUD";
  if (filters.keyword.trim()) query.keyword = filters.keyword.trim();
  if (filters.customerId.trim()) query.customer_id = filters.customerId.trim();
  if (filters.orderId.trim()) query.order_id = filters.orderId.trim();
  if (filters.providerType) query.provider_type = filters.providerType;
  if (filters.providerAccountId) query.provider_account_id = filters.providerAccountId;
  if (filters.syncStatus) query.sync_status = filters.syncStatus;
  return query;
}

async function loadServices() {
  loading.value = true;
  try {
    const [serviceData, accountData] = await Promise.all([fetchServices(buildQuery()), fetchProviderAccounts()]);
    services.value = serviceData.items;
    total.value = serviceData.total;
    providerAccounts.value = accountData;
  } finally {
    loading.value = false;
  }
}

function readRouteQueryValue(value: unknown) {
  if (Array.isArray(value)) return String(value[0] ?? "");
  if (value === undefined || value === null) return "";
  return String(value);
}

function syncFiltersFromRoute() {
  filters.keyword = readRouteQueryValue(route.query.keyword);
  filters.customerId = readRouteQueryValue(route.query.customerId || route.query.uid);
  filters.orderId = readRouteQueryValue(route.query.orderId);
  filters.providerType = readRouteQueryValue(route.query.providerType).toUpperCase();
  filters.providerAccountId = readRouteQueryValue(route.query.providerAccountId || route.query.accountId);
  filters.syncStatus = readRouteQueryValue(route.query.syncStatus).toUpperCase();
  pagination.page = 1;
}

function applyFilters() {
  pagination.page = 1;
  void loadServices();
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.customerId = "";
  filters.orderId = "";
  filters.providerType = "";
  filters.providerAccountId = "";
  filters.syncStatus = "";
  pagination.page = 1;
  void router.push({ path: "/services/list", query: {} });
}

function handleSelectionChange(rows: ServiceRecord[]) {
  selectedRows.value = rows;
}

function providerAccountLabel(id?: number) {
  if (!id) return "未绑定接口账户";
  const matched = providerAccounts.value.find(item => item.id === id);
  if (!matched) return `接口账户 #${id}`;
  return matched.name || matched.sourceName || matched.baseUrl;
}

function providerLabel(value: string) {
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "财务上下游",
    LOCAL: "本地模块",
    RESOURCE: "资源池",
    MANUAL: "手动资源",
    MOCK: "演示资源"
  };
  return mapping[value] ?? value;
}

function serviceStatusLabel(value: string) {
  const mapping: Record<string, string> = {
    PENDING: "待开通",
    ACTIVE: "运行中",
    SUSPENDED: "已暂停",
    PROVISIONING: "处理中",
    TERMINATED: "已终止",
    FAILED: "异常"
  };
  return mapping[value] ?? value;
}

function serviceStatusType(value: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    ACTIVE: "success",
    SUSPENDED: "warning",
    PROVISIONING: "primary",
    TERMINATED: "info",
    FAILED: "danger"
  };
  return mapping[value] ?? "info";
}

function syncStatusLabel(value: string) {
  const mapping: Record<string, string> = {
    SUCCESS: "同步成功",
    FAILED: "同步失败",
    IDLE: "未同步",
    PENDING: "待同步"
  };
  return mapping[value] ?? value;
}

function syncStatusType(value: string) {
  const mapping: Record<string, string> = {
    SUCCESS: "success",
    FAILED: "danger",
    IDLE: "info",
    PENDING: "warning"
  };
  return mapping[value] ?? "info";
}

function serviceActionLabel(action: ServiceBatchAction) {
  const mapping: Record<ServiceBatchAction, string> = {
    activate: "开通服务",
    suspend: "暂停服务",
    reboot: "重启服务",
    terminate: "终止服务"
  };
  return mapping[action];
}

function serviceActionDescription(action: ServiceBatchAction) {
  const mapping: Record<ServiceBatchAction, string> = {
    activate: "把服务重新切回可用状态，适合人工放行、恢复暂停服务或修正状态时使用。",
    suspend: "立即暂停服务，保留业务记录和资源上下文，适合欠费或人工封停场景。",
    reboot: "触发服务重启动作，适合实例类服务的故障恢复。",
    terminate: "将服务标记为终止并停止后续交付，适合关闭业务或清理异常服务。"
  };
  return mapping[action];
}

function canRunAction(row: ServiceRecord, action: ServiceBatchAction) {
  switch (action) {
    case "activate":
      return row.status !== "ACTIVE" && row.status !== "TERMINATED";
    case "suspend":
      return row.status === "ACTIVE" || row.status === "PROVISIONING";
    case "reboot":
      return row.status === "ACTIVE";
    case "terminate":
      return row.status !== "TERMINATED";
    default:
      return true;
  }
}

function openDetail(row: ServiceRecord) {
  void router.push(`/services/detail/${row.id}`);
}

function openCustomerWorkbench(row: ServiceRecord) {
  void router.push(`/customer/detail/${row.customerId}`);
}

function openOrderWorkbench(row: ServiceRecord) {
  if (!row.orderId) return;
  void router.push(`/orders/detail/${row.orderId}`);
}

function openInvoiceWorkbench(row: ServiceRecord) {
  if (!row.invoiceId) return;
  void router.push(`/billing/invoices/${row.invoiceId}`);
}

function openTicketWorkbench(row: ServiceRecord) {
  void router.push({
    path: "/tickets/list",
    query: {
      serviceId: String(row.id),
      customerId: String(row.customerId)
    }
  });
}

function openAutomationWorkbench(row: ServiceRecord) {
  void router.push({
    path: "/providers/automation",
    query: {
      serviceId: String(row.id),
      orderId: row.orderId ? String(row.orderId) : undefined,
      invoiceId: row.invoiceId ? String(row.invoiceId) : undefined,
      channel: row.providerType || undefined
    }
  });
}

function openResourcesWorkbench(row: ServiceRecord) {
  void router.push({
    path: "/providers/resources",
    query: {
      providerType: row.providerType || undefined,
      accountId: row.providerAccountId ? String(row.providerAccountId) : undefined,
      serviceId: String(row.id)
    }
  });
}

function openFinanceWorkbench(row: ServiceRecord) {
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId: String(row.customerId)
    }
  });
}

async function syncSelectedServices(rows?: ServiceRecord[]) {
  const baseRows = rows ?? (selectedRows.value.length > 0 ? selectedRows.value : currentPageMofangRows.value);
  const targetRows = baseRows.filter(item => item.providerType === "MOFANG_CLOUD");
  if (targetRows.length === 0) {
    ElMessage.info(selectedRows.value.length > 0 ? "选中的服务里没有可同步的实例渠道服务" : "当前页没有可同步的实例渠道服务");
    return;
  }

  syncLoading.value = true;
  try {
    let successCount = 0;
    for (const item of targetRows) {
      await syncMofangService(item.id);
      successCount += 1;
    }
    const ignoredCount = baseRows.length - targetRows.length;
    ElMessage.success(
      ignoredCount > 0
        ? `已完成 ${successCount} 台实例服务同步，忽略 ${ignoredCount} 台非实例渠道服务`
        : `已完成 ${successCount} 台实例服务同步`
    );
    await loadServices();
  } finally {
    syncLoading.value = false;
  }
}

function openBatchAction(action: ServiceBatchAction, rows?: ServiceRecord[]) {
  const baseRows = rows && rows.length > 0 ? rows : selectedRows.value;
  if (baseRows.length === 0) {
    ElMessage.info("请先选择服务");
    return;
  }
  const actionableRows = baseRows.filter(item => canRunAction(item, action));
  if (actionableRows.length === 0) {
    ElMessage.info(`当前选中的服务都不适合执行“${serviceActionLabel(action)}”`);
    return;
  }
  batchActionForm.action = action;
  batchActionRows.value = actionableRows;
  batchActionVisible.value = true;
}

async function submitBatchAction() {
  if (batchActionRows.value.length === 0) return;
  batchActionLoading.value = true;
  let success = 0;
  const failed: ServiceRecord[] = [];
  try {
    for (const row of batchActionRows.value) {
      try {
        await runServiceAction(row.id, batchActionForm.action);
        success += 1;
      } catch {
        failed.push(row);
      }
    }
    if (success > 0) {
      ElMessage.success(
        failed.length > 0
          ? `${serviceActionLabel(batchActionForm.action)}完成 ${success} 条，失败 ${failed.length} 条`
          : `已完成 ${success} 条服务的${serviceActionLabel(batchActionForm.action)}`
      );
    } else {
      ElMessage.error(`${serviceActionLabel(batchActionForm.action)}失败`);
    }
    batchActionRows.value = failed;
    if (failed.length === 0) {
      batchActionVisible.value = false;
      selectedRows.value = [];
    }
    await loadServices();
  } finally {
    batchActionLoading.value = false;
  }
}

function closeBatchAction() {
  batchActionVisible.value = false;
  batchActionRows.value = [];
}

function handleRowCommand(row: ServiceRecord, command: string | number | object) {
  switch (String(command)) {
    case "customer":
      openCustomerWorkbench(row);
      return;
    case "order":
      openOrderWorkbench(row);
      return;
    case "invoice":
      openInvoiceWorkbench(row);
      return;
    case "tickets":
      openTicketWorkbench(row);
      return;
    case "automation":
      openAutomationWorkbench(row);
      return;
    case "resources":
      openResourcesWorkbench(row);
      return;
    case "finance":
      openFinanceWorkbench(row);
      return;
    case "sync":
      void syncSelectedServices([row]);
      return;
    case "activate":
    case "suspend":
    case "reboot":
    case "terminate":
      openBatchAction(String(command) as ServiceBatchAction, [row]);
      return;
    default:
      return;
  }
}

async function copySelectedServiceNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择服务");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.serviceNo).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个服务编号`);
}

function exportRows(rows: ServiceRecord[], filename: string) {
  downloadCsv(
    filename,
    ["服务编号", "订单 ID", "客户 ID", "远端资源 ID", "产品名称", "渠道", "接口账户", "服务状态", "公网 IP", "同步状态", "最近同步", "下次到期"],
    rows.map(item => [
      item.serviceNo,
      String(item.orderId || ""),
      String(item.customerId || ""),
      item.providerResourceId,
      item.productName,
      providerLabel(item.providerType),
      providerAccountLabel(item.providerAccountId),
      serviceStatusLabel(item.status),
      item.ipAddress || "-",
      syncStatusLabel(item.syncStatus),
      item.lastSyncAt || "-",
      item.nextDueAt || "-"
    ])
  );
}

function exportCurrent() {
  exportRows(services.value, `services-page-${pagination.page}.csv`);
  ElMessage.success("当前页服务已导出");
}

function exportSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择服务");
    return;
  }
  exportRows(selectedRows.value, "services-selected.csv");
  ElMessage.success("已导出选中服务");
}

watch(activeTab, () => {
  pagination.page = 1;
  void loadServices();
});

onMounted(() => {
  syncFiltersFromRoute();
  void loadServices();
});

watch(
  () => route.fullPath,
  () => {
    syncFiltersFromRoute();
    void loadServices();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 服务"
      title="服务列表"
      subtitle="集中查看服务状态、订单归属、工单入口、渠道资源、自动化任务和批量生命周期动作。"
    >
      <template #actions>
        <el-button @click="loadServices">刷新列表</el-button>
        <el-button type="primary" :loading="syncLoading" :disabled="syncActionDisabled" @click="syncSelectedServices()">
          {{ syncActionLabel }}
        </el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>服务总数</span><strong>{{ total }}</strong></div>
          <div class="summary-pill"><span>运行中</span><strong>{{ activeCount }}</strong></div>
          <div class="summary-pill"><span>已暂停</span><strong>{{ suspendedCount }}</strong></div>
          <div class="summary-pill"><span>已终止</span><strong>{{ terminatedCount }}</strong></div>
          <div class="summary-pill"><span>同步失败</span><strong>{{ syncFailedCount }}</strong></div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div v-if="currentContextSummary" class="panel-card" style="margin-bottom: 16px;">
          <div class="section-card__head">
            <strong>当前服务上下文</strong>
            <span class="section-card__meta">{{ currentContextSummary }}</span>
          </div>
          <div class="action-group">
            <el-button v-if="filters.customerId" plain @click="router.push(`/customer/detail/${filters.customerId}`)">查看客户</el-button>
            <el-button v-if="filters.orderId" plain @click="router.push(`/orders/detail/${filters.orderId}`)">查看订单</el-button>
            <el-button plain @click="resetFilters">退出上下文</el-button>
          </div>
        </div>

        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="搜索服务号、产品名、远端资源 ID 或公网 IP" clearable />
          <el-input v-model="filters.customerId" placeholder="客户 ID" clearable />
          <el-input v-model="filters.orderId" placeholder="订单 ID" clearable />
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div class="filter-bar filter-bar--compact">
          <el-select v-model="filters.providerType" placeholder="自动化渠道" clearable>
            <el-option label="魔方云" value="MOFANG_CLOUD" />
            <el-option label="财务上下游" value="ZJMF_API" />
            <el-option label="本地模块" value="LOCAL" />
            <el-option label="资源池" value="RESOURCE" />
            <el-option label="手动资源" value="MANUAL" />
          </el-select>
          <el-select v-model="filters.providerAccountId" placeholder="接口账户" clearable>
            <el-option
              v-for="item in providerAccounts"
              :key="item.id"
              :label="providerAccountLabel(item.id)"
              :value="String(item.id)"
            />
          </el-select>
          <el-select v-model="filters.syncStatus" placeholder="同步状态" clearable>
            <el-option label="同步成功" value="SUCCESS" />
            <el-option label="同步失败" value="FAILED" />
            <el-option label="待同步" value="PENDING" />
            <el-option label="未同步" value="IDLE" />
          </el-select>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>服务工作台</strong>
          <span>当前页 {{ services.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchAction('activate')">开通服务</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchAction('suspend')">暂停服务</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchAction('reboot')">重启服务</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openBatchAction('terminate')">终止服务</el-button>
          <el-button plain @click="copySelectedServiceNos">复制服务号</el-button>
          <el-button plain @click="exportSelected">导出选中</el-button>
          <el-button plain @click="exportCurrent">导出当前页</el-button>
        </div>
      </div>

      <el-table
        :data="services"
        border
        stripe
        empty-text="暂无符合条件的服务"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column prop="serviceNo" label="服务编号" min-width="170" />
        <el-table-column label="客户 / 订单" min-width="180">
          <template #default="{ row }">
            <div class="inline-actions" style="flex-direction: column; align-items: flex-start; gap: 4px;">
              <el-button type="primary" link @click="openCustomerWorkbench(row)">客户 #{{ row.customerId }}</el-button>
              <el-button v-if="row.orderId" type="primary" link @click="openOrderWorkbench(row)">订单 #{{ row.orderId }}</el-button>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="productName" label="产品名称" min-width="240" show-overflow-tooltip />
        <el-table-column label="自动化渠道" min-width="140">
          <template #default="{ row }">{{ providerLabel(row.providerType) }}</template>
        </el-table-column>
        <el-table-column label="接口账户" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">{{ providerAccountLabel(row.providerAccountId) }}</template>
        </el-table-column>
        <el-table-column label="服务状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="serviceStatusType(row.status)" effect="light">{{ serviceStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="地域 / IP" min-width="190">
          <template #default="{ row }">
            <div style="display: grid; gap: 4px;">
              <span>{{ row.regionName || "-" }}</span>
              <small style="color: var(--el-text-color-secondary)">{{ row.ipAddress || "无公网 IP" }}</small>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="同步状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="syncStatusType(row.syncStatus)" effect="light">{{ syncStatusLabel(row.syncStatus) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastSyncAt" label="最近同步" min-width="180" />
        <el-table-column prop="nextDueAt" label="下次到期" min-width="160" />
        <el-table-column label="操作" min-width="320" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetail(row)">工作台</el-button>
              <el-button type="primary" link @click="openTicketWorkbench(row)">工单</el-button>
              <el-button type="primary" link @click="openAutomationWorkbench(row)">任务</el-button>
              <el-dropdown @command="handleRowCommand(row, $event)">
                <el-button type="primary" link>更多</el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="customer">客户详情</el-dropdown-item>
                    <el-dropdown-item command="order" :disabled="!row.orderId">订单工作台</el-dropdown-item>
                    <el-dropdown-item command="invoice" :disabled="!row.invoiceId">账单工作台</el-dropdown-item>
                    <el-dropdown-item command="finance">客户资金台账</el-dropdown-item>
                    <el-dropdown-item command="resources">渠道资源</el-dropdown-item>
                    <el-dropdown-item command="automation">自动化任务</el-dropdown-item>
                    <el-dropdown-item command="sync" :disabled="row.providerType !== 'MOFANG_CLOUD'">立即同步</el-dropdown-item>
                    <el-dropdown-item command="activate" :disabled="!canRunAction(row, 'activate')">开通服务</el-dropdown-item>
                    <el-dropdown-item command="suspend" :disabled="!canRunAction(row, 'suspend')">暂停服务</el-dropdown-item>
                    <el-dropdown-item command="reboot" :disabled="!canRunAction(row, 'reboot')">重启服务</el-dropdown-item>
                    <el-dropdown-item command="terminate" :disabled="!canRunAction(row, 'terminate')">终止服务</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
          @current-change="(page:number) => { pagination.page = page; void loadServices(); }"
          @size-change="(limit:number) => { pagination.limit = limit; pagination.page = 1; void loadServices(); }"
        />
      </div>
    </PageWorkbench>

    <el-dialog v-model="batchActionVisible" :title="batchActionLabel" width="560px" @closed="closeBatchAction">
      <el-alert
        :title="`将对 ${batchActionTargetCount} 条服务执行“${batchActionLabel}”`"
        :description="serviceActionDescription(batchActionForm.action)"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 16px"
      />

      <div class="panel-card" style="display: grid; gap: 10px;">
        <div class="section-card__head">
          <strong>本次服务清单</strong>
          <span class="section-card__meta">执行失败的服务会保留在这里，方便你继续处理。</span>
        </div>
        <div v-for="item in batchActionRows" :key="item.id" class="summary-pill">
          <span>{{ item.serviceNo }} / {{ item.productName }}</span>
          <strong>{{ serviceStatusLabel(item.status) }}</strong>
        </div>
      </div>

      <template #footer>
        <el-button @click="closeBatchAction">取消</el-button>
        <el-button type="primary" :loading="batchActionLoading" @click="submitBatchAction">确认执行</el-button>
      </template>
    </el-dialog>
  </div>
</template>
