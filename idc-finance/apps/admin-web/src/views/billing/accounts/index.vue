<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import { formatPaymentChannel } from "@/utils/business";
import {
  adjustCustomerWallet,
  fetchAccountTransactions,
  fetchCustomerWallet,
  fetchCustomers,
  type AccountTransactionRecord,
  type AdjustWalletRequest,
  type Customer,
  type CustomerWalletResponse
} from "@/api/admin";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(false);
const customers = ref<Customer[]>([]);
const transactions = ref<AccountTransactionRecord[]>([]);
const total = ref(0);
const walletLoading = ref(false);
const walletDetail = ref<CustomerWalletResponse | null>(null);
const dialogVisible = ref(false);
const submitting = ref(false);

const filters = reactive({
  customerId: undefined as number | undefined,
  keyword: "",
  transactionType: "",
  direction: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const form = reactive<AdjustWalletRequest>({
  target: "BALANCE",
  operation: "INCREASE",
  amount: 0,
  summary: "",
  remark: ""
});

const currentCustomer = computed(
  () => customers.value.find(item => item.id === filters.customerId) ?? null
);
const pageAmount = computed(() => transactions.value.reduce((sum, item) => sum + item.amount, 0));
const dialogTitle = computed(() => (form.target === "BALANCE" ? "调整账户余额" : "调整授信额度"));
const adjustPreview = computed(() => {
  const target = form.target === "BALANCE" ? "账户余额" : "授信额度";
  const action =
    form.operation === "INCREASE" ? "增加" : form.operation === "DECREASE" ? "扣减" : "直接设定";
  const description =
    form.target === "BALANCE"
      ? "适用于线下充值、手工扣款、余额纠偏等场景，提交后会立即生成一笔资金台账流水。"
      : "适用于开通授信、缩减授信和额度校正等场景，提交后会立即生成一笔授信调整流水。";
  return {
    title: `${action}${target}`,
    description
  };
});
const routeActionConsumed = ref(false);

function formatMoney(value: number) {
  return `¥${value.toFixed(2)}`;
}

function directionLabel(value: string) {
  return (
    {
      IN: "收入",
      OUT: "支出",
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

function transactionTypeLabel(value: string) {
  return (
    {
      RECHARGE: "充值",
      CONSUME: "余额支付",
      REFUND: "退款回退",
      ADJUSTMENT: "手工调整",
      CREDIT_LIMIT: "授信调整"
    }[value] ?? value
  );
}

function defaultSummary() {
  if (form.target === "BALANCE" && form.operation === "INCREASE") return "线下充值补充余额";
  if (form.target === "BALANCE" && form.operation === "DECREASE") return "人工扣减账户余额";
  if (form.target === "BALANCE") return "人工设定账户余额";
  if (form.operation === "INCREASE") return "增加客户授信额度";
  if (form.operation === "DECREASE") return "扣减客户授信额度";
  return "人工设定授信额度";
}

function hydrateFiltersFromRoute() {
  if (route.query.customerId) {
    const customerId = Number(route.query.customerId);
    if (Number.isFinite(customerId) && customerId > 0) {
      filters.customerId = customerId;
    }
  }
  if (route.query.keyword) {
    filters.keyword = String(route.query.keyword);
  }
  if (route.query.transactionType) {
    filters.transactionType = String(route.query.transactionType);
  }
  if (route.query.direction) {
    filters.direction = String(route.query.direction);
  }
}

function maybeOpenActionFromRoute() {
  if (routeActionConsumed.value || !filters.customerId) return;
  const action = String(route.query.action ?? "").toLowerCase();
  if (!action) return;
  if (action === "recharge") {
    openAdjustDialog("BALANCE", "INCREASE");
    routeActionConsumed.value = true;
    return;
  }
  if (action === "deduct") {
    openAdjustDialog("BALANCE", "DECREASE");
    routeActionConsumed.value = true;
    return;
  }
  if (action === "credit") {
    openAdjustDialog("CREDIT_LIMIT", "SET");
    routeActionConsumed.value = true;
  }
}

function buildQuery() {
  return {
    page: pagination.page,
    limit: pagination.limit,
    customerId: filters.customerId,
    keyword: filters.keyword || undefined,
    transactionType: filters.transactionType || undefined,
    direction: filters.direction || undefined,
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

async function loadWallet() {
  if (!filters.customerId) {
    walletDetail.value = null;
    return;
  }
  walletLoading.value = true;
  try {
    walletDetail.value = await fetchCustomerWallet(filters.customerId);
  } finally {
    walletLoading.value = false;
  }
}

function applyFilters() {
  pagination.page = 1;
  void Promise.all([loadTransactions(), loadWallet()]);
}

function resetFilters() {
  filters.customerId = undefined;
  filters.keyword = "";
  filters.transactionType = "";
  filters.direction = "";
  pagination.page = 1;
  void Promise.all([loadTransactions(), loadWallet()]);
}

function openAdjustDialog(target = "BALANCE", operation = "INCREASE") {
  if (!filters.customerId) {
    ElMessage.warning("请先选择一个客户");
    return;
  }
  form.target = target;
  form.operation = operation;
  form.amount = 0;
  form.summary = "";
  form.remark = "";
  dialogVisible.value = true;
}

function openRechargeWorkbench() {
  const query: Record<string, string> = {};
  if (filters.customerId) {
    query.customerId = String(filters.customerId);
  }
  void router.push({ path: "/billing/recharges", query });
}

function openAdjustmentWorkbench() {
  const query: Record<string, string> = {};
  if (filters.customerId) {
    query.customerId = String(filters.customerId);
  }
  void router.push({ path: "/billing/adjustments", query });
}

async function submitAdjust() {
  if (!filters.customerId) return;
  if (form.amount < 0 || (form.operation !== "SET" && form.amount <= 0)) {
    ElMessage.warning("请输入有效金额");
    return;
  }
  if (!(form.summary ?? "").trim()) {
    form.summary = defaultSummary();
  }
  submitting.value = true;
  try {
    const result = await adjustCustomerWallet(filters.customerId, form);
    dialogVisible.value = false;
    ElMessage.success(`已生成台账流水 ${result.transaction.transactionNo}`);
    await Promise.all([loadTransactions(), loadWallet()]);
  } finally {
    submitting.value = false;
  }
}

watch(
  () => filters.customerId,
  () => {
    pagination.page = 1;
    void Promise.all([loadTransactions(), loadWallet()]);
  }
);

onMounted(async () => {
  hydrateFiltersFromRoute();
  await loadCustomers();
  await loadTransactions();
  await loadWallet();
  maybeOpenActionFromRoute();
});
</script>

<template>
  <div v-loading="loading || walletLoading">
    <PageWorkbench
      eyebrow="财务 / 台账"
      title="资金台账"
      subtitle="统一查看客户余额、授信额度、余额支付和手工调整流水。"
    >
      <template #actions>
        <el-button @click="loadTransactions">刷新台账</el-button>
        <el-button plain @click="openRechargeWorkbench">充值记录</el-button>
        <el-button plain @click="openAdjustmentWorkbench">资金调整</el-button>
        <el-button type="primary" plain @click="openAdjustDialog('BALANCE', 'INCREASE')">线下充值</el-button>
        <el-button type="warning" plain @click="openAdjustDialog('BALANCE', 'DECREASE')">扣减余额</el-button>
        <el-button type="primary" @click="openAdjustDialog('CREDIT_LIMIT', 'SET')">调整授信</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>流水总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatMoney(pageAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>选中客户</span>
            <strong>{{ currentCustomer?.name ?? "全部客户" }}</strong>
          </div>
          <div class="summary-pill">
            <span>余额</span>
            <strong>{{ walletDetail ? formatMoney(walletDetail.wallet.balance) : "-" }}</strong>
          </div>
          <div class="summary-pill">
            <span>可用授信</span>
            <strong>{{ walletDetail ? formatMoney(walletDetail.wallet.availableCredit) : "-" }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-select v-model="filters.customerId" placeholder="选择客户" clearable filterable>
            <el-option v-for="item in customers" :key="item.id" :label="`${item.name} / ${item.customerNo}`" :value="item.id" />
          </el-select>
          <el-input v-model="filters.keyword" placeholder="搜索流水号、摘要、账单号" clearable />
          <el-select v-model="filters.transactionType" placeholder="流水类型" clearable>
            <el-option label="充值" value="RECHARGE" />
            <el-option label="余额支付" value="CONSUME" />
            <el-option label="退款回退" value="REFUND" />
            <el-option label="手工调整" value="ADJUSTMENT" />
            <el-option label="授信调整" value="CREDIT_LIMIT" />
          </el-select>
          <el-select v-model="filters.direction" placeholder="方向" clearable>
            <el-option label="收入" value="IN" />
            <el-option label="支出" value="OUT" />
            <el-option label="平移" value="FLAT" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>
      </template>

      <div class="portal-grid portal-grid--two" style="margin-bottom: 16px">
        <div class="panel-card">
          <div class="section-card__head">
            <strong>钱包概览</strong>
            <span class="section-card__meta">{{ currentCustomer?.customerNo ?? "请选择客户" }}</span>
          </div>
          <el-empty v-if="!walletDetail" description="选择客户后查看余额与授信" :image-size="72" />
          <el-descriptions v-else :column="2" border>
            <el-descriptions-item label="客户名称">{{ walletDetail.wallet.customerName }}</el-descriptions-item>
            <el-descriptions-item label="客户编号">{{ walletDetail.wallet.customerNo }}</el-descriptions-item>
            <el-descriptions-item label="账户余额">{{ formatMoney(walletDetail.wallet.balance) }}</el-descriptions-item>
            <el-descriptions-item label="授信额度">{{ formatMoney(walletDetail.wallet.creditLimit) }}</el-descriptions-item>
            <el-descriptions-item label="已用授信">{{ formatMoney(walletDetail.wallet.creditUsed) }}</el-descriptions-item>
            <el-descriptions-item label="可用授信">{{ formatMoney(walletDetail.wallet.availableCredit) }}</el-descriptions-item>
            <el-descriptions-item label="更新时间" :span="2">{{ walletDetail.wallet.updatedAt || "-" }}</el-descriptions-item>
          </el-descriptions>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>最近流水</strong>
            <span class="section-card__meta">最近 20 条</span>
          </div>
          <el-table v-if="walletDetail" :data="walletDetail.transactions.slice(0, 6)" border stripe size="small" empty-text="暂无流水">
            <el-table-column prop="transactionNo" label="流水号" min-width="140" />
            <el-table-column label="类型" min-width="110">
              <template #default="{ row }">{{ transactionTypeLabel(row.transactionType) }}</template>
            </el-table-column>
            <el-table-column label="渠道" min-width="110">
              <template #default="{ row }">{{ formatPaymentChannel(localeStore.locale, row.channel || "SYSTEM") }}</template>
            </el-table-column>
            <el-table-column label="金额" min-width="110">
              <template #default="{ row }">{{ formatMoney(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="occurredAt" label="时间" min-width="160" />
          </el-table>
          <el-empty v-else description="暂无客户上下文" :image-size="72" />
        </div>
      </div>

      <el-table :data="transactions" border stripe empty-text="暂无匹配的台账流水">
        <el-table-column prop="transactionNo" label="流水号" min-width="150" />
        <el-table-column label="客户" min-width="180">
          <template #default="{ row }">{{ row.customerName }} / {{ row.customerNo }}</template>
        </el-table-column>
        <el-table-column label="类型" min-width="120">
          <template #default="{ row }">{{ transactionTypeLabel(row.transactionType) }}</template>
        </el-table-column>
        <el-table-column label="方向" min-width="90">
          <template #default="{ row }">
            <el-tag :type="directionTag(row.direction)" effect="light">{{ directionLabel(row.direction) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="金额" min-width="110">
          <template #default="{ row }">{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="变更前后" min-width="220">
          <template #default="{ row }">
            余额 {{ formatMoney(row.balanceBefore) }} → {{ formatMoney(row.balanceAfter) }}
            <br />
            授信 {{ formatMoney(row.creditBefore) }} → {{ formatMoney(row.creditAfter) }}
          </template>
        </el-table-column>
        <el-table-column label="关联单据" min-width="180">
          <template #default="{ row }">
            <div>{{ row.invoiceNo || "-" }}</div>
            <div style="color: #94a3b8">{{ row.orderNo || row.paymentNo || row.refundNo || "-" }}</div>
          </template>
        </el-table-column>
        <el-table-column label="渠道" min-width="110">
          <template #default="{ row }">{{ formatPaymentChannel(localeStore.locale, row.channel || "SYSTEM") }}</template>
        </el-table-column>
        <el-table-column prop="summary" label="摘要" min-width="220" show-overflow-tooltip />
        <el-table-column label="操作人" min-width="120">
          <template #default="{ row }">{{ row.operatorName || "-" }}</template>
        </el-table-column>
        <el-table-column prop="occurredAt" label="发生时间" min-width="170" />
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

    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="560px">
      <el-form label-position="top">
        <el-alert
          :title="adjustPreview.title"
          :description="adjustPreview.description"
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="调整目标">
          <el-radio-group v-model="form.target">
            <el-radio-button label="BALANCE">账户余额</el-radio-button>
            <el-radio-button label="CREDIT_LIMIT">授信额度</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="操作方式">
          <el-radio-group v-model="form.operation">
            <el-radio-button label="INCREASE">增加</el-radio-button>
            <el-radio-button label="DECREASE">扣减</el-radio-button>
            <el-radio-button label="SET">直接设为</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="金额">
          <el-input-number v-model="form.amount" :min="0" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item label="摘要">
          <el-input v-model="form.summary" placeholder="例如：线下收款补充余额 / 调整授信额度" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="记录本次调整原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitAdjust">
          {{ form.target === "BALANCE" ? "确认记账" : "确认调整" }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>
