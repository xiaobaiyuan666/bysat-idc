<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { isLegacyMojibake, localeOptions, pickLabel, resolveMetaTitle } from "@/locales";
import { useLocaleStore, useSettingsStore, useUserStore } from "@/store";

const router = useRouter();
const route = useRoute();
const localeStore = useLocaleStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();

const pageTitle = computed(() =>
  resolveMetaTitle(route.meta, localeStore.locale, settingsStore.title, settingsStore.titleEn)
);

const portalLabel = computed(() => pickLabel(localeStore.locale, "客户控制台", "Client Portal"));
const subtitle = computed(() =>
  pickLabel(
    localeStore.locale,
    "统一管理服务、财务、订单、工单与账户设置",
    "Manage services, finance, orders, tickets, and account settings in one place"
  )
);
const languageLabel = computed(() => pickLabel(localeStore.locale, "语言", "Language"));
const logoutLabel = computed(() => pickLabel(localeStore.locale, "退出登录", "Sign out"));
const displayName = computed(() => {
  if (!userStore.displayName || isLegacyMojibake(userStore.displayName)) {
    return pickLabel(localeStore.locale, "演示客户", "Demo Client");
  }
  return userStore.displayName;
});

function handleLogout() {
  userStore.logout();
  void router.push("/login");
}
</script>

<template>
  <header class="shell-navbar portal-navbar">
    <div class="navbar-title">
      <div class="navbar-title__top">
        <span class="meta-chip meta-chip--success">{{ portalLabel }}</span>
        <span class="navbar-title__heading">{{ pageTitle }}</span>
      </div>
      <div class="navbar-title__meta">
        <span>{{ subtitle }}</span>
        <span>www.bysat.com</span>
      </div>
    </div>

    <div class="portal-navbar__actions">
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
            {{ item.shortLabel }}
          </el-button>
        </el-button-group>
      </div>

      <div class="user-panel">
        <div class="user-panel__badge">BY</div>
        <div class="user-panel__info">
          <strong>{{ displayName }}</strong>
          <span>{{ settingsStore.subtitle }}</span>
        </div>
      </div>

      <el-button type="primary" plain @click="handleLogout">{{ logoutLabel }}</el-button>
    </div>
  </header>
</template>
