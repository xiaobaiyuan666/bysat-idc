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
    "无穷云IDC业务管理系统\n江苏白猿网络科技有限公司-猿创软件开发\n100% AI 开发，著作权归江苏白猿网络科技有限公司所有",
    "Infinity Cloud IDC Management System\nJiangsu Baiyuan Network Technology Co., Ltd. - Yuanchuang Software Development\n100% AI-developed, copyright owned by Jiangsu Baiyuan Network Technology Co., Ltd."
  )
);

function handleSelect(path: string) {
  if (path && path !== route.path) {
    router.push(path);
  }
}
</script>

<template>
  <aside class="portal-sidebar">
    <div class="portal-brand">
      <div class="portal-brand__mark">∞</div>
      <div>
        <div class="portal-brand__title">{{ brandTitle }}</div>
        <div class="portal-brand__subtitle">{{ brandSubtitle }}</div>
      </div>
    </div>

    <div class="portal-sidebar__body">
      <el-menu
        :default-active="route.path"
        background-color="transparent"
        text-color="#dbeafe"
        active-text-color="#ffffff"
        @select="handleSelect"
      >
        <MenuTreeItem v-for="item in rootMenus" :key="item.id" :item="item" @select="handleSelect" />
      </el-menu>
    </div>

    <div class="portal-sidebar__footer" style="white-space: pre-line">
      {{ footerText }}
    </div>
  </aside>
</template>
