<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import { downloadCsv } from "@/utils/download";
import { createCustomer, fetchCustomers, type Customer } from "@/api/admin";

const router = useRouter();
const loading = ref(false);
const customers = ref<Customer[]>([]);
const total = ref(0);
const dialogVisible = ref(false);
const submitting = ref(false);
const advancedVisible = ref(false);
const selectedRows = ref<Customer[]>([]);
const activeTab = ref("ALL");

const keyword = ref("");
const status = ref("");
const identityStatus = ref("");
const customerType = ref("");
const groupKeyword = ref("");
const salesKeyword = ref("");

const form = ref({
  name: "",
  email: "",
  mobile: "",
  type: "COMPANY",
  groupName: "云业务客户",
  levelName: "标准",
  salesOwner: "",
  remarks: ""
});

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: customers.value.length },
  { key: "ACTIVE", label: "正常", count: customers.value.filter(item => item.status === "ACTIVE").length },
  { key: "DISABLED", label: "停用", count: customers.value.filter(item => item.status === "DISABLED").length },
  {
    key: "PENDING_IDENTITY",
    label: "待实名审核",
    count: customers.value.filter(item => item.identity?.verifyStatus === "PENDING").length
  },
  {
    key: "APPROVED_IDENTITY",
    label: "已实名",
    count: customers.value.filter(item => item.identity?.verifyStatus === "APPROVED").length
  }
]);

const filteredCustomers = computed(() =>
  customers.value.filter(item => {
    const normalizedKeyword = keyword.value.trim();
    const keywordMatched =
      !normalizedKeyword ||
      item.customerNo.includes(normalizedKeyword) ||
      item.name.includes(normalizedKeyword) ||
      item.email.includes(normalizedKeyword) ||
      item.mobile.includes(normalizedKeyword);

    const statusMatched = !status.value || item.status === status.value;
    const identityMatched = !identityStatus.value || item.identity?.verifyStatus === identityStatus.value;
    const typeMatched = !customerType.value || item.type === customerType.value;
    const groupMatched = !groupKeyword.value || item.groupName.includes(groupKeyword.value);
    const salesMatched = !salesKeyword.value || item.salesOwner.includes(salesKeyword.value);

    const tabMatched =
      activeTab.value === "ALL" ||
      item.status === activeTab.value ||
      (activeTab.value === "PENDING_IDENTITY" && item.identity?.verifyStatus === "PENDING") ||
      (activeTab.value === "APPROVED_IDENTITY" && item.identity?.verifyStatus === "APPROVED");

    return keywordMatched && statusMatched && identityMatched && typeMatched && groupMatched && salesMatched && tabMatched;
  })
);

const approvedCount = computed(
  () => customers.value.filter(item => item.identity?.verifyStatus === "APPROVED").length
);
const pendingCount = computed(
  () => customers.value.filter(item => item.identity?.verifyStatus === "PENDING").length
);
const activeCount = computed(() => customers.value.filter(item => item.status === "ACTIVE").length);

async function loadCustomers() {
  loading.value = true;
  try {
    const data = await fetchCustomers();
    customers.value = data.items;
    total.value = data.total;
  } finally {
    loading.value = false;
  }
}

function resetForm() {
  form.value = {
    name: "",
    email: "",
    mobile: "",
    type: "COMPANY",
    groupName: "云业务客户",
    levelName: "标准",
    salesOwner: "",
    remarks: ""
  };
}

async function handleCreate() {
  submitting.value = true;
  try {
    await createCustomer(form.value);
    ElMessage.success("客户已创建");
    dialogVisible.value = false;
    resetForm();
    await loadCustomers();
  } finally {
    submitting.value = false;
  }
}

function resetFilters() {
  activeTab.value = "ALL";
  keyword.value = "";
  status.value = "";
  identityStatus.value = "";
  customerType.value = "";
  groupKeyword.value = "";
  salesKeyword.value = "";
}

function handleSelectionChange(rows: Customer[]) {
  selectedRows.value = rows;
}

async function copySelectedCustomerNos() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择客户");
    return;
  }
  await navigator.clipboard.writeText(selectedRows.value.map(item => item.customerNo).join("\n"));
  ElMessage.success(`已复制 ${selectedRows.value.length} 个客户编号`);
}

function exportCustomers(rows: Customer[], filename: string) {
  downloadCsv(
    filename,
    ["客户编号", "客户名称", "类型", "邮箱", "手机号", "分组", "等级", "销售归属", "客户状态", "实名状态"],
    rows.map(item => [
      item.customerNo,
      item.name,
      customerTypeLabel(item.type),
      item.email,
      item.mobile,
      item.groupName,
      item.levelName,
      item.salesOwner,
      customerStatusLabel(item.status),
      identityLabel(item.identity?.verifyStatus)
    ])
  );
}

function exportCurrent() {
  exportCustomers(filteredCustomers.value, "customers-current.csv");
  ElMessage.success("已导出当前列表");
}

function exportSelected() {
  if (selectedRows.value.length === 0) {
    ElMessage.info("请先选择客户");
    return;
  }
  exportCustomers(selectedRows.value, "customers-selected.csv");
  ElMessage.success("已导出选中客户");
}

function openDetail(row: Customer) {
  void router.push(`/customer/detail/${row.id}`);
}

function customerStatusLabel(value: string) {
  return value === "ACTIVE" ? "正常" : "停用";
}

function customerStatusType(value: string) {
  return value === "ACTIVE" ? "success" : "danger";
}

function customerTypeLabel(value: string) {
  return value === "COMPANY" ? "企业客户" : "个人客户";
}

function identityLabel(value?: string) {
  switch (value) {
    case "APPROVED":
      return "已实名";
    case "PENDING":
      return "待审核";
    case "REJECTED":
      return "已驳回";
    default:
      return "未提交";
  }
}

function identityType(value?: string) {
  switch (value) {
    case "APPROVED":
      return "success";
    case "PENDING":
      return "warning";
    case "REJECTED":
      return "danger";
    default:
      return "info";
  }
}

onMounted(() => {
  void loadCustomers();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="客户"
      title="客户列表"
      subtitle="按客户状态、实名状态、客户类型和销售归属进行筛选，作为客户工作台的统一入口。"
    >
      <template #actions>
        <el-button @click="loadCustomers">刷新列表</el-button>
        <el-button type="primary" @click="dialogVisible = true">新增客户</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>客户总数</span><strong>{{ total }}</strong></div>
          <div class="summary-pill"><span>正常客户</span><strong>{{ activeCount }}</strong></div>
          <div class="summary-pill"><span>已实名</span><strong>{{ approvedCount }}</strong></div>
          <div class="summary-pill"><span>待实名审核</span><strong>{{ pendingCount }}</strong></div>
          <div class="summary-pill"><span>当前结果</span><strong>{{ filteredCustomers.length }}</strong></div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />

        <div class="filter-bar">
          <el-input v-model="keyword" placeholder="搜索客户编号、客户名称、邮箱或手机号" clearable />
          <el-select v-model="status" placeholder="客户状态" clearable>
            <el-option label="正常" value="ACTIVE" />
            <el-option label="停用" value="DISABLED" />
          </el-select>
          <el-select v-model="identityStatus" placeholder="实名状态" clearable>
            <el-option label="待审核" value="PENDING" />
            <el-option label="已实名" value="APPROVED" />
            <el-option label="已驳回" value="REJECTED" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="loadCustomers">应用筛选</el-button>
            <el-button plain @click="advancedVisible = !advancedVisible">
              {{ advancedVisible ? "收起高级筛选" : "展开高级筛选" }}
            </el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>

        <div v-if="advancedVisible" class="filter-bar filter-bar--compact">
          <el-select v-model="customerType" placeholder="客户类型" clearable>
            <el-option label="企业客户" value="COMPANY" />
            <el-option label="个人客户" value="PERSONAL" />
          </el-select>
          <el-input v-model="groupKeyword" placeholder="客户分组关键字" clearable />
          <el-input v-model="salesKeyword" placeholder="销售归属关键字" clearable />
        </div>
      </template>

      <div class="table-toolbar">
        <div class="table-toolbar__meta">
          <strong>运营列表</strong>
          <span>显示 {{ filteredCustomers.length }} 条记录</span>
          <span>已选 {{ selectedRows.length }} 条</span>
        </div>
        <div class="action-group">
          <el-button plain @click="copySelectedCustomerNos">复制客户编号</el-button>
          <el-button plain @click="exportSelected">导出选中</el-button>
          <el-button plain @click="exportCurrent">导出当前列表</el-button>
        </div>
      </div>

      <el-table
        :data="filteredCustomers"
        border
        stripe
        empty-text="暂无符合条件的客户"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="48" />
        <el-table-column prop="customerNo" label="客户编号" min-width="150" />
        <el-table-column prop="name" label="客户名称" min-width="180" />
        <el-table-column label="类型" min-width="120">
          <template #default="{ row }">{{ customerTypeLabel(row.type) }}</template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="mobile" label="手机号" min-width="140" />
        <el-table-column prop="groupName" label="客户分组" min-width="140">
          <template #default="{ row }">
            <el-tag effect="plain">{{ row.groupName || "-" }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="levelName" label="客户等级" min-width="120">
          <template #default="{ row }">
            <el-tag type="primary" effect="plain">{{ row.levelName || "-" }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="salesOwner" label="销售归属" min-width="160" />
        <el-table-column label="客户状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="customerStatusType(row.status)" effect="light">
              {{ customerStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="实名状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="identityType(row.identity?.verifyStatus)" effect="light">
              {{ identityLabel(row.identity?.verifyStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="120" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">进入工作台</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="page-table-count">共 {{ filteredCustomers.length }} 条客户记录</div>

      <el-dialog v-model="dialogVisible" title="新增客户" width="560px">
        <el-form label-position="top">
          <el-form-item label="客户名称">
            <el-input v-model="form.name" />
          </el-form-item>
          <div class="filter-bar filter-bar--compact">
            <el-form-item label="邮箱" style="flex: 1">
              <el-input v-model="form.email" />
            </el-form-item>
            <el-form-item label="手机号" style="flex: 1">
              <el-input v-model="form.mobile" />
            </el-form-item>
          </div>
          <div class="filter-bar filter-bar--compact">
            <el-form-item label="客户类型" style="flex: 1">
              <el-select v-model="form.type">
                <el-option label="企业客户" value="COMPANY" />
                <el-option label="个人客户" value="PERSONAL" />
              </el-select>
            </el-form-item>
            <el-form-item label="销售归属" style="flex: 1">
              <el-input v-model="form.salesOwner" />
            </el-form-item>
          </div>
          <div class="filter-bar filter-bar--compact">
            <el-form-item label="客户分组" style="flex: 1">
              <el-input v-model="form.groupName" />
            </el-form-item>
            <el-form-item label="客户等级" style="flex: 1">
              <el-input v-model="form.levelName" />
            </el-form-item>
          </div>
          <el-form-item label="备注">
            <el-input v-model="form.remarks" type="textarea" />
          </el-form-item>
        </el-form>
        <template #footer>
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" :loading="submitting" @click="handleCreate">保存客户</el-button>
        </template>
      </el-dialog>
    </PageWorkbench>
  </div>
</template>
