<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  fetchAutomationTasks,
  fetchOrderRequests,
  fetchWorkbenchOverview,
  type AutomationTask,
  type MetricCard,
  type OrderRequestRecord,
  type StatusBucket,
  type TrendPoint,
  type WorkbenchResponse
} from "@/api/admin";
import { formatMoney } from "@/utils/business";

type ShortcutCard = {
  key: string;
  title: string;
  value: number;
  hint: string;
  action: string;
  tone: string;
  onClick: () => void;
};

const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(false);
const activeTodoTab = ref("identity");
const overview = ref<WorkbenchResponse | null>(null);
const pendingOrderRequests = ref<OrderRequestRecord[]>([]);
const failedAutomationTasks = ref<AutomationTask[]>([]);
const blockedAutomationTasks = ref<AutomationTask[]>([]);
const automationSummary = ref({
  total: 0,
  running: 0,
  success: 0,
  failed: 0,
  blocked: 0
});

const trendMax = computed(() => {
  if (!overview.value?.revenueTrend.length) return 1;
  return Math.max(...overview.value.revenueTrend.map(item => item.amount), 1);
});

const serviceStatusTotal = computed(() =>
  (overview.value?.serviceStatus ?? []).reduce((total: number, item: StatusBucket) => total + item.count, 0)
);

const pendingOrderRequestCount = computed(() => pendingOrderRequests.value.length);
const automationAlertRows = computed(() => [...failedAutomationTasks.value, ...blockedAutomationTasks.value].slice(0, 6));
const automationAlertCount = computed(() => automationSummary.value.failed + automationSummary.value.blocked);

const copy = computed(() => ({
  eyebrow: pickLabel(localeStore.locale, "工作台", "Workbench"),
  title: pickLabel(localeStore.locale, "运营工作台", "Operations Workbench"),
  subtitle: pickLabel(
    localeStore.locale,
    "把订单申请、自动化异常、逾期账单、续费风险、工单和最近系统动作集中到一个后台总控台，管理员登录后就能直接进入处理状态。",
    "Bring order requests, automation alerts, overdue invoices, renewal risks, tickets, and recent activity into one admin control tower."
  ),
  identityReview: pickLabel(localeStore.locale, "实名审核", "Identity Review"),
  orderRequests: pickLabel(localeStore.locale, "订单申请", "Order Requests"),
  taskCenter: pickLabel(localeStore.locale, "任务中心", "Task Center"),
  refresh: pickLabel(localeStore.locale, "刷新工作台", "Refresh"),
  quickActions: pickLabel(localeStore.locale, "运营快捷入口", "Operations Shortcuts"),
  quickActionsDesc: pickLabel(
    localeStore.locale,
    "首页直接放置最常处理的入口，避免在菜单之间来回切换。",
    "Keep the highest-frequency workbenches on the first screen."
  ),
  revenueTrend: pickLabel(localeStore.locale, "近 7 日收款趋势", "Revenue Trend (7 days)"),
  revenueTrendDesc: pickLabel(
    localeStore.locale,
    "按支付流水聚合金额与笔数，便于观察每日回款节奏。",
    "Aggregated by payment flows to show daily collection trends."
  ),
  serviceStatus: pickLabel(localeStore.locale, "服务状态分布", "Service Status Distribution"),
  serviceStatusDesc: pickLabel(
    localeStore.locale,
    "按服务生命周期聚合，快速识别运行、暂停和终止服务占比。",
    "Grouped by lifecycle status to spot active, suspended, and terminated services."
  ),
  serviceTotal: pickLabel(localeStore.locale, "服务总量", "Total services"),
  todos: pickLabel(localeStore.locale, "运营待办", "Operations To-do"),
  todosDesc: pickLabel(
    localeStore.locale,
    "把实名、订单申请、逾期账单、即将到期服务、自动化异常和待跟进工单集中在一个面板处理。",
    "Review identities, order requests, overdue invoices, expiring services, automation alerts, and open tickets in one place."
  ),
  audits: pickLabel(localeStore.locale, "最近系统动作", "Recent System Activity"),
  auditsDesc: pickLabel(
    localeStore.locale,
    "来自审计日志，串联客户、订单、账单、服务和自动化动作。",
    "Pulled from audit logs across customers, orders, invoices, services, and automation."
  ),
  noIdentity: pickLabel(localeStore.locale, "暂无待审实名", "No pending identities"),
  noOrderRequests: pickLabel(localeStore.locale, "暂无待处理订单申请", "No pending order requests"),
  noInvoice: pickLabel(localeStore.locale, "暂无逾期账单", "No overdue invoices"),
  noService: pickLabel(localeStore.locale, "暂无即将到期服务", "No expiring services"),
  noAutomation: pickLabel(localeStore.locale, "暂无自动化异常", "No automation alerts"),
  noTicket: pickLabel(localeStore.locale, "暂无待跟进工单", "No open tickets"),
  noAudit: pickLabel(localeStore.locale, "暂无审计记录", "No audit logs"),
  openCustomer: pickLabel(localeStore.locale, "打开客户", "Open Customer"),
  openInvoice: pickLabel(localeStore.locale, "打开账单", "Open Invoice"),
  openService: pickLabel(localeStore.locale, "打开服务", "Open Service"),
  openRequest: pickLabel(localeStore.locale, "打开申请", "Open Request"),
  openTask: pickLabel(localeStore.locale, "打开任务", "Open Task"),
  openTicket: pickLabel(localeStore.locale, "打开工单", "Open Ticket"),
  identityTab: pickLabel(localeStore.locale, "实名审核", "Identity Review"),
  requestTab: pickLabel(localeStore.locale, "订单申请", "Order Requests"),
  invoiceTab: pickLabel(localeStore.locale, "逾期账单", "Overdue Invoices"),
  serviceTab: pickLabel(localeStore.locale, "即将到期服务", "Expiring Services"),
  automationTab: pickLabel(localeStore.locale, "自动化异常", "Automation Alerts"),
  ticketTab: pickLabel(localeStore.locale, "工单待跟进", "Open Tickets")
}));

const summaryCards = computed<MetricCard[]>(() => {
  const base = overview.value?.summaryCards ?? [];
  return [
    ...base,
    {
      key: "pendingOrderRequests",
      label: pickLabel(localeStore.locale, "待处理申请", "Pending Requests"),
      value: String(pendingOrderRequestCount.value),
      hint: pickLabel(localeStore.locale, "取消、续费、改价等申请待运营处理", "Pending cancel, renew, and price requests"),
      tone: "warning"
    },
    {
      key: "automationAlerts",
      label: pickLabel(localeStore.locale, "自动化异常", "Automation Alerts"),
      value: String(automationAlertCount.value),
      hint: pickLabel(localeStore.locale, "失败与阻塞任务需要人工介入", "Failed and blocked tasks need attention"),
      tone: "danger"
    }
  ];
});

const shortcutCards = computed<ShortcutCard[]>(() => [
  {
    key: "requests",
    title: pickLabel(localeStore.locale, "订单申请待办", "Order Requests"),
    value: pendingOrderRequestCount.value,
    hint: pickLabel(localeStore.locale, "集中处理取消、续费、改价等订单申请", "Review cancel, renew, and price requests"),
    action: pickLabel(localeStore.locale, "打开申请工作台", "Open Request Workbench"),
    tone: "warning",
    onClick: () => goOrderRequestWorkbench({ status: "PENDING" })
  },
  {
    key: "automation",
    title: pickLabel(localeStore.locale, "自动化异常", "Automation Alerts"),
    value: automationAlertCount.value,
    hint: pickLabel(localeStore.locale, "失败与阻塞任务优先跟进", "Prioritize failed and blocked tasks"),
    action: pickLabel(localeStore.locale, "打开任务中心", "Open Task Center"),
    tone: "danger",
    onClick: () => goAutomationWorkbench({ status: "FAILED" })
  },
  {
    key: "invoices",
    title: pickLabel(localeStore.locale, "逾期账单", "Overdue Invoices"),
    value: overview.value?.overdueInvoices.length ?? 0,
    hint: pickLabel(localeStore.locale, "快速切到逾期账单与催收链路", "Jump into overdue invoice follow-up"),
    action: pickLabel(localeStore.locale, "打开账单工作台", "Open Invoice Workbench"),
    tone: "warning",
    onClick: () => router.push({ path: "/billing/invoices", query: { status: "UNPAID" } })
  },
  {
    key: "tickets",
    title: pickLabel(localeStore.locale, "工单待办", "Ticket Queue"),
    value: overview.value?.openTickets.length ?? 0,
    hint: pickLabel(localeStore.locale, "客服与技术待处理工单入口", "Support and technical queues"),
    action: pickLabel(localeStore.locale, "打开工单中心", "Open Tickets"),
    tone: "primary",
    onClick: () => router.push("/tickets/list")
  },
  {
    key: "services",
    title: pickLabel(localeStore.locale, "即将到期服务", "Expiring Services"),
    value: overview.value?.expiringServices.length ?? 0,
    hint: pickLabel(localeStore.locale, "需要续费提醒与续费跟进", "Renewal reminder and retention follow-up"),
    action: pickLabel(localeStore.locale, "打开服务列表", "Open Services"),
    tone: "success",
    onClick: () => router.push("/services/list")
  },
  {
    key: "providers",
    title: pickLabel(localeStore.locale, "接口与上游", "Providers & Upstream"),
    value: automationSummary.value.total,
    hint: pickLabel(localeStore.locale, "查看账号、资源、同步和自动化任务", "Accounts, resources, sync, and automation"),
    action: pickLabel(localeStore.locale, "打开接口账户", "Open Provider Accounts"),
    tone: "primary",
    onClick: () => router.push("/providers/accounts")
  }
]);

function formatCurrency(value: number) {
  return formatMoney(localeStore.locale, value);
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

function goTicketDetail(ticketId: number) {
  void router.push(`/tickets/${ticketId}`);
}

function goOrderDetail(orderId: number) {
  void router.push(`/orders/detail/${orderId}`);
}

function goOrderRequestWorkbench(options?: { status?: string; orderId?: number; keyword?: string }) {
  void router.push({
    path: "/orders/requests",
    query: {
      status: options?.status,
      orderId: options?.orderId ? String(options.orderId) : undefined,
      keyword: options?.keyword
    }
  });
}

function goAutomationWorkbench(options?: {
  status?: string;
  orderId?: number;
  serviceId?: number;
  invoiceId?: number;
  channel?: string;
}) {
  void router.push({
    path: "/providers/automation",
    query: {
      status: options?.status,
      orderId: options?.orderId ? String(options.orderId) : undefined,
      serviceId: options?.serviceId ? String(options.serviceId) : undefined,
      invoiceId: options?.invoiceId ? String(options.invoiceId) : undefined,
      channel: options?.channel
    }
  });
}

async function loadOverview() {
  loading.value = true;
  try {
    const [overviewData, pendingRequestData, automationSummaryData, failedTaskData, blockedTaskData] = await Promise.all([
      fetchWorkbenchOverview(),
      fetchOrderRequests({ status: "PENDING", page: 1, limit: 6 }),
      fetchAutomationTasks({ limit: 1 }),
      fetchAutomationTasks({ status: "FAILED", limit: 6 }),
      fetchAutomationTasks({ status: "BLOCKED", limit: 6 })
    ]);
    overview.value = overviewData;
    pendingOrderRequests.value = pendingRequestData.items;
    automationSummary.value = automationSummaryData.summary;
    failedAutomationTasks.value = failedTaskData.items;
    blockedAutomationTasks.value = blockedTaskData.items;
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
        <el-button @click="goOrderRequestWorkbench({ status: 'PENDING' })">{{ copy.orderRequests }}</el-button>
        <el-button @click="router.push('/providers/automation')">{{ copy.taskCenter }}</el-button>
        <el-button type="primary" @click="loadOverview">{{ copy.refresh }}</el-button>
      </template>

      <template #metrics>
        <div class="metric-grid">
          <div v-for="card in summaryCards" :key="card.key" class="metric-card" :class="toneClass(card.tone)">
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

      <div class="page-card section-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">{{ copy.quickActions }}</h2>
            <p class="page-subtitle">{{ copy.quickActionsDesc }}</p>
          </div>
        </div>
        <div class="shortcut-grid">
          <button
            v-for="card in shortcutCards"
            :key="card.key"
            type="button"
            class="shortcut-card"
            :class="toneClass(card.tone)"
            @click="card.onClick"
          >
            <div class="shortcut-card__head">
              <strong>{{ card.title }}</strong>
              <span>{{ card.value }}</span>
            </div>
            <div class="shortcut-card__hint">{{ card.hint }}</div>
            <div class="shortcut-card__action">{{ card.action }}</div>
          </button>
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

          <el-tab-pane :label="`${copy.requestTab} (${pendingOrderRequests.length})`" name="request">
            <el-table :data="pendingOrderRequests" border stripe :empty-text="copy.noOrderRequests">
              <el-table-column prop="requestNo" label="申请单号" min-width="150" />
              <el-table-column label="类型" min-width="120">
                <template #default="{ row }">
                  {{ row.type === "CANCEL" ? "取消申请" : row.type === "RENEW" ? "续费请求" : "改价申请" }}
                </template>
              </el-table-column>
              <el-table-column prop="customerName" label="客户名称" min-width="160" />
              <el-table-column prop="productName" label="商品名称" min-width="220" show-overflow-tooltip />
              <el-table-column prop="summary" label="摘要" min-width="220" show-overflow-tooltip />
              <el-table-column prop="createdAt" label="申请时间" min-width="180" />
              <el-table-column label="操作" min-width="220" fixed="right">
                <template #default="{ row }">
                  <div class="inline-actions">
                    <el-button
                      type="primary"
                      link
                      @click="goOrderRequestWorkbench({ orderId: row.orderId, keyword: row.requestNo })"
                    >
                      {{ copy.openRequest }}
                    </el-button>
                    <el-button type="primary" link @click="goOrderDetail(row.orderId)">订单</el-button>
                  </div>
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

          <el-tab-pane :label="`${copy.automationTab} (${automationAlertCount})`" name="automation">
            <el-table :data="automationAlertRows" border stripe :empty-text="copy.noAutomation">
              <el-table-column prop="taskNo" label="任务编号" min-width="150" />
              <el-table-column prop="title" label="任务标题" min-width="220" show-overflow-tooltip />
              <el-table-column prop="channel" label="渠道" min-width="120" />
              <el-table-column prop="status" label="状态" min-width="120" />
              <el-table-column prop="message" label="结果说明" min-width="260" show-overflow-tooltip />
              <el-table-column prop="createdAt" label="创建时间" min-width="180" />
              <el-table-column label="操作" min-width="220" fixed="right">
                <template #default="{ row }">
                  <div class="inline-actions">
                    <el-button
                      type="primary"
                      link
                      @click="
                        goAutomationWorkbench({
                          status: row.status,
                          orderId: row.orderId || undefined,
                          serviceId: row.serviceId || undefined,
                          invoiceId: row.invoiceId || undefined,
                          channel: row.channel || undefined
                        })
                      "
                    >
                      {{ copy.openTask }}
                    </el-button>
                    <el-button v-if="row.serviceId" type="primary" link @click="goServiceDetail(row.serviceId)">服务</el-button>
                    <el-button v-else-if="row.orderId" type="primary" link @click="goOrderDetail(row.orderId)">订单</el-button>
                  </div>
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
              <el-table-column label="操作" min-width="140" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" link @click="goTicketDetail(row.ticketId)">{{ copy.openTicket }}</el-button>
                </template>
              </el-table-column>
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
.risk-card,
.shortcut-card {
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

.risk-grid,
.shortcut-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 16px;
}

.shortcut-card {
  text-align: left;
  cursor: pointer;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.shortcut-card:hover {
  transform: translateY(-1px);
  box-shadow: 0 12px 28px rgba(22, 55, 111, 0.08);
}

.shortcut-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  color: #16376f;
}

.shortcut-card__head strong {
  font-size: 15px;
}

.shortcut-card__head span {
  font-size: 24px;
  font-weight: 700;
}

.shortcut-card__hint {
  margin-top: 10px;
  color: #60728e;
  font-size: 13px;
  line-height: 1.6;
}

.shortcut-card__action {
  margin-top: 14px;
  color: #2563eb;
  font-size: 13px;
  font-weight: 600;
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
