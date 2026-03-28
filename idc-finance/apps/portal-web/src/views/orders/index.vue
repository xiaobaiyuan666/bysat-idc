<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import {
  createOrderRequest,
  loadOrderDetail,
  loadOrders,
  type CreatePortalOrderRequestPayload,
  type PortalOrder,
  type PortalOrderDetailResponse
} from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalInvoiceStatus,
  formatPortalMoney,
  formatPortalOrderStatus,
  formatPortalServiceStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();

const loading = ref(false);
const detailLoading = ref(false);
const requestLoading = ref(false);
const error = ref("");
const orders = ref<PortalOrder[]>([]);
const detail = ref<PortalOrderDetailResponse | null>(null);
const detailVisible = ref(false);
const requestVisible = ref(false);
const isRouteDetail = computed(() => Number.isFinite(Number(route.params.id)) && Number(route.params.id) > 0);
const detailSize = computed(() => (isRouteDetail.value ? "calc(100vw - 280px)" : "72%"));

const filters = reactive({
  keyword: "",
  status: ""
});

const requestForm = reactive<CreatePortalOrderRequestPayload>({
  type: "RENEW",
  reason: "",
  requestedBillingCycle: "",
  summary: ""
});

const filteredOrders = computed(() =>
  orders.value.filter(item => {
    const text = filters.keyword.trim().toLowerCase();
    const matchesKeyword =
      !text || item.orderNo.toLowerCase().includes(text) || item.productName.toLowerCase().includes(text);
    const matchesStatus = !filters.status || item.status === filters.status;
    return matchesKeyword && matchesStatus;
  })
);

const stats = computed(() => ({
  total: orders.value.length,
  pending: orders.value.filter(item => item.status === "PENDING").length,
  active: orders.value.filter(item => item.status === "ACTIVE" || item.status === "COMPLETED").length,
  requests: detail.value?.requests.length ?? 0
}));

const pendingAmount = computed(() =>
  orders.value.filter(item => item.status === "PENDING").reduce((sum, item) => sum + item.amount, 0)
);

const latestOrders = computed(() => filteredOrders.value.slice(0, 5));

const detailStats = computed(() => ({
  invoiceCount: detail.value?.invoices.length ?? 0,
  serviceCount: detail.value?.services.length ?? 0,
  requestCount: detail.value?.requests.length ?? 0,
  pendingRequestCount: detail.value?.requests.filter(item => item.status === "PENDING").length ?? 0
}));

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "订单中心", "Orders"),
  title: pickLabel(localeStore.locale, "订单中心", "Orders"),
  subtitle: pickLabel(
    localeStore.locale,
    "统一查看订单生命周期、关联账单与服务，并在门户侧提交续费、取消或人工调整申请。",
    "Review order lifecycle, linked invoices and services, and submit renew, cancel, or adjustment requests."
  ),
  total: pickLabel(localeStore.locale, "订单总数", "Total Orders"),
  pending: pickLabel(localeStore.locale, "待支付", "Pending"),
  active: pickLabel(localeStore.locale, "已生效", "Active"),
  requests: pickLabel(localeStore.locale, "当前申请数", "Request Count"),
  pendingAmount: pickLabel(localeStore.locale, "待支付金额", "Pending Amount"),
  latest: pickLabel(localeStore.locale, "最近订单", "Recent Orders"),
  latestDesc: pickLabel(localeStore.locale, "优先处理待支付和新近生成的订单。", "Handle pending and newly created orders first."),
  requestDesk: pickLabel(localeStore.locale, "申请工作台", "Request Desk"),
  requestDeskDesc: pickLabel(localeStore.locale, "从订单视角统一处理续费、取消和调整申请。", "Handle renew, cancel, and adjustment requests from the order perspective."),
  keyword: pickLabel(localeStore.locale, "搜索订单号 / 商品", "Search order or product"),
  status: pickLabel(localeStore.locale, "订单状态", "Status"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  orderNo: pickLabel(localeStore.locale, "订单号", "Order No."),
  productName: pickLabel(localeStore.locale, "商品", "Product"),
  cycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  statusCol: pickLabel(localeStore.locale, "状态", "Status"),
  createdAt: pickLabel(localeStore.locale, "创建时间", "Created At"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  detail: pickLabel(localeStore.locale, "详情", "Detail"),
  invoices: pickLabel(localeStore.locale, "账单", "Invoices"),
  tickets: pickLabel(localeStore.locale, "工单", "Tickets"),
  request: pickLabel(localeStore.locale, "提交申请", "Submit Request"),
  empty: pickLabel(localeStore.locale, "暂无匹配的订单记录。", "No matching orders."),
  detailTitle: pickLabel(localeStore.locale, "订单详情", "Order Detail"),
  backToList: pickLabel(localeStore.locale, "返回订单列表", "Back to Orders"),
  linkedInvoices: pickLabel(localeStore.locale, "关联账单", "Linked Invoices"),
  linkedServices: pickLabel(localeStore.locale, "已开通服务", "Provisioned Services"),
  requestRecords: pickLabel(localeStore.locale, "订单申请", "Order Requests"),
  linkedSummary: pickLabel(localeStore.locale, "关联摘要", "Linked Summary"),
  createRequestTitle: pickLabel(localeStore.locale, "提交订单申请", "Create Order Request"),
  requestType: pickLabel(localeStore.locale, "申请类型", "Request Type"),
  requestReason: pickLabel(localeStore.locale, "申请原因", "Reason"),
  requestedCycle: pickLabel(localeStore.locale, "期望周期", "Requested Cycle"),
  summary: pickLabel(localeStore.locale, "摘要", "Summary"),
  submit: pickLabel(localeStore.locale, "提交", "Submit"),
  cancel: pickLabel(localeStore.locale, "取消", "Cancel"),
  noRequest: pickLabel(localeStore.locale, "暂无申请记录", "No requests"),
  invoiceStatus: pickLabel(localeStore.locale, "账单状态", "Invoice Status"),
  serviceStatus: pickLabel(localeStore.locale, "服务状态", "Service Status"),
  openInvoice: pickLabel(localeStore.locale, "打开账单", "Open Invoice"),
  openService: pickLabel(localeStore.locale, "打开服务", "Open Service"),
  openTicket: pickLabel(localeStore.locale, "提交工单", "Create Ticket"),
  pendingRequests: pickLabel(localeStore.locale, "待处理申请", "Pending Requests"),
  renew: pickLabel(localeStore.locale, "续费", "Renew"),
  cancelType: pickLabel(localeStore.locale, "取消", "Cancel"),
  adjustType: pickLabel(localeStore.locale, "调整", "Adjustment"),
  requestSuccess: pickLabel(localeStore.locale, "订单申请已提交", "Order request submitted"),
  loadError: pickLabel(localeStore.locale, "订单列表加载失败", "Failed to load orders")
}));

function requestTypeLabel(type: string) {
  const mapping: Record<string, string> = {
    RENEW: copy.value.renew,
    CANCEL: copy.value.cancelType,
    PRICE_ADJUST: copy.value.adjustType
  };
  return mapping[type] ?? type;
}

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
  const keyword = typeof route.query.orderNo === "string" ? route.query.orderNo : "";
  const status = typeof route.query.status === "string" ? route.query.status : "";
  filters.keyword = keyword;
  filters.status = status;
}

async function openDetail(orderId: number, replacePath = true) {
  detailLoading.value = true;
  detailVisible.value = true;
  try {
    detail.value = await loadOrderDetail(orderId);
    if (replacePath && route.path !== `/orders/${orderId}`) {
      void router.replace(`/orders/${orderId}`);
    }
  } finally {
    detailLoading.value = false;
  }
}

function closeDetail() {
  detailVisible.value = false;
  detail.value = null;
  if (route.path !== "/orders") {
    void router.replace("/orders");
  }
}

function openTickets(orderNo: string) {
  void router.push({ path: "/tickets", query: { keyword: orderNo } });
}

function openTicketForService() {
  const service = detail.value?.services[0];
  if (!service) {
    openTickets(detail.value?.order.orderNo || "");
    return;
  }
  void router.push({
    path: "/tickets",
    query: {
      action: "create",
      serviceId: String(service.id),
      title: detail.value?.order.productName || undefined
    }
  });
}

function openInvoices(orderNo: string) {
  void router.push({ path: "/invoices", query: { orderNo } });
}

function openRequestDialog(order: PortalOrder) {
  requestForm.type = "RENEW";
  requestForm.reason = "";
  requestForm.requestedBillingCycle = order.billingCycle;
  requestForm.summary = `${order.productName} 续费/调整申请`;
  detail.value = detail.value?.order.id === order.id
    ? detail.value
    : {
        order,
        invoices: [],
        services: [],
        requests: [],
        auditLogs: []
      };
  requestVisible.value = true;
}

async function submitRequest() {
  if (!detail.value) return;
  requestLoading.value = true;
  try {
    detail.value = await createOrderRequest(detail.value.order.id, requestForm);
    requestVisible.value = false;
    detailVisible.value = true;
    ElMessage.success(copy.value.requestSuccess);
  } finally {
    requestLoading.value = false;
  }
}

async function syncDetailFromRoute() {
  const rawId = Number(route.params.id);
  if (Number.isFinite(rawId) && rawId > 0) {
    await openDetail(rawId, false);
    return;
  }
  detailVisible.value = false;
  detail.value = null;
}

onMounted(async () => {
  applyRouteFilters();
  await fetchOrders();
  await syncDetailFromRoute();
});

watch(
  () => [route.query.orderNo, route.query.status],
  () => {
    applyRouteFilters();
  }
);

watch(
  () => route.params.id,
  async () => {
    await syncDetailFromRoute();
  }
);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <section v-if="!isRouteDetail" class="portal-card">
      <div class="portal-card-head">
        <div>
          <div class="portal-badge">{{ copy.badge }}</div>
          <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
          <p class="portal-subtitle">{{ copy.subtitle }}</p>
        </div>
      </div>

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <article class="portal-stat"><h3>{{ copy.total }}</h3><strong>{{ stats.total }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.pending }}</h3><strong>{{ stats.pending }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.active }}</h3><strong>{{ stats.active }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.pendingAmount }}</h3><strong>{{ formatPortalMoney(localeStore.locale, pendingAmount) }}</strong></article>
      </div>

      <section class="portal-grid portal-grid--two" style="margin-top: 20px">
        <article class="portal-card" style="padding: 0">
          <div style="padding: 20px 20px 0">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.latest }}</h2>
                <div class="portal-panel__meta">{{ copy.latestDesc }}</div>
              </div>
            </div>
          </div>
          <div v-if="latestOrders.length" class="portal-list" style="padding: 18px 20px 20px">
            <div v-for="item in latestOrders" :key="item.id" class="portal-list-item">
              <div class="portal-list-item__meta">
                <div class="portal-list-item__title">{{ item.productName }}</div>
                <div class="portal-list-item__desc">{{ copy.orderNo }}：{{ item.orderNo }}</div>
                <div class="portal-list-item__desc">{{ copy.cycle }}：{{ formatPortalBillingCycle(localeStore.locale, item.billingCycle) }}</div>
              </div>
              <div class="portal-toolbar">
                <el-tag :type="portalTagTypeByStatus(item.status)">
                  {{ formatPortalOrderStatus(localeStore.locale, item.status) }}
                </el-tag>
                <strong>{{ formatPortalMoney(localeStore.locale, item.amount) }}</strong>
                <el-button type="primary" link @click="openDetail(item.id)">{{ copy.detail }}</el-button>
              </div>
            </div>
          </div>
          <div v-else class="portal-empty-state" style="margin: 18px 20px 20px">{{ copy.empty }}</div>
        </article>

        <article class="portal-card" style="padding: 0">
          <div style="padding: 20px 20px 0">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.requestDesk }}</h2>
                <div class="portal-panel__meta">{{ copy.requestDeskDesc }}</div>
              </div>
            </div>
          </div>
          <div class="portal-summary" style="margin: 18px 20px 20px">
            <div class="portal-summary-row">
              <span>{{ copy.pending }}</span>
              <strong>{{ stats.pending }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.requests }}</span>
              <strong>{{ stats.requests }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.visible }}</span>
              <strong>{{ filteredOrders.length }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.active }}</span>
              <strong>{{ stats.active }}</strong>
            </div>
          </div>
        </article>
      </section>

      <div class="portal-toolbar" style="margin-top: 20px">
        <el-input v-model="filters.keyword" :placeholder="copy.keyword" clearable />
        <el-select v-model="filters.status" :placeholder="copy.status" clearable style="width: 160px">
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'PENDING')" value="PENDING" />
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'ACTIVE')" value="ACTIVE" />
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'COMPLETED')" value="COMPLETED" />
          <el-option :label="formatPortalOrderStatus(localeStore.locale, 'CANCELLED')" value="CANCELLED" />
        </el-select>
        <div class="portal-summary" style="margin-left: auto; min-width: 180px">
          <div class="portal-summary-row">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredOrders.length }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section v-if="!isRouteDetail" class="portal-table-card">
      <el-table v-if="filteredOrders.length" :data="filteredOrders" border>
        <el-table-column prop="orderNo" :label="copy.orderNo" min-width="160" />
        <el-table-column prop="productName" :label="copy.productName" min-width="220" />
        <el-table-column :label="copy.cycle" min-width="120">
          <template #default="{ row }">{{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
        </el-table-column>
        <el-table-column :label="copy.statusCol" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalOrderStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" :label="copy.createdAt" min-width="160" />
        <el-table-column :label="copy.action" min-width="260" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row.id)">{{ copy.detail }}</el-button>
            <el-button link @click="openInvoices(row.orderNo)">{{ copy.invoices }}</el-button>
            <el-button link @click="openTickets(row.orderNo)">{{ copy.tickets }}</el-button>
            <el-button link @click="openRequestDialog(row)">{{ copy.request }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state">{{ copy.empty }}</div>
    </section>

    <el-drawer
      v-model="detailVisible"
      :class="{ 'portal-detail-drawer--page': isRouteDetail }"
      :size="detailSize"
      :modal="!isRouteDetail"
      :close-on-click-modal="!isRouteDetail"
      :lock-scroll="!isRouteDetail"
      :with-header="!isRouteDetail"
      :show-close="!isRouteDetail"
      :title="copy.detailTitle"
      @close="closeDetail"
    >
      <div v-loading="detailLoading" class="portal-shell" v-if="detail">
        <section class="portal-card">
          <div class="portal-card-head" style="margin-bottom: 18px">
            <div>
              <h2 class="portal-panel__title">{{ detail.order.orderNo }}</h2>
              <div class="portal-panel__meta">{{ detail.order.productName }}</div>
            </div>
            <div class="portal-toolbar">
              <el-button plain @click="closeDetail">{{ copy.backToList }}</el-button>
              <el-button
                v-if="detail.invoices[0]"
                plain
                @click="router.push(`/invoices/${detail.invoices[0].id}`)"
              >
                {{ copy.openInvoice }}
              </el-button>
              <el-button
                v-if="detail.services[0]"
                plain
                @click="router.push(`/services/${detail.services[0].id}`)"
              >
                {{ copy.openService }}
              </el-button>
              <el-button plain @click="openTicketForService">{{ copy.openTicket }}</el-button>
              <el-button type="primary" plain @click="requestVisible = true">{{ copy.request }}</el-button>
            </div>
          </div>
          <div class="portal-grid portal-grid--four">
            <article class="portal-stat"><h3>{{ copy.orderNo }}</h3><strong>{{ detail.order.orderNo }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.productName }}</h3><strong>{{ detail.order.productName }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.amount }}</h3><strong>{{ formatPortalMoney(localeStore.locale, detail.order.amount) }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.statusCol }}</h3><strong>{{ formatPortalOrderStatus(localeStore.locale, detail.order.status) }}</strong></article>
          </div>
        </section>

        <section class="portal-card">
          <div class="portal-card-head">
            <h2 class="portal-panel__title">{{ copy.linkedSummary }}</h2>
          </div>
          <div class="portal-grid portal-grid--four" style="margin-top: 18px">
            <article class="portal-stat"><h3>{{ copy.linkedInvoices }}</h3><strong>{{ detailStats.invoiceCount }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.linkedServices }}</h3><strong>{{ detailStats.serviceCount }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.requestRecords }}</h3><strong>{{ detailStats.requestCount }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.pendingRequests }}</h3><strong>{{ detailStats.pendingRequestCount }}</strong></article>
          </div>
        </section>

        <section class="portal-grid portal-grid--two">
          <article class="portal-card">
            <div class="portal-card-head"><h2 class="portal-panel__title">{{ copy.linkedInvoices }}</h2></div>
            <el-table v-if="detail.invoices.length" :data="detail.invoices" border style="margin-top: 18px">
              <el-table-column prop="invoiceNo" :label="copy.invoices" min-width="140" />
              <el-table-column :label="copy.amount" min-width="120">
                <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.totalAmount) }}</template>
              </el-table-column>
              <el-table-column :label="copy.invoiceStatus" min-width="120">
                <template #default="{ row }">
                  <el-tag :type="portalTagTypeByStatus(row.status)">
                    {{ formatPortalInvoiceStatus(localeStore.locale, row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column min-width="110" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="router.push(`/invoices/${row.id}`)">{{ copy.openInvoice }}</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noRequest }}</div>
          </article>

          <article class="portal-card">
            <div class="portal-card-head"><h2 class="portal-panel__title">{{ copy.linkedServices }}</h2></div>
            <el-table v-if="detail.services.length" :data="detail.services" border style="margin-top: 18px">
              <el-table-column prop="serviceNo" :label="pickLabel(localeStore.locale, '服务号', 'Service No.')" min-width="140" />
              <el-table-column prop="productName" :label="copy.productName" min-width="180" />
              <el-table-column prop="nextDueAt" :label="pickLabel(localeStore.locale, '到期时间', 'Next Due')" min-width="120" />
              <el-table-column :label="copy.serviceStatus" min-width="120">
                <template #default="{ row }">
                  <el-tag :type="portalTagTypeByStatus(row.status)">
                    {{ formatPortalServiceStatus(localeStore.locale, row.status) }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column min-width="110" fixed="right">
                <template #default="{ row }">
                  <el-button link type="primary" @click="router.push(`/services/${row.id}`)">{{ copy.openService }}</el-button>
                </template>
              </el-table-column>
            </el-table>
            <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noRequest }}</div>
          </article>
        </section>

        <section class="portal-card">
          <div class="portal-card-head">
            <h2 class="portal-panel__title">{{ copy.requestRecords }}</h2>
            <el-button type="primary" plain @click="requestVisible = true">{{ copy.request }}</el-button>
          </div>
          <el-table v-if="detail.requests.length" :data="detail.requests" border style="margin-top: 18px">
            <el-table-column prop="requestNo" :label="pickLabel(localeStore.locale, '申请号', 'Request No.')" min-width="150" />
            <el-table-column :label="copy.requestType" min-width="120">
              <template #default="{ row }">{{ requestTypeLabel(row.type) }}</template>
            </el-table-column>
            <el-table-column prop="summary" :label="copy.summary" min-width="180" />
            <el-table-column prop="status" :label="copy.statusCol" min-width="120" />
            <el-table-column prop="createdAt" :label="copy.createdAt" min-width="160" />
          </el-table>
          <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noRequest }}</div>
        </section>
      </div>
    </el-drawer>

    <el-dialog v-model="requestVisible" :title="copy.createRequestTitle" width="520px">
      <el-form label-position="top">
        <el-form-item :label="copy.requestType">
          <el-select v-model="requestForm.type" style="width: 100%">
            <el-option :label="copy.renew" value="RENEW" />
            <el-option :label="copy.cancelType" value="CANCEL" />
            <el-option :label="copy.adjustType" value="PRICE_ADJUST" />
          </el-select>
        </el-form-item>
        <el-form-item :label="copy.summary">
          <el-input v-model="requestForm.summary" />
        </el-form-item>
        <el-form-item :label="copy.requestedCycle">
          <el-input v-model="requestForm.requestedBillingCycle" />
        </el-form-item>
        <el-form-item :label="copy.requestReason">
          <el-input v-model="requestForm.reason" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="requestVisible = false">{{ copy.cancel }}</el-button>
        <el-button type="primary" :loading="requestLoading" @click="submitRequest">{{ copy.submit }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>
