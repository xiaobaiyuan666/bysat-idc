<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { loadDashboard, type PortalDashboard, type PortalInvoice, type PortalOrder, type PortalService } from "@/api/portal";
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

const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(false);
const error = ref("");
const dashboard = ref<PortalDashboard | null>(null);

const services = computed(() => dashboard.value?.services ?? []);
const orders = computed(() => dashboard.value?.orders ?? []);
const invoices = computed(() => dashboard.value?.invoices ?? []);
const tickets = computed(() => dashboard.value?.tickets ?? []);
const wallet = computed(() => dashboard.value?.wallet);

const activeServices = computed(() => services.value.filter(item => item.status === "ACTIVE"));
const dueSoonServices = computed(() =>
  services.value
    .filter(item => {
      if (!item.nextDueAt) return false;
      const target = new Date(item.nextDueAt);
      if (Number.isNaN(target.getTime())) return false;
      const diff = target.getTime() - Date.now();
      return diff >= 0 && diff <= 7 * 24 * 60 * 60 * 1000;
    })
    .slice(0, 5)
);
const unpaidInvoices = computed(() => invoices.value.filter(item => item.status === "UNPAID").slice(0, 6));
const recentOrders = computed(() => orders.value.slice(0, 6));
const openTicketItems = computed(() => tickets.value.filter(item => item.status !== "CLOSED").slice(0, 6));
const unpaidAmount = computed(() =>
  invoices.value.filter(item => item.status === "UNPAID").reduce((sum, item) => sum + item.totalAmount, 0)
);
const providerCount = computed(() => new Set(services.value.map(item => item.providerType).filter(Boolean)).size);
const providerMix = computed(() => {
  const counters = new Map<string, number>();
  for (const item of services.value) {
    const key = item.providerType || "LOCAL";
    counters.set(key, (counters.get(key) || 0) + 1);
  }
  return Array.from(counters.entries())
    .map(([name, count]) => ({ name, count }))
    .sort((a, b) => b.count - a.count)
    .slice(0, 5);
});
const activityFeed = computed(() => {
  const items: Array<{ time: string; title: string; subtitle: string; action: () => void }> = [];

  for (const service of dueSoonServices.value.slice(0, 3)) {
    items.push({
      time: service.nextDueAt || service.createdAt || "",
      title: pickLabel(localeStore.locale, "服务即将到期", "Service Expiring Soon"),
      subtitle: `${service.productName} / ${service.serviceNo}`,
      action: () => openService(service)
    });
  }

  for (const invoice of unpaidInvoices.value.slice(0, 3)) {
    items.push({
      time: invoice.dueAt || "",
      title: pickLabel(localeStore.locale, "账单待支付", "Invoice Pending"),
      subtitle: `${invoice.invoiceNo} / ${formatPortalMoney(localeStore.locale, invoice.totalAmount)}`,
      action: () => router.push(`/invoices/${invoice.id}`)
    });
  }

  for (const ticket of openTicketItems.value.slice(0, 3)) {
    items.push({
      time: ticket.updatedAt || "",
      title: pickLabel(localeStore.locale, "支持工单待跟进", "Ticket Needs Follow-up"),
      subtitle: `${ticket.no} / ${ticket.title}`,
      action: () => openTicketsCenter({ keyword: ticket.no })
    });
  }

  for (const order of recentOrders.value.slice(0, 2)) {
    items.push({
      time: order.createdAt || "",
      title: pickLabel(localeStore.locale, "新近订单", "Recent Order"),
      subtitle: `${order.orderNo} / ${order.productName}`,
      action: () => router.push(`/orders/${order.id}`)
    });
  }

  return items
    .filter(item => item.time)
    .sort((a, b) => String(b.time).localeCompare(String(a.time)))
    .slice(0, 8);
});
const recommendedActions = computed(() => {
  const actions: Array<{ title: string; description: string; action: () => void }> = [];
  if (dueSoonServices.value[0]) {
    actions.push({
      title: pickLabel(localeStore.locale, "优先续费服务", "Handle Renewals"),
      description: `${dueSoonServices.value[0].productName} / ${dueSoonServices.value[0].nextDueAt || "-"}`,
      action: () => openService(dueSoonServices.value[0])
    });
  }
  if (unpaidInvoices.value[0]) {
    actions.push({
      title: pickLabel(localeStore.locale, "处理待支付账单", "Pay Invoices"),
      description: `${unpaidInvoices.value[0].invoiceNo} / ${formatPortalMoney(localeStore.locale, unpaidInvoices.value[0].totalAmount)}`,
      action: () => router.push(`/invoices/${unpaidInvoices.value[0].id}`)
    });
  }
  if (openTicketItems.value[0]) {
    actions.push({
      title: pickLabel(localeStore.locale, "跟进支持工单", "Follow Tickets"),
      description: openTicketItems.value[0].title,
      action: () => openTicketsCenter({ keyword: openTicketItems.value[0].no })
    });
  }
  actions.push({
    title: pickLabel(localeStore.locale, "查看客户资料", "Review Account"),
    description: pickLabel(localeStore.locale, "维护联系人、实名和安全设置", "Manage contacts, identity, and security"),
    action: () => router.push("/account/profile")
  });
  return actions.slice(0, 4);
});
const quickActions = computed(() => [
  {
    title: copy.value.goServices,
    description: pickLabel(localeStore.locale, "查看所有实例、到期时间和状态。", "Review services, due dates, and status."),
    action: () => router.push("/services")
  },
  {
    title: copy.value.goInvoices,
    description: pickLabel(localeStore.locale, "优先处理未支付账单和财务压力。", "Handle unpaid invoices and finance pressure."),
    action: () => router.push("/invoices")
  },
  {
    title: copy.value.goTickets,
    description: pickLabel(localeStore.locale, "跟进待处理和待回复工单。", "Follow up on open and waiting tickets."),
    action: () => router.push("/tickets")
  },
  {
    title: copy.value.goAccount,
    description: pickLabel(localeStore.locale, "维护联系人、实名和安全设置。", "Manage contacts, identity, and security settings."),
    action: () => router.push("/account/profile")
  }
]);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "客户控制台", "Client Console"),
  title: pickLabel(localeStore.locale, "云业务客户门户", "Cloud Client Portal"),
  subtitle: pickLabel(
    localeStore.locale,
    "把服务、账单、订单、工单和账户状态集中到一个工作台里，优先处理最常用、最紧急的客户动作。",
    "Bring services, invoices, orders, tickets, and account status into one workbench."
  ),
  declaration: pickLabel(
    localeStore.locale,
    "无穷云IDC业务管理系统由江苏白猿网络科技有限公司 - 猿创软件开发 100% AI 开发，全部著作权归白猿科技所有。",
    "Infinity Cloud IDC Management System is 100% AI-developed by Jiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development."
  ),
  serviceCount: pickLabel(localeStore.locale, "运行中服务", "Active Services"),
  invoicePressure: pickLabel(localeStore.locale, "待支付账单", "Unpaid Invoices"),
  ticketOpen: pickLabel(localeStore.locale, "处理中工单", "Open Tickets"),
  availableCredit: pickLabel(localeStore.locale, "可用授信", "Available Credit"),
  dueSoonTitle: pickLabel(localeStore.locale, "即将到期服务", "Services Due Soon"),
  dueSoonDesc: pickLabel(localeStore.locale, "优先处理 7 天内到期服务，减少续费遗漏。", "Prioritize services due within 7 days."),
  orderTitle: pickLabel(localeStore.locale, "最近订单", "Recent Orders"),
  invoiceTitle: pickLabel(localeStore.locale, "待支付账单", "Unpaid Invoices"),
  ticketTitle: pickLabel(localeStore.locale, "待跟进工单", "Open Tickets"),
  walletTitle: pickLabel(localeStore.locale, "资金概览", "Funding Overview"),
  walletDesc: pickLabel(
    localeStore.locale,
    "直接查看余额、授信和当前账务压力，并快速跳到财务中心处理。",
    "View balance, credit, and billing pressure in one place."
  ),
  actionTitle: pickLabel(localeStore.locale, "高频入口", "Quick Actions"),
  actionDesc: pickLabel(localeStore.locale, "把最常用的客户动作收敛到一屏完成。", "Complete the most common customer actions from one screen."),
  queueTitle: pickLabel(localeStore.locale, "业务待办", "Action Queues"),
  queueDesc: pickLabel(localeStore.locale, "把续费、财务和支持三条链的待办收敛到同一块区域。", "Keep renewal, finance, and support queues in one area."),
  nextTitle: pickLabel(localeStore.locale, "推荐动作", "Recommended Actions"),
  nextDesc: pickLabel(localeStore.locale, "根据当前服务、账单和工单状态，优先完成下一步操作。", "Prioritize the next best actions based on service, invoice, and ticket status."),
  activityTitle: pickLabel(localeStore.locale, "业务活动", "Activity Feed"),
  activityDesc: pickLabel(localeStore.locale, "把服务、账单、订单和工单的关键变动放到一条时间线上。", "Keep important service, invoice, order, and ticket changes in one timeline."),
  renewQueue: pickLabel(localeStore.locale, "续费队列", "Renewal Queue"),
  financeQueue: pickLabel(localeStore.locale, "财务队列", "Finance Queue"),
  supportQueue: pickLabel(localeStore.locale, "支持队列", "Support Queue"),
  providerMix: pickLabel(localeStore.locale, "渠道分布", "Provider Mix"),
  providerMixDesc: pickLabel(localeStore.locale, "当前服务主要落在哪些资源渠道。", "See which upstream channels your services are mainly using."),
  riskTitle: pickLabel(localeStore.locale, "风险信号", "Risk Signals"),
  riskDesc: pickLabel(localeStore.locale, "优先跟进到期、未支付和支持工单。", "Prioritize renewals, unpaid invoices, and support follow-up."),
  heroTitle: pickLabel(localeStore.locale, "今日运营总览", "Today's Overview"),
  heroDesc: pickLabel(localeStore.locale, "把服务、财务与支持状态汇总到一块，先处理最重要的客户动作。", "Summarize service, finance, and support status in one place and handle the most important actions first."),
  heroAmount: pickLabel(localeStore.locale, "待支付金额", "Outstanding Amount"),
  heroProvider: pickLabel(localeStore.locale, "对接渠道数", "Providers"),
  walletDesk: pickLabel(localeStore.locale, "资金与授信", "Funding & Credit"),
  walletDeskDesc: pickLabel(localeStore.locale, "从控制台首页快速掌握余额、授信和资金承压情况。", "Understand balance, credit, and funding pressure directly from the console."),
  dueSoonSignal: pickLabel(localeStore.locale, "即将到期服务", "Services Expiring Soon"),
  unpaidSignal: pickLabel(localeStore.locale, "待支付账单", "Unpaid Invoices"),
  ticketSignal: pickLabel(localeStore.locale, "待跟进工单", "Open Tickets"),
  noData: pickLabel(localeStore.locale, "暂无数据", "No data"),
  loadError: pickLabel(localeStore.locale, "控制台数据加载失败", "Failed to load dashboard"),
  open: pickLabel(localeStore.locale, "进入", "Open"),
  goServices: pickLabel(localeStore.locale, "服务中心", "Services"),
  goInvoices: pickLabel(localeStore.locale, "财务中心", "Finance"),
  goOrders: pickLabel(localeStore.locale, "订单中心", "Orders"),
  goTickets: pickLabel(localeStore.locale, "工单中心", "Tickets"),
  goAccount: pickLabel(localeStore.locale, "账户中心", "Account"),
  serviceNo: pickLabel(localeStore.locale, "服务号", "Service No."),
  nextDueAt: pickLabel(localeStore.locale, "到期时间", "Next Due"),
  productName: pickLabel(localeStore.locale, "商品", "Product"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  orderNo: pickLabel(localeStore.locale, "订单号", "Order No."),
  invoiceNo: pickLabel(localeStore.locale, "账单号", "Invoice No."),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  dueAt: pickLabel(localeStore.locale, "到期时间", "Due At"),
  ticketNo: pickLabel(localeStore.locale, "工单号", "Ticket No."),
  titleCol: pickLabel(localeStore.locale, "标题", "Title"),
  updatedAt: pickLabel(localeStore.locale, "更新时间", "Updated At"),
  balance: pickLabel(localeStore.locale, "账户余额", "Balance"),
  creditLimit: pickLabel(localeStore.locale, "授信额度", "Credit Limit"),
  creditUsed: pickLabel(localeStore.locale, "已用授信", "Credit Used"),
  availableCreditLabel: pickLabel(localeStore.locale, "可用授信", "Available Credit")
}));

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

function openService(item: PortalService) {
  void router.push(`/services/${item.id}`);
}

function openFinance(query?: Record<string, string>) {
  void router.push({ path: "/invoices", query });
}

function openOrders(query?: Record<string, string>) {
  void router.push({ path: "/orders", query });
}

function openTicketsCenter(query?: Record<string, string>) {
  void router.push({ path: "/tickets", query });
}

function orderStatus(item: PortalOrder) {
  return formatPortalOrderStatus(localeStore.locale, item.status);
}

function invoiceStatus(item: PortalInvoice) {
  return formatPortalInvoiceStatus(localeStore.locale, item.status);
}

onMounted(() => {
  void fetchDashboard();
});
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

        <div class="portal-toolbar" style="margin-top: 20px">
          <el-button type="primary" plain @click="router.push('/services')">{{ copy.goServices }}</el-button>
          <el-button plain @click="router.push('/invoices')">{{ copy.goInvoices }}</el-button>
        <el-button plain @click="router.push('/orders')">{{ copy.goOrders }}</el-button>
        <el-button plain @click="router.push('/tickets')">{{ copy.goTickets }}</el-button>
          <el-button plain @click="router.push('/account/profile')">{{ copy.goAccount }}</el-button>
        </div>

        <section class="portal-console-hero" style="margin-top: 20px">
          <div class="portal-console-hero__body">
            <div class="portal-kicker">{{ copy.heroTitle }}</div>
            <h2>{{ copy.title }}</h2>
            <p>{{ copy.heroDesc }}</p>
            <div class="portal-console-hero__metrics">
              <div class="portal-mini-card">
                <span class="portal-mini-card__label">{{ copy.heroAmount }}</span>
                <strong class="portal-mini-card__value">{{ formatPortalMoney(localeStore.locale, unpaidAmount) }}</strong>
              </div>
              <div class="portal-mini-card">
                <span class="portal-mini-card__label">{{ copy.heroProvider }}</span>
                <strong class="portal-mini-card__value">{{ providerCount }}</strong>
              </div>
            </div>
          </div>
          <div class="portal-console-hero__aside">
            <div class="portal-mini-card">
              <span class="portal-mini-card__label">{{ copy.serviceCount }}</span>
              <strong class="portal-mini-card__value">{{ activeServices.length }}</strong>
              <span class="portal-mini-card__hint">{{ copy.goServices }}</span>
            </div>
            <div class="portal-mini-card">
              <span class="portal-mini-card__label">{{ copy.invoicePressure }}</span>
              <strong class="portal-mini-card__value">{{ unpaidInvoices.length }}</strong>
              <span class="portal-mini-card__hint">{{ copy.goInvoices }}</span>
            </div>
            <div class="portal-mini-card">
              <span class="portal-mini-card__label">{{ copy.ticketOpen }}</span>
              <strong class="portal-mini-card__value">{{ openTicketItems.length }}</strong>
              <span class="portal-mini-card__hint">{{ copy.goTickets }}</span>
            </div>
          </div>
        </section>

      <section class="portal-card" style="margin-top: 20px; padding: 0">
        <div style="padding: 20px 20px 0">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.walletDesk }}</h2>
              <div class="portal-panel__meta">{{ copy.walletDeskDesc }}</div>
            </div>
            <el-button link type="primary" @click="router.push('/wallet/transactions')">{{ copy.open }}</el-button>
          </div>
        </div>
        <div class="portal-grid portal-grid--four" style="padding: 18px 20px 20px">
          <article class="portal-stat">
            <h3>{{ copy.balance }}</h3>
            <strong>{{ formatPortalMoney(localeStore.locale, wallet?.balance ?? 0) }}</strong>
          </article>
          <article class="portal-stat">
            <h3>{{ copy.creditLimit }}</h3>
            <strong>{{ formatPortalMoney(localeStore.locale, wallet?.creditLimit ?? 0) }}</strong>
          </article>
          <article class="portal-stat">
            <h3>{{ copy.creditUsed }}</h3>
            <strong>{{ formatPortalMoney(localeStore.locale, wallet?.creditUsed ?? 0) }}</strong>
          </article>
          <article class="portal-stat">
            <h3>{{ copy.availableCreditLabel }}</h3>
            <strong>{{ formatPortalMoney(localeStore.locale, wallet?.availableCredit ?? 0) }}</strong>
          </article>
        </div>
      </section>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.actionTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.actionDesc }}</div>
          </div>
        </div>
        <div class="portal-actions-grid" style="margin-top: 18px">
          <button
            v-for="item in quickActions"
            :key="item.title"
            type="button"
            class="portal-action-card"
            @click="item.action()"
          >
            <strong>{{ item.title }}</strong>
            <span>{{ item.description }}</span>
          </button>
        </div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.riskTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.riskDesc }}</div>
          </div>
        </div>
        <div class="portal-hero-grid" style="margin-top: 18px">
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.dueSoonSignal }}</span>
            <strong class="portal-mini-card__value">{{ dueSoonServices.length }}</strong>
            <span class="portal-mini-card__hint">{{ copy.dueSoonDesc }}</span>
          </div>
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.unpaidSignal }}</span>
            <strong class="portal-mini-card__value">{{ unpaidInvoices.length }}</strong>
            <span class="portal-mini-card__hint">{{ copy.invoiceTitle }}</span>
          </div>
          <div class="portal-mini-card">
            <span class="portal-mini-card__label">{{ copy.ticketSignal }}</span>
            <strong class="portal-mini-card__value">{{ openTicketItems.length }}</strong>
            <span class="portal-mini-card__hint">{{ copy.ticketTitle }}</span>
          </div>
        </div>
      </article>
    </section>

    <section class="portal-grid portal-grid--two">
      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.nextTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.nextDesc }}</div>
          </div>
        </div>
        <div class="portal-actions-grid" style="margin-top: 18px">
          <button
            v-for="item in recommendedActions"
            :key="item.title"
            type="button"
            class="portal-action-card"
            @click="item.action()"
          >
            <strong>{{ item.title }}</strong>
            <span>{{ item.description }}</span>
          </button>
        </div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.queueTitle }}</h2>
            <div class="portal-panel__meta">{{ copy.queueDesc }}</div>
          </div>
        </div>
        <div class="portal-lane-grid" style="margin-top: 18px">
          <div class="portal-lane">
            <div class="portal-lane__title">{{ copy.renewQueue }}</div>
            <div v-if="dueSoonServices.length" class="portal-lane__list">
              <button
                v-for="item in dueSoonServices.slice(0, 4)"
                :key="item.id"
                type="button"
                class="portal-lane__item"
                @click="openService(item)"
              >
                <strong>{{ item.productName }}</strong>
                <span>{{ item.nextDueAt || "-" }}</span>
              </button>
            </div>
            <div v-else class="portal-empty-state">{{ copy.noData }}</div>
          </div>
          <div class="portal-lane">
            <div class="portal-lane__title">{{ copy.financeQueue }}</div>
            <div v-if="unpaidInvoices.length" class="portal-lane__list">
              <button
                v-for="item in unpaidInvoices.slice(0, 4)"
                :key="item.id"
                type="button"
                class="portal-lane__item"
                @click="router.push(`/invoices/${item.id}`)"
              >
                <strong>{{ item.invoiceNo }}</strong>
                <span>{{ formatPortalMoney(localeStore.locale, item.totalAmount) }}</span>
              </button>
            </div>
            <div v-else class="portal-empty-state">{{ copy.noData }}</div>
          </div>
          <div class="portal-lane">
            <div class="portal-lane__title">{{ copy.supportQueue }}</div>
            <div v-if="openTicketItems.length" class="portal-lane__list">
              <button
                v-for="item in openTicketItems.slice(0, 4)"
                :key="item.no"
                type="button"
                class="portal-lane__item"
                @click="openTicketsCenter({ keyword: item.no })"
              >
                <strong>{{ item.title }}</strong>
                <span>{{ formatPortalTicketStatus(localeStore.locale, item.status) }}</span>
              </button>
            </div>
            <div v-else class="portal-empty-state">{{ copy.noData }}</div>
          </div>
        </div>
      </article>

      <article class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.providerMix }}</h2>
            <div class="portal-panel__meta">{{ copy.providerMixDesc }}</div>
          </div>
        </div>
        <div v-if="providerMix.length" class="portal-list" style="margin-top: 18px">
          <div v-for="item in providerMix" :key="item.name" class="portal-list-item">
            <div class="portal-list-item__meta">
              <div class="portal-list-item__title">{{ item.name }}</div>
              <div class="portal-list-item__desc">{{ copy.serviceCount }}</div>
            </div>
            <strong>{{ item.count }}</strong>
          </div>
        </div>
        <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noData }}</div>
      </article>
    </section>

    <section class="portal-card">
      <div class="portal-card-head">
        <div>
          <h2 class="portal-panel__title">{{ copy.activityTitle }}</h2>
          <div class="portal-panel__meta">{{ copy.activityDesc }}</div>
        </div>
      </div>
      <div v-if="activityFeed.length" class="portal-list" style="margin-top: 18px">
        <button
          v-for="item in activityFeed"
          :key="`${item.time}-${item.title}-${item.subtitle}`"
          type="button"
          class="portal-lane__item"
          @click="item.action()"
        >
          <div class="portal-list-item__meta">
            <div class="portal-list-item__title">{{ item.title }}</div>
            <div class="portal-list-item__desc">{{ item.subtitle }}</div>
          </div>
          <span>{{ item.time }}</span>
        </button>
      </div>
      <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noData }}</div>
    </section>

    <section class="portal-grid portal-grid--three">
      <article class="portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.orderTitle }}</h2>
          </div>
          <el-button link type="primary" @click="router.push('/orders')">{{ copy.open }}</el-button>
        </div>
        <el-table v-if="recentOrders.length" :data="recentOrders" border>
          <el-table-column prop="orderNo" :label="copy.orderNo" min-width="140" />
          <el-table-column prop="productName" :label="copy.productName" min-width="160" />
          <el-table-column :label="copy.amount" min-width="110">
            <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
          </el-table-column>
          <el-table-column :label="copy.status" min-width="110">
            <template #default="{ row }">
              <el-tag :type="portalTagTypeByStatus(row.status)">{{ orderStatus(row) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column width="80" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="router.push(`/orders/${row.id}`)">{{ copy.open }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state" style="margin: 20px">{{ copy.noData }}</div>
      </article>

      <article class="portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.invoiceTitle }}</h2>
          </div>
          <el-button link type="primary" @click="router.push('/invoices')">{{ copy.open }}</el-button>
        </div>
        <el-table v-if="unpaidInvoices.length" :data="unpaidInvoices" border>
          <el-table-column prop="invoiceNo" :label="copy.invoiceNo" min-width="140" />
          <el-table-column prop="productName" :label="copy.productName" min-width="160" />
          <el-table-column :label="copy.amount" min-width="110">
            <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.totalAmount) }}</template>
          </el-table-column>
          <el-table-column :label="copy.dueAt" min-width="130">
            <template #default="{ row }">{{ row.dueAt || "-" }}</template>
          </el-table-column>
          <el-table-column :label="copy.status" min-width="110">
            <template #default="{ row }">
              <el-tag :type="portalTagTypeByStatus(row.status)">{{ invoiceStatus(row) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column width="80" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="router.push(`/invoices/${row.id}`)">{{ copy.open }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state" style="margin: 20px">{{ copy.noData }}</div>
      </article>

      <article class="portal-table-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.ticketTitle }}</h2>
          </div>
          <el-button link type="primary" @click="router.push('/tickets')">{{ copy.open }}</el-button>
        </div>
        <el-table v-if="openTicketItems.length" :data="openTicketItems" border>
          <el-table-column prop="no" :label="copy.ticketNo" min-width="140" />
          <el-table-column prop="title" :label="copy.titleCol" min-width="170" show-overflow-tooltip />
          <el-table-column :label="copy.status" min-width="120">
            <template #default="{ row }">
              <el-tag :type="portalTagTypeByStatus(row.status)">
                {{ formatPortalTicketStatus(localeStore.locale, row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="updatedAt" :label="copy.updatedAt" min-width="150" />
          <el-table-column width="80" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openTicketsCenter({ keyword: row.no })">{{ copy.open }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state" style="margin: 20px">{{ copy.noData }}</div>
      </article>
    </section>
  </div>
</template>
