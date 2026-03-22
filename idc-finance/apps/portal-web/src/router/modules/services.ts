import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/services",
  name: "services",
  component: () => import("@/views/services/index.vue"),
  meta: { title: "服务中心", titleEn: "Services", requiresAuth: true }
} satisfies RouteRecordRaw;
