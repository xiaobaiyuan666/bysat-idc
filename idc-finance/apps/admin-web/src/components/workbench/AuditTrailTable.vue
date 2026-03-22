<script setup lang="ts">
import { computed } from "vue";
import type { AuditLog } from "@/api/admin";
import { useLocaleStore } from "@/store";
import {
  fieldLabel,
  formatAuditAction,
  formatAuditDescription,
  formatAuditReason,
  formatFieldValue
} from "@/utils/business";

const props = withDefaults(defineProps<{ items: AuditLog[]; emptyText?: string }>(), {
  emptyText: "暂无变更记录"
});

const localeStore = useLocaleStore();

const copy = computed(() => ({
  createdAt: localeStore.locale === "en-US" ? "Time" : "时间",
  actor: localeStore.locale === "en-US" ? "Operator" : "操作人",
  action: localeStore.locale === "en-US" ? "Action" : "动作",
  description: localeStore.locale === "en-US" ? "Description" : "说明",
  reason: localeStore.locale === "en-US" ? "Reason" : "变更原因",
  summary: localeStore.locale === "en-US" ? "Change Summary" : "变更摘要"
}));

function asRecord(value: unknown): Record<string, unknown> {
  if (value && typeof value === "object" && !Array.isArray(value)) {
    return value as Record<string, unknown>;
  }
  return {};
}

function changeReason(item: AuditLog) {
  return formatAuditReason(asRecord(item.payload).reason);
}

function changeSummary(item: AuditLog) {
  const payload = asRecord(item.payload);
  const before = asRecord(payload.before);
  const after = asRecord(payload.after);
  const keys = Array.from(new Set([...Object.keys(before), ...Object.keys(after)]));
  return keys
    .filter(
      key =>
        formatFieldValue(localeStore.locale, key, before[key], item.action) !==
        formatFieldValue(localeStore.locale, key, after[key], item.action)
    )
    .map(key => ({
      key,
      label: fieldLabel(localeStore.locale, key),
      before: formatFieldValue(localeStore.locale, key, before[key], item.action),
      after: formatFieldValue(localeStore.locale, key, after[key], item.action)
    }));
}
</script>

<template>
  <el-table :data="props.items" border stripe :empty-text="props.emptyText">
    <el-table-column prop="createdAt" :label="copy.createdAt" min-width="170" />
    <el-table-column prop="actor" :label="copy.actor" min-width="120" />
    <el-table-column :label="copy.action" min-width="180">
      <template #default="{ row }">
        {{ formatAuditAction(localeStore.locale, (row as AuditLog).action) }}
      </template>
    </el-table-column>
    <el-table-column :label="copy.description" min-width="220">
      <template #default="{ row }">
        {{ formatAuditDescription(localeStore.locale, (row as AuditLog).action, (row as AuditLog).description) }}
      </template>
    </el-table-column>
    <el-table-column :label="copy.reason" min-width="220">
      <template #default="{ row }">
        <span>{{ changeReason(row as AuditLog) }}</span>
      </template>
    </el-table-column>
    <el-table-column :label="copy.summary" min-width="360">
      <template #default="{ row }">
        <div class="audit-change-list">
          <div v-for="item in changeSummary(row as AuditLog)" :key="item.key" class="audit-change-item">
            <span class="audit-change-item__label">{{ item.label }}</span>
            <span class="audit-change-item__before">{{ item.before }}</span>
            <span class="audit-change-item__arrow">-></span>
            <span class="audit-change-item__after">{{ item.after }}</span>
          </div>
          <span v-if="changeSummary(row as AuditLog).length === 0" class="audit-change-item__empty">-</span>
        </div>
      </template>
    </el-table-column>
  </el-table>
</template>
