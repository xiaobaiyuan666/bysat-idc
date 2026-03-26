<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { loadOrders, type PortalOrder } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalMoney,
  formatPortalOrderStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();
const loading = ref(true);
const error = ref("");
const orders = ref<PortalOrder[]>([]);
const keyword = ref("");
const status = ref("");

const filteredOrders = computed(() =>
  orders.value.filter(item => {
    const matchedKeyword =
      !keyword.value || item.orderNo.includes(keyword.value) || item.productName.includes(keyword.value);
    const matchedStatus = !status.value || item.status === status.value;
    return matchedKeyword && matchedStatus;
  })
);

const pendingCount = computed(() => orders.value.filter(item => item.status === "PENDING").length);
const activeCount = computed(
  () => orders.value.filter(item => item.status === "ACTIVE" || item.status === "COMPLETED").length
);

const quickActions = computed(() => [
  { label: pickLabel(localeStore.locale, "前往商城", "Open Store"), path: "/store", type: "primary" as const },
  { label: pickLabel(localeStore.locale, "处理账单", "Invoices"), path: "/invoices", type: "danger" as const },
  { label: pickLabel(localeStore.locale, "查看服务", "Services"), path: "/services", type: "info" as const }
]);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "订单中心", "Orders"),
  title: pickLabel(localeStore.locale, "订单中心", "Orders"),
  subtitle: pickLabel(
    localeStore.locale,
    "查看下单记录、计费周期和订单状态，便于对照账单与服务。",
    "Review order history, billing cycles, and fulfillment status."
  ),
  total: pickLabel(localeStore.locale, "订单总数", "Total Orders"),
  pending: pickLabel(localeStore.locale, "待支付", "Pending"),
  active: pickLabel(localeStore.locale, "已生效", "Active"),
  listTitle: pickLabel(localeStore.locale, "订单列表", "Order List"),
  listDesc: pickLabel(localeStore.locale, "支持按编号、产品名称和状态筛选。", "Filter by order number, product name, and status."),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索订单编号或产品名称", "Search order no. or product"),
  statusPlaceholder: pickLabel(localeStore.locale, "订单状态", "Order Status"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  orderNo: pickLabel(localeStore.locale, "订单编号", "Order No."),
  productName: pickLabel(localeStore.locale, "产品名称", "Product"),
  cycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  createdAt: pickLabel(localeStore.locale, "创建时间", "Created At"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  openInvoices: pickLabel(localeStore.locale, "关联账单", "Invoices"),
  empty: pickLabel(localeStore.locale, "暂无匹配的订单记录。", "No matching orders."),
  loadError: pickLabel(localeStore.locale, "订单列表加载失败", "Failed to load orders")
}));

async function fetchOrders() {
  loading.value = true;
  error.value = "";
  try {
    orders.value = await loadOrders();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

function applyRouteFilters() {
  const orderNo = typeof route.query.orderNo === "string" ? route.query.orderNo : "";
  const keywordQuery = typeof route.query.keyword === "string" ? route.query.keyword : "";
  const statusQuery = typeof route.query.status === "string" ? route.query.status : "";
  keyword.value = orderNo || keywordQuery;
  status.value = statusQuery;
}

onMounted(() => {
  applyRouteFilters();
  void fetchOrders();
});

watch(
  () => [route.query.orderNo, route.query.keyword, route.query.status],
  () => {
    applyRouteFilters();
  }
);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <div class="portal-badge">{{ copy.badge }}</div>
          <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="portal-toolbar" style="margin-top: 18px">
        <el-button
          v-for="item in quickActions"
          :key="item.path"
          :type="item.type"
          plain
          @click="router.push(item.path)"
        >
          {{ item.label }}
        </el-button>
      </div>

      <div class="portal-grid portal-grid--three" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.total }}</h3>
          <strong>{{ orders.length }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.pending }}</h3>
          <strong>{{ pendingCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.active }}</h3>
          <strong>{{ activeCount }}</strong>
        </div>
      </div>
    </section>

    <section class="portal-card portal-table-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.listTitle }}</h2>
          <div class="portal-panel__meta">{{ copy.listDesc }}</div>
        </div>
      </div>

      <div class="portal-toolbar" style="margin: 18px 20px 0">
        <el-input v-model="keyword" :placeholder="copy.keywordPlaceholder" clearable />
        <el-select v-model="status" :placeholder="copy.statusPlaceholder" clearable>
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'PENDING')" value="PENDING" />
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'ACTIVE')" value="ACTIVE" />
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'COMPLETED')" value="COMPLETED" />
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'CANCELLED')" value="CANCELLED" />
        </el-select>
        <div class="portal-summary" style="margin: 0; padding: 12px 16px">
          <div class="portal-summary-row" style="font-size: 12px">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredOrders.length }}</strong>
          </div>
        </div>
      </div>

      <el-table v-if="filteredOrders.length" :data="filteredOrders" border>
        <el-table-column prop="orderNo" :label="copy.orderNo" min-width="170" />
        <el-table-column prop="productName" :label="copy.productName" min-width="220" />
        <el-table-column :label="copy.cycle" min-width="120">
          <template #default="{ row }">{{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
        </el-table-column>
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalOrderStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="copy.createdAt" min-width="180" />
        <el-table-column :label="copy.action" min-width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="router.push(`/invoices?keyword=${row.orderNo}`)">
              {{ copy.openInvoices }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>
  </div>
</template>
