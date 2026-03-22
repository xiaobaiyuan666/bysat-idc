<script setup lang="ts">
import { computed } from "vue";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";
import type { MenuItem } from "@/store/modules/permission";

const props = defineProps<{
  item: MenuItem;
}>();

const localeStore = useLocaleStore();
const displayTitle = computed(() =>
  pickLabel(localeStore.locale, props.item.title, props.item.titleEn || props.item.title)
);
</script>

<template>
  <el-sub-menu v-if="item.children && item.children.length" :index="item.path">
    <template #title>
      <span class="menu-node">
        <span class="menu-node__dot"></span>
        <span>{{ displayTitle }}</span>
      </span>
    </template>
    <MenuTreeItem v-for="child in item.children" :key="child.id" :item="child" />
  </el-sub-menu>
  <el-menu-item v-else :index="item.path">
    <span class="menu-node">
      <span class="menu-node__dot"></span>
      <span>{{ displayTitle }}</span>
    </span>
  </el-menu-item>
</template>
