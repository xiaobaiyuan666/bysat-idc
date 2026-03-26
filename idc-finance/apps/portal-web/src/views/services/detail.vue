<script setup lang="ts">
import { computed, onMounted, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import {
  loadServiceDetail,
  type PortalServiceChangeOrder,
  type PortalServiceDetailResponse
} from "@/api/portal";
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

const loading = ref(true);
const error = ref("");
const detail = ref<PortalServiceDetailResponse | null>(null);

const copy = computed(() => ({
  title: pickLabel(localeStore.locale, "服务工作台", "Service Workbench"),
  subtitle: pickLabel(
    localeStore.locale,
    "围绕当前服务统一查看账单、工单、资源快照和改配记录，减少在多个页面之间来回切换。",
    "Review billing, tickets, resource snapshots, and change records for this service in one place."
  ),
  back: pickLabel(localeStore.locale, "返回服务列表", "Back to Services"),
  orders: pickLabel(localeStore.locale, "查看订单", "Open Orders"),
  invoices: pickLabel(localeStore.locale, "查看账单", "Open Invoices"),
  ticketCreate: pickLabel(localeStore.locale, "提交工单", "Create Ticket"),
  ticketList: pickLabel(localeStore.locale, "查看工单", "View Tickets"),
  payInvoice: pickLabel(localeStore.locale, "去账单处理", "Open Invoice"),
  serviceNo: pickLabel(localeStore.locale, "服务编号", "Service No."),
  product: pickLabel(localeStore.locale, "商品", "Product"),
  status: pickLabel(localeStore.locale, "状态", "Status"),
  nextDueAt: pickLabel(localeStore.locale, "下次到期", "Next Due"),
  providerType: pickLabel(localeStore.locale, "渠道类型", "Provider"),
  providerResourceId: pickLabel(localeStore.locale, "远端资源", "Remote ID"),
  syncStatus: pickLabel(localeStore.locale, "同步状态", "Sync Status"),
  syncMessage: pickLabel(localeStore.locale, "同步说明", "Sync Message"),
  orderNo: pickLabel(localeStore.locale, "关联订单", "Order"),
  invoiceNo: pickLabel(localeStore.locale, "关联账单", "Invoice"),
  billingCycle: pickLabel(localeStore.locale, "计费周期", "Billing Cycle"),
  createdAt: pickLabel(localeStore.locale, "开通时间", "Created At"),
  updatedAt: pickLabel(localeStore.locale, "更新时间", "Updated At"),
  region: pickLabel(localeStore.locale, "地域", "Region"),
  ipAddress: pickLabel(localeStore.locale, "公网 IP", "Public IP"),
  overview: pickLabel(localeStore.locale, "服务概览", "Overview"),
  business: pickLabel(localeStore.locale, "业务联动", "Business Links"),
  resource: pickLabel(localeStore.locale, "资源快照", "Resource Snapshot"),
  config: pickLabel(localeStore.locale, "配置项快照", "Configuration"),
  changes: pickLabel(localeStore.locale, "改配记录", "Change Orders"),
  changePending: pickLabel(localeStore.locale, "待处理改配", "Pending Changes"),
  changeCount: pickLabel(localeStore.locale, "改配总数", "Change Count"),
  businessHint: pickLabel(
    localeStore.locale,
    "围绕当前服务快速跳转到订单、账单和工单中心。",
    "Quick jumps to orders, invoices, and tickets for this service."
  ),
  invoiceHint: pickLabel(
    localeStore.locale,
    "如果账单未支付，优先前往账单中心处理收款，再继续变更或工单流程。",
    "If the invoice is unpaid, handle billing first before continuing changes or support."
  ),
  orderCenter: pickLabel(localeStore.locale, "订单中心", "Order Center"),
  invoiceCenter: pickLabel(localeStore.locale, "账单中心", "Invoice Center"),
  ticketCenter: pickLabel(localeStore.locale, "工单中心", "Ticket Center"),
  resourceCenter: pickLabel(localeStore.locale, "资源概况", "Resource Overview"),
  hostname: pickLabel(localeStore.locale, "主机名", "Hostname"),
  os: pickLabel(localeStore.locale, "系统", "OS"),
  cpu: pickLabel(localeStore.locale, "CPU", "CPU"),
  memory: pickLabel(localeStore.locale, "内存", "Memory"),
  systemDisk: pickLabel(localeStore.locale, "系统盘", "System Disk"),
  dataDisk: pickLabel(localeStore.locale, "数据盘", "Data Disk"),
  bandwidth: pickLabel(localeStore.locale, "带宽", "Bandwidth"),
  loginUser: pickLabel(localeStore.locale, "登录用户", "Login User"),
  securityGroup: pickLabel(localeStore.locale, "安全组", "Security Group"),
  ipv6: pickLabel(localeStore.locale, "IPv6", "IPv6"),
  configName: pickLabel(localeStore.locale, "配置项", "Option"),
  configValue: pickLabel(localeStore.locale, "当前值", "Value"),
  priceDelta: pickLabel(localeStore.locale, "差价", "Delta"),
  changeAction: pickLabel(localeStore.locale, "改配动作", "Action"),
  changeTitle: pickLabel(localeStore.locale, "标题", "Title"),
  amount: pickLabel(localeStore.locale, "金额", "Amount"),
  invoiceStatus: pickLabel(localeStore.locale, "账单状态", "Invoice Status"),
  executionStatus: pickLabel(localeStore.locale, "执行状态", "Execution"),
  changeCreatedAt: pickLabel(localeStore.locale, "创建时间", "Created At"),
  action: pickLabel(localeStore.locale, "操作", "Action"),
  openInvoice: pickLabel(localeStore.locale, "查看账单", "Open Invoice"),
  openTickets: pickLabel(localeStore.locale, "查看关联工单", "Open Tickets"),
  emptyConfig: pickLabel(localeStore.locale, "当前服务没有配置项快照。", "No configuration snapshot."),
  emptyChanges: pickLabel(localeStore.locale, "当前服务还没有改配记录。", "No change orders."),
  loadError: pickLabel(localeStore.locale, "服务详情加载失败", "Failed to load service detail")
}));

const pendingChangeCount = computed(
  () =>
    detail.value?.changeOrders.filter(
      item => item.status === "UNPAID" || item.executionStatus === "WAITING_PAYMENT"
    ).length ?? 0
);

const linkedInvoiceUnpaid = computed(() => detail.value?.invoice?.status === "UNPAID");

function providerTypeLabel(type?: string) {
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: pickLabel(localeStore.locale, "魔方云", "Mofang Cloud"),
    ZJMF_API: pickLabel(localeStore.locale, "上游财务", "Upstream Finance"),
    WHMCS: "WHMCS",
    RESOURCE: pickLabel(localeStore.locale, "资源池", "Resource Pool"),
    MANUAL: pickLabel(localeStore.locale, "人工交付", "Manual Delivery"),
    LOCAL: pickLabel(localeStore.locale, "本地模块", "Local")
  };
  return mapping[type || ""] ?? (type || "-");
}

function changeActionLabel(action: string) {
  const mapping: Record<string, [string, string]> = {
    "add-ipv4": ["新增 IPv4", "Add IPv4"],
    "add-ipv6": ["新增 IPv6", "Add IPv6"],
    "add-disk": ["新增数据盘", "Add Disk"],
    "resize-disk": ["扩容磁盘", "Resize Disk"],
    "create-snapshot": ["创建快照", "Create Snapshot"],
    "create-backup": ["创建备份", "Create Backup"]
  };
  const pair = mapping[action];
  return pair ? pickLabel(localeStore.locale, pair[0], pair[1]) : action || "-";
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

function executionTagType(status: string) {
  const mapping: Record<string, "success" | "warning" | "danger" | "info" | "primary"> = {
    WAITING_PAYMENT: "warning",
    PAID: "info",
    EXECUTING: "primary",
    EXECUTED: "success",
    EXECUTE_FAILED: "danger",
    EXECUTE_BLOCKED: "warning",
    REFUNDED: "info"
  };
  return mapping[status] ?? "info";
}

function invoiceStatusTag(status: string) {
  return portalTagTypeByStatus(status);
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

function openOrderWorkbench() {
  const orderNo = detail.value?.order?.orderNo;
  void router.push(orderNo ? { path: "/orders", query: { orderNo } } : "/orders");
}

function openInvoiceWorkbench(invoiceNo?: string) {
  const targetInvoiceNo = invoiceNo || detail.value?.invoice?.invoiceNo;
  const serviceId = detail.value?.service?.id;
  void router.push(
    targetInvoiceNo
      ? { path: "/invoices", query: { invoiceNo: targetInvoiceNo, serviceId: serviceId ? String(serviceId) : undefined } }
      : serviceId
        ? { path: "/invoices", query: { serviceId: String(serviceId) } }
        : "/invoices"
  );
}

function openTicketCreate() {
  const serviceId = detail.value?.service?.id;
  void router.push(
    serviceId
      ? {
          path: "/tickets",
          query: {
            action: "create",
            serviceId: String(serviceId),
            title: detail.value?.service?.productName || undefined
          }
        }
      : "/tickets"
  );
}

function openTicketList() {
  const serviceId = detail.value?.service?.id;
  void router.push(serviceId ? { path: "/tickets", query: { serviceId: String(serviceId) } } : "/tickets");
}

function openChangeInvoice(change: PortalServiceChangeOrder) {
  openInvoiceWorkbench(change.invoiceNo);
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
            <div class="portal-badge">{{ detail.service.serviceNo }}</div>
            <h1 class="portal-title" style="margin-top: 12px">{{ copy.title }}</h1>
            <p class="portal-subtitle">{{ copy.subtitle }}</p>
          </div>
          <div class="portal-toolbar" style="margin: 0">
            <el-button @click="router.push('/services')">{{ copy.back }}</el-button>
            <el-button plain @click="openOrderWorkbench">{{ copy.orders }}</el-button>
            <el-button plain @click="openInvoiceWorkbench()">{{ copy.invoices }}</el-button>
            <el-button type="primary" plain @click="openTicketCreate">{{ copy.ticketCreate }}</el-button>
          </div>
        </div>

        <div class="portal-grid portal-grid--four" style="margin-top: 20px">
          <div class="portal-stat">
            <h3>{{ copy.serviceNo }}</h3>
            <strong>{{ detail.service.serviceNo }}</strong>
          </div>
          <div class="portal-stat">
            <h3>{{ copy.status }}</h3>
            <strong>{{ formatPortalServiceStatus(localeStore.locale, detail.service.status) }}</strong>
          </div>
          <div class="portal-stat">
            <h3>{{ copy.nextDueAt }}</h3>
            <strong>{{ detail.service.nextDueAt || "-" }}</strong>
          </div>
          <div class="portal-stat">
            <h3>{{ copy.changePending }}</h3>
            <strong>{{ pendingChangeCount }}/{{ detail.changeOrders.length }}</strong>
          </div>
        </div>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.business }}</h2>
            <div class="portal-panel__meta">{{ copy.businessHint }}</div>
          </div>
        </div>
        <div class="service-link-grid">
          <div class="service-link-card">
            <div class="service-link-card__label">{{ copy.orderCenter }}</div>
            <strong>{{ detail.order?.orderNo || "-" }}</strong>
            <div class="service-link-card__meta">{{ detail.order?.productName || detail.service.productName }}</div>
            <div class="service-link-card__actions">
              <el-button plain @click="openOrderWorkbench">{{ copy.orders }}</el-button>
            </div>
          </div>
          <div class="service-link-card">
            <div class="service-link-card__label">{{ copy.invoiceCenter }}</div>
            <strong>{{ detail.invoice?.invoiceNo || "-" }}</strong>
            <div class="service-link-card__meta">
              <el-tag
                v-if="detail.invoice?.status"
                :type="portalTagTypeByStatus(detail.invoice.status)"
                effect="light"
              >
                {{ formatPortalInvoiceStatus(localeStore.locale, detail.invoice.status) }}
              </el-tag>
            </div>
            <div class="service-link-card__actions">
              <el-button :type="linkedInvoiceUnpaid ? 'warning' : 'primary'" plain @click="openInvoiceWorkbench()">
                {{ linkedInvoiceUnpaid ? copy.payInvoice : copy.invoices }}
              </el-button>
            </div>
          </div>
          <div class="service-link-card">
            <div class="service-link-card__label">{{ copy.ticketCenter }}</div>
            <strong>{{ copy.changeCount }}</strong>
            <div class="service-link-card__meta">{{ pendingChangeCount }} / {{ detail.changeOrders.length }}</div>
            <div class="service-link-card__actions">
              <el-button plain @click="openTicketList">{{ copy.openTickets }}</el-button>
              <el-button type="primary" plain @click="openTicketCreate">{{ copy.ticketCreate }}</el-button>
            </div>
          </div>
          <div class="service-link-card">
            <div class="service-link-card__label">{{ copy.resourceCenter }}</div>
            <strong>{{ providerTypeLabel(detail.service.providerType) }}</strong>
            <div class="service-link-card__meta">{{ detail.service.providerResourceId || "-" }}</div>
            <div class="service-link-card__actions">
              <el-button plain @click="openInvoiceWorkbench()">{{ copy.invoices }}</el-button>
              <el-button plain @click="openTicketList">{{ copy.ticketList }}</el-button>
            </div>
          </div>
        </div>
        <el-alert
          style="margin-top: 16px"
          :title="copy.invoiceHint"
          :type="linkedInvoiceUnpaid ? 'warning' : 'info'"
          :closable="false"
          show-icon
        />
      </section>

      <section class="portal-card">
        <div class="portal-grid portal-grid--two">
          <div class="portal-summary">
            <div class="portal-card-head" style="padding: 0 0 16px">
              <div>
                <h2 class="portal-panel__title">{{ copy.overview }}</h2>
              </div>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.product }}</span>
              <strong>{{ detail.service.productName }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.status }}</span>
              <strong>
                <el-tag :type="portalTagTypeByStatus(detail.service.status)" effect="light">
                  {{ formatPortalServiceStatus(localeStore.locale, detail.service.status) }}
                </el-tag>
              </strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.providerType }}</span>
              <strong>{{ providerTypeLabel(detail.service.providerType) }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.providerResourceId }}</span>
              <strong>{{ detail.service.providerResourceId || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.region }}</span>
              <strong>{{ detail.service.regionName || detail.service.resourceSnapshot.regionName || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.ipAddress }}</span>
              <strong>{{ detail.service.ipAddress || detail.service.resourceSnapshot.publicIpv4 || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.syncStatus }}</span>
              <strong>{{ detail.service.syncStatus || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.syncMessage }}</span>
              <strong>{{ detail.service.syncMessage || "-" }}</strong>
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
              <span>{{ copy.billingCycle }}</span>
              <strong>
                {{ formatPortalBillingCycle(localeStore.locale, detail.invoice?.billingCycle || detail.order?.billingCycle) }}
              </strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.createdAt }}</span>
              <strong>{{ detail.service.createdAt || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.updatedAt }}</span>
              <strong>{{ detail.service.updatedAt || "-" }}</strong>
            </div>
          </div>

          <div class="portal-summary">
            <div class="portal-card-head" style="padding: 0 0 16px">
              <div>
                <h2 class="portal-panel__title">{{ copy.resource }}</h2>
              </div>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.hostname }}</span>
              <strong>{{ detail.service.resourceSnapshot.hostname || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.os }}</span>
              <strong>{{ detail.service.resourceSnapshot.operatingSystem || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.cpu }}</span>
              <strong>{{ detail.service.resourceSnapshot.cpuCores || 0 }} Core</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.memory }}</span>
              <strong>{{ detail.service.resourceSnapshot.memoryGB || 0 }} GB</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.systemDisk }}</span>
              <strong>{{ detail.service.resourceSnapshot.systemDiskGB || 0 }} GB</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.dataDisk }}</span>
              <strong>{{ detail.service.resourceSnapshot.dataDiskGB || 0 }} GB</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.bandwidth }}</span>
              <strong>{{ detail.service.resourceSnapshot.bandwidthMbps || 0 }} Mbps</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.loginUser }}</span>
              <strong>{{ detail.service.resourceSnapshot.loginUsername || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.securityGroup }}</span>
              <strong>{{ detail.service.resourceSnapshot.securityGroup || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.ipAddress }}</span>
              <strong>{{ detail.service.resourceSnapshot.publicIpv4 || "-" }}</strong>
            </div>
            <div class="portal-summary-row">
              <span>{{ copy.ipv6 }}</span>
              <strong>{{ detail.service.resourceSnapshot.publicIpv6 || "-" }}</strong>
            </div>
          </div>
        </div>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.config }}</h2>
            <div class="portal-panel__meta">
              {{ formatPortalConfiguration(detail.service.configuration, localeStore.locale) }}
            </div>
          </div>
        </div>
        <el-table v-if="detail.service.configuration?.length" :data="detail.service.configuration" border>
          <el-table-column prop="name" :label="copy.configName" min-width="180" />
          <el-table-column :label="copy.configValue" min-width="220">
            <template #default="{ row }">{{ row.valueLabel || row.value || "-" }}</template>
          </el-table-column>
          <el-table-column :label="copy.priceDelta" min-width="120">
            <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.priceDelta) }}</template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state">{{ copy.emptyConfig }}</div>
      </section>

      <section class="portal-card">
        <div class="portal-card-head">
          <div>
            <h2 class="portal-panel__title">{{ copy.changes }}</h2>
            <div class="portal-panel__meta">
              {{ detail.invoice?.status === "UNPAID" ? copy.payInvoice : formatPortalInvoiceStatus(localeStore.locale, detail.invoice?.status) }}
            </div>
          </div>
        </div>

        <el-table v-if="detail.changeOrders.length" :data="detail.changeOrders" border>
          <el-table-column prop="invoiceNo" :label="copy.invoiceNo" min-width="160" />
          <el-table-column prop="orderNo" :label="copy.orderNo" min-width="160" />
          <el-table-column :label="copy.changeAction" min-width="140">
            <template #default="{ row }">{{ changeActionLabel(row.actionName) }}</template>
          </el-table-column>
          <el-table-column prop="title" :label="copy.changeTitle" min-width="220" show-overflow-tooltip />
          <el-table-column :label="copy.amount" min-width="120">
            <template #default="{ row }">{{ formatPortalMoney(localeStore.locale, row.amount) }}</template>
          </el-table-column>
          <el-table-column :label="copy.invoiceStatus" min-width="120">
            <template #default="{ row }">
              <el-tag :type="invoiceStatusTag(row.status)" effect="light">
                {{ formatPortalInvoiceStatus(localeStore.locale, row.status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="copy.executionStatus" min-width="140">
            <template #default="{ row }">
              <el-tag :type="executionTagType(row.executionStatus)" effect="light">
                {{ executionLabel(row.executionStatus) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="copy.billingCycle" min-width="120">
            <template #default="{ row }">
              {{ formatPortalBillingCycle(localeStore.locale, row.billingCycle) }}
            </template>
          </el-table-column>
          <el-table-column prop="createdAt" :label="copy.changeCreatedAt" min-width="180" />
          <el-table-column :label="copy.action" min-width="160" fixed="right">
            <template #default="{ row }">
              <div class="service-table-actions">
                <el-button type="primary" link @click="openChangeInvoice(row)">
                  {{ copy.openInvoice }}
                </el-button>
                <el-button type="primary" link @click="openTicketList">
                  {{ copy.ticketList }}
                </el-button>
              </div>
            </template>
          </el-table-column>
        </el-table>
        <div v-else class="portal-empty-state">{{ copy.emptyChanges }}</div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.service-link-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 16px;
}

.service-link-card {
  padding: 18px;
  border-radius: 20px;
  border: 1px solid rgba(15, 23, 42, 0.08);
  background:
    radial-gradient(circle at top right, rgba(14, 165, 233, 0.08), transparent 45%),
    linear-gradient(180deg, rgba(255, 255, 255, 0.96), rgba(248, 250, 252, 0.94));
  display: grid;
  gap: 10px;
}

.service-link-card strong {
  font-size: 20px;
  color: #0f172a;
}

.service-link-card__label {
  font-size: 13px;
  color: #64748b;
}

.service-link-card__meta {
  min-height: 24px;
  color: #475569;
  font-size: 13px;
}

.service-link-card__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.service-table-actions {
  display: flex;
  gap: 10px;
  align-items: center;
}

@media (max-width: 1280px) {
  .service-link-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 768px) {
  .service-link-grid {
    grid-template-columns: 1fr;
  }
}
</style>
