import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/orders",
  name: "orders",
  component: () => import("@/views/orders/index.vue"),
  meta: { title: "订单中心", titleEn: "Orders", requiresAuth: true }
} satisfies RouteRecordRaw;
