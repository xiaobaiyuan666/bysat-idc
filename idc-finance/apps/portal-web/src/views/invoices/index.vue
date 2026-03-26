<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import { loadInvoices, loadServices, payInvoice, type PortalInvoice, type PortalService } from "@/api/portal";
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

const loading = ref(true);
const error = ref("");
const payingId = ref<number | null>(null);
const invoices = ref<PortalInvoice[]>([]);
const services = ref<PortalService[]>([]);
const keyword = ref("");
const status = ref("");
const serviceFilterId = ref<number | null>(null);

const serviceMap = computed(() => new Map(services.value.map(item => [item.id, item])));
const serviceByInvoiceId = computed(() => {
  const map = new Map<number, PortalService>();
  for (const service of services.value) {
    if (service.invoiceId) {
      map.set(service.invoiceId, service);
    }
  }
  return map;
});

const serviceOptions = computed(() =>
  services.value.map(item => ({
    label: `${item.serviceNo} / ${item.productName}`,
    value: item.id
  }))
);

const filteredInvoices = computed(() =>
  invoices.value.filter(item => {
    const matchedKeyword =
      !keyword.value ||
      item.invoiceNo.includes(keyword.value) ||
      item.productName.includes(keyword.value) ||
      item.orderNo.includes(keyword.value);
    const matchedStatus = !status.value || item.status === status.value;
    const matchedService = !serviceFilterId.value || findLinkedService(item)?.id === serviceFilterId.value;
    return matchedKeyword && matchedStatus && matchedService;
  })
);

const unpaidCount = computed(() => invoices.value.filter(item => item.status === "UNPAID").length);
const overdueCount = computed(
  () =>
    invoices.value.filter(item => item.status === "UNPAID" && item.dueAt < new Date().toISOString()).length
);
const linkedServiceCount = computed(() => invoices.value.filter(item => serviceByInvoiceId.value.has(item.id)).length);

const quickActions = computed(() => [
  { label: pickLabel(localeStore.locale, "钱包余额", "Wallet"), path: "/wallet", type: "primary" as const },
  { label: pickLabel(localeStore.locale, "订单中心", "Orders"), path: "/orders", type: "info" as const },
  { label: pickLabel(localeStore.locale, "服务中心", "Services"), path: "/services", type: "success" as const },
  { label: pickLabel(localeStore.locale, "工单中心", "Tickets"), path: "/tickets", type: "warning" as const }
]);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "账单中心", "Invoices"),
  title: pickLabel(localeStore.locale, "账单中心", "Invoices"),
  subtitle: pickLabel(
    localeStore.locale,
    "集中查看未支付、已支付和已退款账单，并直接联动到服务和工单处理页面。",
    "Review invoices and jump straight into related services and tickets."
  ),
  total: pickLabel(localeStore.locale, "账单总数", "Total Invoices"),
  unpaid: pickLabel(localeStore.locale, "未支付", "Unpaid"),
  overdue: pickLabel(localeStore.locale, "逾期风险", "Overdue"),
  linkedServices: pickLabel(localeStore.locale, "已关联服务", "Linked Services"),
  listTitle: pickLabel(localeStore.locale, "账单列表", "Invoice List"),
  listDesc: pickLabel(
    localeStore.locale,
    "支持按账单号、订单号、商品名、账单状态和关联服务过滤。",
    "Filter by invoice, order, product, status, and linked service."
  ),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索账单号、订单号或商品名", "Search invoice, order, or product"),
  statusPlaceholder: pickLabel(localeStore.locale, "账单状态", "Invoice Status"),
  servicePlaceholder: pickLabel(localeStore.locale, "关联服务", "Linked Service"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  invoiceNo: pickLabel(localeStore.locale, "账单编号", "Invoice No."),
  productName: pickLabel(localeStore.locale, "商品名称", "Product"),
  linkedService: pickLabel(localeStore.locale, "关联服务", "Linked Service"),
  cycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  dueAt: pickLabel(localeStore.locale, "到期时间", "Due At"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  payNow: pickLabel(localeStore.locale, "立即支付", "Pay Now"),
  service: pickLabel(localeStore.locale, "查看服务", "Service"),
  ticket: pickLabel(localeStore.locale, "提交工单", "Ticket"),
  order: pickLabel(localeStore.locale, "查看订单", "Order"),
  noAction: pickLabel(localeStore.locale, "无可用操作", "No Action"),
  empty: pickLabel(localeStore.locale, "暂无匹配的账单记录。", "No matching invoices."),
  noService: pickLabel(localeStore.locale, "未绑定服务", "No linked service"),
  loadError: pickLabel(localeStore.locale, "账单列表加载失败", "Failed to load invoices"),
  payError: pickLabel(localeStore.locale, "支付失败", "Payment failed"),
  payRule: pickLabel(
    localeStore.locale,
    "支付待处理账单时，系统会优先尝试余额抵扣；余额不足时再按在线支付处理。",
    "The system tries wallet balance first and falls back to online payment when needed."
  )
}));

function findLinkedService(invoice: PortalInvoice) {
  return serviceByInvoiceId.value.get(invoice.id) ?? null;
}

function goToOrder(invoice: PortalInvoice) {
  void router.push({ path: "/orders", query: { orderNo: invoice.orderNo } });
}

function goToService(invoice: PortalInvoice) {
  const linkedService = findLinkedService(invoice);
  if (!linkedService) return;
  void router.push(`/services/${linkedService.id}`);
}

function goToTicket(invoice: PortalInvoice) {
  const linkedService = findLinkedService(invoice);
  void router.push(
    linkedService
      ? {
          path: "/tickets",
          query: {
            action: "create",
            serviceId: String(linkedService.id)
          }
        }
      : {
          path: "/tickets",
          query: {
            action: "create"
          }
        }
  );
}

async function fetchData() {
  loading.value = true;
  error.value = "";
  try {
    const [invoiceItems, serviceItems] = await Promise.all([loadInvoices(), loadServices()]);
    invoices.value = invoiceItems;
    services.value = serviceItems;
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

async function handlePay(row: PortalInvoice) {
  payingId.value = row.id;
  try {
    const result = await payInvoice(row.id);
    const channelLabel = result.payment
      ? formatPortalPaymentChannel(localeStore.locale, result.payment.channel)
      : pickLabel(localeStore.locale, "系统默认渠道", "Default channel");
    ElMessage.success(
      result.service
        ? pickLabel(
            localeStore.locale,
            `账单 ${result.invoice.invoiceNo} 支付成功，支付单 ${result.payment?.paymentNo || "-"}，渠道 ${channelLabel}，服务 ${result.service.serviceNo} 已激活。`,
            `Invoice ${result.invoice.invoiceNo} paid successfully. Payment ${result.payment?.paymentNo || "-"} via ${channelLabel}. Service ${result.service.serviceNo} is now active.`
          )
        : pickLabel(
            localeStore.locale,
            `账单 ${result.invoice.invoiceNo} 支付成功，支付单 ${result.payment?.paymentNo || "-"}，渠道 ${channelLabel}。`,
            `Invoice ${result.invoice.invoiceNo} paid successfully. Payment ${result.payment?.paymentNo || "-"} via ${channelLabel}.`
          )
    );
    await fetchData();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : copy.value.payError);
  } finally {
    payingId.value = null;
  }
}

function applyRouteFilters() {
  const invoiceNo = typeof route.query.invoiceNo === "string" ? route.query.invoiceNo : "";
  const keywordQuery = typeof route.query.keyword === "string" ? route.query.keyword : "";
  const statusQuery = typeof route.query.status === "string" ? route.query.status : "";
  const serviceIdQuery = typeof route.query.serviceId === "string" ? Number(route.query.serviceId) : NaN;

  keyword.value = invoiceNo || keywordQuery;
  status.value = statusQuery;
  serviceFilterId.value = Number.isFinite(serviceIdQuery) && serviceIdQuery > 0 ? serviceIdQuery : null;
}

onMounted(() => {
  applyRouteFilters();
  void fetchData();
});

watch(
  () => [route.query.invoiceNo, route.query.keyword, route.query.status, route.query.serviceId],
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
        <el-button v-for="item in quickActions" :key="item.path" :type="item.type" plain @click="router.push(item.path)">
          {{ item.label }}
        </el-button>
      </div>

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <div class="portal-stat">
          <h3>{{ copy.total }}</h3>
          <strong>{{ invoices.length }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.unpaid }}</h3>
          <strong>{{ unpaidCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.overdue }}</h3>
          <strong>{{ overdueCount }}</strong>
        </div>
        <div class="portal-stat">
          <h3>{{ copy.linkedServices }}</h3>
          <strong>{{ linkedServiceCount }}</strong>
        </div>
      </div>

      <el-alert :title="copy.payRule" type="info" :closable="false" show-icon style="margin-top: 20px" />
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
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'UNPAID')" value="UNPAID" />
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'PAID')" value="PAID" />
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'REFUNDED')" value="REFUNDED" />
        </el-select>
        <el-select v-model="serviceFilterId" :placeholder="copy.servicePlaceholder" clearable>
          <el-option v-for="item in serviceOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <div class="portal-summary" style="margin: 0; padding: 12px 16px">
          <div class="portal-summary-row" style="font-size: 12px">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredInvoices.length }}</strong>
          </div>
        </div>
      </div>

      <el-table v-if="filteredInvoices.length" :data="filteredInvoices" border>
        <el-table-column prop="invoiceNo" :label="copy.invoiceNo" min-width="170" />
        <el-table-column prop="productName" :label="copy.productName" min-width="180" />
        <el-table-column :label="copy.linkedService" min-width="210">
          <template #default="{ row }">
            <template v-if="findLinkedService(row)">
              <div>{{ findLinkedService(row)?.serviceNo }}</div>
              <small>{{ formatPortalServiceStatus(localeStore.locale, findLinkedService(row)?.status || "") }}</small>
            </template>
            <span v-else style="color: #94a3b8">{{ copy.noService }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="copy.cycle" min-width="120">
          <template #default="{ row }">{{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.totalAmount) }}</template>
        </el-table-column>
        <el-table-column prop="dueAt" :label="copy.dueAt" min-width="160" />
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)" effect="light">
              {{ formatPortalInvoiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.action" min-width="260" fixed="right">
          <template #default="{ row }">
            <div class="portal-actions">
              <el-button
                v-if="row.status === 'UNPAID'"
                type="primary"
                link
                :loading="payingId === row.id"
                @click="handlePay(row)"
              >
                {{ copy.payNow }}
              </el-button>
              <el-button v-if="findLinkedService(row)" type="primary" link @click="goToService(row)">
                {{ copy.service }}
              </el-button>
              <el-button type="primary" link @click="goToTicket(row)">
                {{ copy.ticket }}
              </el-button>
              <el-button type="primary" link @click="goToOrder(row)">
                {{ copy.order }}
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>
  </div>
</template>

<style scoped>
.portal-actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}
</style>
