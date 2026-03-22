<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { ElMessage } from "element-plus";
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

const loading = ref(false);
const syncLoading = ref(false);
const selectedRemoteIds = ref<string[]>([]);
const selectedAccountId = ref<number>(0);
const health = ref<MofangHealthResponse | null>(null);
const instances = ref<MofangInstanceSummary[]>([]);
const syncLogs = ref<MofangSyncLogItem[]>([]);
const accounts = ref<ProviderAccount[]>([]);

const mofangAccounts = computed(() =>
  accounts.value.filter(item => item.providerType === "MOFANG_CLOUD" && item.status === "ACTIVE")
);

const activeCount = computed(() => instances.value.filter(item => item.status === "ACTIVE").length);
const suspendedCount = computed(() => instances.value.filter(item => item.status === "SUSPENDED").length);
const provisioningCount = computed(() => instances.value.filter(item => item.status === "PROVISIONING").length);

async function loadAccounts() {
  accounts.value = await fetchProviderAccounts("MOFANG_CLOUD");
  if (!selectedAccountId.value && mofangAccounts.value.length > 0) {
    selectedAccountId.value = mofangAccounts.value[0].id;
  }
}

async function loadHealth() {
  if (!selectedAccountId.value) {
    health.value = null;
    return;
  }
  health.value = await fetchMofangHealth(selectedAccountId.value);
}

async function loadInstances() {
  if (!selectedAccountId.value) {
    instances.value = [];
    return;
  }
  const result = await fetchMofangInstances(selectedAccountId.value);
  instances.value = result.items;
}

async function loadSyncLogs() {
  const result = await fetchMofangSyncLogs({ limit: 20 });
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
  if (!selectedAccountId.value) {
    ElMessage.warning("请先选择魔方云接口账户");
    return;
  }

  syncLoading.value = true;
  try {
    const result = await pullMofangSync({
      providerAccountId: selectedAccountId.value,
      includeResources: true,
      remoteIds: remoteIds && remoteIds.length > 0 ? remoteIds : undefined
    });
    ElMessage.success(
      `同步完成：处理 ${result.summary.processedCount} 台，新增 ${result.summary.createdServices} 台，更新 ${result.summary.updatedServices} 台`
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
    ElMessage.info("当前实例列表未返回控制台地址，请先同步到本地服务后从服务详情进入。");
    return;
  }
  window.open(url, "_blank", "noopener,noreferrer");
}

function healthType() {
  return health.value?.connected ? "success" : "danger";
}

function statusType(status: string) {
  const mapping: Record<string, string> = {
    ACTIVE: "success",
    SUSPENDED: "warning",
    PROVISIONING: "info",
    TERMINATED: "danger",
    FAILED: "danger"
  };
  return mapping[status] ?? "info";
}

watch(selectedAccountId, () => {
  void refreshAll();
});

onMounted(async () => {
  await loadAccounts();
  await refreshAll();
});
</script>

<template>
  <div class="page-card workbench-shell" v-loading="loading">
    <div class="workbench-hero">
      <div>
        <h1 class="page-title">魔方云</h1>
        <p class="page-subtitle">按接口账户查看远端实例、执行同步，并把服务和资源明细回写到本地系统。</p>
      </div>
      <div class="action-group">
        <el-select v-model="selectedAccountId" placeholder="选择魔方云接口账户" style="width: 320px">
          <el-option
            v-for="item in mofangAccounts"
            :key="item.id"
            :label="`${item.name} / ${item.baseUrl}`"
            :value="item.id"
          />
        </el-select>
        <el-button @click="refreshAll">刷新状态</el-button>
        <el-button
          type="primary"
          :loading="syncLoading"
          @click="syncRemote(selectedRemoteIds.length > 0 ? selectedRemoteIds : undefined)"
        >
          {{ selectedRemoteIds.length > 0 ? "同步选中实例" : "全量同步" }}
        </el-button>
      </div>
    </div>

    <div class="summary-strip">
      <div class="summary-pill"><span>接口账户</span><strong>{{ selectedAccountId || "-" }}</strong></div>
      <div class="summary-pill"><span>远端实例总数</span><strong>{{ instances.length }}</strong></div>
      <div class="summary-pill"><span>运行中</span><strong>{{ activeCount }}</strong></div>
      <div class="summary-pill"><span>已暂停</span><strong>{{ suspendedCount }}</strong></div>
      <div class="summary-pill"><span>处理中</span><strong>{{ provisioningCount }}</strong></div>
      <div class="summary-pill"><span>连接状态</span><strong>{{ health?.connected ? "正常" : "异常" }}</strong></div>
    </div>

    <div class="portal-grid portal-grid--two" style="margin-top: 16px">
      <div class="panel-card">
        <div class="section-card__head">
          <strong>连接健康</strong>
          <el-tag :type="healthType()" effect="light">
            {{ health?.connected ? "已连接" : "未连接" }}
          </el-tag>
        </div>
        <el-descriptions v-if="health" :column="2" border>
          <el-descriptions-item label="基础地址">{{ health.baseUrl || "-" }}</el-descriptions-item>
          <el-descriptions-item label="实际地址">{{ health.activeUrl || "-" }}</el-descriptions-item>
          <el-descriptions-item label="认证方式">{{ health.authMode || "-" }}</el-descriptions-item>
          <el-descriptions-item label="最近认证">{{ health.lastAuthAt || "-" }}</el-descriptions-item>
          <el-descriptions-item :span="2" label="状态说明">{{ health.message }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="panel-card">
        <div class="section-card__head">
          <strong>同步说明</strong>
          <span class="section-card__meta">魔方云应作为商品自动化渠道之一，不同商品绑定不同接口账户。</span>
        </div>
        <el-alert
          title="这里看到的是接口账户维度的远端实例。真正自动开通、同步、资源动作，都应落在绑定该接口账户的商品和服务上。"
          type="info"
          :closable="false"
          show-icon
        />
      </div>
    </div>

    <div class="page-card" style="margin-top: 16px">
      <div class="page-header">
        <div>
          <h2 class="section-title">远端实例列表</h2>
          <p class="page-subtitle">支持按接口账户局部同步，也支持对单台实例进行手工拉取。</p>
        </div>
      </div>

      <el-table :data="instances" border stripe @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="48" />
        <el-table-column prop="remoteId" label="远端编号" min-width="120" />
        <el-table-column prop="name" label="实例名称" min-width="220" />
        <el-table-column label="状态" min-width="110">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="region" label="地区/节点" min-width="180" />
        <el-table-column prop="ipAddress" label="公网 IP" min-width="140" />
        <el-table-column prop="expiresAt" label="到期时间" min-width="180" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="syncRemote([row.remoteId])">同步这台</el-button>
            <el-button type="primary" link @click="openConsole(row.consoleUrl)">控制台</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <div class="page-card" style="margin-top: 16px">
      <div class="page-header">
        <div>
          <h2 class="section-title">最近同步日志</h2>
          <p class="page-subtitle">用于确认服务、资源、同步链是否已回写到本地。</p>
        </div>
      </div>

      <el-table :data="syncLogs" border stripe empty-text="暂无同步日志">
        <el-table-column prop="createdAt" label="时间" min-width="180" />
        <el-table-column prop="action" label="动作" min-width="160" />
        <el-table-column prop="resourceId" label="远端实例" min-width="140" />
        <el-table-column prop="serviceId" label="本地服务 ID" min-width="120" />
        <el-table-column prop="status" label="状态" min-width="110" />
        <el-table-column prop="message" label="消息" min-width="260" />
      </el-table>
    </div>
  </div>
</template>
