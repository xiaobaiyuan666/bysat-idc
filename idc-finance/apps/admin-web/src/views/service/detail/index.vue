<script setup lang="ts">
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import AutomationTaskPanel from "@/components/workbench/AutomationTaskPanel.vue";
import AuditTrailTable from "@/components/workbench/AuditTrailTable.vue";
import ContextTabs from "@/components/workbench/ContextTabs.vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { useLocaleStore } from "@/store";
import {
  createServiceChangeOrder,
  fetchMofangServiceResources,
  fetchMofangSyncLogs,
  fetchProviderAccounts,
  fetchServiceDetail,
  runMofangServiceResourceAction,
  runServiceAction,
  syncMofangService,
  updateServiceRecord,
  type AuditLog,
  type MofangBackup,
  type MofangDisk,
  type MofangResourceActionResponse,
  type MofangServiceResourcesResponse,
  type MofangSnapshot,
  type MofangSyncLogItem,
  type ProviderAccount,
  type ServiceDetailResponse
} from "@/api/admin";
import {
  formatBillingCycle,
  formatChangeOrderAction,
  formatChangeOrderExecution,
  formatInvoiceStatus,
  formatLastAction,
  formatMoney,
  formatProviderType,
  formatResourceType,
  formatServiceStatus,
  formatSyncLogAction,
  formatSyncStatus,
  providerTypeOptions as createProviderTypeOptions,
  serviceStatusOptions as createServiceStatusOptions,
  syncStatusOptions as createSyncStatusOptions
} from "@/utils/business";

type ResourceDialogAction =
  | "add-ipv4"
  | "add-ipv6"
  | "add-disk"
  | "resize-disk"
  | "create-snapshot"
  | "create-backup"
  | "";

type ResourceActionName =
  | "add-ipv4"
  | "add-ipv6"
  | "add-disk"
  | "resize-disk"
  | "create-snapshot"
  | "create-backup"
  | "delete-snapshot"
  | "restore-snapshot"
  | "delete-backup"
  | "restore-backup";

type ResourceActionFollowUp = {
  action: ResourceActionName;
  title: string;
  payload: Record<string, unknown>;
  result: MofangResourceActionResponse;
  hasBillingImpact: boolean;
  periodAmount: number;
  summary: string;
};

type ResourceTimelineItem = {
  id: string;
  source: "FOLLOW_UP" | "SYNC" | "CHANGE_ORDER" | "AUDIT";
  createdAt: string;
  category: string;
  title: string;
  statusText: string;
  statusType: "success" | "warning" | "danger" | "info" | "primary";
  summary: string;
  invoiceId?: number;
  orderId?: number;
  syncLog?: MofangSyncLogItem;
  changeOrder?: ServiceDetailResponse["changeOrders"][number];
  auditLog?: AuditLog;
};

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const syncLoading = ref(false);
const actionLoading = ref("");
const resourceLoading = ref("");
const changeOrderLoading = ref(false);
const retryingSyncLogId = ref(0);
const detail = ref<ServiceDetailResponse | null>(null);
const resources = ref<MofangServiceResourcesResponse | null>(null);
const providerAccounts = ref<ProviderAccount[]>([]);
const syncLogs = ref<MofangSyncLogItem[]>([]);
const resourceActionFollowUp = ref<ResourceActionFollowUp | null>(null);

const editVisible = ref(false);
const syncReviewVisible = ref(false);
const resourceFollowUpVisible = ref(false);
const resetVisible = ref(false);
const reinstallVisible = ref(false);
const resourceVisible = ref(false);
const resourceAction = ref<ResourceDialogAction>("");
const resourceTargetId = ref("");
const reviewingSyncLog = ref<MofangSyncLogItem | null>(null);

const savingEdit = ref(false);
const savingSyncReview = ref(false);
const resetForm = reactive({ password: "" });
const reinstallForm = reactive({ imageName: "Ubuntu 24.04.1 LTS" });
const resourceForm = reactive({ count: 1, sizeGb: 10, driver: "virtio", name: "" });
const syncReviewForm = reactive({
  targetStatus: "SUCCESS",
  syncMessage: "",
  reason: ""
});
const editForm = reactive({
  providerType: "MOFANG_CLOUD",
  providerAccountId: 0,
  providerResourceId: "",
  regionName: "",
  ipAddress: "",
  nextDueAt: "",
  status: "PENDING",
  syncStatus: "PENDING",
  syncMessage: "",
  reason: ""
});

const providerTypeOptions = computed(() => createProviderTypeOptions(localeStore.locale));
const serviceStatusOptions = computed(() => createServiceStatusOptions(localeStore.locale));
const syncStatusOptions = computed(() => createSyncStatusOptions(localeStore.locale));
const _legacyProviderTypeOptions = [
  { label: "本地模块", value: "LOCAL" },
  { label: "魔方云", value: "MOFANG_CLOUD" },
  { label: "上下游财务", value: "ZJMF_API" },
  { label: "资源池", value: "RESOURCE" },
  { label: "手动资源", value: "MANUAL" }
];

const _legacyServiceStatusOptions = [
  { label: "待开通", value: "PENDING" },
  { label: "运行中", value: "ACTIVE" },
  { label: "已暂停", value: "SUSPENDED" },
  { label: "已终止", value: "TERMINATED" }
];

const _legacySyncStatusOptions = [
  { label: "待同步", value: "PENDING" },
  { label: "执行中", value: "RUNNING" },
  { label: "同步成功", value: "SUCCESS" },
  { label: "同步失败", value: "FAILED" }
];

const isMofangService = computed(() => detail.value?.service.providerType === "MOFANG_CLOUD");
const customerId = computed(() => detail.value?.service.customerId ?? 0);
const serviceAccount = computed(() => {
  const accountId = detail.value?.service.providerAccountId || 0;
  if (!accountId) return null;
  return providerAccounts.value.find(item => item.id === accountId) ?? null;
});
const serviceEditImpact = computed(() => {
  switch (editForm.status) {
    case "ACTIVE":
      return {
        type: "success" as const,
        title: "将服务改为运行中",
        description: "适合人工恢复服务或校正运行状态，订单展示也会同步回到生效。"
      };
    case "TERMINATED":
      return {
        type: "warning" as const,
        title: "将服务改为已终止",
        description: "适合人工下线或清退业务，订单展示会按账单状态回退。"
      };
    case "SUSPENDED":
      return {
        type: "info" as const,
        title: "将服务改为已暂停",
        description: "适合欠费、人工冻结或待恢复场景，账单本身不会直接作废。"
      };
    default:
      return {
        type: "info" as const,
        title: "将服务改为待开通",
        description: "适合资源尚未正式交付的场景，后续仍可继续同步或人工校正。"
      };
  }
});

const editProviderAccounts = computed(() => {
  const activeAccounts = providerAccounts.value.filter(item => item.status === "ACTIVE");
  if (editForm.providerType === "MOFANG_CLOUD") {
    return activeAccounts.filter(item => item.providerType === "MOFANG_CLOUD");
  }
  if (editForm.providerType === "ZJMF_API") {
    return activeAccounts.filter(item => item.providerType === "ZJMF_API");
  }
  if (editForm.providerType === "RESOURCE") {
    return activeAccounts.filter(item => item.providerType === "RESOURCE");
  }
  return activeAccounts;
});

const serviceEditGuide = computed(() => {
  if (editForm.providerType === "MOFANG_CLOUD") {
    return {
      type: "success" as const,
      title: "当前服务绑定魔方云渠道",
      description: "适合需要直接同步实例、磁盘、快照和 IP 的云产品。修改接口账户或远端资源 ID 后，后续资源动作会直接下发到对应魔方云账户。"
    };
  }
  if (editForm.providerType === "ZJMF_API") {
    return {
      type: "warning" as const,
      title: "当前服务绑定财务上下游渠道",
      description: "适合走上游财务链路的商品。修改后要重点核对上游产品、账单和服务编号是否一致。"
    };
  }
  if (editForm.providerType === "RESOURCE") {
    return {
      type: "info" as const,
      title: "当前服务绑定资源池渠道",
      description: "适合从库存或资源池分配现成资源。通常不需要远端自动开通，但需要维护资源编号和到期时间。"
    };
  }
  if (editForm.providerType === "MANUAL") {
    return {
      type: "info" as const,
      title: "当前服务绑定手工资源渠道",
      description: "适合代维、临时业务或人工交付场景。建议在变更原因里写清楚处理背景，方便客服和运维追踪。"
    };
  }
  return {
    type: "info" as const,
    title: "当前服务绑定本地模块",
    description: "适合完全由本地系统维护的服务对象。修改后不会自动对接外部平台，但会联动本地订单和账单状态。"
  };
});

const resourceActionGuide = computed(() => {
  switch (resourceAction.value) {
    case "add-ipv4":
      return {
        title: "新增 IPv4",
        description: "为当前实例追加公网 IPv4。系统会把新增地址同步回本地资源快照。",
        submitText: "确认新增 IPv4"
      };
    case "add-ipv6":
      return {
        title: "新增 IPv6",
        description: "为当前实例追加 IPv6 地址。适合双栈或海外网络场景。",
        submitText: "确认新增 IPv6"
      };
    case "add-disk":
      return {
        title: "新增数据盘",
        description: "为当前实例增加新的数据盘。建议先确认磁盘类型、容量和驱动。",
        submitText: "确认新增数据盘"
      };
    case "resize-disk":
      return {
        title: "扩容数据盘",
        description: "扩容会直接作用于远端磁盘。请确保填写的是扩容后的总容量，而不是增加量。",
        submitText: "确认扩容磁盘"
      };
    case "create-snapshot":
      return {
        title: "创建快照",
        description: "创建快照前建议确认业务低峰期，避免在高写入时段做保护动作。",
        submitText: "确认创建快照"
      };
    case "create-backup":
      return {
        title: "创建备份",
        description: "备份适合长期保留数据。系统会把创建结果同步到服务资源页。",
        submitText: "确认创建备份"
      };
    default:
      return {
        title: "资源动作",
        description: "资源动作会直接下发到魔方云，并记录到同步日志和自动化任务。",
        submitText: "确认执行"
      };
  }
});

const resourceDialogMeta = computed(() => {
  switch (resourceAction.value) {
    case "add-ipv4":
      return {
        countLabel: "新增 IPv4 数量",
        sizeLabel: "容量 GB",
        nameLabel: "名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "请输入要新增的 IPv4 数量",
        sizePlaceholder: "",
        namePlaceholder: ""
      };
    case "add-ipv6":
      return {
        countLabel: "新增 IPv6 数量",
        sizeLabel: "容量 GB",
        nameLabel: "名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "请输入要新增的 IPv6 数量",
        sizePlaceholder: "",
        namePlaceholder: ""
      };
    case "add-disk":
      return {
        countLabel: "数量",
        sizeLabel: "新数据盘容量 GB",
        nameLabel: "名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "",
        sizePlaceholder: "请输入新数据盘容量",
        namePlaceholder: ""
      };
    case "resize-disk":
      return {
        countLabel: "数量",
        sizeLabel: "扩容后总容量 GB",
        nameLabel: "名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "",
        sizePlaceholder: "请输入扩容后的总容量",
        namePlaceholder: ""
      };
    case "create-snapshot":
      return {
        countLabel: "数量",
        sizeLabel: "容量 GB",
        nameLabel: "快照名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "",
        sizePlaceholder: "",
        namePlaceholder: "例如：升级前快照"
      };
    case "create-backup":
      return {
        countLabel: "数量",
        sizeLabel: "容量 GB",
        nameLabel: "备份名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "",
        sizePlaceholder: "",
        namePlaceholder: "例如：周度备份"
      };
    default:
      return {
        countLabel: "数量",
        sizeLabel: "容量 GB",
        nameLabel: "名称",
        driverLabel: "磁盘驱动",
        countPlaceholder: "",
        sizePlaceholder: "",
        namePlaceholder: ""
      };
  }
});

const resourceBillingImpact = computed(() => {
  const cycle = detail.value?.invoice?.billingCycle || detail.value?.order?.billingCycle || "monthly";
  const factorMap: Record<string, number> = {
    monthly: 1,
    quarterly: 3,
    semiannually: 6,
    annually: 12,
    annual: 12,
    biennially: 24,
    triennially: 36
  };
  const factor = factorMap[String(cycle).toLowerCase()] || 1;
  const count = Number(resourceForm.count || 0);
  const sizeGb = Number(resourceForm.sizeGb || 0);
  let monthlyDelta = 0;
  let summary = "当前动作不产生建议报价。";

  switch (resourceAction.value) {
    case "add-ipv4":
      monthlyDelta = count * 35;
      summary = `建议按 ${count} 个 IPv4 生成改配报价。`;
      break;
    case "add-ipv6":
      monthlyDelta = count * 15;
      summary = `建议按 ${count} 个 IPv6 生成改配报价。`;
      break;
    case "add-disk":
      monthlyDelta = sizeGb * 0.8;
      summary = `建议按新增 ${sizeGb} GB 数据盘生成改配报价。`;
      break;
    case "resize-disk":
      monthlyDelta = sizeGb * 0.6;
      summary = `建议按扩容后的容量 ${sizeGb} GB 重新核算本期价格。`;
      break;
    case "create-snapshot":
      monthlyDelta = 20;
      summary = "建议按单次快照服务费生成改配报价。";
      break;
    case "create-backup":
      monthlyDelta = 35;
      summary = "建议按单次备份服务费生成改配报价。";
      break;
  }

  const periodAmount = Math.max(0, Number((monthlyDelta * factor).toFixed(2)));
  return {
    cycle,
    factor,
    monthlyDelta,
    periodAmount,
    summary,
    hasImpact: periodAmount > 0,
    orderNo: detail.value?.order?.orderNo || "-",
    invoiceNo: detail.value?.invoice?.invoiceNo || "-"
  };
});

const changeOrderSummary = computed(() => {
  const items = detail.value?.changeOrders ?? [];
  return {
    total: items.length,
    unpaid: items.filter(item => item.status === "UNPAID").length,
    paid: items.filter(item => item.status === "PAID").length,
    refunded: items.filter(item => item.status === "REFUNDED").length,
    waitingExecution: items.filter(item => item.executionStatus === "WAITING_PAYMENT" || item.executionStatus === "PENDING").length,
    failedExecution: items.filter(item => item.executionStatus === "FAILED").length
  };
});

const syncLogSummary = computed(() => ({
  total: syncLogs.value.length,
  success: syncLogs.value.filter(item => item.status === "SUCCESS").length,
  failed: syncLogs.value.filter(item => item.status === "FAILED").length,
  running: syncLogs.value.filter(item => item.status === "RUNNING").length,
  latestAction: syncLogs.value[0]?.action || "-"
}));

const resourceSummary = computed(() => ({
  totalIps: resources.value?.ipAddresses.length ?? 0,
  ipv4: resources.value?.ipAddresses.filter(item => item.version === "IPv4").length ?? 0,
  ipv6: resources.value?.ipAddresses.filter(item => item.version === "IPv6").length ?? 0,
  diskCount: resources.value?.disks.length ?? 0,
  diskCapacityGb: resources.value?.disks.reduce((sum, item) => sum + Number(item.sizeGb || 0), 0) ?? 0,
  snapshotCount: resources.value?.snapshots.length ?? 0,
  backupCount: resources.value?.backups.length ?? 0
}));

const resourceFollowUpSummary = computed(() => {
  const followUp = resourceActionFollowUp.value;
  if (!followUp || !detail.value) return null;
  return {
    actionTitle: followUp.title,
    message: followUp.result.message,
    remoteId: followUp.result.remoteId || detail.value.service.providerResourceId || "-",
    resourceId: followUp.result.resourceId || "-",
    syncOperation: followUp.result.syncItem?.operation || "-",
    syncMessage: followUp.result.syncItem?.message || detail.value.service.syncMessage || "-",
    hasBillingImpact: followUp.hasBillingImpact,
    periodAmount: followUp.periodAmount,
    summary: followUp.summary,
    powered: followUp.result.poweredOff
      ? "已自动关机"
      : followUp.result.poweredOn
        ? "已自动开机"
        : "未触发电源动作"
  };
});

function containsResourceKeyword(text: string) {
  const source = text.toLowerCase();
  return ["resource", "sync", "snapshot", "backup", "disk", "ip", "服务", "同步", "资源", "改配"].some(keyword =>
    source.includes(keyword)
  );
}

const resourceActionHistory = computed<ResourceTimelineItem[]>(() => {
  const rows: ResourceTimelineItem[] = [];

  if (resourceActionFollowUp.value) {
    rows.push({
      id: `followup-${resourceActionFollowUp.value.result.action}-${resourceActionFollowUp.value.result.resourceId || resourceActionFollowUp.value.result.remoteId}`,
      source: "FOLLOW_UP",
      createdAt: new Date().toISOString(),
      category: "本次动作",
      title: resourceActionFollowUp.value.title,
      statusText: "待收口",
      statusType: "primary",
      summary: resourceActionFollowUp.value.result.message
    });
  }

  for (const row of syncLogs.value) {
    rows.push({
      id: `sync-${row.id}`,
      source: "SYNC",
      createdAt: row.createdAt,
      category: "同步日志",
      title: formatSyncLogAction(localeStore.locale, row.action),
      statusText: formatSyncStatus(localeStore.locale, row.status),
      statusType: syncLogType(row.status) as "success" | "warning" | "danger" | "info" | "primary",
      summary: row.message || `${formatResourceType(localeStore.locale, row.resourceType)} / ${row.resourceId || "-"}`,
      invoiceId: row.invoiceId,
      orderId: row.orderId,
      syncLog: row
    });
  }

  for (const row of detail.value?.changeOrders ?? []) {
    rows.push({
      id: `change-${row.id}`,
      source: "CHANGE_ORDER",
      createdAt: row.createdAt,
      category: "改配单",
      title: row.title || changeOrderActionLabel(row.actionName),
      statusText: `${changeOrderStatusLabel(row.status)} / ${changeOrderExecutionLabel(row.executionStatus)}`,
      statusType:
        row.executionStatus === "FAILED" || row.executionStatus === "EXECUTE_FAILED"
          ? "danger"
          : row.status === "PAID"
            ? "success"
            : row.status === "UNPAID"
              ? "warning"
              : "info",
      summary: row.executionMessage || `账单 ${row.invoiceNo} / 订单 ${row.orderNo}`,
      invoiceId: row.invoiceId,
      orderId: row.orderId,
      changeOrder: row
    });
  }

  for (const row of detail.value?.auditLogs ?? []) {
    const payloadText = row.payload ? JSON.stringify(row.payload) : "";
    if (!containsResourceKeyword(`${row.action} ${row.description} ${payloadText}`)) continue;
    rows.push({
      id: `audit-${row.id}`,
      source: "AUDIT",
      createdAt: row.createdAt,
      category: "审计记录",
      title: row.action || "服务审计",
      statusText: row.target || "已记录",
      statusType: "info",
      summary: row.description || payloadText || "服务动作已写入审计",
      auditLog: row
    });
  }

  return rows.sort((a, b) => String(b.createdAt).localeCompare(String(a.createdAt)));
});

const resourceActionHistorySummary = computed(() => ({
  total: resourceActionHistory.value.length,
  sync: resourceActionHistory.value.filter(item => item.source === "SYNC").length,
  changeOrders: resourceActionHistory.value.filter(item => item.source === "CHANGE_ORDER").length,
  audit: resourceActionHistory.value.filter(item => item.source === "AUDIT").length,
  pending: resourceActionHistory.value.filter(item => item.statusType === "warning" || item.statusType === "primary").length,
  failed: resourceActionHistory.value.filter(item => item.statusType === "danger").length
}));

const syncLogFilter = ref<"ALL" | "FAILED" | "RUNNING" | "SUCCESS">("ALL");

const filteredSyncLogs = computed(() => {
  if (syncLogFilter.value === "ALL") return syncLogs.value;
  return syncLogs.value.filter(item => item.status === syncLogFilter.value);
});

const channelStrategyGuide = computed(() => {
  if (isMofangService.value) {
    return {
      type: "success" as const,
      title: "当前服务已接入云资源渠道",
      description:
        "适合直接处理实例生命周期、密码重置、重装系统、磁盘/IP 改配与资源同步。资源动作会联动自动化任务和同步日志。"
    };
  }
  switch (detail.value?.service.providerType) {
    case "ZJMF_API":
      return {
        type: "warning" as const,
        title: "当前服务走上游财务渠道",
        description:
          "重点是核对上游商品、账单、服务编号和同步状态。后台动作以订单、账单、工单和自动化协同为主，不建议误用云资源动作。"
      };
    case "RESOURCE":
      return {
        type: "info" as const,
        title: "当前服务走资源池交付",
        description:
          "适合从库存或资源池分配现成资源。后台要优先核对库存编号、到期时间、交付记录和人工处理说明。"
      };
    case "MANUAL":
      return {
        type: "info" as const,
        title: "当前服务走手工交付",
        description:
          "适合代维或人工处理类业务。建议通过工单、审计、改配单和账单工作台收口，不依赖远端资源动作。"
      };
    default:
      return {
        type: "info" as const,
        title: "当前服务走本地模块",
        description:
          "适合完全由本地系统维护的服务对象。重点关注本地状态、续费、改配、审计和工单协同。"
      };
  }
});

const channelCapabilitySummary = computed(() => ({
  canLifecycle: Boolean(detail.value),
  canReboot: isMofangService.value,
  canCredential: isMofangService.value,
  canSync: isMofangService.value,
  canResourceAction: isMofangService.value && Boolean(resources.value),
  canOpenChangeOrders: (detail.value?.changeOrders?.length ?? 0) > 0,
  providerLabel: providerTypeLabel(detail.value?.service.providerType || "LOCAL")
}));

const contextTabs = computed(() => [
  { key: "customer", label: "客户", to: customerId.value ? `/customer/detail/${customerId.value}` : undefined },
  { key: "service", label: "服务工作台", active: true, badge: detail.value?.service.serviceNo },
  {
    key: "invoice",
    label: "账单",
    to: detail.value?.invoice ? `/billing/invoices/${detail.value.invoice.id}` : undefined,
    badge: detail.value?.invoice?.invoiceNo
  },
  {
    key: "order",
    label: "订单",
    to: detail.value?.order ? `/orders/detail/${detail.value.order.id}` : undefined,
    badge: detail.value?.order?.orderNo
  }
]);

function renderValue(value?: string | number | null) {
  return value === undefined || value === null || value === "" ? "-" : String(value);
}

function formatBoolean(value: boolean) {
  return value ? "是" : "否";
}

function serviceStatusLabel(status: string) {
  return formatServiceStatus(localeStore.locale, status);
  const mapping: Record<string, string> = {
    PENDING: "待开通",
    ACTIVE: "运行中",
    SUSPENDED: "已暂停",
    PROVISIONING: "处理中",
    TERMINATED: "已终止",
    FAILED: "异常"
  };
  return mapping[status] ?? status;
}

function serviceStatusType(status: string) {
  const mapping: Record<string, string> = {
    PENDING: "warning",
    ACTIVE: "success",
    SUSPENDED: "info",
    PROVISIONING: "primary",
    TERMINATED: "danger",
    FAILED: "danger"
  };
  return mapping[status] ?? "info";
}

function providerTypeLabel(type: string) {
  return formatProviderType(localeStore.locale, type);
  const mapping: Record<string, string> = {
    MOFANG_CLOUD: "魔方云",
    ZJMF_API: "上下游财务",
    RESOURCE: "资源池",
    MANUAL: "手动资源",
    LOCAL: "本地模块"
  };
  return mapping[type] ?? type;
}

function changeOrderActionLabel(action: string) {
  return formatChangeOrderAction(localeStore.locale, action);
  const mapping: Record<string, string> = {
    "add-ipv4": "新增 IPv4",
    "add-ipv6": "新增 IPv6",
    "add-disk": "新增数据盘",
    "resize-disk": "扩容数据盘",
    "create-snapshot": "创建快照",
    "create-backup": "创建备份"
  };
  return mapping[action] ?? action;
}

function changeOrderStatusLabel(status: string) {
  return formatInvoiceStatus(localeStore.locale, status);
  const mapping: Record<string, string> = {
    UNPAID: "待支付",
    PAID: "已支付",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

function changeOrderExecutionLabel(status: string) {
  return formatChangeOrderExecution(localeStore.locale, status);
  const mapping: Record<string, string> = {
    WAITING_PAYMENT: "待支付",
    PAID: "待回写",
    EXECUTING: "执行中",
    EXECUTED: "已执行",
    EXECUTE_FAILED: "执行失败",
    EXECUTE_BLOCKED: "执行阻塞",
    REFUNDED: "已退款"
  };
  return mapping[status] ?? status;
}

function changeOrderExecutionType(status: string) {
  const mapping: Record<string, string> = {
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

function syncLogType(status: string) {
  const mapping: Record<string, string> = {
    SUCCESS: "success",
    FAILED: "danger",
    PENDING: "warning",
    RUNNING: "primary"
  };
  return mapping[status] ?? "info";
}

async function loadProviderAccounts() {
  providerAccounts.value = await fetchProviderAccounts();
}

async function loadDetail() {
  loading.value = true;
  try {
    detail.value = await fetchServiceDetail(route.params.id as string);
    resources.value = null;
    syncLogs.value = [];
    if (detail.value.service.providerType === "MOFANG_CLOUD") {
      const [resourceData, logData] = await Promise.all([
        fetchMofangServiceResources(route.params.id as string),
        fetchMofangSyncLogs({ serviceId: detail.value.service.id, limit: 20 })
      ]);
      resources.value = resourceData;
      syncLogs.value = logData.items;
    }
  } finally {
    loading.value = false;
  }
}

function openEditDialog() {
  if (!detail.value) return;
  editForm.providerType = detail.value.service.providerType || "LOCAL";
  editForm.providerAccountId = detail.value.service.providerAccountId || 0;
  editForm.providerResourceId = detail.value.service.providerResourceId || "";
  editForm.regionName = detail.value.service.regionName || "";
  editForm.ipAddress = detail.value.service.ipAddress || "";
  editForm.nextDueAt = detail.value.service.nextDueAt || "";
  editForm.status = detail.value.service.status || "PENDING";
  editForm.syncStatus = detail.value.service.syncStatus || "PENDING";
  editForm.syncMessage = detail.value.service.syncMessage || "";
  editForm.reason = "";
  editVisible.value = true;
}

async function submitServiceEdit() {
  if (!detail.value) return;
  savingEdit.value = true;
  try {
    await updateServiceRecord(detail.value.service.id, {
      providerType: editForm.providerType,
      providerAccountId: editForm.providerAccountId || 0,
      providerResourceId: editForm.providerResourceId,
      regionName: editForm.regionName,
      ipAddress: editForm.ipAddress,
      nextDueAt: editForm.nextDueAt,
      status: editForm.status,
      syncStatus: editForm.syncStatus,
      syncMessage: editForm.syncMessage,
      reason: editForm.reason
    });
    editVisible.value = false;
    ElMessage.success("服务信息已更新，变更记录已写入审计");
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "服务人工调整失败");
  } finally {
    savingEdit.value = false;
  }
}

async function handleSimpleAction(action: "activate" | "suspend" | "terminate" | "reboot") {
  if (!detail.value) return;
  actionLoading.value = action;
  try {
    const result = await runServiceAction(detail.value.service.id, action);
    ElMessage.success(`服务 ${result.serviceNo} 已执行 ${action}`);
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "服务动作执行失败");
  } finally {
    actionLoading.value = "";
  }
}

async function handleResetPassword() {
  if (!detail.value) return;
  actionLoading.value = "reset-password";
  try {
    await runServiceAction(detail.value.service.id, "reset-password", { password: resetForm.password.trim() });
    ElMessage.success("重置密码任务已提交");
    resetVisible.value = false;
    resetForm.password = "";
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "重置密码失败");
  } finally {
    actionLoading.value = "";
  }
}

async function handleReinstall() {
  if (!detail.value) return;
  actionLoading.value = "reinstall";
  try {
    await runServiceAction(detail.value.service.id, "reinstall", { imageName: reinstallForm.imageName.trim() });
    ElMessage.success("重装系统任务已提交");
    reinstallVisible.value = false;
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "重装系统失败");
  } finally {
    actionLoading.value = "";
  }
}

async function handleSync() {
  if (!detail.value) return;
  syncLoading.value = true;
  try {
    const result = await syncMofangService(detail.value.service.id);
    ElMessage.success(result.message);
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "同步失败");
  } finally {
    syncLoading.value = false;
  }
}

function openOrderWorkbench() {
  if (!detail.value?.order) return;
  void router.push(`/orders/detail/${detail.value.order.id}`);
}

function openCustomerWorkbench() {
  if (!detail.value) return;
  void router.push(`/customer/detail/${detail.value.service.customerId}`);
}

function openCustomerFinanceWorkbench() {
  if (!detail.value) return;
  void router.push({
    path: "/billing/accounts",
    query: {
      customerId: String(detail.value.service.customerId)
    }
  });
}

function openPaymentsWorkbench(invoiceId?: number) {
  if (!detail.value) return;
  void router.push({
    path: "/billing/payments",
    query: {
      customerId: String(detail.value.service.customerId),
      invoiceId: invoiceId ? String(invoiceId) : detail.value.invoice?.id ? String(detail.value.invoice.id) : undefined
    }
  });
}

function openRefundsWorkbench(invoiceId?: number) {
  if (!detail.value) return;
  void router.push({
    path: "/billing/refunds",
    query: {
      customerId: String(detail.value.service.customerId),
      invoiceId: invoiceId ? String(invoiceId) : detail.value.invoice?.id ? String(detail.value.invoice.id) : undefined
    }
  });
}

function openOrderRequestsWorkbench(orderId?: number) {
  const id = orderId || detail.value?.order?.id;
  if (!id) return;
  void router.push({
    path: "/orders/requests",
    query: {
      orderId: String(id)
    }
  });
}

function openInvoiceWorkbench(invoiceId?: number) {
  const id = invoiceId || detail.value?.invoice?.id;
  if (!id) return;
  void router.push(`/billing/invoices/${id}`);
}

function openTicketCreateWorkbench(overrides?: { title?: string; content?: string }) {
  if (!detail.value) return;
  void router.push({
    path: "/tickets/create",
    query: {
      customerId: String(detail.value.service.customerId),
      serviceId: String(detail.value.service.id),
      title: overrides?.title || `${detail.value.service.serviceNo} / ${detail.value.service.productName}`,
      content: overrides?.content
    }
  });
}

function openTicketListWorkbench() {
  if (!detail.value) return;
  void router.push({
    path: "/tickets/list",
    query: {
      serviceId: String(detail.value.service.id),
      customerId: String(detail.value.service.customerId)
    }
  });
}

function openProviderResourcesWorkbench() {
  if (!detail.value) return;
  void router.push({
    path: "/providers/resources",
    query: {
      providerType: detail.value.service.providerType,
      accountId: detail.value.service.providerAccountId ? String(detail.value.service.providerAccountId) : undefined,
      serviceId: String(detail.value.service.id)
    }
  });
}

function openProviderAccountsWorkbench() {
  if (!detail.value) return;
  void router.push({
    path: "/providers/accounts",
    query: {
      providerType: detail.value.service.providerType,
      accountId: detail.value.service.providerAccountId ? String(detail.value.service.providerAccountId) : undefined
    }
  });
}

function openProviderAutomationWorkbench() {
  if (!detail.value) return;
  void router.push({
    path: "/providers/automation",
    query: {
      channel: detail.value.service.providerType,
      serviceId: String(detail.value.service.id)
    }
  });
}

function openChangeOrdersWorkbench(filters?: { status?: string; executionStatus?: string }) {
  if (!detail.value) return;
  void router.push({
    path: "/orders/change-orders",
    query: {
      serviceId: String(detail.value.service.id),
      status: filters?.status || undefined,
      executionStatus: filters?.executionStatus || undefined
    }
  });
}

function openChangeOrderFailureTicket(row: {
  orderNo: string;
  invoiceNo: string;
  title: string;
  executionStatus: string;
  executionMessage?: string;
}) {
  if (!detail.value) return;
  openTicketCreateWorkbench({
    title: `${detail.value.service.serviceNo} / 改配异常 / ${row.title || row.orderNo}`,
    content: [
      `服务号：${detail.value.service.serviceNo}`,
      `订单号：${row.orderNo || "-"}`,
      `账单号：${row.invoiceNo || "-"}`,
      `执行状态：${row.executionStatus || "-"}`,
      `异常回执：${row.executionMessage || "请管理员核对自动化任务和资源状态"}`
    ].join("\n")
  });
}

function openSyncFailureTicket(row: MofangSyncLogItem) {
  if (!detail.value) return;
  openTicketCreateWorkbench({
    title: `${detail.value.service.serviceNo} / 同步异常 / ${formatSyncLogAction(localeStore.locale, row.action)}`,
    content: [
      `服务号：${detail.value.service.serviceNo}`,
      `资源类型：${formatResourceType(localeStore.locale, row.resourceType)}`,
      `资源编号：${row.resourceId || "-"}`,
      `同步状态：${formatSyncStatus(localeStore.locale, row.status)}`,
      `失败回执：${row.message || "请管理员核对自动化任务与远端资源"}`
    ].join("\n")
  });
}

function openSyncReviewDialog(row: MofangSyncLogItem) {
  if (!detail.value) return;
  reviewingSyncLog.value = row;
  syncReviewForm.targetStatus = row.status === "RUNNING" ? "RUNNING" : "SUCCESS";
  syncReviewForm.syncMessage = row.message
    ? `人工复核：${row.message}`
    : `${formatSyncLogAction(localeStore.locale, row.action)} 已人工复核`;
  syncReviewForm.reason = `围绕同步日志 ${formatSyncLogAction(localeStore.locale, row.action)} 做人工收口`;
  syncReviewVisible.value = true;
}

function buildSyntheticSyncLogFromFollowUp(): MofangSyncLogItem | null {
  if (!detail.value || !resourceActionFollowUp.value) return null;
  const followUp = resourceActionFollowUp.value;
  return {
    id: 0,
    providerType: detail.value.service.providerType,
    action: followUp.action,
    resourceType: inferResourceTypeFromAction(followUp.action),
    resourceId: followUp.result.resourceId || followUp.result.remoteId || "",
    serviceId: detail.value.service.id,
    status: detail.value.service.syncStatus || followUp.result.syncItem?.status || "PENDING",
    message: followUp.result.syncItem?.message || followUp.result.message || "",
    createdAt: new Date().toISOString()
  };
}

function openResourceFollowUpSyncReview() {
  const target = syncLogs.value[0] || buildSyntheticSyncLogFromFollowUp();
  if (!target) return;
  resourceFollowUpVisible.value = false;
  openSyncReviewDialog(target);
}

async function retrySyncLog(row: MofangSyncLogItem) {
  if (!detail.value) return;
  try {
    await ElMessageBox.confirm(
      `将重新发起当前服务的同步流程，用于重试“${formatSyncLogAction(localeStore.locale, row.action)}”。继续吗？`,
      "重新同步确认",
      {
        type: "warning",
        confirmButtonText: "重新同步",
        cancelButtonText: "取消"
      }
    );
  } catch {
    return;
  }

  retryingSyncLogId.value = row.id;
  try {
    const result = await syncMofangService(detail.value.service.id);
    ElMessage.success(result.message || "同步任务已重新发起");
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "同步重试失败");
  } finally {
    retryingSyncLogId.value = 0;
  }
}

async function submitSyncReview() {
  if (!detail.value || !reviewingSyncLog.value) return;
  savingSyncReview.value = true;
  try {
    await updateServiceRecord(detail.value.service.id, {
      syncStatus: syncReviewForm.targetStatus,
      syncMessage: syncReviewForm.syncMessage.trim(),
      reason:
        syncReviewForm.reason.trim() ||
        `围绕同步日志 ${formatSyncLogAction(localeStore.locale, reviewingSyncLog.value.action)} 做人工收口`
    });
    syncReviewVisible.value = false;
    ElMessage.success("同步状态已人工收口");
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "同步人工收口失败");
  } finally {
    savingSyncReview.value = false;
  }
}

function openProviderAutomationContext(params?: { orderId?: number; invoiceId?: number }) {
  if (!detail.value) return;
  void router.push({
    path: "/providers/automation",
    query: {
      channel: detail.value.service.providerType,
      serviceId: String(detail.value.service.id),
      orderId: params?.orderId ? String(params.orderId) : undefined,
      invoiceId: params?.invoiceId ? String(params.invoiceId) : undefined
    }
  });
}

function inferResourceTypeFromAction(action: ResourceActionName) {
  switch (action) {
    case "add-ipv4":
    case "add-ipv6":
      return "IP";
    case "add-disk":
    case "resize-disk":
      return "DISK";
    case "create-snapshot":
    case "delete-snapshot":
    case "restore-snapshot":
      return "SNAPSHOT";
    case "create-backup":
    case "delete-backup":
    case "restore-backup":
      return "BACKUP";
    default:
      return "INSTANCE";
  }
}

function openResourceDialog(action: ResourceDialogAction, payload?: { id?: string; sizeGb?: number }) {
  resourceAction.value = action;
  resourceTargetId.value = payload?.id ?? "";
  resourceForm.count = 1;
  resourceForm.sizeGb = payload?.sizeGb ?? 10;
  resourceForm.driver = "virtio";
  resourceForm.name = "";
  resourceVisible.value = true;
}

function closeResourceDialog() {
  resourceVisible.value = false;
  resourceAction.value = "";
  resourceTargetId.value = "";
}

function createResourceActionSnapshot(action: ResourceActionName, title: string, payload: Record<string, unknown>) {
  return {
    action,
    title,
    payload,
    hasBillingImpact: resourceBillingImpact.value.hasImpact,
    periodAmount: resourceBillingImpact.value.periodAmount,
    summary: resourceBillingImpact.value.summary
  };
}

async function finalizeResourceAction(
  result: MofangResourceActionResponse,
  snapshot: ReturnType<typeof createResourceActionSnapshot>,
  closeDialog = true
) {
  if (closeDialog) closeResourceDialog();
  await loadDetail();
  resourceActionFollowUp.value = {
    ...snapshot,
    result
  };
  resourceFollowUpVisible.value = true;
}

function buildResourcePayload() {
  switch (resourceAction.value) {
    case "add-ipv4":
    case "add-ipv6":
      return { count: resourceForm.count };
    case "add-disk":
      return { sizeGb: resourceForm.sizeGb, driver: resourceForm.driver };
    case "resize-disk":
      return { diskId: resourceTargetId.value, sizeGb: resourceForm.sizeGb };
    case "create-snapshot":
    case "create-backup":
      return { diskId: resourceTargetId.value, name: resourceForm.name.trim() };
    default:
      return {};
  }
}

async function submitResourceDialog() {
  if (!detail.value || !resourceAction.value) return;
  const payload = buildResourcePayload();
  const snapshot = createResourceActionSnapshot(resourceAction.value, resourceActionGuide.value.title, payload);
  resourceLoading.value = resourceAction.value;
  try {
    const result = await runMofangServiceResourceAction(detail.value.service.id, resourceAction.value, payload);
    ElMessage.success(result.message);
    await finalizeResourceAction(result, snapshot);
  } catch (error: any) {
    ElMessage.error(error?.message ?? "资源动作执行失败");
  } finally {
    resourceLoading.value = "";
  }
}

async function createChangeOrderForResourceAction(snapshot: {
  action: string;
  title: string;
  payload: Record<string, unknown>;
  hasBillingImpact: boolean;
  periodAmount: number;
}) {
  if (!detail.value || !snapshot.hasBillingImpact) return;
  changeOrderLoading.value = true;
  try {
    const result = await createServiceChangeOrder(detail.value.service.id, {
      actionName: snapshot.action,
      title: `${snapshot.title} 改配单`,
      billingCycle: detail.value.invoice?.billingCycle || detail.value.order?.billingCycle || resourceBillingImpact.value.cycle || "monthly",
      amount: Number(snapshot.periodAmount.toFixed(2)),
      reason: `服务改配：${snapshot.title}`,
      payload: snapshot.payload
    });
    ElMessage.success(`已生成改配单：${result.invoice.invoiceNo}`);
    resourceFollowUpVisible.value = false;
    closeResourceDialog();
    await router.push(`/billing/invoices/${result.invoice.id}`);
  } catch (error: any) {
    ElMessage.error(error?.message ?? "生成改配单失败");
  } finally {
    changeOrderLoading.value = false;
  }
}

async function handleCreateServiceChangeOrder() {
  if (!detail.value || !resourceAction.value || !resourceBillingImpact.value.hasImpact) return;
  await createChangeOrderForResourceAction(
    createResourceActionSnapshot(resourceAction.value, resourceActionGuide.value.title, buildResourcePayload())
  );
}

async function handleCreateServiceChangeOrderFromFollowUp() {
  if (!resourceActionFollowUp.value) return;
  await createChangeOrderForResourceAction(resourceActionFollowUp.value);
}

async function confirmResourceAction(
  action: "delete-snapshot" | "restore-snapshot" | "delete-backup" | "restore-backup",
  payload: { snapshotId?: string; backupId?: string },
  label: string
) {
  if (!detail.value) return;
  try {
    await ElMessageBox.confirm(`确认执行：${label}`, "资源动作确认", { type: "warning" });
  } catch {
    return;
  }

  resourceLoading.value = action;
  try {
    const snapshot = {
      action,
      title: label,
      payload,
      hasBillingImpact: false,
      periodAmount: 0,
      summary: "当前动作不产生改配费用，建议重点核对自动化任务和资源同步状态。"
    };
    const result = await runMofangServiceResourceAction(detail.value.service.id, action, payload);
    ElMessage.success(result.message);
    await finalizeResourceAction(result, snapshot, false);
  } catch (error: any) {
    ElMessage.error(error?.message ?? "资源动作执行失败");
  } finally {
    resourceLoading.value = "";
  }
}

watch(
  () => editForm.providerType,
  providerType => {
    if (!providerType) return;
    const availableIds = new Set(editProviderAccounts.value.map(item => item.id));
    if (!availableIds.has(editForm.providerAccountId)) {
      editForm.providerAccountId = editProviderAccounts.value[0]?.id || 0;
    }
  }
);

watch(() => route.params.id, () => void loadDetail());

onMounted(async () => {
  await loadProviderAccounts();
  await loadDetail();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="业务 / 服务"
      title="服务工作台"
      subtitle="集中处理服务对象、接口账户、资源动作、同步日志和自动化任务。"
    >
      <template #actions>
        <el-button @click="router.push('/services/list')">返回列表</el-button>
        <el-button @click="loadDetail">刷新详情</el-button>
        <el-button v-if="detail" type="primary" plain @click="openEditDialog">人工调整服务</el-button>
        <el-button v-if="isMofangService" type="primary" plain :loading="syncLoading" @click="handleSync">
          拉取信息
        </el-button>
      </template>

      <template #context>
        <ContextTabs v-if="detail" :items="contextTabs" />
      </template>

      <template #metrics>
        <div v-if="detail" class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">服务编号</span>
            <strong class="detail-kpi-card__value">{{ detail.service.serviceNo }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">状态</span>
            <el-tag :type="serviceStatusType(detail.service.status)" effect="light">
              {{ serviceStatusLabel(detail.service.status) }}
            </el-tag>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">远端实例</span>
            <strong class="detail-kpi-card__value">{{ renderValue(detail.service.providerResourceId) }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">公网 IP</span>
            <strong class="detail-kpi-card__value">{{ renderValue(detail.service.ipAddress) }}</strong>
          </div>
        </div>
      </template>

      <template v-if="detail">
        <el-tabs>
          <el-tab-pane label="服务概览">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>基础信息</strong>
                  <span class="section-card__meta">{{ providerTypeLabel(detail.service.providerType) }}</span>
                </div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="客户">{{ renderValue(detail.order?.customerName) }}</el-descriptions-item>
                  <el-descriptions-item label="产品">{{ detail.service.productName }}</el-descriptions-item>
                  <el-descriptions-item label="接口账户">
                    {{ serviceAccount ? `${serviceAccount.name} / ${serviceAccount.baseUrl}` : renderValue(detail.service.providerAccountId) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="渠道">{{ providerTypeLabel(detail.service.providerType) }}</el-descriptions-item>
                  <el-descriptions-item label="区域">
                    {{ renderValue(detail.service.regionName || detail.service.resourceSnapshot.regionName) }}
                  </el-descriptions-item>
                  <el-descriptions-item label="最近动作">{{ formatLastAction(localeStore.locale, detail.service.lastAction || "") }}</el-descriptions-item>
                  <el-descriptions-item label="同步状态">{{ formatSyncStatus(localeStore.locale, detail.service.syncStatus || "") }}</el-descriptions-item>
                  <el-descriptions-item label="同步时间">{{ renderValue(detail.service.lastSyncAt) }}</el-descriptions-item>
                  <el-descriptions-item label="下次到期">{{ renderValue(detail.service.nextDueAt) }}</el-descriptions-item>
                  <el-descriptions-item label="账单">{{ renderValue(detail.invoice?.invoiceNo) }}</el-descriptions-item>
                </el-descriptions>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>{{ isMofangService ? "实例动作" : "服务动作" }}</strong>
                  <span class="section-card__meta">
                    {{ isMofangService ? "电源、密码、重装和同步" : "状态维护、人工校正和业务协同" }}
                  </span>
                </div>
                <div class="command-group__actions">
                  <el-button type="success" :loading="actionLoading === 'activate'" @click="handleSimpleAction('activate')">
                    恢复运行
                  </el-button>
                  <el-button type="warning" :loading="actionLoading === 'suspend'" @click="handleSimpleAction('suspend')">
                    暂停服务
                  </el-button>
                  <el-button type="danger" :loading="actionLoading === 'terminate'" @click="handleSimpleAction('terminate')">
                    终止服务
                  </el-button>
                  <el-button
                    v-if="isMofangService"
                    type="primary"
                    plain
                    :loading="actionLoading === 'reboot'"
                    @click="handleSimpleAction('reboot')"
                  >
                    重启实例
                  </el-button>
                  <el-button v-if="isMofangService" type="primary" plain @click="resetVisible = true">重置密码</el-button>
                  <el-button v-if="isMofangService" type="primary" plain @click="reinstallVisible = true">重装系统</el-button>
                  <el-button plain @click="openTicketCreateWorkbench">
                    创建工单
                  </el-button>
                </div>
                <el-alert
                  style="margin-top: 16px"
                  type="info"
                  :closable="false"
                  show-icon
                  :title="
                    isMofangService
                      ? '服务动作会同时联动本地状态、自动化任务和渠道资源同步记录。'
                      : '当前服务以本地状态、账单、工单和自动化协同为主。'
                  "
                />
              </div>
            </div>

            <div class="portal-grid portal-grid--two" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>业务联动</strong>
                  <span class="section-card__meta">围绕订单、账单和工单继续处理当前服务</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill"><span>客户</span><strong>{{ detail.order?.customerName || `客户 #${detail.service.customerId}` }}</strong></div>
                  <div class="summary-pill"><span>订单</span><strong>{{ detail.order?.orderNo || "-" }}</strong></div>
                  <div class="summary-pill"><span>账单</span><strong>{{ detail.invoice?.invoiceNo || "-" }}</strong></div>
                  <div class="summary-pill"><span>服务</span><strong>{{ detail.service.serviceNo }}</strong></div>
                </div>
                <div class="command-group__actions" style="margin-top: 16px">
                  <el-button plain @click="openCustomerWorkbench">客户详情</el-button>
                  <el-button plain :disabled="!detail.order" @click="openOrderWorkbench">订单工作台</el-button>
                  <el-button type="primary" plain :disabled="!detail.invoice" @click="openInvoiceWorkbench">账单工作台</el-button>
                  <el-button plain @click="openCustomerFinanceWorkbench">资金台账</el-button>
                  <el-button plain @click="openTicketListWorkbench">工单中心</el-button>
                  <el-button plain @click="openTicketCreateWorkbench">代客建单</el-button>
                </div>
                <div class="command-group__actions" style="margin-top: 12px">
                  <el-button plain :disabled="!detail.order" @click="openOrderRequestsWorkbench">订单申请</el-button>
                  <el-button plain @click="openPaymentsWorkbench">支付记录</el-button>
                  <el-button plain @click="openRefundsWorkbench">退款记录</el-button>
                </div>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>渠道联动</strong>
                  <span class="section-card__meta">从服务直接进入渠道资源、自动化和接口账户工作台</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill"><span>渠道</span><strong>{{ providerTypeLabel(detail.service.providerType) }}</strong></div>
                  <div class="summary-pill"><span>账户</span><strong>{{ detail.service.providerAccountId || "-" }}</strong></div>
                  <div class="summary-pill"><span>远端资源</span><strong>{{ detail.service.providerResourceId || "-" }}</strong></div>
                </div>
                <div class="command-group__actions" style="margin-top: 16px">
                  <el-button type="primary" plain @click="openProviderResourcesWorkbench">渠道资源</el-button>
                  <el-button plain @click="openProviderAutomationWorkbench">自动化任务</el-button>
                  <el-button plain @click="openProviderAccountsWorkbench">接口账户</el-button>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="渠道动作中心">
            <div class="portal-grid portal-grid--three">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>渠道策略</strong>
                  <span class="section-card__meta">{{ channelCapabilitySummary.providerLabel }}</span>
                </div>
                <el-alert
                  :type="channelStrategyGuide.type"
                  :title="channelStrategyGuide.title"
                  :description="channelStrategyGuide.description"
                  :closable="false"
                  show-icon
                />
                <div class="summary-strip" style="margin-top: 16px">
                  <div class="summary-pill"><span>接口账户</span><strong>{{ serviceAccount?.name || detail.service.providerAccountId || "-" }}</strong></div>
                  <div class="summary-pill"><span>远端资源</span><strong>{{ detail.service.providerResourceId || "-" }}</strong></div>
                  <div class="summary-pill"><span>同步状态</span><strong>{{ formatSyncStatus(localeStore.locale, detail.service.syncStatus || "") }}</strong></div>
                </div>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>生命周期动作</strong>
                  <span class="section-card__meta">状态变更与实例控制</span>
                </div>
                <div class="command-group__actions">
                  <el-button type="success" :loading="actionLoading === 'activate'" @click="handleSimpleAction('activate')">
                    恢复运行
                  </el-button>
                  <el-button type="warning" :loading="actionLoading === 'suspend'" @click="handleSimpleAction('suspend')">
                    暂停服务
                  </el-button>
                  <el-button type="danger" :loading="actionLoading === 'terminate'" @click="handleSimpleAction('terminate')">
                    终止服务
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="!channelCapabilitySummary.canReboot"
                    :loading="actionLoading === 'reboot'"
                    @click="handleSimpleAction('reboot')"
                  >
                    重启实例
                  </el-button>
                </div>
                <div class="section-card__meta" style="margin-top: 16px">
                  {{
                    channelCapabilitySummary.canReboot
                      ? "当前渠道支持实例级生命周期动作，执行后会联动自动化任务和同步状态。"
                      : "当前渠道以本地状态维护为主，不提供实例级重启能力。"
                  }}
                </div>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>系统与同步</strong>
                  <span class="section-card__meta">凭据、镜像、资源回拉</span>
                </div>
                <div class="command-group__actions">
                  <el-button type="primary" plain :disabled="!channelCapabilitySummary.canCredential" @click="resetVisible = true">
                    重置密码
                  </el-button>
                  <el-button type="primary" plain :disabled="!channelCapabilitySummary.canCredential" @click="reinstallVisible = true">
                    重装系统
                  </el-button>
                  <el-button type="primary" plain :disabled="!channelCapabilitySummary.canSync" :loading="syncLoading" @click="handleSync">
                    拉取信息
                  </el-button>
                  <el-button plain @click="openProviderAutomationWorkbench">查看自动化</el-button>
                </div>
                <div class="section-card__meta" style="margin-top: 16px">
                  {{
                    channelCapabilitySummary.canSync
                      ? "当前渠道支持从上游回拉资源与状态信息。"
                      : "当前渠道以人工维护和本地协同为主，建议通过工单、订单和审计记录处理。"
                  }}
                </div>
              </div>
            </div>

            <div class="portal-grid portal-grid--two" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>改配与资源收口</strong>
                  <span class="section-card__meta">资源动作、改配单、账单收口</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill"><span>改配总数</span><strong>{{ changeOrderSummary.total }}</strong></div>
                  <div class="summary-pill"><span>待支付</span><strong>{{ changeOrderSummary.unpaid }}</strong></div>
                  <div class="summary-pill"><span>待执行</span><strong>{{ changeOrderSummary.waitingExecution }}</strong></div>
                  <div class="summary-pill"><span>执行异常</span><strong>{{ changeOrderSummary.failedExecution }}</strong></div>
                </div>
                <el-alert
                  style="margin-top: 12px"
                  :type="channelCapabilitySummary.canResourceAction ? 'info' : 'warning'"
                  :closable="false"
                  show-icon
                  :title="
                    channelCapabilitySummary.canResourceAction
                      ? '当前渠道支持资源级扩容与数据保护动作，建议改配后立即回到账单和改配单工作台收口。'
                      : '当前渠道不提供资源级动作，建议通过改配单、账单、工单和人工调整完成业务收口。'
                  "
                />
                <div class="command-group__actions" style="margin-top: 16px">
                  <el-button type="primary" plain :disabled="!channelCapabilitySummary.canResourceAction" @click="openResourceDialog('add-ipv4')">
                    新增 IPv4
                  </el-button>
                  <el-button type="primary" plain :disabled="!channelCapabilitySummary.canResourceAction" @click="openResourceDialog('add-disk')">
                    新增数据盘
                  </el-button>
                  <el-button plain @click="openChangeOrdersWorkbench()">改配单工作台</el-button>
                  <el-button
                    plain
                    :disabled="changeOrderSummary.unpaid === 0"
                    @click="openChangeOrdersWorkbench({ status: 'UNPAID' })"
                  >
                    待支付改配
                  </el-button>
                  <el-button
                    plain
                    :disabled="changeOrderSummary.failedExecution === 0"
                    @click="openChangeOrdersWorkbench({ executionStatus: 'EXECUTE_FAILED' })"
                  >
                    执行异常
                  </el-button>
                  <el-button type="primary" plain :disabled="!detail.invoice" @click="openInvoiceWorkbench">账单工作台</el-button>
                </div>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>渠道联查入口</strong>
                  <span class="section-card__meta">接口账户、资源、任务与工单协同</span>
                </div>
                <div class="summary-strip">
                  <div class="summary-pill"><span>工单入口</span><strong>已接通</strong></div>
                  <div class="summary-pill"><span>自动化</span><strong>{{ detail.service.providerType || "-" }}</strong></div>
                  <div class="summary-pill"><span>同步日志</span><strong>{{ syncLogSummary.total }}</strong></div>
                  <div class="summary-pill"><span>失败日志</span><strong>{{ syncLogSummary.failed }}</strong></div>
                </div>
                <div class="command-group__actions" style="margin-top: 16px">
                  <el-button plain @click="openProviderAccountsWorkbench">接口账户</el-button>
                  <el-button type="primary" plain @click="openProviderResourcesWorkbench">渠道资源</el-button>
                  <el-button plain @click="openProviderAutomationWorkbench">自动化任务</el-button>
                  <el-button plain @click="openTicketListWorkbench">工单中心</el-button>
                </div>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="资源快照">
            <div class="portal-grid portal-grid--two">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>资源规格</strong>
                </div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="地域">{{ renderValue(detail.service.resourceSnapshot.regionName) }}</el-descriptions-item>
                  <el-descriptions-item label="可用区">{{ renderValue(detail.service.resourceSnapshot.zoneName) }}</el-descriptions-item>
                  <el-descriptions-item label="主机名">{{ renderValue(detail.service.resourceSnapshot.hostname) }}</el-descriptions-item>
                  <el-descriptions-item label="系统镜像">{{ renderValue(detail.service.resourceSnapshot.operatingSystem) }}</el-descriptions-item>
                  <el-descriptions-item label="登录用户">{{ renderValue(detail.service.resourceSnapshot.loginUsername) }}</el-descriptions-item>
                  <el-descriptions-item label="密码提示">{{ renderValue(detail.service.resourceSnapshot.passwordHint) }}</el-descriptions-item>
                  <el-descriptions-item label="安全组">{{ renderValue(detail.service.resourceSnapshot.securityGroup) }}</el-descriptions-item>
                  <el-descriptions-item label="公网 IPv4">{{ renderValue(detail.service.resourceSnapshot.publicIpv4) }}</el-descriptions-item>
                  <el-descriptions-item label="公网 IPv6">{{ renderValue(detail.service.resourceSnapshot.publicIpv6) }}</el-descriptions-item>
                  <el-descriptions-item label="CPU">{{ renderValue(detail.service.resourceSnapshot.cpuCores) }} 核</el-descriptions-item>
                  <el-descriptions-item label="内存">{{ renderValue(detail.service.resourceSnapshot.memoryGB) }} GB</el-descriptions-item>
                  <el-descriptions-item label="带宽">{{ renderValue(detail.service.resourceSnapshot.bandwidthMbps) }} Mbps</el-descriptions-item>
                  <el-descriptions-item label="系统盘">{{ renderValue(detail.service.resourceSnapshot.systemDiskGB) }} GB</el-descriptions-item>
                  <el-descriptions-item label="数据盘">{{ renderValue(detail.service.resourceSnapshot.dataDiskGB) }} GB</el-descriptions-item>
                </el-descriptions>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>配置项快照</strong>
                  <span class="section-card__meta">来自下单与同步</span>
                </div>
                <el-table :data="detail.service.configuration" border stripe empty-text="当前没有配置项快照">
                  <el-table-column prop="name" label="配置项" min-width="180" />
                  <el-table-column prop="valueLabel" label="当前值" min-width="180" />
                  <el-table-column label="加价" min-width="120">
                    <template #default="{ row }">¥{{ Number(row.priceDelta || 0).toFixed(2) }}</template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane v-if="resources" label="渠道资源动作">
            <div class="summary-strip">
              <div class="summary-pill"><span>IP</span><strong>{{ resourceSummary.totalIps }}</strong></div>
              <div class="summary-pill"><span>IPv4</span><strong>{{ resourceSummary.ipv4 }}</strong></div>
              <div class="summary-pill"><span>IPv6</span><strong>{{ resourceSummary.ipv6 }}</strong></div>
              <div class="summary-pill"><span>磁盘</span><strong>{{ resourceSummary.diskCount }}</strong></div>
              <div class="summary-pill"><span>磁盘容量</span><strong>{{ resourceSummary.diskCapacityGb }} GB</strong></div>
              <div class="summary-pill"><span>快照</span><strong>{{ resourceSummary.snapshotCount }}</strong></div>
              <div class="summary-pill"><span>备份</span><strong>{{ resourceSummary.backupCount }}</strong></div>
            </div>

            <el-alert
              style="margin-top: 16px"
              type="info"
              :closable="false"
              show-icon
              title="这里的操作会直接调用当前云资源渠道接口，并把结果写入同步日志和自动化任务。"
              description="建议按“网络扩容、存储扩容、数据保护”顺序操作。执行前先核对当前渠道、实例编号、磁盘列表和公网地址。"
            />

            <div class="portal-grid portal-grid--three" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>网络扩容</strong>
                  <span class="section-card__meta">增加公网地址</span>
                </div>
                <div class="section-card__meta">适合客户额外购买 IPv4 / IPv6，执行后会自动刷新本地 IP 资源记录。</div>
                <div class="command-group__actions" style="margin-top: 16px">
                  <el-button size="small" type="primary" plain @click="openResourceDialog('add-ipv4')">新增 IPv4</el-button>
                  <el-button size="small" type="primary" plain @click="openResourceDialog('add-ipv6')">新增 IPv6</el-button>
                </div>
              </div>
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>存储扩容</strong>
                  <span class="section-card__meta">新增或扩容磁盘</span>
                </div>
                <div class="section-card__meta">适合订单升级、容量扩充。新增磁盘后，可继续在下方为磁盘创建快照或备份。</div>
                <div class="command-group__actions" style="margin-top: 16px">
                  <el-button size="small" type="primary" plain @click="openResourceDialog('add-disk')">新增数据盘</el-button>
                </div>
              </div>
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>数据保护</strong>
                  <span class="section-card__meta">快照与备份</span>
                </div>
                <div class="section-card__meta">快照适合短期回滚，备份适合长期保留。具体创建入口在下方磁盘列表中选择。</div>
                <div class="summary-strip" style="margin-top: 16px">
                  <div class="summary-pill"><span>可恢复快照</span><strong>{{ resources.snapshots.length }}</strong></div>
                  <div class="summary-pill"><span>可恢复备份</span><strong>{{ resources.backups.length }}</strong></div>
                </div>
              </div>
            </div>

            <div class="panel-card" style="margin-top: 16px">
              <div class="section-card__head">
                <strong>改配后的财务收口</strong>
                <span class="section-card__meta">方便客服和财务继续处理</span>
              </div>
              <div class="summary-strip">
                <div class="summary-pill"><span>关联订单</span><strong>{{ detail.order?.orderNo || "-" }}</strong></div>
                <div class="summary-pill"><span>关联账单</span><strong>{{ detail.invoice?.invoiceNo || "-" }}</strong></div>
                <div class="summary-pill">
                  <span>当前周期</span>
                  <strong>{{ detail.invoice?.billingCycle || detail.order?.billingCycle ? formatBillingCycle(localeStore.locale, detail.invoice?.billingCycle || detail.order?.billingCycle || "monthly") : "-" }}</strong>
                </div>
              </div>
              <el-alert
                style="margin-top: 12px"
                type="info"
                :closable="false"
                show-icon
                title="资源动作执行成功后，建议立即核对改配费用并回到订单/账单工作台收口。"
                description="如果是客户主动升级，优先处理账单；如果是后台补配，建议在变更原因里写清楚，并同步更新订单和账单金额。"
              />
              <div class="command-group__actions" style="margin-top: 16px">
                <el-button v-if="detail.order" plain @click="router.push(`/orders/detail/${detail.order.id}`)">查看订单工作台</el-button>
                <el-button v-if="detail.invoice" type="primary" plain @click="router.push(`/billing/invoices/${detail.invoice.id}`)">查看账单工作台</el-button>
              </div>
            </div>

            <div class="portal-grid portal-grid--two" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>IP 地址</strong>
                  <div class="inline-actions">
                    <el-button size="small" type="primary" plain @click="openResourceDialog('add-ipv4')">新增 IPv4</el-button>
                    <el-button size="small" type="primary" plain @click="openResourceDialog('add-ipv6')">新增 IPv6</el-button>
                  </div>
                </div>
                <el-table :data="resources.ipAddresses" border stripe empty-text="暂无 IP 地址">
                  <el-table-column prop="address" label="地址" min-width="170" />
                  <el-table-column prop="version" label="版本" min-width="90" />
                  <el-table-column prop="gateway" label="网关" min-width="130" />
                  <el-table-column prop="bandwidthMbps" label="带宽" min-width="90" />
                  <el-table-column label="主 IP" min-width="90">
                    <template #default="{ row }">{{ formatBoolean(row.isPrimary) }}</template>
                  </el-table-column>
                </el-table>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>磁盘资源</strong>
                  <div class="inline-actions">
                    <el-button size="small" type="primary" plain @click="openResourceDialog('add-disk')">新增数据盘</el-button>
                  </div>
                </div>
                <el-table :data="resources.disks" border stripe empty-text="暂无磁盘资源">
                  <el-table-column prop="name" label="名称" min-width="160" />
                  <el-table-column prop="diskType" label="类型" min-width="100" />
                  <el-table-column prop="sizeGb" label="容量(GB)" min-width="100" />
                  <el-table-column prop="deviceName" label="设备名" min-width="120" />
                  <el-table-column label="系统盘" min-width="90">
                    <template #default="{ row }">{{ formatBoolean(row.isSystem) }}</template>
                  </el-table-column>
                  <el-table-column label="操作" min-width="220" fixed="right">
                    <template #default="{ row }">
                      <div class="inline-actions">
                        <el-button
                          size="small"
                          link
                          type="primary"
                          @click="openResourceDialog('create-snapshot', { id: (row as MofangDisk).providerDiskId })"
                        >
                          创建快照
                        </el-button>
                        <el-button
                          size="small"
                          link
                          type="primary"
                          @click="openResourceDialog('create-backup', { id: (row as MofangDisk).providerDiskId })"
                        >
                          创建备份
                        </el-button>
                        <el-button
                          v-if="!(row as MofangDisk).isSystem"
                          size="small"
                          link
                          type="warning"
                          @click="openResourceDialog('resize-disk', { id: (row as MofangDisk).providerDiskId, sizeGb: (row as MofangDisk).sizeGb + 10 })"
                        >
                          扩容
                        </el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </div>

            <div class="portal-grid portal-grid--two" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>快照</strong>
                  <span class="section-card__meta">支持恢复和删除</span>
                </div>
                <el-table :data="resources.snapshots" border stripe empty-text="暂无快照">
                  <el-table-column prop="name" label="名称" min-width="170" />
                  <el-table-column prop="sizeGb" label="容量(GB)" min-width="100" />
                  <el-table-column prop="status" label="状态" min-width="100" />
                  <el-table-column label="操作" min-width="180">
                    <template #default="{ row }">
                      <div class="inline-actions">
                        <el-button
                          size="small"
                          link
                          type="primary"
                          @click="confirmResourceAction('restore-snapshot', { snapshotId: (row as MofangSnapshot).providerSnapshotId }, `恢复快照 ${(row as MofangSnapshot).name}`)"
                        >
                          恢复
                        </el-button>
                        <el-button
                          size="small"
                          link
                          type="danger"
                          @click="confirmResourceAction('delete-snapshot', { snapshotId: (row as MofangSnapshot).providerSnapshotId }, `删除快照 ${(row as MofangSnapshot).name}`)"
                        >
                          删除
                        </el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>备份</strong>
                  <span class="section-card__meta">支持恢复和删除</span>
                </div>
                <el-table :data="resources.backups" border stripe empty-text="暂无备份">
                  <el-table-column prop="name" label="名称" min-width="170" />
                  <el-table-column prop="sizeGb" label="容量(GB)" min-width="100" />
                  <el-table-column prop="status" label="状态" min-width="100" />
                  <el-table-column prop="expiresAt" label="过期时间" min-width="160" />
                  <el-table-column label="操作" min-width="180">
                    <template #default="{ row }">
                      <div class="inline-actions">
                        <el-button
                          size="small"
                          link
                          type="primary"
                          @click="confirmResourceAction('restore-backup', { backupId: (row as MofangBackup).providerBackupId }, `恢复备份 ${(row as MofangBackup).name}`)"
                        >
                          恢复
                        </el-button>
                        <el-button
                          size="small"
                          link
                          type="danger"
                          @click="confirmResourceAction('delete-backup', { backupId: (row as MofangBackup).providerBackupId }, `删除备份 ${(row as MofangBackup).name}`)"
                        >
                          删除
                        </el-button>
                      </div>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane label="同步日志">
            <div class="summary-strip" style="margin-bottom: 16px">
              <div class="summary-pill">
                <span>日志总数</span>
                <strong>{{ syncLogSummary.total }}</strong>
              </div>
              <div class="summary-pill">
                <span>成功</span>
                <strong>{{ syncLogSummary.success }}</strong>
              </div>
              <div class="summary-pill">
                <span>失败</span>
                <strong>{{ syncLogSummary.failed }}</strong>
              </div>
              <div class="summary-pill">
                <span>执行中</span>
                <strong>{{ syncLogSummary.running }}</strong>
              </div>
              <div class="summary-pill">
                <span>最近动作</span>
                <strong>{{ formatSyncLogAction(localeStore.locale, syncLogSummary.latestAction) }}</strong>
              </div>
            </div>
            <div class="table-toolbar">
              <div class="table-toolbar__meta">
                <strong>同步协同工作台</strong>
                <span>围绕当前服务的同步回执、资源状态和财务收口继续处理</span>
              </div>
              <div class="inline-actions">
                <el-button plain :loading="syncLoading" @click="handleSync">重新同步服务</el-button>
                <el-button plain @click="syncLogFilter = 'ALL'">全部日志</el-button>
                <el-button plain @click="syncLogFilter = 'FAILED'">失败日志</el-button>
                <el-button plain @click="syncLogFilter = 'RUNNING'">执行中</el-button>
                <el-button plain @click="syncLogFilter = 'SUCCESS'">成功日志</el-button>
                <el-button plain @click="openCustomerFinanceWorkbench">客户财务</el-button>
                <el-button plain @click="openTicketListWorkbench">工单中心</el-button>
                <el-button plain @click="openProviderResourcesWorkbench">渠道资源</el-button>
                <el-button plain @click="openProviderAutomationWorkbench">自动化任务</el-button>
              </div>
            </div>
            <el-table :data="filteredSyncLogs" border stripe empty-text="当前服务暂无同步日志">
              <el-table-column prop="createdAt" label="时间" min-width="180" />
              <el-table-column label="动作" min-width="140">
                <template #default="{ row }">{{ formatSyncLogAction(localeStore.locale, row.action) }}</template>
              </el-table-column>
              <el-table-column label="资源类型" min-width="140">
                <template #default="{ row }">{{ formatResourceType(localeStore.locale, row.resourceType) }}</template>
              </el-table-column>
              <el-table-column prop="resourceId" label="资源编号" min-width="140" />
              <el-table-column label="状态" min-width="100">
                <template #default="{ row }">
                  <el-tag :type="syncLogType(row.status)" effect="light">{{ formatSyncStatus(localeStore.locale, row.status) }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="message" label="回执" min-width="320" show-overflow-tooltip />
              <el-table-column label="操作" min-width="320" fixed="right">
                <template #default="{ row }">
                  <div class="inline-actions">
                    <el-button v-if="row.invoiceId || detail?.invoice" type="primary" link @click="openInvoiceWorkbench(row.invoiceId)">
                      账单
                    </el-button>
                    <el-button v-if="row.invoiceId || detail?.invoice" type="primary" link @click="openPaymentsWorkbench(row.invoiceId)">
                      支付
                    </el-button>
                    <el-button v-if="row.invoiceId || detail?.invoice" type="primary" link @click="openRefundsWorkbench(row.invoiceId)">
                      退款
                    </el-button>
                    <el-button type="primary" link @click="openCustomerFinanceWorkbench">财务</el-button>
                    <el-button type="primary" link @click="openProviderResourcesWorkbench">资源</el-button>
                    <el-button
                      v-if="row.status === 'FAILED'"
                      type="warning"
                      link
                      :loading="retryingSyncLogId === row.id"
                      @click="retrySyncLog(row)"
                    >
                      重试同步
                    </el-button>
                    <el-button
                      v-if="row.status !== 'SUCCESS'"
                      type="primary"
                      link
                      @click="openSyncReviewDialog(row)"
                    >
                      人工收口
                    </el-button>
                    <el-button
                      v-if="row.status === 'FAILED'"
                      type="warning"
                      link
                      @click="openSyncFailureTicket(row)"
                    >
                      异常工单
                    </el-button>
                    <el-button
                      type="primary"
                      link
                      @click="openProviderAutomationContext({ orderId: row.orderId, invoiceId: row.invoiceId })"
                    >
                      任务
                    </el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="资源动作追踪">
            <div class="summary-strip" style="margin-bottom: 16px">
              <div class="summary-pill"><span>总记录</span><strong>{{ resourceActionHistorySummary.total }}</strong></div>
              <div class="summary-pill"><span>同步日志</span><strong>{{ resourceActionHistorySummary.sync }}</strong></div>
              <div class="summary-pill"><span>改配单</span><strong>{{ resourceActionHistorySummary.changeOrders }}</strong></div>
              <div class="summary-pill"><span>审计</span><strong>{{ resourceActionHistorySummary.audit }}</strong></div>
              <div class="summary-pill"><span>待收口</span><strong>{{ resourceActionHistorySummary.pending }}</strong></div>
              <div class="summary-pill"><span>异常</span><strong>{{ resourceActionHistorySummary.failed }}</strong></div>
            </div>
            <div class="table-toolbar">
              <div class="table-toolbar__meta">
                <strong>资源动作追踪台</strong>
                <span>把资源动作、同步、改配和审计合并到一处查看，方便客服和运维连续处理。</span>
              </div>
              <div class="inline-actions">
                <el-button plain @click="openProviderResourcesWorkbench">渠道资源</el-button>
                <el-button plain @click="openProviderAutomationWorkbench">自动化任务</el-button>
                <el-button plain @click="openTicketListWorkbench">工单中心</el-button>
              </div>
            </div>
            <el-table :data="resourceActionHistory" border stripe empty-text="当前服务暂无资源动作追踪记录">
              <el-table-column prop="createdAt" label="时间" min-width="180" />
              <el-table-column prop="category" label="来源" min-width="110" />
              <el-table-column prop="title" label="动作/标题" min-width="220" show-overflow-tooltip />
              <el-table-column label="状态" min-width="150">
                <template #default="{ row }">
                  <el-tag :type="row.statusType" effect="light">{{ row.statusText }}</el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="summary" label="摘要" min-width="320" show-overflow-tooltip />
              <el-table-column label="操作" min-width="360" fixed="right">
                <template #default="{ row }">
                  <div class="inline-actions">
                    <el-button v-if="row.orderId" type="primary" link @click="router.push(`/orders/detail/${row.orderId}`)">订单</el-button>
                    <el-button v-if="row.invoiceId" type="primary" link @click="router.push(`/billing/invoices/${row.invoiceId}`)">账单</el-button>
                    <el-button v-if="row.invoiceId" type="primary" link @click="openPaymentsWorkbench(row.invoiceId)">支付</el-button>
                    <el-button v-if="row.invoiceId" type="primary" link @click="openRefundsWorkbench(row.invoiceId)">退款</el-button>
                    <el-button
                      v-if="row.syncLog && row.syncLog.status === 'FAILED'"
                      type="warning"
                      link
                      @click="retrySyncLog(row.syncLog)"
                    >
                      重试同步
                    </el-button>
                    <el-button
                      v-if="row.syncLog && row.syncLog.status !== 'SUCCESS'"
                      type="primary"
                      link
                      @click="openSyncReviewDialog(row.syncLog)"
                    >
                      人工收口
                    </el-button>
                    <el-button
                      v-if="row.changeOrder && (row.changeOrder.executionStatus === 'FAILED' || row.changeOrder.executionStatus === 'EXECUTE_FAILED')"
                      type="warning"
                      link
                      @click="openChangeOrderFailureTicket(row.changeOrder)"
                    >
                      异常工单
                    </el-button>
                    <el-button type="primary" link @click="openProviderResourcesWorkbench">资源</el-button>
                    <el-button type="primary" link @click="openProviderAutomationWorkbench">任务</el-button>
                  </div>
                </template>
              </el-table-column>
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="改配记录">
            <div class="panel-card">
              <div class="section-card__head">
                <strong>改配单执行记录</strong>
                <span class="section-card__meta">生成改配单、收款执行、退款回退都在这里查看</span>
              </div>
              <div class="summary-strip" style="margin-bottom: 16px">
                <div class="summary-pill"><span>改配总数</span><strong>{{ changeOrderSummary.total }}</strong></div>
                <div class="summary-pill"><span>待支付</span><strong>{{ changeOrderSummary.unpaid }}</strong></div>
                <div class="summary-pill"><span>已支付</span><strong>{{ changeOrderSummary.paid }}</strong></div>
                <div class="summary-pill"><span>已退款</span><strong>{{ changeOrderSummary.refunded }}</strong></div>
                <div class="summary-pill"><span>执行异常</span><strong>{{ changeOrderSummary.failedExecution }}</strong></div>
              </div>
              <div class="table-toolbar">
                <div class="table-toolbar__meta">
                  <strong>改配收口工作台</strong>
                  <span>围绕待支付、执行异常、账单与工单协同处理当前服务的改配单</span>
                </div>
                <div class="inline-actions">
                  <el-button plain @click="openChangeOrdersWorkbench()">全部改配单</el-button>
                  <el-button plain :disabled="changeOrderSummary.unpaid === 0" @click="openChangeOrdersWorkbench({ status: 'UNPAID' })">
                    待支付改配
                  </el-button>
                  <el-button
                    plain
                    :disabled="changeOrderSummary.failedExecution === 0"
                    @click="openChangeOrdersWorkbench({ executionStatus: 'EXECUTE_FAILED' })"
                  >
                    执行异常
                  </el-button>
                  <el-button plain @click="openCustomerFinanceWorkbench">客户财务</el-button>
                  <el-button plain @click="openTicketListWorkbench">工单中心</el-button>
                </div>
              </div>
              <div class="inline-actions" style="margin-bottom: 12px">
                <el-button type="primary" link @click="router.push(`/orders/change-orders?serviceId=${detail.service.id}`)">打开改配单工作台</el-button>
              </div>
              <el-table :data="detail.changeOrders" border stripe empty-text="当前服务还没有改配记录">
                <el-table-column prop="invoiceNo" label="改配账单" min-width="150" />
                <el-table-column prop="orderNo" label="改配订单" min-width="150" />
                <el-table-column label="动作" min-width="130">
                  <template #default="{ row }">{{ changeOrderActionLabel(row.actionName) }}</template>
                </el-table-column>
                <el-table-column prop="title" label="标题" min-width="200" show-overflow-tooltip />
                <el-table-column label="金额" min-width="120">
                  <template #default="{ row }">{{ formatMoney(localeStore.locale, Number(row.amount || 0)) }}</template>
                </el-table-column>
                <el-table-column label="支付状态" min-width="110">
                  <template #default="{ row }">
                    <el-tag :type="row.status === 'PAID' ? 'success' : row.status === 'REFUNDED' ? 'info' : 'warning'" effect="light">
                      {{ changeOrderStatusLabel(row.status) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column label="执行状态" min-width="120">
                  <template #default="{ row }">
                    <el-tag :type="changeOrderExecutionType(row.executionStatus)" effect="light">
                      {{ changeOrderExecutionLabel(row.executionStatus) }}
                    </el-tag>
                  </template>
                </el-table-column>
                <el-table-column prop="createdAt" label="创建时间" min-width="160" />
                <el-table-column label="执行结果" min-width="260" show-overflow-tooltip>
                  <template #default="{ row }">{{ row.executionMessage || "-" }}</template>
                </el-table-column>
                <el-table-column label="操作" min-width="320" fixed="right">
                  <template #default="{ row }">
                    <div class="inline-actions">
                      <el-button type="primary" link @click="router.push(`/orders/detail/${row.orderId}`)">订单</el-button>
                      <el-button type="primary" link @click="router.push(`/billing/invoices/${row.invoiceId}`)">账单</el-button>
                      <el-button type="primary" link @click="openPaymentsWorkbench(row.invoiceId)">支付</el-button>
                      <el-button type="primary" link @click="openRefundsWorkbench(row.invoiceId)">退款</el-button>
                      <el-button
                        v-if="row.executionStatus === 'FAILED' || row.executionStatus === 'EXECUTE_FAILED'"
                        type="warning"
                        link
                        @click="openChangeOrderFailureTicket(row)"
                      >
                        异常工单
                      </el-button>
                      <el-button type="primary" link @click="openProviderAutomationContext({ orderId: row.orderId, invoiceId: row.invoiceId })">
                        任务
                      </el-button>
                    </div>
                  </template>
                </el-table-column>
              </el-table>
            </div>
          </el-tab-pane>

          <el-tab-pane label="自动化任务">
            <AutomationTaskPanel title="服务自动化任务" :service-id="detail.service.id" />
          </el-tab-pane>

          <el-tab-pane label="变更记录">
            <AuditTrailTable :items="detail.auditLogs" empty-text="暂无服务变更记录" />
          </el-tab-pane>
        </el-tabs>
      </template>
    </PageWorkbench>

    <el-dialog v-model="editVisible" title="人工调整服务" width="560px">
      <el-form label-position="top">
        <el-alert
          :title="serviceEditGuide.title"
          :description="serviceEditGuide.description"
          :type="serviceEditGuide.type"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-alert
          :title="serviceEditImpact.title"
          :description="serviceEditImpact.description"
          :type="serviceEditImpact.type"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item label="自动化渠道">
          <el-select v-model="editForm.providerType" style="width: 100%">
            <el-option v-for="item in providerTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="接口账户">
          <el-select v-model="editForm.providerAccountId" style="width: 100%">
            <el-option label="未绑定" :value="0" />
            <el-option
              v-for="item in editProviderAccounts"
              :key="item.id"
              :label="`${item.name} / ${item.baseUrl}`"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="远端资源 ID">
          <el-input v-model="editForm.providerResourceId" placeholder="例如 590" />
        </el-form-item>
        <el-form-item label="地域">
          <el-input v-model="editForm.regionName" placeholder="例如 陕西西安" />
        </el-form-item>
        <el-form-item label="IP 地址">
          <el-input v-model="editForm.ipAddress" placeholder="例如 203.0.113.10" />
        </el-form-item>
        <el-form-item label="下次到期">
          <el-date-picker
            v-model="editForm.nextDueAt"
            type="date"
            value-format="YYYY-MM-DD"
            placeholder="请选择下次到期时间"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="服务状态">
          <el-select v-model="editForm.status" style="width: 100%">
            <el-option v-for="item in serviceStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="同步状态">
          <el-select v-model="editForm.syncStatus" style="width: 100%">
            <el-option v-for="item in syncStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="同步说明">
          <el-input v-model="editForm.syncMessage" type="textarea" :rows="2" placeholder="填写本次人工变更的同步说明" />
        </el-form-item>
        <el-form-item label="变更原因">
          <el-input v-model="editForm.reason" type="textarea" :rows="3" placeholder="例如：客户线下补开、人工纠正渠道绑定、修正同步状态" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingEdit" @click="submitServiceEdit">保存调整</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="syncReviewVisible" title="同步人工收口" width="520px">
      <el-form v-if="reviewingSyncLog" label-position="top">
        <el-alert
          title="用于人工核对远端资源、自动化任务和本地服务状态后，手动完成同步闭环。"
          :description="`当前收口对象：${formatSyncLogAction(localeStore.locale, reviewingSyncLog.action)} / ${formatResourceType(localeStore.locale, reviewingSyncLog.resourceType)} / ${reviewingSyncLog.resourceId || '-'}`"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <div class="summary-strip" style="margin-bottom: 16px">
          <div class="summary-pill">
            <span>日志状态</span>
            <strong>{{ formatSyncStatus(localeStore.locale, reviewingSyncLog.status) }}</strong>
          </div>
          <div class="summary-pill">
            <span>当前服务同步</span>
            <strong>{{ formatSyncStatus(localeStore.locale, detail?.service.syncStatus || 'PENDING') }}</strong>
          </div>
          <div class="summary-pill">
            <span>服务号</span>
            <strong>{{ detail?.service.serviceNo || "-" }}</strong>
          </div>
        </div>
        <el-form-item label="收口后的同步状态">
          <el-select v-model="syncReviewForm.targetStatus" style="width: 100%">
            <el-option v-for="item in syncStatusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="同步说明">
          <el-input
            v-model="syncReviewForm.syncMessage"
            type="textarea"
            :rows="3"
            placeholder="例如：已人工核对远端资源状态正常，本地同步状态改为成功"
          />
        </el-form-item>
        <el-form-item label="变更原因">
          <el-input
            v-model="syncReviewForm.reason"
            type="textarea"
            :rows="3"
            placeholder="例如：重试后远端实例已存在，人工确认收口"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="syncReviewVisible = false">取消</el-button>
        <el-button type="primary" :loading="savingSyncReview" @click="submitSyncReview">确认收口</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resourceFollowUpVisible" title="资源动作收口面板" width="620px">
      <template v-if="resourceFollowUpSummary">
        <el-alert
          :title="`${resourceFollowUpSummary.actionTitle} 已提交`"
          :description="resourceFollowUpSummary.message"
          type="success"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <div class="summary-strip" style="margin-bottom: 16px">
          <div class="summary-pill">
            <span>服务号</span>
            <strong>{{ detail?.service.serviceNo || "-" }}</strong>
          </div>
          <div class="summary-pill">
            <span>远端实例</span>
            <strong>{{ resourceFollowUpSummary.remoteId }}</strong>
          </div>
          <div class="summary-pill">
            <span>资源对象</span>
            <strong>{{ resourceFollowUpSummary.resourceId }}</strong>
          </div>
          <div class="summary-pill">
            <span>同步动作</span>
            <strong>{{ resourceFollowUpSummary.syncOperation }}</strong>
          </div>
          <div class="summary-pill">
            <span>执行附带</span>
            <strong>{{ resourceFollowUpSummary.powered }}</strong>
          </div>
        </div>
        <el-alert
          type="info"
          :closable="false"
          show-icon
          :title="`同步回执：${resourceFollowUpSummary.syncMessage}`"
          description="建议按“自动化任务 -> 渠道资源 -> 人工收口/改配单”的顺序继续处理。"
          style="margin-bottom: 16px"
        />
        <el-alert
          v-if="resourceFollowUpSummary.hasBillingImpact"
          type="warning"
          :closable="false"
          show-icon
          :title="`建议改配金额 ${formatMoney(localeStore.locale, resourceFollowUpSummary.periodAmount)}`"
          :description="resourceFollowUpSummary.summary"
          style="margin-bottom: 16px"
        />
      </template>
      <template #footer>
        <el-button @click="resourceFollowUpVisible = false">关闭</el-button>
        <el-button plain @click="openProviderResourcesWorkbench">渠道资源</el-button>
        <el-button plain @click="openProviderAutomationWorkbench">自动化任务</el-button>
        <el-button plain :loading="syncLoading" @click="handleSync">重新同步服务</el-button>
        <el-button type="primary" plain @click="openResourceFollowUpSyncReview">人工收口</el-button>
        <el-button
          v-if="resourceActionFollowUp?.hasBillingImpact"
          type="warning"
          plain
          :loading="changeOrderLoading"
          @click="handleCreateServiceChangeOrderFromFollowUp"
        >
          生成改配单
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resetVisible" title="重置密码" width="420px">
      <el-form label-position="top">
        <el-form-item label="新密码">
          <el-input v-model="resetForm.password" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resetVisible = false">取消</el-button>
        <el-button type="primary" :loading="actionLoading === 'reset-password'" @click="handleResetPassword">确认重置</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="reinstallVisible" title="重装系统" width="420px">
      <el-form label-position="top">
        <el-form-item label="系统镜像">
          <el-input v-model="reinstallForm.imageName" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reinstallVisible = false">取消</el-button>
        <el-button type="primary" :loading="actionLoading === 'reinstall'" @click="handleReinstall">确认重装</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="resourceVisible" :title="resourceActionGuide.title" width="420px">
      <el-form label-position="top">
        <el-alert
          :title="resourceActionGuide.title"
          :description="resourceActionGuide.description"
          type="info"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-alert
          v-if="resourceBillingImpact.hasImpact"
          title="建议同步生成改配费用"
          :description="`${resourceBillingImpact.summary} 当前按 ${formatBillingCycle(localeStore.locale, resourceBillingImpact.cycle)} 周期建议金额 ${formatMoney(localeStore.locale, resourceBillingImpact.periodAmount)}。执行完成后建议回到账单工作台继续处理。`"
          type="warning"
          :closable="false"
          show-icon
          style="margin-bottom: 16px"
        />
        <el-form-item v-if="resourceAction === 'add-ipv4' || resourceAction === 'add-ipv6'" :label="resourceDialogMeta.countLabel">
          <el-input-number v-model="resourceForm.count" :min="1" :max="10" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="resourceAction === 'add-disk' || resourceAction === 'resize-disk'" :label="resourceDialogMeta.sizeLabel">
          <el-input-number v-model="resourceForm.sizeGb" :min="10" style="width: 100%" />
        </el-form-item>
        <el-form-item v-if="resourceAction === 'add-disk'" :label="resourceDialogMeta.driverLabel">
          <el-select v-model="resourceForm.driver">
            <el-option label="virtio" value="virtio" />
            <el-option label="scsi" value="scsi" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="resourceAction === 'create-snapshot' || resourceAction === 'create-backup'" :label="resourceDialogMeta.nameLabel">
          <el-input v-model="resourceForm.name" :placeholder="resourceDialogMeta.namePlaceholder" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="closeResourceDialog">取消</el-button>
        <el-button v-if="detail?.invoice" plain @click="router.push(`/billing/invoices/${detail.invoice.id}`)">去账单工作台</el-button>
        <el-button
          v-if="resourceBillingImpact.hasImpact"
          type="warning"
          plain
          :loading="changeOrderLoading"
          @click="handleCreateServiceChangeOrder"
        >
          生成改配单
        </el-button>
        <el-button type="primary" :loading="Boolean(resourceLoading)" @click="submitResourceDialog">{{ resourceActionGuide.submitText }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>
