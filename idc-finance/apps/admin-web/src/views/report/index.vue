<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchReportOverview,
  type ReportOverviewResponse,
  type StatusBucket,
  type TrendPoint
} from "@/api/admin";

const loading = ref(false);
const overview = ref<ReportOverviewResponse | null>(null);

const revenueMax = computed(() => {
  if (!overview.value?.revenueTrend.length && !overview.value?.refundTrend.length) return 1;
  return Math.max(
    ...[...(overview.value?.revenueTrend ?? []), ...(overview.value?.refundTrend ?? [])].map(item => item.amount),
    1
  );
});

function formatCurrency(value: number) {
  return `¥${value.toFixed(2)}`;
}

function barWidth(item: TrendPoint) {
  return `${Math.max((item.amount / revenueMax.value) * 100, item.amount > 0 ? 10 : 4)}%`;
}

function ratioWidth(item: StatusBucket, total: number) {
  if (!total) return "0%";
  return `${Math.max((item.count / total) * 100, item.count > 0 ? 8 : 0)}%`;
}

function totalCount(items: StatusBucket[]) {
  return items.reduce((total, item) => total + item.count, 0);
}

async function loadOverview() {
  loading.value = true;
  try {
    overview.value = await fetchReportOverview();
  } finally {
    loading.value = false;
  }
}

onMounted(loadOverview);
</script>

<template>
  <div v-loading="loading" class="report-page">
    <PageWorkbench
      eyebrow="功能 / 报表"
      title="统计报表"
      subtitle="聚合收款、退款、账单状态、服务结构、支付渠道、账期偏好和客户应收，用于运营与财务复盘。"
    >
      <template #actions>
        <el-button type="primary" @click="loadOverview">刷新报表</el-button>
      </template>

      <template #metrics>
        <div class="headline-grid">
          <div v-for="item in overview?.headline ?? []" :key="item.key" class="headline-card">
            <div class="headline-card__label">{{ item.label }}</div>
            <div class="headline-card__value">{{ item.value }}</div>
            <div class="headline-card__hint">{{ item.hint }}</div>
          </div>
        </div>
      </template>

      <div class="sub-grid">
        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">近 30 日收款趋势</h2>
              <p class="page-subtitle">按支付流水聚合金额与笔数，反映收款节奏和峰谷。</p>
            </div>
          </div>
          <div class="trend-list">
            <div v-for="item in overview?.revenueTrend ?? []" :key="item.label" class="trend-item">
              <div class="trend-item__meta">
                <strong>{{ item.label }}</strong>
                <span>{{ formatCurrency(item.amount) }} / {{ item.count }} 笔</span>
              </div>
              <div class="trend-item__bar">
                <span class="trend-item__fill" :style="{ width: barWidth(item) }"></span>
              </div>
            </div>
          </div>
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">近 30 日退款趋势</h2>
              <p class="page-subtitle">按退款单聚合金额与笔数，辅助财务观察退款压力。</p>
            </div>
          </div>
          <div class="trend-list">
            <div v-for="item in overview?.refundTrend ?? []" :key="item.label" class="trend-item">
              <div class="trend-item__meta">
                <strong>{{ item.label }}</strong>
                <span>{{ formatCurrency(item.amount) }} / {{ item.count }} 笔</span>
              </div>
              <div class="trend-item__bar trend-item__bar--warning">
                <span class="trend-item__fill trend-item__fill--warning" :style="{ width: barWidth(item) }"></span>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="sub-grid">
        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">账单与服务结构</h2>
              <p class="page-subtitle">观察账单状态与服务生命周期的总体分布。</p>
            </div>
          </div>
          <div class="ratio-panel">
            <div>
              <div class="ratio-title">账单状态</div>
              <div v-for="item in overview?.invoiceStatus ?? []" :key="item.name" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} / {{ formatCurrency(item.amount) }}</span>
                </div>
                <div class="ratio-item__bar">
                  <span :style="{ width: ratioWidth(item, totalCount(overview?.invoiceStatus ?? [])) }"></span>
                </div>
              </div>
            </div>

            <div>
              <div class="ratio-title">服务状态</div>
              <div v-for="item in overview?.serviceStatus ?? []" :key="item.name" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} 台</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--success">
                  <span :style="{ width: ratioWidth(item, totalCount(overview?.serviceStatus ?? [])) }"></span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">渠道与客户结构</h2>
              <p class="page-subtitle">看支付渠道、账期偏好和客户分层，便于配置产品与催收策略。</p>
            </div>
          </div>
          <div class="ratio-panel">
            <div>
              <div class="ratio-title">支付渠道</div>
              <div v-for="item in overview?.paymentChannels ?? []" :key="item.name" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} / {{ formatCurrency(item.amount) }}</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--violet">
                  <span :style="{ width: ratioWidth(item, totalCount(overview?.paymentChannels ?? [])) }"></span>
                </div>
              </div>
            </div>

            <div>
              <div class="ratio-title">账期结构</div>
              <div v-for="item in overview?.billingCycles ?? []" :key="item.name" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} / {{ formatCurrency(item.amount) }}</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--orange">
                  <span :style="{ width: ratioWidth(item, totalCount(overview?.billingCycles ?? [])) }"></span>
                </div>
              </div>
            </div>

            <div>
              <div class="ratio-title">客户分组</div>
              <div v-for="item in overview?.customerGroups ?? []" :key="item.name" class="ratio-item">
                <div class="ratio-item__meta">
                  <span>{{ item.name }}</span>
                  <span>{{ item.count }} 客户</span>
                </div>
                <div class="ratio-item__bar ratio-item__bar--teal">
                  <span :style="{ width: ratioWidth(item, totalCount(overview?.customerGroups ?? [])) }"></span>
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
              <h2 class="section-title">热销产品</h2>
              <p class="page-subtitle">按订单金额聚合的产品排行，可用于销售和活动判断。</p>
            </div>
          </div>
          <el-table :data="overview?.topProducts ?? []" border stripe empty-text="暂无产品排行">
            <el-table-column type="index" label="#" width="60" />
            <el-table-column prop="name" label="产品名称" min-width="220" />
            <el-table-column prop="count" label="订单数" min-width="120" />
            <el-table-column label="累计金额" min-width="160">
              <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
            </el-table-column>
          </el-table>
        </div>

        <div class="page-card section-card">
          <div class="page-header">
            <div>
              <h2 class="section-title">客户应收排行</h2>
              <p class="page-subtitle">按未支付账单聚合，方便财务催收和客服跟进。</p>
            </div>
          </div>
          <el-table :data="overview?.topReceivables ?? []" border stripe empty-text="暂无应收排行">
            <el-table-column type="index" label="#" width="60" />
            <el-table-column prop="name" label="客户名称" min-width="220" />
            <el-table-column prop="count" label="未付账单数" min-width="120" />
            <el-table-column label="应收金额" min-width="160">
              <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="extra" label="最近到期时间" min-width="180" />
          </el-table>
        </div>
      </div>
    </PageWorkbench>
  </div>
</template>

<style scoped>
.report-page {
  display: grid;
  gap: 16px;
}

.headline-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 14px;
}

.headline-card {
  border-radius: 16px;
  border: 1px solid #e6edf7;
  background: #fff;
  padding: 16px;
}

.headline-card__label {
  color: #5e7093;
  font-size: 13px;
}

.headline-card__value {
  margin-top: 10px;
  font-size: 28px;
  font-weight: 700;
  color: #16376f;
}

.headline-card__hint {
  margin-top: 10px;
  color: #7083a6;
  font-size: 12px;
  line-height: 1.6;
}

.section-card {
  padding: 20px;
}

.trend-list {
  display: grid;
  gap: 12px;
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
  gap: 20px;
}

.ratio-title {
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

.ratio-item__bar--teal span {
  background: linear-gradient(90deg, #0f766e 0%, #2dd4bf 100%);
}
</style>
