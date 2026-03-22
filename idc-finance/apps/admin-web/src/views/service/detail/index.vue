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
  fetchMofangInstanceDetail,
  fetchMofangServiceResources,
  fetchMofangSyncLogs,
  fetchProviderAccounts,
  fetchServiceDetail,
  runMofangServiceResourceAction,
  runMofangInstanceAction,
  runServiceAction,
  syncMofangService,
  updateServiceRecord,
  type MofangBackup,
  type MofangDisk,
  type MofangInstanceActionResponse,
  type MofangInstanceDetail,
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
import { buildVncPageUrl } from "@/utils/vnc";

type ResourceDialogAction =
  | "add-ipv4"
  | "add-ipv6"
  | "add-disk"
  | "resize-disk"
  | "create-snapshot"
  | "create-backup"
  | "";

type RemoteInstanceAction =
  | "suspend"
  | "power-on"
  | "power-off"
  | "reboot"
  | "hard-power-off"
  | "hard-reboot"
  | "reset-password"
  | "reinstall"
  | "unsuspend"
  | "get-vnc"
  | "rescue-start"
  | "rescue-stop"
  | "lock"
  | "unlock";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();

const loading = ref(false);
const syncLoading = ref(false);
const actionLoading = ref("");
const resourceLoading = ref("");
const changeOrderLoading = ref(false);
const detail = ref<ServiceDetailResponse | null>(null);
const resources = ref<MofangServiceResourcesResponse | null>(null);
const instanceDetail = ref<MofangInstanceDetail | null>(null);
const providerAccounts = ref<ProviderAccount[]>([]);
const syncLogs = ref<MofangSyncLogItem[]>([]);
const remoteActionLoading = ref("");
const consoleVisible = ref(false);
const consolePayload = ref<Record<string, unknown> | null>(null);

const editVisible = ref(false);
const resetVisible = ref(false);
const reinstallVisible = ref(false);
const resourceVisible = ref(false);
const resourceAction = ref<ResourceDialogAction>("");
const resourceTargetId = ref("");

const savingEdit = ref(false);
const resetForm = reactive({ password: "" });
const reinstallForm = reactive({ imageName: "Ubuntu 24.04.1 LTS" });
const resourceForm = reactive({ count: 1, sizeGb: 10, driver: "virtio", name: "" });
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
const isRemoteManagedService = computed(() => isMofangService.value && Boolean(detail.value?.service.providerResourceId));
const customerId = computed(() => detail.value?.service.customerId ?? 0);
const serviceAccount = computed(() => {
  const accountId = detail.value?.service.providerAccountId || 0;
  if (!accountId) return null;
  return providerAccounts.value.find(item => item.id === accountId) ?? null;
});
const remoteRaw = computed(() => toRecord(instanceDetail.value?.raw));
const remoteUser = computed(() => toRecord(remoteRaw.value.user));
const remotePowerState = computed(() => {
  const rawStatus = String(pickString(remoteRaw.value, "status") || instanceDetail.value?.status || "").toLowerCase();
  if (["active", "on", "running"].includes(rawStatus)) return "on";
  if (["suspended", "suspend"].includes(rawStatus)) return "suspend";
  if (["off", "shutdown", "stopped"].includes(rawStatus)) return "off";
  return rawStatus || "unknown";
});
const remoteConsoleEnabled = computed(() => pickBoolean(remoteRaw.value, "vnc") !== false);
const remoteLocked = computed(() => pickBoolean(remoteRaw.value, "lock") === true);
const remoteInRescue = computed(() => pickBoolean(remoteRaw.value, "rescue") === true);
const remoteActionDisabled = computed(() => {
  const state = remotePowerState.value;
  const suspended = state === "suspend";
  const running = state === "on";
  const poweredOff = state === "off";

  return {
    getVnc: !remoteConsoleEnabled.value,
    powerOn: running || suspended,
    powerOff: !running,
    reboot: !running,
    hardPowerOff: !running,
    hardReboot: !running,
    suspend: suspended,
    unsuspend: !suspended,
    rescueStart: suspended || remoteInRescue.value,
    rescueStop: suspended || !remoteInRescue.value,
    lock: remoteLocked.value,
    unlock: !remoteLocked.value,
    resetPassword: suspended || poweredOff,
    reinstall: suspended
  };
});
const remoteConsoleSummary = computed(() => {
  if (!consolePayload.value) return [];
  return [
    {
      label: "控制台地址",
      value: renderValue(
        pickString(consolePayload.value, "vnc_url_https") ||
          pickString(consolePayload.value, "vnc_url") ||
          pickString(consolePayload.value, "vnc_url_http")
      )
    },
    { label: "实例地址", value: renderValue(pickString(consolePayload.value, "ip")) },
    { label: "控制台口令", value: renderValue(pickString(consolePayload.value, "password")) },
    { label: "VNC 密码", value: renderValue(pickString(consolePayload.value, "vnc_pass")) },
    { label: "连接令牌", value: renderValue(pickString(consolePayload.value, "token")) },
    { label: "路径编号", value: renderValue(pickNumber(consolePayload.value, "path")) }
  ];
});
const consoleLaunchUrl = computed(() =>
  buildVncPageUrl(consolePayload.value, {
    ip: pickString(consolePayload.value ?? {}, "ip") || detail.value?.service.ipAddress,
    title:
      instanceDetail.value?.name ||
      pickString(consolePayload.value ?? {}, "hostname") ||
      detail.value?.service.serviceNo
  })
);

function openConsoleWindow(url: string) {
  const opened = window.open(url, "_blank", "noopener,noreferrer");
  return Boolean(opened);
}
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

function toRecord(value: unknown): Record<string, unknown> {
  return value && typeof value === "object" && !Array.isArray(value) ? (value as Record<string, unknown>) : {};
}

function pickString(source: Record<string, unknown>, key: string) {
  const value = source[key];
  return typeof value === "string" ? value : "";
}

function pickNumber(source: Record<string, unknown>, key: string) {
  const value = source[key];
  if (typeof value === "number" && Number.isFinite(value)) return value;
  if (typeof value === "string" && value.trim() !== "") {
    const parsed = Number(value);
    return Number.isFinite(parsed) ? parsed : undefined;
  }
  return undefined;
}

function pickBoolean(source: Record<string, unknown>, key: string) {
  const value = source[key];
  if (typeof value === "boolean") return value;
  if (typeof value === "number") return value > 0;
  if (typeof value === "string") {
    if (["1", "true", "yes", "on"].includes(value.toLowerCase())) return true;
    if (["0", "false", "no", "off"].includes(value.toLowerCase())) return false;
  }
  return undefined;
}

function renderValue(value?: string | number | null) {
  return value === undefined || value === null || value === "" ? "-" : String(value);
}

function formatBoolean(value: boolean) {
  return value ? "是" : "否";
}

function formatOptionalBoolean(value?: boolean) {
  return value === undefined ? "-" : value ? "是" : "否";
}

function formatBandwidth(inbound?: number, outbound?: number) {
  if (inbound === undefined && outbound === undefined) return "-";
  return `${renderValue(inbound)} / ${renderValue(outbound)} Mbps`;
}

function formatTrafficQuota(quota?: number) {
  return quota === undefined ? "-" : `${quota} GB`;
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

function remoteStatusLabel(status: string) {
  const mapping: Record<string, string> = {
    on: "运行中",
    off: "已关机",
    suspend: "已暂停",
    unknown: "未知",
    ACTIVE: "运行中",
    SUSPENDED: "已暂停",
    UNKNOWN: "未知"
  };
  return mapping[status] ?? status ?? "-";
}

function remoteStatusType(status: string) {
  const mapping: Record<string, string> = {
    on: "success",
    off: "info",
    suspend: "warning",
    unknown: "danger",
    ACTIVE: "success",
    SUSPENDED: "warning",
    UNKNOWN: "danger"
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
    instanceDetail.value = null;
    syncLogs.value = [];
    consolePayload.value = null;
    if (detail.value.service.providerType === "MOFANG_CLOUD") {
      const accountId = detail.value.service.providerAccountId || undefined;
      const remoteId = detail.value.service.providerResourceId || "";
      const tasks: [
        Promise<MofangServiceResourcesResponse>,
        Promise<{ items: MofangSyncLogItem[]; total: number }>,
        Promise<MofangInstanceDetail | null>
      ] = [
        fetchMofangServiceResources(route.params.id as string),
        fetchMofangSyncLogs({ serviceId: detail.value.service.id, limit: 20 }),
        remoteId ? fetchMofangInstanceDetail(remoteId, accountId) : Promise.resolve(null)
      ];
      const [resourceData, logData, instanceData] = await Promise.all(tasks);
      resources.value = resourceData;
      syncLogs.value = logData.items;
      instanceDetail.value = instanceData;
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
  const password = resetForm.password.trim();
  if (!password) {
    ElMessage.warning("请输入新的实例密码");
    return;
  }
  actionLoading.value = "reset-password";
  try {
    if (isRemoteManagedService.value) {
      const accountId = detail.value.service.providerAccountId || undefined;
      const result = await runMofangInstanceAction(
        detail.value.service.providerResourceId,
        "reset-password",
        { password },
        accountId
      );
      ElMessage.success(result.message || "远端实例密码重置任务已提交");
      try {
        await syncMofangService(detail.value.service.id);
      } catch {
        // Ignore sync lag and still refresh detail below.
      }
    } else {
      await runServiceAction(detail.value.service.id, "reset-password", { password });
      ElMessage.success("重置密码任务已提交");
    }
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
  const imageName = reinstallForm.imageName.trim();
  if (!imageName) {
    ElMessage.warning("请输入系统镜像名称");
    return;
  }
  actionLoading.value = "reinstall";
  try {
    if (isRemoteManagedService.value) {
      await ElMessageBox.confirm(
        `确认对实例 ${detail.value.service.providerResourceId} 发起重装系统？当前镜像将切换为 ${imageName}。`,
        "远端实例重装确认",
        { type: "warning" }
      );
      const accountId = detail.value.service.providerAccountId || undefined;
      const result = await runMofangInstanceAction(
        detail.value.service.providerResourceId,
        "reinstall",
        { imageName },
        accountId
      );
      ElMessage.success(result.message || "远端实例重装任务已提交");
      try {
        await syncMofangService(detail.value.service.id);
      } catch {
        // Ignore sync lag and still refresh detail below.
      }
    } else {
      await runServiceAction(detail.value.service.id, "reinstall", { imageName });
      ElMessage.success("重装系统任务已提交");
    }
    reinstallVisible.value = false;
    await loadDetail();
  } catch (error: any) {
    if (error === "cancel" || error === "close") return;
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

async function handleRemoteInstanceAction(
  action: RemoteInstanceAction,
  options?: {
    label?: string;
    confirmText?: string;
    syncAfter?: boolean;
  }
) {
  if (!detail.value?.service.providerResourceId) return;

  if (options?.confirmText) {
    try {
      await ElMessageBox.confirm(options.confirmText, "远端实例操作确认", { type: "warning" });
    } catch {
      return;
    }
  }

  remoteActionLoading.value = action;
  try {
    const accountId = detail.value.service.providerAccountId || undefined;
    const result = await runMofangInstanceAction(detail.value.service.providerResourceId, action, undefined, accountId);
    if (action === "get-vnc") {
      consolePayload.value = result.response ?? null;
      const consoleUrl = buildVncPageUrl(consolePayload.value, {
        ip: detail.value.service.ipAddress,
        title: instanceDetail.value?.name || detail.value.service.serviceNo
      });
      if (consoleUrl && openConsoleWindow(consoleUrl)) {
        ElMessage.success(result.message || "已打开控制台");
        return;
      }
      consoleVisible.value = true;
      ElMessage.success(result.message || "已获取控制台连接信息");
      return;
    }

    ElMessage.success(result.message || `${options?.label || action} 已提交`);
    if (options?.syncAfter !== false) {
      try {
        await syncMofangService(detail.value.service.id);
      } catch {
        // Ignore pull-sync lag and still refresh local detail.
      }
    }
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? `${options?.label || action} 失败`);
  } finally {
    remoteActionLoading.value = "";
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
  resourceLoading.value = resourceAction.value;
  try {
    const result = await runMofangServiceResourceAction(detail.value.service.id, resourceAction.value, buildResourcePayload());
    ElMessage.success(result.message);
    closeResourceDialog();
    await loadDetail();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "资源动作执行失败");
  } finally {
    resourceLoading.value = "";
  }
}

async function handleCreateServiceChangeOrder() {
  if (!detail.value || !resourceAction.value || !resourceBillingImpact.value.hasImpact) return;
  changeOrderLoading.value = true;
  try {
    const result = await createServiceChangeOrder(detail.value.service.id, {
      actionName: resourceAction.value,
      title: `${resourceActionGuide.value.title} 改配单`,
      billingCycle: detail.value.invoice?.billingCycle || detail.value.order?.billingCycle || resourceBillingImpact.value.cycle || "monthly",
      amount: Number(resourceBillingImpact.value.periodAmount.toFixed(2)),
      reason: `服务改配：${resourceActionGuide.value.title}`,
      payload: buildResourcePayload()
    });
    ElMessage.success(`已生成改配单：${result.invoice.invoiceNo}`);
    closeResourceDialog();
    await router.push(`/billing/invoices/${result.invoice.id}`);
  } catch (error: any) {
    ElMessage.error(error?.message ?? "生成改配单失败");
  } finally {
    changeOrderLoading.value = false;
  }
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
    const result = await runMofangServiceResourceAction(detail.value.service.id, action, payload);
    ElMessage.success(result.message);
    await loadDetail();
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
                  <strong>实例动作</strong>
                  <span class="section-card__meta">电源、密码、重装和同步</span>
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
                  <el-button type="primary" plain :loading="actionLoading === 'reboot'" @click="handleSimpleAction('reboot')">
                    重启实例
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="isRemoteManagedService && remoteActionDisabled.resetPassword"
                    @click="resetVisible = true"
                  >
                    重置密码
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="isRemoteManagedService && remoteActionDisabled.reinstall"
                    @click="reinstallVisible = true"
                  >
                    重装系统
                  </el-button>
                </div>
                <el-alert
                  style="margin-top: 16px"
                  type="info"
                  :closable="false"
                  show-icon
                  title="服务动作会同时联动本地状态、自动化任务和魔方云资源同步记录。"
                />
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

          <el-tab-pane v-if="instanceDetail" label="远端实例">
            <div class="summary-strip">
              <div class="summary-pill"><span>上游实例</span><strong>{{ instanceDetail.remoteId }}</strong></div>
              <div class="summary-pill"><span>节点</span><strong>{{ renderValue(pickString(remoteRaw, "node_name") || pickString(remoteRaw, "kvmid")) }}</strong></div>
              <div class="summary-pill"><span>上游状态</span><strong>{{ remoteStatusLabel(pickString(remoteRaw, "status") || instanceDetail.status) }}</strong></div>
              <div class="summary-pill"><span>登录用户</span><strong>{{ renderValue(pickString(remoteRaw, "osuser") || detail.service.resourceSnapshot.loginUsername) }}</strong></div>
            </div>

            <el-alert
              style="margin-top: 16px"
              type="warning"
              :closable="false"
              show-icon
              title="以下动作会直接操作真实魔方云实例。"
              description="建议先核对实例编号、节点、当前上游状态，再执行开关机、救援、锁定或控制台操作。"
            />

            <div class="portal-grid portal-grid--two" style="margin-top: 16px">
              <div class="panel-card">
                <div class="section-card__head">
                  <strong>远端运行信息</strong>
                  <span class="section-card__meta">直接来自魔方云实例详情</span>
                </div>
                <el-descriptions :column="2" border>
                  <el-descriptions-item label="实例名称">{{ renderValue(instanceDetail.name) }}</el-descriptions-item>
                  <el-descriptions-item label="状态">
                    <el-tag :type="remoteStatusType(pickString(remoteRaw, 'status') || instanceDetail.status)" effect="light">
                      {{ remoteStatusLabel(pickString(remoteRaw, "status") || instanceDetail.status) }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="主机名">{{ renderValue(pickString(remoteRaw, "hostname") || detail.service.resourceSnapshot.hostname) }}</el-descriptions-item>
                  <el-descriptions-item label="UUID">{{ renderValue(pickString(remoteRaw, "uuid")) }}</el-descriptions-item>
                  <el-descriptions-item label="地域">{{ renderValue(instanceDetail.region || pickString(remoteRaw, "area_name") || detail.service.regionName) }}</el-descriptions-item>
                  <el-descriptions-item label="节点">{{ renderValue(pickString(remoteRaw, "node_name")) }}</el-descriptions-item>
                  <el-descriptions-item label="主 IP">{{ renderValue(pickString(remoteRaw, "mainip") || instanceDetail.ipAddress || detail.service.ipAddress) }}</el-descriptions-item>
                  <el-descriptions-item label="网络类型">{{ renderValue(pickString(remoteRaw, "network_type")) }}</el-descriptions-item>
                  <el-descriptions-item label="登录用户">{{ renderValue(pickString(remoteRaw, "osuser") || detail.service.resourceSnapshot.loginUsername) }}</el-descriptions-item>
                  <el-descriptions-item label="系统镜像">{{ renderValue(pickString(remoteRaw, "operate_system") || detail.service.resourceSnapshot.operatingSystem) }}</el-descriptions-item>
                  <el-descriptions-item label="CPU / 内存">{{ `${renderValue(pickNumber(remoteRaw, "cpu"))} 核 / ${renderValue(pickNumber(remoteRaw, "memory"))} GB` }}</el-descriptions-item>
                  <el-descriptions-item label="带宽">{{ formatBandwidth(pickNumber(remoteRaw, "in_bw"), pickNumber(remoteRaw, "out_bw")) }}</el-descriptions-item>
                  <el-descriptions-item label="流量配额">{{ formatTrafficQuota(pickNumber(remoteRaw, "traffic_quota")) }}</el-descriptions-item>
                  <el-descriptions-item label="IPv4 / IPv6">{{ `${renderValue(pickNumber(remoteRaw, "ip_num"))} / ${renderValue(pickNumber(remoteRaw, "ipv6_num"))}` }}</el-descriptions-item>
                  <el-descriptions-item label="控制台">{{ formatOptionalBoolean(pickBoolean(remoteRaw, "vnc")) }}</el-descriptions-item>
                  <el-descriptions-item label="实例锁定">{{ formatOptionalBoolean(pickBoolean(remoteRaw, "lock")) }}</el-descriptions-item>
                  <el-descriptions-item label="救援模式">{{ formatOptionalBoolean(pickBoolean(remoteRaw, "rescue")) }}</el-descriptions-item>
                  <el-descriptions-item label="开通时间">{{ renderValue(pickString(remoteRaw, "create_time")) }}</el-descriptions-item>
                  <el-descriptions-item label="上游客户">
                    {{ renderValue(pickString(remoteUser, "email") || pickString(remoteUser, "username") || pickString(remoteRaw, "username")) }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>

              <div class="panel-card">
                <div class="section-card__head">
                  <strong>远端控制动作</strong>
                  <span class="section-card__meta">控制台、电源、救援、锁定</span>
                </div>
                <div class="command-group__actions">
                  <el-button
                    type="primary"
                    plain
                    :disabled="remoteActionDisabled.getVnc"
                    :loading="remoteActionLoading === 'get-vnc'"
                    @click="handleRemoteInstanceAction('get-vnc', { label: '获取控制台', syncAfter: false })"
                  >
                    获取控制台
                  </el-button>
                  <el-button
                    type="success"
                    plain
                    :disabled="remoteActionDisabled.powerOn"
                    :loading="remoteActionLoading === 'power-on'"
                    @click="handleRemoteInstanceAction('power-on', { label: '实例开机' })"
                  >
                    开机
                  </el-button>
                  <el-button
                    type="warning"
                    plain
                    :disabled="remoteActionDisabled.powerOff"
                    :loading="remoteActionLoading === 'power-off'"
                    @click="handleRemoteInstanceAction('power-off', { label: '实例关机', confirmText: '确认对当前实例执行关机？' })"
                  >
                    关机
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="remoteActionDisabled.reboot"
                    :loading="remoteActionLoading === 'reboot'"
                    @click="handleRemoteInstanceAction('reboot', { label: '实例重启', confirmText: '确认重启当前实例？' })"
                  >
                    重启实例
                  </el-button>
                  <el-button
                    type="danger"
                    plain
                    :disabled="remoteActionDisabled.hardPowerOff"
                    :loading="remoteActionLoading === 'hard-power-off'"
                    @click="handleRemoteInstanceAction('hard-power-off', { label: '强制关机', confirmText: '确认强制关闭当前实例？这会直接中断运行中的业务。' })"
                  >
                    强制关机
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="remoteActionDisabled.hardReboot"
                    :loading="remoteActionLoading === 'hard-reboot'"
                    @click="handleRemoteInstanceAction('hard-reboot', { label: '强制重启', confirmText: '确认强制重启当前实例？' })"
                  >
                    强制重启
                  </el-button>
                  <el-button
                    type="warning"
                    plain
                    :disabled="remoteActionDisabled.suspend"
                    :loading="remoteActionLoading === 'suspend'"
                    @click="handleRemoteInstanceAction('suspend', { label: '暂停实例', confirmText: '确认暂停当前实例？暂停后实例将无法继续运行。' })"
                  >
                    暂停实例
                  </el-button>
                  <el-button
                    type="warning"
                    plain
                    :disabled="remoteActionDisabled.unsuspend"
                    :loading="remoteActionLoading === 'unsuspend'"
                    @click="handleRemoteInstanceAction('unsuspend', { label: '解除暂停' })"
                  >
                    解除暂停
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="remoteActionDisabled.rescueStart"
                    :loading="remoteActionLoading === 'rescue-start'"
                    @click="handleRemoteInstanceAction('rescue-start', { label: '进入救援', confirmText: '确认让实例进入救援模式？' })"
                  >
                    进入救援
                  </el-button>
                  <el-button
                    type="primary"
                    plain
                    :disabled="remoteActionDisabled.rescueStop"
                    :loading="remoteActionLoading === 'rescue-stop'"
                    @click="handleRemoteInstanceAction('rescue-stop', { label: '退出救援' })"
                  >
                    退出救援
                  </el-button>
                  <el-button
                    type="danger"
                    plain
                    :disabled="remoteActionDisabled.lock"
                    :loading="remoteActionLoading === 'lock'"
                    @click="handleRemoteInstanceAction('lock', { label: '锁定实例', confirmText: '确认锁定当前实例？锁定后部分操作会被限制。' })"
                  >
                    锁定实例
                  </el-button>
                  <el-button
                    type="success"
                    plain
                    :disabled="remoteActionDisabled.unlock"
                    :loading="remoteActionLoading === 'unlock'"
                    @click="handleRemoteInstanceAction('unlock', { label: '解除锁定' })"
                  >
                    解除锁定
                  </el-button>
                  <el-button type="primary" plain :disabled="remoteActionDisabled.resetPassword" @click="resetVisible = true">
                    远端重置密码
                  </el-button>
                  <el-button type="primary" plain :disabled="remoteActionDisabled.reinstall" @click="reinstallVisible = true">
                    远端重装系统
                  </el-button>
                </div>

                <el-alert
                  style="margin-top: 16px"
                  type="info"
                  :closable="false"
                  show-icon
                  title="远端动作执行成功后，页面会自动拉取同步，方便你立即看到服务状态和资源快照变化。"
                />

                <el-descriptions v-if="consolePayload" :column="1" border style="margin-top: 16px">
                  <el-descriptions-item v-for="item in remoteConsoleSummary" :key="item.label" :label="item.label">
                    {{ item.value }}
                  </el-descriptions-item>
                </el-descriptions>
              </div>
            </div>
          </el-tab-pane>

          <el-tab-pane v-if="resources" label="魔方云资源">
            <div class="summary-strip">
              <div class="summary-pill"><span>IP</span><strong>{{ resources.ipAddresses.length }}</strong></div>
              <div class="summary-pill"><span>磁盘</span><strong>{{ resources.disks.length }}</strong></div>
              <div class="summary-pill"><span>快照</span><strong>{{ resources.snapshots.length }}</strong></div>
              <div class="summary-pill"><span>备份</span><strong>{{ resources.backups.length }}</strong></div>
            </div>

            <el-alert
              style="margin-top: 16px"
              type="info"
              :closable="false"
              show-icon
              title="这里的操作会直接调用魔方云接口，并把结果写入同步日志和自动化任务。"
              description="建议按“网络扩容、存储扩容、数据保护”顺序操作。做扩容或保护前，先核对当前实例编号、磁盘列表和公网地址。"
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
            <el-table :data="syncLogs" border stripe empty-text="当前服务暂无同步日志">
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
            </el-table>
          </el-tab-pane>

          <el-tab-pane label="改配记录">
            <div class="panel-card">
              <div class="section-card__head">
                <strong>改配单执行记录</strong>
                <span class="section-card__meta">生成改配单、收款执行、退款回退都在这里查看</span>
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
                <el-table-column label="操作" min-width="150" fixed="right">
                  <template #default="{ row }">
                    <div class="inline-actions">
                      <el-button type="primary" link @click="router.push(`/orders/detail/${row.orderId}`)">订单</el-button>
                      <el-button type="primary" link @click="router.push(`/billing/invoices/${row.invoiceId}`)">账单</el-button>
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

    <el-dialog v-model="resetVisible" title="重置密码" width="420px">
      <el-form label-position="top">
        <el-alert
          v-if="isRemoteManagedService"
          type="warning"
          :closable="false"
          show-icon
          title="该操作会直接下发到真实魔方云实例。"
          style="margin-bottom: 16px"
        />
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
        <el-alert
          v-if="isRemoteManagedService"
          type="warning"
          :closable="false"
          show-icon
          title="重装会直接作用于真实实例，请确认镜像名称和业务窗口。"
          style="margin-bottom: 16px"
        />
        <el-form-item label="系统镜像">
          <el-input v-model="reinstallForm.imageName" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reinstallVisible = false">取消</el-button>
        <el-button type="primary" :loading="actionLoading === 'reinstall'" @click="handleReinstall">确认重装</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="consoleVisible" title="控制台连接信息" width="620px">
      <el-alert
        type="warning"
        :closable="false"
        show-icon
        title="以下信息包含真实实例控制台凭证，请按最小范围使用。"
        style="margin-bottom: 16px"
      />
      <el-descriptions :column="1" border>
        <el-descriptions-item v-for="item in remoteConsoleSummary" :key="item.label" :label="item.label">
          {{ item.value }}
        </el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="consoleVisible = false">关闭</el-button>
        <el-button v-if="consoleLaunchUrl" type="primary" @click="openConsoleWindow(consoleLaunchUrl)">
          打开控制台
        </el-button>
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
