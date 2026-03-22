<script setup lang="ts">
import { Delete, Edit, Plus } from "@element-plus/icons-vue";
import { ElMessage, ElMessageBox } from "element-plus";
import { reactive, ref } from "vue";

import { http } from "@/api/http";
import { customerTypeMap, getLabel, portalRoleMap } from "@/utils/maps";

const props = defineProps<{
  loading: boolean;
  portalUsers: any[];
  customers: any[];
}>();

const emit = defineEmits<{
  updated: [];
}>();

const dialogVisible = ref(false);
const editingId = ref("");
const form = reactive({
  customerId: "",
  name: "",
  email: "",
  phone: "",
  role: "OWNER",
  password: "",
  isOwner: true,
  isActive: true,
});

function resetForm() {
  Object.assign(form, {
    customerId: props.customers[0]?.id ?? "",
    name: "",
    email: "",
    phone: "",
    role: "OWNER",
    password: "",
    isOwner: true,
    isActive: true,
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
    customerId: row.customerId,
    name: row.name,
    email: row.email,
    phone: row.phone || "",
    role: row.role,
    password: "",
    isOwner: row.isOwner,
    isActive: row.isActive,
  });
  dialogVisible.value = true;
}

async function submitForm() {
  const payload = {
    kind: "portalUser",
    ...form,
    ...(editingId.value ? { id: editingId.value } : {}),
  };

  if (editingId.value) {
    await http.put("/architecture", payload);
    ElMessage.success("门户账号已更新");
  } else {
    await http.post("/architecture", payload);
    ElMessage.success("门户账号已创建");
  }

  dialogVisible.value = false;
  emit("updated");
}

async function removeItem(row: any) {
  await ElMessageBox.confirm(`确认删除门户账号 ${row.email} 吗？`, "删除确认", {
    type: "warning",
  });

  await http.delete("/architecture", {
    data: {
      kind: "portalUser",
      id: row.id,
    },
  });

  ElMessage.success("门户账号已删除");
  emit("updated");
}

resetForm();
</script>

<template>
  <el-card class="page-card">
    <template #header>
      <div class="card-header">
        <div>
          <div class="card-title">门户联系人</div>
          <div class="card-subtitle">客户前台登录账号、多联系人和角色体系。</div>
        </div>
        <el-button type="primary" :icon="Plus" @click="openCreate">新增账号</el-button>
      </div>
    </template>

    <el-table v-loading="loading" :data="portalUsers" stripe>
      <el-table-column label="联系人" min-width="180">
        <template #default="{ row }">
          <div style="font-weight: 700">{{ row.name }}</div>
          <div class="muted-line">{{ row.email }}</div>
        </template>
      </el-table-column>
      <el-table-column label="所属客户" min-width="180">
        <template #default="{ row }">
          <div>{{ row.customer?.name }}</div>
          <div class="muted-line">{{ getLabel(customerTypeMap, row.customer?.type) }}</div>
        </template>
      </el-table-column>
      <el-table-column label="角色" width="120">
        <template #default="{ row }">
          {{ getLabel(portalRoleMap, row.role) }}
        </template>
      </el-table-column>
      <el-table-column label="主账号" width="90">
        <template #default="{ row }">
          <el-tag :type="row.isOwner ? 'success' : 'info'">{{ row.isOwner ? "是" : "否" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="状态" width="90">
        <template #default="{ row }">
          <el-tag :type="row.isActive ? 'success' : 'info'">{{ row.isActive ? "启用" : "停用" }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="150" fixed="right">
        <template #default="{ row }">
          <div class="row-actions">
            <el-button text type="primary" :icon="Edit" @click="openEdit(row)">编辑</el-button>
            <el-button text type="danger" :icon="Delete" @click="removeItem(row)">删除</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="dialogVisible" :title="editingId ? '编辑门户账号' : '新增门户账号'" width="620px">
      <el-form label-position="top">
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="所属客户">
              <el-select v-model="form.customerId" style="width: 100%">
                <el-option
                  v-for="item in customers"
                  :key="item.id"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="门户角色">
              <el-select v-model="form.role" style="width: 100%">
                <el-option label="主账号" value="OWNER" />
                <el-option label="技术联系人" value="TECH" />
                <el-option label="财务联系人" value="BILLING" />
                <el-option label="只读成员" value="READONLY" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="联系人姓名">
              <el-input v-model="form.name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="手机号">
              <el-input v-model="form.phone" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-row :gutter="12">
          <el-col :span="12">
            <el-form-item label="邮箱">
              <el-input v-model="form.email" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item :label="editingId ? '重置密码（留空则不修改）' : '登录密码'">
              <el-input v-model="form.password" type="password" show-password />
            </el-form-item>
          </el-col>
        </el-row>
        <div class="switch-row">
          <el-switch v-model="form.isOwner" active-text="设为主账号" inactive-text="普通账号" />
          <el-switch v-model="form.isActive" active-text="启用" inactive-text="停用" />
        </div>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitForm">{{ editingId ? "保存修改" : "确认创建" }}</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<style scoped>
.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.card-title {
  font-size: 18px;
  font-weight: 700;
}

.card-subtitle,
.muted-line {
  margin-top: 6px;
  color: var(--text-secondary);
  font-size: 12px;
}

.row-actions {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
}

.switch-row {
  display: flex;
  flex-wrap: wrap;
  gap: 20px;
}
</style>
