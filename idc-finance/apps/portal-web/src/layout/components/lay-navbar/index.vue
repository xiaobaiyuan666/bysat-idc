<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { isLegacyMojibake, localeOptions, pickLabel, resolveMetaTitle } from "@/locales";
import { useLocaleStore, useUserStore } from "@/store";

const router = useRouter();
const route = useRoute();
const localeStore = useLocaleStore();
const userStore = useUserStore();

const title = computed(() => resolveMetaTitle(route.meta, localeStore.locale, "客户中心", "Client Area"));
const subtitle = computed(() =>
  pickLabel(localeStore.locale, "客户自助中心 · 订单 / 服务 / 账单 / 工单 / 钱包", "Self-service center · Orders / Services / Invoices / Tickets / Wallet")
);
const logoutLabel = computed(() => pickLabel(localeStore.locale, "退出登录", "Sign out"));
const languageLabel = computed(() => pickLabel(localeStore.locale, "语言", "Language"));
const displayName = computed(() => {
  if (!userStore.displayName || isLegacyMojibake(userStore.displayName)) {
    return pickLabel(localeStore.locale, "演示客户", "Demo Client");
  }
  return userStore.displayName;
});

function handleLogout() {
  userStore.logout();
  router.push("/login");
}
</script>

<template>
  <header class="portal-navbar">
    <div>
      <div class="portal-navbar__title">{{ title }}</div>
      <div class="portal-navbar__subtitle">{{ subtitle }}</div>
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
            {{ item.label }}
          </el-button>
        </el-button-group>
      </div>
      <el-tag type="info" effect="plain">{{ displayName }}</el-tag>
      <el-button type="primary" plain @click="handleLogout">{{ logoutLabel }}</el-button>
    </div>
  </header>
</template>
