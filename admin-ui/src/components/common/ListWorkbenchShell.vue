<script setup lang="ts">
const props = withDefaults(
  defineProps<{
    title: string;
    subtitle: string;
    resultCount?: number;
    selectedCount?: number;
    advancedOpen?: boolean;
    showAdvancedToggle?: boolean;
    advancedButtonText?: string;
  }>(),
  {
    resultCount: 0,
    selectedCount: 0,
    advancedOpen: false,
    showAdvancedToggle: false,
    advancedButtonText: "高级筛选",
  },
);

const emit = defineEmits<{
  "update:advancedOpen": [value: boolean];
}>();

function toggleAdvanced() {
  emit("update:advancedOpen", !props.advancedOpen);
}
</script>

<template>
  <div class="page-body">
    <div class="toolbar-row">
      <div>
        <h1 class="page-title">{{ title }}</h1>
        <div class="page-subtitle">{{ subtitle }}</div>
      </div>
      <div class="toolbar-actions">
        <slot name="actions" />
      </div>
    </div>

    <div class="summary-strip">
      <slot name="summary" />
    </div>

    <el-card class="page-card">
      <div class="filter-bar">
        <div class="filter-group">
          <slot name="filters" />
        </div>

        <div class="filter-group">
          <slot name="filterMeta">
            <el-tag effect="plain">当前结果 {{ resultCount }} 条</el-tag>
          </slot>
          <el-button v-if="showAdvancedToggle" plain @click="toggleAdvanced">
            {{ advancedOpen ? "收起筛选" : advancedButtonText }}
          </el-button>
        </div>
      </div>

      <div v-if="advancedOpen" class="advanced-filter-panel">
        <slot name="advanced" />
      </div>

      <div v-if="selectedCount > 0" class="selection-toolbar">
        <div class="selection-toolbar__meta">
          已选择 <strong>{{ selectedCount }}</strong> 项
        </div>
        <div class="selection-toolbar__actions">
          <slot name="selection" />
        </div>
      </div>

      <div class="list-table-shell">
        <slot />
      </div>
    </el-card>
  </div>
</template>
