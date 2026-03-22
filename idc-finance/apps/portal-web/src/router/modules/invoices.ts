import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/invoices",
  name: "invoices",
  component: () => import("@/views/invoices/index.vue"),
  meta: { title: "账单中心", titleEn: "Invoices", requiresAuth: true }
} satisfies RouteRecordRaw;
