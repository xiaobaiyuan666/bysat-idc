<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createAdmin,
  fetchAdmins,
  fetchRoles,
  updateAdmin,
  type AdminMember,
  type RoleRecord,
  type SaveAdminRequest
} from "@/api/admin";

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const admins = ref<AdminMember[]>([]);
const roles = ref<RoleRecord[]>([]);
const selectedAdmin = ref<AdminMember | null>(null);

const filters = reactive({
  keyword: "",
  status: "",
  roleId: ""
});

const form = reactive<SaveAdminRequest>({
  username: "",
  displayName: "",
  email: "",
  mobile: "",
  status: "ACTIVE",
  roleIds: [],
  remarks: ""
});

const roleOptions = computed(() =>
  roles.value
    .filter(item => item.status === "ACTIVE")
    .map(item => ({
      label: `${item.name} (${item.code})`,
      value: item.id
    }))
);

const filteredAdmins = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();
  return admins.value.filter(item => {
    const keywordMatched =
      !keyword ||
      item.username.toLowerCase().includes(keyword) ||
      item.displayName.toLowerCase().includes(keyword) ||
      item.email.toLowerCase().includes(keyword) ||
      item.mobile.toLowerCase().includes(keyword) ||
      item.roleNames.join(" ").toLowerCase().includes(keyword);
    const statusMatched = !filters.status || item.status === filters.status;
    const roleMatched = !filters.roleId || item.roleIds.includes(Number(filters.roleId));
    return keywordMatched && statusMatched && roleMatched;
  });
});

const summary = computed(() => ({
  total: admins.value.length,
  active: admins.value.filter(item => item.status === "ACTIVE").length,
  disabled: admins.value.filter(item => item.status !== "ACTIVE").length,
  superAdmins: admins.value.filter(item => item.roles.includes("super-admin")).length
}));

const selectedRoleRecords = computed(() => {
  if (!selectedAdmin.value) return [];
  const roleIds = new Set(selectedAdmin.value.roleIds);
  return roles.value.filter(item => roleIds.has(item.id));
});

function statusLabel(value: string) {
  return value === "ACTIVE" ? "正常" : "停用";
}

function statusType(value: string) {
  return value === "ACTIVE" ? "success" : "danger";
}

function resetForm() {
  editingId.value = null;
  form.username = "";
  form.displayName = "";
  form.email = "";
  form.mobile = "";
  form.status = "ACTIVE";
  form.roleIds = [];
  form.remarks = "";
}

function openCreate() {
  resetForm();
  if (roleOptions.value.length > 0) {
    form.roleIds = [roleOptions.value[0].value];
  }
  dialogVisible.value = true;
}

function openEdit(row: AdminMember) {
  editingId.value = row.id;
  form.username = row.username;
  form.displayName = row.displayName;
  form.email = row.email || "";
  form.mobile = row.mobile || "";
  form.status = row.status;
  form.roleIds = [...row.roleIds];
  form.remarks = row.remarks || "";
  dialogVisible.value = true;
}

function selectAdmin(row?: AdminMember | null) {
  selectedAdmin.value = row ?? null;
}

async function loadData() {
  loading.value = true;
  try {
    const [adminData, roleData] = await Promise.all([fetchAdmins(), fetchRoles()]);
    admins.value = adminData;
    roles.value = roleData;
    if (!selectedAdmin.value || !adminData.some(item => item.id === selectedAdmin.value?.id)) {
      selectedAdmin.value = adminData[0] ?? null;
    } else {
      selectedAdmin.value = adminData.find(item => item.id === selectedAdmin.value?.id) ?? null;
    }
  } finally {
    loading.value = false;
  }
}

async function saveAdmin() {
  saving.value = true;
  try {
    if (editingId.value) {
      await updateAdmin(editingId.value, form);
      ElMessage.success("管理员已更新");
    } else {
      await createAdmin(form);
      ElMessage.success("管理员已创建");
    }
    dialogVisible.value = false;
    await loadData();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "管理员保存失败");
  } finally {
    saving.value = false;
  }
}

async function toggleStatus(row: AdminMember) {
  try {
    await updateAdmin(row.id, {
      username: row.username,
      displayName: row.displayName,
      email: row.email,
      mobile: row.mobile,
      status: row.status === "ACTIVE" ? "DISABLED" : "ACTIVE",
      roleIds: row.roleIds,
      remarks: row.remarks
    });
    ElMessage.success(row.status === "ACTIVE" ? "管理员已停用" : "管理员已启用");
    await loadData();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "管理员状态更新失败");
  }
}

onMounted(() => {
  void loadData();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="系统 / 管理员"
      title="管理员工作台"
      subtitle="集中维护后台账号、角色绑定、启停状态和最近登录信息。"
    >
      <template #actions>
        <el-button @click="loadData">刷新</el-button>
        <el-button type="primary" @click="openCreate">新增管理员</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>管理员总数</span><strong>{{ summary.total }}</strong></div>
          <div class="summary-pill"><span>正常账号</span><strong>{{ summary.active }}</strong></div>
          <div class="summary-pill"><span>停用账号</span><strong>{{ summary.disabled }}</strong></div>
          <div class="summary-pill"><span>超级管理员</span><strong>{{ summary.superAdmins }}</strong></div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-input v-model="filters.keyword" clearable placeholder="搜索账号、姓名、邮箱、手机号或角色" />
          <el-select v-model="filters.status" clearable placeholder="状态">
            <el-option label="正常" value="ACTIVE" />
            <el-option label="停用" value="DISABLED" />
          </el-select>
          <el-select v-model="filters.roleId" clearable placeholder="角色">
            <el-option v-for="item in roles" :key="item.id" :label="item.name" :value="String(item.id)" />
          </el-select>
        </div>
      </template>

      <div class="system-grid">
        <div class="panel-card">
          <div class="section-card__head">
            <strong>管理员列表</strong>
            <span class="section-card__meta">当前 {{ filteredAdmins.length }} 个账号</span>
          </div>
          <el-table :data="filteredAdmins" border stripe highlight-current-row @current-change="selectAdmin">
            <el-table-column prop="username" label="登录账号" min-width="140" />
            <el-table-column prop="displayName" label="显示名称" min-width="140" />
            <el-table-column label="角色" min-width="220">
              <template #default="{ row }">
                <div class="tag-wrap">
                  <el-tag v-for="item in row.roleNames" :key="item" effect="plain">{{ item }}</el-tag>
                </div>
              </template>
            </el-table-column>
            <el-table-column prop="mobile" label="手机号" min-width="140" />
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="lastLoginAt" label="最近登录" min-width="160" />
            <el-table-column label="操作" min-width="170" fixed="right">
              <template #default="{ row }">
                <div class="inline-actions">
                  <el-button link type="primary" @click="openEdit(row)">编辑</el-button>
                  <el-button link :type="row.status === 'ACTIVE' ? 'danger' : 'success'" @click="toggleStatus(row)">
                    {{ row.status === "ACTIVE" ? "停用" : "启用" }}
                  </el-button>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>{{ selectedAdmin ? `${selectedAdmin.displayName} 工作台` : "账号详情" }}</strong>
          </div>
          <el-empty v-if="!selectedAdmin" description="请选择一个管理员账号" :image-size="80" />
          <template v-else>
            <div class="summary-strip">
              <div class="summary-pill"><span>绑定角色</span><strong>{{ selectedAdmin.roleIds.length }}</strong></div>
              <div class="summary-pill"><span>状态</span><strong>{{ statusLabel(selectedAdmin.status) }}</strong></div>
              <div class="summary-pill"><span>最近登录 IP</span><strong>{{ selectedAdmin.lastLoginIp || "-" }}</strong></div>
            </div>

            <el-descriptions :column="2" border style="margin-top: 16px">
              <el-descriptions-item label="登录账号">{{ selectedAdmin.username }}</el-descriptions-item>
              <el-descriptions-item label="显示名称">{{ selectedAdmin.displayName }}</el-descriptions-item>
              <el-descriptions-item label="邮箱">{{ selectedAdmin.email || "-" }}</el-descriptions-item>
              <el-descriptions-item label="手机号">{{ selectedAdmin.mobile || "-" }}</el-descriptions-item>
              <el-descriptions-item label="最近登录">{{ selectedAdmin.lastLoginAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="登录 IP">{{ selectedAdmin.lastLoginIp || "-" }}</el-descriptions-item>
              <el-descriptions-item label="备注" :span="2">{{ selectedAdmin.remarks || "暂无备注" }}</el-descriptions-item>
            </el-descriptions>

            <div class="system-stack">
              <div>
                <div class="section-card__head">
                  <strong>角色清单</strong>
                </div>
                <el-table :data="selectedRoleRecords" border stripe size="small" empty-text="当前账号暂无角色">
                  <el-table-column prop="name" label="角色名称" min-width="160" />
                  <el-table-column prop="code" label="角色编码" min-width="160" />
                  <el-table-column prop="users" label="绑定人数" width="90" />
                  <el-table-column label="权限数" width="90">
                    <template #default="{ row }">{{ row.permissions.length }}</template>
                  </el-table-column>
                </el-table>
              </div>

              <el-alert
                title="管理员状态变化会写入审计日志，角色权限由角色工作台统一维护。"
                type="info"
                :closable="false"
                show-icon
              />
            </div>
          </template>
        </div>
      </div>
    </PageWorkbench>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑管理员' : '新增管理员'" width="720px" destroy-on-close>
      <el-form label-position="top">
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="登录账号" style="flex: 1">
            <el-input v-model="form.username" :disabled="Boolean(editingId)" />
          </el-form-item>
          <el-form-item label="显示名称" style="flex: 1">
            <el-input v-model="form.displayName" />
          </el-form-item>
        </div>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="邮箱" style="flex: 1">
            <el-input v-model="form.email" />
          </el-form-item>
          <el-form-item label="手机号" style="flex: 1">
            <el-input v-model="form.mobile" />
          </el-form-item>
        </div>
        <div class="filter-bar filter-bar--compact">
          <el-form-item label="账号状态" style="flex: 1">
            <el-select v-model="form.status">
              <el-option label="正常" value="ACTIVE" />
              <el-option label="停用" value="DISABLED" />
            </el-select>
          </el-form-item>
          <el-form-item label="绑定角色" style="flex: 1">
            <el-select v-model="form.roleIds" multiple collapse-tags collapse-tags-tooltip>
              <el-option v-for="item in roleOptions" :key="item.value" :label="item.label" :value="item.value" />
            </el-select>
          </el-form-item>
        </div>
        <el-form-item label="备注">
          <el-input v-model="form.remarks" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveAdmin">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.system-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.15fr) minmax(360px, 0.85fr);
  gap: 16px;
}

.system-stack {
  display: grid;
  gap: 16px;
  margin-top: 16px;
}

.tag-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 1280px) {
  .system-grid {
    grid-template-columns: 1fr;
  }
}
</style>
