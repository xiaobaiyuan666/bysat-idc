<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { loadDashboard, type PortalDashboard } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalInvoiceStatus,
  formatPortalMoney,
  formatPortalOrderStatus,
  formatPortalServiceStatus,
  formatPortalTicketStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const loading = ref(true);
const error = ref("");
const dashboard = ref<PortalDashboard | null>(null);

const stats = computed(() => dashboard.value?.stats ?? []);
const services = computed(() => dashboard.value?.services ?? []);
const orders = computed(() => dashboard.value?.orders ?? []);
const invoices = computed(() => dashboard.value?.invoices ?? []);
const tickets = computed(() => dashboard.value?.tickets ?? []);
const wallet = computed(() => dashboard.value?.wallet);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "客户门户", "Client Portal"),
  title: pickLabel(localeStore.locale, "控制台", "Console"),
  subtitle: pickLabel(
    localeStore.locale,
    "集中查看当前客户的服务、订单、账单、工单和资金状态，方便快速处理日常业务。",
    "A single dashboard for services, orders, invoices, tickets, and account balance."
  ),
  loadError: pickLabel(localeStore.locale, "控制台数据加载失败", "Failed to load dashboard"),
  serviceOverview: pickLabel(localeStore.locale, "服务概览", "Service Overview"),
  serviceOverviewDesc: pickLabel(
    localeStore.locale,
    "按当前已开通服务展示，便于快速判断运行状态和到期时间。",
    "Review service status and due dates at a glance."
  ),
  walletInfo: pickLabel(localeStore.locale, "钱包信息", "Wallet"),
  walletInfoDesc: pickLabel(
    localeStore.locale,
    "余额与信用额度可直接用于账单支付。",
    "Balance and credit limit can be used to pay invoices."
  ),
  recentOrders: pickLabel(localeStore.locale, "最近订单", "Recent Orders"),
  recentOrdersDesc: pickLabel(
    localeStore.locale,
    "显示最近下单记录，方便跟进付款和开通。",
    "Track the latest orders and activation progress."
  ),
  openTickets: pickLabel(localeStore.locale, "待处理工单", "Open Tickets"),
  openTicketsDesc: pickLabel(
    localeStore.locale,
    "集中查看当前仍在处理中的支持请求。",
    "Review support requests that still need attention."
  ),
  unpaidInvoices: pickLabel(localeStore.locale, "待支付账单", "Unpaid Invoices"),
  unpaidInvoicesDesc: pickLabel(
    localeStore.locale,
    "支付成功后可自动激活或续费对应服务。",
    "Pay invoices to activate or renew linked services."
  ),
  emptyServices: pickLabel(localeStore.locale, "暂无服务记录，新订单支付后会在这里显示。", "No services yet."),
  emptyOrders: pickLabel(localeStore.locale, "暂无订单记录。", "No orders yet."),
  emptyTickets: pickLabel(localeStore.locale, "暂无工单记录。", "No tickets yet."),
  emptyInvoices: pickLabel(localeStore.locale, "暂无待支付账单。", "No unpaid invoices."),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  creditLimit: pickLabel(localeStore.locale, "信用额度", "Credit Limit"),
  serviceNo: pickLabel(localeStore.locale, "服务编号", "Service No."),
  productName: pickLabel(localeStore.locale, "产品名称", "Product"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  nextDueAt: pickLabel(localeStore.locale, "下次到期", "Next Due"),
  orderNo: pickLabel(localeStore.locale, "订单编号", "Order No."),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  ticketNo: pickLabel(localeStore.locale, "工单编号", "Ticket No."),
  titleCol: pickLabel(localeStore.locale, "标题", "Title"),
  updatedAt: pickLabel(localeStore.locale, "更新时间", "Updated At"),
  invoiceNo: pickLabel(localeStore.locale, "账单编号", "Invoice No."),
  dueAt: pickLabel(localeStore.locale, "到期时间", "Due At")
}));

function money(value: number | string | undefined) {
  return formatPortalMoney(localeStore.locale, value);
}

async function fetchDashboard() {
  loading.value = true;
  error.value = "";
  try {
    dashboard.value = await loadDashboard();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

onMounted(fetchDashboard);
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

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <article v-for="item in stats" :key="item.label" class="portal-stat">
          <h3>{{ item.label }}</h3>
          <strong>{{ item.value }}</strong>
        </article>
      </div>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.serviceOverview }}</h2>
            <div class="portal-panel__meta">{{ copy.serviceOverviewDesc }}</div>
          </div>
          <el-tag type="primary" effect="plain">{{ services.length }}</el-tag>
        </div>

        <el-table v-if="services.length" :data="services" border>
          <el-table-column prop="serviceNo" :label="copy.serviceNo" min-width="170" />
          <el-table-column prop="productName" :label="copy.productName" min-width="220" />
          <el-table-column :label="copy.status" min-width="120">
            <template #default="{ row }">
              <el-tag :type="portalTagTypeByStatus(row.status)">
                {{ formatPortalServiceStatus(localeStore.locale, row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="nextDueAt" :label="copy.nextDueAt" min-width="160" />
        </el-table>
        <div v-else class="portal-empty-state">
          {{ copy.emptyServices }}
        </div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.walletInfo }}</h2>
            <div class="portal-panel__meta">{{ copy.walletInfoDesc }}</div>
          </div>
        </div>

        <div class="portal-summary" style="margin-top: 18px">
          <div class="portal-summary-row">
            <span>{{ copy.balance }}</span>
            <strong>{{ money(wallet?.balance) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.creditLimit }}</span>
            <strong>{{ money(wallet?.creditLimit) }}</strong>
          </div>
        </div>
      </article>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.recentOrders }}</h2>
            <div class="portal-panel__meta">{{ copy.recentOrdersDesc }}</div>
          </div>
          <el-tag type="info" effect="plain">{{ orders.length }}</el-tag>
        </div>

        <el-table v-if="orders.length" :data="orders" border>
          <el-table-column prop="orderNo" :label="copy.orderNo" min-width="170" />
          <el-table-column prop="productName" :label="copy.productName" min-width="200" />
          <el-table-column :label="copy.amount" min-width="120">
            <template #default="{ row }">{{ money(row.amount) }}</template>
          </el-table-column>
          <el-table-column :label="copy.status" min-width="120">
            <template #default="{ row }">
              <el-tag :type="portalTagTypeByStatus(row.status)">
                {{ formatPortalOrderStatus(localeStore.locale, row.status) }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state">
          {{ copy.emptyOrders }}
        </div>
      </article>

      <article class="portal-card portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.openTickets }}</h2>
            <div class="portal-panel__meta">{{ copy.openTicketsDesc }}</div>
          </div>
          <el-tag type="warning" effect="plain">{{ tickets.length }}</el-tag>
        </div>

        <el-table v-if="tickets.length" :data="tickets" border>
          <el-table-column prop="no" :label="copy.ticketNo" min-width="170" />
          <el-table-column prop="title" :label="copy.titleCol" min-width="220" />
          <el-table-column :label="copy.status" min-width="120">
            <template #default="{ row }">
              <el-tag :type="portalTagTypeByStatus(row.status)">
                {{ formatPortalTicketStatus(localeStore.locale, row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="updatedAt" :label="copy.updatedAt" min-width="180" />
        </el-table>
        <div v-else class="portal-empty-state">
          {{ copy.emptyTickets }}
        </div>
      </article>
    </section>

    <section class="portal-card portal-table-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.unpaidInvoices }}</h2>
          <div class="portal-panel__meta">{{ copy.unpaidInvoicesDesc }}</div>
        </div>
        <el-tag type="danger" effect="plain">{{ invoices.length }}</el-tag>
      </div>

      <el-table v-if="invoices.length" :data="invoices" border>
        <el-table-column prop="invoiceNo" :label="copy.invoiceNo" min-width="170" />
        <el-table-column prop="productName" :label="copy.productName" min-width="180" />
        <el-table-column :label="copy.amount" min-width="120">
          <template #default="{ row }">{{ money(row.totalAmount) }}</template>
        </el-table-column>
        <el-table-column prop="dueAt" :label="copy.dueAt" min-width="160" />
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalInvoiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state">
        {{ copy.emptyInvoices }}
      </div>
    </section>
  </div>
</template>
