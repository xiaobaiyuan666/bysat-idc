import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/system",
  redirect: "/system/admins",
  meta: {
    title: "\u7cfb\u7edf",
    titleEn: "Settings",
    requiresAuth: true
  },
  children: [
    {
      path: "admins",
      name: "system-admins",
      component: () => import("@/views/system/admins/index.vue"),
      meta: { title: "\u7ba1\u7406\u5458", titleEn: "Admins", requiresAuth: true }
    },
    {
      path: "roles",
      name: "system-roles",
      component: () => import("@/views/system/roles/index.vue"),
      meta: { title: "\u89d2\u8272\u6743\u9650", titleEn: "Roles & Permissions", requiresAuth: true }
    },
    {
      path: "menus",
      name: "system-menus",
      component: () => import("@/views/system/menus/index.vue"),
      meta: { title: "\u83dc\u5355\u7ba1\u7406", titleEn: "Menus", requiresAuth: true }
    },
    {
      path: "audit",
      name: "system-audit",
      component: () => import("@/views/system/audit/index.vue"),
      meta: { title: "\u5ba1\u8ba1\u65e5\u5fd7", titleEn: "Audit Logs", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
