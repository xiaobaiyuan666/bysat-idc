<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import {
  claimTicket,
  fetchAdmins,
  fetchTicketDepartments,
  fetchTickets,
  type TicketDepartment,
  type TicketQuery,
  type TicketRecord
} from "@/api/admin";

type TabKey = "ALL" | "OPEN" | "PROCESSING" | "WAITING_CUSTOMER" | "CLOSED";

const route = useRoute();
const router = useRouter();

const text = {
  eyebrow: "\u8d22\u52a1 / \u5de5\u5355",
  title: "\u5de5\u5355\u4e2d\u5fc3",
  subtitle:
    "\u6309\u7167 IDC \u8d22\u52a1\u7cfb\u7edf\u7684\u5de5\u5355\u8ddf\u8fdb\u65b9\u5f0f\uff0c\u96c6\u4e2d\u5904\u7406\u7528\u6237\u63d0\u5355\u3001\u56de\u590d\u3001\u5206\u914d\u548c\u5173\u5355\u8fdb\u5ea6\u3002",
  refresh: "\u5237\u65b0\u5217\u8868",
  create: "\u65b0\u5efa\u5de5\u5355",
  settings: "\u5de5\u5355\u914d\u7f6e",
  total: "\u5de5\u5355\u603b\u6570",
  open: "\u5f85\u5904\u7406",
  processing: "\u5904\u7406\u4e2d",
  waitingCustomer: "\u5f85\u5ba2\u6237\u56de\u590d",
  closed: "\u5df2\u5173\u95ed",
  keywordPlaceholder: "\u641c\u7d22\u5de5\u5355\u53f7/\u6807\u9898/\u5ba2\u6237/\u670d\u52a1\u53f7",
  customerIdPlaceholder: "\u5ba2\u6237 ID",
  serviceIdPlaceholder: "\u670d\u52a1 ID",
  priorityPlaceholder: "\u4f18\u5148\u7ea7",
  departmentPlaceholder: "\u90e8\u95e8",
  assigneePlaceholder: "\u8d1f\u8d23\u4eba",
  apply: "\u5e94\u7528\u7b5b\u9009",
  reset: "\u91cd\u7f6e",
  currentPage: "\u5f53\u524d\u9875",
  selected: "\u5f53\u524d\u72b6\u6001",
  ticketNo: "\u5de5\u5355\u53f7",
  customer: "\u5ba2\u6237",
  service: "\u5173\u8054\u670d\u52a1",
  titleCol: "\u6807\u9898",
  priority: "\u4f18\u5148\u7ea7",
  status: "\u72b6\u6001",
  sla: "SLA",
  department: "\u90e8\u95e8",
  assignedAdmin: "\u8d1f\u8d23\u4eba",
  lastReplyAt: "\u6700\u540e\u56de\u590d",
  updatedAt: "\u66f4\u65b0\u65f6\u95f4",
  action: "\u64cd\u4f5c",
  claim: "\u63a5\u5355",
  openWorkbench: "\u8fdb\u5165\u5de5\u4f5c\u53f0",
  empty: "\u6682\u65e0\u7b26\u5408\u6761\u4ef6\u7684\u5de5\u5355"
} as const;

const loading = ref(false);
const tickets = ref<TicketRecord[]>([]);
const total = ref(0);
const activeTab = ref<TabKey>("ALL");
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

const pagination = reactive({
  page: 1,
  limit: 20
});

const filters = reactive({
  keyword: "",
  customerId: "",
  serviceId: "",
  priority: "",
  departmentName: "",
  assignedAdminId: ""
});

const statusTabs = computed(() => [
  { key: "ALL", label: "\u5168\u90e8", count: total.value },
  { key: "OPEN", label: text.open, count: tickets.value.filter(item => item.status === "OPEN").length },
  { key: "PROCESSING", label: text.processing, count: tickets.value.filter(item => item.status === "PROCESSING").length },
  {
    key: "WAITING_CUSTOMER",
    label: text.waitingCustomer,
    count: tickets.value.filter(item => item.status === "WAITING_CUSTOMER").length
  },
  { key: "CLOSED", label: text.closed, count: tickets.value.filter(item => item.status === "CLOSED").length }
]);

function buildQuery(): TicketQuery {
  const query: TicketQuery = {
    page: pagination.page,
    limit: pagination.limit
  };
  if (activeTab.value !== "ALL") query.status = activeTab.value;
  if (filters.keyword.trim()) query.keyword = filters.keyword.trim();
  if (filters.customerId.trim()) query.customerId = filters.customerId.trim();
  if (filters.serviceId.trim()) query.serviceId = filters.serviceId.trim();
  if (filters.priority) query.priority = filters.priority;
  if (filters.departmentName) query.departmentName = filters.departmentName;
  if (filters.assignedAdminId) query.assignedAdminId = filters.assignedAdminId;
  return query;
}

function applyRouteFilters() {
  const statusQuery = typeof route.query.status === "string" ? route.query.status : "";
  const keywordQuery = typeof route.query.keyword === "string" ? route.query.keyword : "";
  const customerIdQuery = typeof route.query.customerId === "string" ? route.query.customerId : "";
  const serviceIdQuery = typeof route.query.serviceId === "string" ? route.query.serviceId : "";
  const priorityQuery = typeof route.query.priority === "string" ? route.query.priority : "";
  const departmentQuery = typeof route.query.departmentName === "string" ? route.query.departmentName : "";
  const assignedAdminIdQuery =
    typeof route.query.assignedAdminId === "string" ? route.query.assignedAdminId : "";

  activeTab.value = ["OPEN", "PROCESSING", "WAITING_CUSTOMER", "CLOSED"].includes(statusQuery)
    ? (statusQuery as TabKey)
    : "ALL";
  filters.keyword = keywordQuery;
  filters.customerId = customerIdQuery;
  filters.serviceId = serviceIdQuery;
  filters.priority = priorityQuery;
  filters.departmentName = departmentQuery;
  filters.assignedAdminId = assignedAdminIdQuery;
}

async function loadData() {
  loading.value = true;
  try {
    const result = await fetchTickets(buildQuery());
    tickets.value = result.items;
    total.value = result.total;
  } finally {
    loading.value = false;
  }
}

async function handleClaim(row: TicketRecord) {
  await claimTicket(row.id);
  ElMessage.success("\u5de5\u5355\u5df2\u63a5\u5355");
  await loadData();
}

async function loadBaseOptions() {
  const [departmentResult, adminItems] = await Promise.all([fetchTicketDepartments(), fetchAdmins()]);
  departments.value = departmentResult.items.filter(item => item.enabled);
  admins.value = adminItems;
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.customerId = "";
  filters.serviceId = "";
  filters.priority = "";
  filters.departmentName = "";
  filters.assignedAdminId = "";
  pagination.page = 1;
  void loadData();
}

function applyFilters() {
  pagination.page = 1;
  void loadData();
}

function openDetail(row: TicketRecord) {
  void router.push(`/tickets/${row.id}`);
}

function openServiceDetail(row: TicketRecord) {
  if (!row.serviceId) return;
  void router.push(`/services/${row.serviceId}`);
}

function statusLabel(status: string) {
  const mapping: Record<string, string> = {
    OPEN: text.open,
    PROCESSING: text.processing,
    WAITING_CUSTOMER: text.waitingCustomer,
    CLOSED: text.closed
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

function slaLabel(row: TicketRecord) {
  if (row.slaPaused) return "\u6682\u505c\u8ba1\u65f6";
  const mapping: Record<string, string> = {
    ON_TRACK: "\u6b63\u5e38",
    AT_RISK: "\u5373\u5c06\u8d85\u65f6",
    BREACHED: "\u5df2\u8d85\u65f6",
    MET: "\u5df2\u5b8c\u6210",
    UNKNOWN: "\u672a\u8ba1\u7b97"
  };
  return mapping[row.slaStatus] || row.slaStatus || "-";
}

function slaType(row: TicketRecord) {
  if (row.slaPaused) return "info";
  const mapping: Record<string, string> = {
    ON_TRACK: "success",
    AT_RISK: "warning",
    BREACHED: "danger",
    MET: "success",
    UNKNOWN: "info"
  };
  return mapping[row.slaStatus] ?? "info";
}

function handlePageChange(page: number) {
  pagination.page = page;
  void loadData();
}

function handleLimitChange(limit: number) {
  pagination.limit = limit;
  pagination.page = 1;
  void loadData();
}

watch(activeTab, () => {
  pagination.page = 1;
  void loadData();
});

onMounted(() => {
  applyRouteFilters();
  void loadBaseOptions();
  void loadData();
});

watch(
  () => [
    route.query.status,
    route.query.keyword,
    route.query.customerId,
    route.query.serviceId,
    route.query.priority,
    route.query.departmentName,
    route.query.assignedAdminId
  ],
  () => {
    applyRouteFilters();
    pagination.page = 1;
    void loadData();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench :eyebrow="text.eyebrow" :title="text.title" :subtitle="text.subtitle">
      <template #actions>
        <el-button type="primary" @click="router.push('/tickets/create')">{{ text.create }}</el-button>
        <el-button @click="router.push('/tickets/settings')">{{ text.settings }}</el-button>
        <el-button @click="loadData">{{ text.refresh }}</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>{{ text.total }}</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.open }}</span>
            <strong>{{ tickets.filter(item => item.status === "OPEN").length }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.processing }}</span>
            <strong>{{ tickets.filter(item => item.status === "PROCESSING").length }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.waitingCustomer }}</span>
            <strong>{{ tickets.filter(item => item.status === "WAITING_CUSTOMER").length }}</strong>
          </div>
          <div class="summary-pill">
            <span>{{ text.closed }}</span>
            <strong>{{ tickets.filter(item => item.status === "CLOSED").length }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="filters.keyword" :placeholder="text.keywordPlaceholder" clearable />
          <el-input v-model="filters.customerId" :placeholder="text.customerIdPlaceholder" clearable />
          <el-input v-model="filters.serviceId" :placeholder="text.serviceIdPlaceholder" clearable />
          <el-select v-model="filters.priority" :placeholder="text.priorityPlaceholder" clearable>
            <el-option :label="priorityLabel('LOW')" value="LOW" />
            <el-option :label="priorityLabel('NORMAL')" value="NORMAL" />
            <el-option :label="priorityLabel('HIGH')" value="HIGH" />
            <el-option :label="priorityLabel('URGENT')" value="URGENT" />
          </el-select>
          <el-select v-model="filters.departmentName" :placeholder="text.departmentPlaceholder" clearable>
            <el-option
              v-for="item in departments"
              :key="item.key"
              :label="item.name"
              :value="item.name"
            />
          </el-select>
          <el-select v-model="filters.assignedAdminId" :placeholder="text.assigneePlaceholder" clearable>
            <el-option
              v-for="item in admins"
              :key="item.id"
              :label="item.displayName || item.username"
              :value="String(item.id)"
            />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">{{ text.apply }}</el-button>
            <el-button plain @click="resetFilters">{{ text.reset }}</el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>{{ text.currentPage }} {{ tickets.length }} {{ text.selected }}</strong>
          <span>{{ activeTab === "ALL" ? "ALL" : statusLabel(activeTab) }}</span>
        </div>
      </div>

      <el-table :data="tickets" border stripe :empty-text="text.empty">
        <el-table-column prop="ticketNo" :label="text.ticketNo" min-width="160" />
        <el-table-column :label="text.customer" min-width="170">
          <template #default="{ row }">
            <div>{{ row.customerName }}</div>
            <small>ID: {{ row.customerId }}</small>
          </template>
        </el-table-column>
        <el-table-column :label="text.service" min-width="180">
          <template #default="{ row }">
            <div>{{ row.serviceNo || "-" }}</div>
            <small>{{ row.productName || "-" }}</small>
          </template>
        </el-table-column>
        <el-table-column prop="title" :label="text.titleCol" min-width="220" show-overflow-tooltip />
        <el-table-column :label="text.priority" min-width="100">
          <template #default="{ row }">
            <el-tag :type="priorityType(row.priority)" effect="light">{{ priorityLabel(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="text.status" min-width="130">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="text.sla" min-width="150">
          <template #default="{ row }">
            <div class="inline-tags">
              <el-tag :type="slaType(row)" effect="light">{{ slaLabel(row) }}</el-tag>
              <small v-if="row.slaDeadlineAt">{{ row.slaDeadlineAt }}</small>
              <small v-else-if="row.autoCloseAt">{{ `Auto: ${row.autoCloseAt}` }}</small>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="departmentName" :label="text.department" min-width="140" />
        <el-table-column prop="assignedAdminName" :label="text.assignedAdmin" min-width="140" />
        <el-table-column prop="lastReplyAt" :label="text.lastReplyAt" min-width="170" />
        <el-table-column prop="updatedAt" :label="text.updatedAt" min-width="170" />
        <el-table-column :label="text.action" min-width="130" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button
                v-if="!row.assignedAdminId && row.status !== 'CLOSED'"
                type="warning"
                link
                @click="handleClaim(row)"
              >
                {{ text.claim }}
              </el-button>
              <el-button v-if="row.serviceId" type="primary" link @click="openServiceDetail(row)">查看服务</el-button>
              <el-button type="primary" link @click="openDetail(row)">{{ text.openWorkbench }}</el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div class="table-pagination">
        <el-pagination
          background
          layout="total, sizes, prev, pager, next"
          :total="total"
          :current-page="pagination.page"
          :page-size="pagination.limit"
          :page-sizes="[20, 50, 100]"
          @current-change="handlePageChange"
          @size-change="handleLimitChange"
        />
      </div>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.inline-tags {
  display: grid;
  gap: 4px;
}

.inline-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.inline-tags small {
  color: var(--el-text-color-secondary);
  line-height: 1.4;
}
</style>
