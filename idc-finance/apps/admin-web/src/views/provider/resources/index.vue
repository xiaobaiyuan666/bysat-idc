<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchMofangHealth,
  fetchMofangInstances,
  fetchMofangSyncLogs,
  fetchProviderAccounts,
  pullMofangSync,
  type MofangHealthResponse,
  type MofangInstanceSummary,
  type MofangSyncLogItem,
  type ProviderAccount
} from "@/api/admin";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const syncLoading = ref(false);
const selectedRemoteIds = ref<string[]>([]);
const selectedAccountId = ref<number>(0);
const selectedProviderType = ref<string>("");
const focusServiceId = ref<number>(0);
const health = ref<MofangHealthResponse | null>(null);
const instances = ref<MofangInstanceSummary[]>([]);
const syncLogs = ref<MofangSyncLogItem[]>([]);
const accounts = ref<ProviderAccount[]>([]);

const activeAccounts = computed(() => accounts.value.filter(item => item.status === "ACTIVE"));
const selectedAccount = computed(() => activeAccounts.value.find(item => item.id === selectedAccountId.value) ?? null);
const supportsRemoteResources = computed(() => selectedAccount.value?.providerType === "MOFANG_CLOUD");

const filteredAccounts = computed(() => {
  if (!selectedProviderType.value) return activeAccounts.value;
  return activeAccounts.value.filter(item => item.providerType === selectedProviderType.value);
});

const providerSummary = computed(() => ({
  total: activeAccounts.value.length,
  mofang: activeAccounts.value.filter(item => item.providerType === "MOFANG_CLOUD").length,
  upstream: activeAccounts.value.filter(item => item.providerType === "ZJMF_API").length,
  whmcs: activeAccounts.value.filter(item => item.providerType === "WHMCS").length,
  manual: activeAccounts.value.filter(item => item.providerType === "MANUAL").length
}));

const activeCount = computed(() => instances.value.filter(item => item.status === "ACTIVE").length);
const suspendedCount = computed(() => instances.value.filter(item => item.status === "SUSPENDED").length);
const provisioningCount = computed(() => instances.value.filter(item => item.status === "PROVISIONING").length);

function providerTypeLabel(value: string) {
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上游财务",
    WHMCS: "WHMCS",
    RESOURCE: "资源池",
    MANUAL: "人工交付"
  };
  return mapping[value] ?? (value || "-");
}

function providerTypeTag(value: string) {
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: "success",
    ZJMF_API: "primary",
    WHMCS: "warning",
    RESOURCE: "info",
    MANUAL: ""
  };
  return mapping[value] ?? "info";
}

function statusTagType(status: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "success",
    SUSPENDED: "warning",
    PROVISIONING: "info",
    TERMINATED: "danger",
    FAILED: "danger"
  };
  return mapping[status] ?? "info";
}

function syncStatusType(status: string) {
  const mapping: Record<string, string> = {
    SUCCESS: "success",
    FAILED: "danger",
    PENDING: "warning",
    RUNNING: "primary"
  };
  return mapping[status] ?? "info";
}

function syncStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    SUCCESS: "成功",
    FAILED: "失败",
    PENDING: "待执行",
    RUNNING: "执行中"
  };
  return mapping[status] ?? status;
}

function syncActionLabel(action: string) {
  const mapping: Record<string, string> = {
    pull_service: "拉取服务",
    pull_resource: "拉取资源",
    update_service: "更新服务",
    update_resource: "更新资源",
    create_service: "创建服务",
    create_resource: "创建资源"
  };
  return mapping[action] ?? action;
}

function healthType() {
  return health.value?.connected ? "success" : "danger";
}

function ensureSelectedAccount() {
  if (selectedAccountId.value && filteredAccounts.value.some(item => item.id === selectedAccountId.value)) return;
  selectedAccountId.value = filteredAccounts.value[0]?.id ?? 0;
}

function applyRouteContext() {
  const providerTypeQuery = typeof route.query.providerType === "string" ? route.query.providerType : "";
  const accountIdQuery = typeof route.query.accountId === "string" ? Number(route.query.accountId) : NaN;
  const serviceIdQuery = typeof route.query.serviceId === "string" ? Number(route.query.serviceId) : NaN;

  selectedProviderType.value = providerTypeQuery || "";
  focusServiceId.value = Number.isFinite(serviceIdQuery) && serviceIdQuery > 0 ? serviceIdQuery : 0;

  if (Number.isFinite(accountIdQuery) && accountIdQuery > 0) {
    selectedAccountId.value = accountIdQuery;
  }
}

async function loadAccounts() {
  accounts.value = await fetchProviderAccounts();
  ensureSelectedAccount();
}

async function loadHealth() {
  if (!selectedAccountId.value || !supportsRemoteResources.value) {
    health.value = null;
    return;
  }
  health.value = await fetchMofangHealth(selectedAccountId.value);
}

async function loadInstances() {
  if (!selectedAccountId.value || !supportsRemoteResources.value) {
    instances.value = [];
    return;
  }
  const result = await fetchMofangInstances(selectedAccountId.value);
  instances.value = result.items;
}

async function loadSyncLogs() {
  const result = await fetchMofangSyncLogs({
    limit: 20,
    serviceId: focusServiceId.value || undefined
  });
  syncLogs.value = result.items;
}

async function refreshAll() {
  loading.value = true;
  try {
    await Promise.all([loadHealth(), loadInstances(), loadSyncLogs()]);
  } finally {
    loading.value = false;
  }
}

async function syncRemote(remoteIds?: string[]) {
  if (!selectedAccount.value) {
    ElMessage.warning("请先选择接口账户");
    return;
  }
  if (!supportsRemoteResources.value) {
    ElMessage.info("当前账户类型没有远端实例资源面板，请在接口账户、商品和服务工作台继续管理。");
    return;
  }

  syncLoading.value = true;
  try {
    const result = await pullMofangSync({
      providerAccountId: selectedAccount.value.id,
      includeResources: true,
      remoteIds: remoteIds && remoteIds.length > 0 ? remoteIds : undefined
    });
    ElMessage.success(
      `渠道资源同步完成：处理 ${result.summary.processedCount} 台，新增 ${result.summary.createdServices} 台，更新 ${result.summary.updatedServices} 台`
    );
    selectedRemoteIds.value = [];
    await refreshAll();
  } finally {
    syncLoading.value = false;
  }
}

function handleSelectionChange(rows: MofangInstanceSummary[]) {
  selectedRemoteIds.value = rows.map(item => item.remoteId);
}

function openConsole(url?: string) {
  if (!url) {
    ElMessage.info("当前资源没有返回控制台地址，请先同步到本地服务后从服务详情进入。");
    return;
  }
  window.open(url, "_blank", "noopener,noreferrer");
}

function openAccountWorkbench() {
  const query: Record<string, string> = {};
  if (selectedProviderType.value) query.providerType = selectedProviderType.value;
  if (selectedAccountId.value) query.accountId = String(selectedAccountId.value);
  void router.push({ path: "/providers/accounts", query });
}

function openServiceWorkbench(serviceId?: number) {
  if (!serviceId) return;
  void router.push(`/services/detail/${serviceId}`);
}

function openAutomationForService(serviceId?: number) {
  if (!serviceId) return;
  void router.push({
    path: "/providers/automation",
    query: {
      serviceId: String(serviceId),
      channel: selectedAccount.value?.providerType || selectedProviderType.value || undefined
    }
  });
}

function openAutomationForAccount() {
  if (!selectedAccount.value) return;
  void router.push({
    path: "/providers/automation",
    query: {
      sourceType: "provider",
      sourceId: String(selectedAccount.value.id),
      channel: selectedAccount.value.providerType
    }
  });
}

watch(selectedProviderType, () => {
  ensureSelectedAccount();
});

watch(selectedAccountId, () => {
  void refreshAll();
});

onMounted(async () => {
  applyRouteContext();
  await loadAccounts();
  await refreshAll();
});

watch(
  () => [route.query.providerType, route.query.accountId, route.query.serviceId],
  () => {
    applyRouteContext();
    ensureSelectedAccount();
    void refreshAll();
  }
);
</script>

<template>
  <PageWorkbench
    eyebrow="接口与上游 / 资源"
    title="渠道资源"
    subtitle="按接口账户查看远端实例、同步日志和资源拉取结果。魔方云走远端实例面板，其他渠道继续通过商品、服务和自动化工作台承接。"
  >
    <template #actions>
      <el-select v-model="selectedProviderType" clearable placeholder="按渠道筛选账户" style="width: 180px">
        <el-option label="魔方云" value="MOFANG_CLOUD" />
        <el-option label="上游财务" value="ZJMF_API" />
        <el-option label="WHMCS" value="WHMCS" />
        <el-option label="资源池" value="RESOURCE" />
        <el-option label="人工交付" value="MANUAL" />
      </el-select>
      <el-select v-model="selectedAccountId" placeholder="选择接口账户" style="width: 340px">
        <el-option
          v-for="item in filteredAccounts"
          :key="item.id"
          :label="`${item.name} / ${providerTypeLabel(item.providerType)}`"
          :value="item.id"
        />
      </el-select>
      <el-button @click="openAccountWorkbench">接口账户</el-button>
      <el-button @click="openAutomationForAccount" :disabled="!selectedAccount">自动化任务</el-button>
      <el-button @click="refreshAll">刷新</el-button>
      <el-button
        type="primary"
        :disabled="!selectedAccount || !supportsRemoteResources"
        :loading="syncLoading"
        @click="syncRemote(selectedRemoteIds.length ? selectedRemoteIds : undefined)"
      >
        {{ selectedRemoteIds.length ? "同步选中资源" : "全量拉取资源" }}
      </el-button>
    </template>

    <template #metrics>
      <div class="summary-strip">
        <div class="summary-pill"><span>活跃接口账户</span><strong>{{ providerSummary.total }}</strong></div>
        <div class="summary-pill"><span>魔方云</span><strong>{{ providerSummary.mofang }}</strong></div>
        <div class="summary-pill"><span>上游财务</span><strong>{{ providerSummary.upstream }}</strong></div>
        <div class="summary-pill"><span>WHMCS</span><strong>{{ providerSummary.whmcs }}</strong></div>
        <div class="summary-pill"><span>人工交付</span><strong>{{ providerSummary.manual }}</strong></div>
      </div>
    </template>

    <div class="page-card" v-if="focusServiceId">
      <div class="context-panel">
        <div>
          <strong>当前正围绕服务 {{ focusServiceId }} 查看渠道资源</strong>
          <p class="page-subtitle" style="margin-top: 6px">
            同步日志已按该服务过滤，你可以直接回到服务工作台，或继续查看同一服务的自动化任务。
          </p>
        </div>
        <div class="inline-actions">
          <el-button plain @click="openServiceWorkbench(focusServiceId)">查看服务</el-button>
          <el-button plain @click="openAutomationForService(focusServiceId)">查看自动化任务</el-button>
          <el-button type="primary" plain @click="router.push({ path: '/providers/resources', query: { providerType: selectedProviderType || undefined, accountId: selectedAccountId || undefined } })">
            退出服务上下文
          </el-button>
        </div>
      </div>
    </div>

    <div class="portal-grid portal-grid--two" v-loading="loading">
      <div class="panel-card">
        <div class="section-card__head">
          <strong>当前账户</strong>
          <el-tag v-if="selectedAccount" :type="providerTypeTag(selectedAccount.providerType)" effect="light">
            {{ providerTypeLabel(selectedAccount.providerType) }}
          </el-tag>
        </div>

        <el-descriptions v-if="selectedAccount" :column="2" border>
          <el-descriptions-item label="账户名称">{{ selectedAccount.name }}</el-descriptions-item>
          <el-descriptions-item label="账户状态">{{ selectedAccount.status }}</el-descriptions-item>
          <el-descriptions-item label="基础地址">{{ selectedAccount.baseUrl || "-" }}</el-descriptions-item>
          <el-descriptions-item label="渠道类型">{{ providerTypeLabel(selectedAccount.providerType) }}</el-descriptions-item>
          <el-descriptions-item label="来源名称">{{ selectedAccount.sourceName || "-" }}</el-descriptions-item>
          <el-descriptions-item label="联系方式">{{ selectedAccount.contactWay || "-" }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="说明">{{ selectedAccount.description || "未填写说明" }}</el-descriptions-item>
        </el-descriptions>
        <el-empty v-else description="暂无可用接口账户" />
      </div>

      <div class="panel-card">
        <div class="section-card__head">
          <strong>资源面板说明</strong>
          <span class="section-card__meta">渠道资源只是自动化链路的一环，商品绑定哪个账户、服务落在哪个渠道，才是后台主线。</span>
        </div>
        <el-alert
          v-if="supportsRemoteResources"
          title="当前账户支持远端实例资源面板"
          type="success"
          :closable="false"
          show-icon
        >
          当前会展示远端实例列表、连接状态和同步日志，适合排查实例是否已同步到本地服务。
        </el-alert>
        <el-alert
          v-else
          title="当前账户类型没有远端实例面板"
          type="info"
          :closable="false"
          show-icon
        >
          上游财务、WHMCS、资源池和人工交付账户，主要通过接口账户页、商品工作台、服务工作台和自动化任务中心配合管理。
        </el-alert>
      </div>
    </div>

    <div class="portal-grid portal-grid--two" style="margin-top: 16px" v-loading="loading">
      <div class="panel-card">
        <div class="section-card__head">
          <strong>连接健康</strong>
          <el-tag :type="healthType()" effect="light">
            {{ health?.connected ? "连接正常" : "未连接" }}
          </el-tag>
        </div>

        <el-descriptions v-if="supportsRemoteResources && health" :column="2" border>
          <el-descriptions-item label="基础地址">{{ health.baseUrl || "-" }}</el-descriptions-item>
          <el-descriptions-item label="实际地址">{{ health.activeUrl || "-" }}</el-descriptions-item>
          <el-descriptions-item label="认证方式">{{ health.authMode || "-" }}</el-descriptions-item>
          <el-descriptions-item label="最近认证">{{ health.lastAuthAt || "-" }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="状态说明">{{ health.message }}</el-descriptions-item>
        </el-descriptions>
        <el-empty v-else description="当前账户没有远端连接检查" />
      </div>

      <div class="panel-card">
        <div class="section-card__head">
          <strong>资源摘要</strong>
          <span class="section-card__meta">这里汇总当前账户拉取到的远端资源状态。</span>
        </div>
        <div class="summary-strip summary-strip--compact">
          <div class="summary-pill"><span>远端资源总数</span><strong>{{ instances.length }}</strong></div>
          <div class="summary-pill"><span>运行中</span><strong>{{ activeCount }}</strong></div>
          <div class="summary-pill"><span>已暂停</span><strong>{{ suspendedCount }}</strong></div>
          <div class="summary-pill"><span>处理中</span><strong>{{ provisioningCount }}</strong></div>
        </div>
      </div>
    </div>

    <div class="page-card" style="margin-top: 16px" v-loading="loading">
      <div class="page-header">
        <div>
          <h2 class="section-title">远端资源列表</h2>
          <p class="page-subtitle">仅对支持实例面板的渠道展示远端资源，选中后可单独同步到本地。</p>
        </div>
      </div>

      <el-table
        v-if="supportsRemoteResources"
        :data="instances"
        border
        stripe
        empty-text="当前账户暂无远端资源"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column prop="remoteId" label="远端编号" min-width="120" />
        <el-table-column prop="name" label="资源名称" min-width="220" />
        <el-table-column label="状态" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" effect="light">{{ row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="region" label="区域 / 节点" min-width="180" />
        <el-table-column prop="ipAddress" label="公网 IP" min-width="140" />
        <el-table-column prop="expiresAt" label="到期时间" min-width="180" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="syncRemote([row.remoteId])">同步这台</el-button>
              <el-button type="primary" link @click="openConsole(row.consoleUrl)">控制台</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-else description="当前账户类型没有远端资源列表，请改在接口账户页和服务工作台管理。" />
    </div>

    <div class="page-card" style="margin-top: 16px" v-loading="loading">
      <div class="page-header">
        <div>
          <h2 class="section-title">最近同步日志</h2>
          <p class="page-subtitle">用于确认资源拉取、服务落库、资源明细回写和失败原因。</p>
        </div>
      </div>

      <el-table :data="syncLogs" border stripe empty-text="暂无同步日志">
        <el-table-column prop="createdAt" label="时间" min-width="180" />
        <el-table-column label="渠道" min-width="120">
          <template #default="{ row }">{{ providerTypeLabel(row.providerType) }}</template>
        </el-table-column>
        <el-table-column label="动作" min-width="160">
          <template #default="{ row }">{{ syncActionLabel(row.action) }}</template>
        </el-table-column>
        <el-table-column prop="resourceType" label="资源类型" min-width="140" />
        <el-table-column prop="resourceId" label="远端编号" min-width="140" />
        <el-table-column prop="serviceId" label="本地服务 ID" min-width="120" />
        <el-table-column label="状态" min-width="110">
          <template #default="{ row }">
            <el-tag :type="syncStatusType(row.status)" effect="light">{{ syncStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="消息" min-width="260" show-overflow-tooltip />
        <el-table-column label="操作" min-width="240" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button v-if="row.serviceId" type="primary" link @click="openServiceWorkbench(row.serviceId)">服务</el-button>
              <el-button v-if="row.serviceId" type="primary" link @click="openAutomationForService(row.serviceId)">任务</el-button>
              <el-button type="primary" link @click="router.push({ path: '/providers/accounts', query: { providerType: row.providerType } })">
                账户
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </PageWorkbench>
</template>

<style scoped>
.context-panel {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.inline-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}
</style>
