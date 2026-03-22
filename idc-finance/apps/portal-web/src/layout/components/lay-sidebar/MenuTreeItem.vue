<script setup lang="ts">
import { computed } from "vue";
import { pickLabel } from "@/locales";
import { useLocaleStore } from "@/store";

type MenuNode = {
  id?: number;
  title: string;
  titleEn?: string;
  path: string;
  children?: MenuNode[];
};

const props = defineProps<{
  item: MenuNode;
}>();

const emit = defineEmits<{
  select: [path: string];
}>();

const localeStore = useLocaleStore();
const hasChildren = computed(() => Boolean(props.item.children?.length));
const displayTitle = computed(() =>
  pickLabel(localeStore.locale, props.item.title, props.item.titleEn || props.item.title)
);

function handleSelect(path: string) {
  emit("select", path);
}
</script>

<template>
  <el-sub-menu v-if="hasChildren" :index="item.path">
    <template #title>{{ displayTitle }}</template>
    <MenuTreeItem
      v-for="child in item.children ?? []"
      :key="child.id"
      :item="child"
      @select="handleSelect"
    />
  </el-sub-menu>
  <el-menu-item v-else :index="item.path">
    {{ displayTitle }}
  </el-menu-item>
</template>
