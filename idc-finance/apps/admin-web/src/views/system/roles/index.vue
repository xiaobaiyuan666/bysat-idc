<script setup lang="ts">
import { computed, nextTick, onMounted, reactive, ref } from "vue";
import { ElMessage } from "element-plus";
import type { MenuItem } from "@/store/modules/permission";
import { usePermissionStore } from "@/store";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import {
  createRole,
  fetchRoles,
  updateRole,
  type RoleRecord,
  type SaveRoleRequest
} from "@/api/admin";

const permissionStore = usePermissionStore();

const loading = ref(false);
const saving = ref(false);
const dialogVisible = ref(false);
const editingId = ref<number | null>(null);
const treeRef = ref<any>();
const roles = ref<RoleRecord[]>([]);
const selectedRole = ref<RoleRecord | null>(null);

const filters = reactive({
  keyword: "",
  status: ""
});

const form = reactive<SaveRoleRequest>({
  name: "",
  code: "",
  description: "",
  status: "ACTIVE",
  menuIds: [],
  permissions: []
});

const flatMenus = computed(() => flattenMenus(permissionStore.menus));
const permissionOptions = computed(() => permissionStore.permissions);

const filteredRoles = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();
  return roles.value.filter(item => {
    const keywordMatched =
      !keyword ||
      item.name.toLowerCase().includes(keyword) ||
      item.code.toLowerCase().includes(keyword) ||
      item.description.toLowerCase().includes(keyword);
    const statusMatched = !filters.status || item.status === filters.status;
    return keywordMatched && statusMatched;
  });
});

const summary = computed(() => ({
  total: roles.value.length,
  active: roles.value.filter(item => item.status === "ACTIVE").length,
  disabled: roles.value.filter(item => item.status !== "ACTIVE").length,
  permissions: roles.value.reduce((sum, item) => sum + item.permissions.length, 0)
}));

const selectedMenuItems = computed(() => {
  if (!selectedRole.value) return [];
  const menuIdSet = new Set(selectedRole.value.menuIds);
  return flatMenus.value.filter(item => menuIdSet.has(item.id));
});

function flattenMenus(nodes: MenuItem[]) {
  const results: MenuItem[] = [];
  const walk = (items: MenuItem[]) => {
    items.forEach(item => {
      results.push(item);
      if (item.children?.length) {
        walk(item.children);
      }
    });
  };
  walk(nodes);
  return results;
}

function statusLabel(value: string) {
  return value === "ACTIVE" ? "启用" : "停用";
}

function statusType(value: string) {
  return value === "ACTIVE" ? "success" : "info";
}

function resetForm() {
  editingId.value = null;
  form.name = "";
  form.code = "";
  form.description = "";
  form.status = "ACTIVE";
  form.menuIds = [];
  form.permissions = [];
}

function selectRole(row?: RoleRecord | null) {
  selectedRole.value = row ?? null;
}

function checkedMenuIds() {
  const checkedKeys = treeRef.value?.getCheckedKeys?.(false) ?? [];
  const halfCheckedKeys = treeRef.value?.getHalfCheckedKeys?.() ?? [];
  return Array.from(new Set([...checkedKeys, ...halfCheckedKeys])).map(item => Number(item));
}

async function syncTreeCheckedKeys() {
  await nextTick();
  treeRef.value?.setCheckedKeys?.(form.menuIds);
}

function openCreate() {
  resetForm();
  dialogVisible.value = true;
  void syncTreeCheckedKeys();
}

function openEdit(row: RoleRecord) {
  editingId.value = row.id;
  form.name = row.name;
  form.code = row.code;
  form.description = row.description || "";
  form.status = row.status;
  form.menuIds = [...row.menuIds];
  form.permissions = [...row.permissions];
  dialogVisible.value = true;
  void syncTreeCheckedKeys();
}

async function loadData() {
  loading.value = true;
  try {
    if (!permissionStore.menus.length || !permissionStore.permissions.length) {
      await permissionStore.load();
    }
    const roleData = await fetchRoles();
    roles.value = roleData;
    if (!selectedRole.value || !roleData.some(item => item.id === selectedRole.value?.id)) {
      selectedRole.value = roleData[0] ?? null;
    } else {
      selectedRole.value = roleData.find(item => item.id === selectedRole.value?.id) ?? null;
    }
  } finally {
    loading.value = false;
  }
}

async function saveRole() {
  saving.value = true;
  try {
    const payload: SaveRoleRequest = {
      ...form,
      menuIds: checkedMenuIds(),
      permissions: [...form.permissions]
    };
    if (editingId.value) {
      await updateRole(editingId.value, payload);
      ElMessage.success("角色已更新");
    } else {
      await createRole(payload);
      ElMessage.success("角色已创建");
    }
    dialogVisible.value = false;
    await loadData();
  } catch (error: any) {
    ElMessage.error(error?.message ?? "角色保存失败");
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  void loadData();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="系统 / 角色权限"
      title="角色权限工作台"
      subtitle="集中维护后台角色、菜单覆盖、权限点集合和绑定人数。"
    >
      <template #actions>
        <el-button @click="loadData">刷新</el-button>
        <el-button type="primary" @click="openCreate">新增角色</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>角色数量</span><strong>{{ summary.total }}</strong></div>
          <div class="summary-pill"><span>启用角色</span><strong>{{ summary.active }}</strong></div>
          <div class="summary-pill"><span>停用角色</span><strong>{{ summary.disabled }}</strong></div>
          <div class="summary-pill"><span>权限总数</span><strong>{{ summary.permissions }}</strong></div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-input v-model="filters.keyword" clearable placeholder="搜索角色名称、编码或描述" />
          <el-select v-model="filters.status" clearable placeholder="状态">
            <el-option label="启用" value="ACTIVE" />
            <el-option label="停用" value="DISABLED" />
          </el-select>
        </div>
      </template>

      <div class="system-grid">
        <div class="panel-card">
          <div class="section-card__head">
            <strong>角色列表</strong>
            <span class="section-card__meta">当前 {{ filteredRoles.length }} 个角色</span>
          </div>
          <el-table :data="filteredRoles" border stripe highlight-current-row @current-change="selectRole">
            <el-table-column prop="name" label="角色名称" min-width="160" />
            <el-table-column prop="code" label="角色编码" min-width="160" />
            <el-table-column label="状态" width="90">
              <template #default="{ row }">
                <el-tag :type="statusType(row.status)" effect="light">{{ statusLabel(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="users" label="绑定人数" width="90" />
            <el-table-column label="菜单数" width="90">
              <template #default="{ row }">{{ row.menuIds.length }}</template>
            </el-table-column>
            <el-table-column label="权限数" width="90">
              <template #default="{ row }">{{ row.permissions.length }}</template>
            </el-table-column>
            <el-table-column prop="updatedAt" label="更新时间" min-width="160" />
            <el-table-column label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="openEdit(row)">编辑</el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>{{ selectedRole ? `${selectedRole.name} 工作台` : "角色详情" }}</strong>
          </div>
          <el-empty v-if="!selectedRole" description="请选择一个角色" :image-size="80" />
          <template v-else>
            <div class="summary-strip">
              <div class="summary-pill"><span>绑定人数</span><strong>{{ selectedRole.users }}</strong></div>
              <div class="summary-pill"><span>菜单节点</span><strong>{{ selectedRole.menuIds.length }}</strong></div>
              <div class="summary-pill"><span>权限点</span><strong>{{ selectedRole.permissions.length }}</strong></div>
            </div>

            <el-descriptions :column="2" border style="margin-top: 16px">
              <el-descriptions-item label="角色编码">{{ selectedRole.code }}</el-descriptions-item>
              <el-descriptions-item label="状态">{{ statusLabel(selectedRole.status) }}</el-descriptions-item>
              <el-descriptions-item label="创建时间">{{ selectedRole.createdAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="更新时间">{{ selectedRole.updatedAt || "-" }}</el-descriptions-item>
              <el-descriptions-item label="说明" :span="2">
                {{ selectedRole.description || "暂无角色说明" }}
              </el-descriptions-item>
            </el-descriptions>

            <div class="system-stack">
              <div>
                <div class="section-card__head">
                  <strong>菜单覆盖</strong>
                  <span class="section-card__meta">已选 {{ selectedMenuItems.length }} 个菜单节点</span>
                </div>
                <div class="tag-wrap">
                  <el-tag v-for="item in selectedMenuItems.slice(0, 16)" :key="item.id" effect="plain">
                    {{ item.title }}
                  </el-tag>
                </div>
              </div>

              <div>
                <div class="section-card__head">
                  <strong>权限点</strong>
                </div>
                <div class="tag-wrap">
                  <el-tag v-for="item in selectedRole.permissions.slice(0, 20)" :key="item" effect="plain">
                    {{ item }}
                  </el-tag>
                </div>
              </div>
            </div>
          </template>
        </div>
      </div>
    </PageWorkbench>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑角色' : '新增角色'" width="980px" destroy-on-close>
      <div class="role-dialog">
        <div class="panel-card">
          <el-form label-position="top">
            <div class="filter-bar filter-bar--compact">
              <el-form-item label="角色名称" style="flex: 1">
                <el-input v-model="form.name" />
              </el-form-item>
              <el-form-item label="角色编码" style="flex: 1">
                <el-input v-model="form.code" :disabled="Boolean(editingId)" />
              </el-form-item>
            </div>
            <div class="filter-bar filter-bar--compact">
              <el-form-item label="状态" style="flex: 1">
                <el-select v-model="form.status">
                  <el-option label="启用" value="ACTIVE" />
                  <el-option label="停用" value="DISABLED" />
                </el-select>
              </el-form-item>
            </div>
            <el-form-item label="角色说明">
              <el-input v-model="form.description" type="textarea" :rows="4" />
            </el-form-item>
            <el-form-item label="额外权限点">
              <el-select v-model="form.permissions" multiple collapse-tags collapse-tags-tooltip>
                <el-option v-for="item in permissionOptions" :key="item" :label="item" :value="item" />
              </el-select>
            </el-form-item>
          </el-form>
        </div>

        <div class="panel-card">
          <div class="section-card__head">
            <strong>菜单授权</strong>
            <span class="section-card__meta">勾选菜单会自动附带其权限点</span>
          </div>
          <el-tree
            ref="treeRef"
            :data="permissionStore.menus"
            node-key="id"
            default-expand-all
            show-checkbox
            :props="{ label: 'title', children: 'children' }"
          />
        </div>
      </div>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveRole">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style scoped>
.system-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.08fr) minmax(360px, 0.92fr);
  gap: 16px;
}

.system-stack {
  display: grid;
  gap: 16px;
  margin-top: 16px;
}

.role-dialog {
  display: grid;
  grid-template-columns: minmax(320px, 0.95fr) minmax(320px, 1.05fr);
  gap: 16px;
}

.tag-wrap {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

@media (max-width: 1280px) {
  .system-grid,
  .role-dialog {
    grid-template-columns: 1fr;
  }
}
</style>
