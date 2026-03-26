<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import {
  fetchOrders,
  fetchOrderRequests,
  fetchProviderAccounts,
  updatePendingOrder,
  type OrderQuery,
  type OrderRecord,
  type ProviderAccount
} from "@/api/admin";

type TabKey = "ALL" | "PENDING" | "ACTIVE" | "COMPLETED" | "CANCELLED";
type LifecycleStatus = Exclude<TabKey, "ALL">;

const route = useRoute();
const router = useRouter();

const loading = ref(false);
const advancedVisible = ref(false);
const selectedRows = ref<OrderRecord[]>([]);
const orders = ref<OrderRecord[]>([]);
const total = ref(0);
const providerAccounts = ref<ProviderAccount[]>([]);
const pendingRequestCount = ref(0);
const activeTab = ref<TabKey>("ALL");
const quickActionVisible = ref(false);
const quickActionLoading = ref(false);
const quickActionRows = ref<OrderRecord[]>([]);

const pagination = reactive({
  page: 1,
  limit: 20,
  sort: "create_time",
  order: "desc"
});

const filters = reactive({
  orderNo: "",
  productName: "",
  customerId: "",
  amount: "",
  payment: "",
  payStatus: "",
  dateRange: [] as string[]
});

const quickActionForm = reactive({
  status: "ACTIVE" as LifecycleStatus,
  reason: ""
});

const paymentOptions = [
  { label: "全部支付方式", value: "" },
  { label: "未支付", value: "UNPAID" },
  { label: "线上支付", value: "ONLINE" },
  { label: "线下收款", value: "OFFLINE" },
  { label: "余额支付", value: "BALANCE" }
];

const payStatusOptions = [
  { label: "全部支付状态", value: "" },
  { label: "未支付", value: "UNPAID" },
  { label: "已支付", value: "PAID" },
  { label: "已退款", value: "REFUNDED" }
];

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: total.value },
  { key: "PENDING", label: "待支付", count: orders.value.filter(item => item.status === "PENDING").length },
  { key: "ACTIVE", label: "待开通", count: orders.value.filter(item => item.status === "ACTIVE").length },
  { key: "COMPLETED", label: "已完成", count: orders.value.filter(item => item.status === "COMPLETED").length },
  { key: "CANCELLED", label: "已取消", count: orders.value.filter(item => item.status === "CANCELLED").length }
]);

const totalAmount = computed(() => orders.value.reduce((sum, item) => sum + item.amount, 0));
const paidCount = computed(() => orders.value.filter(item => item.payStatus === "PAID").length);
const unpaidCount = computed(() => orders.value.filter(item => item.payStatus === "UNPAID").length);
const refundedCount = computed(() => orders.value.filter(item => item.payStatus === "REFUNDED").length);
const automatedCount = computed(() => orders.value.filter(item => item.automationType !== "LOCAL").length);
const quickActionTargetCount = computed(() => quickActionRows.value.length);
const quickActionStatusLabel = computed(() => statusLabel(quickActionForm.status));

const lifecycleActionOptions = [
  { label: "转待支付", value: "PENDING" as LifecycleStatus },
  { label: "转待开通", value: "ACTIVE" as LifecycleStatus },
  { label: "标记完成", value: "COMPLETED" as LifecycleStatus },
  { label: "标记取消", value: "CANCELLED" as LifecycleStatus }
];

function buildQuery(): OrderQuery {
  const query: OrderQuery = {
    page: pagination.page,
    limit: pagination.limit,
    sort: pagination.sort,
    order: pagination.order
  };

  if (activeTab.value !== "ALL") query.status = activeTab.value;
  if (filters.orderNo.trim()) query.ordernum = filters.orderNo.trim();
  if (filters.productName.trim()) query.product_name = filters.productName.trim();
  if (filters.customerId.trim()) query.uid = filters.customerId.trim();
  if (filters.amount.trim()) query.amount = filters.amount.trim();
  if (filters.payStatus) query.pay_status = filters.payStatus;
  if (filters.payment && filters.payment !== "UNPAID") query.payment = filters.payment;
  if (filters.payment === "UNPAID") query.pay_status = "UNPAID";
  if (filters.dateRange.length === 2) {
    query.start_time = filters.dateRange[0];
    query.end_time = filters.dateRange[1];
  }

  return query;
}

async function loadOrders() {
  loading.value = true;
  try {
    const [orderData, accounts, pendingRequests] = await Promise.all([
      fetchOrders(buildQuery()),
      fetchProviderAccounts(),
      fetchOrderRequests({ status: "PENDING", limit: 1, page: 1 })
    ]);
    orders.value = orderData.items;
    total.value = orderData.total;
    providerAccounts.value = accounts;
    pendingRequestCount.value = pendingRequests.total;
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.orderNo = "";
  filters.productName = "";
  filters.customerId = "";
  filters.amount = "";
  filters.payment = "";
  filters.payStatus = "";
  filters.dateRange = [];
  pagination.page = 1;
  void router.push({ path: "/orders/list", query: {} });
}

function applyFilters() {
  pagination.page = 1;
  void loadOrders();
}

function handlePageChange(page: number) {
  pagination.page = page;
  void loadOrders();
}

function handleLimitChange(limit: number) {
  pagination.limit = limit;
  pagination.page = 1;
  void loadOrders();
}

function handleSelectionChange(rows: OrderRecord[]) {
  selectedRows.value = rows;
}

function providerAccountLabel(id?: number) {
  if (!id) return "未绑定接口账户";
  const matched = providerAccounts.value.find(item => item.id === id);
  if (!matched) return `接口账户 #${id}`;
  return matched.name || matched.sourceName || matched.baseUrl;
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
    PENDING: "待支付",
    ACTIVE: "待开通",
    COMPLETED: "已完成",
    CANCELLED: "已取消"
  };
  return mapping[value] ?? value;
}

function payStatusLabel(value: string) {
  const mapping: Record<string, string> = {
    PAID: "已支付",
    UNPAID: "未支付",
    REFUNDED: "已退款"
  };
  return mapping[value] ?? value;
}

function paymentLabel(value: string) {
  const mapping: Record<string, string> = {
    ONLINE: "线上支付",
    OFFLINE: "线下收款",
    BALANCE: "余额支付",
    UNPAID: "未支付"
  };
  return mapping[value] ?? (value || "未支付");
}

function automationLabel(value: string) {
  const mapping: Record<string, string> = {
    LOCAL: "本地模块",
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "财务上下游",
    RESOURCE: "资源池",
    MANUAL: "手动资源"
  };
  return mapping[value] ?? value;
}

function statusTagType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    ACTIVE: "primary",
    COMPLETED: "success",
    CANCELLED: "info"
  };
  return mapping[status] ?? "info";
}

function payStatusTagType(status: string) {
  if (status === "PAID") return "success";
  if (status === "REFUNDED") return "info";
  return "warning";
}

function defaultReasonForStatus(status: LifecycleStatus) {
  const mapping: Record<LifecycleStatus, string> = {
    PENDING: "后台回退为待支付，等待重新收款或人工复核",
    ACTIVE: "后台确认订单进入待开通流程",
    COMPLETED: "后台确认订单已完成交付",
    CANCELLED: "后台关闭订单并终止后续交付"
  };
  return mapping[status];
}

function openCustomerWorkbench(row: OrderRecord) {
  void router.push(`/customer/detail/${row.customerId}`);
}

function openInvoiceWorkbench(row: OrderRecord) {
  void router.push({
    path: "/billing/invoices",
    query: {
      orderNo: row.orderNo,
      customerId: String(row.customerId)
    }
  });
}

function openServiceWorkbench(row: OrderRecord) {
  void router.push({
    path: "/services/list",
    query: {
      orderId: String(row.id),
      customerId: String(row.customerId)
    }
  });
}

function openFinanceWorkbench(row: OrderRecord) {
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId: String(row.customerId)
    }
  });
}

function openProviderWorkbench(row: OrderRecord) {
  void router.push({
    path: "/providers/accounts",
    query: {
      providerType: row.automationType || undefined,
      accountId: row.providerAccountId ? String(row.providerAccountId) : undefined
    }
  });
}

function openAutomationWorkbench(row: OrderRecord) {
  void router.push({
    path: "/providers/automation",
    query: {
      orderId: String(row.id),
      channel: row.automationType || undefined
    }
  });
}

function openOrderRequestsWorkbench(orderId?: number) {
  void router.push({
    path: "/orders/requests",
    query: {
      orderId: orderId ? String(orderId) : undefined
    }
  });
}

function openQuickAction(status: LifecycleStatus, rows?: OrderRecord[]) {
  const targets = rows && rows.length > 0 ? rows : selectedRows.value;
  if (targets.length === 0) {
    ElMessage.info("请先选择订单");
    return;
  }
  quickActionRows.value = [...targets];
  quickActionForm.status = status;
  quickActionForm.reason = defaultReasonForStatus(status);
  quickActionVisible.value = true;
}

async function submitQuickAction() {
  if (quickActionRows.value.length === 0) return;
  quickActionLoading.value = true;
  let success = 0;
  const failed: OrderRecord[] = [];
  try {
    for (const row of quickActionRows.value) {
      try {
        await updatePendingOrder(row.id, {
          status: quickActionForm.status,
          reason: quickActionForm.reason.trim() || defaultReasonForStatus(quickActionForm.status)
        });
        success += 1;
      } catch {
        failed.push(row);
      }
    }
    if (success > 0) {
      ElMessage.success(
        failed.length > 0
          ? `已处理 ${success} 条订单，失败 ${failed.length} 条`
          : `已将 ${success} 条订单更新为${quickActionStatusLabel.value}`
      );
    } else {
      ElMessage.error("订单状态处理失败");
    }
    if (failed.length > 0) {
      quickActionRows.value = failed;
      quickActionForm.reason = `${quickActionForm.reason.trim()}\n失败订单：${failed.map(item => item.orderNo).join("、")}`.trim();
    } else {
      quickActionVisible.value = false;
      quickActionRows.value = [];
      selectedRows.value = [];
      await loadOrders();
    }
  } finally {
    quickActionLoading.value = false;
    if (failed.length > 0) {
      await loadOrders();
    }
  }
}

function closeQuickAction() {
  quickActionVisible.value = false;
  quickActionRows.value = [];
  quickActionForm.reason = "";
}

function handleRowAction(row: OrderRecord, command: string) {
  switch (command) {
    case "customer":
      openCustomerWorkbench(row);
      return;
    case "invoice":
      openInvoiceWorkbench(row);
      return;
    case "service":
      openServiceWorkbench(row);
      return;
    case "finance":
      openFinanceWorkbench(row);
      return;
    case "provider":
      openProviderWorkbench(row);
      return;
    case "tasks":
      openAutomationWorkbench(row);
      return;
    case "requests":
      openOrderRequestsWorkbench(row.id);
      return;
    case "pending":
      openQuickAction("PENDING", [row]);
      return;
    case "active":
      openQuickAction("ACTIVE", [row]);
      return;
    case "completed":
      openQuickAction("COMPLETED", [row]);
      return;
    case "cancelled":
      openQuickAction("CANCELLED", [row]);
      return;
    default:
      return;
  }
}

function handleRowDropdownCommand(row: OrderRecord, command: string | number | object) {
  handleRowAction(row, String(command));
}

async function copySelectedOrderNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择订单");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.orderNo).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个订单编号`);
}

function exportRows(rows: OrderRecord[], filename: string) {
  downloadCsv(
    filename,
    [
      "订单编号",
      "客户 ID",
      "客户名称",
      "商品名称",
      "自动化渠道",
      "接口账户",
      "支付方式",
      "支付状态",
      "计费周期",
      "金额",
      "订单状态",
      "创建时间"
    ],
    rows.map(item => [
      item.orderNo,
      String(item.customerId),
      item.customerName,
      item.productName,
      automationLabel(item.automationType),
      providerAccountLabel(item.providerAccountId),
      paymentLabel(item.payment),
      payStatusLabel(item.payStatus),
      formatCycle(item.billingCycle),
      formatCurrency(item.amount),
      statusLabel(item.status),
      item.createdAt
    ])
  );
}

function exportCurrent() {
  exportRows(orders.value, `orders-page-${pagination.page}.csv`);
  ElMessage.success("当前页订单已导出");
}

function exportSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择订单");
    return;
  }
  exportRows(selectedRows.value, "orders-selected.csv");
  ElMessage.success("已导出选中订单");
}

function openDetail(row: OrderRecord) {
  void router.push(`/orders/detail/${row.id}`);
}

function readRouteQueryValue(value: unknown) {
  if (Array.isArray(value)) return String(value[0] ?? "");
  if (value === undefined || value === null) return "";
  return String(value);
}

function syncFiltersFromRoute() {
  const status = readRouteQueryValue(route.query.status).toUpperCase();
  activeTab.value = (status as TabKey) || "ALL";
  filters.orderNo = readRouteQueryValue(route.query.orderNo || route.query.ordernum);
  filters.productName = readRouteQueryValue(route.query.productName || route.query.product_name);
  filters.customerId = readRouteQueryValue(route.query.customerId || route.query.uid);
  filters.amount = readRouteQueryValue(route.query.amount);
  filters.payment = readRouteQueryValue(route.query.payment).toUpperCase();
  filters.payStatus = readRouteQueryValue(route.query.payStatus || route.query.pay_status).toUpperCase();
  const startTime = readRouteQueryValue(route.query.start_time);
  const endTime = readRouteQueryValue(route.query.end_time);
  filters.dateRange = startTime && endTime ? [startTime, endTime] : [];
  pagination.page = 1;
}

onMounted(() => {
  syncFiltersFromRoute();
  void loadOrders();
});

watch(activeTab, () => {
  pagination.page = 1;
  void loadOrders();
});

watch(
  () => route.fullPath,
  () => {
    syncFiltersFromRoute();
    void loadOrders();
  }
);
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 订单"
      title="产品订单"
      subtitle="按照魔方财务的订单列表逻辑，统一查看订单、支付、自动化渠道、接口账户和后台处置动作。"
    >
      <template #actions>
        <el-button @click="loadOrders">刷新列表</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill">
            <span>订单总数</span>
            <strong>{{ total }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前页金额</span>
            <strong>{{ formatCurrency(totalAmount) }}</strong>
          </div>
          <div class="summary-pill">
            <span>已支付</span>
            <strong>{{ paidCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>未支付</span>
            <strong>{{ unpaidCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>已退款</span>
            <strong>{{ refundedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>自动化订单</span>
            <strong>{{ automatedCount }}</strong>
          </div>
          <div class="summary-pill">
            <span>待处理申请</span>
            <strong>{{ pendingRequestCount }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="filters.orderNo" placeholder="订单号" clearable />
          <el-input v-model="filters.productName" placeholder="商品名称" clearable />
          <el-input v-model="filters.customerId" placeholder="客户 ID" clearable />
          <div class="action-group">
            <el-button type="primary" @click="applyFilters">应用筛选</el-button>
            <el-button plain @click="advancedVisible = !advancedVisible">
              {{ advancedVisible ? "收起高级筛选" : "展开高级筛选" }}
            </el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div v-if="advancedVisible" class="filter-bar filter-bar--compact">
          <el-select v-model="filters.payment" placeholder="支付方式" clearable>
            <el-option v-for="item in paymentOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-select v-model="filters.payStatus" placeholder="支付状态" clearable>
            <el-option v-for="item in payStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
          <el-input v-model="filters.amount" placeholder="订单金额" clearable />
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>运营列表</strong>
          <span>当前页 {{ orders.length }} 条</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain @click="openOrderRequestsWorkbench()">订单申请</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openQuickAction('PENDING')">转待支付</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openQuickAction('ACTIVE')">转待开通</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openQuickAction('COMPLETED')">标记完成</el-button>
          <el-button plain :disabled="selectedRows.length === 0" @click="openQuickAction('CANCELLED')">标记取消</el-button>
          <el-button plain :disabled="selectedRows.length !== 1" @click="selectedRows[0] && openInvoiceWorkbench(selectedRows[0])">
            账单联查
          </el-button>
          <el-button plain :disabled="selectedRows.length !== 1" @click="selectedRows[0] && openServiceWorkbench(selectedRows[0])">
            服务联查
          </el-button>
          <el-button plain @click="copySelectedOrderNos">复制订单号</el-button>
          <el-button plain @click="exportSelected">导出选中</el-button>
          <el-button plain @click="exportCurrent">导出当前页</el-button>
        </div>
      </div>

      <el-table
        :data="orders"
        border
        stripe
        empty-text="暂无符合条件的订单"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column prop="orderNo" label="订单编号" min-width="170" />
        <el-table-column prop="customerId" label="客户 ID" min-width="90" />
        <el-table-column label="客户名称" min-width="170">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomerWorkbench(row)">
              {{ row.customerName }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="productName" label="商品名称" min-width="240" show-overflow-tooltip />
        <el-table-column label="自动化渠道" min-width="140">
          <template #default="{ row }">{{ automationLabel(row.automationType) }}</template>
        </el-table-column>
        <el-table-column label="接口账户" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">{{ providerAccountLabel(row.providerAccountId) }}</template>
        </el-table-column>
        <el-table-column label="支付方式" min-width="120">
          <template #default="{ row }">{{ paymentLabel(row.payment) }}</template>
        </el-table-column>
        <el-table-column label="支付状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="payStatusTagType(row.payStatus)" effect="light">{{ payStatusLabel(row.payStatus) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="计费周期" min-width="100">
          <template #default="{ row }">{{ formatCycle(row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column label="订单金额" min-width="120">
          <template #default="{ row }">{{ formatCurrency(row.amount) }}</template>
        </el-table-column>
        <el-table-column label="订单状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" min-width="180" />
        <el-table-column label="操作" min-width="320" fixed="right">
          <template #default="{ row }">
            <div class="inline-actions">
              <el-button type="primary" link @click="openDetail(row)">进入工作台</el-button>
              <el-button type="primary" link @click="openCustomerWorkbench(row)">客户</el-button>
              <el-button type="primary" link @click="openInvoiceWorkbench(row)">账单</el-button>
              <el-button type="primary" link @click="openServiceWorkbench(row)">服务</el-button>
              <el-button type="primary" link @click="openOrderRequestsWorkbench(row.id)">申请</el-button>
              <el-dropdown @command="handleRowDropdownCommand(row, $event)">
                <el-button type="primary" link>
                  更多
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item command="invoice">账单联查</el-dropdown-item>
                    <el-dropdown-item command="service">服务联查</el-dropdown-item>
                    <el-dropdown-item command="finance">客户财务</el-dropdown-item>
                    <el-dropdown-item command="provider" :disabled="!row.providerAccountId">接口账户</el-dropdown-item>
                    <el-dropdown-item command="tasks">任务中心</el-dropdown-item>
                    <el-dropdown-item command="requests">订单申请</el-dropdown-item>
                    <el-dropdown-item command="pending">转待支付</el-dropdown-item>
                    <el-dropdown-item command="active">转待开通</el-dropdown-item>
                    <el-dropdown-item command="completed">标记完成</el-dropdown-item>
                    <el-dropdown-item command="cancelled">标记取消</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
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
          @current-change="handlePageChange"
          @size-change="handleLimitChange"
        />
      </div>
    </PageWorkbench>

    <el-dialog v-model="quickActionVisible" title="订单批量处置" width="520px" @closed="closeQuickAction">
      <el-form label-position="top">
        <el-alert
          :title="`即将把 ${quickActionTargetCount} 条订单更新为${quickActionStatusLabel}`"
          description="这个动作会调用订单人工调整链路，并写入审计记录；如果订单已挂服务，后续交付状态也会跟着回写。"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="目标状态">
          <el-select v-model="quickActionForm.status" style="width: 100%">
            <el-option v-for="item in lifecycleActionOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="变更原因">
          <el-input
            v-model="quickActionForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请输入后台处理原因，客户改单、人工收款、取消关闭等都应记录清楚"
          />
        </el-form-item>
        <el-form-item label="本次涉及订单">
          <div class="table-toolbar__meta" style="line-height: 1.8">
            <span v-for="item in quickActionRows" :key="item.id">{{ item.orderNo }} / {{ item.customerName }}</span>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeQuickAction">取消</el-button>
        <el-button type="primary" :loading="quickActionLoading" @click="submitQuickAction">确认处置</el-button>
      </template>
    </el-dialog>
  </div>
</template>

