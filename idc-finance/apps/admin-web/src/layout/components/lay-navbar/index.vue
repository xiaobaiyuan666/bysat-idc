<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { isLegacyMojibake, localeOptions, pickLabel, resolveMetaTitle } from "@/locales";
import { useLocaleStore, useUserStore } from "@/store";

const route = useRoute();
const router = useRouter();
const localeStore = useLocaleStore();
const userStore = useUserStore();

const breadcrumbTitles = computed(() =>
  route.matched.map(item => resolveMetaTitle(item.meta, localeStore.locale)).filter(Boolean)
);

const currentTitle = computed(() => {
  const titles = breadcrumbTitles.value;
  return titles[titles.length - 1] || pickLabel(localeStore.locale, "工作台", "Workbench");
});
const currentSection = computed(() => breadcrumbTitles.value[0] || pickLabel(localeStore.locale, "后台", "Admin"));
const displayName = computed(() => {
  if (!userStore.displayName || isLegacyMojibake(userStore.displayName)) {
    return pickLabel(localeStore.locale, "系统管理员", "System Admin");
  }
  return userStore.displayName;
});
const userInitial = computed(() => displayName.value.slice(0, 1).toUpperCase());
const mysqlStatus = computed(() => pickLabel(localeStore.locale, "MySQL 在线", "MySQL Online"));
const productSubtitle = computed(() =>
  pickLabel(
    localeStore.locale,
    "江苏白猿网络科技有限公司 · 猿创软件业务组",
    "Jiangsu Baiyuan Network Technology Co., Ltd. · Yuanchuang Software Business Group"
  )
);
const roleLabel = computed(() => pickLabel(localeStore.locale, "系统管理员", "System Admin"));
const logoutLabel = computed(() => pickLabel(localeStore.locale, "退出登录", "Sign out"));
const languageLabel = computed(() => pickLabel(localeStore.locale, "语言", "Language"));

function handleLogout() {
  userStore.logout();
  router.push("/login");
}
</script>

<template>
  <header class="shell-navbar">
    <div class="navbar-title">
      <div class="navbar-title__top">
        <span class="navbar-title__heading">{{ currentTitle }}</span>
        <span class="meta-chip">{{ currentSection }}</span>
        <span class="meta-chip meta-chip--success">{{ mysqlStatus }}</span>
      </div>
      <div class="navbar-title__meta">
        <span>{{ productSubtitle }}</span>
        <span v-if="breadcrumbTitles.length > 1">/ {{ breadcrumbTitles.join(" / ") }}</span>
      </div>
    </div>

    <div class="user-panel">
      <div class="locale-switch">
        <span class="locale-switch__label">{{ languageLabel }}</span>
        <el-button-group>
          <el-button
            v-for="item in localeOptions"
            :key="item.value"
            size="small"
            :type="localeStore.locale === item.value ? 'primary' : 'default'"
            @click="localeStore.setLocale(item.value)"
          >
            {{ item.label }}
          </el-button>
        </el-button-group>
      </div>
      <div class="user-panel__badge">{{ userInitial }}</div>
      <div>
        <div class="user-panel__name">{{ displayName }}</div>
        <div class="user-panel__role">{{ roleLabel }}</div>
      </div>
      <el-button type="primary" plain @click="handleLogout">{{ logoutLabel }}</el-button>
    </div>
  </header>
</template>
