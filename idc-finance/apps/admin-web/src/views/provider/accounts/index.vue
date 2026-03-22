<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage, ElMessageBox } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createProviderAccount,
  deleteProviderAccount,
  fetchProviderAccountHealth,
  fetchProviderAccounts,
  updateProviderAccount,
  type MofangHealthResponse,
  type ProviderAccount
} from "@/api/admin";

type ProviderType = "MOFANG_CLOUD" | "ZJMF_API" | "WHMCS" | "RESOURCE" | "MANUAL";

const loading = ref(false);
const saving = ref(false);
const healthLoading = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const filterType = ref<string>("");
const accounts = ref<ProviderAccount[]>([]);
const selectedHealth = ref<MofangHealthResponse | null>(null);

const providerTypeOptions: Array<{ label: string; value: ProviderType }> = [
  { label: "魔方云", value: "MOFANG_CLOUD" },
  { label: "上下游财务", value: "ZJMF_API" },
  { label: "WHMCS", value: "WHMCS" },
  { label: "资源池", value: "RESOURCE" },
  { label: "手动资源", value: "MANUAL" }
];

const form = reactive<Partial<ProviderAccount>>({
  providerType: "MOFANG_CLOUD",
  name: "",
  baseUrl: "",
  username: "",
  password: "",
  sourceName: "",
  contactWay: "",
  description: "",
  accountMode: "API",
  lang: "zh-cn",
  listPath: "/v1/clouds",
  detailPath: "/v1/clouds/:id",
  insecureSkipVerify: true,
  autoUpdate: true,
  status: "ACTIVE",
  extraConfig: ""
});

const filteredAccounts = computed(() => {
  if (!filterType.value) return accounts.value;
  return accounts.value.filter(item => item.providerType === filterType.value);
});

function resetForm() {
  editingId.value = null;
  form.providerType = "MOFANG_CLOUD";
  form.name = "";
  form.baseUrl = "";
  form.username = "";
  form.password = "";
  form.sourceName = "";
  form.contactWay = "";
  form.description = "";
  form.accountMode = "API";
  form.lang = "zh-cn";
  form.listPath = "/v1/clouds";
  form.detailPath = "/v1/clouds/:id";
  form.insecureSkipVerify = true;
  form.autoUpdate = true;
  form.status = "ACTIVE";
  form.extraConfig = "";
}

function applyProviderDefaults(type: string) {
  if (type === "MOFANG_CLOUD") {
    form.listPath = "/v1/clouds";
    form.detailPath = "/v1/clouds/:id";
    form.lang = "zh-cn";
    form.accountMode = "API";
    return;
  }
  if (type === "ZJMF_API") {
    form.listPath = "";
    form.detailPath = "";
    form.lang = "zh-cn";
    form.accountMode = "API";
  }
}

async function loadAccounts() {
  loading.value = true;
  try {
    accounts.value = await fetchProviderAccounts();
  } finally {
    loading.value = false;
  }
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
}

function openEdit(account: ProviderAccount) {
  editingId.value = account.id;
  Object.assign(form, account);
  dialogVisible.value = true;
}

async function saveAccount() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateProviderAccount(editingId.value, form);
      ElMessage.success("接口账户已更新");
    } else {
      await createProviderAccount(form);
      ElMessage.success("接口账户已创建");
    }
    dialogVisible.value = false;
    await loadAccounts();
  } finally {
    saving.value = false;
  }
}

async function removeAccount(account: ProviderAccount) {
  await ElMessageBox.confirm(`确认删除接口账户“${account.name}”吗？`, "删除确认", {
    type: "warning"
  });
  await deleteProviderAccount(account.id);
  ElMessage.success("接口账户已删除");
  await loadAccounts();
}

async function checkHealth(account: ProviderAccount) {
  healthLoading.value = true;
  try {
    selectedHealth.value = await fetchProviderAccountHealth(account.id);
    ElMessage.success(`${account.name} 已完成连接检测`);
  } finally {
    healthLoading.value = false;
  }
}

function statusType(status: string) {
  if (status === "ACTIVE") return "success";
  if (status === "DISABLED") return "warning";
  return "info";
}

onMounted(() => {
  void loadAccounts();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="资源与商店"
      title="接口账户"
      subtitle="统一管理多套魔方云、上下游财务、资源池接口，再由商品和服务绑定具体接口账户。"
    >
      <template #actions>
        <el-select v-model="filterType" clearable placeholder="筛选接口类型" style="width: 180px">
          <el-option v-for="item in providerTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
        </el-select>
        <el-button @click="loadAccounts">刷新</el-button>
        <el-button type="primary" @click="openCreate">新增接口账户</el-button>
      </template>

      <template #metrics>
        <div class="detail-kpi-grid">
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">接口总数</span>
            <strong class="detail-kpi-card__value">{{ accounts.length }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">魔方云</span>
            <strong class="detail-kpi-card__value">{{ accounts.filter(item => item.providerType === "MOFANG_CLOUD").length }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">上下游财务</span>
            <strong class="detail-kpi-card__value">{{ accounts.filter(item => item.providerType === "ZJMF_API").length }}</strong>
          </div>
          <div class="detail-kpi-card">
            <span class="detail-kpi-card__label">启用中</span>
            <strong class="detail-kpi-card__value">{{ accounts.filter(item => item.status === "ACTIVE").length }}</strong>
          </div>
        </div>
      </template>

      <div class="portal-grid portal-grid--two">
        <div class="panel-card">
          <div class="section-card__head">
            <strong>接口账户列表</strong>
            <span class="section-card__meta">商品自动化、上游映射、服务同步都应绑定到具体接口账户。</span>
          </div>
          <el-table :data="filteredAccounts" border stripe>
            <el-table-column prop="name" label="名称" min-width="180" />
            <el-table-column prop="providerType" label="类型" min-width="140" />
            <el-table-column prop="sourceName" label="来源名称" min-width="180" />
            <el-table-column prop="baseUrl" label="接口地址" min-width="240" show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" min-width="140" />
            <el-table-column prop="productCount" label="商品数" width="90" />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" effect="light">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" :loading="healthLoading" @click="checkHealth(row)">检测</el-button>
                <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
                <el-button link type="danger" @click="removeAccount(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>连接检测</strong>
            <span class="section-card__meta">先在这里确认接口可用，再到商品工作台绑定具体账户。</span>
          </div>
          <el-empty v-if="!selectedHealth" description="请选择一条接口账户执行检测" />
          <el-descriptions v-else :column="1" border>
            <el-descriptions-item label="基础地址">{{ selectedHealth.baseUrl || "-" }}</el-descriptions-item>
            <el-descriptions-item label="实际地址">{{ selectedHealth.activeUrl || "-" }}</el-descriptions-item>
            <el-descriptions-item label="认证方式">{{ selectedHealth.authMode || "-" }}</el-descriptions-item>
            <el-descriptions-item label="连接状态">
              <el-tag :type="selectedHealth.connected ? 'success' : 'danger'" effect="light">
                {{ selectedHealth.connected ? "已连接" : "异常" }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="最近认证">{{ selectedHealth.lastAuthAt || "-" }}</el-descriptions-item>
            <el-descriptions-item label="说明">{{ selectedHealth.message }}</el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
    </PageWorkbench>

    <el-dialog
      v-model="dialogVisible"
      :title="editingId ? '编辑接口账户' : '新增接口账户'"
      width="760px"
      destroy-on-close
    >
      <el-form label-position="top">
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="接口类型" style="flex: 1">
            <el-select v-model="form.providerType" @change="applyProviderDefaults">
              <el-option v-for="item in providerTypeOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
          <el-form-item label="接口名称" style="flex: 1">
            <el-input v-model="form.name" />
          </el-form-item>
        </div>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="来源名称" style="flex: 1">
            <el-input v-model="form.sourceName" />
          </el-form-item>
          <el-form-item label="状态" style="flex: 1">
            <el-select v-model="form.status">
              <el-option label="启用" value="ACTIVE" />
              <el-option label="停用" value="DISABLED" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item label="接口地址">
          <el-input v-model="form.baseUrl" placeholder="例如 https://cloud.example.com/entry/" />
        </el-form-item>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="用户名" style="flex: 1">
            <el-input v-model="form.username" />
          </el-form-item>
          <el-form-item label="密码" style="flex: 1">
            <el-input v-model="form.password" type="password" show-password />
          </el-form-item>
        </div>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="列表路径" style="flex: 1">
            <el-input v-model="form.listPath" />
          </el-form-item>
          <el-form-item label="详情路径" style="flex: 1">
            <el-input v-model="form.detailPath" />
          </el-form-item>
        </div>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="账号模式" style="flex: 1">
            <el-input v-model="form.accountMode" />
          </el-form-item>
          <el-form-item label="语言" style="flex: 1">
            <el-input v-model="form.lang" />
          </el-form-item>
        </div>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="联系信息" style="flex: 1">
            <el-input v-model="form.contactWay" />
          </el-form-item>
          <el-form-item label="自动更新" style="width: 160px">
            <el-switch v-model="form.autoUpdate" />
          </el-form-item>
          <el-form-item label="忽略证书" style="width: 160px">
            <el-switch v-model="form.insecureSkipVerify" />
          </el-form-item>
        </div>
        <el-form-item label="备注">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="扩展配置">
          <el-input v-model="form.extraConfig" type="textarea" :rows="3" placeholder="需要时可放 JSON 或附加说明" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveAccount">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
