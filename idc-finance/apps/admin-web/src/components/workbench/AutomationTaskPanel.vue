<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRouter } from "vue-router";
import { fetchAutomationTasks, retryAutomationTask, type AutomationTask } from "@/api/admin";
import { useLocaleStore } from "@/store";
import { formatAutomationTaskStatus, formatAutomationTaskTitle, formatAutomationTaskType } from "@/utils/business";

const props = withDefaults(
  defineProps<{
    title?: string;
    serviceId?: number;
    orderId?: number;
    invoiceId?: number;
    sourceType?: string;
    sourceId?: number;
    limit?: number;
  }>(),
  {
    title: "自动化任务",
    serviceId: 0,
    orderId: 0,
    invoiceId: 0,
    sourceType: "",
    sourceId: 0,
    limit: 6
  }
);

const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(false);
const retryingTaskId = ref(0);
const tasks = ref<AutomationTask[]>([]);

const copy = computed(() => ({
  refresh: localeStore.locale === "en-US" ? "Refresh" : "刷新",
  openTaskCenter: localeStore.locale === "en-US" ? "Open Task Center" : "打开任务中心",
  total: localeStore.locale === "en-US" ? "Total" : "任务总数",
  running: localeStore.locale === "en-US" ? "Running" : "执行中",
  success: localeStore.locale === "en-US" ? "Success" : "成功",
  failed: localeStore.locale === "en-US" ? "Failed" : "失败",
  blocked: localeStore.locale === "en-US" ? "Blocked" : "阻塞",
  taskNo: localeStore.locale === "en-US" ? "Task No." : "任务编号",
  type: localeStore.locale === "en-US" ? "Type" : "类型",
  title: localeStore.locale === "en-US" ? "Title" : "标题",
  status: localeStore.locale === "en-US" ? "Status" : "状态",
  operator: localeStore.locale === "en-US" ? "Operator" : "执行方",
  createdAt: localeStore.locale === "en-US" ? "Created At" : "创建时间",
  message: localeStore.locale === "en-US" ? "Message" : "结果消息",
  actions: localeStore.locale === "en-US" ? "Actions" : "操作",
  retry: localeStore.locale === "en-US" ? "Retry" : "重试",
  empty: localeStore.locale === "en-US" ? "No automation tasks for this object." : "当前对象暂无自动化任务",
  retrySuccess: localeStore.locale === "en-US" ? "Retry has been queued" : "任务重试已重新入队",
  retryFailed: localeStore.locale === "en-US" ? "Task retry failed" : "任务重试失败"
}));

const summary = computed(() => {
  const result = { total: tasks.value.length, running: 0, success: 0, failed: 0, blocked: 0 };
  for (const item of tasks.value) {
    if (item.status === "RUNNING") result.running += 1;
    if (item.status === "SUCCESS") result.success += 1;
    if (item.status === "FAILED") result.failed += 1;
    if (item.status === "BLOCKED") result.blocked += 1;
  }
  return result;
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

function canRetry(task: AutomationTask) {
  return ["FAILED", "BLOCKED", "PENDING"].includes(task.status);
}

async function loadTasks() {
  loading.value = true;
  try {
    const response = await fetchAutomationTasks({
      serviceId: props.serviceId || undefined,
      orderId: props.orderId || undefined,
      invoiceId: props.invoiceId || undefined,
      sourceType: props.sourceType || undefined,
      sourceId: props.sourceId || undefined,
      limit: props.limit
    });
    tasks.value = response.items;
  } finally {
    loading.value = false;
  }
}

async function retryTask(task: AutomationTask) {
  retryingTaskId.value = task.id;
  try {
    const result = await retryAutomationTask(task.id);
    ElMessage.success(
      result.triggeredTask
        ? `${copy.value.retrySuccess}：${result.triggeredTask.taskNo}`
        : result.message || copy.value.retrySuccess
    );
    await loadTasks();
  } catch (error: any) {
    ElMessage.error(error?.message ?? copy.value.retryFailed);
  } finally {
    retryingTaskId.value = 0;
  }
}

function openTaskCenter() {
  const query: Record<string, string> = {};
  if (props.serviceId) query.serviceId = String(props.serviceId);
  if (props.orderId) query.orderId = String(props.orderId);
  if (props.invoiceId) query.invoiceId = String(props.invoiceId);
  if (props.sourceType) query.sourceType = props.sourceType;
  if (props.sourceId) query.sourceId = String(props.sourceId);
  router.push({ path: "/providers/automation", query });
}

watch(
  () => [props.serviceId, props.orderId, props.invoiceId, props.sourceType, props.sourceId].join("|"),
  () => void loadTasks()
);

onMounted(() => void loadTasks());
</script>

<template>
  <div class="panel-card" v-loading="loading">
    <div class="section-card__head">
      <strong>{{ title }}</strong>
      <div class="inline-actions">
        <el-button size="small" @click="loadTasks">{{ copy.refresh }}</el-button>
        <el-button size="small" type="primary" plain @click="openTaskCenter">{{ copy.openTaskCenter }}</el-button>
      </div>
    </div>

    <div class="summary-strip" style="margin-bottom: 14px">
      <div class="summary-pill">
        <span>{{ copy.total }}</span>
        <strong>{{ summary.total }}</strong>
      </div>
      <div class="summary-pill">
        <span>{{ copy.running }}</span>
        <strong>{{ summary.running }}</strong>
      </div>
      <div class="summary-pill">
        <span>{{ copy.success }}</span>
        <strong>{{ summary.success }}</strong>
      </div>
      <div class="summary-pill">
        <span>{{ copy.failed }}</span>
        <strong>{{ summary.failed }}</strong>
      </div>
      <div class="summary-pill">
        <span>{{ copy.blocked }}</span>
        <strong>{{ summary.blocked }}</strong>
      </div>
    </div>

    <el-table :data="tasks" border stripe :empty-text="copy.empty">
      <el-table-column prop="taskNo" :label="copy.taskNo" min-width="170" />
      <el-table-column :label="copy.type" min-width="120">
        <template #default="{ row }">{{ formatAutomationTaskType(localeStore.locale, row.taskType) }}</template>
      </el-table-column>
      <el-table-column :label="copy.title" min-width="220" show-overflow-tooltip>
        <template #default="{ row }">
          {{ formatAutomationTaskTitle(localeStore.locale, row.taskType, row.title, row.actionName) }}
        </template>
      </el-table-column>
      <el-table-column :label="copy.status" min-width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" effect="light">{{ formatAutomationTaskStatus(localeStore.locale, row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="operatorName" :label="copy.operator" min-width="120" />
      <el-table-column prop="createdAt" :label="copy.createdAt" min-width="170" />
      <el-table-column prop="message" :label="copy.message" min-width="240" show-overflow-tooltip />
      <el-table-column :label="copy.actions" min-width="140" fixed="right">
        <template #default="{ row }">
          <div class="inline-actions">
            <el-button
              v-if="canRetry(row)"
              size="small"
              type="warning"
              link
              :loading="retryingTaskId === row.id"
              @click="retryTask(row)"
            >
              {{ copy.retry }}
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
