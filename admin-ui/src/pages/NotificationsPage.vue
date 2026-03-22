<script setup lang="ts">
import { ElMessage } from "element-plus";
import { onMounted, reactive, ref } from "vue";

import { http } from "@/api/http";
import { formatDateTime } from "@/utils/format";
import {
  getLabel,
  getStatusTagType,
  notificationChannelMap,
  notificationPriorityMap,
  notificationStatusMap,
} from "@/utils/maps";

const loading = ref(false);
const creating = ref(false);
const processing = ref(false);
const dialogVisible = ref(false);

const summary = ref<any>({
  templateCount: 0,
  totalMessages: 0,
  pendingMessages: 0,
  sentMessages: 0,
  failedMessages: 0,
  pendingTasks: 0,
});

const templates = ref<any[]>([]);
const messages = ref<any[]>([]);
const tasks = ref<any[]>([]);

const form = reactive({
  customerId: "",
  channel: "SYSTEM",
  priority: "NORMAL",
  recipient: "",
  recipientName: "",
  subject: "",
  content: "",
  module: "notification",
  relatedType: "",
  relatedId: "",
});

function resetForm() {
  Object.assign(form, {
    customerId: "",
    channel: "SYSTEM",
    priority: "NORMAL",
    recipient: "",
    recipientName: "",
    subject: "",
    content: "",
    module: "notification",
    relatedType: "",
    relatedId: "",
  });
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/notifications");
    summary.value = data.data.summary;
    templates.value = data.data.templates;
    messages.value = data.data.messages;
    tasks.value = data.data.tasks;
  } finally {
    loading.value = false;
  }
}

async function submitForm() {
  creating.value = true;
  try {
    await http.post("/notifications", form);
    ElMessage.success("通知已加入发送队列");
    dialogVisible.value = false;
    resetForm();
    await loadData();
  } finally {
    creating.value = false;
  }
}

async function processQueue() {
  processing.value = true;
  try {
    const { data } = await http.post("/notifications/process", {
      limit: 20,
    });
    ElMessage.success(`已处理 ${data.data.processed} 条通知任务`);
    await loadData();
  } finally {
    processing.value = false;
  }
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">通知中心</h1>
        <div class="page-subtitle">
          统一管理站内通知、邮件、短信和 Webhook 队列，并查看模板、投递状态和异步任务执行结果。
        </div>
      </div>
      <div class="toolbar-actions">
        <el-button :loading="processing" @click="processQueue">执行队列</el-button>
        <el-button type="primary" @click="dialogVisible = true">创建通知</el-button>
      </div>
    </div>

    <div class="metric-grid">
      <div class="metric-tile">
        <div class="metric-label">模板数量</div>
        <div class="metric-value">{{ summary.templateCount }}</div>
        <div class="metric-hint">预置支付、订单、工单和服务通知模板。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">通知总数</div>
        <div class="metric-value">{{ summary.totalMessages }}</div>
        <div class="metric-hint">系统近期生成的所有通知消息。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">待投递</div>
        <div class="metric-value">{{ summary.pendingMessages }}</div>
        <div class="metric-hint">仍在队列中等待发送的通知。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">待执行任务</div>
        <div class="metric-value">{{ summary.pendingTasks }}</div>
        <div class="metric-hint">通知异步任务队列中的待执行任务数量。</div>
      </div>
    </div>

    <el-card class="page-card">
      <template #header>
        <span>通知模板</span>
      </template>
      <el-table v-loading="loading" :data="templates" stripe>
        <el-table-column prop="name" label="模板名称" min-width="180" />
        <el-table-column prop="code" label="模板编码" min-width="180" />
        <el-table-column label="渠道" width="140">
          <template #default="{ row }">
            {{ getLabel(notificationChannelMap, row.channel) }}
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="标题模板" min-width="220" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.isActive ? 'success' : 'info'">
              {{ row.isActive ? "启用" : "停用" }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="page-card">
      <template #header>
        <span>通知消息</span>
      </template>
      <el-table v-loading="loading" :data="messages" stripe>
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="接收对象" min-width="220">
          <template #default="{ row }">
            <div style="font-weight: 700">{{ row.recipientName || row.customer?.name || row.recipient }}</div>
            <div class="muted-line">{{ row.recipient }}</div>
          </template>
        </el-table-column>
        <el-table-column label="渠道 / 优先级" width="180">
          <template #default="{ row }">
            <div>{{ getLabel(notificationChannelMap, row.channel) }}</div>
            <div class="muted-line">{{ getLabel(notificationPriorityMap, row.priority) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="subject" label="标题" min-width="220" />
        <el-table-column prop="content" label="内容" min-width="260" show-overflow-tooltip />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getLabel(notificationStatusMap, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card class="page-card">
      <template #header>
        <span>队列任务</span>
      </template>
      <el-table v-loading="loading" :data="tasks" stripe>
        <el-table-column label="任务时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column prop="jobType" label="任务类型" width="160" />
        <el-table-column prop="module" label="模块" width="120" />
        <el-table-column label="目标通知" min-width="220">
          <template #default="{ row }">
            {{ row.notification?.recipient || "-" }}
          </template>
        </el-table-column>
        <el-table-column prop="result" label="执行结果" min-width="220" show-overflow-tooltip />
        <el-table-column prop="errorMessage" label="错误信息" min-width="180" show-overflow-tooltip />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ row.status }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="创建通知" width="620px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="通知渠道">
              <el-select v-model="form.channel" style="width: 100%">
                <el-option
                  v-for="(label, value) in notificationChannelMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="优先级">
              <el-select v-model="form.priority" style="width: 100%">
                <el-option
                  v-for="(label, value) in notificationPriorityMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="接收地址">
              <el-input v-model="form.recipient" placeholder="邮箱、手机号或 Webhook 地址" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="接收方名称">
              <el-input v-model="form.recipientName" placeholder="例如 StarGalaxy 财务" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="模块">
              <el-input v-model="form.module" placeholder="例如 payment / ticket / notice" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="标题">
              <el-input v-model="form.subject" placeholder="站内通知可留空" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="通知内容">
          <el-input v-model="form.content" type="textarea" :rows="5" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="submitForm">加入队列</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.toolbar-actions {
  display: flex;
  gap: 12px;
}

.muted-line {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
