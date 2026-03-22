<script setup lang="ts">
import {
  Plus,
  Refresh,
  Search,
  DocumentCopy,
  Download,
} from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { computed, onMounted, reactive, ref, watch } from "vue";
import { useRouter } from "vue-router";

import ListWorkbenchShell from "@/components/common/ListWorkbenchShell.vue";
import { http } from "@/api/http";
import { useAuthStore } from "@/stores/auth";
import { downloadCsv } from "@/utils/download";
import { formatCurrency, formatDate } from "@/utils/format";
import { getLabel, getStatusTagType, invoiceStatusMap, invoiceTypeMap } from "@/utils/maps";

const authStore = useAuthStore();
const router = useRouter();

const loading = ref(false);
const creating = ref(false);
const actionSubmitting = ref(false);
const createDialogVisible = ref(false);
const actionDialogVisible = ref(false);
const advancedOpen = ref(false);
const keyword = ref("");
const statusFilter = ref("ALL");

const invoices = ref<any[]>([]);
const customers = ref<any[]>([]);
const services = ref<any[]>([]);
const taxProfiles = ref<any[]>([]);
const selectedRows = ref<any[]>([]);
const billingSetting = ref<any>({
  defaultTaxRate: 13,
  invoicePrefix: "INV",
});
const summary = ref({
  invoiceCount: 0,
  draftCount: 0,
  openCount: 0,
  overdueCount: 0,
  paidCount: 0,
  voidCount: 0,
  totalIssued: 0,
  totalTax: 0,
  totalOutstanding: 0,
});

const filters = reactive({
  customerId: "ALL",
  type: "ALL",
  amountMin: undefined as number | undefined,
  amountMax: undefined as number | undefined,
  dueRange: [] as [Date, Date] | [],
  taxKeyword: "",
});

const createForm = reactive({
  customerId: "",
  serviceId: "",
  taxProfileId: "",
  type: "MANUAL",
  subtotal: 0,
  taxRate: 13,
  dueDate: "",
  remark: "",
  issueNow: true,
});

const actionForm = reactive({
  invoiceId: "",
  action: "issue",
  dueDate: "",
  reason: "",
});

const canManageInvoices = computed(() =>
  Boolean(authStore.user?.permissions.includes("invoices.manage")),
);

const customerServices = computed(() =>
  services.value.filter((service) => service.customerId === createForm.customerId),
);

const selectedDraftCount = computed(
  () => selectedRows.value.filter((item) => item.status === "DRAFT").length,
);

const selectedVoidableCount = computed(
  () =>
    selectedRows.value.filter((item) => item.status !== "VOID" && Number(item.paidAmount) === 0)
      .length,
);

const createPreview = computed(() => {
  const subtotal = Math.round(Number(createForm.subtotal || 0) * 100);
  const taxAmount = Math.round((subtotal * Number(createForm.taxRate || 0)) / 100);
  return {
    subtotal,
    taxAmount,
    totalAmount: subtotal + taxAmount,
  };
});

const statusCards = computed(() => [
  { key: "ALL", label: "全部账单", value: summary.value.invoiceCount, hint: "当前系统全部账单" },
  { key: "ISSUED", label: "待支付", value: summary.value.openCount, hint: "已签发且未结清账单" },
  { key: "OVERDUE", label: "已逾期", value: summary.value.overdueCount, hint: "超过到期日仍未结清" },
  { key: "PAID", label: "已支付", value: summary.value.paidCount, hint: "已完成收款入账" },
  { key: "DRAFT", label: "草稿账单", value: summary.value.draftCount, hint: "尚未签发的账单" },
  { key: "VOID", label: "已作废", value: summary.value.voidCount, hint: "已关闭不可再收款" },
]);

const summaryTotals = computed(() => [
  {
    label: "总开票额",
    value: formatCurrency(summary.value.totalIssued),
    hint: "全部非作废账单累计金额",
  },
  {
    label: "总税额",
    value: formatCurrency(summary.value.totalTax),
    hint: "全部账单累计税额",
  },
  {
    label: "待收金额",
    value: formatCurrency(summary.value.totalOutstanding),
    hint: "全部未结清账单待收合计",
  },
]);

const filteredInvoices = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLowerCase();
  const normalizedTaxKeyword = filters.taxKeyword.trim().toLowerCase();
  const [dueFrom, dueTo] = filters.dueRange;

  return invoices.value.filter((item) => {
    const matchStatus = statusFilter.value === "ALL" || item.status === statusFilter.value;
    const matchCustomer =
      filters.customerId === "ALL" || item.customer?.id === filters.customerId;
    const matchType = filters.type === "ALL" || item.type === filters.type;
    const matchKeyword =
      normalizedKeyword.length === 0 ||
      [
        item.invoiceNo,
        item.customer?.name,
        item.service?.name,
        item.taxProfileName,
        item.remark,
      ]
        .some((field) => String(field ?? "").toLowerCase().includes(normalizedKeyword));
    const matchTaxKeyword =
      normalizedTaxKeyword.length === 0 ||
      String(item.taxProfileName ?? "").toLowerCase().includes(normalizedTaxKeyword);

    const totalAmountYuan = item.totalAmount / 100;
    const matchAmountMin =
      typeof filters.amountMin !== "number" || totalAmountYuan >= filters.amountMin;
    const matchAmountMax =
      typeof filters.amountMax !== "number" || totalAmountYuan <= filters.amountMax;

    const dueDate = item.dueDate ? new Date(item.dueDate) : null;
    const matchDueFrom = !dueFrom || !dueDate || dueDate >= dueFrom;
    const matchDueTo = !dueTo || !dueDate || dueDate <= dueTo;

    return (
      matchStatus &&
      matchCustomer &&
      matchType &&
      matchKeyword &&
      matchTaxKeyword &&
      matchAmountMin &&
      matchAmountMax &&
      matchDueFrom &&
      matchDueTo
    );
  });
});

function toDateValue(value?: string | Date | null) {
  if (!value) {
    return "";
  }

  return new Date(value).toISOString().slice(0, 10);
}

function unpaidAmount(row: any) {
  return Math.max(row.totalAmount - row.paidAmount, 0);
}

function resetCreateForm() {
  Object.assign(createForm, {
    customerId: customers.value[0]?.id || "",
    serviceId: "",
    taxProfileId: taxProfiles.value.find((item) => item.isDefault)?.id || "",
    type: "MANUAL",
    subtotal: 0,
    taxRate: billingSetting.value.defaultTaxRate ?? 13,
    dueDate: toDateValue(new Date(Date.now() + 7 * 24 * 60 * 60 * 1000)),
    remark: "",
    issueNow: true,
  });
}

function resetAdvancedFilters() {
  Object.assign(filters, {
    customerId: "ALL",
    type: "ALL",
    amountMin: undefined,
    amountMax: undefined,
    dueRange: [],
    taxKeyword: "",
  });
}

function handleSelectionChange(rows: any[]) {
  selectedRows.value = rows;
}

watch(
  () => createForm.customerId,
  () => {
    if (
      createForm.serviceId &&
      !customerServices.value.some((service) => service.id === createForm.serviceId)
    ) {
      createForm.serviceId = "";
    }
  },
);

watch(
  () => createForm.taxProfileId,
  (value) => {
    const profile = taxProfiles.value.find((item) => item.id === value);
    createForm.taxRate = profile?.taxRate ?? billingSetting.value.defaultTaxRate ?? 13;
  },
);

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/invoices");
    invoices.value = data.data.invoices;
    customers.value = data.data.customers;
    services.value = data.data.services;
    taxProfiles.value = data.data.taxProfiles || [];
    billingSetting.value = data.data.billingSetting || billingSetting.value;
    summary.value = data.data.summary;

    if (
      !createForm.customerId ||
      !customers.value.some((item) => item.id === createForm.customerId)
    ) {
      resetCreateForm();
    }
  } finally {
    loading.value = false;
  }
}

function openCreateDialog() {
  resetCreateForm();
  createDialogVisible.value = true;
}

function openIssueDialog(row: any) {
  actionForm.invoiceId = row.id;
  actionForm.action = "issue";
  actionForm.dueDate = toDateValue(row.dueDate);
  actionForm.reason = "";
  actionDialogVisible.value = true;
}

function openVoidDialog(row: any) {
  actionForm.invoiceId = row.id;
  actionForm.action = "void";
  actionForm.dueDate = "";
  actionForm.reason = "";
  actionDialogVisible.value = true;
}

async function submitCreate() {
  creating.value = true;
  try {
    await http.post("/invoices", createForm);
    ElMessage.success(createForm.issueNow ? "账单已创建并签发" : "草稿账单已创建");
    createDialogVisible.value = false;
    await loadData();
  } finally {
    creating.value = false;
  }
}

async function executeInvoiceAction(
  invoiceId: string,
  action: "issue" | "void",
  payload?: { dueDate?: string; reason?: string },
  options?: { silent?: boolean; reloadAfter?: boolean },
) {
  const { silent = false, reloadAfter = true } = options ?? {};

  try {
    await http.post(`/invoices/${invoiceId}/action`, {
      action,
      dueDate: payload?.dueDate || "",
      reason: payload?.reason || "",
    });

    if (!silent) {
      ElMessage.success(action === "issue" ? "账单已签发" : "账单已作废");
    }

    if (reloadAfter) {
      await loadData();
    }

    return true;
  } catch {
    return false;
  }
}

async function submitAction() {
  actionSubmitting.value = true;
  try {
    const ok = await executeInvoiceAction(actionForm.invoiceId, actionForm.action as "issue" | "void", {
      dueDate: actionForm.dueDate,
      reason: actionForm.reason,
    });

    if (ok) {
      actionDialogVisible.value = false;
    }
  } finally {
    actionSubmitting.value = false;
  }
}

async function batchIssueSelected() {
  if (!selectedDraftCount.value) {
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确认批量签发 ${selectedDraftCount.value} 张草稿账单吗？`,
      "批量签发确认",
      { type: "warning" },
    );
  } catch {
    return;
  }

  let successCount = 0;
  let failedCount = 0;
  for (const invoice of selectedRows.value.filter((item) => item.status === "DRAFT")) {
    const ok = await executeInvoiceAction(
      invoice.id,
      "issue",
      {
        dueDate: toDateValue(invoice.dueDate),
      },
      { silent: true, reloadAfter: false },
    );
    if (ok) {
      successCount += 1;
    } else {
      failedCount += 1;
    }
  }

  await loadData();
  ElMessage.success(`批量签发完成，成功 ${successCount} 张，失败 ${failedCount} 张`);
}

async function batchVoidSelected() {
  if (!selectedVoidableCount.value) {
    return;
  }

  try {
    await ElMessageBox.confirm(
      `确认批量作废 ${selectedVoidableCount.value} 张未收款账单吗？`,
      "批量作废确认",
      { type: "warning" },
    );
  } catch {
    return;
  }

  let successCount = 0;
  let failedCount = 0;
  for (const invoice of selectedRows.value.filter(
    (item) => item.status !== "VOID" && Number(item.paidAmount) === 0,
  )) {
    const ok = await executeInvoiceAction(
      invoice.id,
      "void",
      { reason: "后台批量作废" },
      { silent: true, reloadAfter: false },
    );
    if (ok) {
      successCount += 1;
    } else {
      failedCount += 1;
    }
  }

  await loadData();
  ElMessage.success(`批量作废完成，成功 ${successCount} 张，失败 ${failedCount} 张`);
}

function exportAllInvoices() {
  window.open("/api/invoices/export", "_blank");
}

async function copySelectedInvoiceNos() {
  const invoiceNos = selectedRows.value.map((item) => item.invoiceNo).join("\n");
  await navigator.clipboard.writeText(invoiceNos);
  ElMessage.success("已复制所选账单号");
}

function exportSelectedInvoices() {
  const rows = (selectedRows.value.length > 0 ? selectedRows.value : filteredInvoices.value).map(
    (item) => [
      item.invoiceNo,
      item.customer?.name || "-",
      item.service?.name || "-",
      getLabel(invoiceTypeMap, item.type),
      formatCurrency(item.subtotal),
      formatCurrency(item.taxAmount),
      formatCurrency(item.totalAmount),
      formatCurrency(item.paidAmount),
      formatCurrency(unpaidAmount(item)),
      formatDate(item.dueDate),
      getLabel(invoiceStatusMap, item.status),
    ],
  );

  downloadCsv(
    `invoices-${new Date().toISOString().slice(0, 10)}.csv`,
    [
      "账单号",
      "客户",
      "关联服务",
      "类型",
      "小计",
      "税额",
      "总金额",
      "已收金额",
      "待收金额",
      "到期日",
      "状态",
    ],
    rows,
  );
  ElMessage.success("账单数据已导出");
}

onMounted(async () => {
  await authStore.ensureUser();
  await loadData();
});
</script>

<template>
  <ListWorkbenchShell
    v-model:advancedOpen="advancedOpen"
    title="账单管理"
    subtitle="账单页已经升级为财务工作台，可在同一页面处理创建、签发、作废、筛选、批量动作和详情跳转。"
    :result-count="filteredInvoices.length"
    :selected-count="selectedRows.length"
    :show-advanced-toggle="true"
    advanced-button-text="高级筛选"
  >
    <template #actions>
      <el-button :icon="Refresh" @click="loadData">刷新数据</el-button>
      <el-button @click="exportAllInvoices">导出账单</el-button>
      <el-button v-if="canManageInvoices" type="primary" :icon="Plus" @click="openCreateDialog">
        创建账单
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
        placeholder="搜索账单号、客户名、服务名、备注、税率档案"
        style="width: 360px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>
      <el-select v-model="filters.customerId" filterable style="width: 180px">
        <el-option label="全部客户" value="ALL" />
        <el-option
          v-for="item in customers"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
    </template>

    <template #filterMeta>
      <el-tag effect="plain">当前结果 {{ filteredInvoices.length }} 条</el-tag>
      <el-tag effect="plain">待收金额 {{ formatCurrency(summary.totalOutstanding) }}</el-tag>
    </template>

    <template #advanced>
      <div class="advanced-filter-grid">
        <el-select v-model="filters.type">
          <el-option label="全部类型" value="ALL" />
          <el-option
            v-for="(label, value) in invoiceTypeMap"
            :key="value"
            :label="label"
            :value="value"
          />
        </el-select>
        <el-input
          v-model="filters.taxKeyword"
          clearable
          placeholder="税率档案关键词"
        />
        <el-input-number
          v-model="filters.amountMin"
          :precision="2"
          :min="0"
          style="width: 100%"
          placeholder="最小金额"
        />
        <el-input-number
          v-model="filters.amountMax"
          :precision="2"
          :min="0"
          style="width: 100%"
          placeholder="最大金额"
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
      <el-button plain :icon="DocumentCopy" @click="copySelectedInvoiceNos">复制账单号</el-button>
      <el-button plain :icon="Download" @click="exportSelectedInvoices">导出所选账单</el-button>
      <el-button
        plain
        type="primary"
        :disabled="!selectedDraftCount || !canManageInvoices"
        @click="batchIssueSelected"
      >
        批量签发
      </el-button>
      <el-button
        plain
        type="danger"
        :disabled="!selectedVoidableCount || !canManageInvoices"
        @click="batchVoidSelected"
      >
        批量作废
      </el-button>
    </template>

    <el-table
      v-loading="loading"
      :data="filteredInvoices"
      stripe
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="账单编号" min-width="190">
        <template #default="{ row }">
          <div class="primary-line">{{ row.invoiceNo }}</div>
          <div class="muted-line">{{ getLabel(invoiceTypeMap, row.type) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="客户 / 服务" min-width="220">
        <template #default="{ row }">
          <div>{{ row.customer.name }}</div>
          <div class="muted-line">{{ row.service?.name || "未关联服务" }}</div>
        </template>
      </el-table-column>
      <el-table-column label="账单金额" width="210">
        <template #default="{ row }">
          <div>{{ formatCurrency(row.totalAmount) }}</div>
          <div class="muted-line">
            小计 {{ formatCurrency(row.subtotal) }} / 税额 {{ formatCurrency(row.taxAmount) }}
          </div>
          <div class="muted-line">
            {{
              row.taxProfileName
                ? `${row.taxProfileName} / ${row.taxRate}%`
                : `自定义税率 / ${row.taxRate}%`
            }}
          </div>
        </template>
      </el-table-column>
      <el-table-column label="收款进度" width="180">
        <template #default="{ row }">
          <div>已收 {{ formatCurrency(row.paidAmount) }}</div>
          <div class="muted-line">待收 {{ formatCurrency(unpaidAmount(row)) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="到期日" width="120">
        <template #default="{ row }">
          {{ formatDate(row.dueDate) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="120">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)">
            {{ getLabel(invoiceStatusMap, row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="210" fixed="right">
        <template #default="{ row }">
          <div class="table-actions">
            <el-button text type="primary" @click="router.push(`/invoices/${row.id}`)">详情</el-button>
            <el-button
              v-if="canManageInvoices && row.status === 'DRAFT'"
              text
              type="primary"
              @click="openIssueDialog(row)"
            >
              签发
            </el-button>
            <el-button
              v-if="canManageInvoices && row.status !== 'VOID' && row.paidAmount === 0"
              text
              type="danger"
              @click="openVoidDialog(row)"
            >
              作废
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="createDialogVisible" title="创建账单" width="680px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="客户">
              <el-select v-model="createForm.customerId" style="width: 100%">
                <el-option
                  v-for="item in customers"
                  :key="item.id"
                  :label="`${item.name} / ${item.customerNo}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="关联服务">
              <el-select v-model="createForm.serviceId" clearable filterable style="width: 100%">
                <el-option
                  v-for="item in customerServices"
                  :key="item.id"
                  :label="`${item.name} / ${item.serviceNo}`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="账单类型">
              <el-select v-model="createForm.type" style="width: 100%">
                <el-option label="手工账单" value="MANUAL" />
                <el-option label="余额账单" value="CREDIT" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="税率档案">
              <el-select v-model="createForm.taxProfileId" clearable style="width: 100%">
                <el-option
                  v-for="item in taxProfiles"
                  :key="item.id"
                  :label="`${item.name} / ${item.taxRate}%`"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="税率">
              <el-input-number v-model="createForm.taxRate" :min="0" :max="100" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="账单小计">
              <el-input-number v-model="createForm.subtotal" :precision="2" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="到期日期">
              <el-input v-model="createForm.dueDate" placeholder="YYYY-MM-DD" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" :rows="4" />
        </el-form-item>

        <el-card class="summary-card">
          <div class="detail-meta-list">
            <div class="detail-meta-item">
              <label>账单前缀</label>
              <strong>{{ billingSetting.invoicePrefix }}</strong>
            </div>
            <div class="detail-meta-item">
              <label>预览金额</label>
              <strong>
                小计 {{ formatCurrency(createPreview.subtotal) }} / 税额
                {{ formatCurrency(createPreview.taxAmount) }} / 总额
                {{ formatCurrency(createPreview.totalAmount) }}
              </strong>
            </div>
          </div>
        </el-card>

        <el-form-item style="margin-top: 16px">
          <el-switch v-model="createForm.issueNow" inline-prompt active-text="立即签发" inactive-text="仅草稿" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="creating" @click="submitCreate">确认提交</el-button>
      </template>
    </el-dialog>

    <el-dialog
      v-model="actionDialogVisible"
      :title="actionForm.action === 'issue' ? '签发账单' : '作废账单'"
      width="520px"
    >
      <el-form label-position="top">
        <el-form-item v-if="actionForm.action === 'issue'" label="到期日期">
          <el-input v-model="actionForm.dueDate" placeholder="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item v-else label="作废原因">
          <el-input v-model="actionForm.reason" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="actionDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="actionSubmitting" @click="submitAction">确认提交</el-button>
      </template>
    </el-dialog>
  </ListWorkbenchShell>
</template>
