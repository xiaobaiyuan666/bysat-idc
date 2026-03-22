<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  fetchWorkbenchOverview,
  type StatusBucket,
  type TrendPoint,
  type WorkbenchResponse
} from "@/api/admin";

const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(false);
const activeTodoTab = ref("identity");
const overview = ref<WorkbenchResponse | null>(null);

const trendMax = computed(() => {
  if (!overview.value?.revenueTrend.length) return 1;
  return Math.max(...overview.value.revenueTrend.map(item => item.amount), 1);
});

const serviceStatusTotal = computed(() =>
  (overview.value?.serviceStatus ?? []).reduce((total: number, item: StatusBucket) => total + item.count, 0)
);

const copy = computed(() => ({
  eyebrow: pickLabel(localeStore.locale, "工作台", "Workbench"),
  title: pickLabel(localeStore.locale, "运营工作台", "Operations Workbench"),
  subtitle: pickLabel(
    localeStore.locale,
    "汇总待办、风险、收款趋势、服务状态和最近系统动作，作为后台每日运营入口。",
    "Daily operations entry for pending items, risks, revenue trends, service status, and recent system activity."
  ),
  identityReview: pickLabel(localeStore.locale, "实名审核", "Identity Review"),
  taskCenter: pickLabel(localeStore.locale, "任务中心", "Task Center"),
  refresh: pickLabel(localeStore.locale, "刷新工作台", "Refresh"),
  revenueTrend: pickLabel(localeStore.locale, "近 7 日收款趋势", "Revenue Trend (7 days)"),
  revenueTrendDesc: pickLabel(localeStore.locale, "按支付流水聚合金额与笔数，便于观察每日回款节奏。", "Aggregated by payment flows to show daily collection trends."),
  serviceStatus: pickLabel(localeStore.locale, "服务状态分布", "Service Status Distribution"),
  serviceStatusDesc: pickLabel(localeStore.locale, "按服务生命周期聚合，快速识别运行、暂停和终止服务占比。", "Grouped by lifecycle status to spot active, suspended, and terminated services."),
  serviceTotal: pickLabel(localeStore.locale, "服务总量", "Total services"),
  todos: pickLabel(localeStore.locale, "运营待办", "Operations To-do"),
  todosDesc: pickLabel(localeStore.locale, "把实名、逾期账单、即将到期服务和待跟进工单集中在一个面板处理。", "Review identity submissions, overdue invoices, expiring services, and open tickets in one place."),
  audits: pickLabel(localeStore.locale, "最近系统动作", "Recent System Activity"),
  auditsDesc: pickLabel(localeStore.locale, "来自审计日志，串联客户、订单、账单、服务和自动化动作。", "Pulled from audit logs across customers, orders, invoices, services, and automation."),
  noIdentity: pickLabel(localeStore.locale, "暂无待审实名", "No pending identities"),
  noInvoice: pickLabel(localeStore.locale, "暂无逾期账单", "No overdue invoices"),
  noService: pickLabel(localeStore.locale, "暂无即将到期服务", "No expiring services"),
  noTicket: pickLabel(localeStore.locale, "暂无待跟进工单", "No open tickets"),
  noAudit: pickLabel(localeStore.locale, "暂无审计记录", "No audit logs"),
  openCustomer: pickLabel(localeStore.locale, "打开客户", "Open Customer"),
  openInvoice: pickLabel(localeStore.locale, "打开账单", "Open Invoice"),
  openService: pickLabel(localeStore.locale, "打开服务", "Open Service"),
  identityTab: pickLabel(localeStore.locale, "实名审核", "Identity Review"),
  invoiceTab: pickLabel(localeStore.locale, "逾期账单", "Overdue Invoices"),
  serviceTab: pickLabel(localeStore.locale, "即将到期服务", "Expiring Services"),
  ticketTab: pickLabel(localeStore.locale, "工单待跟进", "Open Tickets")
}));

function formatCurrency(value: number) {
  const prefix = localeStore.locale === "en-US" ? "$" : "¥";
  return `${prefix}${value.toFixed(2)}`;
}

function percentage(value: number, total: number) {
  if (!total) return 0;
  return Math.round((value / total) * 100);
}

function toneClass(tone: string) {
  return `tone-${tone || "primary"}`;
}

function trendWidth(item: TrendPoint) {
  return `${Math.max((item.amount / trendMax.value) * 100, item.amount > 0 ? 12 : 4)}%`;
}

function serviceWidth(item: StatusBucket) {
  return `${Math.max(percentage(item.count, serviceStatusTotal.value), item.count > 0 ? 10 : 0)}%`;
}

function goCustomerDetail(customerId: number) {
  void router.push(`/customer/detail/${customerId}`);
}

function goInvoiceDetail(invoiceId: number) {
  void router.push(`/billing/invoices/${invoiceId}`);
}

function goServiceDetail(serviceId: number) {
  void router.push(`/services/detail/${serviceId}`);
}

async function loadOverview() {
  loading.value = true;
  try {
    overview.value = await fetchWorkbenchOverview();
  } finally {
    loading.value = false;
  }
}

onMounted(loadOverview);
</script>

<template>
  <div v-loading="loading" class="workbench-page">
    <PageWorkbench :eyebrow="copy.eyebrow" :title="copy.title" :subtitle="copy.subtitle">
      <template #actions>
        <el-button @click="router.push('/customer/identities')">{{ copy.identityReview }}</el-button>
        <el-button @click="router.push('/providers/automation')">{{ copy.taskCenter }}</el-button>
        <el-button type="primary" @click="loadOverview">{{ copy.refresh }}</el-button>
      </template>

      <template #metrics>
        <div class="metric-grid">
          <div
            v-for="card in overview?.summaryCards ?? []"
            :key="card.key"
            class="metric-card"
            :class="toneClass(card.tone)"
          >
            <div class="metric-card__label">{{ card.label }}</div>
            <div class="metric-card__value">{{ card.value }}</div>
            <div class="metric-card__hint">{{ card.hint }}</div>
          </div>
        </div>
      </template>

      <div class="risk-grid">
        <div
          v-for="card in overview?.riskCards ?? []"
          :key="card.key"
          class="risk-card"
          :class="toneClass(card.tone)"
        >
          <div class="risk-card__label">{{ card.label }}</div>
          <div class="risk-card__value">{{ card.value }}</div>
          <div class="risk-card__hint">{{ card.hint }}</div>
        </div>
      </div>

      <div class="sub-grid">
        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">{{ copy.revenueTrend }}</h2>
              <p class="page-subtitle">{{ copy.revenueTrendDesc }}</p>
            </div>
          </div>
          <div class="trend-list">
            <div v-for="item in overview?.revenueTrend ?? []" :key="item.label" class="trend-item">
              <div class="trend-item__meta">
                <strong>{{ item.label }}</strong>
                <span>{{ formatCurrency(item.amount) }} / {{ item.count }}</span>
              </div>
              <div class="trend-item__bar">
                <span class="trend-item__fill" :style="{ width: trendWidth(item) }"></span>
              </div>
            </div>
          </div>
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">{{ copy.serviceStatus }}</h2>
              <p class="page-subtitle">{{ copy.serviceStatusDesc }}</p>
            </div>
          </div>
          <div class="status-stack">
            <div v-for="item in overview?.serviceStatus ?? []" :key="item.name" class="status-stack__item">
              <div class="status-stack__meta">
                <strong>{{ item.name }}</strong>
                <span>{{ item.count }}</span>
              </div>
              <div class="status-stack__bar">
                <span class="status-stack__fill" :style="{ width: serviceWidth(item) }"></span>
              </div>
            </div>
          </div>
          <div class="status-total">{{ copy.serviceTotal }}：{{ serviceStatusTotal }}</div>
        </div>
      </div>

      <div class="page-card section-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">{{ copy.todos }}</h2>
            <p class="page-subtitle">{{ copy.todosDesc }}</p>
          </div>
        </div>

        <el-tabs v-model="activeTodoTab">
          <el-tab-pane :label="`${copy.identityTab} (${overview?.pendingIdentities.length ?? 0})`" name="identity">
            <el-table :data="overview?.pendingIdentities ?? []" border stripe :empty-text="copy.noIdentity">
              <el-table-column prop="customerNo" label="客户编号" min-width="140" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column prop="subjectName" label="实名主体" min-width="220" />
              <el-table-column prop="submittedAt" label="提交时间" min-width="180" />
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="goCustomerDetail(row.customerId)">{{ copy.openCustomer }}</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane :label="`${copy.invoiceTab} (${overview?.overdueInvoices.length ?? 0})`" name="invoice">
            <el-table :data="overview?.overdueInvoices ?? []" border stripe :empty-text="copy.noInvoice">
              <el-table-column prop="invoiceNo" label="账单编号" min-width="160" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column label="应收金额" min-width="140">
                <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="dueAt" label="到期时间" min-width="180" />
              <el-table-column prop="daysOverdue" label="逾期天数" min-width="120" />
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="goInvoiceDetail(row.invoiceId)">{{ copy.openInvoice }}</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane :label="`${copy.serviceTab} (${overview?.expiringServices.length ?? 0})`" name="service">
            <el-table :data="overview?.expiringServices ?? []" border stripe :empty-text="copy.noService">
              <el-table-column prop="serviceNo" label="服务编号" min-width="160" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column prop="productName" label="产品名称" min-width="220" />
              <el-table-column prop="status" label="状态" min-width="120" />
              <el-table-column prop="nextDueAt" label="到期时间" min-width="180" />
              <el-table-column prop="daysRemaining" label="剩余天数" min-width="120" />
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="goServiceDetail(row.serviceId)">{{ copy.openService }}</el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane :label="`${copy.ticketTab} (${overview?.openTickets.length ?? 0})`" name="ticket">
            <el-table :data="overview?.openTickets ?? []" border stripe :empty-text="copy.noTicket">
              <el-table-column prop="ticketNo" label="工单编号" min-width="160" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column prop="title" label="工单标题" min-width="260" />
              <el-table-column prop="status" label="状态" min-width="120" />
              <el-table-column prop="updatedAt" label="最后更新时间" min-width="180" />
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>

      <div class="page-card section-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">{{ copy.audits }}</h2>
            <p class="page-subtitle">{{ copy.auditsDesc }}</p>
          </div>
        </div>
        <el-table :data="overview?.recentAudits ?? []" border stripe :empty-text="copy.noAudit">
          <el-table-column prop="createdAt" label="时间" min-width="180" />
          <el-table-column prop="actor" label="操作人" min-width="140" />
          <el-table-column prop="action" label="动作" min-width="180" />
          <el-table-column prop="target" label="目标对象" min-width="180" />
          <el-table-column prop="description" label="说明" min-width="280" />
        </el-table>
      </div>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.workbench-page {
  display: grid;
  gap: 16px;
}

.metric-card,
.risk-card {
  border-radius: 16px;
  border: 1px solid #e6edf7;
  background: #fff;
  padding: 16px;
}

.metric-card__label,
.risk-card__label {
  color: #5e7093;
  font-size: 13px;
}

.metric-card__value,
.risk-card__value {
  margin-top: 10px;
  font-size: 28px;
  font-weight: 700;
  color: #16376f;
}

.metric-card__hint,
.risk-card__hint {
  margin-top: 10px;
  color: #7083a6;
  font-size: 12px;
  line-height: 1.6;
}

.risk-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.section-card {
  padding: 20px;
}

.trend-list,
.status-stack {
  display: grid;
  gap: 12px;
}

.trend-item__meta,
.status-stack__meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  color: #526684;
  font-size: 13px;
}

.trend-item__bar,
.status-stack__bar {
  height: 10px;
  border-radius: 999px;
  background: #edf3fb;
  overflow: hidden;
}

.trend-item__fill,
.status-stack__fill {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #2563eb 0%, #60a5fa 100%);
}

.status-total {
  margin-top: 12px;
  color: #60728e;
  font-size: 13px;
}

.tone-primary {
  background: linear-gradient(180deg, #f5f9ff 0%, #ffffff 100%);
}

.tone-success {
  background: linear-gradient(180deg, #f2fbf6 0%, #ffffff 100%);
}

.tone-warning {
  background: linear-gradient(180deg, #fff8ef 0%, #ffffff 100%);
}

.tone-danger {
  background: linear-gradient(180deg, #fff5f5 0%, #ffffff 100%);
}
</style>
