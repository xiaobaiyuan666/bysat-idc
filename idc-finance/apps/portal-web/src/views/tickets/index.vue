<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { loadTickets, type PortalTicket } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import { formatPortalTicketStatus, portalTagTypeByStatus } from "@/utils/business";

const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const tickets = ref<PortalTicket[]>([]);
const keyword = ref("");
const status = ref("");

const filteredTickets = computed(() =>
  tickets.value.filter(item => {
    const matchedKeyword =
      !keyword.value || item.no.includes(keyword.value) || item.title.includes(keyword.value);
    const matchedStatus = !status.value || item.status === status.value;
    return matchedKeyword && matchedStatus;
  })
);

const openCount = computed(() => tickets.value.filter(item => item.status === "OPEN").length);
const closedCount = computed(() => tickets.value.filter(item => item.status === "CLOSED").length);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "工单中心", "Tickets"),
  title: pickLabel(localeStore.locale, "工单中心", "Tickets"),
  subtitle: pickLabel(
    localeStore.locale,
    "查看客户提交的支持请求和处理进度，后续可在此基础上扩展回复与评价。",
    "Review support requests and ticket progress."
  ),
  total: pickLabel(localeStore.locale, "工单总数", "Total Tickets"),
  open: pickLabel(localeStore.locale, "处理中", "Open"),
  closed: pickLabel(localeStore.locale, "已关闭", "Closed"),
  listTitle: pickLabel(localeStore.locale, "工单列表", "Ticket List"),
  listDesc: pickLabel(
    localeStore.locale,
    "支持按编号、标题和状态筛选。",
    "Filter by ticket number, title, and status."
  ),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索工单编号或标题", "Search ticket no. or title"),
  statusPlaceholder: pickLabel(localeStore.locale, "工单状态", "Ticket Status"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  ticketNo: pickLabel(localeStore.locale, "工单编号", "Ticket No."),
  titleCol: pickLabel(localeStore.locale, "标题", "Title"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  updatedAt: pickLabel(localeStore.locale, "更新时间", "Updated At"),
  empty: pickLabel(localeStore.locale, "暂无匹配的工单记录。", "No matching tickets."),
  loadError: pickLabel(localeStore.locale, "工单列表加载失败", "Failed to load tickets")
}));

async function fetchTickets() {
  loading.value = true;
  error.value = "";
  try {
    tickets.value = await loadTickets();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchTickets);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <div class="portal-badge">{{ copy.badge }}</div>
          <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="portal-grid portal-grid--three" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.total }}</h3>
          <strong>{{ tickets.length }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.open }}</h3>
          <strong>{{ openCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.closed }}</h3>
          <strong>{{ closedCount }}</strong>
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
        <el-input v-model="keyword" :placeholder="copy.keywordPlaceholder" clearable />
        <el-select v-model="status" :placeholder="copy.statusPlaceholder" clearable>
          <el-option :label="formatPortalTicketStatus(localeStore.locale, 'OPEN')" value="OPEN" />
          <el-option :label="formatPortalTicketStatus(localeStore.locale, 'CLOSED')" value="CLOSED" />
        </el-select>
        <div class="portal-summary" style="margin: 0; padding: 12px 16px">
          <div class="portal-summary-row" style="font-size: 12px">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredTickets.length }}</strong>
          </div>
        </div>
      </div>

      <el-table v-if="filteredTickets.length" :data="filteredTickets" border>
        <el-table-column prop="no" :label="copy.ticketNo" min-width="170" />
        <el-table-column prop="title" :label="copy.titleCol" min-width="240" />
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalTicketStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updatedAt" :label="copy.updatedAt" min-width="180" />
      </el-table>
      <div v-else class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>
  </div>
</template>
