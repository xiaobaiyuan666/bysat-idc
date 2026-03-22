<script setup lang="ts">
withDefaults(
  defineProps<{
    title: string;
    subtitle?: string;
    eyebrow?: string;
  }>(),
  {
    subtitle: "",
    eyebrow: ""
  }
);
</script>

<template>
  <div class="page-card workbench-shell">
    <div class="workbench-hero">
      <div class="workbench-page__titles">
        <div v-if="eyebrow" class="workbench-page__eyebrow">{{ eyebrow }}</div>
        <h1 class="page-title">{{ title }}</h1>
        <p v-if="subtitle" class="page-subtitle">{{ subtitle }}</p>
      </div>
      <div v-if="$slots.actions" class="action-group">
        <slot name="actions" />
      </div>
    </div>

    <div v-if="$slots.context" class="workbench-page__context">
      <slot name="context" />
    </div>

    <div v-if="$slots.metrics" class="workbench-page__metrics">
      <slot name="metrics" />
    </div>

    <div v-if="$slots.filters" class="workbench-page__filters">
      <slot name="filters" />
    </div>

    <slot />
  </div>
</template>
