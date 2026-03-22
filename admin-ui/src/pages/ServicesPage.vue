<script setup lang="ts">
import {
  Monitor,
  Refresh,
  Search,
  Upload,
  VideoPause,
  VideoPlay,
  DocumentCopy,
  Download,
} from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";

import ListWorkbenchShell from "@/components/common/ListWorkbenchShell.vue";
import { http } from "@/api/http";
import { downloadCsv } from "@/utils/download";
import { formatCurrency, formatDate } from "@/utils/format";
import {
  cycleMap,
  getLabel,
  getStatusTagType,
  providerActionMap,
  providerTypeMap,
  serviceStatusMap,
} from "@/utils/maps";

type ServiceRow = {
  id: string;
  name: string;
  serviceNo: string;
  hostname?: string | null;
  region?: string | null;
  status: string;
  monthlyCost: number;
  nextDueDate?: string | null;
  providerType: string;
  providerResourceId?: string | null;
  ipAddress?: string | null;
  billingCycle: string;
  configSnapshot?: string | null;
  customer: { id: string; name: string };
  product: { name: string };
  plan?: { name?: string | null } | null;
  vpcNetwork?: { name?: string | null } | null;
  _count: {
    ipAddresses: number;
    disks: number;
    snapshots: number;
    backups: number;
    securityGroups: number;
  };
};

const router = useRouter();

const loading = ref(false);
const providerSyncLoading = ref(false);
const advancedOpen = ref(false);
const statusFilter = ref("ALL");
const keyword = ref("");
const services = ref<ServiceRow[]>([]);
const selectedRows = ref<ServiceRow[]>([]);

const filters = reactive({
  showLiveOnly: true,
  providerType: "ALL",
  billingCycle: "ALL",
  region: "ALL",
  customerId: "ALL",
  monthlyCostMin: undefined as number | undefined,
  monthlyCostMax: undefined as number | undefined,
  dueRange: [] as [Date, Date] | [],
});

function parseConfigSnapshot(value?: string | null) {
  if (!value) {
    return {};
  }

  try {
    const parsed = JSON.parse(value) as Record<string, unknown>;
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

function isLiveMofangService(service: ServiceRow) {
  if (service.providerType !== "MOFANG_CLOUD") {
    return false;
  }

  const config = parseConfigSnapshot(service.configSnapshot);
  return config.source === "mofang-pull-sync";
}

const uniqueRegions = computed(() =>
  Array.from(new Set(services.value.map((item) => item.region).filter(Boolean))).sort(),
);

const uniqueCustomers = computed(() =>
  services.value
    .map((item) => item.customer)
    .filter(
      (item, index, array) => array.findIndex((target) => target.id === item.id) === index,
    )
    .sort((a, b) => a.name.localeCompare(b.name, "zh-CN")),
);

const liveServiceCount = computed(() => services.value.filter(isLiveMofangService).length);

const statusCards = computed(() => {
  const active = services.value.filter((item) => item.status === "ACTIVE").length;
  const suspended = services.value.filter((item) => item.status === "SUSPENDED").length;
  const overdue = services.value.filter((item) => item.status === "OVERDUE").length;

  return [
    { key: "ALL", label: "全部业务", value: services.value.length, hint: "系统内所有服务实例" },
    { key: "ACTIVE", label: "运行中", value: active, hint: "当前正常提供服务" },
    { key: "SUSPENDED", label: "已暂停", value: suspended, hint: "需要恢复或续费处理" },
    { key: "OVERDUE", label: "已逾期", value: overdue, hint: "已进入催缴或停机链路" },
  ];
});

const summaryTotals = computed(() => {
  const totalMonthlyCost = services.value.reduce((sum, item) => sum + item.monthlyCost, 0);
  const manualServiceCount = services.value.filter((item) => item.providerType === "MANUAL").length;

  return [
    {
      label: "正式魔方云",
      value: liveServiceCount.value,
      hint: "真实魔方云同步实例",
    },
    {
      label: "人工服务",
      value: manualServiceCount,
      hint: "本地人工交付实例",
    },
    {
      label: "月度成本",
      value: formatCurrency(totalMonthlyCost),
      hint: "全部服务当前计费成本",
    },
  ];
});

const filteredServices = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLowerCase();
  const [dueFrom, dueTo] = filters.dueRange;

  return services.value.filter((service) => {
    const matchLive = !filters.showLiveOnly || isLiveMofangService(service);
    const matchStatus = statusFilter.value === "ALL" || service.status === statusFilter.value;
    const matchProvider =
      filters.providerType === "ALL" || service.providerType === filters.providerType;
    const matchCycle =
      filters.billingCycle === "ALL" || service.billingCycle === filters.billingCycle;
    const matchRegion = filters.region === "ALL" || service.region === filters.region;
    const matchCustomer =
      filters.customerId === "ALL" || service.customer.id === filters.customerId;
    const matchKeyword =
      normalizedKeyword.length === 0 ||
      [
        service.name,
        service.serviceNo,
        service.hostname,
        service.ipAddress,
        service.customer.name,
        service.product.name,
        service.providerResourceId,
      ]
        .some((field) => String(field ?? "").toLowerCase().includes(normalizedKeyword));

    const monthlyCostYuan = service.monthlyCost / 100;
    const matchCostMin =
      typeof filters.monthlyCostMin !== "number" || monthlyCostYuan >= filters.monthlyCostMin;
    const matchCostMax =
      typeof filters.monthlyCostMax !== "number" || monthlyCostYuan <= filters.monthlyCostMax;

    const nextDueDate = service.nextDueDate ? new Date(service.nextDueDate) : null;
    const matchDueFrom = !dueFrom || !nextDueDate || nextDueDate >= dueFrom;
    const matchDueTo = !dueTo || !nextDueDate || nextDueDate <= dueTo;

    return (
      matchLive &&
      matchStatus &&
      matchProvider &&
      matchCycle &&
      matchRegion &&
      matchCustomer &&
      matchKeyword &&
      matchCostMin &&
      matchCostMax &&
      matchDueFrom &&
      matchDueTo
    );
  });
});

const overdueOrSuspendedCount = computed(
  () =>
    filteredServices.value.filter((item) =>
      ["SUSPENDED", "OVERDUE", "EXPIRED"].includes(item.status),
    ).length,
);

function resetAdvancedFilters() {
  Object.assign(filters, {
    showLiveOnly: true,
    providerType: "ALL",
    billingCycle: "ALL",
    region: "ALL",
    customerId: "ALL",
    monthlyCostMin: undefined,
    monthlyCostMax: undefined,
    dueRange: [],
  });
}

function handleSelectionChange(rows: ServiceRow[]) {
  selectedRows.value = rows;
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/services");
    services.value = data.data;
  } finally {
    loading.value = false;
  }
}

async function showProviderResult(providerResult: any, action: string) {
  if (!providerResult) {
    return;
  }

  if (action === "getVnc" && providerResult.consoleUrl) {
    window.open(providerResult.consoleUrl, "_blank", "noopener,noreferrer");
  }

  const extraLines = [
    providerResult.message ? `返回消息：${providerResult.message}` : "",
    providerResult.taskId ? `任务 ID：${providerResult.taskId}` : "",
    providerResult.consoleUrl
      ? `控制台入口：<a href="${providerResult.consoleUrl}" target="_blank" rel="noreferrer">打开 VNC 控制台</a>`
      : "",
  ].filter(Boolean);

  if (extraLines.length === 0) {
    return;
  }

  await ElMessageBox.alert(extraLines.join("<br/>"), "云平台回执", {
    dangerouslyUseHTMLString: true,
    confirmButtonText: "知道了",
  });
}

async function executeAction(
  serviceId: string,
  action: string,
  options?: {
    silent?: boolean;
    reloadAfter?: boolean;
    showProviderDialog?: boolean;
  },
) {
  const { silent = false, reloadAfter = true, showProviderDialog = true } = options ?? {};

  try {
    const { data } = await http.post(`/services/${serviceId}/action`, {
      action,
    });

    if (!silent) {
      ElMessage.success(`已执行${getLabel(providerActionMap, action)}`);
    }

    if (showProviderDialog) {
      await showProviderResult(data.providerResult, action);
    }

    if (reloadAfter) {
      await loadData();
    }

    return true;
  } catch (error: any) {
    const providerResult = error?.response?.data?.providerResult;
    const message = error?.response?.data?.message || `执行${getLabel(providerActionMap, action)}失败`;

    if (!silent) {
      ElMessage.error(message);
    }

    if (showProviderDialog && providerResult) {
      await showProviderResult(providerResult, action);
    }

    return false;
  }
}

async function runBatchAction(action: "sync" | "powerOn" | "powerOff") {
  if (!selectedRows.value.length) {
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确认对已选 ${selectedRows.value.length} 台服务执行“${getLabel(providerActionMap, action)}”吗？`,
      "批量操作确认",
      {
        type: "warning",
      },
    );
  } catch {
    return;
  }

  let successCount = 0;
  let failedCount = 0;

  for (const service of selectedRows.value) {
    const ok = await executeAction(service.id, action, {
      silent: true,
      reloadAfter: false,
      showProviderDialog: false,
    });
    if (ok) {
      successCount += 1;
    } else {
      failedCount += 1;
    }
  }

  await loadData();
  ElMessage.success(
    `${getLabel(providerActionMap, action)}已执行，成功 ${successCount} 台，失败 ${failedCount} 台`,
  );
}

async function runProviderSync() {
  providerSyncLoading.value = true;
  try {
    const { data } = await http.post("/services/sync/provider", {
      limit: 100,
      includeResources: true,
    });
    const summary = data.data.summary;
    await ElMessageBox.alert(
      [
        `远端实例：${summary.remoteCount}`,
        `本次处理：${summary.processedCount}`,
        `新增服务：${summary.createdServices}`,
        `更新服务：${summary.updatedServices}`,
        `失败服务：${summary.failedServices}`,
        `新增门户账号：${summary.createdPortalUsers}`,
        `新增导入订单：${summary.createdImportOrders}`,
        `写入门户通知：${summary.queuedPortalNotifications}`,
        `同步 IP：${summary.syncedIps}`,
        `同步磁盘：${summary.syncedDisks}`,
        `同步快照：${summary.syncedSnapshots}`,
        `同步备份：${summary.syncedBackups}`,
      ].join("<br/>"),
      "魔方云同步结果",
      {
        dangerouslyUseHTMLString: true,
        confirmButtonText: "知道了",
      },
    );
    await loadData();
  } catch (error: any) {
    const message = error?.response?.data?.message || "同步魔方云失败";
    ElMessage.error(message);
  } finally {
    providerSyncLoading.value = false;
  }
}

async function copySelectedServiceNos() {
  const serviceNos = selectedRows.value.map((item) => item.serviceNo).join("\n");
  await navigator.clipboard.writeText(serviceNos);
  ElMessage.success("已复制所选服务编号");
}

function exportServices() {
  const rows = (selectedRows.value.length > 0 ? selectedRows.value : filteredServices.value).map(
    (item) => [
      item.serviceNo,
      item.name,
      item.customer.name,
      item.product.name,
      getLabel(providerTypeMap, item.providerType),
      item.providerResourceId || "-",
      item.ipAddress || "-",
      item.region || "-",
      getLabel(cycleMap, item.billingCycle),
      formatCurrency(item.monthlyCost),
      formatDate(item.nextDueDate),
      getLabel(serviceStatusMap, item.status),
    ],
  );

  downloadCsv(
    `services-${new Date().toISOString().slice(0, 10)}.csv`,
    [
      "服务编号",
      "实例名称",
      "客户",
      "产品",
      "接入方式",
      "远端资源 ID",
      "IP 地址",
      "地域",
      "计费周期",
      "月费",
      "下次到期",
      "状态",
    ],
    rows,
  );
  ElMessage.success("服务数据已导出");
}

onMounted(loadData);
</script>

<template>
  <ListWorkbenchShell
    v-model:advancedOpen="advancedOpen"
    title="业务列表"
    subtitle="业务列表负责管理已经开通的服务实例，列表页提供筛选、批量动作和高频运维入口，复杂操作统一进入服务详情工作台处理。"
    :result-count="filteredServices.length"
    :selected-count="selectedRows.length"
    :show-advanced-toggle="true"
    advanced-button-text="高级筛选"
  >
    <template #actions>
      <el-button :icon="Refresh" @click="loadData">刷新数据</el-button>
      <el-button
        type="primary"
        :icon="Upload"
        :loading="providerSyncLoading"
        @click="runProviderSync"
      >
        同步魔方云
      </el-button>
    </template>

    <template #summary>
      <button
        v-for="card in statusCards"
        :key="card.key"
        class="status-tab"
        :class="{ 'is-active': statusFilter === card.key }"
        @click="statusFilter = card.key"
      >
        <div class="status-tab-label">{{ card.label }}</div>
        <div class="status-tab-value">{{ card.value }}</div>
        <div class="muted-line">{{ card.hint }}</div>
      </button>

      <div v-for="card in summaryTotals" :key="card.label" class="summary-card">
        <div class="status-tab-label">{{ card.label }}</div>
        <h3>{{ card.value }}</h3>
        <p>{{ card.hint }}</p>
      </div>
    </template>

    <template #filters>
      <el-input
        v-model="keyword"
        clearable
        placeholder="搜索服务名、编号、客户名、IP、主机名、远端资源 ID"
        style="width: 380px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-switch
        v-model="filters.showLiveOnly"
        inline-prompt
        active-text="正式"
        inactive-text="全部"
      />
    </template>

    <template #filterMeta>
      <el-tag effect="plain">正式实例 {{ liveServiceCount }}</el-tag>
      <el-tag effect="plain">风险业务 {{ overdueOrSuspendedCount }}</el-tag>
      <el-tag effect="plain">当前结果 {{ filteredServices.length }}</el-tag>
    </template>

    <template #advanced>
      <div class="advanced-filter-grid">
        <el-select v-model="filters.providerType">
          <el-option label="全部接入方式" value="ALL" />
          <el-option
            v-for="(label, value) in providerTypeMap"
            :key="value"
            :label="label"
            :value="value"
          />
        </el-select>
        <el-select v-model="filters.billingCycle">
          <el-option label="全部周期" value="ALL" />
          <el-option
            v-for="(label, value) in cycleMap"
            :key="value"
            :label="label"
            :value="value"
          />
        </el-select>
        <el-select v-model="filters.region">
          <el-option label="全部地域" value="ALL" />
          <el-option
            v-for="region in uniqueRegions"
            :key="region"
            :label="region"
            :value="region"
          />
        </el-select>
        <el-select v-model="filters.customerId" filterable>
          <el-option label="全部客户" value="ALL" />
          <el-option
            v-for="customer in uniqueCustomers"
            :key="customer.id"
            :label="customer.name"
            :value="customer.id"
          />
        </el-select>
        <el-input-number
          v-model="filters.monthlyCostMin"
          :precision="2"
          :min="0"
          style="width: 100%"
          placeholder="最小月费"
        />
        <el-input-number
          v-model="filters.monthlyCostMax"
          :precision="2"
          :min="0"
          style="width: 100%"
          placeholder="最大月费"
        />
        <el-date-picker
          v-model="filters.dueRange"
          type="daterange"
          range-separator="至"
          start-placeholder="到期开始"
          end-placeholder="到期结束"
          style="width: 100%"
        />
        <div class="inline-actions">
          <el-button @click="resetAdvancedFilters">清空高级筛选</el-button>
        </div>
      </div>
    </template>

    <template #selection>
      <el-button plain :icon="DocumentCopy" @click="copySelectedServiceNos">复制服务编号</el-button>
      <el-button plain :icon="Download" @click="exportServices">导出所选服务</el-button>
      <el-button plain @click="runBatchAction('sync')">批量同步</el-button>
      <el-button plain type="success" @click="runBatchAction('powerOn')">批量开机</el-button>
      <el-button plain type="warning" @click="runBatchAction('powerOff')">批量关机</el-button>
    </template>

    <el-table
      v-loading="loading"
      :data="filteredServices"
      stripe
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="服务信息" min-width="260">
        <template #default="{ row }">
          <div class="primary-line">{{ row.name }}</div>
          <div class="muted-line">
            {{ row.serviceNo }} / {{ row.customer.name }} / {{ row.product.name }}
          </div>
          <div class="muted-line">{{ row.plan?.name || "未绑定售卖方案" }}</div>
        </template>
      </el-table-column>

      <el-table-column label="接入信息" min-width="220">
        <template #default="{ row }">
          <div>{{ row.hostname || "-" }}</div>
          <div class="muted-line">{{ row.ipAddress || "未分配 IP" }}</div>
          <div class="muted-line">资源 ID：{{ row.providerResourceId || "-" }}</div>
        </template>
      </el-table-column>

      <el-table-column label="计费信息" width="170">
        <template #default="{ row }">
          <div>{{ getLabel(cycleMap, row.billingCycle) }}</div>
          <div class="muted-line">{{ formatCurrency(row.monthlyCost) }}</div>
        </template>
      </el-table-column>

      <el-table-column label="资源归属" min-width="180">
        <template #default="{ row }">
          <div>{{ row.region || "-" }}</div>
          <div class="muted-line">{{ row.vpcNetwork?.name || "未绑定 VPC" }}</div>
        </template>
      </el-table-column>

      <el-table-column label="资源统计" width="220">
        <template #default="{ row }">
          <div>IP {{ row._count.ipAddresses }} / 磁盘 {{ row._count.disks }}</div>
          <div class="muted-line">
            快照 {{ row._count.snapshots }} / 备份 {{ row._count.backups }} / 安全组
            {{ row._count.securityGroups }}
          </div>
        </template>
      </el-table-column>

      <el-table-column label="接入方式" width="120">
        <template #default="{ row }">
          {{ getLabel(providerTypeMap, row.providerType) }}
        </template>
      </el-table-column>

      <el-table-column label="下次到期" width="130">
        <template #default="{ row }">
          {{ formatDate(row.nextDueDate) }}
        </template>
      </el-table-column>

      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)">
            {{ getLabel(serviceStatusMap, row.status) }}
          </el-tag>
        </template>
      </el-table-column>

      <el-table-column label="操作" min-width="320" fixed="right">
        <template #default="{ row }">
          <div class="table-actions">
            <el-button text type="primary" :icon="Monitor" @click="router.push(`/services/${row.id}`)">
              详情
            </el-button>
            <el-button text :icon="Refresh" @click="executeAction(row.id, 'sync')">同步</el-button>
            <el-button text :icon="VideoPlay" @click="executeAction(row.id, 'powerOn')">开机</el-button>
            <el-button text :icon="VideoPause" @click="executeAction(row.id, 'powerOff')">关机</el-button>
            <el-button text @click="executeAction(row.id, 'getVnc')">VNC</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </ListWorkbenchShell>
</template>
