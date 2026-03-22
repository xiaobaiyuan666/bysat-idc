<script setup lang="ts">
import { ArrowDown, Fold, Menu, Refresh, SwitchButton } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";

import { useAppStoreHook } from "@/stores/modules/app";
import { useAuthStore } from "@/stores/auth";
import { getLabel, roleMap } from "@/utils/maps";

const route = useRoute();
const router = useRouter();
const appStore = useAppStoreHook();
const authStore = useAuthStore();

const parentTitle = computed(() => {
  const activeMenu = String(route.meta.activeMenu ?? "");
  if (!activeMenu || activeMenu === route.path) {
    return "";
  }

  const parent = router.getRoutes().find((item) => item.path === activeMenu);
  return String(parent?.meta?.title ?? "");
});

async function onLogout() {
  await authStore.logout();
  ElMessage.success("已退出登录");
  router.push("/login");
}

function refreshCurrent() {
  router.replace({
    path: "/dashboard",
  }).then(() => {
    router.replace(route.fullPath);
  });
}
</script>

<template>
  <header class="pure-navbar">
    <div class="pure-navbar__left">
      <el-button circle text @click="appStore.toggleSidebar()">
        <el-icon><component :is="appStore.sidebar.opened ? Fold : Menu" /></el-icon>
      </el-button>

      <div class="pure-navbar__copy">
        <div class="pure-navbar__breadcrumb">
          <span v-if="parentTitle">{{ parentTitle }}</span>
          <span v-if="parentTitle" class="divider">/</span>
          <span>{{ route.meta.title ?? "管理后台" }}</span>
        </div>
        <div class="pure-navbar__title">{{ route.meta.title ?? "管理后台" }}</div>
        <div class="pure-navbar__subtitle">
          {{
            route.meta.subtitle ??
            "面向 IDC 公有云售卖、财务结算、服务生命周期和资源交付的一体化后台。"
          }}
        </div>
      </div>
    </div>

    <div class="pure-navbar__right">
      <el-button text :icon="Refresh" @click="refreshCurrent">刷新页面</el-button>
      <el-tag type="success" effect="light">正式魔方云已接入</el-tag>
      <el-dropdown trigger="click">
        <div class="pure-navbar__user">
          <div class="pure-navbar__avatar">{{ authStore.user?.name?.slice(0, 1) }}</div>
          <div class="pure-navbar__user-copy">
            <div class="pure-navbar__user-name">{{ authStore.user?.name }}</div>
            <div class="pure-navbar__user-role">{{ getLabel(roleMap, authStore.user?.role) }}</div>
          </div>
          <el-icon><ArrowDown /></el-icon>
        </div>

        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="onLogout">
              <el-icon><SwitchButton /></el-icon>
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>
  </header>
</template>

<style scoped>
.pure-navbar {
  height: 82px;
  border-bottom: 1px solid var(--border-color);
  background: rgba(255, 255, 255, 0.94);
  backdrop-filter: blur(18px);
  box-shadow: 0 14px 28px rgba(18, 27, 44, 0.04);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.pure-navbar__left,
.pure-navbar__right {
  display: flex;
  align-items: center;
  gap: 14px;
}

.pure-navbar__breadcrumb {
  color: var(--text-secondary);
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.pure-navbar__title {
  margin-top: 6px;
  font-size: 22px;
  font-weight: 700;
  color: var(--text-main);
}

.pure-navbar__subtitle {
  margin-top: 4px;
  font-size: 13px;
  color: var(--text-secondary);
}

.pure-navbar__user {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 18px;
  background: rgba(255, 255, 255, 0.96);
  border: 1px solid rgba(10, 35, 66, 0.08);
  cursor: pointer;
}

.pure-navbar__avatar {
  width: 38px;
  height: 38px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, rgba(47, 128, 255, 0.16), rgba(53, 198, 167, 0.12));
  color: var(--text-main);
  font-weight: 700;
}

.pure-navbar__user-name {
  color: var(--text-main);
  font-weight: 600;
}

.pure-navbar__user-role {
  color: var(--text-secondary);
  font-size: 12px;
}
</style>
