import { defineStore } from "pinia";
import { markRaw } from "vue";
import type { RouteRecordNormalized } from "vue-router";

import type { MenuNode } from "@/layout/types";

const groupMap = {
  workbench: {
    title: "工作台",
    subtitle: "总览、待办与运营视图",
  },
  business: {
    title: "客户与业务",
    subtitle: "客户、订单、服务与工单",
  },
  finance: {
    title: "财务",
    subtitle: "计费、账单、收款与对账",
  },
  cloud: {
    title: "商品与资源",
    subtitle: "产品、云架构与资源中心",
  },
  platform: {
    title: "平台",
    subtitle: "通知、审计与平台配置",
  },
} as const;

function hasRouteAccess(route: RouteRecordNormalized, permissions: string[]) {
  const routePermissions = route.meta?.permissions;
  if (!Array.isArray(routePermissions) || routePermissions.length === 0) {
    return true;
  }

  return routePermissions.some((permission) => permissions.includes(String(permission)));
}

export const usePermissionStore = defineStore("pure-permission", {
  state: () => ({
    wholeMenus: [] as MenuNode[],
  }),
  actions: {
    syncMenus(routes: RouteRecordNormalized[], permissions: string[]) {
      const visibleRoutes = routes
        .filter((route) => route.meta?.title && route.meta?.nav !== false)
        .filter((route) => route.path !== "/login")
        .filter((route) => route.path.startsWith("/"))
        .filter((route) => hasRouteAccess(route, permissions))
        .sort((a, b) => Number(a.meta?.navOrder ?? 999) - Number(b.meta?.navOrder ?? 999));

      this.wholeMenus = Object.entries(groupMap)
        .map(([groupKey, group]) => ({
          index: `/__group/${groupKey}`,
          path: `/__group/${groupKey}`,
          title: group.title,
          subtitle: group.subtitle,
          children: visibleRoutes
            .filter((route) => route.meta?.navGroup === groupKey)
            .map((route) => ({
              index: route.path,
              path: route.path,
              title: String(route.meta?.title ?? ""),
              subtitle: String(route.meta?.subtitle ?? ""),
              icon: route.meta?.icon
                ? (markRaw(route.meta.icon as object) as MenuNode["icon"])
                : undefined,
              activePath: String(route.meta?.activeMenu ?? route.path),
            })),
        }))
        .filter((group) => group.children && group.children.length > 0);
    },
  },
});

export function usePermissionStoreHook() {
  return usePermissionStore();
}
