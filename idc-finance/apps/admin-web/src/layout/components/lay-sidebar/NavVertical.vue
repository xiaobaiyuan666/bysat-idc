<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { pickLabel } from "@/locales";
import { useLocaleStore, usePermissionStore, useSettingsStore } from "@/store";
import MenuTreeItem from "./MenuTreeItem.vue";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();
const permissionStore = usePermissionStore();
const settingsStore = useSettingsStore();

const activeMenu = computed(() => route.path);
const brandTitle = computed(() => pickLabel(localeStore.locale, settingsStore.title, settingsStore.titleEn));
const brandSubtitle = computed(() => pickLabel(localeStore.locale, settingsStore.subtitle, settingsStore.subtitleEn));
const footerText = computed(() =>
  pickLabel(
    localeStore.locale,
    "系统 100% AI 开发\n著作权归江苏白猿网络科技有限公司所有\n官网：www.bysat.com",
    "100% AI-developed system\nCopyright owned by Jiangsu Baiyuan Network Technology Co., Ltd.\nWebsite: www.bysat.com"
  )
);

function handleSelect(path: string) {
  router.push(path);
}
</script>

<template>
  <aside class="shell-sidebar">
    <div class="menu-brand">
      <div class="menu-brand__logo">BY</div>
      <div class="menu-brand__text">
        <div class="menu-brand__title">{{ brandTitle }}</div>
        <div class="menu-brand__subtitle">{{ brandSubtitle }}</div>
      </div>
    </div>

    <el-scrollbar class="sidebar-scroll">
      <el-menu
        class="sidebar-menu"
        :default-active="activeMenu"
        background-color="transparent"
        text-color="#dbeafe"
        active-text-color="#ffffff"
        @select="handleSelect"
      >
        <MenuTreeItem v-for="item in permissionStore.menus" :key="item.id" :item="item" />
      </el-menu>
    </el-scrollbar>

    <div class="portal-sidebar__footer" style="white-space: pre-line">
      {{ footerText }}
    </div>
  </aside>
</template>
