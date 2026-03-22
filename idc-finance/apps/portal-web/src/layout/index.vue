<template>
  <div class="portal-layout">
    <NavVertical />
    <div class="portal-main">
      <LayNavbar />
      <LayTag />
      <LayContent />
    </div>
  </div>
</template>

<script setup lang="ts">
import NavVertical from "./components/lay-sidebar/NavVertical.vue";
import LayNavbar from "./components/lay-navbar/index.vue";
import LayTag from "./components/lay-tag/index.vue";
import LayContent from "./components/lay-content/index.vue";
</script>
