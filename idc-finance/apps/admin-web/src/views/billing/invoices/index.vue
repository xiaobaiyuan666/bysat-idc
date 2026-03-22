<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import { fetchInvoices, type InvoiceQuery, type InvoiceRecord } from "@/api/admin";

type TabKey = "ALL" | "UNPAID" | "PAID" | "REFUNDED";

const router = useRouter();

const loading = ref(false);
const invoices = ref<InvoiceRecord[]>([]);
const total = ref(0);
const selectedRows = ref<InvoiceRecord[]>([]);
const activeTab = ref<TabKey>("ALL");
const advancedVisible = ref(false);

const filters = reactive({
  keyword: "",
  orderNo: "",
  productName: "",
  cycle: ""
});

const pagination = reactive({
  page: 1,
  limit: 20
});

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: total.value },
  { key: "UNPAID", label: "未支付", count: invoices.value.filter(item => item.status === "UNPAID").length },
  { key: "PAID", label: "已支付", count: invoices.value.filter(item => item.status === "PAID").length },
  { key: "REFUNDED", label: "已退款", count: invoices.value.filter(item => item.status === "REFUNDED").length }
]);

const unpaidCount = computed(() => invoices.value.filter(item => item.status === "UNPAID").length);
const paidCount = computed(() => invoices.value.filter(item => item.status === "PAID").length);
const refundedCount = computed(() => invoices.value.filter(item => item.status === "REFUNDED").length);
const currentReceivable = computed(() =>
  invoices.value
    .filter(item => item.status === "UNPAID")
    .reduce((sum, item) => sum + item.totalAmount, 0)
);

function buildQuery(): InvoiceQuery {
  const query: InvoiceQuery = {
    page: pagination.page,
    limit: pagination.limit,
    sort: "created_at",
    order: "desc"
  };
  if (activeTab.value !== "ALL") query.status = activeTab.value;
  if (filters.keyword.trim()) query.invoice_no = filters.keyword.trim();
  if (filters.orderNo.trim()) query.order_no = filters.orderNo.trim();
  if (filters.productName.trim()) query.product_name = filters.productName.trim();
  if (filters.cycle) query.billing_cycle = filters.cycle;
  return query;
}

async function loadInvoices() {
  loading.value = true;
  try {
    const data = await fetchInvoices(buildQuery());
    invoices.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.orderNo = "";
  filters.productName = "";
  filters.cycle = "";
  pagination.page = 1;
  void loadInvoices();
}

function applyFilters() {
  pagination.page = 1;
  void loadInvoices();
}

function handleSelectionChange(rows: InvoiceRecord[]) {
  selectedRows.value = rows;
}

function formatCurrency(value: number) {
  return `¥${value.toFixed(2)}`;
}

function formatCycle(value: string) {
  const mapping: Record<string, string> = {
    monthly: "月付",
    quarterly: "季付",
    semiannual: "半年付",
    semiannually: "半年付",
    annual: "年付",
    annually: "年付",
    onetime: "一次性"
  };
  return mapping[value] ?? value;
}

function statusLabel(value: string) {
  const mapping: Record<string, string> = {
    UNPAID: "未支付",
    PAID: "已支付",
    REFUNDED: "已退款"
  };
  return mapping[value] ?? value;
}

function statusTagType(value: string) {
  const mapping: Record<string, string> = {
    UNPAID: "warning",
    PAID: "success",
    REFUNDED: "info"
  };
  return mapping[value] ?? "info";
}

async function copySelectedInvoiceNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择账单");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.invoiceNo).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个账单编号`);
}

function exportRows(rows: InvoiceRecord[], filename: string) {
  downloadCsv(
    filename,
    ["账单编号", "订单编号", "商品名称", "计费周期", "账单金额", "账单状态", "到期时间", "支付时间"],
    rows.map(item => [
      item.invoiceNo,
      item.orderNo,
      item.productName,
      formatCycle(item.billingCycle),
      formatCurrency(item.totalAmount),
      statusLabel(item.status),
      item.dueAt || "-",
      item.paidAt || "-"
    ])
  );
}

function exportCurrent() {
  exportRows(invoices.value, `invoices-page-${pagination.page}.csv`);
  ElMessage.success("当前页账单已导出");
}

function exportSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择账单");
    return;
  }
  exportRows(selectedRows.value, "invoices-selected.csv");
  ElMessage.success("已导出选中账单");
}

function openDetail(row: InvoiceRecord) {
  void router.push(`/billing/invoices/${row.id}`);
}

watch(activeTab, () => {
  pagination.page = 1;
  void loadInvoices();
});

onMounted(() => {
  void loadInvoices();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="财务 / 账单"
      title="账单列表"
      subtitle="集中查看应收账单、支付状态、到期时间和退款状态，作为收款与退款的统一入口。"
    >
      <template #actions>
        <el-button @click="loadInvoices">刷新列表</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>账单总数</span><strong>{{ total }}</strong></div>
          <div class="summary-pill">
            <span>未支付</span>
            <strong>{{ unpaidCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>已支付</span>
            <strong>{{ paidCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>已退款</span>
            <strong>{{ refundedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前应收</span>
            <strong>{{ formatCurrency(currentReceivable) }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="filters.keyword" placeholder="账单编号" clearable />
          <el-select v-model="filters.cycle" placeholder="计费周期" clearable>
            <el-option label="月付" value="monthly" />
            <el-option label="季付" value="quarterly" />
            <el-option label="半年付" value="semiannual" />
            <el-option label="年付" value="annual" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="advancedVisible = !advancedVisible">
              {{ advancedVisible ? "收起高级筛选" : "展开高级筛选" }}
            </el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div v-if="advancedVisible" class="filter-bar filter-bar--compact">
          <el-input v-model="filters.orderNo" placeholder="订单编号" clearable />
          <el-input v-model="filters.productName" placeholder="商品名称" clearable />
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>财务列表</strong>
          <span>当前页 {{ invoices.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain @click="copySelectedInvoiceNos">复制账单号</el-button>
          <el-button plain @click="exportSelected">导出选中</el-button>
          <el-button plain @click="exportCurrent">导出当前页</el-button>
        </div>
      </div>

      <el-table
        :data="invoices"
        border
        stripe
        empty-text="暂无符合条件的账单"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column prop="invoiceNo" label="账单编号" min-width="170" />
        <el-table-column prop="orderNo" label="订单编号" min-width="170" />
        <el-table-column prop="productName" label="商品名称" min-width="240" show-overflow-tooltip />
        <el-table-column label="计费周期" min-width="100">
          <template #default="{ row }">{{ formatCycle(row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column label="账单金额" min-width="120">
          <template #default="{ row }">{{ formatCurrency(row.totalAmount) }}</template>
        </el-table-column>
        <el-table-column label="账单状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="dueAt" label="到期时间" min-width="180" />
        <el-table-column label="支付时间" min-width="180">
          <template #default="{ row }">{{ row.paidAt || "-" }}</template>
        </el-table-column>
        <el-table-column label="操作" min-width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row)">进入工作台</el-button>
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
          @current-change="(page:number) => { pagination.page = page; void loadInvoices(); }"
          @size-change="(limit:number) => { pagination.limit = limit; pagination.page = 1; void loadInvoices(); }"
        />
      </div>
    </PageWorkbench>
  </div>
</template>
