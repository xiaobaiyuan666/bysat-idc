import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/console",
  name: "console",
  component: () => import("@/views/console/index.vue"),
  meta: { title: "控制台", titleEn: "Console", requiresAuth: true }
} satisfies RouteRecordRaw;
