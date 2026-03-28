import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/store",
  name: "store",
  component: () => import("@/views/store/index.vue"),
  meta: { title: "云产品商城", titleEn: "Marketplace", requiresAuth: true }
} satisfies RouteRecordRaw;
