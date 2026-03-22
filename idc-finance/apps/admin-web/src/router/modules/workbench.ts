import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/workbench",
  name: "workbench",
  component: () => import("@/views/workbench/index.vue"),
  meta: {
    title: "工作台",
    titleEn: "Workbench",
    requiresAuth: true
  }
} satisfies RouteRecordRaw;
