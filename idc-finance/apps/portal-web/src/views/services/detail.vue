<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { ElMessage } from "element-plus";
import { useRoute, useRouter } from "vue-router";
import { loadServiceDetail, type PortalServiceChangeOrder, type PortalServiceDetailResponse } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalConfiguration,
  formatPortalInvoiceStatus,
  formatPortalMoney,
  formatPortalServiceStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const localeStore = useLocaleStore();
const route = useRoute();
const router = useRouter();

const loading = ref(false);
const error = ref("");
const detail = ref<PortalServiceDetailResponse | null>(null);

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "服务工作台", "Service Workbench"),
  title: pickLabel(localeStore.locale, "服务详情", "Service Detail"),
  subtitle: pickLabel(
    localeStore.locale,
    "在一个页面里统一查看当前服务的资源快照、配置、账务、工单与改配记录，避免在多个页面之间来回切换。",
    "Review resource snapshot, configuration, billing, tickets, and change records in one page."
  ),
  back: pickLabel(localeStore.locale, "返回服务中心", "Back to Services"),
  openOrders: pickLabel(localeStore.locale, "订单中心", "Orders"),
  openInvoices: pickLabel(localeStore.locale, "财务中心", "Finance"),
  openTickets: pickLabel(localeStore.locale, "工单中心", "Tickets"),
  newTicket: pickLabel(localeStore.locale, "提交工单", "Create Ticket"),
  overview: pickLabel(localeStore.locale, "服务概览", "Overview"),
  resources: pickLabel(localeStore.locale, "资源快照", "Resource Snapshot"),
  access: pickLabel(localeStore.locale, "访问与登录", "Access"),
  accessDesc: pickLabel(localeStore.locale, "集中查看登录账号、密码提示和访问地址。", "Review login account, password hint, and access endpoints."),
  network: pickLabel(localeStore.locale, "网络与安全", "Network & Security"),
  networkDesc: pickLabel(localeStore.locale, "快速查看公网地址、地域和安全组信息。", "Quickly review public addresses, region, and security group."),
  timeline: pickLabel(localeStore.locale, "实例时间线", "Instance Timeline"),
  timelineDesc: pickLabel(localeStore.locale, "按时间查看创建、同步、到期和改配变化。", "Review creation, sync, due, and change events in time order."),
  billingDesk: pickLabel(localeStore.locale, "账务工作台", "Billing Desk"),
  billingDeskDesc: pickLabel(localeStore.locale, "围绕订单、账单和到期节点继续处理当前服务。", "Continue handling order, invoice, and due actions around this service."),
  changeDesk: pickLabel(localeStore.locale, "改配工作台", "Change Desk"),
  changeDeskDesc: pickLabel(localeStore.locale, "查看改配总量、待支付和执行异常，并继续跟进。", "Track total changes, waiting payments, and failures from one panel."),
  config: pickLabel(localeStore.locale, "配置项", "Configuration"),
  billing: pickLabel(localeStore.locale, "账务关联", "Billing"),
  changes: pickLabel(localeStore.locale, "改配记录", "Change Orders"),
  noConfig: pickLabel(localeStore.locale, "暂无配置快照", "No configuration snapshot"),
  noChanges: pickLabel(localeStore.locale, "暂无改配记录", "No change orders"),
  serviceNo: pickLabel(localeStore.locale, "服务号", "Service No."),
  productName: pickLabel(localeStore.locale, "商品", "Product"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  provider: pickLabel(localeStore.locale, "渠道", "Provider"),
  nextDueAt: pickLabel(localeStore.locale, "到期时间", "Next Due"),
  orderNo: pickLabel(localeStore.locale, "订单号", "Order No."),
  invoiceNo: pickLabel(localeStore.locale, "账单号", "Invoice No."),
  invoiceStatus: pickLabel(localeStore.locale, "账单状态", "Invoice Status"),
  hostname: pickLabel(localeStore.locale, "主机名", "Hostname"),
  region: pickLabel(localeStore.locale, "地域", "Region"),
  publicIpv4: pickLabel(localeStore.locale, "公网 IPv4", "Public IPv4"),
  publicIpv6: pickLabel(localeStore.locale, "公网 IPv6", "Public IPv6"),
  zone: pickLabel(localeStore.locale, "可用区", "Zone"),
  cpu: pickLabel(localeStore.locale, "CPU", "CPU"),
  memory: pickLabel(localeStore.locale, "内存", "Memory"),
  disk: pickLabel(localeStore.locale, "磁盘", "Disk"),
  bandwidth: pickLabel(localeStore.locale, "带宽", "Bandwidth"),
  loginUser: pickLabel(localeStore.locale, "登录账号", "Login User"),
  passwordHint: pickLabel(localeStore.locale, "密码提示", "Password Hint"),
  securityGroup: pickLabel(localeStore.locale, "安全组", "Security Group"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  execution: pickLabel(localeStore.locale, "执行状态", "Execution"),
  cycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  createdAt: pickLabel(localeStore.locale, "创建时间", "Created At"),
  collaboration: pickLabel(localeStore.locale, "业务协同", "Business Collaboration"),
  collaborationDesc: pickLabel(localeStore.locale, "围绕该服务继续处理订单、账单、工单与改配。", "Continue handling orders, invoices, tickets, and change actions around this service."),
  waitingPayment: pickLabel(localeStore.locale, "待支付改配", "Waiting Payment"),
  executingNow: pickLabel(localeStore.locale, "执行中", "Executing"),
  failedChanges: pickLabel(localeStore.locale, "执行失败", "Failed Changes"),
  configName: pickLabel(localeStore.locale, "配置项", "Configuration"),
  configValue: pickLabel(localeStore.locale, "当前值", "Current Value"),
  configPrice: pickLabel(localeStore.locale, "差价", "Delta"),
  changeTitle: pickLabel(localeStore.locale, "标题", "Title"),
  changeAction: pickLabel(localeStore.locale, "动作", "Action"),
  openInvoice: pickLabel(localeStore.locale, "查看账单", "Open Invoice"),
  loadError: pickLabel(localeStore.locale, "服务详情加载失败", "Failed to load service detail")
}));

const configurationRows = computed(() => detail.value?.service.configuration ?? []);
const changeOrders = computed(() => detail.value?.changeOrders ?? []);
const resource = computed(() => detail.value?.service.resourceSnapshot);
const heroMetrics = computed(() => [
  {
    label: pickLabel(localeStore.locale, "公网 IPv4", "Public IPv4"),
    value: resource.value?.publicIpv4 || "-"
  },
  {
    label: pickLabel(localeStore.locale, "配置项数", "Config Items"),
    value: String(configurationRows.value.length)
  },
  {
    label: pickLabel(localeStore.locale, "改配记录", "Change Orders"),
    value: String(changeOrders.value.length)
  },
  {
    label: pickLabel(localeStore.locale, "账单金额", "Invoice Amount"),
    value: formatPortalMoney(localeStore.locale, detail.value?.invoice?.totalAmount || 0)
  }
]);
const accessRows = computed(() => [
  {
    label: copy.value.loginUser,
    value: resource.value?.loginUsername || "-",
    copyable: Boolean(resource.value?.loginUsername)
  },
  {
    label: copy.value.passwordHint,
    value: resource.value?.passwordHint || "-",
    copyable: Boolean(resource.value?.passwordHint)
  },
  {
    label: copy.value.hostname,
    value: resource.value?.hostname || "-",
    copyable: Boolean(resource.value?.hostname)
  }
]);
const changeStats = computed(() => ({
  total: changeOrders.value.length,
  waitingPayment: changeOrders.value.filter(item => item.executionStatus === "WAITING_PAYMENT").length,
  executing: changeOrders.value.filter(item => item.executionStatus === "EXECUTING").length,
  failed: changeOrders.value.filter(item => item.executionStatus === "EXECUTE_FAILED").length
}));
const networkRows = computed(() => [
  {
    label: copy.value.region,
    value: resource.value?.regionName || "-",
    copyable: false
  },
  {
    label: copy.value.zone,
    value: resource.value?.zoneName || "-",
    copyable: false
  },
  {
    label: copy.value.publicIpv4,
    value: resource.value?.publicIpv4 || "-",
    copyable: Boolean(resource.value?.publicIpv4)
  },
  {
    label: copy.value.publicIpv6,
    value: resource.value?.publicIpv6 || "-",
    copyable: Boolean(resource.value?.publicIpv6)
  },
  {
    label: copy.value.securityGroup,
    value: resource.value?.securityGroup || "-",
    copyable: Boolean(resource.value?.securityGroup)
  }
]);
const serviceTimeline = computed(() => {
  const items: Array<{ time: string; title: string; subtitle: string }> = [];
  if (detail.value?.service.createdAt) {
    items.push({
      time: detail.value.service.createdAt,
      title: pickLabel(localeStore.locale, "服务创建", "Service Created"),
      subtitle: detail.value.service.serviceNo
    });
  }
  if (detail.value?.service.lastSyncAt) {
    items.push({
      time: detail.value.service.lastSyncAt,
      title: pickLabel(localeStore.locale, "最近同步", "Last Sync"),
      subtitle: detail.value.service.syncStatus || "-"
    });
  }
  if (detail.value?.service.nextDueAt) {
    items.push({
      time: detail.value.service.nextDueAt,
      title: pickLabel(localeStore.locale, "下次到期", "Next Due"),
      subtitle: detail.value.invoice?.invoiceNo || detail.value.order?.orderNo || "-"
    });
  }
  for (const item of changeOrders.value.slice(0, 4)) {
    items.push({
      time: item.createdAt,
      title: pickLabel(localeStore.locale, "改配记录", "Change Order"),
      subtitle: `${item.actionName} / ${item.title}`
    });
  }
  return items
    .filter(item => item.time)
    .sort((a, b) => String(b.time).localeCompare(String(a.time)));
});
const resourceHighlights = computed(() => [
  {
    label: copy.value.cpu,
    value: `${resource.value?.cpuCores || 0} Core`
  },
  {
    label: copy.value.memory,
    value: `${resource.value?.memoryGB || 0} GB`
  },
  {
    label: copy.value.disk,
    value: `${(resource.value?.systemDiskGB || 0) + (resource.value?.dataDiskGB || 0)} GB`
  },
  {
    label: copy.value.bandwidth,
    value: `${resource.value?.bandwidthMbps || 0} Mbps`
  }
]);

function executionType(status: string) {
  const mapping: Record<string, "success" | "warning" | "danger" | "info"> = {
    WAITING_PAYMENT: "warning",
    PAID: "info",
    EXECUTING: "info",
    EXECUTED: "success",
    EXECUTE_FAILED: "danger",
    EXECUTE_BLOCKED: "warning",
    REFUNDED: "info"
  };
  return mapping[status] ?? "info";
}

function executionLabel(status: string) {
  const mapping: Record<string, [string, string]> = {
    WAITING_PAYMENT: ["待支付", "Waiting Payment"],
    PAID: ["待执行", "Paid"],
    EXECUTING: ["执行中", "Executing"],
    EXECUTED: ["已执行", "Executed"],
    EXECUTE_FAILED: ["执行失败", "Execution Failed"],
    EXECUTE_BLOCKED: ["执行阻塞", "Blocked"],
    REFUNDED: ["已退款", "Refunded"]
  };
  const pair = mapping[status];
  return pair ? pickLabel(localeStore.locale, pair[0], pair[1]) : status || "-";
}

async function fetchDetail() {
  loading.value = true;
  error.value = "";
  try {
    detail.value = await loadServiceDetail(String(route.params.id));
  } catch (err) {
    detail.value = null;
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

function openInvoice(invoiceNo?: string) {
  if (invoiceNo) {
    void router.push({ path: "/invoices", query: { invoiceNo, serviceId: String(detail.value?.service.id ?? "") } });
    return;
  }
  void router.push({ path: "/invoices", query: { serviceId: String(detail.value?.service.id ?? "") } });
}

function openTicketCenter(action = "") {
  const serviceId = detail.value?.service.id;
  void router.push({
    path: "/tickets",
    query: {
      serviceId: serviceId ? String(serviceId) : undefined,
      action: action || undefined,
      title: action === "create" ? detail.value?.service.productName : undefined
    }
  });
}

function openOrderCenter() {
  const orderNo = detail.value?.order?.orderNo;
  void router.push({ path: "/orders", query: { orderNo: orderNo || undefined } });
}

function openChangeInvoice(item: PortalServiceChangeOrder) {
  openInvoice(item.invoiceNo);
}

async function copyValue(value: string) {
  if (!value || value === "-") return;
  await navigator.clipboard.writeText(value);
  ElMessage.success(pickLabel(localeStore.locale, "已复制到剪贴板", "Copied to clipboard"));
}

onMounted(() => {
  void fetchDetail();
});

watch(
  () => route.params.id,
  () => {
    void fetchDetail();
  }
);
</script>

<template>
  <div class="portal-shell" v-loading="loading">
    <el-alert v-if="error" :title="error" type="error" :closable="false" show-icon />

    <template v-if="detail">
      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <div class="portal-badge">{{ copy.badge }}</div>
            <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
            <p class="portal-subtitle">{{ copy.subtitle }}</p>
          </div>
          <div class="portal-toolbar">
            <el-button @click="router.push('/services')">{{ copy.back }}</el-button>
            <el-button plain @click="openOrderCenter">{{ copy.openOrders }}</el-button>
            <el-button plain @click="openInvoice()">{{ copy.openInvoices }}</el-button>
            <el-button plain @click="openTicketCenter()">{{ copy.openTickets }}</el-button>
            <el-button type="primary" plain @click="openTicketCenter('create')">{{ copy.newTicket }}</el-button>
          </div>
        </div>

        <div class="portal-grid portal-grid--four" style="margin-top: 20px">
          <article class="portal-stat">
            <h3>{{ copy.serviceNo }}</h3>
            <strong>{{ detail.service.serviceNo }}</strong>
          </article>
          <article class="portal-stat">
            <h3>{{ copy.productName }}</h3>
            <strong>{{ detail.service.productName }}</strong>
          </article>
          <article class="portal-stat">
            <h3>{{ copy.status }}</h3>
            <strong>{{ formatPortalServiceStatus(localeStore.locale, detail.service.status) }}</strong>
          </article>
          <article class="portal-stat">
            <h3>{{ copy.nextDueAt }}</h3>
            <strong>{{ detail.service.nextDueAt || "-" }}</strong>
          </article>
        </div>

        <div class="portal-hero-grid" style="margin-top: 18px">
          <div v-for="item in heroMetrics" :key="item.label" class="portal-mini-card">
            <span class="portal-mini-card__label">{{ item.label }}</span>
            <strong class="portal-mini-card__value">{{ item.value }}</strong>
          </div>
        </div>
      </section>

      <section class="portal-grid portal-grid--three">
        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.overview }}</h2>
            </div>
          </div>
          <div class="portal-summary" style="margin-top: 18px">
            <div class="portal-summary-row">
              <span>{{ copy.provider }}</span>
              <strong>{{ detail.service.providerType || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.orderNo }}</span>
              <strong>{{ detail.order?.orderNo || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.invoiceNo }}</span>
              <strong>{{ detail.invoice?.invoiceNo || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.invoiceStatus }}</span>
              <strong>{{ formatPortalInvoiceStatus(localeStore.locale, detail.invoice?.status || "") }}</strong>
            </div>
          </div>
        </article>

        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.collaboration }}</h2>
              <div class="portal-panel__meta">{{ copy.collaborationDesc }}</div>
            </div>
          </div>
          <div class="portal-actions-grid" style="margin-top: 18px; grid-template-columns: 1fr">
            <button type="button" class="portal-action-card" @click="openOrderCenter">
              <strong>{{ copy.openOrders }}</strong>
              <span>{{ detail.order?.orderNo || "-" }}</span>
            </button>
            <button type="button" class="portal-action-card" @click="openInvoice()">
              <strong>{{ copy.openInvoices }}</strong>
              <span>{{ detail.invoice?.invoiceNo || "-" }}</span>
            </button>
            <button type="button" class="portal-action-card" @click="openTicketCenter()">
              <strong>{{ copy.openTickets }}</strong>
              <span>{{ pickLabel(localeStore.locale, "查看围绕该服务的支持记录。", "Review tickets linked to this service.") }}</span>
            </button>
            <button type="button" class="portal-action-card" @click="openTicketCenter('create')">
              <strong>{{ copy.newTicket }}</strong>
              <span>{{ pickLabel(localeStore.locale, "直接带上服务上下文创建工单。", "Create a ticket with service context attached.") }}</span>
            </button>
          </div>
        </article>

        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.resources }}</h2>
            </div>
          </div>
          <div class="portal-hero-grid" style="margin-top: 18px">
            <div v-for="item in resourceHighlights" :key="item.label" class="portal-mini-card">
              <span class="portal-mini-card__label">{{ item.label }}</span>
              <strong class="portal-mini-card__value">{{ item.value }}</strong>
            </div>
          </div>
          <div class="portal-summary" style="margin-top: 18px">
            <div class="portal-summary-row">
              <span>{{ copy.hostname }}</span>
              <strong>{{ resource?.hostname || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.region }}</span>
              <strong>{{ resource?.regionName || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.publicIpv4 }}</span>
              <strong>{{ resource?.publicIpv4 || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.publicIpv6 }}</span>
              <strong>{{ resource?.publicIpv6 || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ pickLabel(localeStore.locale, "操作系统", "Operating System") }}</span>
              <strong>{{ resource?.operatingSystem || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ pickLabel(localeStore.locale, "登录账号", "Login User") }}</span>
              <strong>{{ resource?.loginUsername || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ pickLabel(localeStore.locale, "安全组", "Security Group") }}</span>
              <strong>{{ resource?.securityGroup || "-" }}</strong>
            </div>
          </div>
        </article>
      </section>

      <section class="portal-grid portal-grid--two">
        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.access }}</h2>
              <div class="portal-panel__meta">{{ copy.accessDesc }}</div>
            </div>
          </div>
          <div class="portal-summary" style="margin-top: 18px">
            <div v-for="item in accessRows" :key="item.label" class="portal-summary-row">
              <span>{{ item.label }}</span>
              <div class="portal-toolbar" style="gap: 8px">
                <strong>{{ item.value }}</strong>
                <el-button v-if="item.copyable" link type="primary" @click="copyValue(item.value)">
                  {{ pickLabel(localeStore.locale, "复制", "Copy") }}
                </el-button>
              </div>
            </div>
          </div>
        </article>

        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.network }}</h2>
              <div class="portal-panel__meta">{{ copy.networkDesc }}</div>
            </div>
          </div>
          <div class="portal-summary" style="margin-top: 18px">
            <div v-for="item in networkRows" :key="item.label" class="portal-summary-row">
              <span>{{ item.label }}</span>
              <div class="portal-toolbar" style="gap: 8px">
                <strong>{{ item.value }}</strong>
                <el-button v-if="item.copyable" link type="primary" @click="copyValue(item.value)">
                  {{ pickLabel(localeStore.locale, "复制", "Copy") }}
                </el-button>
              </div>
            </div>
          </div>
        </article>
      </section>

      <section class="portal-grid portal-grid--two">
        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.billingDesk }}</h2>
              <div class="portal-panel__meta">{{ copy.billingDeskDesc }}</div>
            </div>
          </div>
          <div class="portal-actions-grid" style="margin-top: 18px">
            <button type="button" class="portal-action-card" @click="openOrderCenter">
              <strong>{{ copy.openOrders }}</strong>
              <span>{{ detail.order?.orderNo || "-" }}</span>
            </button>
            <button type="button" class="portal-action-card" @click="openInvoice()">
              <strong>{{ copy.openInvoices }}</strong>
              <span>{{ detail.invoice?.invoiceNo || "-" }}</span>
            </button>
            <button type="button" class="portal-action-card" @click="router.push('/wallet/transactions')">
              <strong>{{ pickLabel(localeStore.locale, "钱包中心", "Wallet") }}</strong>
              <span>{{ pickLabel(localeStore.locale, "继续查看余额、授信和资金流水。", "Continue reviewing balance, credit, and transactions.") }}</span>
            </button>
            <button type="button" class="portal-action-card" @click="openTicketCenter('create')">
              <strong>{{ copy.newTicket }}</strong>
              <span>{{ pickLabel(localeStore.locale, "围绕账务或实例问题继续提单。", "Create a ticket for finance or instance issues.") }}</span>
            </button>
          </div>
        </article>

        <article class="portal-card">
          <div class="portal-card-head">
            <div>
              <h2 class="portal-panel__title">{{ copy.changeDesk }}</h2>
              <div class="portal-panel__meta">{{ copy.changeDeskDesc }}</div>
            </div>
          </div>
          <div class="portal-grid portal-grid--four" style="margin-top: 18px">
            <article class="portal-stat"><h3>{{ copy.changes }}</h3><strong>{{ changeStats.total }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.waitingPayment }}</h3><strong>{{ changeStats.waitingPayment }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.executingNow }}</h3><strong>{{ changeStats.executing }}</strong></article>
            <article class="portal-stat"><h3>{{ copy.failedChanges }}</h3><strong>{{ changeStats.failed }}</strong></article>
          </div>
        </article>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.timeline }}</h2>
            <div class="portal-panel__meta">{{ copy.timelineDesc }}</div>
          </div>
        </div>
        <div v-if="serviceTimeline.length" class="portal-list" style="margin-top: 18px">
          <div v-for="item in serviceTimeline" :key="`${item.time}-${item.title}-${item.subtitle}`" class="portal-list-item">
            <div class="portal-list-item__meta">
              <div class="portal-list-item__title">{{ item.title }}</div>
              <div class="portal-list-item__desc">{{ item.subtitle }}</div>
            </div>
            <strong>{{ item.time }}</strong>
          </div>
        </div>
        <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noChanges }}</div>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.config }}</h2>
          </div>
        </div>
        <el-table v-if="configurationRows.length" :data="configurationRows" border style="margin-top: 18px">
          <el-table-column prop="name" :label="copy.configName" min-width="180" />
          <el-table-column :label="copy.configValue" min-width="220">
            <template #default="{ row }">{{ formatPortalConfiguration([row], localeStore.locale) }}</template>
          </el-table-column>
          <el-table-column :label="copy.configPrice" min-width="120">
            <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.priceDelta) }}</template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noConfig }}</div>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.billing }}</h2>
          </div>
        </div>
        <div class="portal-summary" style="margin-top: 18px">
          <div class="portal-summary-row">
            <span>{{ copy.invoiceNo }}</span>
            <strong>{{ detail.invoice?.invoiceNo || "-" }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.amount }}</span>
            <strong>{{ formatPortalMoney(localeStore.locale, detail.invoice?.totalAmount || 0) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.cycle }}</span>
            <strong>{{ formatPortalBillingCycle(localeStore.locale, detail.invoice?.billingCycle) }}</strong>
          </div>
          <div class="portal-summary-row">
            <span>{{ copy.invoiceStatus }}</span>
            <el-tag :type="portalTagTypeByStatus(detail.invoice?.status)">
              {{ formatPortalInvoiceStatus(localeStore.locale, detail.invoice?.status) }}
            </el-tag>
          </div>
        </div>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.changes }}</h2>
          </div>
        </div>
        <el-table v-if="changeOrders.length" :data="changeOrders" border style="margin-top: 18px">
          <el-table-column prop="title" :label="copy.changeTitle" min-width="220" />
          <el-table-column prop="actionName" :label="copy.changeAction" min-width="140" />
          <el-table-column :label="copy.amount" min-width="120">
            <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
          </el-table-column>
          <el-table-column :label="copy.cycle" min-width="120">
            <template #default="{ row }">{{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}</template>
          </el-table-column>
          <el-table-column :label="copy.execution" min-width="120">
            <template #default="{ row }">
              <el-tag :type="executionType(row.executionStatus)">{{ executionLabel(row.executionStatus) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" :label="copy.createdAt" min-width="160" />
          <el-table-column :label="copy.action" min-width="120" fixed="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openChangeInvoice(row)">{{ copy.openInvoice }}</el-button>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state" style="margin-top: 18px">{{ copy.noChanges }}</div>
      </section>
    </template>
  </div>
</template>
