import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/reports",
  redirect: "/reports/overview",
  meta: {
    title: "功能",
    titleEn: "Reports",
    requiresAuth: true
  },
  children: [
    {
      path: "overview",
      name: "report-overview",
      component: () => import("@/views/report/index.vue"),
      meta: { title: "统计报表", titleEn: "Reports", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
