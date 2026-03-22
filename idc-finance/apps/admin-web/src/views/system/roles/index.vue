<script setup lang="ts">
import { computed, onMounted, ref } from "vue";
import { fetchRoles } from "@/api/admin";

type RoleRow = {
  id: number;
  name: string;
  code: string;
  users: number;
};

const items = ref<RoleRow[]>([]);
const loading = ref(false);

const totalRoles = computed(() => items.value.length);
const totalBindings = computed(() => items.value.reduce((sum, item) => sum + item.users, 0));
const maxBindings = computed(() => (items.value.length ? Math.max(...items.value.map(item => item.users)) : 0));

onMounted(async () => {
  loading.value = true;
  try {
    items.value = await fetchRoles();
  } finally {
    loading.value = false;
  }
});
</script>

<template>
  <div class="page-card" v-loading="loading">
    <div class="page-header">
      <div>
        <h1 class="page-title">角色权限</h1>
        <p class="page-subtitle">
          这里先把后台角色主体和绑定用户数列清，后续再扩展到菜单编辑器和按钮权限。
        </p>
      </div>
    </div>

    <div class="summary-strip">
      <div class="summary-pill"><span>角色数量</span><strong>{{ totalRoles }}</strong></div>
      <div class="summary-pill"><span>绑定用户总数</span><strong>{{ totalBindings }}</strong></div>
      <div class="summary-pill"><span>最高绑定数</span><strong>{{ maxBindings }}</strong></div>
    </div>

    <el-table :data="items" border>
      <el-table-column prop="name" label="角色名称" min-width="180" />
      <el-table-column prop="code" label="角色编码" min-width="220" />
      <el-table-column prop="users" label="绑定用户数" min-width="140" />
    </el-table>
  </div>
</template>
