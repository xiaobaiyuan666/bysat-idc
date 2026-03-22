<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { fetchAdmins } from "@/api/admin";

type AdminRow = {
  id: number;
  username: string;
  displayName: string;
  status: string;
  roles: string[];
};

const items = ref<AdminRow[]>([]);
const loading = ref(false);

const totalAdmins = computed(() => items.value.length);
const enabledCount = computed(() => items.value.filter(item => item.status === "ACTIVE").length);
const superAdminCount = computed(() =>
  items.value.filter(item => item.roles.includes("super-admin")).length
);

function statusLabel(value: string) {
  return value === "ACTIVE" ? "正常" : "停用";
}

function statusType(value: string) {
  return value === "ACTIVE" ? "success" : "danger";
}

onMounted(async () => {
  loading.value = true;
  try {
    items.value = await fetchAdmins();
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="page-card" v-loading="loading">
    <div class="page-header">
      <div>
        <h1 class="page-title">管理员</h1>
        <p class="page-subtitle">
          管理员账号、状态和角色绑定的基础页，后续可继续扩展操作审计和权限分配。
        </p>
      </div>
    </div>

    <div class="summary-strip">
      <div class="summary-pill"><span>管理员总数</span><strong>{{ totalAdmins }}</strong></div>
      <div class="summary-pill"><span>正常账号</span><strong>{{ enabledCount }}</strong></div>
      <div class="summary-pill"><span>超级管理员</span><strong>{{ superAdminCount }}</strong></div>
    </div>

    <el-table :data="items" border>
      <el-table-column prop="username" label="登录账号" min-width="160" />
      <el-table-column prop="displayName" label="显示名称" min-width="160" />
      <el-table-column label="状态" min-width="120">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" effect="light">
            {{ statusLabel(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="roles" label="角色" min-width="260">
        <template #default="{ row }">
          <div style="display: flex; gap: 8px; flex-wrap: wrap">
            <el-tag v-for="role in row.roles" :key="role" effect="plain">{{ role }}</el-tag>
          </div>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>
