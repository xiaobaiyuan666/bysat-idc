<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { loadDashboard, type PortalDashboard } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalInvoiceStatus,
  formatPortalMoney,
  formatPortalOrderStatus,
  formatPortalTicketStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const router = useRouter();
const loading = ref(true);
const error = ref("");
const dashboard = ref<PortalDashboard | null>(null);

const stats = computed(() => dashboard.value?.stats ?? []);
const services = computed(() => dashboard.value?.services ?? []);
const orders = computed(() => dashboard.value?.orders ?? []);
const invoices = computed(() => dashboard.value?.invoices ?? []);
const tickets = computed(() => dashboard.value?.tickets ?? []);
const wallet = computed(() => dashboard.value?.wallet);

const activeServices = computed(() => services.value.filter(item => item.status === "ACTIVE").length);
const unpaidTotal = computed(() => invoices.value.reduce((total, item) => total + Number(item.totalAmount || 0), 0));
const dueSoonServices = computed(() =>
  services.value
    .filter(item => {
      const nextDue = item.nextDueAt ? new Date(item.nextDueAt) : null;
      if (!nextDue || Number.isNaN(nextDue.getTime())) return false;
      const diff = nextDue.getTime() - Date.now();
      return diff >= 0 && diff <= 7 * 24 * 60 * 60 * 1000;
    })
    .slice(0, 5)
);

const quickActions = computed(() => [
  { label: pickLabel(localeStore.locale, "服务中心", "Services"), path: "/services", type: "primary" as const },
  { label: pickLabel(localeStore.locale, "待支付账单", "Unpaid Invoices"), path: "/invoices", type: "danger" as const },
  { label: pickLabel(localeStore.locale, "提交工单", "Open Ticket"), path: "/tickets", type: "warning" as const },
  { label: pickLabel(localeStore.locale, "钱包与流水", "Wallet"), path: "/wallet", type: "success" as const },
  { label: pickLabel(localeStore.locale, "进入商城", "Store"), path: "/store", type: "info" as const }
]);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "白猿科技客户中心", "BYSAT Client Portal"),
  title: pickLabel(localeStore.locale, "客户控制台", "Client Console"),
  subtitle: pickLabel(
    localeStore.locale,
    "把服务、账单、订单、工单和钱包集中在同一张工作台里，方便客户快速续费、支付、报障和查看业务状态。",
    "A single workbench for services, invoices, orders, tickets, and wallet actions."
  ),
  declaration: pickLabel(
    localeStore.locale,
    "本系统由江苏白猿网络科技有限公司 · 猿创软件业务组 100% AI 开发，白猿科技享有 100% 著作权。官网：www.bysat.com",
    "This system is 100% AI-developed by Jiangsu Baiyuan Network Technology Co., Ltd., with full copyright retained by the company."
  ),
  quickActionTitle: pickLabel(localeStore.locale, "快捷入口", "Quick Actions"),
  quickActionDesc: pickLabel(localeStore.locale, "常用动作直接进入，不再在菜单里来回找。", "Jump straight to the most-used self-service actions."),
  activeServiceTitle: pickLabel(localeStore.locale, "运行中服务", "Active Services"),
  unpaidTotalTitle: pickLabel(localeStore.locale, "待支付金额", "Unpaid Amount"),
  ticketTitle: pickLabel(localeStore.locale, "待跟进工单", "Open Tickets"),
  availableCreditTitle: pickLabel(localeStore.locale, "可用授信", "Available Credit"),
  dueSoonTitle: pickLabel(localeStore.locale, "近期到期服务", "Services Due Soon"),
  dueSoonDesc: pickLabel(localeStore.locale, "优先处理 7 天内到期的业务，避免影响服务连续性。", "Prioritize services due within 7 days."),
  walletPanelTitle: pickLabel(localeStore.locale, "资金概览", "Funding Overview"),
  walletPanelDesc: pickLabel(localeStore.locale, "客户侧可直接看到余额、授信和当前待支付压力。", "Balance, credit, and invoice pressure in one place."),
  goWallet: pickLabel(localeStore.locale, "前往钱包", "Open Wallet"),
  goInvoices: pickLabel(localeStore.locale, "处理账单", "Handle Invoices"),
  goAccount: pickLabel(localeStore.locale, "账户资料", "Account Profile"),
  loadError: pickLabel(localeStore.locale, "控制台数据加载失败", "Failed to load dashboard"),
  recentOrders: pickLabel(localeStore.locale, "最近订单", "Recent Orders"),
  recentOrdersDesc: pickLabel(localeStore.locale, "显示最近下单记录，方便跟进付款和开通。", "Track the latest orders and activation progress."),
  openTickets: pickLabel(localeStore.locale, "待处理工单", "Open Tickets"),
  openTicketsDesc: pickLabel(localeStore.locale, "集中查看当前仍在处理中的支持请求。", "Review support requests that still need attention."),
  unpaidInvoices: pickLabel(localeStore.locale, "待支付账单", "Unpaid Invoices"),
  unpaidInvoicesDesc: pickLabel(localeStore.locale, "支付成功后可自动激活或续费对应服务。", "Pay invoices to activate or renew linked services."),
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
  dueAt: pickLabel(localeStore.locale, "到期时间", "Due At"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  detail: pickLabel(localeStore.locale, "进入", "Open")
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

      <el-alert :title="copy.declaration" type="info" :closable="false" show-icon style="margin-top: 20px" />

      <div class="portal-card-head" style="margin-top: 20px">
        <div>
          <h2 class="portal-panel__title">{{ copy.quickActionTitle }}</h2>
          <div class="portal-panel__meta">{{ copy.quickActionDesc }}</div>
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

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <article v-for="item in stats" :key="item.label" class="portal-stat">
          <h3>{{ item.label }}</h3>
          <strong>{{ item.value }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.activeServiceTitle }}</h3>
          <strong>{{ activeServices }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.unpaidTotalTitle }}</h3>
          <strong>{{ money(unpaidTotal) }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.ticketTitle }}</h3>
          <strong>{{ tickets.length }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.availableCreditTitle }}</h3>
          <strong>{{ money(wallet?.availableCredit) }}</strong>
        </article>
      </div>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.dueSoonTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.dueSoonDesc }}</div>
          </div>
          <el-tag type="warning" effect="plain">{{ dueSoonServices.length }}</el-tag>
        </div>

        <el-table v-if="dueSoonServices.length" :data="dueSoonServices" border>
          <el-table-column prop="serviceNo" :label="copy.serviceNo" min-width="170" />
          <el-table-column prop="productName" :label="copy.productName" min-width="220" />
          <el-table-column prop="nextDueAt" :label="copy.nextDueAt" min-width="160" />
          <el-table-column :label="copy.action" min-width="120" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="router.push(`/services/${row.id}`)">
                {{ copy.detail }}
              </el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state">
          {{ copy.emptyServices }}
        </div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.walletPanelTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.walletPanelDesc }}</div>
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
          <div class="portal-summary-row">
            <span>{{ copy.availableCreditTitle }}</span>
            <strong>{{ money(wallet?.availableCredit) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.unpaidTotalTitle }}</span>
            <strong>{{ money(unpaidTotal) }}</strong>
          </div>
        </div>

        <div class="portal-toolbar" style="margin-top: 18px">
          <el-button type="primary" plain @click="router.push('/wallet')">{{ copy.goWallet }}</el-button>
          <el-button type="danger" plain @click="router.push('/invoices')">{{ copy.goInvoices }}</el-button>
          <el-button plain @click="router.push('/account')">{{ copy.goAccount }}</el-button>
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
          <el-table-column :label="copy.action" min-width="120" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="router.push(`/orders?orderNo=${row.orderNo}`)">
                {{ copy.detail }}
              </el-button>
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
          <el-table-column :label="copy.action" min-width="120" fixed="right">
            <template #default>
              <el-button type="primary" link @click="router.push('/tickets')">
                {{ copy.detail }}
              </el-button>
            </template>
          </el-table-column>
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
        <el-table-column :label="copy.action" min-width="120" fixed="right">
          <template #default>
            <el-button type="primary" link @click="router.push('/invoices')">
              {{ copy.detail }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state">
        {{ copy.emptyInvoices }}
      </div>
    </section>
  </div>
</template>
