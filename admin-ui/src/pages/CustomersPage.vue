<script setup lang="ts">
import {
  DocumentCopy,
  Edit,
  Money,
  Plus,
  Refresh,
  Search,
  Download,
} from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";

import ListWorkbenchShell from "@/components/common/ListWorkbenchShell.vue";
import { http } from "@/api/http";
import { downloadCsv } from "@/utils/download";
import { formatCurrency, formatDate, formatDateTime } from "@/utils/format";
import {
  creditTransactionTypeMap,
  customerStatusMap,
  customerTypeMap,
  getLabel,
  getStatusTagType,
  roleMap,
} from "@/utils/maps";

const router = useRouter();

const loading = ref(false);
const dialogVisible = ref(false);
const balanceVisible = ref(false);
const ledgerVisible = ref(false);
const advancedOpen = ref(false);
const editingId = ref("");
const statusFilter = ref("ALL");
const keyword = ref("");
const customers = ref<any[]>([]);
const selectedRows = ref<any[]>([]);
const targetCustomer = ref<any | null>(null);
const ledger = ref<any | null>(null);

const filters = reactive({
  type: "ALL",
  levelKeyword: "",
  risk: "ALL",
  minBalance: undefined as number | undefined,
  maxBalance: undefined as number | undefined,
  createdRange: [] as [Date, Date] | [],
});

const form = reactive({
  name: "",
  companyName: "",
  email: "",
  phone: "",
  contactQQ: "",
  contactWechat: "",
  type: "COMPANY",
  status: "ACTIVE",
  level: "standard",
  tags: "",
  notes: "",
});

const balanceForm = reactive({
  amount: 0,
  description: "",
});

const statusCards = computed(() => {
  const total = customers.value.length;
  const active = customers.value.filter((item) => item.status === "ACTIVE").length;
  const suspended = customers.value.filter((item) => item.status === "SUSPENDED").length;
  const overdue = customers.value.filter((item) => item.status === "OVERDUE").length;

  return [
    { key: "ALL", label: "全部客户", value: total, hint: "系统内全部客户对象" },
    { key: "ACTIVE", label: "正常客户", value: active, hint: "可正常下单与续费" },
    { key: "OVERDUE", label: "逾期客户", value: overdue, hint: "需要财务继续跟进" },
    { key: "SUSPENDED", label: "暂停客户", value: suspended, hint: "账户受限或已停用" },
  ];
});

const summaryTotals = computed(() => {
  const totalBilled = customers.value.reduce((sum, item) => sum + item.totalBilled, 0);
  const totalBalance = customers.value.reduce((sum, item) => sum + item.creditBalance, 0);

  return [
    {
      label: "累计开票",
      value: formatCurrency(totalBilled),
      hint: "全部客户历史账单金额",
    },
    {
      label: "客户余额池",
      value: formatCurrency(totalBalance),
      hint: "可用于自动续费和手工扣费",
    },
  ];
});

const filteredCustomers = computed(() => {
  const normalizedKeyword = keyword.value.trim().toLowerCase();
  const normalizedLevelKeyword = filters.levelKeyword.trim().toLowerCase();
  const [createdFrom, createdTo] = filters.createdRange;

  return customers.value.filter((item) => {
    const matchStatus = statusFilter.value === "ALL" || item.status === statusFilter.value;
    const matchType = filters.type === "ALL" || item.type === filters.type;
    const matchKeyword =
      normalizedKeyword.length === 0 ||
      [
        item.name,
        item.customerNo,
        item.companyName,
        item.email,
        item.phone,
        item.tags,
      ]
        .some((field) => String(field ?? "").toLowerCase().includes(normalizedKeyword));
    const matchLevel =
      normalizedLevelKeyword.length === 0 ||
      String(item.level ?? "").toLowerCase().includes(normalizedLevelKeyword);

    const balanceYuan = item.creditBalance / 100;
    const matchMinBalance =
      typeof filters.minBalance !== "number" || balanceYuan >= filters.minBalance;
    const matchMaxBalance =
      typeof filters.maxBalance !== "number" || balanceYuan <= filters.maxBalance;

    const createdAt = new Date(item.createdAt);
    const matchCreatedFrom = !createdFrom || createdAt >= createdFrom;
    const matchCreatedTo = !createdTo || createdAt <= createdTo;

    let matchRisk = true;
    if (filters.risk === "NEGATIVE_BALANCE") {
      matchRisk = item.creditBalance < 0;
    } else if (filters.risk === "OPEN_RECEIVABLE") {
      matchRisk = item.totalBilled > item.totalPaid;
    } else if (filters.risk === "NO_SERVICE") {
      matchRisk = item._count.services === 0;
    }

    return (
      matchStatus &&
      matchType &&
      matchKeyword &&
      matchLevel &&
      matchMinBalance &&
      matchMaxBalance &&
      matchCreatedFrom &&
      matchCreatedTo &&
      matchRisk
    );
  });
});

const openReceivableCount = computed(
  () => filteredCustomers.value.filter((item) => item.totalBilled > item.totalPaid).length,
);

function resetForm() {
  Object.assign(form, {
    name: "",
    companyName: "",
    email: "",
    phone: "",
    contactQQ: "",
    contactWechat: "",
    type: "COMPANY",
    status: "ACTIVE",
    level: "standard",
    tags: "",
    notes: "",
  });
}

function resetAdvancedFilters() {
  Object.assign(filters, {
    type: "ALL",
    levelKeyword: "",
    risk: "ALL",
    minBalance: undefined,
    maxBalance: undefined,
    createdRange: [],
  });
}

function openCreate() {
  editingId.value = "";
  resetForm();
  dialogVisible.value = true;
}

function openEdit(row: any) {
  editingId.value = row.id;
  Object.assign(form, {
    name: row.name,
    companyName: row.companyName || "",
    email: row.email,
    phone: row.phone || "",
    contactQQ: row.contactQQ || "",
    contactWechat: row.contactWechat || "",
    type: row.type,
    status: row.status,
    level: row.level,
    tags: row.tags || "",
    notes: row.notes || "",
  });
  dialogVisible.value = true;
}

function handleSelectionChange(rows: any[]) {
  selectedRows.value = rows;
}

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/customers");
    customers.value = data.data;
  } finally {
    loading.value = false;
  }
}

async function submitForm() {
  if (editingId.value) {
    await http.put(`/customers/${editingId.value}`, form);
    ElMessage.success("客户已更新");
  } else {
    await http.post("/customers", form);
    ElMessage.success("客户已创建");
  }

  dialogVisible.value = false;
  resetForm();
  await loadData();
}

function openBalanceDialog(row: any) {
  targetCustomer.value = row;
  balanceForm.amount = 0;
  balanceForm.description = "";
  balanceVisible.value = true;
}

async function submitBalance() {
  if (!targetCustomer.value) {
    return;
  }

  await http.post(`/customers/${targetCustomer.value.id}/balance`, balanceForm);
  ElMessage.success("客户余额已更新");
  balanceVisible.value = false;
  await loadData();
}

async function openLedger(row: any) {
  targetCustomer.value = row;
  const { data } = await http.get(`/customers/${row.id}/ledger`);
  ledger.value = data.data;
  ledgerVisible.value = true;
}

async function copySelectedCustomerNos() {
  const customerNos = selectedRows.value.map((item) => item.customerNo).join("\n");
  await navigator.clipboard.writeText(customerNos);
  ElMessage.success("已复制所选客户编号");
}

function exportCustomers() {
  const rows = (selectedRows.value.length > 0 ? selectedRows.value : filteredCustomers.value).map(
    (item) => [
      item.customerNo,
      item.name,
      getLabel(customerTypeMap, item.type),
      getLabel(customerStatusMap, item.status),
      item.email,
      item.phone || "-",
      formatCurrency(item.creditBalance),
      formatCurrency(item.totalBilled),
      formatCurrency(item.totalPaid),
      item.level || "-",
      formatDate(item.createdAt),
    ],
  );

  downloadCsv(
    `customers-${new Date().toISOString().slice(0, 10)}.csv`,
    ["客户编号", "客户名称", "类型", "状态", "邮箱", "电话", "余额", "累计开票", "累计实收", "等级", "创建日期"],
    rows,
  );
  ElMessage.success("客户数据已导出");
}

onMounted(loadData);
</script>

<template>
  <ListWorkbenchShell
    v-model:advancedOpen="advancedOpen"
    title="客户列表"
    subtitle="客户列表已经升级成完整的客户工作台入口，可从这里进入客户详情，继续查看订单、服务、账单、流水、工单和审计记录。"
    :result-count="filteredCustomers.length"
    :selected-count="selectedRows.length"
    :show-advanced-toggle="true"
    advanced-button-text="高级筛选"
  >
    <template #actions>
      <el-button :icon="Refresh" @click="loadData">刷新数据</el-button>
      <el-button type="primary" :icon="Plus" @click="openCreate">新增客户</el-button>
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
        placeholder="搜索客户名称、编号、邮箱、电话、标签"
        style="width: 360px"
      >
        <template #prefix>
          <el-icon><Search /></el-icon>
        </template>
      </el-input>

      <el-select v-model="filters.type" style="width: 160px">
        <el-option label="全部类型" value="ALL" />
        <el-option
          v-for="(label, value) in customerTypeMap"
          :key="value"
          :label="label"
          :value="value"
        />
      </el-select>
    </template>

    <template #filterMeta>
      <el-tag effect="plain">当前结果 {{ filteredCustomers.length }} 条</el-tag>
      <el-tag effect="plain">待跟进应收 {{ openReceivableCount }} 家</el-tag>
    </template>

    <template #advanced>
      <div class="advanced-filter-grid">
        <el-input
          v-model="filters.levelKeyword"
          clearable
          placeholder="等级关键词，例如 KA / 代理"
        />
        <el-select v-model="filters.risk">
          <el-option label="全部风险" value="ALL" />
          <el-option label="负余额客户" value="NEGATIVE_BALANCE" />
          <el-option label="有未结清应收" value="OPEN_RECEIVABLE" />
          <el-option label="暂无服务" value="NO_SERVICE" />
        </el-select>
        <el-input-number
          v-model="filters.minBalance"
          :precision="2"
          :min="-9999999"
          style="width: 100%"
          placeholder="最小余额"
        />
        <el-input-number
          v-model="filters.maxBalance"
          :precision="2"
          :min="-9999999"
          style="width: 100%"
          placeholder="最大余额"
        />
        <el-date-picker
          v-model="filters.createdRange"
          type="daterange"
          range-separator="至"
          start-placeholder="创建开始"
          end-placeholder="创建结束"
          style="width: 100%"
        />
        <div class="inline-actions">
          <el-button @click="resetAdvancedFilters">清空高级筛选</el-button>
        </div>
      </div>
    </template>

    <template #selection>
      <el-button plain :icon="DocumentCopy" @click="copySelectedCustomerNos">复制客户编号</el-button>
      <el-button plain :icon="Download" @click="exportCustomers">导出所选客户</el-button>
    </template>

    <el-table
      v-loading="loading"
      :data="filteredCustomers"
      stripe
      row-key="id"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" />
      <el-table-column label="客户信息" min-width="260">
        <template #default="{ row }">
          <div class="primary-line">{{ row.name }}</div>
          <div class="muted-line">{{ row.companyName || row.email }}</div>
          <div class="muted-line">{{ row.customerNo }} / 创建于 {{ formatDate(row.createdAt) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="类型 / 等级" width="150">
        <template #default="{ row }">
          <div>{{ getLabel(customerTypeMap, row.type) }}</div>
          <div class="muted-line">{{ row.level }}</div>
        </template>
      </el-table-column>
      <el-table-column label="联系方式" min-width="180">
        <template #default="{ row }">
          <div>{{ row.email }}</div>
          <div class="muted-line">{{ row.phone || "-" }}</div>
        </template>
      </el-table-column>
      <el-table-column label="订单 / 服务 / 工单" width="160">
        <template #default="{ row }">
          {{ row._count.orders }} / {{ row._count.services }} / {{ row._count.tickets }}
        </template>
      </el-table-column>
      <el-table-column label="账户余额" width="130">
        <template #default="{ row }">
          {{ formatCurrency(row.creditBalance) }}
        </template>
      </el-table-column>
      <el-table-column label="累计开票 / 实收" width="170">
        <template #default="{ row }">
          <div>{{ formatCurrency(row.totalBilled) }}</div>
          <div class="muted-line">{{ formatCurrency(row.totalPaid) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="110">
        <template #default="{ row }">
          <el-tag :type="getStatusTagType(row.status)">
            {{ getLabel(customerStatusMap, row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" min-width="260" fixed="right">
        <template #default="{ row }">
          <div class="table-actions">
            <el-button text type="primary" @click="router.push(`/customers/${row.id}`)">详情</el-button>
            <el-button text type="primary" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button text @click="openLedger(row)">流水</el-button>
            <el-button text type="warning" :icon="Money" @click="openBalanceDialog(row)">余额</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑客户' : '新增客户'" width="720px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="客户名称">
              <el-input v-model="form.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="企业名称">
              <el-input v-model="form.companyName" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="联系邮箱">
              <el-input v-model="form.email" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="联系电话">
              <el-input v-model="form.phone" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="QQ">
              <el-input v-model="form.contactQQ" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="微信">
              <el-input v-model="form.contactWechat" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="8">
            <el-form-item label="客户类型">
              <el-select v-model="form.type" style="width: 100%">
                <el-option label="企业" value="COMPANY" />
                <el-option label="个人" value="PERSONAL" />
                <el-option label="代理商" value="RESELLER" />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="账户状态">
              <el-select v-model="form.status" style="width: 100%">
                <el-option
                  v-for="(label, value) in customerStatusMap"
                  :key="value"
                  :label="label"
                  :value="value"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="客户等级">
              <el-input v-model="form.level" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="标签">
          <el-input v-model="form.tags" placeholder="例如 KA、代理商、重点客户" />
        </el-form-item>
        <el-form-item label="备注说明">
          <el-input v-model="form.notes" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">
          {{ editingId ? "保存修改" : "确认创建" }}
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="balanceVisible" title="调整客户余额" width="520px">
      <div v-if="targetCustomer" class="muted-line" style="margin-bottom: 12px">
        {{ targetCustomer.name }} / 当前余额 {{ formatCurrency(targetCustomer.creditBalance) }}
      </div>
      <el-form label-position="top">
        <el-form-item label="调整金额">
          <el-input-number v-model="balanceForm.amount" :precision="2" style="width: 100%" />
          <div class="muted-line">正数表示充值，负数表示扣减，操作会写入余额流水和审计日志。</div>
        </el-form-item>
        <el-form-item label="调整说明">
          <el-input v-model="balanceForm.description" type="textarea" :rows="4" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="balanceVisible = false">取消</el-button>
        <el-button type="primary" @click="submitBalance">确认提交</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="ledgerVisible" title="客户余额流水" size="680px">
      <div v-if="ledger" class="muted-line" style="margin-bottom: 12px">
        {{ ledger.name }} / 当前余额 {{ formatCurrency(ledger.creditBalance) }}
      </div>
      <el-table v-if="ledger" :data="ledger.creditTransactions" stripe>
        <el-table-column label="时间" width="170">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="类型" width="120">
          <template #default="{ row }">
            {{ getLabel(creditTransactionTypeMap, row.type) }}
          </template>
        </el-table-column>
        <el-table-column label="金额" width="130">
          <template #default="{ row }">
            {{ formatCurrency(row.amount) }}
          </template>
        </el-table-column>
        <el-table-column label="变更后余额" width="140">
          <template #default="{ row }">
            {{ formatCurrency(row.balanceAfter) }}
          </template>
        </el-table-column>
        <el-table-column label="操作人" width="140">
          <template #default="{ row }">
            {{ row.operator?.name || "系统" }}
            <div class="muted-line">{{ getLabel(roleMap, row.operator?.role) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="说明" min-width="220" />
      </el-table>
    </el-drawer>
  </ListWorkbenchShell>
</template>
