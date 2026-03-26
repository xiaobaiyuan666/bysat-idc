import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/workbench",
  name: "workbench",
  component: () => import("@/views/workbench/index.vue"),
  meta: {
    title: "\u5de5\u4f5c\u53f0",
    titleEn: "Workbench",
    requiresAuth: true
  }
} satisfies RouteRecordRaw;
