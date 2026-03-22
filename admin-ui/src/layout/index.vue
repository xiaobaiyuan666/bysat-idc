<script setup lang="ts">
import { onMounted, watch } from "vue";
import { useRouter } from "vue-router";

import LayContent from "./components/lay-content/index.vue";
import LayNavbar from "./components/lay-navbar/index.vue";
import NavVertical from "./components/lay-sidebar/NavVertical.vue";
import LayTag from "./components/lay-tag/index.vue";
import { useAppStoreHook } from "@/stores/modules/app";
import { useAuthStore } from "@/stores/auth";
import { usePermissionStoreHook } from "@/stores/modules/permission";

const router = useRouter();
const appStore = useAppStoreHook();
const authStore = useAuthStore();
const permissionStore = usePermissionStoreHook();

function syncMenus() {
  permissionStore.syncMenus(router.getRoutes(), authStore.user?.permissions || []);
}

onMounted(syncMenus);

watch(
  () => authStore.user?.permissions,
  () => {
    syncMenus();
  },
  { deep: true },
);
</script>

<template>
  <div class="pure-layout">
    <NavVertical />

    <div class="pure-layout__main" :class="{ collapsed: !appStore.sidebar.opened }">
      <LayNavbar />
      <LayTag v-if="!appStore.hideTabs" />
      <LayContent />
    </div>
  </div>
</template>

<style scoped>
.pure-layout {
  min-height: 100vh;
  display: flex;
}

.pure-layout__main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}
</style>
