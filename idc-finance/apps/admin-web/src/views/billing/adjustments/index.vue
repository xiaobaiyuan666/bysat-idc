<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  fetchAccountTransactions,
  fetchCustomers,
  type AccountTransactionRecord,
  type Customer
} from "@/api/admin";

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const customers = ref<Customer[]>([]);
const transactions = ref<AccountTransactionRecord[]>([]);
const total = ref(0);

const filters = reactive({
  customerId: undefined as number | undefined,
  transactionType: "",
  direction: "",
  keyword: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const currentCustomer = computed(
  () => customers.value.find(item => item.id === filters.customerId) ?? null
);
const pageAmount = computed(() => transactions.value.reduce((sum, item) => sum + item.amount, 0));
const balanceAdjustCount = computed(
  () => transactions.value.filter(item => item.transactionType === "ADJUSTMENT").length
);
const creditAdjustCount = computed(
  () => transactions.value.filter(item => item.transactionType === "CREDIT_LIMIT").length
);
const decreaseCount = computed(() => transactions.value.filter(item => item.direction === "OUT").length);

function hydrateFiltersFromRoute() {
  const customerId = Number(route.query.customerId);
  if (Number.isFinite(customerId) && customerId > 0) {
    filters.customerId = customerId;
  }
  if (route.query.transactionType) {
    filters.transactionType = String(route.query.transactionType);
  }
  if (route.query.direction) {
    filters.direction = String(route.query.direction);
  }
  if (route.query.keyword) {
    filters.keyword = String(route.query.keyword);
  }
}

function buildQuery() {
  return {
    page: pagination.page,
    limit: pagination.limit,
    customerId: filters.customerId,
    transactionType: filters.transactionType || "ADJUSTMENT,CREDIT_LIMIT",
    direction: filters.direction || undefined,
    keyword: filters.keyword || undefined,
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
  filters.transactionType = "";
  filters.direction = "";
  filters.keyword = "";
  pagination.page = 1;
  void loadTransactions();
}

function formatMoney(value: number) {
  return `¥${value.toFixed(2)}`;
}

function transactionTypeLabel(value: string) {
  return (
    {
      ADJUSTMENT: "余额调整",
      CREDIT_LIMIT: "授信调整"
    }[value] ?? value
  );
}

function directionLabel(value: string) {
  return (
    {
      IN: "增加",
      OUT: "扣减",
      FLAT: "平移"
    }[value] ?? value
  );
}

function directionTag(value: string) {
  return (
    {
      IN: "success",
      OUT: "danger",
      FLAT: "info"
    }[value] ?? "info"
  );
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

function openAdjustmentAction(action: "deduct" | "credit") {
  if (!filters.customerId) {
    ElMessage.warning("请先选择客户，再执行资金调整。");
    return;
  }
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId: String(filters.customerId),
      action
    }
  });
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
      eyebrow="财务 / 调整"
      title="资金调整"
      subtitle="集中查看后台人工余额调整、授信额度调整和对应的变更流水，便于复盘财务操作。"
    >
      <template #actions>
        <el-button @click="loadTransactions">刷新列表</el-button>
        <el-button plain @click="openLedger()">打开总台账</el-button>
        <el-button type="warning" plain @click="openAdjustmentAction('deduct')">扣减余额</el-button>
        <el-button type="primary" @click="openAdjustmentAction('credit')">调整授信</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>调整总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatMoney(pageAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>余额调整</span>
            <strong>{{ balanceAdjustCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>授信调整</span>
            <strong>{{ creditAdjustCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>扣减次数</span>
            <strong>{{ decreaseCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前客户</span>
            <strong>{{ currentCustomer?.name ?? "全部客户" }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-select v-model="filters.customerId" placeholder="选择客户" clearable filterable>
            <el-option v-for="item in customers" :key="item.id" :label="`${item.name} / ${item.customerNo}`" :value="item.id" />
          </el-select>
          <el-select v-model="filters.transactionType" placeholder="调整类型" clearable>
            <el-option label="余额调整" value="ADJUSTMENT" />
            <el-option label="授信调整" value="CREDIT_LIMIT" />
          </el-select>
          <el-select v-model="filters.direction" placeholder="变更方向" clearable>
            <el-option label="增加" value="IN" />
            <el-option label="扣减" value="OUT" />
            <el-option label="平移" value="FLAT" />
          </el-select>
          <el-input v-model="filters.keyword" placeholder="搜索流水号、摘要、备注" clearable style="width: 240px" />
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>调整流水</strong>
          <span>当前页 {{ transactions.length }} 条</span>
        </div>
      </div>

      <el-table :data="transactions" border stripe empty-text="暂无资金调整记录">
        <el-table-column prop="transactionNo" label="流水号" min-width="160" />
        <el-table-column label="客户" min-width="200">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomer(row.customerId)">
              {{ row.customerName }} / {{ row.customerNo }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column label="调整类型" min-width="120">
          <template #default="{ row }">{{ transactionTypeLabel(row.transactionType) }}</template>
        </el-table-column>
        <el-table-column label="变更方向" min-width="100">
          <template #default="{ row }">
            <el-tag :type="directionTag(row.direction)" effect="light">{{ directionLabel(row.direction) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="调整金额" min-width="120">
          <template #default="{ row }">{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="变更前后" min-width="240">
          <template #default="{ row }">
            余额 {{ formatMoney(row.balanceBefore) }} → {{ formatMoney(row.balanceAfter) }}
            <br />
            授信 {{ formatMoney(row.creditBefore) }} → {{ formatMoney(row.creditAfter) }}
          </template>
        </el-table-column>
        <el-table-column prop="summary" label="摘要" min-width="220" show-overflow-tooltip />
        <el-table-column prop="remark" label="备注" min-width="220" show-overflow-tooltip />
        <el-table-column label="操作人" min-width="120">
          <template #default="{ row }">{{ row.operatorName || "-" }}</template>
        </el-table-column>
        <el-table-column prop="occurredAt" label="调整时间" min-width="180" />
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
