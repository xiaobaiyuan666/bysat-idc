import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/reports",
  redirect: "/reports/overview",
  meta: {
    title: "\u62a5\u8868\u4e2d\u5fc3",
    titleEn: "Reports",
    requiresAuth: true
  },
  children: [
    {
      path: "overview",
      name: "report-overview",
      component: () => import("@/views/report/index.vue"),
      meta: { title: "\u8fd0\u8425\u62a5\u8868", titleEn: "Operations Reports", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
