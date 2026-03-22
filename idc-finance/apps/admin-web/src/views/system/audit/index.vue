<script setup lang="ts">
import { onMounted, ref } from "vue";
import { fetchAuditLogs, type AuditLog } from "@/api/admin";

const loading = ref(false);
const items = ref<AuditLog[]>([]);

async function loadAuditLogs() {
  loading.value = true;
  try {
    items.value = await fetchAuditLogs();
  } finally {
    loading.value = false;
  }
}

onMounted(loadAuditLogs);
</script>

<template>
  <div class="page-card" v-loading="loading">
    <div class="page-header">
      <div>
        <h1 class="page-title">审计日志</h1>
        <p class="page-subtitle">所有后台关键写操作都会沉淀到这里，便于回溯责任和排查异常。</p>
      </div>
      <el-button type="primary" @click="loadAuditLogs">刷新日志</el-button>
    </div>

    <el-table :data="items" border>
      <el-table-column prop="actor" label="操作人" min-width="140" />
      <el-table-column prop="action" label="动作" min-width="180" />
      <el-table-column prop="target" label="目标对象" min-width="180" />
      <el-table-column prop="description" label="说明" min-width="280" />
      <el-table-column prop="createdAt" label="时间" min-width="180" />
    </el-table>
  </div>
</template>
