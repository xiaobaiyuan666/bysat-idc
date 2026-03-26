<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchAutomationTaskDetail,
  fetchAutomationTasks,
  retryAutomationTask,
  type AutomationTask,
  type AutomationTaskListResponse
} from "@/api/admin";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const detailLoading = ref(false);
const retryingTaskId = ref(0);
const detailVisible = ref(false);
const response = ref<AutomationTaskListResponse | null>(null);
const detail = ref<AutomationTask | null>(null);

const filters = reactive({
  status: "",
  taskType: "",
  channel: "",
  stage: "",
  sourceType: "",
  sourceId: undefined as number | undefined,
  orderId: undefined as number | undefined,
  invoiceId: undefined as number | undefined,
  serviceId: undefined as number | undefined,
  keyword: ""
});

const statusOptions = [
  { label: "全部", value: "" },
  { label: "待执行", value: "PENDING" },
  { label: "执行中", value: "RUNNING" },
  { label: "成功", value: "SUCCESS" },
  { label: "失败", value: "FAILED" },
  { label: "阻塞", value: "BLOCKED" }
];

const taskTypeOptions = [
  { label: "全部", value: "" },
  { label: "自动开通", value: "AUTO_PROVISION" },
  { label: "服务动作", value: "SERVICE_ACTION" },
  { label: "服务同步", value: "PULL_SYNC_SERVICE" },
  { label: "批量同步", value: "PULL_SYNC_BATCH" },
  { label: "资源动作", value: "RESOURCE_ACTION" },
  { label: "账单动作", value: "INVOICE_ACTION" },
  { label: "自定义任务", value: "CUSTOM" }
];

const channelOptions = [
  { label: "全部", value: "" },
  { label: "魔方云", value: "MOFANG_CLOUD" },
  { label: "上游财务", value: "ZJMF_API" },
  { label: "WHMCS", value: "WHMCS" },
  { label: "资源池", value: "RESOURCE" },
  { label: "人工交付", value: "MANUAL" },
  { label: "本地", value: "LOCAL" }
];

const stageOptions = [
  { label: "全部", value: "" },
  { label: "待处理", value: "PENDING" },
  { label: "执行中", value: "RUNNING" },
  { label: "同步", value: "SYNC" },
  { label: "财务", value: "FINANCE" },
  { label: "支付后开通", value: "AFTER_PAYMENT" },
  { label: "手工操作", value: "MANUAL" },
  { label: "回调", value: "CALLBACK" },
  { label: "完成", value: "DONE" },
  { label: "失败", value: "FAILED" },
  { label: "阻塞", value: "BLOCKED" }
];

const sourceTypeOptions = [
  { label: "全部", value: "" },
  { label: "服务", value: "service" },
  { label: "订单", value: "order" },
  { label: "账单", value: "invoice" },
  { label: "客户", value: "customer" },
  { label: "商品", value: "product" },
  { label: "计划任务", value: "cron" },
  { label: "资源", value: "resource" },
  { label: "渠道", value: "provider" }
];

const summary = computed(
  () => response.value?.summary ?? { total: 0, running: 0, success: 0, failed: 0, blocked: 0 }
);
const tasks = computed(() => response.value?.items ?? []);
const pendingCount = computed(() => tasks.value.filter(item => item.status === "PENDING").length);
const retryableCount = computed(() =>
  tasks.value.filter(item => ["FAILED", "BLOCKED", "PENDING"].includes(item.status)).length
);

const routeContext = computed(() => {
  const entries: Array<{ label: string; value: string }> = [];
  if (filters.serviceId) entries.push({ label: "服务", value: String(filters.serviceId) });
  if (filters.orderId) entries.push({ label: "订单", value: String(filters.orderId) });
  if (filters.invoiceId) entries.push({ label: "账单", value: String(filters.invoiceId) });
  if (filters.sourceType) entries.push({ label: "来源类型", value: sourceTypeLabel(filters.sourceType) });
  if (filters.sourceId) entries.push({ label: "来源 ID", value: String(filters.sourceId) });
  if (filters.channel) entries.push({ label: "渠道", value: channelLabel(filters.channel) });
  if (filters.status) entries.push({ label: "状态", value: statusLabel(filters.status) });
  if (filters.taskType) entries.push({ label: "任务类型", value: taskTypeLabel(filters.taskType) });
  return entries;
});

function statusType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "primary",
    RUNNING: "info",
    SUCCESS: "success",
    FAILED: "danger",
    BLOCKED: "warning"
  };
  return mapping[status] ?? "info";
}

function statusLabel(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "待执行",
    RUNNING: "执行中",
    SUCCESS: "成功",
    FAILED: "失败",
    BLOCKED: "阻塞"
  };
  return mapping[status] ?? status;
}

function taskTypeLabel(type: string) {
  const mapping: Record<string, string> = {
    AUTO_PROVISION: "自动开通",
    SERVICE_ACTION: "服务动作",
    PULL_SYNC_SERVICE: "服务同步",
    PULL_SYNC_BATCH: "批量同步",
    RESOURCE_ACTION: "资源动作",
    INVOICE_ACTION: "账单动作",
    CUSTOM: "自定义任务"
  };
  return mapping[type] ?? type;
}

function channelLabel(channel: string) {
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上游财务",
    WHMCS: "WHMCS",
    RESOURCE: "资源池",
    MANUAL: "人工交付",
    LOCAL: "本地"
  };
  return mapping[channel] ?? channel;
}

function sourceTypeLabel(type: string) {
  const mapping: Record<string, string> = {
    service: "服务",
    order: "订单",
    invoice: "账单",
    customer: "客户",
    product: "商品",
    cron: "计划任务",
    resource: "资源",
    provider: "渠道"
  };
  return mapping[type] ?? (type || "-");
}

function safeJson(value: string) {
  if (!value) return "-";
  try {
    return JSON.stringify(JSON.parse(value), null, 2);
  } catch {
    return value;
  }
}

function canRetry(task: AutomationTask) {
  return ["FAILED", "BLOCKED", "PENDING"].includes(task.status);
}

function readRouteQueryValue(value: unknown) {
  if (Array.isArray(value)) return String(value[0] ?? "");
  if (value === undefined || value === null) return "";
  return String(value);
}

function syncFiltersFromRoute() {
  const serviceId = readRouteQueryValue(route.query.serviceId);
  const orderId = readRouteQueryValue(route.query.orderId);
  const invoiceId = readRouteQueryValue(route.query.invoiceId);
  const sourceId = readRouteQueryValue(route.query.sourceId);

  filters.serviceId = serviceId ? Number(serviceId) : undefined;
  filters.orderId = orderId ? Number(orderId) : undefined;
  filters.invoiceId = invoiceId ? Number(invoiceId) : undefined;
  filters.sourceId = sourceId ? Number(sourceId) : undefined;
  filters.sourceType = readRouteQueryValue(route.query.sourceType);
  filters.stage = readRouteQueryValue(route.query.stage).toUpperCase();
  filters.channel = readRouteQueryValue(route.query.channel).toUpperCase();
  filters.status = readRouteQueryValue(route.query.status).toUpperCase();
  filters.taskType = readRouteQueryValue(route.query.taskType).toUpperCase();
  filters.keyword = readRouteQueryValue(route.query.keyword);
}

async function loadTasks() {
  loading.value = true;
  try {
    response.value = await fetchAutomationTasks({
      status: filters.status || undefined,
      taskType: filters.taskType || undefined,
      channel: filters.channel || undefined,
      stage: filters.stage || undefined,
      sourceType: filters.sourceType || undefined,
      sourceId: filters.sourceId || undefined,
      orderId: filters.orderId || undefined,
      invoiceId: filters.invoiceId || undefined,
      serviceId: filters.serviceId || undefined,
      keyword: filters.keyword || undefined,
      limit: 200
    });
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  filters.status = "";
  filters.taskType = "";
  filters.channel = "";
  filters.stage = "";
  filters.sourceType = "";
  filters.sourceId = undefined;
  filters.orderId = undefined;
  filters.invoiceId = undefined;
  filters.serviceId = undefined;
  filters.keyword = "";
  void router.push({ path: "/providers/automation", query: {} });
}

async function openDetail(task: AutomationTask) {
  detailLoading.value = true;
  try {
    detail.value = await fetchAutomationTaskDetail(task.id);
    detailVisible.value = true;
  } catch {
    ElMessage.error("无法加载任务详情");
  } finally {
    detailLoading.value = false;
  }
}

async function retryTask(task: AutomationTask) {
  retryingTaskId.value = task.id;
  try {
    const result = await retryAutomationTask(task.id);
    ElMessage.success(
      result.triggeredTask
        ? `${result.message}，新任务 ${result.triggeredTask.taskNo} 已创建`
        : result.message
    );
    await loadTasks();
    if (detailVisible.value && detail.value?.id === task.id) {
      await openDetail(task);
    }
  } catch (error: any) {
    ElMessage.error(error?.message ?? "任务重试失败");
  } finally {
    retryingTaskId.value = 0;
  }
}

function jumpToService(serviceId?: number) {
  if (!serviceId) return;
  void router.push(`/services/detail/${serviceId}`);
}

function jumpToOrder(orderId?: number) {
  if (!orderId) return;
  void router.push(`/orders/detail/${orderId}`);
}

function jumpToInvoice(invoiceId?: number) {
  if (!invoiceId) return;
  void router.push(`/billing/invoices/${invoiceId}`);
}

function jumpToAccount(task?: Partial<AutomationTask> | null) {
  if (!task) return;
  const query: Record<string, string> = {};
  const providerType = task.providerType || task.channel;
  if (providerType) query.providerType = providerType;
  if (task.sourceType === "provider" && task.sourceId) query.accountId = String(task.sourceId);
  void router.push({ path: "/providers/accounts", query });
}

function jumpToResources(task?: Partial<AutomationTask> | null) {
  if (!task) return;
  const query: Record<string, string> = {};
  const providerType = task.providerType || task.channel;
  if (providerType) query.providerType = providerType;
  if (task.serviceId) query.serviceId = String(task.serviceId);
  if (task.sourceType === "provider" && task.sourceId) query.accountId = String(task.sourceId);
  void router.push({ path: "/providers/resources", query });
}

function jumpToSettings() {
  void router.push("/providers/automation-settings");
}

onMounted(() => {
  syncFiltersFromRoute();
  void loadTasks();
});

watch(
  () => route.fullPath,
  () => {
    syncFiltersFromRoute();
    void loadTasks();
  }
);
</script>

<template>
  <PageWorkbench
    eyebrow="接口与上游 / 自动化"
    title="自动化任务中心"
    subtitle="集中查看自动开通、同步、资源动作和财务动作任务，统一追踪请求、结果和重试。"
  >
    <template #actions>
      <el-button @click="jumpToSettings">自动化策略</el-button>
      <el-button @click="resetFilters">重置筛选</el-button>
      <el-button type="primary" @click="loadTasks">刷新任务</el-button>
    </template>

    <template #metrics>
      <div class="automation-metrics">
        <div class="automation-metric">
          <span>任务总数</span>
          <strong>{{ summary.total }}</strong>
        </div>
        <div class="automation-metric">
          <span>待执行</span>
          <strong>{{ pendingCount }}</strong>
        </div>
        <div class="automation-metric automation-metric--info">
          <span>执行中</span>
          <strong>{{ summary.running }}</strong>
        </div>
        <div class="automation-metric automation-metric--success">
          <span>成功</span>
          <strong>{{ summary.success }}</strong>
        </div>
        <div class="automation-metric automation-metric--danger">
          <span>失败</span>
          <strong>{{ summary.failed }}</strong>
        </div>
        <div class="automation-metric automation-metric--warning">
          <span>可重试</span>
          <strong>{{ retryableCount }}</strong>
        </div>
      </div>
    </template>

    <template #filters>
      <el-form inline :model="filters" class="automation-filter">
        <el-form-item label="状态">
          <el-select v-model="filters.status" placeholder="全部状态" style="width: 140px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="任务类型">
          <el-select v-model="filters.taskType" placeholder="全部类型" style="width: 160px">
            <el-option v-for="item in taskTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="渠道">
          <el-select v-model="filters.channel" placeholder="全部渠道" style="width: 160px">
            <el-option v-for="item in channelOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="阶段">
          <el-select v-model="filters.stage" clearable filterable placeholder="全部阶段" style="width: 160px">
            <el-option v-for="item in stageOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源类型">
          <el-select v-model="filters.sourceType" clearable filterable placeholder="全部来源" style="width: 160px">
            <el-option v-for="item in sourceTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="来源 ID">
          <el-input-number v-model="filters.sourceId" :min="0" :precision="0" controls-position="right" style="width: 150px" />
        </el-form-item>
        <el-form-item label="订单 ID">
          <el-input-number v-model="filters.orderId" :min="0" :precision="0" controls-position="right" style="width: 150px" />
        </el-form-item>
        <el-form-item label="账单 ID">
          <el-input-number v-model="filters.invoiceId" :min="0" :precision="0" controls-position="right" style="width: 150px" />
        </el-form-item>
        <el-form-item label="服务 ID">
          <el-input-number v-model="filters.serviceId" :min="0" :precision="0" controls-position="right" style="width: 150px" />
        </el-form-item>
        <el-form-item label="关键字">
          <el-input
            v-model="filters.keyword"
            placeholder="任务号、标题、服务号、客户名"
            clearable
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadTasks">查询</el-button>
        </el-form-item>
      </el-form>
    </template>

    <div class="page-card" v-if="routeContext.length">
      <div class="context-panel">
        <div>
          <strong>当前已带入业务上下文</strong>
          <div class="context-tags">
            <el-tag v-for="item in routeContext" :key="`${item.label}-${item.value}`" effect="light">
              {{ item.label }}：{{ item.value }}
            </el-tag>
          </div>
        </div>
        <div class="inline-actions">
          <el-button v-if="filters.serviceId" plain @click="jumpToService(filters.serviceId)">查看服务</el-button>
          <el-button v-if="filters.orderId" plain @click="jumpToOrder(filters.orderId)">查看订单</el-button>
          <el-button v-if="filters.invoiceId" plain @click="jumpToInvoice(filters.invoiceId)">查看账单</el-button>
          <el-button
            v-if="filters.sourceType === 'provider' && filters.sourceId"
            plain
            @click="jumpToAccount({ sourceType: 'provider', sourceId: filters.sourceId, channel: filters.channel })"
          >
            查看接口账户
          </el-button>
          <el-button type="primary" plain @click="resetFilters">退出上下文</el-button>
        </div>
      </div>
    </div>

    <div class="page-card" v-loading="loading">
      <div class="page-header">
        <div>
          <h3 class="page-header__title">任务列表</h3>
          <p class="page-header__desc">失败或阻塞的任务可以直接在这里重新提交，也可以沿着服务、账单和渠道继续追查。</p>
        </div>
      </div>

      <el-table :data="tasks" border stripe empty-text="当前没有自动化任务">
        <el-table-column prop="taskNo" label="任务编号" min-width="180" />
        <el-table-column label="任务类型" min-width="130">
          <template #default="{ row }">{{ taskTypeLabel(row.taskType) }}</template>
        </el-table-column>
        <el-table-column prop="title" label="标题" min-width="220" show-overflow-tooltip />
        <el-table-column label="渠道" min-width="120">
          <template #default="{ row }">{{ channelLabel(row.channel) }}</template>
        </el-table-column>
        <el-table-column label="来源" min-width="140">
          <template #default="{ row }">{{ sourceTypeLabel(row.sourceType) }}</template>
        </el-table-column>
        <el-table-column prop="stage" label="阶段" min-width="120" />
        <el-table-column label="状态" min-width="100">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="serviceNo" label="服务号" min-width="160" />
        <el-table-column prop="customerName" label="客户" min-width="140" />
        <el-table-column prop="operatorName" label="执行方" min-width="120" />
        <el-table-column prop="createdAt" label="创建时间" min-width="170" />
        <el-table-column prop="message" label="执行结果" min-width="260" show-overflow-tooltip />
        <el-table-column label="操作" min-width="360" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetail(row)">详情</el-button>
              <el-button v-if="row.serviceId" type="primary" link @click="jumpToService(row.serviceId)">服务</el-button>
              <el-button v-if="row.orderId" type="primary" link @click="jumpToOrder(row.orderId)">订单</el-button>
              <el-button v-if="row.invoiceId" type="primary" link @click="jumpToInvoice(row.invoiceId)">账单</el-button>
              <el-button v-if="row.sourceType === 'provider' || row.serviceId" type="primary" link @click="jumpToResources(row)">
                渠道资源
              </el-button>
              <el-button v-if="row.sourceType === 'provider'" type="primary" link @click="jumpToAccount(row)">
                接口账户
              </el-button>
              <el-button
                v-if="canRetry(row)"
                type="warning"
                link
                :loading="retryingTaskId === row.id"
                @click="retryTask(row)"
              >
                重试
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-drawer v-model="detailVisible" size="760px" :with-header="false">
      <div v-loading="detailLoading" class="drawer-panel">
        <div class="page-header">
          <div>
            <h3 class="page-header__title">任务详情</h3>
            <p class="page-header__desc">查看请求参数、执行结果和关联对象。</p>
          </div>
          <div class="inline-actions" v-if="detail">
            <el-button v-if="detail.serviceId" plain @click="jumpToService(detail.serviceId)">服务</el-button>
            <el-button v-if="detail.orderId" plain @click="jumpToOrder(detail.orderId)">订单</el-button>
            <el-button v-if="detail.invoiceId" plain @click="jumpToInvoice(detail.invoiceId)">账单</el-button>
            <el-button v-if="detail.sourceType === 'provider' || detail.serviceId" plain @click="jumpToResources(detail)">
              渠道资源
            </el-button>
            <el-button v-if="detail.sourceType === 'provider'" plain @click="jumpToAccount(detail)">接口账户</el-button>
            <el-button
              v-if="canRetry(detail)"
              type="warning"
              plain
              :loading="retryingTaskId === detail.id"
              @click="retryTask(detail)"
            >
              重新提交
            </el-button>
          </div>
        </div>

        <template v-if="detail">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="任务编号">{{ detail.taskNo }}</el-descriptions-item>
            <el-descriptions-item label="任务类型">{{ taskTypeLabel(detail.taskType) }}</el-descriptions-item>
            <el-descriptions-item label="渠道">{{ channelLabel(detail.channel) }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="statusType(detail.status)" effect="light">{{ statusLabel(detail.status) }}</el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="来源类型">{{ sourceTypeLabel(detail.sourceType) }}</el-descriptions-item>
            <el-descriptions-item label="来源 ID">{{ detail.sourceId || "-" }}</el-descriptions-item>
            <el-descriptions-item label="订单 ID">{{ detail.orderId || "-" }}</el-descriptions-item>
            <el-descriptions-item label="账单 ID">{{ detail.invoiceId || "-" }}</el-descriptions-item>
            <el-descriptions-item label="服务 ID">{{ detail.serviceId || "-" }}</el-descriptions-item>
            <el-descriptions-item label="服务号">{{ detail.serviceNo || "-" }}</el-descriptions-item>
            <el-descriptions-item label="资源渠道">{{ channelLabel(detail.providerType || detail.channel) }}</el-descriptions-item>
            <el-descriptions-item label="远端资源 ID">{{ detail.providerResourceId || "-" }}</el-descriptions-item>
            <el-descriptions-item label="执行动作">{{ detail.actionName || "-" }}</el-descriptions-item>
            <el-descriptions-item label="执行方">{{ detail.operatorName || "-" }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ detail.createdAt || "-" }}</el-descriptions-item>
            <el-descriptions-item label="完成时间">{{ detail.finishedAt || "-" }}</el-descriptions-item>
          </el-descriptions>

          <div class="portal-grid portal-grid--two" style="margin-top: 16px">
            <div class="panel-card">
              <div class="section-card__head">
                <strong>关联工作台</strong>
                <span class="section-card__meta">按任务上下文继续排查服务、账单和渠道侧问题。</span>
              </div>
              <div class="inline-actions inline-actions--stack">
                <el-button v-if="detail.serviceId" plain @click="jumpToService(detail.serviceId)">打开服务工作台</el-button>
                <el-button v-if="detail.orderId" plain @click="jumpToOrder(detail.orderId)">打开订单工作台</el-button>
                <el-button v-if="detail.invoiceId" plain @click="jumpToInvoice(detail.invoiceId)">打开账单工作台</el-button>
                <el-button v-if="detail.sourceType === 'provider'" plain @click="jumpToAccount(detail)">打开接口账户</el-button>
                <el-button v-if="detail.sourceType === 'provider' || detail.serviceId" type="primary" plain @click="jumpToResources(detail)">
                  打开渠道资源
                </el-button>
              </div>
            </div>
            <div class="panel-card">
              <div class="section-card__head">
                <strong>执行说明</strong>
              </div>
              <el-alert :title="detail.message || '该任务当前没有额外说明'" type="info" :closable="false" show-icon />
            </div>
          </div>

          <div class="portal-grid portal-grid--two" style="margin-top: 16px">
            <div class="panel-card">
              <div class="section-card__head">
                <strong>请求参数</strong>
              </div>
              <pre class="json-block">{{ safeJson(detail.requestPayload) }}</pre>
            </div>
            <div class="panel-card">
              <div class="section-card__head">
                <strong>执行结果</strong>
              </div>
              <pre class="json-block">{{ safeJson(detail.resultPayload) }}</pre>
            </div>
          </div>
        </template>
      </div>
    </el-drawer>
  </PageWorkbench>
</template>

<style scoped>
.automation-metrics {
  display: grid;
  grid-template-columns: repeat(6, minmax(0, 1fr));
  gap: 12px;
}

.automation-metric {
  padding: 16px 18px;
  border-radius: 18px;
  background: #f6f8fb;
  border: 1px solid #e5ebf5;
  display: grid;
  gap: 8px;
}

.automation-metric span {
  font-size: 13px;
  color: #60708a;
}

.automation-metric strong {
  font-size: 28px;
  line-height: 1;
  color: #1f2a37;
}

.automation-metric--info {
  background: #eef6ff;
  border-color: #cfe3ff;
}

.automation-metric--success {
  background: #edf9f1;
  border-color: #cfead9;
}

.automation-metric--danger {
  background: #fff1f1;
  border-color: #f4cccc;
}

.automation-metric--warning {
  background: #fff8e8;
  border-color: #f1ddaa;
}

.automation-filter {
  display: flex;
  flex-wrap: wrap;
}

.context-panel {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  align-items: flex-start;
  flex-wrap: wrap;
}

.context-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-top: 10px;
}

.inline-actions {
  display: flex;
  gap: 10px;
  align-items: center;
  flex-wrap: wrap;
}

.inline-actions--stack {
  align-items: flex-start;
}

.json-block {
  margin: 0;
  border-radius: 16px;
  padding: 14px;
  background: #0f172a;
  color: #d7e1f3;
  font-size: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-word;
  min-height: 220px;
}

@media (max-width: 1280px) {
  .automation-metrics {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .automation-metrics {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
