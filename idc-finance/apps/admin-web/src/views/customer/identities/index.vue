<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
import { ElMessage, ElMessageBox } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import StatusTabs from "@/components/workbench/StatusTabs.vue";
import {
  fetchCustomerIdentities,
  reviewCustomerIdentity,
  type CustomerIdentityOverview
} from "@/api/admin";

type IdentityTab = "ALL" | "PENDING" | "APPROVED" | "REJECTED";

const router = useRouter();
const loading = ref(false);
const items = ref<CustomerIdentityOverview[]>([]);
const activeTab = ref<IdentityTab>("ALL");

const filters = reactive({
  keyword: "",
  identityType: ""
});

const statusTabs = computed(() => [
  { key: "ALL", label: "全部", count: items.value.length },
  { key: "PENDING", label: "待审核", count: items.value.filter(item => item.verifyStatus === "PENDING").length },
  {
    key: "APPROVED",
    label: "已通过",
    count: items.value.filter(item => item.verifyStatus === "APPROVED").length
  },
  {
    key: "REJECTED",
    label: "已驳回",
    count: items.value.filter(item => item.verifyStatus === "REJECTED").length
  }
]);

const identityTypes = computed(() =>
  Array.from(new Set(items.value.map(item => item.identityType).filter(Boolean)))
);

const filteredItems = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();
  return items.value.filter(item => {
    const tabMatched = activeTab.value === "ALL" || item.verifyStatus === activeTab.value;
    const typeMatched = !filters.identityType || item.identityType === filters.identityType;
    const keywordMatched =
      !keyword ||
      item.customerNo.toLowerCase().includes(keyword) ||
      item.customerName.toLowerCase().includes(keyword) ||
      item.subjectName.toLowerCase().includes(keyword) ||
      item.certNo.toLowerCase().includes(keyword);
    return tabMatched && typeMatched && keywordMatched;
  });
});

async function loadIdentities() {
  loading.value = true;
  try {
    items.value = await fetchCustomerIdentities();
  } finally {
    loading.value = false;
  }
}

function resetFilters() {
  activeTab.value = "ALL";
  filters.keyword = "";
  filters.identityType = "";
}

function statusLabel(value: CustomerIdentityOverview["verifyStatus"]) {
  return value === "APPROVED" ? "已通过" : value === "PENDING" ? "待审核" : "已驳回";
}

function statusType(value: CustomerIdentityOverview["verifyStatus"]) {
  return value === "APPROVED" ? "success" : value === "PENDING" ? "warning" : "danger";
}

function customerTypeLabel(value: string) {
  return value === "COMPANY" ? "企业客户" : "个人客户";
}

async function handleReview(item: CustomerIdentityOverview, status: "APPROVED" | "REJECTED") {
  const approved = status === "APPROVED";
  const result = await ElMessageBox.prompt(
    approved ? "填写审核备注，可留空后直接提交。" : "请输入驳回原因。",
    approved ? "通过实名认证" : "驳回实名认证",
    {
      inputValue: approved ? "资料齐全，审核通过" : "资料不完整，请补充后重新提交",
      inputPlaceholder: "请输入备注",
      confirmButtonText: "提交",
      cancelButtonText: "取消"
    }
  );

  await reviewCustomerIdentity(item.customerId, item.id, {
    status,
    remark: result.value || undefined
  });
  ElMessage.success(approved ? "实名认证已通过" : "实名认证已驳回");
  await loadIdentities();
}

function openCustomer(item: CustomerIdentityOverview) {
  void router.push(`/customer/detail/${item.customerId}`);
}

function openPendingCustomers() {
  void router.push({
    path: "/customer/list",
    query: {
      tab: "PENDING_IDENTITY"
    }
  });
}

onMounted(() => {
  void loadIdentities();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="客户"
      title="实名认证"
      subtitle="集中处理实名档案、审核备注和客户跳转，让审核人员可以直接围绕客户完成处理。"
    >
      <template #actions>
        <el-button @click="openPendingCustomers">待实名客户</el-button>
        <el-button @click="loadIdentities">刷新</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>总记录</span><strong>{{ items.length }}</strong></div>
          <div class="summary-pill">
            <span>待审核</span>
            <strong>{{ items.filter(item => item.verifyStatus === "PENDING").length }}</strong>
          </div>
          <div class="summary-pill">
            <span>已通过</span>
            <strong>{{ items.filter(item => item.verifyStatus === "APPROVED").length }}</strong>
          </div>
          <div class="summary-pill">
            <span>已驳回</span>
            <strong>{{ items.filter(item => item.verifyStatus === "REJECTED").length }}</strong>
          </div>
        </div>
      </template>

      <template #filters>
        <StatusTabs v-model="activeTab" :items="statusTabs" />
        <div class="filter-bar">
          <el-input
            v-model="filters.keyword"
            placeholder="搜索客户编号、客户名称、认证主体或证件号码"
            clearable
          />
          <el-select v-model="filters.identityType" placeholder="认证类型" clearable>
            <el-option v-for="item in identityTypes" :key="item" :label="item" :value="item" />
          </el-select>
          <div class="action-group">
            <el-button type="primary" @click="loadIdentities">刷新数据</el-button>
            <el-button plain @click="resetFilters">重置</el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredItems" border stripe empty-text="暂无实名认证记录">
        <el-table-column prop="customerNo" label="客户编号" min-width="160" />
        <el-table-column prop="customerName" label="客户名称" min-width="180">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomer(row)">{{ row.customerName }}</el-button>
          </template>
        </el-table-column>
        <el-table-column label="客户类型" min-width="120">
          <template #default="{ row }">{{ customerTypeLabel(row.customerType) }}</template>
        </el-table-column>
        <el-table-column prop="identityType" label="认证类型" min-width="120" />
        <el-table-column prop="subjectName" label="认证主体" min-width="220" />
        <el-table-column prop="certNo" label="证件号码" min-width="220" />
        <el-table-column prop="countryCode" label="国家/地区" min-width="120">
          <template #default="{ row }">
            {{ row.countryCode || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="审核状态" min-width="120">
          <template #default="{ row }">
            <el-tag :type="statusType(row.verifyStatus)" effect="light">
              {{ statusLabel(row.verifyStatus) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submittedAt" label="提交时间" min-width="180">
          <template #default="{ row }">
            {{ row.submittedAt || "-" }}
          </template>
        </el-table-column>
        <el-table-column prop="reviewedAt" label="审核时间" min-width="180">
          <template #default="{ row }">
            {{ row.reviewedAt || "-" }}
          </template>
        </el-table-column>
        <el-table-column prop="reviewRemark" label="审核备注" min-width="220">
          <template #default="{ row }">
            {{ row.reviewRemark || "-" }}
          </template>
        </el-table-column>
        <el-table-column label="操作" min-width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openCustomer(row)">客户工作台</el-button>
            <el-button
              link
              type="success"
              :disabled="row.verifyStatus === 'APPROVED'"
              @click="handleReview(row, 'APPROVED')"
            >
              通过
            </el-button>
            <el-button link type="danger" @click="handleReview(row, 'REJECTED')">驳回</el-button>
          </template>
        </el-table-column>
      </el-table>
    </PageWorkbench>
  </div>
</template>
