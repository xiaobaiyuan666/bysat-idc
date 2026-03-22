<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import {
  fetchProviderAccounts,
  fetchServices,
  syncMofangService,
  type ServiceQuery,
  type ProviderAccount,
  type ServiceRecord
} from "@/api/admin";

type TabKey = "ALL" | "ACTIVE" | "SUSPENDED" | "TERMINATED" | "MOFANG_CLOUD";

const router = useRouter();

const loading = ref(false);
const syncLoading = ref(false);
const services = ref<ServiceRecord[]>([]);
const total = ref(0);
const providerAccounts = ref<ProviderAccount[]>([]);
const selectedRows = ref<ServiceRecord[]>([]);
const activeTab = ref<TabKey>("ALL");

const filters = reactive({
  keyword: "",
  providerType: "",
  providerAccountId: "",
  syncStatus: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
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
  if (filters.providerType) query.provider_type = filters.providerType;
  if (filters.providerAccountId) query.provider_account_id = filters.providerAccountId;
  if (filters.syncStatus) query.sync_status = filters.syncStatus;
  return query;
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.providerType = "";
  filters.providerAccountId = "";
  filters.syncStatus = "";
  pagination.page = 1;
  void loadServices();
}

function applyFilters() {
  pagination.page = 1;
  void loadServices();
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
    TERMINATED: "danger",
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

async function syncSelectedServices(rows?: ServiceRecord[]) {
  const targetRows = (rows ?? selectedRows.value).filter(item => item.providerType === "MOFANG_CLOUD");
  if (targetRows.length === 0) {
    ElMessage.info("请先选择魔方云服务");
    return;
  }

  syncLoading.value = true;
  try {
    let successCount = 0;
    for (const item of targetRows) {
      await syncMofangService(item.id);
      successCount += 1;
    }
    ElMessage.success(`已完成 ${successCount} 台服务的同步`);
    await loadServices();
  } finally {
    syncLoading.value = false;
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
    ["服务编号", "远端资源 ID", "产品名称", "渠道", "接口账户", "服务状态", "公网 IP", "同步状态", "最近同步", "下次到期"],
    rows.map(item => [
      item.serviceNo,
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

function openDetail(row: ServiceRecord) {
  void router.push(`/services/detail/${row.id}`);
}

watch(activeTab, () => {
  pagination.page = 1;
  void loadServices();
});

onMounted(() => {
  void loadServices();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 服务"
      title="服务列表"
      subtitle="集中查看服务状态、自动化渠道、接口账户、同步状态和资源入口。"
    >
      <template #actions>
        <el-button @click="loadServices">刷新列表</el-button>
        <el-button type="primary" :loading="syncLoading" @click="syncSelectedServices()">同步选中魔方云服务</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>服务总数</span><strong>{{ total }}</strong></div>
          <div class="summary-pill">
            <span>运行中</span>
            <strong>{{ activeCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>已暂停</span>
            <strong>{{ suspendedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>已终止</span>
            <strong>{{ terminatedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>同步失败</span>
            <strong>{{ syncFailedCount }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="搜索服务号、产品名、远端资源 ID 或公网 IP" clearable />
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
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div class="filter-bar filter-bar--compact">
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
          <strong>业务列表</strong>
          <span>当前页 {{ services.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
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
        <el-table-column prop="providerResourceId" label="远端资源 ID" min-width="140" />
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
        <el-table-column prop="regionName" label="地域" min-width="140" />
        <el-table-column prop="ipAddress" label="公网 IP" min-width="160" />
        <el-table-column label="同步状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="syncStatusType(row.syncStatus)" effect="light">{{ syncStatusLabel(row.syncStatus) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="lastSyncAt" label="最近同步" min-width="180" />
        <el-table-column prop="nextDueAt" label="下次到期" min-width="140" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetail(row)">进入工作台</el-button>
              <el-button
                v-if="row.providerType === 'MOFANG_CLOUD'"
                type="primary"
                link
                :disabled="syncLoading"
                @click="syncSelectedServices([row])"
              >
                立即同步
              </el-button>
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
  </div>
</template>
