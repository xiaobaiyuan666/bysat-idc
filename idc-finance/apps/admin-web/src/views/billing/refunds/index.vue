<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchCustomers,
  fetchRefunds,
  type Customer,
  type RefundRecord
} from "@/api/admin";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const customers = ref<Customer[]>([]);
const refunds = ref<RefundRecord[]>([]);
const total = ref(0);

const filters = reactive({
  customerId: undefined as number | undefined,
  invoiceId: undefined as number | undefined,
  keyword: "",
  status: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const pageAmount = computed(() => refunds.value.reduce((sum, item) => sum + item.amount, 0));
const completedCount = computed(() => refunds.value.filter(item => item.status === "COMPLETED").length);
const customerCount = computed(() => new Set(refunds.value.map(item => item.customerId)).size);

function hydrateFiltersFromRoute() {
  const customerId = Number(route.query.customerId);
  if (Number.isFinite(customerId) && customerId > 0) {
    filters.customerId = customerId;
  }

  const invoiceId = Number(route.query.invoiceId);
  if (Number.isFinite(invoiceId) && invoiceId > 0) {
    filters.invoiceId = invoiceId;
  }

  if (route.query.keyword) {
    filters.keyword = String(route.query.keyword);
  }

  if (route.query.status) {
    filters.status = String(route.query.status);
  }
}

function buildQuery() {
  return {
    page: pagination.page,
    limit: pagination.limit,
    customerId: filters.customerId,
    invoiceId: filters.invoiceId,
    keyword: filters.keyword || undefined,
    status: filters.status || undefined,
    sort: "created_at",
    order: "desc"
  };
}

async function loadCustomers() {
  const data = await fetchCustomers();
  customers.value = data.items;
}

async function loadRefunds() {
  loading.value = true;
  try {
    const data = await fetchRefunds(buildQuery());
    refunds.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function applyFilters() {
  pagination.page = 1;
  void loadRefunds();
}

function resetFilters() {
  filters.customerId = undefined;
  filters.invoiceId = undefined;
  filters.keyword = "";
  filters.status = "";
  pagination.page = 1;
  void loadRefunds();
}

function formatMoney(value: number) {
  return `¥${value.toFixed(2)}`;
}

function statusLabel(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "已退款"
  };
  return mapping[status] ?? status;
}

function statusTag(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "warning"
  };
  return mapping[status] ?? "info";
}

function customerName(customerId: number) {
  return customers.value.find(item => item.id === customerId)?.name ?? `客户 #${customerId}`;
}

function openCustomer(customerId: number) {
  void router.push(`/customer/detail/${customerId}`);
}

function openInvoice(invoiceId: number) {
  void router.push(`/billing/invoices/${invoiceId}`);
}

function openAccounts(refundNo: string, customerId: number) {
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId,
      keyword: refundNo
    }
  });
}

onMounted(async () => {
  hydrateFiltersFromRoute();
  await loadCustomers();
  await loadRefunds();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="财务 / 退款"
      title="退款记录"
      subtitle="集中查看退款单、退款原因和关联账单，作为财务售后与人工退款的统一查询入口。"
    >
      <template #actions>
        <el-button @click="loadRefunds">刷新列表</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>退款总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatMoney(pageAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>已完成退款</span>
            <strong>{{ completedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>涉及客户</span>
            <strong>{{ customerCount }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-select v-model="filters.customerId" placeholder="选择客户" clearable filterable>
            <el-option v-for="item in customers" :key="item.id" :label="item.name" :value="item.id" />
          </el-select>
          <el-input-number
            v-model="filters.invoiceId"
            :min="1"
            :controls="false"
            placeholder="账单 ID"
            style="width: 180px"
          />
          <el-input v-model="filters.keyword" placeholder="退款单号 / 原因" clearable style="width: 220px" />
          <el-select v-model="filters.status" placeholder="退款状态" clearable>
            <el-option label="已退款" value="COMPLETED" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>退款流水</strong>
          <span>当前页 {{ refunds.length }} 条</span>
        </div>
      </div>

      <el-table :data="refunds" border stripe empty-text="暂无退款记录">
        <el-table-column prop="refundNo" label="退款单号" min-width="180" />
        <el-table-column label="客户" min-width="180">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomer(row.customerId)">
              {{ customerName(row.customerId) }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="关联账单" min-width="140">
          <template #default="{ row }">
            <el-button link type="primary" @click="openInvoice(row.invoiceId)">#{{ row.invoiceId }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="金额" min-width="120">
          <template #default="{ row }">{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="退款原因" min-width="280" show-overflow-tooltip />
        <el-table-column prop="createdAt" label="退款时间" min-width="180" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button link type="primary" @click="openInvoice(row.invoiceId)">账单详情</el-button>
              <el-button link type="primary" @click="openAccounts(row.refundNo, row.customerId)">台账联查</el-button>
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
          @current-change="(page:number) => { pagination.page = page; void loadRefunds(); }"
          @size-change="(limit:number) => { pagination.limit = limit; pagination.page = 1; void loadRefunds(); }"
        />
      </div>
    </PageWorkbench>
  </div>
</template>
