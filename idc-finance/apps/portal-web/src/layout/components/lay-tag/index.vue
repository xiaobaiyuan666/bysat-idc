<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { pickLabel } from "@/locales";
import { useLocaleStore, useMultiTagsStore } from "@/store";

const router = useRouter();
const route = useRoute();
const localeStore = useLocaleStore();
const multiTagsStore = useMultiTagsStore();

function isActive(path: string) {
  return route.path === path || route.path.startsWith(`${path}/`);
}

function handleTagClick(path: string) {
  if (path !== route.path) {
    router.push(path);
  }
}
</script>

<template>
  <div class="portal-tags">
    <el-tag
      v-for="item in multiTagsStore.items"
      :key="item.path"
      :type="isActive(item.path) ? 'primary' : 'info'"
      :effect="isActive(item.path) ? 'dark' : 'plain'"
      @click="handleTagClick(item.path)"
    >
      {{ pickLabel(localeStore.locale, item.title, item.titleEn || item.title) }}
    </el-tag>
  </div>
</template>
