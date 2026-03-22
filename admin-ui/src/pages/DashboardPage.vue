<script setup lang="ts">
import { onMounted, ref } from "vue";

import { http } from "@/api/http";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  billingJobStatusMap,
  billingJobTypeMap,
  getLabel,
  getStatusTagType,
  providerActionMap,
  providerSyncStatusMap,
  serviceStatusMap,
} from "@/utils/maps";

const loading = ref(false);
const dashboard = ref<any>({
  metrics: {
    monthlyRevenue: 0,
    outstandingAmount: 0,
    activeServices: 0,
    overdueCustomers: 0,
    totalBalance: 0,
    openRenewalInvoices: 0,
  },
  topCustomers: [],
  renewals: [],
  billingJobs: [],
  providerLogs: [],
  resourceSummary: {
    vpcs: 0,
    ips: 0,
    disks: 0,
    snapshots: 0,
    backups: 0,
    securityGroups: 0,
  },
});

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/dashboard");
    dashboard.value = data.data;
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
        <h1 class="page-title">运营工作台</h1>
        <div class="page-subtitle">
          首页不再只是统计卡片，而是面向运营和财务的统一工作台，集中展示收入、待收款、重点客户、续费队列、计费任务和云平台同步结果。
        </div>
      </div>
      <el-button type="primary" @click="loadData">刷新数据</el-button>
    </div>

    <el-skeleton :loading="loading" animated :rows="12">
      <div class="metric-grid">
        <div class="metric-tile">
          <div class="metric-label">本月实收</div>
          <div class="metric-value">{{ formatCurrency(dashboard.metrics.monthlyRevenue) }}</div>
          <div class="metric-hint">按成功入账的收款记录统计本月已结算金额。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">当前应收</div>
          <div class="metric-value">{{ formatCurrency(dashboard.metrics.outstandingAmount) }}</div>
          <div class="metric-hint">包含待支付、部分支付和已逾期账单。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">运行中服务</div>
          <div class="metric-value">{{ dashboard.metrics.activeServices }}</div>
          <div class="metric-hint">当前处于运行或开通中的服务实例数量。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">逾期客户</div>
          <div class="metric-value">{{ dashboard.metrics.overdueCustomers }}</div>
          <div class="metric-hint">需要财务或客服跟进处理的客户数量。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">客户余额池</div>
          <div class="metric-value">{{ formatCurrency(dashboard.metrics.totalBalance) }}</div>
          <div class="metric-hint">可用于自动续费和手工扣费的客户余额总量。</div>
        </div>
        <div class="metric-tile">
          <div class="metric-label">待处理续费账单</div>
          <div class="metric-value">{{ dashboard.metrics.openRenewalInvoices }}</div>
          <div class="metric-hint">已生成但尚未完成支付的续费账单数量。</div>
        </div>
      </div>

      <el-row :gutter="16">
        <el-col :xl="14" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="toolbar-row" style="margin-bottom: 0">
                <span>重点客户</span>
                <el-tag type="info" effect="plain">按累计开票金额排序</el-tag>
              </div>
            </template>
            <el-table :data="dashboard.topCustomers" stripe>
              <el-table-column prop="name" label="客户" min-width="180" />
              <el-table-column prop="serviceCount" label="服务数" width="100" />
              <el-table-column label="累计开票" width="140">
                <template #default="{ row }">
                  {{ formatCurrency(row.billedAmount) }}
                </template>
              </el-table-column>
              <el-table-column label="当前应收" width="140">
                <template #default="{ row }">
                  {{ formatCurrency(row.outstanding) }}
                </template>
              </el-table-column>
              <el-table-column label="余额" width="140">
                <template #default="{ row }">
                  {{ formatCurrency(row.creditBalance) }}
                </template>
              </el-table-column>
            </el-table>
          </el-card>
        </el-col>

        <el-col :xl="10" :span="24">
          <el-card class="page-card">
            <template #header>
              <div class="toolbar-row" style="margin-bottom: 0">
                <span>近期续费队列</span>
                <el-tag type="warning" effect="light">生命周期重点关注</el-tag>
              </div>
            </template>
            <div class="list-stack">
              <div v-for="service in dashboard.renewals" :key="service.id" class="mini-card">
                <div class="mini-title">{{ service.name }}</div>
                <div class="mini-copy">{{ service.customerName }} / {{ service.productName }}</div>
                <div class="mini-meta">
                  <span>到期 {{ formatDate(service.nextDueDate) }}</span>
                  <span>{{ formatCurrency(service.amount) }}</span>
                </div>
                <el-tag :type="getStatusTagType(service.status)" size="small">
                  {{ getLabel(serviceStatusMap, service.status) }}
                </el-tag>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-row :gutter="16">
        <el-col :xl="12" :span="24">
          <el-card class="page-card">
            <template #header>
              <span>计费任务执行记录</span>
            </template>
            <div class="list-stack">
              <div v-for="job in dashboard.billingJobs" :key="job.id" class="mini-card">
                <div class="mini-title">{{ getLabel(billingJobTypeMap, job.jobType) }}</div>
                <div class="mini-copy">
                  {{ job.customer?.name || "系统任务" }} / {{ job.service?.serviceNo || "无服务对象" }}
                </div>
                <div class="mini-copy">{{ job.message }}</div>
                <div class="mini-meta">
                  <span>{{ formatDateTime(job.executedAt) }}</span>
                  <el-tag :type="getStatusTagType(job.status)">
                    {{ getLabel(billingJobStatusMap, job.status) }}
                  </el-tag>
                </div>
              </div>
            </div>
          </el-card>
        </el-col>

        <el-col :xl="12" :span="24">
          <el-card class="page-card">
            <template #header>
              <span>资源库存概览</span>
            </template>
            <div class="summary-strip">
              <div class="summary-card">
                <div class="metric-label">VPC 网络</div>
                <h3>{{ dashboard.resourceSummary.vpcs }}</h3>
              </div>
              <div class="summary-card">
                <div class="metric-label">IP 地址</div>
                <h3>{{ dashboard.resourceSummary.ips }}</h3>
              </div>
              <div class="summary-card">
                <div class="metric-label">云硬盘</div>
                <h3>{{ dashboard.resourceSummary.disks }}</h3>
              </div>
              <div class="summary-card">
                <div class="metric-label">快照</div>
                <h3>{{ dashboard.resourceSummary.snapshots }}</h3>
              </div>
              <div class="summary-card">
                <div class="metric-label">备份</div>
                <h3>{{ dashboard.resourceSummary.backups }}</h3>
              </div>
              <div class="summary-card">
                <div class="metric-label">安全组</div>
                <h3>{{ dashboard.resourceSummary.securityGroups }}</h3>
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <el-card class="page-card">
        <template #header>
          <span>云平台同步日志</span>
        </template>
        <el-table :data="dashboard.providerLogs" stripe>
          <el-table-column label="动作" width="120">
            <template #default="{ row }">
              {{ getLabel(providerActionMap, row.action) }}
            </template>
          </el-table-column>
          <el-table-column prop="resourceId" label="资源标识" min-width="180" />
          <el-table-column prop="message" label="执行说明" min-width="260" />
          <el-table-column label="同步时间" width="180">
            <template #default="{ row }">
              {{ formatDateTime(row.syncedAt) }}
            </template>
          </el-table-column>
          <el-table-column label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="getStatusTagType(row.status)">
                {{ getLabel(providerSyncStatusMap, row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </el-card>
    </el-skeleton>
  </div>
</template>
