<script setup lang="ts">
import { useRouter, useRoute } from "vue-router";
import { pickLabel } from "@/locales";
import { useLocaleStore, useMultiTagsStore } from "@/store";

const router = useRouter();
const route = useRoute();
const localeStore = useLocaleStore();
const multiTagsStore = useMultiTagsStore();

function navigate(path: string) {
  router.push(path);
}

function isActive(path: string) {
  return path === route.path;
}
</script>

<template>
  <div class="shell-tags">
    <el-tag
      v-for="item in multiTagsStore.items"
      :key="item.path"
      :type="isActive(item.path) ? 'primary' : 'info'"
      effect="plain"
      round
      style="cursor: pointer"
      @click="navigate(item.path)"
    >
      {{ pickLabel(localeStore.locale, item.title, item.titleEn || item.title) }}
    </el-tag>
  </div>
</template>
