<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import {
  closeTicket,
  createTicket,
  loadTicketDetail,
  loadTicketDepartments,
  loadTickets,
  loadServices,
  replyTicket,
  type PortalTicketDepartment,
  type PortalService,
  type PortalTicketDetailResponse,
  type PortalTicketRecord
} from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import { formatPortalServiceStatus } from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();

const loading = ref(true);
const submitLoading = ref(false);
const detailLoading = ref(false);
const createVisible = ref(false);
const detailVisible = ref(false);
const tickets = ref<PortalTicketRecord[]>([]);
const services = ref<PortalService[]>([]);
const departments = ref<PortalTicketDepartment[]>([]);
const detail = ref<PortalTicketDetailResponse | null>(null);

const filters = reactive({
  keyword: "",
  status: "",
  priority: "",
  serviceId: undefined as number | undefined
});

const createForm = reactive({
  serviceId: undefined as number | undefined,
  title: "",
  content: "",
  priority: "NORMAL",
  departmentName: "Technical Support"
});

const replyForm = reactive({
  content: "",
  status: "OPEN"
});

const serviceMap = computed(() => {
  const mapping = new Map<number, PortalService>();
  for (const item of services.value) {
    mapping.set(item.id, item);
  }
  return mapping;
});

const filteredTickets = computed(() =>
  tickets.value.filter(item => {
    const keyword = filters.keyword.trim().toLowerCase();
    const matchesKeyword =
      !keyword ||
      item.ticketNo.toLowerCase().includes(keyword) ||
      item.title.toLowerCase().includes(keyword) ||
      item.serviceNo.toLowerCase().includes(keyword) ||
      item.productName.toLowerCase().includes(keyword);
    const matchesStatus = !filters.status || item.status === filters.status;
    const matchesPriority = !filters.priority || item.priority === filters.priority;
    const matchesService = !filters.serviceId || item.serviceId === filters.serviceId;
    return matchesKeyword && matchesStatus && matchesPriority && matchesService;
  })
);

const openCount = computed(() => tickets.value.filter(item => item.status === "OPEN").length);
const processingCount = computed(() => tickets.value.filter(item => item.status === "PROCESSING").length);
const waitingCount = computed(() => tickets.value.filter(item => item.status === "WAITING_CUSTOMER").length);
const currentLinkedService = computed(() =>
  detail.value?.ticket.serviceId ? serviceMap.value.get(detail.value.ticket.serviceId) ?? null : null
);

const serviceOptions = computed(() =>
  services.value.map(item => ({
    label: `${item.serviceNo} / ${item.productName}`,
    value: item.id
  }))
);

const departmentOptions = computed(() => departments.value.filter(item => item.enabled));

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "工单中心", "Tickets"),
  title: pickLabel(localeStore.locale, "工单中心", "Tickets"),
  subtitle: pickLabel(
    localeStore.locale,
    "围绕服务、账单和沟通记录统一查看问题进度，减少客户在多个页面之间重复跳转。",
    "Track support progress with linked services, billing, and conversation records."
  ),
  create: pickLabel(localeStore.locale, "提交新工单", "New Ticket"),
  refresh: pickLabel(localeStore.locale, "刷新", "Refresh"),
  total: pickLabel(localeStore.locale, "工单总数", "Total"),
  open: pickLabel(localeStore.locale, "待处理", "Open"),
  processing: pickLabel(localeStore.locale, "处理中", "Processing"),
  waiting: pickLabel(localeStore.locale, "待您回复", "Waiting for You"),
  closed: pickLabel(localeStore.locale, "已关闭", "Closed"),
  listTitle: pickLabel(localeStore.locale, "工单列表", "Ticket List"),
  listDesc: pickLabel(
    localeStore.locale,
    "可按工单号、标题、服务、状态和优先级过滤，并从这里直接跳到服务或账单中心。",
    "Filter by ticket, title, service, status, and priority, then jump directly to the related service or invoice."
  ),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索工单号 / 标题 / 服务号", "Search ticket, title, or service"),
  statusPlaceholder: pickLabel(localeStore.locale, "状态", "Status"),
  priorityPlaceholder: pickLabel(localeStore.locale, "优先级", "Priority"),
  servicePlaceholder: pickLabel(localeStore.locale, "关联服务", "Linked Service"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  ticketNo: pickLabel(localeStore.locale, "工单号", "Ticket No."),
  titleCol: pickLabel(localeStore.locale, "标题", "Title"),
  service: pickLabel(localeStore.locale, "关联服务", "Service"),
  business: pickLabel(localeStore.locale, "关联业务", "Business"),
  product: pickLabel(localeStore.locale, "商品", "Product"),
  priority: pickLabel(localeStore.locale, "优先级", "Priority"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  sla: "SLA",
  autoClose: pickLabel(localeStore.locale, "自动关单", "Auto Close"),
  updatedAt: pickLabel(localeStore.locale, "更新时间", "Updated"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  detail: pickLabel(localeStore.locale, "查看详情", "View"),
  serviceDetail: pickLabel(localeStore.locale, "查看服务", "View Service"),
  invoiceDetail: pickLabel(localeStore.locale, "查看账单", "View Invoice"),
  empty: pickLabel(localeStore.locale, "暂无匹配的工单记录。", "No matching tickets."),
  createTitle: pickLabel(localeStore.locale, "提交工单", "Create Ticket"),
  linkedService: pickLabel(localeStore.locale, "关联服务", "Linked Service"),
  linkedServiceHint: pickLabel(localeStore.locale, "不选则为通用工单", "Leave empty for a general ticket"),
  createTicketTitle: pickLabel(localeStore.locale, "工单标题", "Ticket Title"),
  content: pickLabel(localeStore.locale, "问题描述", "Description"),
  department: pickLabel(localeStore.locale, "部门", "Department"),
  submit: pickLabel(localeStore.locale, "提交", "Submit"),
  cancel: pickLabel(localeStore.locale, "取消", "Cancel"),
  detailTitle: pickLabel(localeStore.locale, "工单详情", "Ticket Detail"),
  origin: pickLabel(localeStore.locale, "首次提单", "Original Request"),
  replyTitle: pickLabel(localeStore.locale, "回复工单", "Reply"),
  replyPlaceholder: pickLabel(localeStore.locale, "输入回复内容", "Type your reply"),
  submitReply: pickLabel(localeStore.locale, "发送回复", "Send Reply"),
  closeTicket: pickLabel(localeStore.locale, "关闭工单", "Close Ticket"),
  replyStatus: pickLabel(localeStore.locale, "回复后状态", "Reply Status"),
  conversation: pickLabel(localeStore.locale, "沟通记录", "Conversation"),
  noReplies: pickLabel(localeStore.locale, "暂无回复记录", "No replies yet"),
  assignedAdmin: pickLabel(localeStore.locale, "对接客服", "Assigned Agent"),
  createdAt: pickLabel(localeStore.locale, "提交时间", "Created"),
  relatedBusiness: pickLabel(localeStore.locale, "关联业务", "Related Business"),
  invoiceCenter: pickLabel(localeStore.locale, "账单中心", "Invoice Center"),
  createFromService: pickLabel(localeStore.locale, "继续围绕该服务提单", "Create Another Ticket"),
  serviceStatus: pickLabel(localeStore.locale, "服务状态", "Service Status"),
  providerType: pickLabel(localeStore.locale, "渠道", "Provider"),
  successCreate: pickLabel(localeStore.locale, "工单提交成功", "Ticket submitted"),
  successReply: pickLabel(localeStore.locale, "已发送回复", "Reply sent"),
  successClose: pickLabel(localeStore.locale, "工单已关闭", "Ticket closed")
}));

function ticketStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    OPEN: copy.value.open,
    PROCESSING: copy.value.processing,
    WAITING_CUSTOMER: copy.value.waiting,
    CLOSED: copy.value.closed
  };
  return mapping[status] ?? status;
}

function ticketStatusType(status: string) {
  const mapping: Record<string, string> = {
    OPEN: "warning",
    PROCESSING: "info",
    WAITING_CUSTOMER: "primary",
    CLOSED: "success"
  };
  return mapping[status] ?? "info";
}

function ticketPriorityLabel(priority: string) {
  const mapping: Record<string, string> = {
    LOW: pickLabel(localeStore.locale, "低", "Low"),
    NORMAL: pickLabel(localeStore.locale, "普通", "Normal"),
    HIGH: pickLabel(localeStore.locale, "高", "High"),
    URGENT: pickLabel(localeStore.locale, "紧急", "Urgent")
  };
  return mapping[priority] ?? priority;
}

function ticketPriorityType(priority: string) {
  const mapping: Record<string, string> = {
    LOW: "info",
    NORMAL: "info",
    HIGH: "warning",
    URGENT: "danger"
  };
  return mapping[priority] ?? "info";
}

function replyAuthorLabel(authorType: string) {
  return authorType === "ADMIN"
    ? pickLabel(localeStore.locale, "客服", "Agent")
    : pickLabel(localeStore.locale, "我", "Me");
}

function ticketSlaLabel(ticket: PortalTicketRecord) {
  if (ticket.slaPaused) return pickLabel(localeStore.locale, "暂停计时", "Paused");
  const mapping: Record<string, string> = {
    ON_TRACK: pickLabel(localeStore.locale, "正常", "On Track"),
    AT_RISK: pickLabel(localeStore.locale, "即将超时", "At Risk"),
    BREACHED: pickLabel(localeStore.locale, "已超时", "Breached"),
    MET: pickLabel(localeStore.locale, "已完成", "Completed"),
    UNKNOWN: pickLabel(localeStore.locale, "未计算", "Unknown")
  };
  return mapping[ticket.slaStatus] || ticket.slaStatus || "-";
}

function ticketSlaType(ticket: PortalTicketRecord) {
  if (ticket.slaPaused) return "info";
  const mapping: Record<string, string> = {
    ON_TRACK: "success",
    AT_RISK: "warning",
    BREACHED: "danger",
    MET: "success",
    UNKNOWN: "info"
  };
  return mapping[ticket.slaStatus] ?? "info";
}

function getLinkedService(serviceId?: number) {
  if (!serviceId) return null;
  return serviceMap.value.get(serviceId) ?? null;
}

function goToService(serviceId?: number) {
  if (!serviceId) return;
  void router.push(`/services/${serviceId}`);
}

function goToInvoice(serviceId?: number) {
  if (!serviceId) return;
  void router.push({ path: "/invoices", query: { serviceId: String(serviceId) } });
}

async function fetchData() {
  loading.value = true;
  try {
    const [ticketItems, serviceItems, departmentItems] = await Promise.all([
      loadTickets(),
      loadServices(),
      loadTicketDepartments()
    ]);
    tickets.value = ticketItems;
    services.value = serviceItems;
    departments.value = departmentItems;
    if (departmentItems.length && !departmentItems.find(item => item.name === createForm.departmentName)) {
      const defaultDepartment = departmentItems.find(item => item.isDefault) || departmentItems[0];
      createForm.departmentName = defaultDepartment?.name || "Technical Support";
    }
  } finally {
    loading.value = false;
  }
}

function resetCreateForm() {
  createForm.serviceId = undefined;
  createForm.title = "";
  createForm.content = "";
  createForm.priority = "NORMAL";
  createForm.departmentName = "Technical Support";
}

async function openDetail(ticket: PortalTicketRecord) {
  detailVisible.value = true;
  detailLoading.value = true;
  try {
    detail.value = await loadTicketDetail(ticket.id);
    replyForm.content = "";
    replyForm.status =
      detail.value.ticket.status === "CLOSED" ? "OPEN" : detail.value.ticket.status || "OPEN";
  } finally {
    detailLoading.value = false;
  }
}

async function submitCreate() {
  if (!createForm.title.trim() || !createForm.content.trim()) return;
  submitLoading.value = true;
  try {
    const created = await createTicket({
      serviceId: createForm.serviceId,
      title: createForm.title.trim(),
      content: createForm.content.trim(),
      priority: createForm.priority,
      departmentName: createForm.departmentName
    });
    ElMessage.success(copy.value.successCreate);
    createVisible.value = false;
    resetCreateForm();
    await fetchData();
    await openDetail(created);
  } finally {
    submitLoading.value = false;
  }
}

async function submitReply() {
  if (!detail.value || !replyForm.content.trim()) return;
  submitLoading.value = true;
  try {
    detail.value = await replyTicket(detail.value.ticket.id, {
      content: replyForm.content.trim(),
      status: replyForm.status
    });
    replyForm.content = "";
    ElMessage.success(copy.value.successReply);
    await fetchData();
  } finally {
    submitLoading.value = false;
  }
}

async function handleCloseTicket() {
  if (!detail.value) return;
  submitLoading.value = true;
  try {
    await closeTicket(detail.value.ticket.id);
    ElMessage.success(copy.value.successClose);
    detail.value = await loadTicketDetail(detail.value.ticket.id);
    await fetchData();
  } finally {
    submitLoading.value = false;
  }
}

function openCreateDialog(serviceId?: number) {
  resetCreateForm();
  if (serviceId) {
    createForm.serviceId = serviceId;
  }
  createVisible.value = true;
}

function applyRouteState() {
  const keywordQuery = typeof route.query.keyword === "string" ? route.query.keyword : "";
  const statusQuery = typeof route.query.status === "string" ? route.query.status : "";
  const priorityQuery = typeof route.query.priority === "string" ? route.query.priority : "";
  const actionQuery = typeof route.query.action === "string" ? route.query.action : "";
  const titleQuery = typeof route.query.title === "string" ? route.query.title : "";
  const serviceIdQuery = typeof route.query.serviceId === "string" ? Number(route.query.serviceId) : NaN;

  filters.keyword = keywordQuery;
  filters.status = statusQuery;
  filters.priority = priorityQuery;
  filters.serviceId = Number.isFinite(serviceIdQuery) && serviceIdQuery > 0 ? serviceIdQuery : undefined;

  if (actionQuery === "create") {
    resetCreateForm();
    if (Number.isFinite(serviceIdQuery) && serviceIdQuery > 0) {
      createForm.serviceId = serviceIdQuery;
    }
    if (titleQuery) {
      createForm.title = titleQuery;
    }
    createVisible.value = true;
  }
}

onMounted(async () => {
  applyRouteState();
  await fetchData();
});

watch(
  () => [
    route.query.keyword,
    route.query.status,
    route.query.priority,
    route.query.action,
    route.query.serviceId,
    route.query.title
  ],
  () => {
    applyRouteState();
  }
);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <div class="portal-badge">{{ copy.badge }}</div>
          <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
        <div class="portal-actions">
          <el-button @click="fetchData">{{ copy.refresh }}</el-button>
          <el-button type="primary" @click="openCreateDialog()">{{ copy.create }}</el-button>
        </div>
      </div>

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.total }}</h3>
          <strong>{{ tickets.length }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.open }}</h3>
          <strong>{{ openCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.processing }}</h3>
          <strong>{{ processingCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.waiting }}</h3>
          <strong>{{ waitingCount }}</strong>
        </div>
      </div>
    </section>

    <section class="portal-card portal-table-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.listTitle }}</h2>
          <div class="portal-panel__meta">{{ copy.listDesc }}</div>
        </div>
      </div>

      <div class="portal-toolbar" style="margin: 18px 20px 0">
        <el-input v-model="filters.keyword" :placeholder="copy.keywordPlaceholder" clearable />
        <el-select v-model="filters.status" :placeholder="copy.statusPlaceholder" clearable>
          <el-option :label="ticketStatusLabel('OPEN')" value="OPEN" />
          <el-option :label="ticketStatusLabel('PROCESSING')" value="PROCESSING" />
          <el-option :label="ticketStatusLabel('WAITING_CUSTOMER')" value="WAITING_CUSTOMER" />
          <el-option :label="ticketStatusLabel('CLOSED')" value="CLOSED" />
        </el-select>
        <el-select v-model="filters.priority" :placeholder="copy.priorityPlaceholder" clearable>
          <el-option :label="ticketPriorityLabel('LOW')" value="LOW" />
          <el-option :label="ticketPriorityLabel('NORMAL')" value="NORMAL" />
          <el-option :label="ticketPriorityLabel('HIGH')" value="HIGH" />
          <el-option :label="ticketPriorityLabel('URGENT')" value="URGENT" />
        </el-select>
        <el-select
          v-model="filters.serviceId"
          :placeholder="copy.servicePlaceholder"
          clearable
          filterable
          style="min-width: 220px"
        >
          <el-option v-for="item in serviceOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <div class="portal-summary" style="margin: 0; padding: 12px 16px">
          <div class="portal-summary-row" style="font-size: 12px">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredTickets.length }}</strong>
          </div>
        </div>
      </div>

      <el-table v-if="filteredTickets.length" :data="filteredTickets" border>
        <el-table-column prop="ticketNo" :label="copy.ticketNo" min-width="160" />
        <el-table-column prop="title" :label="copy.titleCol" min-width="240" />
        <el-table-column :label="copy.service" min-width="220">
          <template #default="{ row }">
            <template v-if="row.serviceId">
              <div>{{ row.serviceNo || "-" }}</div>
              <small>{{ row.productName || "-" }}</small>
            </template>
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column :label="copy.business" min-width="180">
          <template #default="{ row }">
            <div class="inline-tags">
              <small>{{ getLinkedService(row.serviceId)?.providerType || "-" }}</small>
              <small>{{ getLinkedService(row.serviceId)?.invoiceId ? copy.invoiceCenter : "-" }}</small>
            </div>
          </template>
        </el-table-column>
        <el-table-column :label="copy.priority" min-width="110">
          <template #default="{ row }">
            <el-tag :type="ticketPriorityType(row.priority)" effect="light">
              {{ ticketPriorityLabel(row.priority) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.status" min-width="130">
          <template #default="{ row }">
            <el-tag :type="ticketStatusType(row.status)" effect="light">
              {{ ticketStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.sla" min-width="160">
          <template #default="{ row }">
            <div class="inline-tags">
              <el-tag :type="ticketSlaType(row)" effect="light">{{ ticketSlaLabel(row) }}</el-tag>
              <small v-if="row.slaDeadlineAt">{{ row.slaDeadlineAt }}</small>
              <small v-else-if="row.autoCloseAt">{{ `${copy.autoClose}: ${row.autoCloseAt}` }}</small>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="copy.updatedAt" min-width="180" />
        <el-table-column :label="copy.action" min-width="200" fixed="right">
          <template #default="{ row }">
            <div class="table-actions">
              <el-button type="primary" link @click="openDetail(row)">{{ copy.detail }}</el-button>
              <el-button v-if="row.serviceId" type="primary" link @click="goToService(row.serviceId)">
                {{ copy.serviceDetail }}
              </el-button>
              <el-button
                v-if="getLinkedService(row.serviceId)?.invoiceId"
                type="primary"
                link
                @click="goToInvoice(row.serviceId)"
              >
                {{ copy.invoiceDetail }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>

    <el-dialog v-model="createVisible" :title="copy.createTitle" width="620px">
      <el-form label-position="top">
        <el-form-item :label="copy.linkedService">
          <el-select
            v-model="createForm.serviceId"
            clearable
            filterable
            style="width: 100%"
            :placeholder="copy.linkedServiceHint"
          >
            <el-option v-for="item in serviceOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <div class="ticket-form-grid">
          <el-form-item :label="copy.priority">
            <el-select v-model="createForm.priority" style="width: 100%">
              <el-option :label="ticketPriorityLabel('LOW')" value="LOW" />
              <el-option :label="ticketPriorityLabel('NORMAL')" value="NORMAL" />
              <el-option :label="ticketPriorityLabel('HIGH')" value="HIGH" />
              <el-option :label="ticketPriorityLabel('URGENT')" value="URGENT" />
            </el-select>
          </el-form-item>
          <el-form-item :label="copy.department">
            <el-select v-model="createForm.departmentName" style="width: 100%">
              <el-option v-for="item in departmentOptions" :key="item.key" :label="item.name" :value="item.name" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item :label="copy.createTicketTitle">
          <el-input v-model="createForm.title" maxlength="120" show-word-limit />
        </el-form-item>
        <el-form-item :label="copy.content">
          <el-input v-model="createForm.content" type="textarea" :rows="6" maxlength="1000" show-word-limit />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createVisible = false">{{ copy.cancel }}</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitCreate">{{ copy.submit }}</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="detailVisible" :title="copy.detailTitle" size="55%">
      <div v-if="detail" v-loading="detailLoading" class="ticket-detail">
        <div class="ticket-detail__head">
          <div>
            <div class="portal-badge">{{ detail.ticket.ticketNo }}</div>
            <h2 class="ticket-detail__title">{{ detail.ticket.title }}</h2>
            <div class="ticket-detail__meta">
              <span>{{ copy.createdAt }}: {{ detail.ticket.createdAt || "-" }}</span>
              <span>{{ copy.assignedAdmin }}: {{ detail.ticket.assignedAdminName || "-" }}</span>
              <span>{{ copy.sla }}: {{ ticketSlaLabel(detail.ticket) }}</span>
              <span v-if="detail.ticket.slaDeadlineAt">{{ detail.ticket.slaDeadlineAt }}</span>
              <span v-else-if="detail.ticket.autoCloseAt">{{ `${copy.autoClose}: ${detail.ticket.autoCloseAt}` }}</span>
            </div>
          </div>
          <div class="ticket-detail__tags">
            <el-tag :type="ticketPriorityType(detail.ticket.priority)" effect="light">
              {{ ticketPriorityLabel(detail.ticket.priority) }}
            </el-tag>
            <el-tag :type="ticketStatusType(detail.ticket.status)" effect="light">
              {{ ticketStatusLabel(detail.ticket.status) }}
            </el-tag>
            <el-tag :type="ticketSlaType(detail.ticket)" effect="light">
              {{ ticketSlaLabel(detail.ticket) }}
            </el-tag>
          </div>
        </div>

        <div class="ticket-detail__grid">
          <div class="ticket-detail__panel">
            <div class="ticket-detail__panel-head">
              <strong>{{ copy.origin }}</strong>
            </div>
            <div class="ticket-detail__content">{{ detail.ticket.content || "-" }}</div>
          </div>

          <div class="ticket-detail__panel">
            <div class="ticket-detail__panel-head">
              <strong>{{ copy.relatedBusiness }}</strong>
            </div>
            <div class="related-business">
              <div class="related-business__row">
                <span>{{ copy.service }}</span>
                <strong>{{ currentLinkedService?.serviceNo || detail.ticket.serviceNo || "-" }}</strong>
              </div>
              <div class="related-business__row">
                <span>{{ copy.product }}</span>
                <strong>{{ currentLinkedService?.productName || detail.ticket.productName || "-" }}</strong>
              </div>
              <div class="related-business__row">
                <span>{{ copy.serviceStatus }}</span>
                <strong>
                  {{ currentLinkedService ? formatPortalServiceStatus(localeStore.locale, currentLinkedService.status) : "-" }}
                </strong>
              </div>
              <div class="related-business__row">
                <span>{{ copy.providerType }}</span>
                <strong>{{ currentLinkedService?.providerType || "-" }}</strong>
              </div>
            </div>
            <div class="ticket-detail__actions">
              <el-button v-if="detail.ticket.serviceId" plain @click="goToService(detail.ticket.serviceId)">
                {{ copy.serviceDetail }}
              </el-button>
              <el-button v-if="currentLinkedService?.invoiceId" plain @click="goToInvoice(detail.ticket.serviceId)">
                {{ copy.invoiceDetail }}
              </el-button>
              <el-button v-if="detail.ticket.serviceId" type="primary" plain @click="openCreateDialog(detail.ticket.serviceId)">
                {{ copy.createFromService }}
              </el-button>
            </div>
          </div>
        </div>

        <div class="ticket-detail__grid">
          <div class="ticket-detail__panel">
            <div class="ticket-detail__panel-head">
              <strong>{{ copy.replyTitle }}</strong>
              <el-button
                v-if="detail.ticket.status !== 'CLOSED'"
                type="danger"
                plain
                size="small"
                :loading="submitLoading"
                @click="handleCloseTicket"
              >
                {{ copy.closeTicket }}
              </el-button>
            </div>
            <el-form v-if="detail.ticket.status !== 'CLOSED'" label-position="top">
              <el-form-item :label="copy.replyStatus">
                <el-select v-model="replyForm.status" style="width: 100%">
                  <el-option :label="ticketStatusLabel('OPEN')" value="OPEN" />
                  <el-option :label="ticketStatusLabel('PROCESSING')" value="PROCESSING" />
                  <el-option :label="ticketStatusLabel('CLOSED')" value="CLOSED" />
                </el-select>
              </el-form-item>
              <el-form-item>
                <el-input v-model="replyForm.content" type="textarea" :rows="5" :placeholder="copy.replyPlaceholder" />
              </el-form-item>
              <div class="portal-actions">
                <el-button type="primary" :loading="submitLoading" @click="submitReply">{{ copy.submitReply }}</el-button>
              </div>
            </el-form>
            <el-empty v-else :description="copy.closed" />
          </div>

          <div class="ticket-detail__panel">
            <div class="ticket-detail__panel-head">
              <strong>{{ copy.conversation }}</strong>
            </div>
            <div v-if="detail.replies.length" class="ticket-replies">
              <div v-for="item in detail.replies" :key="item.id" class="ticket-reply">
                <div class="ticket-reply__meta">
                  <div class="ticket-reply__title">
                    <el-tag :type="item.authorType === 'ADMIN' ? 'primary' : 'success'" effect="light">
                      {{ replyAuthorLabel(item.authorType) }}
                    </el-tag>
                    <strong>{{ item.authorName }}</strong>
                  </div>
                  <span>{{ item.createdAt }}</span>
                </div>
                <div class="ticket-reply__content">{{ item.content }}</div>
              </div>
            </div>
            <el-empty v-else :description="copy.noReplies" />
          </div>
        </div>
      </div>
    </el-drawer>
  </div>
</template>

<style scoped>
.portal-actions {
  display: flex;
  gap: 12px;
  align-items: center;
}

.table-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

.inline-tags {
  display: grid;
  gap: 4px;
}

.inline-tags small {
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}

.ticket-form-grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.ticket-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.ticket-detail__head {
  display: flex;
  justify-content: space-between;
  gap: 16px;
}

.ticket-detail__title {
  margin: 10px 0 8px;
  font-size: 22px;
}

.ticket-detail__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.ticket-detail__tags {
  display: flex;
  gap: 8px;
  align-items: flex-start;
}

.ticket-detail__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.ticket-detail__panel {
  border: 1px solid var(--el-border-color-light);
  border-radius: 18px;
  padding: 18px;
  background: rgba(255, 255, 255, 0.86);
}

.ticket-detail__panel-head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.ticket-detail__content {
  white-space: pre-wrap;
  line-height: 1.8;
}

.ticket-detail__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
  margin-top: 16px;
}

.related-business {
  display: grid;
  gap: 10px;
}

.related-business__row {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  font-size: 13px;
}

.related-business__row span {
  color: var(--el-text-color-secondary);
}

.ticket-replies {
  display: grid;
  gap: 12px;
}

.ticket-reply {
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(247, 249, 252, 0.9);
  border: 1px solid rgba(15, 23, 42, 0.06);
}

.ticket-reply__meta {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 10px;
  font-size: 12px;
  color: var(--el-text-color-secondary);
}

.ticket-reply__title {
  display: flex;
  gap: 8px;
  align-items: center;
}

.ticket-reply__content {
  white-space: pre-wrap;
  line-height: 1.8;
}

@media (max-width: 960px) {
  .ticket-detail__grid,
  .ticket-form-grid {
    grid-template-columns: 1fr;
  }

  .ticket-detail__head {
    flex-direction: column;
  }
}
</style>
