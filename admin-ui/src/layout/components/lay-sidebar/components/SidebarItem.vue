<script setup lang="ts">
import type { MenuNode } from "@/layout/types";

defineProps<{
  item: MenuNode;
}>();
</script>

<template>
  <el-sub-menu
    v-if="item.children?.length"
    :index="item.index"
    popper-class="pure-sidebar-popper"
  >
    <template #title>
      <span>{{ item.title }}</span>
    </template>

    <SidebarItem v-for="child in item.children" :key="child.index" :item="child" />
  </el-sub-menu>

  <el-menu-item v-else :index="item.path">
    <el-icon v-if="item.icon">
      <component :is="item.icon" />
    </el-icon>
    <template #title>{{ item.title }}</template>
  </el-menu-item>
</template>
