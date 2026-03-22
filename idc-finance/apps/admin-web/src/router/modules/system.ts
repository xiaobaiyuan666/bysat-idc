import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/system",
  redirect: "/system/admins",
  meta: {
    title: "设置",
    titleEn: "Settings",
    requiresAuth: true
  },
  children: [
    {
      path: "admins",
      name: "system-admins",
      component: () => import("@/views/system/admins/index.vue"),
      meta: { title: "管理员", titleEn: "Admins", requiresAuth: true }
    },
    {
      path: "roles",
      name: "system-roles",
      component: () => import("@/views/system/roles/index.vue"),
      meta: { title: "角色权限", titleEn: "Roles & Permissions", requiresAuth: true }
    },
    {
      path: "menus",
      name: "system-menus",
      component: () => import("@/views/system/menus/index.vue"),
      meta: { title: "菜单管理", titleEn: "Menus", requiresAuth: true }
    },
    {
      path: "audit",
      name: "system-audit",
      component: () => import("@/views/system/audit/index.vue"),
      meta: { title: "审计日志", titleEn: "Audit Logs", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
