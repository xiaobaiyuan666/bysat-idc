<script setup lang="ts">
import { onMounted, ref } from "vue";

import { http } from "@/api/http";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  getLabel,
  getStatusTagType,
  invoiceStatusMap,
  paymentMethodMap,
} from "@/utils/maps";

const loading = ref(false);
const dataSource = ref<any>({
  summary: {
    grossReceipts: 0,
    totalFees: 0,
    totalRefunded: 0,
    netReceipts: 0,
    outstandingAmount: 0,
    openInvoiceCount: 0,
    overdueInvoiceCount: 0,
    pendingCallbackCount: 0,
    failedCallbackCount: 0,
    pendingRefundCount: 0,
    issueCount: 0,
    negativeBalanceCustomerCount: 0,
  },
  methodSummary: [],
  issues: [],
  openInvoices: [],
  trend: [],
  negativeBalanceCustomers: [],
});

function getIssueTagType(level: string) {
  if (level === "critical") {
    return "danger";
  }

  if (level === "warning") {
    return "warning";
  }

  return "info";
}

function getIssueLevelLabel(level: string) {
  if (level === "critical") {
    return "严重";
  }

  if (level === "warning") {
    return "提醒";
  }

  return "一般";
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/payments/reconciliation");
    dataSource.value = data.data;
  } finally {
    loading.value = false;
  }
}

function exportIssues() {
  window.open("/api/payments/reconciliation/export", "_blank");
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">财务对账</h1>
        <div class="page-subtitle">
          汇总收款、退款、回调、未结清账单和余额异常，形成统一的财务核对视图。
        </div>
      </div>
      <div class="toolbar-actions">
        <el-button @click="loadData">刷新数据</el-button>
        <el-button @click="exportIssues">导出异常</el-button>
      </div>
    </div>

    <div class="metric-grid">
      <div class="metric-tile">
        <div class="metric-label">累计实收</div>
        <div class="metric-value">{{ formatCurrency(dataSource.summary.grossReceipts) }}</div>
        <div class="metric-hint">成功入账和已退款收款单的原始收款总额。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">累计退款</div>
        <div class="metric-value">{{ formatCurrency(dataSource.summary.totalRefunded) }}</div>
        <div class="metric-hint">所有成功退款的累计金额。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">净入账</div>
        <div class="metric-value">{{ formatCurrency(dataSource.summary.netReceipts) }}</div>
        <div class="metric-hint">收款减退款再减通道手续费后的净额。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">待核销应收</div>
        <div class="metric-value">{{ formatCurrency(dataSource.summary.outstandingAmount) }}</div>
        <div class="metric-hint">全部未结清账单的剩余待收金额。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">待处理回调</div>
        <div class="metric-value">{{ dataSource.summary.pendingCallbackCount }}</div>
        <div class="metric-hint">尚未完成入账确认的支付回调数量。</div>
      </div>
      <div class="metric-tile">
        <div class="metric-label">异常项</div>
        <div class="metric-value">{{ dataSource.summary.issueCount }}</div>
        <div class="metric-hint">账单、收款、回调和余额层面的异常总数。</div>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :xl="16" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>渠道汇总</span>
          </template>
          <el-table v-loading="loading" :data="dataSource.methodSummary" stripe>
            <el-table-column label="支付方式" min-width="140">
              <template #default="{ row }">
                {{ getLabel(paymentMethodMap, row.method) }}
              </template>
            </el-table-column>
            <el-table-column label="收款笔数" width="100" prop="paymentCount" />
            <el-table-column label="收款金额" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.paymentAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="退款金额" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.refundAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="手续费" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.feeAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="净额" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.netAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="待处理回调" width="110" prop="pendingCallbacks" />
            <el-table-column label="失败回调" width="100" prop="failedCallbacks" />
            <el-table-column label="最近收款" width="170">
              <template #default="{ row }">
                {{ formatDateTime(row.lastPaymentAt) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xl="8" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>负余额客户</span>
          </template>
          <div class="side-list">
            <div
              v-for="item in dataSource.negativeBalanceCustomers"
              :key="item.id"
              class="side-list-item"
            >
              <div>
                <div class="side-title">{{ item.name }}</div>
                <div class="side-subtitle">{{ item.customerNo }}</div>
              </div>
              <div class="side-amount danger">{{ formatCurrency(item.creditBalance) }}</div>
            </div>
            <el-empty
              v-if="!dataSource.negativeBalanceCustomers.length"
              description="当前没有负余额客户"
            />
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :xl="15" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>对账异常</span>
          </template>
          <el-table v-loading="loading" :data="dataSource.issues" stripe>
            <el-table-column label="级别" width="100">
              <template #default="{ row }">
                <el-tag :type="getIssueTagType(row.level)">
                  {{ getIssueLevelLabel(row.level) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="kind" label="类型" width="110" />
            <el-table-column prop="title" label="异常标题" min-width="190" />
            <el-table-column prop="relatedNo" label="相关单号" min-width="180" />
            <el-table-column label="客户" min-width="140">
              <template #default="{ row }">
                {{ row.customerName || "-" }}
              </template>
            </el-table-column>
            <el-table-column label="涉及金额" width="140">
              <template #default="{ row }">
                {{ formatCurrency(row.amount) }}
              </template>
            </el-table-column>
            <el-table-column prop="detail" label="异常说明" min-width="280" show-overflow-tooltip />
            <el-table-column label="发生时间" width="170">
              <template #default="{ row }">
                {{ formatDateTime(row.createdAt) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>

      <el-col :xl="9" :span="24">
        <el-card class="page-card">
          <template #header>
            <span>近 7 天资金变化</span>
          </template>
          <el-table v-loading="loading" :data="dataSource.trend" stripe>
            <el-table-column prop="date" label="日期" width="100" />
            <el-table-column label="收款" width="120">
              <template #default="{ row }">
                {{ formatCurrency(row.paymentAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="退款" width="120">
              <template #default="{ row }">
                {{ formatCurrency(row.refundAmount) }}
              </template>
            </el-table-column>
            <el-table-column label="净额" width="120">
              <template #default="{ row }">
                {{ formatCurrency(row.netAmount) }}
              </template>
            </el-table-column>
            <el-table-column prop="paymentCount" label="收款笔数" width="90" />
            <el-table-column prop="refundCount" label="退款笔数" width="90" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="page-card">
      <template #header>
        <span>未结清账单</span>
      </template>
      <el-table v-loading="loading" :data="dataSource.openInvoices" stripe>
        <el-table-column prop="invoiceNo" label="账单号" min-width="180" />
        <el-table-column label="客户" min-width="160">
          <template #default="{ row }">
            {{ row.customer.name }}
          </template>
        </el-table-column>
        <el-table-column label="账单金额" width="140">
          <template #default="{ row }">
            {{ formatCurrency(row.totalAmount) }}
          </template>
        </el-table-column>
        <el-table-column label="已收金额" width="140">
          <template #default="{ row }">
            {{ formatCurrency(row.paidAmount) }}
          </template>
        </el-table-column>
        <el-table-column label="待收金额" width="140">
          <template #default="{ row }">
            {{ formatCurrency(row.outstandingAmount) }}
          </template>
        </el-table-column>
        <el-table-column label="到期日" width="120">
          <template #default="{ row }">
            {{ formatDate(row.dueDate) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="110">
          <template #default="{ row }">
            <el-tag :type="getStatusTagType(row.status)">
              {{ getLabel(invoiceStatusMap, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.toolbar-actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.side-list {
  display: grid;
  gap: 12px;
}

.side-list-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 16px;
  background: var(--el-fill-color-light);
}

.side-title {
  font-weight: 700;
}

.side-subtitle {
  margin-top: 4px;
  color: var(--text-secondary);
  font-size: 12px;
}

.side-amount {
  font-weight: 700;
}

.danger {
  color: var(--el-color-danger);
}
</style>
