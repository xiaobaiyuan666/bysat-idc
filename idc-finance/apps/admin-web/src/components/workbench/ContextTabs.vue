<script setup lang="ts">
import { useRouter } from "vue-router";

interface ContextTabItem {
  key: string;
  label: string;
  to?: string;
  active?: boolean;
  badge?: string | number;
  disabled?: boolean;
}

const props = defineProps<{
  items: ContextTabItem[];
}>();

const router = useRouter();

function handleClick(item: ContextTabItem) {
  if (item.disabled || !item.to || item.active) return;
  void router.push(item.to);
}
</script>

<template>
  <div class="context-tabs">
    <button
      v-for="item in props.items"
      :key="item.key"
      type="button"
      class="context-tab"
      :class="{
        'context-tab--active': item.active,
        'context-tab--disabled': item.disabled
      }"
      @click="handleClick(item)"
    >
      <span>{{ item.label }}</span>
      <span v-if="item.badge !== undefined" class="context-tab__badge">{{ item.badge }}</span>
    </button>
  </div>
</template>
