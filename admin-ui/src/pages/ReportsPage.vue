<script setup lang="ts">
import { computed, onMounted, ref } from "vue";

import { http } from "@/api/http";
import { formatCurrency, formatDate } from "@/utils/format";
import { getStatusTagType } from "@/utils/maps";

type ReportsOverview = {
  summary: {
    monthlyRevenue: number;
    monthlyRefunded: number;
    netRevenue: number;
    outstandingAmount: number;
    overdueServiceCount: number;
    activeServiceCount: number;
  };
  dailyRevenue: Array<{
    date: string;
    revenue: number;
    refunds: number;
  }>;
  paymentChannels: Array<{
    method: string;
    amount: number;
    count: number;
  }>;
  invoiceStatusSummary: Array<{
    status: string;
    count: number;
    amount: number;
  }>;
  productRanking: Array<{
    productName: string;
    amount: number;
    count: number;
  }>;
  regionDistribution: Array<{
    region: string;
    count: number;
    amount: number;
  }>;
  overdueServices: Array<{
    id: string;
    serviceNo: string;
    name: string;
    customerName: string;
    productName: string;
    status: string;
    nextDueDate: string | Date | null;
    amount: number;
  }>;
  customerArRanking: Array<{
    id: string;
    name: string;
    outstandingAmount: number;
    openInvoiceCount: number;
  }>;
};

const paymentMethodLabels: Record<string, string> = {
  BALANCE: "余额支付",
  ALIPAY: "支付宝",
  WECHAT: "微信支付",
  BANK_TRANSFER: "银行转账",
  OFFLINE: "线下收款",
  PAYPAL: "PayPal",
  STRIPE: "Stripe",
  OTHER: "其他",
};

const invoiceStatusLabels: Record<string, string> = {
  DRAFT: "草稿",
  ISSUED: "待支付",
  PARTIAL: "部分支付",
  PAID: "已支付",
  OVERDUE: "已逾期",
  VOID: "已作废",
};

const serviceStatusLabels: Record<string, string> = {
  PENDING: "待开通",
  PROVISIONING: "开通中",
  ACTIVE: "运行中",
  SUSPENDED: "已暂停",
  OVERDUE: "已逾期",
  TERMINATED: "已终止",
  EXPIRED: "已到期",
  FAILED: "异常",
};

const loading = ref(false);
const dataSource = ref<ReportsOverview>({
  summary: {
    monthlyRevenue: 0,
    monthlyRefunded: 0,
    netRevenue: 0,
    outstandingAmount: 0,
    overdueServiceCount: 0,
    activeServiceCount: 0,
  },
  dailyRevenue: [],
  paymentChannels: [],
  invoiceStatusSummary: [],
  productRanking: [],
  regionDistribution: [],
  overdueServices: [],
  customerArRanking: [],
});

const trendMax = computed(() =>
  Math.max(
    1,
    ...dataSource.value.dailyRevenue.map((item) =>
      Math.max(item.revenue, item.refunds, item.revenue - item.refunds),
    ),
  ),
);

const productMax = computed(() =>
  Math.max(1, ...dataSource.value.productRanking.map((item) => item.amount)),
);

const regionMax = computed(() =>
  Math.max(1, ...dataSource.value.regionDistribution.map((item) => item.count)),
);

function getPaymentMethodLabel(method: string) {
  return paymentMethodLabels[method] ?? method;
}

function getInvoiceStatusLabel(status: string) {
  return invoiceStatusLabels[status] ?? status;
}

function getServiceStatusLabel(status: string) {
  return serviceStatusLabels[status] ?? status;
}

function getTrendWidth(amount: number) {
  return `${Math.max(8, Math.round((amount / trendMax.value) * 100))}%`;
}

function getProductWidth(amount: number) {
  return `${Math.max(10, Math.round((amount / productMax.value) * 100))}%`;
}

function getRegionWidth(count: number) {
  return `${Math.max(10, Math.round((count / regionMax.value) * 100))}%`;
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/reports/overview");
    dataSource.value = data.data;
  } finally {
    loading.value = false;
  }
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">统计报表</h1>
        <div class="page-subtitle">
          汇总近 30 天收入趋势、账单结构、支付渠道表现、地域资源分布和欠费风险，作为运营和财务的统一分析入口。
        </div>
      </div>
      <div class="toolbar-actions">
        <el-button @click="loadData">刷新报表</el-button>
      </div>
    </div>

    <el-skeleton :loading="loading" animated :rows="10">
      <div class="metric-grid">
        <div class="metric-tile">
          <div class="metric-label">本月实收</div>
          <div class="metric-value">{{ formatCurrency(dataSource.summary.monthlyRevenue) }}</div>
          <div class="metric-hint">按成功入账收款统计的本月累计收入。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">本月退款</div>
          <div class="metric-value">{{ formatCurrency(dataSource.summary.monthlyRefunded) }}</div>
          <div class="metric-hint">已成功处理的退款总额。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">净收入</div>
          <div class="metric-value">{{ formatCurrency(dataSource.summary.netRevenue) }}</div>
          <div class="metric-hint">本月实收减去退款后的净额。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">当前应收</div>
          <div class="metric-value">{{ formatCurrency(dataSource.summary.outstandingAmount) }}</div>
          <div class="metric-hint">全部未结清账单的待收金额。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">活跃服务</div>
          <div class="metric-value">{{ dataSource.summary.activeServiceCount }}</div>
          <div class="metric-hint">当前处于活动状态的服务实例数。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">欠费风险服务</div>
          <div class="metric-value">{{ dataSource.summary.overdueServiceCount }}</div>
          <div class="metric-hint">已逾期、暂停或到期服务的预警数量。</div>
        </div>
      </div>

      <el-row :gutter="16">
        <el-col :xl="14" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>近 14 天收入趋势</h2>
                <el-tag type="info" effect="plain">收入与退款对照</el-tag>
              </div>
            </template>
            <div class="trend-list">
              <div v-for="item in dataSource.dailyRevenue" :key="item.date" class="trend-item">
                <div class="trend-date">{{ item.date }}</div>
                <div class="trend-bars">
                  <div class="trend-bar revenue" :style="{ width: getTrendWidth(item.revenue) }" />
                  <div class="trend-bar refund" :style="{ width: getTrendWidth(item.refunds) }" />
                </div>
                <div class="trend-amounts">
                  <span>收 {{ formatCurrency(item.revenue) }}</span>
                  <span>退 {{ formatCurrency(item.refunds) }}</span>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xl="10" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>支付渠道表现</h2>
                <el-tag type="success" effect="plain">按实收金额排序</el-tag>
              </div>
            </template>
            <el-table :data="dataSource.paymentChannels" stripe>
              <el-table-column label="渠道" min-width="130">
                <template #default="{ row }">
                  {{ getPaymentMethodLabel(row.method) }}
                </template>
              </el-table-column>
              <el-table-column prop="count" label="笔数" width="90" />
              <el-table-column label="金额" width="140">
                <template #default="{ row }">
                  {{ formatCurrency(row.amount) }}
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :xl="8" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>账单状态结构</h2>
                <el-tag type="warning" effect="plain">账单全量统计</el-tag>
              </div>
            </template>
            <div class="status-summary-list">
              <div
                v-for="item in dataSource.invoiceStatusSummary"
                :key="item.status"
                class="status-summary-item"
              >
                <div class="status-summary-header">
                  <el-tag :type="getStatusTagType(item.status)">
                    {{ getInvoiceStatusLabel(item.status) }}
                  </el-tag>
                  <strong>{{ item.count }}</strong>
                </div>
                <div class="muted-line">{{ formatCurrency(item.amount) }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xl="8" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>热销产品排行</h2>
                <el-tag type="info" effect="plain">按收入贡献排序</el-tag>
              </div>
            </template>
            <div class="ranking-list">
              <div v-for="item in dataSource.productRanking" :key="item.productName" class="ranking-item">
                <div class="ranking-meta">
                  <div class="primary-line">{{ item.productName }}</div>
                  <div class="muted-line">{{ item.count }} 笔订单收入</div>
                </div>
                <div class="ranking-bar-track">
                  <div class="ranking-bar" :style="{ width: getProductWidth(item.amount) }" />
                </div>
                <div class="ranking-amount">{{ formatCurrency(item.amount) }}</div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xl="8" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>区域资源分布</h2>
                <el-tag type="success" effect="plain">按服务数量统计</el-tag>
              </div>
            </template>
            <div class="ranking-list">
              <div v-for="item in dataSource.regionDistribution" :key="item.region" class="ranking-item">
                <div class="ranking-meta">
                  <div class="primary-line">{{ item.region }}</div>
                  <div class="muted-line">{{ formatCurrency(item.amount) }} 月费规模</div>
                </div>
                <div class="ranking-bar-track green">
                  <div class="ranking-bar green" :style="{ width: getRegionWidth(item.count) }" />
                </div>
                <div class="ranking-amount">{{ item.count }} 台</div>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :xl="14" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>欠费风险服务</h2>
                <el-tag type="danger" effect="plain">优先跟进</el-tag>
              </div>
            </template>
            <el-table :data="dataSource.overdueServices" stripe>
              <el-table-column prop="serviceNo" label="服务编号" min-width="150" />
              <el-table-column prop="name" label="实例名称" min-width="170" />
              <el-table-column prop="customerName" label="客户" min-width="150" />
              <el-table-column prop="productName" label="产品" min-width="140" />
              <el-table-column label="月费" width="120">
                <template #default="{ row }">
                  {{ formatCurrency(row.amount) }}
                </template>
              </el-table-column>
              <el-table-column label="到期日" width="120">
                <template #default="{ row }">
                  {{ formatDate(row.nextDueDate) }}
                </template>
              </el-table-column>
              <el-table-column label="状态" width="110">
                <template #default="{ row }">
                  <el-tag :type="getStatusTagType(row.status)">
                    {{ getServiceStatusLabel(row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <el-col :xl="10" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="section-heading">
                <h2>客户应收排行</h2>
                <el-tag type="warning" effect="plain">按未结清金额排序</el-tag>
              </div>
            </template>
            <div class="ar-list">
              <div v-for="item in dataSource.customerArRanking" :key="item.id" class="ar-item">
                <div>
                  <div class="primary-line">{{ item.name }}</div>
                  <div class="muted-line">{{ item.openInvoiceCount }} 张未结清账单</div>
                </div>
                <div class="ar-amount">{{ formatCurrency(item.outstandingAmount) }}</div>
              </div>
              <el-empty v-if="!dataSource.customerArRanking.length" description="当前没有应收风险客户" />
            </div>
          </el-card>
        </el-col>
      </el-row>
    </el-skeleton>
  </div>
</template>

<style scoped>
.toolbar-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.trend-list,
.ranking-list,
.ar-list,
.status-summary-list {
  display: grid;
  gap: 14px;
}

.trend-item {
  display: grid;
  grid-template-columns: 76px minmax(0, 1fr) 220px;
  gap: 12px;
  align-items: center;
}

.trend-date {
  font-weight: 700;
  color: var(--text-secondary);
}

.trend-bars {
  display: grid;
  gap: 8px;
}

.trend-bar {
  height: 10px;
  border-radius: 999px;
}

.trend-bar.revenue {
  background: linear-gradient(90deg, rgba(36, 104, 242, 0.85), rgba(47, 128, 255, 0.55));
}

.trend-bar.refund {
  background: linear-gradient(90deg, rgba(240, 68, 56, 0.82), rgba(247, 144, 9, 0.58));
}

.trend-amounts {
  display: flex;
  justify-content: space-between;
  gap: 12px;
  color: var(--text-secondary);
  font-size: 13px;
}

.status-summary-item {
  padding: 14px 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: linear-gradient(180deg, #fff 0%, #f8fbff 100%);
}

.status-summary-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.status-summary-header strong {
  font-size: 22px;
}

.ranking-item,
.ar-item {
  padding: 14px 16px;
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: #fbfcff;
}

.ranking-meta {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
}

.ranking-bar-track {
  margin-top: 10px;
  height: 10px;
  border-radius: 999px;
  background: rgba(36, 104, 242, 0.12);
  overflow: hidden;
}

.ranking-bar-track.green {
  background: rgba(23, 178, 106, 0.14);
}

.ranking-bar {
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, #2468f2 0%, #58a5ff 100%);
}

.ranking-bar.green {
  background: linear-gradient(90deg, #17b26a 0%, #3ccf91 100%);
}

.ranking-amount,
.ar-amount {
  margin-top: 10px;
  font-weight: 700;
}

.ar-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

@media (max-width: 1200px) {
  .trend-item {
    grid-template-columns: 1fr;
  }

  .trend-amounts {
    justify-content: flex-start;
  }
}
</style>
