<script setup lang="ts">
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";

import { http } from "@/api/http";
import { formatDate, formatDateTime } from "@/utils/format";
import {
  getLabel,
  getStatusTagType,
  roleMap,
  ticketPriorityMap,
  ticketStatusMap,
} from "@/utils/maps";

const loading = ref(false);
const createVisible = ref(false);
const replyVisible = ref(false);
const tickets = ref<any[]>([]);
const customers = ref<any[]>([]);
const services = ref<any[]>([]);
const admins = ref<any[]>([]);
const currentTicket = ref<any | null>(null);

const createForm = reactive({
  customerId: "",
  serviceId: "",
  subject: "",
  priority: "NORMAL",
  summary: "",
});

const replyForm = reactive({
  content: "",
  status: "PROCESSING",
  assignedToId: "",
  isInternal: false,
});

const filteredServices = computed(() =>
  services.value.filter((service) =>
    !createForm.customerId || service.customerId === createForm.customerId,
  ),
);

watch(
  () => createForm.customerId,
  () => {
    if (
      createForm.serviceId &&
      !filteredServices.value.some((service) => service.id === createForm.serviceId)
    ) {
      createForm.serviceId = "";
    }
  },
);

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/tickets");
    tickets.value = data.data.tickets;
    customers.value = data.data.customers;
    services.value = data.data.services;
    admins.value = data.data.admins;

    if (!createForm.customerId && customers.value[0]) {
      createForm.customerId = customers.value[0].id;
    }
  } finally {
    loading.value = false;
  }
}

async function submitCreate() {
  await http.post("/tickets", createForm);
  ElMessage.success("工单创建成功");
  createVisible.value = false;
  createForm.subject = "";
  createForm.summary = "";
  createForm.serviceId = "";
  await loadData();
}

function openReplyDialog(ticket: any) {
  currentTicket.value = ticket;
  replyForm.content = "";
  replyForm.status = ticket.status;
  replyForm.assignedToId = ticket.assignedToId || "";
  replyForm.isInternal = false;
  replyVisible.value = true;
}

async function submitReply() {
  if (!currentTicket.value) {
    return;
  }

  await http.post(`/tickets/${currentTicket.value.id}/reply`, replyForm);
  ElMessage.success("工单回复已提交");
  replyVisible.value = false;
  await loadData();
}

function priorityTagType(priority: string) {
  if (priority === "URGENT") {
    return "danger";
  }

  if (priority === "HIGH") {
    return "warning";
  }

  return "info";
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">工单中心</h1>
        <div class="page-subtitle">
          财务跟进、服务异常、客户沟通和内部协作统一在一套工单流程中处理。
        </div>
      </div>
      <el-button type="primary" @click="createVisible = true">新建工单</el-button>
    </div>

    <el-card class="page-card">
      <el-table v-loading="loading" :data="tickets" stripe>
        <el-table-column label="工单信息" min-width="250">
          <template #default="{ row }">
            <div style="font-weight: 700">{{ row.subject }}</div>
            <div style="margin-top: 6px; color: var(--text-secondary); font-size: 13px">
              {{ row.ticketNo }} / {{ row.customer.name }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="关联服务" min-width="180">
          <template #default="{ row }">
            {{ row.service?.name || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="问题摘要" min-width="220">
          <template #default="{ row }">
            {{ row.summary || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="优先级" width="110">
          <template #default="{ row }">
            <el-tag :type="priorityTagType(row.priority)">
              {{ getLabel(ticketPriorityMap, row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="140">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getLabel(ticketStatusMap, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="负责人" width="160">
          <template #default="{ row }">
            {{ row.assignedTo?.name || "-" }}
            <div style="margin-top: 4px; color: var(--text-secondary); font-size: 12px">
              {{ getLabel(roleMap, row.assignedTo?.role) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="最后回复" width="170">
          <template #default="{ row }">
            {{ formatDate(row.lastReplyAt || row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" plain @click="openReplyDialog(row)">
              回复
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="createVisible" title="新建工单" width="620px">
      <el-form label-position="top">
        <el-form-item label="客户">
          <el-select v-model="createForm.customerId" style="width: 100%">
            <el-option
              v-for="item in customers"
              :key="item.id"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="关联服务">
          <el-select v-model="createForm.serviceId" clearable style="width: 100%">
            <el-option
              v-for="item in filteredServices"
              :key="item.id"
              :label="`${item.name} / ${item.customer.name}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="工单标题">
          <el-input v-model="createForm.subject" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="createForm.priority" style="width: 100%">
            <el-option
              v-for="(label, value) in ticketPriorityMap"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="问题说明">
          <el-input v-model="createForm.summary" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">取消</el-button>
        <el-button type="primary" @click="submitCreate">确认创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="replyVisible" title="回复工单" width="720px">
      <div v-if="currentTicket" class="ticket-head">
        <div class="ticket-title">{{ currentTicket.subject }}</div>
        <div class="ticket-meta">
          {{ currentTicket.ticketNo }} / {{ currentTicket.customer.name }} / {{ getLabel(ticketStatusMap, currentTicket.status) }}
        </div>
      </div>

      <el-form label-position="top">
        <el-form-item label="回复内容">
          <el-input v-model="replyForm.content" type="textarea" :rows="5" />
        </el-form-item>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="工单状态">
              <el-select v-model="replyForm.status" style="width: 100%">
                <el-option
                  v-for="(label, value) in ticketStatusMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="负责人">
              <el-select v-model="replyForm.assignedToId" clearable style="width: 100%">
                <el-option
                  v-for="item in admins"
                  :key="item.id"
                  :label="`${item.name} / ${getLabel(roleMap, item.role)}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="内部备注">
              <el-switch v-model="replyForm.isInternal" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>

      <div v-if="currentTicket?.replies?.length" class="reply-history">
        <div class="reply-history-title">回复记录</div>
        <div
          v-for="reply in currentTicket.replies"
          :key="reply.id"
          class="reply-item"
        >
          <div class="reply-top">
            <span>{{ reply.authorName }}</span>
            <span>{{ formatDateTime(reply.createdAt) }}</span>
          </div>
          <div class="reply-body">{{ reply.content }}</div>
          <el-tag v-if="reply.isInternal" size="small" type="warning">内部备注</el-tag>
        </div>
      </div>

      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReply">提交回复</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.ticket-head {
  margin-bottom: 12px;
}

.ticket-title {
  font-weight: 700;
}

.ticket-meta {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 13px;
}

.reply-history {
  margin-top: 8px;
  display: grid;
  gap: 10px;
}

.reply-history-title {
  font-weight: 700;
}

.reply-item {
  padding: 12px 14px;
  border: 1px solid var(--border-color);
  border-radius: 14px;
  background: #fafcff;
}

.reply-top {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--text-secondary);
  font-size: 12px;
}

.reply-body {
  margin-top: 8px;
  line-height: 1.7;
}
</style>
