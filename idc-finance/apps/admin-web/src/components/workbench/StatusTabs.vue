<script setup lang="ts">
interface StatusTabItem {
  key: string;
  label: string;
  count?: number;
}

const activeKey = defineModel<string>({ required: true });

const props = defineProps<{
  items: StatusTabItem[];
}>();

function handleClick(item: StatusTabItem) {
  activeKey.value = item.key;
}
</script>

<template>
  <div class="status-tabs">
    <button
      v-for="item in props.items"
      :key="item.key"
      type="button"
      class="status-tab"
      :class="{ 'status-tab--active': activeKey === item.key }"
      @click="handleClick(item)"
    >
      <span>{{ item.label }}</span>
      <span v-if="item.count !== undefined" class="status-tab__badge">{{ item.count }}</span>
    </button>
  </div>
</template>
