<script setup lang="ts">
import { RouterView } from "vue-router";
</script>

<template>
  <main class="pure-content">
    <RouterView />
  </main>
</template>

<style scoped>
.pure-content {
  min-height: 0;
  padding: 20px 24px 24px;
  background:
    radial-gradient(circle at top right, rgba(47, 128, 255, 0.09), transparent 30%),
    radial-gradient(circle at top left, rgba(53, 198, 167, 0.08), transparent 24%),
    var(--bg-page);
}
</style>
