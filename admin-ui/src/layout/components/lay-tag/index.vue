<script setup lang="ts">
import { Close } from "@element-plus/icons-vue";
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

import { useMultiTagsStoreHook } from "@/stores/modules/multiTags";

const route = useRoute();
const router = useRouter();
const tagsStore = useMultiTagsStoreHook();

const activePath = computed(() => route.path);

function openTag(path: string) {
  if (path === route.path) {
    return;
  }
  router.push(path);
}

function closeTag(path: string) {
  const closingActive = path === route.path;
  tagsStore.removeTag(path);

  if (!closingActive) {
    return;
  }

  const fallback = tagsStore.tags.at(-1);
  router.push(fallback?.path || "/workbench");
}
</script>

<template>
  <div class="pure-tags">
    <div
      v-for="tag in tagsStore.tags"
      :key="tag.path"
      class="pure-tag"
      :class="{ active: activePath === tag.path }"
      @click="openTag(tag.path)"
    >
      <span>{{ tag.title }}</span>
      <el-icon v-if="tag.closable" class="pure-tag__close" @click.stop="closeTag(tag.path)">
        <Close />
      </el-icon>
    </div>
  </div>
</template>

<style scoped>
.pure-tags {
  display: flex;
  gap: 10px;
  padding: 12px 20px 0;
  overflow: auto hidden;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.75), rgba(255, 255, 255, 0));
}

.pure-tag {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 9px 14px;
  border-radius: 14px 14px 0 0;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid var(--border-color);
  border-bottom: none;
  color: var(--text-secondary);
  cursor: pointer;
  white-space: nowrap;
}

.pure-tag.active {
  color: var(--text-main);
  background: #fff;
  box-shadow: 0 14px 28px rgba(18, 27, 44, 0.08);
}

.pure-tag__close {
  font-size: 12px;
}
</style>
