<script setup lang="ts">
import { computed } from "vue";
import { useRoute } from "vue-router";

import SidebarItem from "./components/SidebarItem.vue";
import SidebarLogo from "./components/SidebarLogo.vue";
import { useAppStoreHook } from "@/stores/modules/app";
import { usePermissionStoreHook } from "@/stores/modules/permission";

const route = useRoute();
const appStore = useAppStoreHook();
const permissionStore = usePermissionStoreHook();

const defaultActive = computed(() =>
  String(route.meta.activeMenu ?? route.path),
);
</script>

<template>
  <aside class="pure-sidebar" :class="{ collapsed: !appStore.sidebar.opened }">
    <SidebarLogo :collapse="!appStore.sidebar.opened" />

    <el-scrollbar class="pure-sidebar__scroll">
      <el-menu
        router
        unique-opened
        mode="vertical"
        class="pure-sidebar__menu"
        :default-active="defaultActive"
        :collapse="!appStore.sidebar.opened"
        :collapse-transition="false"
      >
        <SidebarItem
          v-for="item in permissionStore.wholeMenus"
          :key="item.index"
          :item="item"
        />
      </el-menu>
    </el-scrollbar>
  </aside>
</template>

<style scoped>
.pure-sidebar {
  width: 286px;
  height: 100vh;
  overflow: hidden;
  background: linear-gradient(180deg, var(--menu-bg) 0%, var(--menu-bg-2) 100%);
  border-right: 1px solid rgba(255, 255, 255, 0.06);
  box-shadow: 24px 0 48px rgba(9, 18, 34, 0.18);
  transition: width 0.2s ease;
}

.pure-sidebar.collapsed {
  width: 84px;
}

.pure-sidebar__scroll {
  height: calc(100vh - 84px);
}

:deep(.pure-sidebar__menu) {
  background: transparent;
  border-right: none;
  padding: 10px 10px 18px;
}

:deep(.pure-sidebar__menu .el-menu-item),
:deep(.pure-sidebar__menu .el-sub-menu__title) {
  color: rgba(255, 255, 255, 0.74);
  border-radius: 14px;
  margin-bottom: 6px;
}

:deep(.pure-sidebar__menu .el-menu-item:hover),
:deep(.pure-sidebar__menu .el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.06);
  color: white;
}

:deep(.pure-sidebar__menu .el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(47, 128, 255, 0.18), rgba(53, 198, 167, 0.16));
  color: white;
}

:deep(.pure-sidebar__menu .el-sub-menu .el-menu) {
  margin: 4px 0 10px;
  padding: 8px 6px;
  border-radius: 18px;
  background: rgba(7, 14, 28, 0.28);
  box-shadow: inset 0 0 0 1px rgba(255, 255, 255, 0.04);
}

:deep(.pure-sidebar__menu .el-sub-menu .el-menu-item) {
  background: rgba(255, 255, 255, 0.04);
}

:deep(.pure-sidebar__menu .el-sub-menu__title span) {
  font-weight: 600;
}
</style>
