<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchProviderAccounts,
  fetchUpstreamImportHistory,
  fetchUpstreamImportHistoryDetail,
  type ProviderAccount,
  type UpstreamImportHistoryDetail,
  type UpstreamImportHistoryRecord
} from "@/api/admin";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const detailLoading = ref(false);
const accounts = ref<ProviderAccount[]>([]);
const records = ref<UpstreamImportHistoryRecord[]>([]);
const detailVisible = ref(false);
const detail = ref<UpstreamImportHistoryDetail | null>(null);

const filters = reactive({
  providerAccountId: "",
  providerType: "",
  status: "",
  keyword: ""
});

const providerTypeOptions = [
  { label: "上游财务", value: "ZJMF_API" },
  { label: "WHMCS", value: "WHMCS" },
  { label: "魔方云", value: "MOFANG_CLOUD" },
  { label: "资源池", value: "RESOURCE" },
  { label: "手工交付", value: "MANUAL" }
];

const statusOptions = [
  { label: "成功", value: "SUCCESS" },
  { label: "部分成功", value: "PARTIAL" },
  { label: "失败", value: "FAILED" }
];

const summary = computed(() => ({
  total: records.value.length,
  success: records.value.filter(item => item.status === "SUCCESS").length,
  partial: records.value.filter(item => item.status === "PARTIAL").length,
  failed: records.value.filter(item => item.status === "FAILED").length,
  changed: records.value.reduce((sum, item) => sum + item.created + item.updated + item.deactivated, 0)
}));

function syncFiltersFromRoute() {
  filters.providerAccountId = typeof route.query.providerAccountId === "string" ? route.query.providerAccountId : "";
  filters.providerType = typeof route.query.providerType === "string" ? route.query.providerType : "";
  filters.status = typeof route.query.status === "string" ? route.query.status : "";
  filters.keyword = typeof route.query.keyword === "string" ? route.query.keyword : "";
}

async function loadData() {
  loading.value = true;
  try {
    syncFiltersFromRoute();
    const [accountList, history] = await Promise.all([
      fetchProviderAccounts(),
      fetchUpstreamImportHistory({
        providerAccountId: filters.providerAccountId || undefined,
        providerType: filters.providerType || undefined,
        status: filters.status || undefined,
        keyword: filters.keyword || undefined,
        limit: 80
      })
    ]);
    accounts.value = accountList;
    records.value = history.items;
    const routeHistoryId = typeof route.query.historyId === "string" ? route.query.historyId : "";
    if (routeHistoryId) {
      const matched = history.items.find(item => item.historyId === routeHistoryId);
      if (matched) {
        await openDetail(matched);
      }
    }
  } finally {
    loading.value = false;
  }
}

function applyFilters() {
  void router.push({
    path: "/providers/upstream-sync-history",
    query: {
      providerAccountId: filters.providerAccountId || undefined,
      providerType: filters.providerType || undefined,
      status: filters.status || undefined,
      keyword: filters.keyword.trim() || undefined
    }
  });
}

function resetFilters() {
  filters.providerAccountId = "";
  filters.providerType = "";
  filters.status = "";
  filters.keyword = "";
  void router.push({ path: "/providers/upstream-sync-history", query: {} });
}

function providerTypeLabel(value: string) {
  return providerTypeOptions.find(item => item.value === value)?.label ?? value;
}

function statusType(value: string) {
  return (
    {
      SUCCESS: "success",
      PARTIAL: "warning",
      FAILED: "danger"
    }[value] ?? "info"
  );
}

function statusLabel(value: string) {
  return (
    {
      SUCCESS: "成功",
      PARTIAL: "部分成功",
      FAILED: "失败"
    }[value] ?? value
  );
}

function formatRequestedCodes(value: string[]) {
  if (!value?.length) return "全量同步";
  if (value.length <= 3) return value.join("、");
  return `${value.slice(0, 3).join("、")} 等 ${value.length} 个`;
}

async function openDetail(record: UpstreamImportHistoryRecord) {
  detailVisible.value = true;
  detailLoading.value = true;
  detail.value = null;
  try {
    detail.value = await fetchUpstreamImportHistoryDetail(record.historyId);
  } finally {
    detailLoading.value = false;
  }
}

function openAccountWorkbench(record: UpstreamImportHistoryRecord) {
  void router.push({
    path: "/providers/accounts",
    query: {
      accountId: String(record.providerAccountId),
      providerType: record.providerType
    }
  });
}

function openProductWorkbench(record: UpstreamImportHistoryRecord) {
  void router.push({
    path: "/catalog/products",
    query: {
      providerAccountId: String(record.providerAccountId),
      providerType: record.providerType
    }
  });
}

onMounted(() => {
  void loadData();
});

watch(
  () => route.fullPath,
  () => {
    void loadData();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="接口与上游 / 同步记录"
      title="上游同步记录工作台"
      subtitle="集中查看每次上游商品同步的结果、失败项、停用项以及对应接口账户。"
    >
      <template #actions>
        <el-select v-model="filters.providerAccountId" placeholder="接口账户" clearable style="width: 220px">
          <el-option
            v-for="item in accounts.filter(account => ['ZJMF_API', 'WHMCS'].includes(account.providerType))"
            :key="item.id"
            :label="item.name"
            :value="String(item.id)"
          />
        </el-select>
        <el-select v-model="filters.providerType" placeholder="渠道类型" clearable style="width: 160px">
          <el-option v-for="item in providerTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-select v-model="filters.status" placeholder="执行结果" clearable style="width: 150px">
          <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-input v-model="filters.keyword" placeholder="搜索账户、来源或商品编码" clearable style="width: 240px" />
        <el-button @click="applyFilters">查询</el-button>
        <el-button @click="resetFilters">重置</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>记录数</span><strong>{{ summary.total }}</strong></div>
          <div class="summary-pill"><span>成功</span><strong>{{ summary.success }}</strong></div>
          <div class="summary-pill"><span>部分成功</span><strong>{{ summary.partial }}</strong></div>
          <div class="summary-pill"><span>失败</span><strong>{{ summary.failed }}</strong></div>
          <div class="summary-pill"><span>变更商品</span><strong>{{ summary.changed }}</strong></div>
        </div>
      </template>

      <el-table :data="records" border stripe empty-text="当前没有上游同步记录">
        <el-table-column prop="createdAt" label="开始时间" min-width="170" />
        <el-table-column label="接口账户" min-width="220">
          <template #default="{ row }">
            <div class="cell-stack">
              <strong>{{ row.accountName || `账户 #${row.providerAccountId}` }}</strong>
              <span>{{ row.sourceName || "-" }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="渠道" min-width="120">
          <template #default="{ row }">{{ providerTypeLabel(row.providerType) }}</template>
        </el-table-column>
        <el-table-column label="执行结果" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="同步范围" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">
            {{ formatRequestedCodes(row.requestedCodes) }}
          </template>
        </el-table-column>
        <el-table-column label="结果汇总" min-width="220">
          <template #default="{ row }">
            <div class="cell-stack">
              <span>新增 {{ row.created }} / 更新 {{ row.updated }} / 停用 {{ row.deactivated }}</span>
              <span>失败 {{ row.failed }} / 总项 {{ row.total }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="durationMs" label="耗时(ms)" min-width="110" />
        <el-table-column prop="message" label="说明" min-width="220" show-overflow-tooltip />
        <el-table-column label="操作" min-width="220" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetail(row)">详情</el-button>
              <el-button type="primary" link @click="openAccountWorkbench(row)">接口账户</el-button>
              <el-button type="primary" link @click="openProductWorkbench(row)">商品</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </PageWorkbench>

    <el-drawer v-model="detailVisible" size="900px" :with-header="false">
      <div v-loading="detailLoading" class="drawer-body">
        <template v-if="detail">
          <div class="section-card__head">
            <strong>同步记录详情</strong>
            <el-tag :type="statusType(detail.record.status)" effect="light">{{ statusLabel(detail.record.status) }}</el-tag>
          </div>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="记录号">{{ detail.record.historyId }}</el-descriptions-item>
            <el-descriptions-item label="接口账户">{{ detail.record.accountName || detail.record.providerAccountId }}</el-descriptions-item>
            <el-descriptions-item label="渠道">{{ providerTypeLabel(detail.record.providerType) }}</el-descriptions-item>
            <el-descriptions-item label="来源">{{ detail.record.sourceName || "-" }}</el-descriptions-item>
            <el-descriptions-item label="开始时间">{{ detail.record.createdAt }}</el-descriptions-item>
            <el-descriptions-item label="结束时间">{{ detail.record.finishedAt }}</el-descriptions-item>
            <el-descriptions-item label="耗时">{{ detail.record.durationMs }} ms</el-descriptions-item>
            <el-descriptions-item label="同步范围">
              {{ formatRequestedCodes(detail.record.requestedCodes) }}
            </el-descriptions-item>
            <el-descriptions-item label="自动同步价格">{{ detail.record.autoSyncPricing ? "是" : "否" }}</el-descriptions-item>
            <el-descriptions-item label="自动同步配置">{{ detail.record.autoSyncConfig ? "是" : "否" }}</el-descriptions-item>
            <el-descriptions-item label="自动同步模板">{{ detail.record.autoSyncTemplate ? "是" : "否" }}</el-descriptions-item>
            <el-descriptions-item label="缺失自动停用">{{ detail.record.deactivateMissing ? "是" : "否" }}</el-descriptions-item>
            <el-descriptions-item label="说明" :span="2">{{ detail.record.message }}</el-descriptions-item>
          </el-descriptions>

          <div class="section-card__head" style="margin-top: 20px">
            <strong>逐项结果</strong>
            <span class="section-card__meta">可直接定位失败项和自动停用项</span>
          </div>

          <el-table :data="detail.items" border stripe empty-text="当前没有逐项结果">
            <el-table-column prop="remoteProductCode" label="远端编码" min-width="140" />
            <el-table-column prop="remoteProductName" label="远端商品" min-width="220" show-overflow-tooltip />
            <el-table-column prop="groupName" label="分组" min-width="160" show-overflow-tooltip />
            <el-table-column prop="productNo" label="本地商品" min-width="150" />
            <el-table-column prop="operation" label="动作" min-width="100" />
            <el-table-column label="状态" min-width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'SUCCESS' ? 'success' : 'danger'" effect="light">
                  {{ row.status === "SUCCESS" ? "成功" : "失败" }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="message" label="结果说明" min-width="260" show-overflow-tooltip />
          </el-table>
        </template>
      </div>
    </el-drawer>
  </div>
</template>

<style scoped>
.inline-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.cell-stack {
  display: grid;
  gap: 4px;
}

.cell-stack span {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.drawer-body {
  display: grid;
  gap: 16px;
}
</style>
