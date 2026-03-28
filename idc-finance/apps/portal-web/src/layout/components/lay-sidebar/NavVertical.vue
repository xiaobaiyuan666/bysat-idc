<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { pickLabel } from "@/locales";
import { useLocaleStore, usePermissionStore, useSettingsStore } from "@/store";
import MenuTreeItem from "./MenuTreeItem.vue";

const router = useRouter();
const route = useRoute();
const localeStore = useLocaleStore();
const permissionStore = usePermissionStore();
const settingsStore = useSettingsStore();

const rootMenus = computed(() => permissionStore.menus ?? []);
const brandTitle = computed(() => pickLabel(localeStore.locale, settingsStore.title, settingsStore.titleEn));
const brandSubtitle = computed(() => pickLabel(localeStore.locale, settingsStore.subtitle, settingsStore.subtitleEn));
const footerText = computed(() =>
  pickLabel(
    localeStore.locale,
    "系统由江苏白猿网络科技有限公司 - 猿创软件开发 100% AI 开发，全部著作权归白猿科技所有。",
    "100% AI-developed by Jiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development. All copyrights belong to Baiyuan Technology."
  )
);

function handleSelect(path: string) {
  if (path && path !== route.path) {
    void router.push(path);
  }
}
</script>

<template>
  <aside class="shell-sidebar portal-sidebar">
    <div class="menu-brand portal-brand">
      <div class="menu-brand__logo">BY</div>
      <div class="menu-brand__text">
        <div class="menu-brand__title">{{ brandTitle }}</div>
        <div class="menu-brand__subtitle">{{ brandSubtitle }}</div>
      </div>
    </div>

    <div class="sidebar-scroll portal-sidebar__body">
      <div class="portal-sidebar__section">
        {{ pickLabel(localeStore.locale, "门户导航", "Portal Navigation") }}
      </div>
      <el-menu
        class="sidebar-menu"
        :default-active="route.path"
        background-color="transparent"
        text-color="#dbeafe"
        active-text-color="#ffffff"
        @select="handleSelect"
      >
        <MenuTreeItem v-for="item in rootMenus" :key="item.id" :item="item" @select="handleSelect" />
      </el-menu>
    </div>

    <div class="portal-sidebar__footer">
      <div class="portal-sidebar__website">www.bysat.com</div>
      <div>{{ footerText }}</div>
    </div>
  </aside>
</template>
