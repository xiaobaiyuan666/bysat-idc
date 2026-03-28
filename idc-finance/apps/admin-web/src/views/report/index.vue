<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchAutomationTasks,
  fetchReportOverview,
  fetchWorkbenchOverview,
  type AutomationTask,
  type ReportOverviewResponse,
  type StatusBucket,
  type TrendPoint,
  type WorkbenchResponse
} from "@/api/admin";

type TodoTab = "identity" | "invoice" | "service" | "ticket";

const router = useRouter();
const loading = ref(false);
const activeTodoTab = ref<TodoTab>("identity");
const reportOverview = ref<ReportOverviewResponse | null>(null);
const operationsWorkbench = ref<WorkbenchResponse | null>(null);
const automationAlerts = ref<AutomationTask[]>([]);

const trendMax = computed(() => {
  const revenue = reportOverview.value?.revenueTrend ?? [];
  const refund = reportOverview.value?.refundTrend ?? [];
  return Math.max(...[...revenue, ...refund].map(item => item.amount), 1);
});

const customerGroupTotal = computed(() =>
  totalCount(reportOverview.value?.customerGroups ?? [])
);

const quickEntries = computed(() => [
  {
    key: "identities",
    label: "待实名审核",
    value: operationsWorkbench.value?.pendingIdentities.length ?? 0,
    hint: "进入实名审核工作台",
    tone: "warning",
    path: "/customer/identities"
  },
  {
    key: "overdue",
    label: "逾期账单",
    value: operationsWorkbench.value?.overdueInvoices.length ?? 0,
    hint: "进入应收与催收处理",
    tone: "danger",
    path: "/billing/invoices?status=UNPAID"
  },
  {
    key: "expiring",
    label: "即将到期服务",
    value: operationsWorkbench.value?.expiringServices.length ?? 0,
    hint: "进入续费与到期运营",
    tone: "warning",
    path: "/services/list"
  },
  {
    key: "tickets",
    label: "待跟进工单",
    value: operationsWorkbench.value?.openTickets.length ?? 0,
    hint: "进入工单中心",
    tone: "primary",
    path: "/tickets/list"
  },
  {
    key: "automation",
    label: "自动化异常",
    value: automationAlerts.value.length,
    hint: "进入自动化任务中心",
    tone: "danger",
    path: "/providers/automation?status=FAILED"
  }
]);

function formatCurrency(value: number) {
  return `¥${value.toFixed(2)}`;
}

function toneClass(tone: string) {
  return `tone-${tone || "primary"}`;
}

function trendWidth(item: TrendPoint) {
  return `${Math.max((item.amount / trendMax.value) * 100, item.amount > 0 ? 10 : 4)}%`;
}

function ratioWidth(item: StatusBucket, total: number) {
  if (!total) return "0%";
  return `${Math.max((item.count / total) * 100, item.count > 0 ? 8 : 0)}%`;
}

function totalCount(items: StatusBucket[]) {
  return items.reduce((total, item) => total + item.count, 0);
}

function statusTagType(status: string) {
  const normalized = status.toUpperCase();
  if (normalized === "FAILED") return "danger";
  if (normalized === "BLOCKED") return "warning";
  if (normalized === "RUNNING") return "primary";
  if (normalized === "SUCCESS") return "success";
  return "info";
}

function go(path: string) {
  void router.push(path);
}

function goCustomerDetail(customerId: number) {
  void router.push(`/customer/detail/${customerId}`);
}

function goCustomerList(query: Record<string, string>) {
  void router.push({
    path: "/customer/list",
    query
  });
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

function goAutomationTask(task: AutomationTask) {
  const query: Record<string, string> = {
    status: task.status || "FAILED"
  };
  if (task.serviceId) query.serviceId = String(task.serviceId);
  if (task.orderId) query.orderId = String(task.orderId);
  if (task.invoiceId) query.invoiceId = String(task.invoiceId);
  if (task.sourceId) query.sourceId = String(task.sourceId);
  if (task.sourceType) query.sourceType = task.sourceType;
  void router.push({ path: "/providers/automation", query });
}

function goProductOrders(productName: string) {
  void router.push({
    path: "/orders/list",
    query: {
      productName
    }
  });
}

function goProductCatalog(productName: string) {
  void router.push({
    path: "/catalog/products",
    query: {
      keyword: productName
    }
  });
}

async function loadOverview() {
  loading.value = true;
  try {
    const [reportData, workbenchData, failedTasks, blockedTasks] = await Promise.all([
      fetchReportOverview(),
      fetchWorkbenchOverview(),
      fetchAutomationTasks({ status: "FAILED", limit: 6 }),
      fetchAutomationTasks({ status: "BLOCKED", limit: 6 })
    ]);

    reportOverview.value = reportData;
    operationsWorkbench.value = workbenchData;
    automationAlerts.value = [...failedTasks.items, ...blockedTasks.items]
      .sort((left, right) => String(right.createdAt).localeCompare(String(left.createdAt)))
      .slice(0, 8);
  } finally {
    loading.value = false;
  }
}

onMounted(loadOverview);
</script>

<template>
  <div v-loading="loading" class="report-page">
    <PageWorkbench
      eyebrow="后台 / 报表中心"
      title="运营报表中心"
      subtitle="把财务走势、待办风险、客户分层、自动化异常和最近系统动作集中到一个后台入口，方便管理员先看全局再处理细项。"
    >
      <template #actions>
        <el-button @click="go('/billing/invoices')">账单工作台</el-button>
        <el-button @click="go('/services/list')">服务工作台</el-button>
        <el-button @click="go('/providers/automation')">自动化任务</el-button>
        <el-button @click="go('/customer/groups')">客户分组</el-button>
        <el-button type="primary" @click="loadOverview">刷新报表</el-button>
      </template>

      <template #metrics>
        <div class="headline-grid">
          <div v-for="item in reportOverview?.headline ?? []" :key="item.key" class="headline-card">
            <div class="headline-card__label">{{ item.label }}</div>
            <div class="headline-card__value">{{ item.value }}</div>
            <div class="headline-card__hint">{{ item.hint }}</div>
          </div>
        </div>
      </template>

      <div class="quick-grid">
        <button
          v-for="item in quickEntries"
          :key="item.key"
          type="button"
          class="quick-card"
          :class="toneClass(item.tone)"
          @click="go(item.path)"
        >
          <span class="quick-card__label">{{ item.label }}</span>
          <strong class="quick-card__value">{{ item.value }}</strong>
          <span class="quick-card__hint">{{ item.hint }}</span>
        </button>
      </div>

      <div class="sub-grid">
        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">近 30 日收款与退款走势</h2>
              <p class="page-subtitle">按天查看收入、退款和单量波峰，适合财务复盘与运营观察。</p>
            </div>
          </div>

          <div class="trend-layout">
            <div>
              <div class="block-title">收款趋势</div>
              <div class="trend-list">
                <div
                  v-for="item in reportOverview?.revenueTrend ?? []"
                  :key="`revenue-${item.label}`"
                  class="trend-item"
                >
                  <div class="trend-item__meta">
                    <strong>{{ item.label }}</strong>
                    <span>{{ formatCurrency(item.amount) }} / {{ item.count }} 笔</span>
                  </div>
                  <div class="trend-item__bar">
                    <span class="trend-item__fill" :style="{ width: trendWidth(item) }"></span>
                  </div>
                </div>
              </div>
            </div>

            <div>
              <div class="block-title">退款趋势</div>
              <div class="trend-list">
                <div
                  v-for="item in reportOverview?.refundTrend ?? []"
                  :key="`refund-${item.label}`"
                  class="trend-item"
                >
                  <div class="trend-item__meta">
                    <strong>{{ item.label }}</strong>
                    <span>{{ formatCurrency(item.amount) }} / {{ item.count }} 笔</span>
                  </div>
                  <div class="trend-item__bar trend-item__bar--warning">
                    <span
                      class="trend-item__fill trend-item__fill--warning"
                      :style="{ width: trendWidth(item) }"
                    ></span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">财务与服务结构</h2>
              <p class="page-subtitle">集中查看账单状态、服务状态、支付渠道和计费周期，判断当前运营质量。</p>
            </div>
          </div>

          <div class="ratio-panel">
            <div>
              <div class="ratio-title">账单状态</div>
              <div v-for="item in reportOverview?.invoiceStatus ?? []" :key="`invoice-${item.name}`" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} / {{ formatCurrency(item.amount) }}</span>
                </div>
                <div class="ratio-item__bar">
                  <span :style="{ width: ratioWidth(item, totalCount(reportOverview?.invoiceStatus ?? [])) }"></span>
                </div>
              </div>
            </div>

            <div>
              <div class="ratio-title">服务状态</div>
              <div v-for="item in reportOverview?.serviceStatus ?? []" :key="`service-${item.name}`" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} 台</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--success">
                  <span :style="{ width: ratioWidth(item, totalCount(reportOverview?.serviceStatus ?? [])) }"></span>
                </div>
              </div>
            </div>

            <div>
              <div class="ratio-title">支付渠道</div>
              <div
                v-for="item in reportOverview?.paymentChannels ?? []"
                :key="`channel-${item.name}`"
                class="ratio-item"
              >
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} / {{ formatCurrency(item.amount) }}</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--violet">
                  <span :style="{ width: ratioWidth(item, totalCount(reportOverview?.paymentChannels ?? [])) }"></span>
                </div>
              </div>
            </div>

            <div>
              <div class="ratio-title">计费周期</div>
              <div v-for="item in reportOverview?.billingCycles ?? []" :key="`cycle-${item.name}`" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} / {{ formatCurrency(item.amount) }}</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--orange">
                  <span :style="{ width: ratioWidth(item, totalCount(reportOverview?.billingCycles ?? [])) }"></span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="sub-grid">
        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">热销商品</h2>
              <p class="page-subtitle">从报表直接跳商品中心和订单中心，继续做商品运营与价格复盘。</p>
            </div>
          </div>

          <el-table :data="reportOverview?.topProducts ?? []" border stripe empty-text="暂无商品排行">
            <el-table-column type="index" label="#" width="56" />
            <el-table-column prop="name" label="商品名称" min-width="220" />
            <el-table-column prop="count" label="订单数" min-width="100" />
            <el-table-column label="累计金额" min-width="140">
              <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
            </el-table-column>
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="goProductOrders(row.name)">查看订单</el-button>
                <el-button link type="primary" @click="goProductCatalog(row.name)">商品中心</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">客户应收排行</h2>
              <p class="page-subtitle">快速定位应收压力集中的客户，再进入客户或账单工作台处理。</p>
            </div>
          </div>

          <el-table :data="reportOverview?.topReceivables ?? []" border stripe empty-text="暂无应收排行">
            <el-table-column type="index" label="#" width="56" />
            <el-table-column prop="name" label="客户名称" min-width="220" />
            <el-table-column prop="count" label="未付账单数" min-width="120" />
            <el-table-column label="应收金额" min-width="140">
              <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="extra" label="最近到期时间" min-width="180" />
            <el-table-column label="操作" min-width="180" fixed="right">
              <template #default="{ row }">
                <el-button
                  link
                  type="primary"
                  @click="goCustomerList({ keyword: row.name, status: 'ACTIVE' })"
                >
                  查看客户
                </el-button>
                <el-button link type="primary" @click="go('/billing/invoices?status=UNPAID')">
                  账单工作台
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <div class="sub-grid">
        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">客户分组结构</h2>
              <p class="page-subtitle">查看分组覆盖情况，并直接跳到客户分组或客户列表继续处理。</p>
            </div>
            <div class="action-group">
              <el-button plain @click="go('/customer/groups')">客户分组</el-button>
              <el-button plain @click="go('/customer/levels')">客户等级</el-button>
            </div>
          </div>

          <div v-if="(reportOverview?.customerGroups ?? []).length" class="ratio-panel">
            <div
              v-for="item in reportOverview?.customerGroups ?? []"
              :key="`customer-group-${item.name}`"
              class="ratio-item"
            >
              <div class="ratio-item__meta">
                <el-button
                  link
                  type="primary"
                  class="inline-link"
                  @click="goCustomerList({ groupName: item.name })"
                >
                  {{ item.name }}
                </el-button>
                <span>{{ item.count }} 位客户</span>
              </div>
              <div class="ratio-item__bar ratio-item__bar--success">
                <span :style="{ width: ratioWidth(item, customerGroupTotal) }"></span>
              </div>
            </div>
          </div>
          <el-empty v-else description="暂无客户分组结构数据" />
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">自动化异常</h2>
              <p class="page-subtitle">把失败和阻塞任务拉到报表中心，避免只能在任务页里被动发现。</p>
            </div>
          </div>

          <el-table :data="automationAlerts" border stripe empty-text="暂无自动化异常">
            <el-table-column prop="createdAt" label="时间" min-width="170" />
            <el-table-column prop="title" label="任务标题" min-width="220" />
            <el-table-column prop="channel" label="渠道" min-width="120" />
            <el-table-column label="状态" min-width="110">
              <template #default="{ row }">
                <el-tag :type="statusTagType(row.status)" effect="light">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="message" label="结果说明" min-width="240" />
            <el-table-column label="操作" min-width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="goAutomationTask(row)">查看任务</el-button>
                <el-button
                  v-if="row.serviceId"
                  link
                  type="primary"
                  @click="goServiceDetail(row.serviceId)"
                >
                  查看服务
                </el-button>
                <el-button
                  v-else-if="row.invoiceId"
                  link
                  type="primary"
                  @click="goInvoiceDetail(row.invoiceId)"
                >
                  查看账单
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </div>

      <div class="page-card section-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">运营待办</h2>
            <p class="page-subtitle">把待实名、逾期账单、即将到期服务和待跟进工单统一收口到这里。</p>
          </div>
        </div>

        <el-tabs v-model="activeTodoTab">
          <el-tab-pane :label="`待实名审核 (${operationsWorkbench?.pendingIdentities.length ?? 0})`" name="identity">
            <el-table :data="operationsWorkbench?.pendingIdentities ?? []" border stripe empty-text="暂无待实名审核">
              <el-table-column prop="customerNo" label="客户编号" min-width="140" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column prop="subjectName" label="实名主体" min-width="220" />
              <el-table-column prop="submittedAt" label="提交时间" min-width="180" />
              <el-table-column label="操作" min-width="120" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="goCustomerDetail(row.customerId)">
                    查看客户
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane :label="`逾期账单 (${operationsWorkbench?.overdueInvoices.length ?? 0})`" name="invoice">
            <el-table :data="operationsWorkbench?.overdueInvoices ?? []" border stripe empty-text="暂无逾期账单">
              <el-table-column prop="invoiceNo" label="账单编号" min-width="160" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column label="应收金额" min-width="140">
                <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="dueAt" label="到期时间" min-width="180" />
              <el-table-column prop="daysOverdue" label="逾期天数" min-width="120" />
              <el-table-column label="操作" min-width="120" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="goInvoiceDetail(row.invoiceId)">
                    查看账单
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane :label="`即将到期服务 (${operationsWorkbench?.expiringServices.length ?? 0})`" name="service">
            <el-table :data="operationsWorkbench?.expiringServices ?? []" border stripe empty-text="暂无即将到期服务">
              <el-table-column prop="serviceNo" label="服务编号" min-width="160" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column prop="productName" label="商品名称" min-width="220" />
              <el-table-column prop="status" label="状态" min-width="120" />
              <el-table-column prop="nextDueAt" label="到期时间" min-width="180" />
              <el-table-column prop="daysRemaining" label="剩余天数" min-width="120" />
              <el-table-column label="操作" min-width="120" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="goServiceDetail(row.serviceId)">
                    查看服务
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane :label="`待跟进工单 (${operationsWorkbench?.openTickets.length ?? 0})`" name="ticket">
            <el-table :data="operationsWorkbench?.openTickets ?? []" border stripe empty-text="暂无待跟进工单">
              <el-table-column prop="ticketNo" label="工单编号" min-width="160" />
              <el-table-column prop="customerName" label="客户名称" min-width="180" />
              <el-table-column prop="title" label="工单标题" min-width="260" />
              <el-table-column prop="status" label="状态" min-width="120" />
              <el-table-column prop="updatedAt" label="最后更新时间" min-width="180" />
              <el-table-column label="操作" min-width="120" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="goTicketDetail(row.ticketId)">
                    查看工单
                  </el-button>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>
        </el-tabs>
      </div>

      <div class="page-card section-card">
        <div class="page-header">
          <div>
            <h2 class="section-title">最近系统动作</h2>
            <p class="page-subtitle">直接查看审计日志，快速回溯后台最近做过的关键操作。</p>
          </div>
        </div>

        <el-table :data="operationsWorkbench?.recentAudits ?? []" border stripe empty-text="暂无系统动作">
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
.report-page {
  display: grid;
  gap: 16px;
}

.headline-grid,
.quick-grid,
.sub-grid,
.trend-layout {
  display: grid;
  gap: 14px;
}

.headline-grid,
.quick-grid {
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
}

.sub-grid,
.trend-layout {
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
}

.headline-card,
.quick-card {
  border-radius: 16px;
  border: 1px solid #e6edf7;
  background: #fff;
  padding: 16px;
}

.quick-card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  width: 100%;
  cursor: pointer;
  text-align: left;
}

.headline-card__label,
.quick-card__label {
  color: #5e7093;
  font-size: 13px;
}

.headline-card__value,
.quick-card__value {
  font-size: 28px;
  font-weight: 700;
  color: #16376f;
}

.headline-card__hint,
.quick-card__hint {
  color: #7083a6;
  font-size: 12px;
  line-height: 1.6;
}

.section-card {
  padding: 20px;
}

.trend-list {
  display: grid;
  gap: 10px;
}

.trend-item__meta,
.ratio-item__meta {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
  color: #526684;
  font-size: 13px;
}

.trend-item__bar,
.ratio-item__bar {
  height: 10px;
  border-radius: 999px;
  background: #edf3fb;
  overflow: hidden;
}

.trend-item__fill,
.ratio-item__bar span {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: linear-gradient(90deg, #2563eb 0%, #60a5fa 100%);
}

.trend-item__bar--warning .trend-item__fill--warning {
  background: linear-gradient(90deg, #f97316 0%, #fdba74 100%);
}

.ratio-panel {
  display: grid;
  gap: 18px;
}

.ratio-title,
.block-title {
  margin-bottom: 10px;
  font-weight: 700;
  color: #173b72;
}

.ratio-item {
  margin-bottom: 12px;
}

.ratio-item__bar--success span {
  background: linear-gradient(90deg, #0f8a54 0%, #34d399 100%);
}

.ratio-item__bar--violet span {
  background: linear-gradient(90deg, #7c3aed 0%, #a78bfa 100%);
}

.ratio-item__bar--orange span {
  background: linear-gradient(90deg, #ea580c 0%, #fb923c 100%);
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

.inline-link {
  padding: 0;
}
</style>
