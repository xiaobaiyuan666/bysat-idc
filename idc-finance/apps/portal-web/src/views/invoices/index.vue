<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import { loadInvoiceDetail, loadInvoices, payInvoice, type PortalInvoice, type PortalInvoiceDetailResponse } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalInvoiceStatus,
  formatPortalMoney,
  formatPortalPaymentChannel,
  formatPortalServiceStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();

const loading = ref(false);
const detailLoading = ref(false);
const payingId = ref<number | null>(null);
const error = ref("");
const invoices = ref<PortalInvoice[]>([]);
const detail = ref<PortalInvoiceDetailResponse | null>(null);
const detailVisible = ref(false);
const isRouteDetail = computed(() => Number.isFinite(Number(route.params.id)) && Number(route.params.id) > 0);
const detailSize = computed(() => (isRouteDetail.value ? "calc(100vw - 280px)" : "72%"));

const filters = reactive({
  keyword: "",
  status: ""
});

const filteredInvoices = computed(() =>
  invoices.value.filter(item => {
    const text = filters.keyword.trim().toLowerCase();
    const matchesKeyword =
      !text ||
      item.invoiceNo.toLowerCase().includes(text) ||
      item.productName.toLowerCase().includes(text) ||
      item.orderNo.toLowerCase().includes(text);
    const matchesStatus = !filters.status || item.status === filters.status;
    return matchesKeyword && matchesStatus;
  })
);

const stats = computed(() => ({
  total: invoices.value.length,
  unpaid: invoices.value.filter(item => item.status === "UNPAID").length,
  paid: invoices.value.filter(item => item.status === "PAID").length,
  refunded: invoices.value.filter(item => item.status === "REFUNDED").length
}));

const unpaidAmount = computed(() =>
  invoices.value.filter(item => item.status === "UNPAID").reduce((sum, item) => sum + item.totalAmount, 0)
);

const latestInvoices = computed(() => filteredInvoices.value.slice(0, 5));

const detailStats = computed(() => ({
  paymentCount: detail.value?.payments.length ?? 0,
  refundCount: detail.value?.refunds.length ?? 0,
  serviceCount: detail.value?.services.length ?? 0
}));
const settlementMetrics = computed(() => [
  {
    label: copy.value.invoiceNo,
    value: detail.value?.invoice.invoiceNo || "-"
  },
  {
    label: copy.value.amount,
    value: formatPortalMoney(localeStore.locale, detail.value?.invoice.totalAmount || 0)
  },
  {
    label: copy.value.statusCol,
    value: formatPortalInvoiceStatus(localeStore.locale, detail.value?.invoice.status)
  },
  {
    label: copy.value.dueAt,
    value: detail.value?.invoice.dueAt || "-"
  }
]);
const financeTimeline = computed(() => {
  const events: Array<{ time: string; title: string; subtitle: string }> = [];
  if (detail.value?.invoice?.dueAt) {
    events.push({
      time: detail.value.invoice.dueAt,
      title: pickLabel(localeStore.locale, "账单到期", "Invoice Due"),
      subtitle: detail.value.invoice.invoiceNo
    });
  }
  for (const payment of detail.value?.payments ?? []) {
    events.push({
      time: payment.paidAt,
      title: pickLabel(localeStore.locale, "支付入账", "Payment Settled"),
      subtitle: `${payment.paymentNo} / ${formatPortalPaymentChannel(localeStore.locale, payment.channel)}`
    });
  }
  for (const refund of detail.value?.refunds ?? []) {
    events.push({
      time: refund.createdAt,
      title: pickLabel(localeStore.locale, "退款记录", "Refund Recorded"),
      subtitle: `${refund.refundNo} / ${refund.reason || "-"}`
    });
  }
  return events
    .filter(item => item.time)
    .sort((a, b) => String(b.time).localeCompare(String(a.time)));
});

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "财务中心", "Finance"),
  title: pickLabel(localeStore.locale, "账单中心", "Invoices"),
  subtitle: pickLabel(
    localeStore.locale,
    "围绕账单查看支付、退款、服务和订单关联，未支付账单可直接在门户内完成收款动作。",
    "Review payments, refunds, services, and orders around each invoice and pay unpaid invoices directly."
  ),
  total: pickLabel(localeStore.locale, "账单总数", "Total Invoices"),
  unpaid: pickLabel(localeStore.locale, "未支付", "Unpaid"),
  paid: pickLabel(localeStore.locale, "已支付", "Paid"),
  refunded: pickLabel(localeStore.locale, "已退款", "Refunded"),
  unpaidAmount: pickLabel(localeStore.locale, "待支付金额", "Outstanding Amount"),
  latest: pickLabel(localeStore.locale, "待处理账单", "Invoices to Handle"),
  latestDesc: pickLabel(localeStore.locale, "优先处理未支付或即将到期账单。", "Prioritize unpaid and near-due invoices."),
  keyword: pickLabel(localeStore.locale, "搜索账单号 / 订单号 / 商品", "Search invoice, order, or product"),
  status: pickLabel(localeStore.locale, "账单状态", "Status"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  invoiceNo: pickLabel(localeStore.locale, "账单号", "Invoice No."),
  productName: pickLabel(localeStore.locale, "商品", "Product"),
  cycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  dueAt: pickLabel(localeStore.locale, "到期时间", "Due At"),
  statusCol: pickLabel(localeStore.locale, "状态", "Status"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  detail: pickLabel(localeStore.locale, "详情", "Detail"),
  pay: pickLabel(localeStore.locale, "立即支付", "Pay Now"),
  order: pickLabel(localeStore.locale, "订单", "Order"),
  tickets: pickLabel(localeStore.locale, "工单", "Tickets"),
  empty: pickLabel(localeStore.locale, "暂无匹配的账单记录。", "No matching invoices."),
  detailTitle: pickLabel(localeStore.locale, "账单详情", "Invoice Detail"),
  backToList: pickLabel(localeStore.locale, "返回账单列表", "Back to Invoices"),
  linkedSummary: pickLabel(localeStore.locale, "财务摘要", "Finance Summary"),
  financeDesk: pickLabel(localeStore.locale, "财务协同", "Finance Desk"),
  financeDeskDesc: pickLabel(localeStore.locale, "围绕当前账单继续处理订单、服务、工单和支付动作。", "Continue handling order, service, ticket, and payment actions around this invoice."),
  financeTimeline: pickLabel(localeStore.locale, "资金时间线", "Finance Timeline"),
  financeTimelineDesc: pickLabel(localeStore.locale, "按时间查看支付、退款和到期节点。", "Review due, payment, and refund events in time order."),
  payments: pickLabel(localeStore.locale, "支付记录", "Payments"),
  refunds: pickLabel(localeStore.locale, "退款记录", "Refunds"),
  services: pickLabel(localeStore.locale, "关联服务", "Services"),
  serviceStatus: pickLabel(localeStore.locale, "服务状态", "Service Status"),
  openService: pickLabel(localeStore.locale, "打开服务", "Open Service"),
  openOrder: pickLabel(localeStore.locale, "打开订单", "Open Order"),
  openTicket: pickLabel(localeStore.locale, "提交工单", "Create Ticket"),
  paymentSignal: pickLabel(localeStore.locale, "支付笔数", "Payments"),
  refundSignal: pickLabel(localeStore.locale, "退款笔数", "Refunds"),
  serviceSignal: pickLabel(localeStore.locale, "关联服务", "Services"),
  payRule: pickLabel(
    localeStore.locale,
    "支付未付账单时，系统会优先尝试余额抵扣；余额不足时再按在线支付渠道完成入账。",
    "The system tries wallet balance first and falls back to online payment when needed."
  ),
  paidSuccess: pickLabel(localeStore.locale, "账单支付成功", "Invoice paid successfully"),
  loadError: pickLabel(localeStore.locale, "账单列表加载失败", "Failed to load invoices")
}));

async function fetchInvoices() {
  loading.value = true;
  error.value = "";
  try {
    invoices.value = await loadInvoices();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

function applyRouteFilters() {
  const invoiceNo = typeof route.query.invoiceNo === "string" ? route.query.invoiceNo : "";
  const orderNo = typeof route.query.orderNo === "string" ? route.query.orderNo : "";
  const status = typeof route.query.status === "string" ? route.query.status : "";
  filters.keyword = invoiceNo || orderNo;
  filters.status = status;
}

async function openDetail(invoiceId: number, replacePath = true) {
  detailLoading.value = true;
  detailVisible.value = true;
  try {
    detail.value = await loadInvoiceDetail(invoiceId);
    if (replacePath && route.path !== `/invoices/${invoiceId}`) {
      void router.replace(`/invoices/${invoiceId}`);
    }
  } finally {
    detailLoading.value = false;
  }
}

function closeDetail() {
  detailVisible.value = false;
  detail.value = null;
  if (route.path !== "/invoices") {
    void router.replace("/invoices");
  }
}

async function handlePay(row: PortalInvoice) {
  payingId.value = row.id;
  try {
    const result = await payInvoice(row.id);
    ElMessage.success(
      pickLabel(
        localeStore.locale,
        `账单 ${result.invoice.invoiceNo} 支付成功，支付渠道：${formatPortalPaymentChannel(localeStore.locale, result.payment?.channel || "")}`,
        `Invoice ${result.invoice.invoiceNo} paid successfully`
      )
    );
    await fetchInvoices();
    await openDetail(row.id);
  } finally {
    payingId.value = null;
  }
}

function openOrder(orderNo?: string) {
  void router.push({ path: "/orders", query: { orderNo: orderNo || undefined } });
}

function openTickets(invoiceNo?: string) {
  void router.push({ path: "/tickets", query: { keyword: invoiceNo || undefined } });
}

function openTicketForService() {
  const service = detail.value?.services[0];
  if (!service) {
    openTickets(detail.value?.invoice.invoiceNo);
    return;
  }
  void router.push({
    path: "/tickets",
    query: {
      action: "create",
      serviceId: String(service.id),
      title: detail.value?.invoice.productName || undefined
    }
  });
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
  await fetchInvoices();
  await syncDetailFromRoute();
});

watch(
  () => [route.query.invoiceNo, route.query.orderNo, route.query.status],
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

      <el-alert :title="copy.payRule" type="info" :closable="false" show-icon style="margin-top: 20px" />

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <article class="portal-stat"><h3>{{ copy.total }}</h3><strong>{{ stats.total }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.unpaid }}</h3><strong>{{ stats.unpaid }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.paid }}</h3><strong>{{ stats.paid }}</strong></article>
        <article class="portal-stat"><h3>{{ copy.unpaidAmount }}</h3><strong>{{ formatPortalMoney(localeStore.locale, unpaidAmount) }}</strong></article>
      </div>

      <section v-if="!isRouteDetail" class="portal-grid portal-grid--two" style="margin-top: 20px">
        <article class="portal-card" style="padding: 0">
          <div style="padding: 20px 20px 0">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.latest }}</h2>
                <div class="portal-panel__meta">{{ copy.latestDesc }}</div>
              </div>
            </div>
          </div>
          <div v-if="latestInvoices.length" class="portal-list" style="padding: 18px 20px 20px">
            <div v-for="item in latestInvoices" :key="item.id" class="portal-list-item">
              <div class="portal-list-item__meta">
                <div class="portal-list-item__title">{{ item.productName }}</div>
                <div class="portal-list-item__desc">{{ copy.invoiceNo }}：{{ item.invoiceNo }}</div>
                <div class="portal-list-item__desc">{{ copy.dueAt }}：{{ item.dueAt || "-" }}</div>
              </div>
              <div class="portal-toolbar">
                <el-tag :type="portalTagTypeByStatus(item.status)">
                  {{ formatPortalInvoiceStatus(localeStore.locale, item.status) }}
                </el-tag>
                <strong>{{ formatPortalMoney(localeStore.locale, item.totalAmount) }}</strong>
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
                <h2 class="portal-panel__title">{{ copy.linkedSummary }}</h2>
                <div class="portal-panel__meta">{{ pickLabel(localeStore.locale, '从账单视角快速掌握支付、退款和服务交付。', 'Quickly understand payments, refunds, and service delivery from the invoice perspective.') }}</div>
              </div>
            </div>
          </div>
          <div class="portal-summary" style="margin: 18px 20px 20px">
            <div class="portal-summary-row">
              <span>{{ copy.unpaid }}</span>
              <strong>{{ stats.unpaid }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.paid }}</span>
              <strong>{{ stats.paid }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.refunded }}</span>
              <strong>{{ stats.refunded }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.visible }}</span>
              <strong>{{ filteredInvoices.length }}</strong>
            </div>
          </div>
        </article>
      </section>

      <div class="portal-toolbar" style="margin-top: 20px">
        <el-input v-model="filters.keyword" :placeholder="copy.keyword" clearable />
        <el-select v-model="filters.status" :placeholder="copy.status" clearable style="width: 180px">
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'UNPAID')" value="UNPAID" />
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'PAID')" value="PAID" />
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'REFUNDED')" value="REFUNDED" />
        </el-select>
        <div class="portal-summary" style="margin-left: auto; min-width: 180px">
          <div class="portal-summary-row">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredInvoices.length }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section v-if="!isRouteDetail" class="portal-table-card">
      <el-table v-if="filteredInvoices.length" :data="filteredInvoices" border>
        <el-table-column prop="invoiceNo" :label="copy.invoiceNo" min-width="160" />
        <el-table-column prop="productName" :label="copy.productName" min-width="220" />
        <el-table-column :label="copy.cycle" min-width="120">
          <template #default="{ row }">{{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.totalAmount) }}</template>
        </el-table-column>
        <el-table-column prop="dueAt" :label="copy.dueAt" min-width="140" />
        <el-table-column :label="copy.statusCol" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalInvoiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.action" min-width="260" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="openDetail(row.id)">{{ copy.detail }}</el-button>
            <el-button link @click="openOrder(row.orderNo)">{{ copy.order }}</el-button>
            <el-button link @click="openTickets(row.invoiceNo)">{{ copy.tickets }}</el-button>
            <el-button
              v-if="row.status === 'UNPAID'"
              link
              type="danger"
              :loading="payingId === row.id"
              @click="handlePay(row)"
            >
              {{ copy.pay }}
            </el-button>
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
              <h2 class="portal-panel__title">{{ detail.invoice.invoiceNo }}</h2>
              <div class="portal-panel__meta">{{ detail.invoice.productName }}</div>
            </div>
            <div class="portal-toolbar">
              <el-button plain @click="closeDetail">{{ copy.backToList }}</el-button>
              <el-button v-if="detail.order" plain @click="router.push(`/orders/${detail.order.id}`)">
                {{ copy.openOrder }}
              </el-button>
              <el-button v-if="detail.services[0]" plain @click="router.push(`/services/${detail.services[0].id}`)">
                {{ copy.openService }}
              </el-button>
              <el-button plain @click="openTicketForService">{{ copy.openTicket }}</el-button>
              <el-button
                v-if="detail.invoice.status === 'UNPAID'"
                type="primary"
                plain
                :loading="payingId === detail.invoice.id"
                @click="handlePay(detail.invoice)"
              >
                {{ copy.pay }}
              </el-button>
            </div>
          </div>
          <div class="portal-hero-grid">
            <div v-for="item in settlementMetrics" :key="item.label" class="portal-mini-card">
              <span class="portal-mini-card__label">{{ item.label }}</span>
              <strong class="portal-mini-card__value">{{ item.value }}</strong>
            </div>
          </div>
        </section>

        <section class="portal-grid portal-grid--two">
          <article class="portal-card">
            <div class="portal-card-head">
              <h2 class="portal-panel__title">{{ copy.linkedSummary }}</h2>
            </div>
            <div class="portal-grid portal-grid--three" style="margin-top: 18px">
              <article class="portal-stat"><h3>{{ copy.paymentSignal }}</h3><strong>{{ detailStats.paymentCount }}</strong></article>
              <article class="portal-stat"><h3>{{ copy.refundSignal }}</h3><strong>{{ detailStats.refundCount }}</strong></article>
              <article class="portal-stat"><h3>{{ copy.serviceSignal }}</h3><strong>{{ detailStats.serviceCount }}</strong></article>
            </div>
          </article>

          <article class="portal-card">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.financeDesk }}</h2>
                <div class="portal-panel__meta">{{ copy.financeDeskDesc }}</div>
              </div>
            </div>
            <div class="portal-actions-grid" style="margin-top: 18px; grid-template-columns: 1fr">
              <button v-if="detail.order" type="button" class="portal-action-card" @click="router.push(`/orders/${detail.order.id}`)">
                <strong>{{ copy.openOrder }}</strong>
                <span>{{ detail.order.orderNo }}</span>
              </button>
              <button v-if="detail.services[0]" type="button" class="portal-action-card" @click="router.push(`/services/${detail.services[0].id}`)">
                <strong>{{ copy.openService }}</strong>
                <span>{{ detail.services[0].serviceNo }}</span>
              </button>
              <button type="button" class="portal-action-card" @click="openTicketForService">
                <strong>{{ copy.openTicket }}</strong>
                <span>{{ pickLabel(localeStore.locale, "围绕该账单继续提交支持请求。", "Create a support request around this invoice.") }}</span>
              </button>
              <button
                v-if="detail.invoice.status === 'UNPAID'"
                type="button"
                class="portal-action-card"
                @click="handlePay(detail.invoice)"
              >
                <strong>{{ copy.pay }}</strong>
                <span>{{ pickLabel(localeStore.locale, "优先使用余额抵扣，不足时走在线支付。", "Use balance first, then fall back to online payment.") }}</span>
              </button>
            </div>
          </article>
        </section>

        <section class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.financeTimeline }}</h2>
              <div class="portal-panel__meta">{{ copy.financeTimelineDesc }}</div>
            </div>
          </div>
          <div v-if="financeTimeline.length" class="portal-list" style="margin-top: 18px">
            <div v-for="item in financeTimeline" :key="`${item.time}-${item.title}-${item.subtitle}`" class="portal-list-item">
              <div class="portal-list-item__meta">
                <div class="portal-list-item__title">{{ item.title }}</div>
                <div class="portal-list-item__desc">{{ item.subtitle }}</div>
              </div>
              <strong>{{ item.time }}</strong>
            </div>
          </div>
          <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.empty }}</div>
        </section>

        <section class="portal-grid portal-grid--two">
          <article class="portal-card">
            <div class="portal-card-head"><h2 class="portal-panel__title">{{ copy.payments }}</h2></div>
            <el-table v-if="detail.payments.length" :data="detail.payments" border style="margin-top: 18px">
              <el-table-column prop="paymentNo" :label="pickLabel(localeStore.locale, '支付单号', 'Payment No.')" min-width="140" />
              <el-table-column :label="copy.action" min-width="120">
                <template #default="{ row }">{{ formatPortalPaymentChannel(localeStore.locale, row.channel) }}</template>
              </el-table-column>
              <el-table-column :label="copy.amount" min-width="120">
                <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="paidAt" :label="pickLabel(localeStore.locale, '支付时间', 'Paid At')" min-width="160" />
            </el-table>
            <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.empty }}</div>
          </article>

          <article class="portal-card">
            <div class="portal-card-head"><h2 class="portal-panel__title">{{ copy.refunds }}</h2></div>
            <el-table v-if="detail.refunds.length" :data="detail.refunds" border style="margin-top: 18px">
              <el-table-column prop="refundNo" :label="pickLabel(localeStore.locale, '退款单号', 'Refund No.')" min-width="140" />
              <el-table-column :label="copy.amount" min-width="120">
                <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
              </el-table-column>
              <el-table-column prop="reason" :label="pickLabel(localeStore.locale, '原因', 'Reason')" min-width="160" />
              <el-table-column prop="createdAt" :label="pickLabel(localeStore.locale, '创建时间', 'Created At')" min-width="160" />
            </el-table>
            <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.empty }}</div>
          </article>
        </section>

        <section class="portal-card">
          <div class="portal-card-head"><h2 class="portal-panel__title">{{ copy.services }}</h2></div>
          <el-table v-if="detail.services.length" :data="detail.services" border style="margin-top: 18px">
            <el-table-column prop="serviceNo" :label="pickLabel(localeStore.locale, '服务号', 'Service No.')" min-width="140" />
            <el-table-column prop="productName" :label="copy.productName" min-width="180" />
            <el-table-column prop="nextDueAt" :label="pickLabel(localeStore.locale, '到期时间', 'Next Due')" min-width="140" />
            <el-table-column :label="copy.serviceStatus" min-width="120">
              <template #default="{ row }">
                <el-tag :type="portalTagTypeByStatus(row.status)">
                  {{ formatPortalServiceStatus(localeStore.locale, row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column min-width="120" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="router.push(`/services/${row.id}`)">{{ copy.openService }}</el-button>
              </template>
            </el-table-column>
          </el-table>
          <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.empty }}</div>
        </section>
      </div>
    </el-drawer>
  </div>
</template>
