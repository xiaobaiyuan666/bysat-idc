<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import AuditTrailTable from "@/components/workbench/AuditTrailTable.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  claimTicket,
  fetchAdmins,
  fetchTicketDepartments,
  fetchTicketDetail,
  fetchTicketPresets,
  replyTicket,
  updateTicketRecord,
  type TicketDepartment,
  type TicketDetailResponse,
  type TicketPresetReply
} from "@/api/admin";

const route = useRoute();
const router = useRouter();

const text = {
  eyebrow: "\u8d22\u52a1 / \u5de5\u5355",
  title: "\u5de5\u5355\u5de5\u4f5c\u53f0",
  subtitle:
    "\u8fd9\u91cc\u5904\u7406\u5de5\u5355\u57fa\u672c\u4fe1\u606f\u3001\u8d1f\u8d23\u4eba\u5206\u914d\u3001\u72b6\u6001\u8c03\u6574\u3001\u56de\u590d\u8bb0\u5f55\u548c\u5ba1\u8ba1\u8f68\u8ff9\u3002",
  refresh: "\u5237\u65b0",
  back: "\u8fd4\u56de\u5217\u8868",
  ticketInfo: "\u5de5\u5355\u4fe1\u606f",
  basicConfig: "\u57fa\u7840\u914d\u7f6e",
  originalContent: "\u9996\u6b21\u63d0\u5355\u5185\u5bb9",
  replyHistory: "\u6c9f\u901a\u8bb0\u5f55",
  replyAction: "\u56de\u590d\u5de5\u5355",
  auditLogs: "\u53d8\u66f4\u8bb0\u5f55",
  save: "\u4fdd\u5b58\u8c03\u6574",
  claim: "\u63a5\u5355",
  submitReply: "\u53d1\u9001\u56de\u590d",
  emptyReply: "\u6682\u65e0\u56de\u590d\u8bb0\u5f55",
  titleLabel: "\u6807\u9898",
  customer: "\u5ba2\u6237",
  service: "\u5173\u8054\u670d\u52a1",
  department: "\u90e8\u95e8",
  assignee: "\u8d1f\u8d23\u4eba",
  priority: "\u4f18\u5148\u7ea7",
  status: "\u72b6\u6001",
  sla: "SLA",
  autoClose: "\u81ea\u52a8\u5173\u5355",
  latestReply: "\u6700\u65b0\u6458\u8981",
  createdAt: "\u521b\u5efa\u65f6\u95f4",
  closedAt: "\u5173\u5355\u65f6\u95f4",
  source: "\u6765\u6e90",
  visibility: "\u5185\u90e8\u5907\u6ce8",
  quickReply: "\u5feb\u6377\u56de\u590d",
  quickReplyPlaceholder: "\u9009\u62e9\u9884\u8bbe\u56de\u590d\u6a21\u677f",
  replyStatus: "\u56de\u590d\u540e\u72b6\u6001",
  replyPlaceholder: "\u8f93\u5165\u7ed9\u5ba2\u6237\u7684\u56de\u590d\u6216\u5185\u90e8\u5907\u6ce8",
  updateSuccess: "\u5de5\u5355\u5df2\u66f4\u65b0",
  replySuccess: "\u5de5\u5355\u56de\u590d\u5df2\u53d1\u9001"
} as const;

const loading = ref(false);
const saving = ref(false);
const replying = ref(false);
const detail = ref<TicketDetailResponse | null>(null);
const presets = ref<TicketPresetReply[]>([]);
const departments = ref<TicketDepartment[]>([]);
const admins = ref<
  Array<{
    id: number;
    username: string;
    displayName: string;
    status: string;
    roles: string[];
  }>
>([]);

const editForm = reactive({
  title: "",
  status: "OPEN",
  priority: "NORMAL",
  departmentName: "",
  assignedAdminId: 0
});

const replyForm = reactive({
  content: "",
  status: "WAITING_CUSTOMER",
  isInternal: false
});

const selectedPresetKey = ref("");

const pageTitle = computed(() => detail.value?.ticket.ticketNo || text.title);
const availableAdmins = computed(() => {
  const department = departments.value.find(item => item.name === editForm.departmentName);
  const limitIds = department?.adminIds || [];
  if (!limitIds.length) return admins.value;
  return admins.value.filter(item => limitIds.includes(item.id));
});

function statusLabel(status: string) {
  const mapping: Record<string, string> = {
    OPEN: "\u5f85\u5904\u7406",
    PROCESSING: "\u5904\u7406\u4e2d",
    WAITING_CUSTOMER: "\u5f85\u5ba2\u6237\u56de\u590d",
    CLOSED: "\u5df2\u5173\u95ed"
  };
  return mapping[status] ?? status;
}

function statusType(status: string) {
  const mapping: Record<string, string> = {
    OPEN: "warning",
    PROCESSING: "primary",
    WAITING_CUSTOMER: "info",
    CLOSED: "success"
  };
  return mapping[status] ?? "info";
}

function priorityLabel(priority: string) {
  const mapping: Record<string, string> = {
    LOW: "\u4f4e",
    NORMAL: "\u666e\u901a",
    HIGH: "\u9ad8",
    URGENT: "\u7d27\u6025"
  };
  return mapping[priority] ?? priority;
}

function priorityType(priority: string) {
  const mapping: Record<string, string> = {
    LOW: "info",
    NORMAL: "info",
    HIGH: "warning",
    URGENT: "danger"
  };
  return mapping[priority] ?? "info";
}

function authorLabel(authorType: string) {
  return authorType === "ADMIN" ? "\u7ba1\u7406\u5458" : "\u5ba2\u6237";
}

function slaLabel() {
  if (!detail.value) return "-";
  if (detail.value.ticket.slaPaused) return "\u6682\u505c\u8ba1\u65f6";
  const mapping: Record<string, string> = {
    ON_TRACK: "\u6b63\u5e38",
    AT_RISK: "\u5373\u5c06\u8d85\u65f6",
    BREACHED: "\u5df2\u8d85\u65f6",
    MET: "\u5df2\u5b8c\u6210",
    UNKNOWN: "\u672a\u8ba1\u7b97"
  };
  return mapping[detail.value.ticket.slaStatus] || detail.value.ticket.slaStatus || "-";
}

function slaType() {
  if (!detail.value) return "info";
  if (detail.value.ticket.slaPaused) return "info";
  const mapping: Record<string, string> = {
    ON_TRACK: "success",
    AT_RISK: "warning",
    BREACHED: "danger",
    MET: "success",
    UNKNOWN: "info"
  };
  return mapping[detail.value.ticket.slaStatus] ?? "info";
}

function assigneeLabel(id: number) {
  if (!id) return "-";
  const matched = admins.value.find(item => item.id === id);
  return matched?.displayName || matched?.username || `Admin #${id}`;
}

function fillForm() {
  if (!detail.value) return;
  editForm.title = detail.value.ticket.title || "";
  editForm.status = detail.value.ticket.status || "OPEN";
  editForm.priority = detail.value.ticket.priority || "NORMAL";
  editForm.departmentName = detail.value.ticket.departmentName || "";
  editForm.assignedAdminId = detail.value.ticket.assignedAdminId || 0;
}

async function loadData() {
  loading.value = true;
  try {
    const [ticketDetail, adminItems, departmentItems] = await Promise.all([
      fetchTicketDetail(route.params.id as string),
      fetchAdmins(),
      fetchTicketDepartments()
    ]);
    detail.value = ticketDetail;
    admins.value = adminItems;
    departments.value = departmentItems.items;
    const presetResult = await fetchTicketPresets(detail.value.ticket.departmentName || "");
    presets.value = presetResult.items;
    selectedPresetKey.value = "";
    fillForm();
  } finally {
    loading.value = false;
  }
}

async function handleClaim() {
  if (!detail.value) return;
  await claimTicket(detail.value.ticket.id);
  ElMessage.success(text.claim);
  await loadData();
}

async function submitEdit() {
  if (!detail.value) return;
  saving.value = true;
  try {
    await updateTicketRecord(detail.value.ticket.id, {
      title: editForm.title,
      status: editForm.status,
      priority: editForm.priority,
      departmentName: editForm.departmentName,
      assignedAdminId: editForm.assignedAdminId,
      assignedAdminName: editForm.assignedAdminId ? assigneeLabel(editForm.assignedAdminId) : ""
    });
    ElMessage.success(text.updateSuccess);
    await loadData();
  } finally {
    saving.value = false;
  }
}

async function submitReply() {
  if (!detail.value || !replyForm.content.trim()) return;
  replying.value = true;
  try {
    detail.value = await replyTicket(detail.value.ticket.id, {
      content: replyForm.content.trim(),
      status: replyForm.status,
      isInternal: replyForm.isInternal
    });
    replyForm.content = "";
    replyForm.isInternal = false;
    fillForm();
    ElMessage.success(text.replySuccess);
  } finally {
    replying.value = false;
  }
}

function applyPreset(key: string) {
  const matched = presets.value.find(item => item.key === key);
  if (!matched) return;
  replyForm.content = matched.content;
  replyForm.status = matched.status || replyForm.status;
}

onMounted(() => {
  void loadData();
});

watch(
  () => editForm.departmentName,
  () => {
    if (editForm.assignedAdminId && !availableAdmins.value.some(item => item.id === editForm.assignedAdminId)) {
      editForm.assignedAdminId = 0;
    }
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench :eyebrow="text.eyebrow" :title="pageTitle" :subtitle="text.subtitle">
      <template #actions>
        <el-button @click="loadData">{{ text.refresh }}</el-button>
        <el-button
          v-if="detail && !detail.ticket.assignedAdminId && detail.ticket.status !== 'CLOSED'"
          type="warning"
          plain
          @click="handleClaim"
        >
          {{ text.claim }}
        </el-button>
        <el-button plain @click="router.push('/tickets/list')">{{ text.back }}</el-button>
      </template>

      <template #metrics>
        <div v-if="detail" class="summary-strip">
          <div class="summary-pill">
            <span>{{ text.status }}</span>
            <strong>{{ statusLabel(detail.ticket.status) }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.priority }}</span>
            <strong>{{ priorityLabel(detail.ticket.priority) }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.customer }}</span>
            <strong>{{ detail.ticket.customerName }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.service }}</span>
            <strong>{{ detail.ticket.serviceNo || "-" }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.sla }}</span>
            <strong>{{ slaLabel() }}</strong>
          </div>
        </div>
      </template>

      <template v-if="detail">
        <div class="ticket-grid">
          <div class="panel-card">
            <div class="section-card__head">
              <strong>{{ text.ticketInfo }}</strong>
            </div>
            <div class="ticket-info-list">
              <div><span>{{ text.customer }}</span><strong>{{ detail.ticket.customerName }} / ID {{ detail.ticket.customerId }}</strong></div>
              <div><span>{{ text.service }}</span><strong>{{ detail.ticket.serviceNo || "-" }} {{ detail.ticket.productName || "" }}</strong></div>
              <div><span>{{ text.department }}</span><strong>{{ detail.ticket.departmentName || "-" }}</strong></div>
              <div><span>{{ text.assignee }}</span><strong>{{ detail.ticket.assignedAdminName || "-" }}</strong></div>
              <div><span>{{ text.sla }}</span><strong>{{ slaLabel() }} {{ detail.ticket.slaDeadlineAt || "" }}</strong></div>
              <div><span>{{ text.autoClose }}</span><strong>{{ detail.ticket.autoCloseAt || "-" }}</strong></div>
              <div><span>{{ text.source }}</span><strong>{{ detail.ticket.source || "-" }}</strong></div>
              <div><span>{{ text.createdAt }}</span><strong>{{ detail.ticket.createdAt || "-" }}</strong></div>
              <div><span>{{ text.closedAt }}</span><strong>{{ detail.ticket.closedAt || "-" }}</strong></div>
              <div><span>{{ text.latestReply }}</span><strong>{{ detail.ticket.latestReplyExcerpt || "-" }}</strong></div>
            </div>
          </div>

          <div class="panel-card">
            <div class="section-card__head">
              <strong>{{ text.basicConfig }}</strong>
            </div>
            <el-form label-position="top">
              <el-form-item :label="text.titleLabel">
                <el-input v-model="editForm.title" />
              </el-form-item>
              <div class="ticket-form-grid">
                <el-form-item :label="text.status">
                  <el-select v-model="editForm.status">
                    <el-option :label="statusLabel('OPEN')" value="OPEN" />
                    <el-option :label="statusLabel('PROCESSING')" value="PROCESSING" />
                    <el-option :label="statusLabel('WAITING_CUSTOMER')" value="WAITING_CUSTOMER" />
                    <el-option :label="statusLabel('CLOSED')" value="CLOSED" />
                  </el-select>
                </el-form-item>
                <el-form-item :label="text.priority">
                  <el-select v-model="editForm.priority">
                    <el-option :label="priorityLabel('LOW')" value="LOW" />
                    <el-option :label="priorityLabel('NORMAL')" value="NORMAL" />
                    <el-option :label="priorityLabel('HIGH')" value="HIGH" />
                    <el-option :label="priorityLabel('URGENT')" value="URGENT" />
                  </el-select>
                </el-form-item>
              </div>
              <div class="ticket-form-grid">
                <el-form-item :label="text.department">
                  <el-select v-model="editForm.departmentName" filterable allow-create default-first-option>
                    <el-option
                      v-for="item in departments.filter(item => item.enabled)"
                      :key="item.key"
                      :label="item.name"
                      :value="item.name"
                    />
                  </el-select>
                </el-form-item>
                <el-form-item :label="text.assignee">
                  <el-select v-model="editForm.assignedAdminId" clearable>
                    <el-option :label="'-'" :value="0" />
                    <el-option
                      v-for="item in availableAdmins"
                      :key="item.id"
                      :label="item.displayName"
                      :value="item.id"
                    />
                  </el-select>
                </el-form-item>
              </div>
              <div class="action-group">
                <el-button type="primary" :loading="saving" @click="submitEdit">{{ text.save }}</el-button>
              </div>
            </el-form>
          </div>
        </div>

        <div class="ticket-grid">
          <div class="panel-card">
            <div class="section-card__head">
              <strong>{{ text.originalContent }}</strong>
              <div class="reply-item__title">
                <el-tag :type="statusType(detail.ticket.status)" effect="light">{{ statusLabel(detail.ticket.status) }}</el-tag>
                <el-tag :type="slaType()" effect="light">{{ slaLabel() }}</el-tag>
              </div>
            </div>
            <div class="ticket-content">{{ detail.ticket.content || "-" }}</div>
          </div>

          <div class="panel-card">
            <div class="section-card__head">
              <strong>{{ text.replyAction }}</strong>
            </div>
            <el-form label-position="top">
              <div class="ticket-form-grid">
                <el-form-item :label="text.replyStatus">
                  <el-select v-model="replyForm.status">
                    <el-option :label="statusLabel('PROCESSING')" value="PROCESSING" />
                    <el-option :label="statusLabel('WAITING_CUSTOMER')" value="WAITING_CUSTOMER" />
                    <el-option :label="statusLabel('CLOSED')" value="CLOSED" />
                  </el-select>
                </el-form-item>
                <el-form-item :label="text.visibility">
                  <el-switch v-model="replyForm.isInternal" />
                </el-form-item>
              </div>
              <el-form-item :label="text.quickReply">
                <el-select
                  v-model="selectedPresetKey"
                  clearable
                  style="width: 100%"
                  :placeholder="text.quickReplyPlaceholder"
                  @change="applyPreset"
                >
                  <el-option v-for="item in presets" :key="item.key" :label="item.title" :value="item.key" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-input v-model="replyForm.content" type="textarea" :rows="5" :placeholder="text.replyPlaceholder" />
              </el-form-item>
              <div class="action-group">
                <el-button type="primary" :loading="replying" @click="submitReply">{{ text.submitReply }}</el-button>
              </div>
            </el-form>
          </div>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>{{ text.replyHistory }}</strong>
          </div>
          <div v-if="detail.replies.length" class="reply-list">
            <div v-for="item in detail.replies" :key="item.id" class="reply-item">
              <div class="reply-item__meta">
                <div class="reply-item__title">
                  <el-tag size="small" :type="item.authorType === 'ADMIN' ? 'primary' : 'success'" effect="light">
                    {{ authorLabel(item.authorType) }}
                  </el-tag>
                  <strong>{{ item.authorName }}</strong>
                  <el-tag v-if="item.isInternal" size="small" type="warning" effect="light">{{ text.visibility }}</el-tag>
                </div>
                <span>{{ item.createdAt }}</span>
              </div>
              <div class="reply-item__content">{{ item.content }}</div>
            </div>
          </div>
          <el-empty v-else :description="text.emptyReply" />
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>{{ text.auditLogs }}</strong>
          </div>
          <AuditTrailTable :items="detail.auditLogs" :empty-text="'\u6682\u65e0\u53d8\u66f4\u8bb0\u5f55'" />
        </div>
      </template>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.ticket-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
  margin-bottom: 16px;
}

.ticket-info-list {
  display: grid;
  gap: 12px;
}

.ticket-info-list div {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  font-size: 13px;
}

.ticket-info-list span {
  color: var(--el-text-color-secondary);
}

.ticket-info-list strong {
  text-align: right;
}

.ticket-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.ticket-content {
  white-space: pre-wrap;
  line-height: 1.7;
  color: var(--el-text-color-regular);
}

.reply-list {
  display: grid;
  gap: 12px;
}

.reply-item {
  border: 1px solid var(--el-border-color-light);
  border-radius: 14px;
  padding: 14px 16px;
  background: var(--el-fill-color-blank);
}

.reply-item__meta {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 8px;
  color: var(--el-text-color-secondary);
  font-size: 12px;
}

.reply-item__title {
  display: flex;
  align-items: center;
  gap: 8px;
}

.reply-item__content {
  white-space: pre-wrap;
  line-height: 1.7;
}

@media (max-width: 1200px) {
  .ticket-grid {
    grid-template-columns: 1fr;
  }
}
</style>
