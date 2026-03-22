<script setup lang="ts">
import { onMounted, ref } from "vue";

import { http } from "@/api/http";
import { formatDateTime } from "@/utils/format";
import {
  auditActionMap,
  auditModuleMap,
  auditTargetTypeMap,
  getLabel,
  roleMap,
} from "@/utils/maps";

const loading = ref(false);
const audits = ref<any[]>([]);

async function loadData() {
  loading.value = true;
  try {
    const { data } = await http.get("/audits");
    audits.value = data.data;
  } finally {
    loading.value = false;
  }
}

onMounted(loadData);
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">审计日志</h1>
        <div class="page-subtitle">
          记录财务、服务、工单、客户和配置操作，便于合规审计与责任追踪。
        </div>
      </div>
      <el-button @click="loadData">刷新数据</el-button>
    </div>

    <el-card class="page-card">
      <el-table v-loading="loading" :data="audits" stripe>
        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatDateTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作人" min-width="180">
          <template #default="{ row }">
            <div>{{ row.adminUser?.name || "系统" }}</div>
            <div style="margin-top: 6px; color: var(--text-secondary); font-size: 12px">
              {{ getLabel(roleMap, row.adminUser?.role) }}
            </div>
          </template>
        </el-table-column>
        <el-table-column label="模块" width="120">
          <template #default="{ row }">
            {{ getLabel(auditModuleMap, row.module) }}
          </template>
        </el-table-column>
        <el-table-column label="动作" width="140">
          <template #default="{ row }">
            {{ getLabel(auditActionMap, row.action) }}
          </template>
        </el-table-column>
        <el-table-column label="对象" width="140">
          <template #default="{ row }">
            {{ getLabel(auditTargetTypeMap, row.targetType) }}
          </template>
        </el-table-column>
        <el-table-column prop="summary" label="摘要" min-width="280" />
      </el-table>
    </el-card>
  </div>
</template>
