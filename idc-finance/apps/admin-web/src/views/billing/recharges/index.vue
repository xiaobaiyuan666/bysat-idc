<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import { formatPaymentChannel } from "@/utils/business";
import {
  fetchAccountTransactions,
  fetchCustomers,
  type AccountTransactionRecord,
  type Customer
} from "@/api/admin";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const customers = ref<Customer[]>([]);
const transactions = ref<AccountTransactionRecord[]>([]);
const total = ref(0);

const filters = reactive({
  customerId: undefined as number | undefined,
  keyword: "",
  channel: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const currentCustomer = computed(
  () => customers.value.find(item => item.id === filters.customerId) ?? null
);
const pageAmount = computed(() => transactions.value.reduce((sum, item) => sum + item.amount, 0));
const offlineCount = computed(() => transactions.value.filter(item => item.channel === "OFFLINE").length);
const onlineCount = computed(() => transactions.value.filter(item => item.channel === "ONLINE").length);
const systemCount = computed(() => transactions.value.filter(item => item.channel === "SYSTEM").length);

function hydrateFiltersFromRoute() {
  const customerId = Number(route.query.customerId);
  if (Number.isFinite(customerId) && customerId > 0) {
    filters.customerId = customerId;
  }
  if (route.query.keyword) {
    filters.keyword = String(route.query.keyword);
  }
  if (route.query.channel) {
    filters.channel = String(route.query.channel);
  }
}

function buildQuery() {
  return {
    page: pagination.page,
    limit: pagination.limit,
    customerId: filters.customerId,
    keyword: filters.keyword || undefined,
    channel: filters.channel || undefined,
    transactionType: "RECHARGE",
    direction: "IN",
    sort: "occurred_at",
    order: "desc"
  };
}

async function loadCustomers() {
  const data = await fetchCustomers();
  customers.value = data.items;
}

async function loadTransactions() {
  loading.value = true;
  try {
    const data = await fetchAccountTransactions(buildQuery());
    transactions.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function applyFilters() {
  pagination.page = 1;
  void loadTransactions();
}

function resetFilters() {
  filters.customerId = undefined;
  filters.keyword = "";
  filters.channel = "";
  pagination.page = 1;
  void loadTransactions();
}

function formatMoney(value: number) {
  return `¥${value.toFixed(2)}`;
}

function openCustomer(customerId: number) {
  void router.push(`/customer/detail/${customerId}`);
}

function openLedger(keyword?: string) {
  const query: Record<string, string> = {};
  if (filters.customerId) {
    query.customerId = String(filters.customerId);
  }
  if (keyword) {
    query.keyword = keyword;
  }
  void router.push({ path: "/billing/accounts", query });
}

function openRechargeAction() {
  if (!filters.customerId) {
    ElMessage.warning("请先选择客户，再执行线下充值。");
    return;
  }
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId: String(filters.customerId),
      action: "recharge"
    }
  });
}

function relatedReference(row: AccountTransactionRecord) {
  return row.paymentNo || row.invoiceNo || row.orderNo || row.transactionNo;
}

onMounted(async () => {
  hydrateFiltersFromRoute();
  await loadCustomers();
  await loadTransactions();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="财务 / 充值"
      title="充值记录"
      subtitle="单独查看客户充值入账流水，和资金总台账分开，方便财务核对充值、补款和到账情况。"
    >
      <template #actions>
        <el-button @click="loadTransactions">刷新列表</el-button>
        <el-button plain @click="openLedger()">打开总台账</el-button>
        <el-button type="primary" @click="openRechargeAction">线下充值</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>充值总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatMoney(pageAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>线下入账</span>
            <strong>{{ offlineCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>在线入账</span>
            <strong>{{ onlineCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前客户</span>
            <strong>{{ currentCustomer?.name ?? "全部客户" }}</strong>
          </div>
          <div class="summary-pill">
            <span>系统入账</span>
            <strong>{{ systemCount }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-select v-model="filters.customerId" placeholder="选择客户" clearable filterable>
            <el-option v-for="item in customers" :key="item.id" :label="`${item.name} / ${item.customerNo}`" :value="item.id" />
          </el-select>
          <el-input v-model="filters.keyword" placeholder="搜索流水号、支付单号、摘要" clearable style="width: 240px" />
          <el-select v-model="filters.channel" placeholder="到账渠道" clearable>
            <el-option label="线下入账" value="OFFLINE" />
            <el-option label="在线入账" value="ONLINE" />
            <el-option label="系统入账" value="SYSTEM" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>充值流水</strong>
          <span>当前页 {{ transactions.length }} 条</span>
        </div>
      </div>

      <el-table :data="transactions" border stripe empty-text="暂无充值记录">
        <el-table-column prop="transactionNo" label="流水号" min-width="160" />
        <el-table-column label="客户" min-width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomer(row.customerId)">
              {{ row.customerName }} / {{ row.customerNo }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="到账渠道" min-width="120">
          <template #default="{ row }">{{ formatPaymentChannel(localeStore.locale, row.channel || "SYSTEM") }}</template>
        </el-table-column>
        <el-table-column label="充值金额" min-width="120">
          <template #default="{ row }">{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="变更前后" min-width="210">
          <template #default="{ row }">
            余额 {{ formatMoney(row.balanceBefore) }} → {{ formatMoney(row.balanceAfter) }}
          </template>
        </el-table-column>
        <el-table-column label="关联单据" min-width="180">
          <template #default="{ row }">
            <div>{{ row.invoiceNo || "-" }}</div>
            <div style="color: #94a3b8">{{ relatedReference(row) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="summary" label="摘要" min-width="240" show-overflow-tooltip />
        <el-table-column label="操作人" min-width="120">
          <template #default="{ row }">{{ row.operatorName || "-" }}</template>
        </el-table-column>
        <el-table-column prop="occurredAt" label="到账时间" min-width="180" />
        <el-table-column label="操作" min-width="180" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button link type="primary" @click="openCustomer(row.customerId)">客户详情</el-button>
              <el-button link type="primary" @click="openLedger(row.transactionNo)">台账联查</el-button>
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
          @current-change="(page:number) => { pagination.page = page; void loadTransactions(); }"
          @size-change="(limit:number) => { pagination.limit = limit; pagination.page = 1; void loadTransactions(); }"
        />
      </div>
    </PageWorkbench>
  </div>
</template>
