<script setup lang="ts">
import { ArrowDown, Fold, Menu, SwitchButton } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import { computed, ref } from "vue";
import { useRoute, useRouter } from "vue-router";

import { useAuthStore } from "@/stores/auth";
import { getLabel, roleMap } from "@/utils/maps";

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const collapsed = ref(false);

const groupMap = {
  workbench: {
    title: "工作台",
    description: "总览与待办",
  },
  business: {
    title: "客户与业务",
    description: "客户、订单、服务、工单",
  },
  finance: {
    title: "财务",
    description: "计费、账单、收款、对账",
  },
  cloud: {
    title: "商品与资源",
    description: "产品、云架构、资源中心",
  },
  platform: {
    title: "平台",
    description: "通知与审计",
  },
} as const;

const menuGroups = computed(() => {
  const visibleRoutes = router
    .getRoutes()
    .filter((item) => item.meta?.title && item.meta?.nav !== false)
    .filter((item) => item.path !== "/login")
    .filter((item) => item.path.startsWith("/"))
    .filter((item) => {
      const permissions = item.meta?.permissions as string[] | undefined;

      if (!permissions || permissions.length === 0) {
        return true;
      }

      return permissions.some((permission) =>
        authStore.user?.permissions?.includes(permission),
      );
    })
    .sort((a, b) => Number(a.meta?.navOrder ?? 999) - Number(b.meta?.navOrder ?? 999));

  return Object.entries(groupMap)
    .map(([key, group]) => ({
      key,
      ...group,
      items: visibleRoutes.filter((item) => item.meta?.navGroup === key),
    }))
    .filter((group) => group.items.length > 0);
});

const activeMenu = computed(() => String(route.meta.activeMenu ?? route.path));

const parentRoute = computed(() => {
  if (!route.meta.activeMenu || route.meta.activeMenu === route.path) {
    return null;
  }

  return router.getRoutes().find((item) => item.path === route.meta.activeMenu) ?? null;
});

async function onLogout() {
  await authStore.logout();
  ElMessage.success("已退出登录");
  router.push("/login");
}
</script>

<template>
  <el-container class="page-shell">
    <el-aside :width="collapsed ? '84px' : '286px'" class="layout-aside">
      <div class="brand-box">
        <div class="brand-mark">ID</div>
        <div v-if="!collapsed" class="brand-copy">
          <div class="brand-title">IDC 公有云业务系统</div>
          <div class="brand-subtitle">客户、订单、财务、资源、工单统一运营后台</div>
        </div>
      </div>

      <el-scrollbar class="menu-scroll">
        <el-menu
          :default-active="activeMenu"
          :collapse="collapsed"
          unique-opened
          router
          class="layout-menu"
        >
          <el-menu-item-group v-for="group in menuGroups" :key="group.key" class="menu-group">
            <template #title>
              <div v-if="!collapsed" class="menu-group-title">
                <span>{{ group.title }}</span>
                <small>{{ group.description }}</small>
              </div>
            </template>

            <el-menu-item v-for="item in group.items" :key="item.path" :index="item.path">
              <el-icon>
                <component :is="item.meta.icon" />
              </el-icon>
              <template #title>{{ item.meta.title }}</template>
            </el-menu-item>
          </el-menu-item-group>
        </el-menu>
      </el-scrollbar>
    </el-aside>

    <el-container>
      <el-header class="layout-header">
        <div class="header-left">
          <el-button circle text @click="collapsed = !collapsed">
            <el-icon><component :is="collapsed ? Menu : Fold" /></el-icon>
          </el-button>
          <div>
            <div class="header-breadcrumb">
              <span v-if="parentRoute">{{ parentRoute.meta.title }}</span>
              <span v-if="parentRoute" class="divider">/</span>
              <span>{{ route.meta.title ?? "管理后台" }}</span>
            </div>
            <div class="header-title">{{ route.meta.title ?? "管理后台" }}</div>
            <div class="header-subtitle">
              {{
                route.meta.subtitle ??
                "面向 IDC 公有云销售、财务结算、服务生命周期和资源交付的一体化后台。"
              }}
            </div>
          </div>
        </div>

        <div class="header-right">
          <el-tag type="success" effect="light">正式魔方云已接入</el-tag>
          <el-dropdown trigger="click">
            <div class="user-pill">
              <div class="user-avatar">{{ authStore.user?.name?.slice(0, 1) }}</div>
              <div class="user-copy">
                <div class="user-name">{{ authStore.user?.name }}</div>
                <div class="user-role">{{ getLabel(roleMap, authStore.user?.role) }}</div>
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
      </el-header>

      <el-main class="layout-main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<style scoped>
.layout-aside {
  background: linear-gradient(180deg, var(--menu-bg) 0%, var(--menu-bg-2) 100%);
  color: white;
  transition: width 0.2s ease;
}

.brand-box {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 22px 18px 18px;
}

.brand-mark {
  width: 44px;
  height: 44px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  color: white;
  font-weight: 800;
  background: linear-gradient(135deg, #2f80ff, #35c6a7);
}

.brand-title {
  font-size: 16px;
  font-weight: 700;
}

.brand-subtitle {
  margin-top: 4px;
  color: rgba(255, 255, 255, 0.58);
  font-size: 12px;
  line-height: 1.6;
}

.menu-scroll {
  height: calc(100vh - 84px);
}

:deep(.layout-menu) {
  background: transparent;
  border-right: none;
}

:deep(.layout-menu .el-sub-menu__title),
:deep(.layout-menu .el-menu-item),
:deep(.layout-menu .el-menu-item-group__title) {
  color: rgba(255, 255, 255, 0.74);
}

:deep(.layout-menu .el-menu-item:hover),
:deep(.layout-menu .el-sub-menu__title:hover) {
  background: rgba(255, 255, 255, 0.06);
  color: white;
}

:deep(.layout-menu .el-menu-item.is-active) {
  background: linear-gradient(90deg, rgba(47, 128, 255, 0.18), rgba(53, 198, 167, 0.16));
  color: white;
}

.menu-group-title {
  display: flex;
  flex-direction: column;
  gap: 2px;
  line-height: 1.4;
}

.menu-group-title small {
  color: rgba(255, 255, 255, 0.42);
}

.layout-header {
  height: 82px;
  border-bottom: 1px solid var(--line);
  background: rgba(255, 255, 255, 0.86);
  backdrop-filter: blur(12px);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 14px;
}

.header-breadcrumb {
  color: var(--muted);
  font-size: 12px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.header-breadcrumb .divider {
  opacity: 0.6;
}

.header-title {
  margin-top: 6px;
  font-size: 22px;
  font-weight: 700;
  color: var(--text);
}

.header-subtitle {
  margin-top: 4px;
  font-size: 13px;
  color: var(--muted);
}

.header-right {
  display: flex;
  align-items: center;
  gap: 14px;
}

.user-pill {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 14px;
  border-radius: 18px;
  background: white;
  border: 1px solid rgba(10, 35, 66, 0.08);
  cursor: pointer;
}

.user-avatar {
  width: 38px;
  height: 38px;
  border-radius: 14px;
  display: grid;
  place-items: center;
  background: linear-gradient(135deg, rgba(47, 128, 255, 0.16), rgba(53, 198, 167, 0.12));
  color: var(--text);
  font-weight: 700;
}

.user-copy {
  line-height: 1.4;
}

.user-name {
  color: var(--text);
  font-weight: 600;
}

.user-role {
  color: var(--muted);
  font-size: 12px;
}

.layout-main {
  padding: 24px;
  background:
    radial-gradient(circle at top right, rgba(47, 128, 255, 0.09), transparent 30%),
    radial-gradient(circle at top left, rgba(53, 198, 167, 0.08), transparent 24%),
    var(--bg);
}
</style>
