<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { loadServices, type PortalService, type PortalServiceConfigSelection } from "@/api/portal";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import {
  formatPortalBillingCycle,
  formatPortalConfiguration,
  formatPortalServiceStatus,
  portalTagTypeByStatus
} from "@/utils/business";

const router = useRouter();
const localeStore = useLocaleStore();
const loading = ref(false);
const error = ref("");
const services = ref<PortalService[]>([]);
const keyword = ref("");
const status = ref("");
const providerType = ref("");
const dueFilter = ref("");

const providerOptions = computed(() =>
  Array.from(new Set(services.value.map(item => item.providerType).filter(Boolean))) as string[]
);

const filteredServices = computed(() =>
  services.value.filter(item => {
    const text = keyword.value.trim().toLowerCase();
    const matchesKeyword =
      !text ||
      item.serviceNo.toLowerCase().includes(text) ||
      item.productName.toLowerCase().includes(text) ||
      String(item.providerType ?? "").toLowerCase().includes(text);
    const matchesStatus = !status.value || item.status === status.value;
    const matchesProvider = !providerType.value || item.providerType === providerType.value;
    const matchesDue =
      !dueFilter.value ||
      (dueFilter.value === "7D" && withinDays(item.nextDueAt, 7)) ||
      (dueFilter.value === "30D" && withinDays(item.nextDueAt, 30));
    return matchesKeyword && matchesStatus && matchesProvider && matchesDue;
  })
);

const stats = computed(() => ({
  total: services.value.length,
  active: services.value.filter(item => item.status === "ACTIVE").length,
  dueSoon: services.value.filter(item => withinDays(item.nextDueAt, 7)).length,
  suspended: services.value.filter(item => item.status === "SUSPENDED").length
}));

const latestServices = computed(() =>
  filteredServices.value
    .slice()
    .sort((a, b) => String(a.nextDueAt || "").localeCompare(String(b.nextDueAt || "")))
    .slice(0, 5)
);

const providerSummary = computed(() => {
  const summary = new Map<string, number>();
  for (const item of filteredServices.value) {
    const key = item.providerType || "LOCAL";
    summary.set(key, (summary.get(key) || 0) + 1);
  }
  return Array.from(summary.entries()).map(([provider, count]) => ({ provider, count }));
});

const copy = computed(() => ({
  badge: pickLabel(localeStore.locale, "服务中心", "Services"),
  title: pickLabel(localeStore.locale, "服务中心", "Services"),
  subtitle: pickLabel(
    localeStore.locale,
    "围绕已购服务统一查看资源、账单、订单、工单与改配记录，是客户侧的核心工作台。",
    "Review resources, invoices, orders, tickets, and change records from one service workbench."
  ),
  total: pickLabel(localeStore.locale, "服务总数", "Total Services"),
  active: pickLabel(localeStore.locale, "运行中", "Active"),
  dueSoon: pickLabel(localeStore.locale, "7 天内到期", "Due in 7 Days"),
  suspended: pickLabel(localeStore.locale, "已暂停", "Suspended"),
  keyword: pickLabel(localeStore.locale, "搜索服务号 / 商品 / 渠道", "Search service, product, or provider"),
  status: pickLabel(localeStore.locale, "服务状态", "Status"),
  provider: pickLabel(localeStore.locale, "渠道类型", "Provider"),
  due: pickLabel(localeStore.locale, "到期筛选", "Due Filter"),
  due7: pickLabel(localeStore.locale, "7 天内", "Within 7 days"),
  due30: pickLabel(localeStore.locale, "30 天内", "Within 30 days"),
  visible: pickLabel(localeStore.locale, "当前显示", "Visible"),
  reset: pickLabel(localeStore.locale, "重置筛选", "Reset"),
  latest: pickLabel(localeStore.locale, "近期关注服务", "Focus Services"),
  latestDesc: pickLabel(localeStore.locale, "优先处理即将到期和状态异常的服务。", "Prioritize services that are expiring soon or need attention."),
  providerSummary: pickLabel(localeStore.locale, "渠道分布", "Provider Mix"),
  providerSummaryDesc: pickLabel(localeStore.locale, "从渠道视角快速判断当前服务分布。", "Quickly understand service distribution by provider."),
  serviceNo: pickLabel(localeStore.locale, "服务号", "Service No."),
  productName: pickLabel(localeStore.locale, "商品", "Product"),
  providerTypeLabel: pickLabel(localeStore.locale, "渠道", "Provider"),
  configuration: pickLabel(localeStore.locale, "配置摘要", "Configuration"),
  billingCycle: pickLabel(localeStore.locale, "周期", "Billing Cycle"),
  nextDueAt: pickLabel(localeStore.locale, "到期时间", "Next Due"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  detail: pickLabel(localeStore.locale, "详情", "Detail"),
  order: pickLabel(localeStore.locale, "订单", "Order"),
  invoice: pickLabel(localeStore.locale, "账单", "Invoice"),
  ticket: pickLabel(localeStore.locale, "工单", "Ticket"),
  empty: pickLabel(localeStore.locale, "暂无匹配的服务记录。", "No matching services."),
  loadError: pickLabel(localeStore.locale, "服务列表加载失败", "Failed to load services")
}));

function withinDays(value: string, days: number) {
  if (!value) return false;
  const target = new Date(value);
  if (Number.isNaN(target.getTime())) return false;
  const diff = target.getTime() - Date.now();
  return diff >= 0 && diff <= days * 24 * 60 * 60 * 1000;
}

async function fetchServices() {
  loading.value = true;
  error.value = "";
  try {
    services.value = await loadServices();
  } catch (err) {
    error.value = err instanceof Error ? err.message : copy.value.loadError;
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  keyword.value = "";
  status.value = "";
  providerType.value = "";
  dueFilter.value = "";
}

function goDetail(item: PortalService) {
  void router.push(`/services/${item.id}`);
}

function goOrders(item: PortalService) {
  if (item.orderId) {
    void router.push(`/orders/${item.orderId}`);
    return;
  }
  void router.push({ path: "/orders", query: { keyword: item.productName } });
}

function goInvoices(item: PortalService) {
  if (item.invoiceId) {
    void router.push(`/invoices/${item.invoiceId}`);
    return;
  }
  void router.push({ path: "/invoices", query: { serviceId: String(item.id) } });
}

function goTickets(item: PortalService) {
  void router.push({ path: "/tickets", query: { serviceId: String(item.id) } });
}

function billingCycleLabel(item: PortalService) {
  const cycle =
    item.configuration?.find((config: PortalServiceConfigSelection) => config.code === "billing_cycle")?.value ||
    "monthly";
  return formatPortalBillingCycle(localeStore.locale, cycle);
}

onMounted(() => {
  void fetchServices();
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

      <div class="portal-grid portal-grid--four" style="margin-top: 20px">
        <article class="portal-stat">
          <h3>{{ copy.total }}</h3>
          <strong>{{ stats.total }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.active }}</h3>
          <strong>{{ stats.active }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.dueSoon }}</h3>
          <strong>{{ stats.dueSoon }}</strong>
        </article>
        <article class="portal-stat">
          <h3>{{ copy.suspended }}</h3>
          <strong>{{ stats.suspended }}</strong>
        </article>
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
          <div v-if="latestServices.length" class="portal-list" style="padding: 18px 20px 20px">
            <div v-for="item in latestServices" :key="item.id" class="portal-list-item">
              <div class="portal-list-item__meta">
                <div class="portal-list-item__title">{{ item.productName }}</div>
                <div class="portal-list-item__desc">{{ copy.serviceNo }}：{{ item.serviceNo }}</div>
                <div class="portal-list-item__desc">{{ copy.nextDueAt }}：{{ item.nextDueAt || "-" }}</div>
              </div>
              <div class="portal-toolbar">
                <el-tag :type="portalTagTypeByStatus(item.status)">
                  {{ formatPortalServiceStatus(localeStore.locale, item.status) }}
                </el-tag>
                <el-button type="primary" link @click="goDetail(item)">{{ copy.detail }}</el-button>
              </div>
            </div>
          </div>
          <div v-else class="portal-empty-state" style="margin: 18px 20px 20px">{{ copy.empty }}</div>
        </article>

        <article class="portal-card" style="padding: 0">
          <div style="padding: 20px 20px 0">
            <div class="portal-card-head">
              <div>
                <h2 class="portal-panel__title">{{ copy.providerSummary }}</h2>
                <div class="portal-panel__meta">{{ copy.providerSummaryDesc }}</div>
              </div>
            </div>
          </div>
          <div v-if="providerSummary.length" class="portal-list" style="padding: 18px 20px 20px">
            <div v-for="item in providerSummary" :key="item.provider" class="portal-list-item">
              <div class="portal-list-item__meta">
                <div class="portal-list-item__title">{{ item.provider }}</div>
                <div class="portal-list-item__desc">{{ copy.visible }}：{{ item.count }}</div>
              </div>
              <strong>{{ item.count }}</strong>
            </div>
          </div>
          <div v-else class="portal-empty-state" style="margin: 18px 20px 20px">{{ copy.empty }}</div>
        </article>
      </section>

      <div class="portal-toolbar" style="margin-top: 20px">
        <el-input v-model="keyword" :placeholder="copy.keyword" clearable />
        <el-select v-model="status" :placeholder="copy.status" clearable style="width: 160px">
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'PENDING')" value="PENDING" />
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'ACTIVE')" value="ACTIVE" />
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'SUSPENDED')" value="SUSPENDED" />
          <el-option :label="formatPortalServiceStatus(localeStore.locale, 'TERMINATED')" value="TERMINATED" />
        </el-select>
        <el-select v-model="providerType" :placeholder="copy.provider" clearable style="width: 180px">
          <el-option v-for="item in providerOptions" :key="item" :label="item" :value="item" />
        </el-select>
        <el-select v-model="dueFilter" :placeholder="copy.due" clearable style="width: 160px">
          <el-option :label="copy.due7" value="7D" />
          <el-option :label="copy.due30" value="30D" />
        </el-select>
        <el-button @click="resetFilters">{{ copy.reset }}</el-button>
        <div class="portal-summary" style="margin-left: auto; min-width: 180px">
          <div class="portal-summary-row">
            <span>{{ copy.visible }}</span>
            <strong>{{ filteredServices.length }}</strong>
          </div>
        </div>
      </div>
    </section>

    <section class="portal-table-card">
      <el-table v-if="filteredServices.length" :data="filteredServices" border>
        <el-table-column prop="serviceNo" :label="copy.serviceNo" min-width="160" />
        <el-table-column prop="productName" :label="copy.productName" min-width="220" />
        <el-table-column prop="providerType" :label="copy.providerTypeLabel" min-width="140" />
        <el-table-column :label="copy.configuration" min-width="220" show-overflow-tooltip>
          <template #default="{ row }">{{ formatPortalConfiguration(row.configuration, localeStore.locale) }}</template>
        </el-table-column>
        <el-table-column :label="copy.billingCycle" min-width="120">
          <template #default="{ row }">{{ billingCycleLabel(row) }}</template>
        </el-table-column>
        <el-table-column :label="copy.nextDueAt" min-width="140">
          <template #default="{ row }">{{ row.nextDueAt || "-" }}</template>
        </el-table-column>
        <el-table-column :label="copy.status" min-width="120">
          <template #default="{ row }">
            <el-tag :type="portalTagTypeByStatus(row.status)">
              {{ formatPortalServiceStatus(localeStore.locale, row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="copy.action" min-width="240" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="goDetail(row)">{{ copy.detail }}</el-button>
            <el-button link @click="goOrders(row)">{{ copy.order }}</el-button>
            <el-button link @click="goInvoices(row)">{{ copy.invoice }}</el-button>
            <el-button link @click="goTickets(row)">{{ copy.ticket }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div v-else class="portal-empty-state" style="margin-bottom: 12px">{{ copy.empty }}</div>
    </section>
  </div>
</template>
