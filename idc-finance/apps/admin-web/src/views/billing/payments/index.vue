<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import { formatPaymentChannel } from "@/utils/business";
import {
  fetchCustomers,
  fetchPayments,
  type Customer,
  type PaymentRecord
} from "@/api/admin";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const customers = ref<Customer[]>([]);
const payments = ref<PaymentRecord[]>([]);
const total = ref(0);

const filters = reactive({
  customerId: undefined as number | undefined,
  invoiceId: undefined as number | undefined,
  keyword: "",
  channel: "",
  status: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const pageAmount = computed(() => payments.value.reduce((sum, item) => sum + item.amount, 0));
const onlineCount = computed(() => payments.value.filter(item => item.channel === "ONLINE").length);
const offlineCount = computed(() => payments.value.filter(item => item.channel === "OFFLINE").length);
const portalCount = computed(() => payments.value.filter(item => item.source === "PORTAL").length);

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

  if (route.query.channel) {
    filters.channel = String(route.query.channel);
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
    channel: filters.channel || undefined,
    status: filters.status || undefined,
    sort: "paid_at",
    order: "desc"
  };
}

async function loadCustomers() {
  const data = await fetchCustomers();
  customers.value = data.items;
}

async function loadPayments() {
  loading.value = true;
  try {
    const data = await fetchPayments(buildQuery());
    payments.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function applyFilters() {
  pagination.page = 1;
  void loadPayments();
}

function resetFilters() {
  filters.customerId = undefined;
  filters.invoiceId = undefined;
  filters.keyword = "";
  filters.channel = "";
  filters.status = "";
  pagination.page = 1;
  void loadPayments();
}

function formatMoney(value: number) {
  return `¥${value.toFixed(2)}`;
}

function statusLabel(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "已收款"
  };
  return mapping[status] ?? status;
}

function statusTag(status: string) {
  const mapping: Record<string, string> = {
    COMPLETED: "success"
  };
  return mapping[status] ?? "info";
}

function sourceLabel(source: string) {
  const mapping: Record<string, string> = {
    PORTAL: "客户端",
    ADMIN: "后台",
    SYSTEM: "系统"
  };
  return mapping[source] ?? source;
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

function openAccounts(paymentNo: string, customerId: number) {
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId,
      keyword: paymentNo
    }
  });
}

onMounted(async () => {
  hydrateFiltersFromRoute();
  await loadCustomers();
  await loadPayments();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="财务 / 支付"
      title="支付记录"
      subtitle="集中查看后台收款、客户在线支付和余额支付流水，支持按客户、账单和渠道联查。"
    >
      <template #actions>
        <el-button @click="loadPayments">刷新列表</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>支付总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatMoney(pageAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>在线支付</span>
            <strong>{{ onlineCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>线下支付</span>
            <strong>{{ offlineCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>客户端来源</span>
            <strong>{{ portalCount }}</strong>
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
          <el-input v-model="filters.keyword" placeholder="支付单号 / 交易号" clearable style="width: 220px" />
          <el-select v-model="filters.channel" placeholder="支付渠道" clearable>
            <el-option label="在线支付" value="ONLINE" />
            <el-option label="线下支付" value="OFFLINE" />
            <el-option label="余额支付" value="BALANCE" />
            <el-option label="系统入账" value="SYSTEM" />
          </el-select>
          <el-select v-model="filters.status" placeholder="支付状态" clearable>
            <el-option label="已收款" value="COMPLETED" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>支付流水</strong>
          <span>当前页 {{ payments.length }} 条</span>
        </div>
      </div>

      <el-table :data="payments" border stripe empty-text="暂无支付记录">
        <el-table-column prop="paymentNo" label="支付单号" min-width="180" />
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
        <el-table-column label="渠道" min-width="120">
          <template #default="{ row }">{{ formatPaymentChannel(localeStore.locale, row.channel || "SYSTEM") }}</template>
        </el-table-column>
        <el-table-column label="金额" min-width="120">
          <template #default="{ row }">{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="来源" min-width="120">
          <template #default="{ row }">{{ sourceLabel(row.source) }}</template>
        </el-table-column>
        <el-table-column label="状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="statusTag(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="tradeNo" label="交易流水" min-width="220" show-overflow-tooltip />
        <el-table-column prop="operator" label="登记人" min-width="120" />
        <el-table-column prop="paidAt" label="支付时间" min-width="180" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button link type="primary" @click="openInvoice(row.invoiceId)">账单详情</el-button>
              <el-button link type="primary" @click="openAccounts(row.paymentNo, row.customerId)">台账联查</el-button>
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
          @current-change="(page:number) => { pagination.page = page; void loadPayments(); }"
          @size-change="(limit:number) => { pagination.limit = limit; pagination.page = 1; void loadPayments(); }"
        />
      </div>
    </PageWorkbench>
  </div>
</template>
