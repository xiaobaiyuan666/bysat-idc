<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createProviderAccount,
  deleteProviderAccount,
  fetchAutomationTasks,
  fetchMofangInstances,
  fetchProducts,
  fetchProviderAccountHealth,
  fetchProviderAccounts,
  fetchServices,
  fetchUpstreamProducts,
  importUpstreamProducts,
  pullMofangSync,
  updateProviderAccount,
  type AutomationTask,
  type MofangHealthResponse,
  type MofangInstanceSummary,
  type Product,
  type ProviderAccount,
  type ServiceRecord,
  type UpstreamCatalogGroup
} from "@/api/admin";

type ProviderType = "MOFANG_CLOUD" | "ZJMF_API" | "WHMCS" | "RESOURCE" | "MANUAL";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const saving = ref(false);
const workbenchLoading = ref(false);
const healthLoading = ref(false);
const syncLoading = ref(false);
const upstreamSyncLoading = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);

const accounts = ref<ProviderAccount[]>([]);
const selectedAccountId = ref<number | null>(null);
const selectedHealth = ref<MofangHealthResponse | null>(null);
const relatedProducts = ref<Product[]>([]);
const relatedServices = ref<ServiceRecord[]>([]);
const relatedTasks = ref<AutomationTask[]>([]);
const remoteInstances = ref<MofangInstanceSummary[]>([]);
const remoteCatalogGroups = ref<UpstreamCatalogGroup[]>([]);

const filters = reactive({
  providerType: "",
  keyword: ""
});

const form = reactive<Partial<ProviderAccount>>({
  providerType: "MOFANG_CLOUD",
  name: "",
  baseUrl: "",
  username: "",
  password: "",
  sourceName: "",
  contactWay: "",
  description: "",
  accountMode: "API",
  lang: "zh-cn",
  listPath: "/v1/clouds",
  detailPath: "/v1/clouds/:id",
  insecureSkipVerify: true,
  autoUpdate: true,
  status: "ACTIVE",
  extraConfig: "{}"
});

const providerTypeOptions: Array<{ label: string; value: ProviderType }> = [
  { label: "魔方云", value: "MOFANG_CLOUD" },
  { label: "上游财务", value: "ZJMF_API" },
  { label: "WHMCS", value: "WHMCS" },
  { label: "资源池", value: "RESOURCE" },
  { label: "人工交付", value: "MANUAL" }
];

const filteredAccounts = computed(() =>
  accounts.value.filter(item => {
    const keyword = filters.keyword.trim().toLowerCase();
    const matchesType = !filters.providerType || item.providerType === filters.providerType;
    const matchesKeyword =
      !keyword ||
      String(item.name).toLowerCase().includes(keyword) ||
      String(item.baseUrl).toLowerCase().includes(keyword) ||
      String(item.sourceName).toLowerCase().includes(keyword);
    return matchesType && matchesKeyword;
  })
);

const selectedAccount = computed(() => accounts.value.find(item => item.id === selectedAccountId.value) ?? null);

const remoteProductPreview = computed(() =>
  remoteCatalogGroups.value.flatMap(group =>
    group.products.slice(0, 6).map(product => ({
      groupName: group.groupName,
      remoteProductCode: product.remoteProductCode,
      name: product.name,
      productType: product.productType
    }))
  )
);

const summary = computed(() => ({
  total: accounts.value.length,
  active: accounts.value.filter(item => item.status === "ACTIVE").length,
  mofang: accounts.value.filter(item => item.providerType === "MOFANG_CLOUD").length,
  upstream: accounts.value.filter(item => item.providerType === "ZJMF_API").length
}));

const workbenchSummary = computed(() => ({
  products: relatedProducts.value.length,
  services: relatedServices.value.length,
  tasks: relatedTasks.value.length,
  remoteProducts: remoteCatalogGroups.value.reduce((total, group) => total + group.products.length, 0),
  remoteGroups: remoteCatalogGroups.value.length,
  remoteInstances: remoteInstances.value.length
}));

function providerTypeLabel(value: string) {
  return (
    {
      MOFANG_CLOUD: "魔方云",
      ZJMF_API: "上游财务",
      WHMCS: "WHMCS",
      RESOURCE: "资源池",
      MANUAL: "人工交付"
    }[value] ?? value
  );
}

function providerCapabilityLabel(value: string) {
  return (
    {
      MOFANG_CLOUD: "实例拉取 / 资源同步 / 资源动作",
      ZJMF_API: "商品目录 / 模板 / 价格同步",
      WHMCS: "上游商品目录 / 套餐映射",
      RESOURCE: "资源库存 / 手工分配",
      MANUAL: "手工交付 / 工单协同"
    }[value] ?? "通用接口能力"
  );
}

function accountStatusType(status: string) {
  return status === "ACTIVE" ? "success" : status === "DISABLED" ? "warning" : "info";
}

function serviceStatusType(status: string) {
  return (
    {
      PENDING: "warning",
      ACTIVE: "success",
      SUSPENDED: "warning",
      PROVISIONING: "primary",
      TERMINATED: "danger",
      FAILED: "danger"
    }[status] ?? "info"
  );
}

function taskStatusType(status: string) {
  return (
    {
      PENDING: "warning",
      RUNNING: "primary",
      SUCCESS: "success",
      FAILED: "danger",
      BLOCKED: "warning"
    }[status] ?? "info"
  );
}

function taskStatusLabel(status: string) {
  return (
    {
      PENDING: "待执行",
      RUNNING: "执行中",
      SUCCESS: "成功",
      FAILED: "失败",
      BLOCKED: "阻塞"
    }[status] ?? status
  );
}

function resetForm() {
  editingId.value = null;
  form.providerType = "MOFANG_CLOUD";
  form.name = "";
  form.baseUrl = "";
  form.username = "";
  form.password = "";
  form.sourceName = "";
  form.contactWay = "";
  form.description = "";
  form.accountMode = "API";
  form.lang = "zh-cn";
  form.listPath = "/v1/clouds";
  form.detailPath = "/v1/clouds/:id";
  form.insecureSkipVerify = true;
  form.autoUpdate = true;
  form.status = "ACTIVE";
  form.extraConfig = "{}";
}

function applyProviderDefaults(providerType: string) {
  if (providerType === "MOFANG_CLOUD") {
    form.accountMode = "API";
    form.listPath = "/v1/clouds";
    form.detailPath = "/v1/clouds/:id";
    form.lang = "zh-cn";
    return;
  }
  if (providerType === "ZJMF_API") {
    form.accountMode = "UPSTREAM_FINANCE";
    form.listPath = "";
    form.detailPath = "";
    form.lang = "zh-cn";
    return;
  }
  form.listPath = "";
  form.detailPath = "";
}

function syncFiltersFromRoute() {
  filters.providerType = typeof route.query.providerType === "string" ? route.query.providerType : "";
}

function applyRouteSelection() {
  const routeAccountId = Number(route.query.accountId || 0);
  const pool = filters.providerType
    ? filteredAccounts.value.filter(item => item.providerType === filters.providerType)
    : filteredAccounts.value;

  if (routeAccountId > 0) {
    const matched = pool.find(item => item.id === routeAccountId) ?? accounts.value.find(item => item.id === routeAccountId);
    if (matched) {
      selectedAccountId.value = matched.id;
      return matched;
    }
  }

  if (!selectedAccountId.value || !pool.some(item => item.id === selectedAccountId.value)) {
    const fallback = pool[0] ?? filteredAccounts.value[0] ?? accounts.value[0] ?? null;
    selectedAccountId.value = fallback?.id ?? null;
    return fallback;
  }

  return selectedAccount.value;
}

async function loadAccounts() {
  loading.value = true;
  try {
    syncFiltersFromRoute();
    accounts.value = await fetchProviderAccounts();
    const next = applyRouteSelection();
    if (next) {
      await openWorkbench(next);
    }
  } finally {
    loading.value = false;
  }
}

async function openWorkbench(account?: ProviderAccount | null) {
  if (!account) return;
  selectedAccountId.value = account.id;
  workbenchLoading.value = true;
  try {
    const [productsData, servicesData, taskData, healthData] = await Promise.all([
      fetchProducts(),
      fetchServices({ provider_account_id: account.id, limit: 20, sort: "created_at", order: "desc" }),
      fetchAutomationTasks({
        sourceType: "provider",
        sourceId: account.id,
        channel: account.providerType,
        limit: 20
      }),
      fetchProviderAccountHealth(account.id).catch(() => null)
    ]);

    relatedProducts.value = productsData.items.filter(
      item =>
        item.automationConfig.providerAccountId === account.id ||
        item.upstreamMapping.providerAccountId === account.id
    );
    relatedServices.value = servicesData.items;
    relatedTasks.value = taskData.items;
    selectedHealth.value = healthData;

    if (account.providerType === "MOFANG_CLOUD") {
      const remote = await fetchMofangInstances(account.id);
      remoteInstances.value = remote.items;
      remoteCatalogGroups.value = [];
    } else if (account.providerType === "ZJMF_API" || account.providerType === "WHMCS") {
      const remote = await fetchUpstreamProducts(account.id);
      remoteCatalogGroups.value = remote.groups;
      remoteInstances.value = [];
    } else {
      remoteInstances.value = [];
      remoteCatalogGroups.value = [];
    }
  } finally {
    workbenchLoading.value = false;
  }
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(record: ProviderAccount) {
  editingId.value = record.id;
  Object.assign(form, record);
  dialogVisible.value = true;
}

async function saveAccount() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateProviderAccount(editingId.value, form);
      ElMessage.success("接口账户已更新");
    } else {
      await createProviderAccount(form);
      ElMessage.success("接口账户已创建");
    }
    dialogVisible.value = false;
    await loadAccounts();
  } finally {
    saving.value = false;
  }
}

async function removeAccount(record: ProviderAccount) {
  await ElMessageBox.confirm(`确认删除接口账户“${record.name}”吗？`, "删除确认", { type: "warning" });
  await deleteProviderAccount(record.id);
  ElMessage.success("接口账户已删除");
  if (selectedAccountId.value === record.id) {
    selectedAccountId.value = null;
  }
  await loadAccounts();
}

async function checkHealth(record: ProviderAccount) {
  healthLoading.value = true;
  try {
    selectedAccountId.value = record.id;
    selectedHealth.value = await fetchProviderAccountHealth(record.id);
    ElMessage.success(`已完成 ${record.name} 的连接检测`);
  } finally {
    healthLoading.value = false;
  }
}

async function syncMofangAccount() {
  if (!selectedAccount.value || selectedAccount.value.providerType !== "MOFANG_CLOUD") return;
  syncLoading.value = true;
  try {
    const result = await pullMofangSync({
      providerAccountId: selectedAccount.value.id,
      includeResources: true
    });
    ElMessage.success(
      `渠道资源同步完成：处理 ${result.summary.processedCount} 台，新增 ${result.summary.createdServices} 台，更新 ${result.summary.updatedServices} 台`
    );
    await openWorkbench(selectedAccount.value);
  } finally {
    syncLoading.value = false;
  }
}

async function syncUpstreamCatalog() {
  if (!selectedAccount.value || !["ZJMF_API", "WHMCS"].includes(selectedAccount.value.providerType)) return;
  upstreamSyncLoading.value = true;
  try {
    const result = await importUpstreamProducts({
      providerAccountId: selectedAccount.value.id,
      providerType: selectedAccount.value.providerType,
      importAll: true,
      remoteProductCodes: [],
      autoSyncPricing: true,
      autoSyncConfig: true,
      autoSyncTemplate: true,
      deactivateMissing: true
    });
    ElMessage.success(`上游商品同步完成：新增 ${result.created}，更新 ${result.updated}，停用 ${result.deactivatedCount || 0}，失败 ${result.failed}`);
    await openWorkbench(selectedAccount.value);
    if (result.historyId) {
      void router.push({
        path: "/providers/upstream-sync-history",
        query: {
          providerAccountId: String(selectedAccount.value.id),
          providerType: selectedAccount.value.providerType,
          historyId: result.historyId
        }
      });
    }
  } finally {
    upstreamSyncLoading.value = false;
  }
}

function jumpToProducts() {
  if (!selectedAccount.value) return;
  void router.push({
    path: "/catalog/products",
    query: {
      providerAccountId: String(selectedAccount.value.id),
      providerType: selectedAccount.value.providerType
    }
  });
}

function jumpToServices() {
  if (!selectedAccount.value) return;
  void router.push({
    path: "/services/list",
    query: {
      providerAccountId: String(selectedAccount.value.id),
      providerType: selectedAccount.value.providerType
    }
  });
}

function jumpToTasks() {
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

function jumpToResources() {
  if (!selectedAccount.value) return;
  void router.push({
    path: "/providers/resources",
    query: {
      providerType: selectedAccount.value.providerType,
      accountId: String(selectedAccount.value.id)
    }
  });
}

function jumpToUpstreamHistory() {
  if (!selectedAccount.value) return;
  void router.push({
    path: "/providers/upstream-sync-history",
    query: {
      providerAccountId: String(selectedAccount.value.id),
      providerType: selectedAccount.value.providerType
    }
  });
}

function openProductDetail(productId?: number) {
  if (!productId) return;
  void router.push(`/catalog/products/${productId}`);
}

function openServiceDetail(serviceId?: number) {
  if (!serviceId) return;
  void router.push(`/services/detail/${serviceId}`);
}

function openTaskWorkbench(task: Partial<AutomationTask>) {
  const query: Record<string, string> = {};
  if (task.taskNo) query.keyword = task.taskNo;
  if (task.serviceId) query.serviceId = String(task.serviceId);
  if (task.orderId) query.orderId = String(task.orderId);
  if (task.invoiceId) query.invoiceId = String(task.invoiceId);
  if (task.channel) query.channel = task.channel;
  void router.push({ path: "/providers/automation", query });
}

watch(
  () => form.providerType,
  value => {
    if (value) {
      applyProviderDefaults(value);
    }
  }
);

watch(
  () => route.fullPath,
  async () => {
    syncFiltersFromRoute();
    const next = applyRouteSelection();
    if (next) {
      await openWorkbench(next);
    }
  }
);

watch(
  () => [filters.providerType, filters.keyword],
  async () => {
    const next = applyRouteSelection();
    if (next) {
      await openWorkbench(next);
    }
  }
);

onMounted(() => {
  void loadAccounts();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="接口与上游 / 账户"
      title="接口账户工作台"
      subtitle="统一管理魔方云、上游财务、WHMCS、资源池和人工交付账户，并直接查看它们影响到的商品、服务和自动化任务。"
    >
      <template #actions>
        <el-input v-model="filters.keyword" placeholder="搜索名称、来源或地址" clearable style="width: 240px" />
        <el-select v-model="filters.providerType" placeholder="接口类型" clearable style="width: 180px">
          <el-option v-for="item in providerTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-button @click="loadAccounts">刷新</el-button>
        <el-button type="primary" @click="openCreate">新增接口账户</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>账户总数</span><strong>{{ summary.total }}</strong></div>
          <div class="summary-pill"><span>启用中</span><strong>{{ summary.active }}</strong></div>
          <div class="summary-pill"><span>魔方云</span><strong>{{ summary.mofang }}</strong></div>
          <div class="summary-pill"><span>上游财务</span><strong>{{ summary.upstream }}</strong></div>
        </div>
      </template>

      <div class="provider-grid">
        <section class="panel-card">
          <div class="section-card__head">
            <strong>账户列表</strong>
            <span class="section-card__meta">当前 {{ filteredAccounts.length }} 条</span>
          </div>

          <el-table
            :data="filteredAccounts"
            border
            stripe
            highlight-current-row
            row-key="id"
            :current-row-key="selectedAccountId || undefined"
            @current-change="openWorkbench"
          >
            <el-table-column prop="name" label="账户名称" min-width="180" />
            <el-table-column label="类型" min-width="120">
              <template #default="{ row }">{{ providerTypeLabel(row.providerType) }}</template>
            </el-table-column>
            <el-table-column prop="baseUrl" label="接口地址" min-width="240" show-overflow-tooltip />
            <el-table-column prop="productCount" label="商品数" width="90" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="accountStatusType(row.status)" effect="light">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" min-width="220" fixed="right">
              <template #default="{ row }">
                <div class="inline-actions">
                  <el-button link type="primary" @click="openWorkbench(row)">工作台</el-button>
                  <el-button link type="primary" :loading="healthLoading" @click="checkHealth(row)">检测</el-button>
                  <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
                  <el-button link type="danger" @click="removeAccount(row)">删除</el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </section>

        <section class="panel-card" v-loading="workbenchLoading || healthLoading">
          <div class="section-card__head">
            <strong>{{ selectedAccount ? `${selectedAccount.name} 工作台` : "账户工作台" }}</strong>
            <div v-if="selectedAccount" class="inline-actions">
              <el-button plain @click="jumpToProducts">商品</el-button>
              <el-button plain @click="jumpToServices">服务</el-button>
              <el-button plain @click="jumpToTasks">任务</el-button>
              <el-button plain @click="jumpToResources">资源/日志</el-button>
              <el-button
                v-if="selectedAccount.providerType === 'MOFANG_CLOUD'"
                type="primary"
                :loading="syncLoading"
                @click="syncMofangAccount"
              >
                拉取渠道资源
              </el-button>
              <el-button
                v-if="['ZJMF_API', 'WHMCS'].includes(selectedAccount.providerType)"
                type="primary"
                :loading="upstreamSyncLoading"
                @click="syncUpstreamCatalog"
              >
                同步上游商品
              </el-button>
              <el-button v-if="['ZJMF_API', 'WHMCS'].includes(selectedAccount.providerType)" plain @click="jumpToUpstreamHistory">
                上游同步记录
              </el-button>
            </div>
          </div>

          <el-empty v-if="!selectedAccount" description="请先选择一套接口账户" :image-size="80" />

          <template v-else>
            <div class="summary-strip">
              <div class="summary-pill"><span>关联商品</span><strong>{{ workbenchSummary.products }}</strong></div>
              <div class="summary-pill"><span>关联服务</span><strong>{{ workbenchSummary.services }}</strong></div>
              <div class="summary-pill"><span>自动化任务</span><strong>{{ workbenchSummary.tasks }}</strong></div>
              <div v-if="selectedAccount.providerType === 'MOFANG_CLOUD'" class="summary-pill">
                <span>远端实例</span>
                <strong>{{ workbenchSummary.remoteInstances }}</strong>
              </div>
              <template v-else-if="['ZJMF_API', 'WHMCS'].includes(selectedAccount.providerType)">
                <div class="summary-pill"><span>远端分组</span><strong>{{ workbenchSummary.remoteGroups }}</strong></div>
                <div class="summary-pill"><span>远端商品</span><strong>{{ workbenchSummary.remoteProducts }}</strong></div>
              </template>
            </div>

            <el-descriptions :column="2" border style="margin-top: 16px">
              <el-descriptions-item label="账户类型">{{ providerTypeLabel(selectedAccount.providerType) }}</el-descriptions-item>
              <el-descriptions-item label="渠道能力">{{ providerCapabilityLabel(selectedAccount.providerType) }}</el-descriptions-item>
              <el-descriptions-item label="来源名称">{{ selectedAccount.sourceName || "-" }}</el-descriptions-item>
              <el-descriptions-item label="账户模式">{{ selectedAccount.accountMode || "-" }}</el-descriptions-item>
              <el-descriptions-item label="最近认证">{{ selectedHealth?.lastAuthAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="连接状态">
                <el-tag :type="selectedHealth?.connected ? 'success' : 'danger'" effect="light">
                  {{ selectedHealth?.connected ? "已连接" : "未连接" }}
                </el-tag>
              </el-descriptions-item>
              <el-descriptions-item label="生效地址">{{ selectedHealth?.activeUrl || selectedAccount.baseUrl }}</el-descriptions-item>
              <el-descriptions-item label="检测说明" :span="2">
                {{ selectedHealth?.message || "尚未执行连接检测" }}
              </el-descriptions-item>
            </el-descriptions>

            <div class="provider-stack">
              <div>
                <div class="section-card__head">
                  <strong>关联商品</strong>
                  <span class="section-card__meta">自动化绑定和上游映射都会计入</span>
                </div>
                <el-table :data="relatedProducts.slice(0, 8)" border stripe size="small" empty-text="暂无关联商品">
                  <el-table-column prop="productNo" label="商品编号" min-width="140" />
                  <el-table-column prop="name" label="商品名称" min-width="220" show-overflow-tooltip />
                  <el-table-column prop="status" label="状态" min-width="90" />
                  <el-table-column label="操作" min-width="120" fixed="right">
                    <template #default="{ row }">
                      <el-button type="primary" link @click="openProductDetail(row.id)">商品工作台</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>

              <div>
                <div class="section-card__head">
                  <strong>关联服务</strong>
                  <span class="section-card__meta">按接口账户过滤</span>
                </div>
                <el-table :data="relatedServices.slice(0, 8)" border stripe size="small" empty-text="暂无关联服务">
                  <el-table-column prop="serviceNo" label="服务编号" min-width="150" />
                  <el-table-column prop="productName" label="产品" min-width="220" show-overflow-tooltip />
                  <el-table-column label="状态" min-width="90">
                    <template #default="{ row }">
                      <el-tag :type="serviceStatusType(row.status)" effect="light">{{ row.status }}</el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" min-width="160" fixed="right">
                    <template #default="{ row }">
                      <div class="inline-actions">
                        <el-button type="primary" link @click="openServiceDetail(row.id)">服务工作台</el-button>
                        <el-button type="primary" link @click="openTaskWorkbench({ serviceId: row.id, channel: selectedAccount?.providerType })">
                          任务
                        </el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
              </div>

              <div>
                <div class="section-card__head">
                  <strong>账户级任务</strong>
                  <span class="section-card__meta">同步、交付和重试都会汇总到这里</span>
                </div>
                <el-table :data="relatedTasks.slice(0, 8)" border stripe size="small" empty-text="暂无账户级任务">
                  <el-table-column prop="taskNo" label="任务号" min-width="150" />
                  <el-table-column prop="title" label="标题" min-width="220" show-overflow-tooltip />
                  <el-table-column label="状态" min-width="100">
                    <template #default="{ row }">
                      <el-tag :type="taskStatusType(row.status)" effect="light">{{ taskStatusLabel(row.status) }}</el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column label="操作" min-width="220" fixed="right">
                    <template #default="{ row }">
                      <div class="inline-actions">
                        <el-button type="primary" link @click="openTaskWorkbench(row)">任务中心</el-button>
                        <el-button v-if="row.serviceId" type="primary" link @click="openServiceDetail(row.serviceId)">服务</el-button>
                        <el-button v-if="row.invoiceId" type="primary" link @click="router.push(`/billing/invoices/${row.invoiceId}`)">账单</el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
              </div>

              <div v-if="selectedAccount.providerType === 'MOFANG_CLOUD'">
                <div class="section-card__head">
                  <strong>渠道资源预览</strong>
                  <span class="section-card__meta">展示当前账户最近拉取到的云实例</span>
                </div>
                <el-table :data="remoteInstances.slice(0, 8)" border stripe size="small" empty-text="暂无远端实例">
                  <el-table-column prop="remoteId" label="实例 ID" min-width="110" />
                  <el-table-column prop="name" label="实例名称" min-width="180" show-overflow-tooltip />
                  <el-table-column prop="region" label="地域" min-width="120" />
                  <el-table-column prop="ipAddress" label="IP" min-width="140" />
                  <el-table-column label="操作" min-width="120" fixed="right">
                    <template #default>
                      <el-button type="primary" link @click="jumpToResources">资源页</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>

              <div v-else-if="['ZJMF_API', 'WHMCS'].includes(selectedAccount.providerType)">
                <div class="section-card__head">
                  <strong>远端商品预览</strong>
                  <span class="section-card__meta">展示当前接口拉取到的上游目录</span>
                </div>
                <el-table :data="remoteProductPreview.slice(0, 12)" border stripe size="small" empty-text="暂无远端商品">
                  <el-table-column prop="remoteProductCode" label="远端编码" min-width="110" />
                  <el-table-column prop="name" label="商品名称" min-width="220" show-overflow-tooltip />
                  <el-table-column prop="groupName" label="分组" min-width="160" show-overflow-tooltip />
                  <el-table-column prop="productType" label="类型" min-width="100" />
                  <el-table-column label="操作" min-width="120" fixed="right">
                    <template #default>
                      <el-button type="primary" link @click="jumpToProducts">商品列表</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </template>
        </section>
      </div>
    </PageWorkbench>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑接口账户' : '新增接口账户'" width="760px" destroy-on-close>
      <el-form label-position="top">
        <div class="form-grid">
          <el-form-item label="接口类型">
            <el-select v-model="form.providerType" style="width: 100%">
              <el-option v-for="item in providerTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="账户名称">
            <el-input v-model="form.name" placeholder="例如：小白云上游 / 魔方云主账号" />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="接口地址">
            <el-input v-model="form.baseUrl" placeholder="https://example.com/" />
          </el-form-item>
          <el-form-item label="来源名称">
            <el-input v-model="form.sourceName" placeholder="用于业务展示的来源名称" />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="用户名">
            <el-input v-model="form.username" />
          </el-form-item>
          <el-form-item label="密码 / 密钥">
            <el-input v-model="form.password" show-password />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="账户模式">
            <el-input v-model="form.accountMode" placeholder="例如 API / UPSTREAM_FINANCE" />
          </el-form-item>
          <el-form-item label="语言">
            <el-input v-model="form.lang" placeholder="zh-cn" />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="列表接口路径">
            <el-input v-model="form.listPath" placeholder="/v1/clouds" />
          </el-form-item>
          <el-form-item label="详情接口路径">
            <el-input v-model="form.detailPath" placeholder="/v1/clouds/:id" />
          </el-form-item>
        </div>

        <div class="form-grid">
          <el-form-item label="联系方式">
            <el-input v-model="form.contactWay" />
          </el-form-item>
          <el-form-item label="状态">
            <el-select v-model="form.status" style="width: 100%">
              <el-option label="ACTIVE" value="ACTIVE" />
              <el-option label="DISABLED" value="DISABLED" />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="说明">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>

        <el-form-item label="扩展配置 JSON">
          <el-input v-model="form.extraConfig" type="textarea" :rows="3" />
        </el-form-item>

        <div class="form-grid form-grid--switch">
          <el-form-item label="允许跳过 TLS 校验">
            <el-switch v-model="form.insecureSkipVerify" />
          </el-form-item>
          <el-form-item label="启用自动更新">
            <el-switch v-model="form.autoUpdate" />
          </el-form-item>
        </div>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveAccount">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.provider-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 1.35fr);
  gap: 20px;
}

.provider-stack {
  display: grid;
  gap: 16px;
  margin-top: 16px;
}

.inline-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.form-grid--switch {
  align-items: center;
}

@media (max-width: 1280px) {
  .provider-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
