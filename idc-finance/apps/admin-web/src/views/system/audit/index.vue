<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import PageWorkbench from "@/components/workbench/PageWorkbench.vue";
import { fetchAuditLogs, type AuditLog } from "@/api/admin";

const loading = ref(false);
const drawerVisible = ref(false);
const items = ref<AuditLog[]>([]);
const selectedLog = ref<AuditLog | null>(null);

const filters = reactive({
  keyword: "",
  actor: "",
  action: ""
});

const filteredItems = computed(() => {
  const keyword = filters.keyword.trim().toLowerCase();
  return items.value.filter(item => {
    const keywordMatched =
      !keyword ||
      item.actor.toLowerCase().includes(keyword) ||
      item.action.toLowerCase().includes(keyword) ||
      item.target?.toLowerCase().includes(keyword) ||
      item.description.toLowerCase().includes(keyword);
    const actorMatched = !filters.actor || item.actor === filters.actor;
    const actionMatched = !filters.action || item.action === filters.action;
    return keywordMatched && actorMatched && actionMatched;
  });
});

const actorOptions = computed(() =>
  Array.from(new Set(items.value.map(item => item.actor).filter(Boolean))).sort((a, b) => a.localeCompare(b))
);

const actionOptions = computed(() =>
  Array.from(new Set(items.value.map(item => item.action).filter(Boolean))).sort((a, b) => a.localeCompare(b))
);

const summary = computed(() => ({
  total: items.value.length,
  system: items.value.filter(item => item.actor === "系统管理员").length,
  write: items.value.filter(item => item.action.includes(".create") || item.action.includes(".update")).length,
  current: filteredItems.value.length
}));

async function loadAuditLogs() {
  loading.value = true;
  try {
    items.value = await fetchAuditLogs();
  } finally {
    loading.value = false;
  }
}

function openLog(row: AuditLog) {
  selectedLog.value = row;
  drawerVisible.value = true;
}

onMounted(() => {
  void loadAuditLogs();
});
</script>

<template>
  <div v-loading="loading">
    <PageWorkbench
      eyebrow="系统 / 审计"
      title="审计日志工作台"
      subtitle="集中查看后台关键写操作，支持按操作人、动作和关键字追踪。"
    >
      <template #actions>
        <el-button type="primary" @click="loadAuditLogs">刷新日志</el-button>
      </template>

      <template #metrics>
        <div class="summary-strip">
          <div class="summary-pill"><span>日志总数</span><strong>{{ summary.total }}</strong></div>
          <div class="summary-pill"><span>系统管理员操作</span><strong>{{ summary.system }}</strong></div>
          <div class="summary-pill"><span>写操作</span><strong>{{ summary.write }}</strong></div>
          <div class="summary-pill"><span>当前筛选</span><strong>{{ summary.current }}</strong></div>
        </div>
      </template>

      <template #filters>
        <div class="filter-bar">
          <el-input v-model="filters.keyword" clearable placeholder="搜索操作人、动作、目标或说明" />
          <el-select v-model="filters.actor" clearable placeholder="操作人">
            <el-option v-for="item in actorOptions" :key="item" :label="item" :value="item" />
          </el-select>
          <el-select v-model="filters.action" clearable filterable placeholder="动作">
            <el-option v-for="item in actionOptions" :key="item" :label="item" :value="item" />
          </el-select>
        </div>
      </template>

      <div class="panel-card">
        <div class="section-card__head">
          <strong>日志列表</strong>
          <span class="section-card__meta">当前 {{ filteredItems.length }} 条</span>
        </div>
        <el-table :data="filteredItems" border stripe>
          <el-table-column prop="actor" label="操作人" min-width="140" />
          <el-table-column prop="action" label="动作" min-width="200" />
          <el-table-column prop="target" label="目标对象" min-width="180" />
          <el-table-column prop="description" label="说明" min-width="320" show-overflow-tooltip />
          <el-table-column prop="createdAt" label="时间" min-width="180" />
          <el-table-column label="操作" width="100" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="openLog(row)">详情</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </PageWorkbench>

    <el-drawer v-model="drawerVisible" title="审计日志详情" size="520px">
      <el-empty v-if="!selectedLog" description="暂无日志详情" />
      <template v-else>
        <el-descriptions :column="1" border>
          <el-descriptions-item label="操作人">{{ selectedLog.actor }}</el-descriptions-item>
          <el-descriptions-item label="动作">{{ selectedLog.action }}</el-descriptions-item>
          <el-descriptions-item label="目标对象">{{ selectedLog.target || "-" }}</el-descriptions-item>
          <el-descriptions-item label="时间">{{ selectedLog.createdAt }}</el-descriptions-item>
          <el-descriptions-item label="说明">{{ selectedLog.description }}</el-descriptions-item>
        </el-descriptions>

        <div class="panel-card" style="margin-top: 16px">
          <div class="section-card__head">
            <strong>Payload</strong>
          </div>
          <pre class="audit-payload">{{ JSON.stringify(selectedLog.payload || {}, null, 2) }}</pre>
        </div>
      </template>
    </el-drawer>
  </div>
</template>

<style scoped>
.audit-payload {
  margin: 0;
  white-space: pre-wrap;
  word-break: break-all;
  font-size: 12px;
  line-height: 1.6;
}
</style>
