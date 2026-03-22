<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { ElMessage } from "element-plus";
import { loadInvoices, payInvoice, type PortalInvoice } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalInvoiceStatus,
  formatPortalMoney,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const payingId = ref<number | null>(null);
const invoices = ref<PortalInvoice[]>([]);
const keyword = ref("");
const status = ref("");

const filteredInvoices = computed(() =>
  invoices.value.filter(item => {
    const matchedKeyword =
      !keyword.value || item.invoiceNo.includes(keyword.value) || item.productName.includes(keyword.value);
    const matchedStatus = !status.value || item.status === status.value;
    return matchedKeyword && matchedStatus;
  })
);

const unpaidCount = computed(() => invoices.value.filter(item => item.status === "UNPAID").length);
const overdueCount = computed(
  () =>
    invoices.value.filter(item => item.status === "UNPAID" && item.dueAt < new Date().toISOString()).length
);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "账单中心", "Invoices"),
  title: pickLabel(localeStore.locale, "账单中心", "Invoices"),
  subtitle: pickLabel(
    localeStore.locale,
    "集中查看未支付、已支付和退款账单，并支持直接支付待处理账单。",
    "Review unpaid, paid, and refunded invoices and pay open items directly."
  ),
  total: pickLabel(localeStore.locale, "账单总数", "Total Invoices"),
  unpaid: pickLabel(localeStore.locale, "未支付", "Unpaid"),
  overdue: pickLabel(localeStore.locale, "逾期风险", "Overdue"),
  listTitle: pickLabel(localeStore.locale, "账单列表", "Invoice List"),
  listDesc: pickLabel(
    localeStore.locale,
    "支持按编号、产品名称和状态筛选。",
    "Filter by invoice number, product name, and status."
  ),
  keywordPlaceholder: pickLabel(localeStore.locale, "搜索账单编号或产品名称", "Search invoice no. or product"),
  statusPlaceholder: pickLabel(localeStore.locale, "账单状态", "Invoice Status"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  invoiceNo: pickLabel(localeStore.locale, "账单编号", "Invoice No."),
  productName: pickLabel(localeStore.locale, "产品名称", "Product"),
  cycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  dueAt: pickLabel(localeStore.locale, "到期时间", "Due At"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  payNow: pickLabel(localeStore.locale, "立即支付", "Pay Now"),
  noAction: pickLabel(localeStore.locale, "无需操作", "No Action"),
  empty: pickLabel(localeStore.locale, "暂无匹配的账单记录。", "No matching invoices."),
  loadError: pickLabel(localeStore.locale, "账单列表加载失败", "Failed to load invoices"),
  paySuccess: pickLabel(localeStore.locale, "支付成功", "Payment successful"),
  payError: pickLabel(localeStore.locale, "支付失败", "Payment failed")
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

async function handlePay(row: PortalInvoice) {
  payingId.value = row.id;
  try {
    const result = await payInvoice(row.id);
    ElMessage.success(
      result.service
        ? pickLabel(
            localeStore.locale,
            `账单 ${result.invoice.invoiceNo} 支付成功，服务 ${result.service.serviceNo} 已激活`,
            `Invoice ${result.invoice.invoiceNo} paid successfully. Service ${result.service.serviceNo} is now active.`
          )
        : pickLabel(
            localeStore.locale,
            `账单 ${result.invoice.invoiceNo} 支付成功`,
            `Invoice ${result.invoice.invoiceNo} paid successfully.`
          )
    );
    await fetchInvoices();
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : copy.value.payError);
  } finally {
    payingId.value = null;
  }
}

onMounted(fetchInvoices);
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

      <div class="portal-grid portal-grid--three" style="margin-top: 20px">
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
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'UNPAID')" value="UNPAID" />
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'PAID')" value="PAID" />
          <el-option :label="formatPortalInvoiceStatus(localeStore.locale, 'REFUNDED')" value="REFUNDED" />
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
        <el-table-column :label="copy.cycle" min-width="120">
          <template #default="{ row }">{{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}</template>
        </el-table-column>
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.totalAmount) }}</template>
        </el-table-column>
        <el-table-column prop="dueAt" :label="copy.dueAt" min-width="160" />
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalInvoiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.action" min-width="140" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'UNPAID'"
              type="primary"
              link
              :loading="payingId === row.id"
              @click="handlePay(row)"
            >
              {{ copy.payNow }}
            </el-button>
            <span v-else style="color: #94a3b8">{{ copy.noAction }}</span>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state" style="margin: 20px">
        {{ copy.empty }}
      </div>
    </section>
  </div>
</template>
