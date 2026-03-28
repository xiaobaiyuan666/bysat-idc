import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/orders",
  component: () => import("@/layout/components/route-view/index.vue"),
  meta: { requiresAuth: true },
  children: [
    {
      path: "",
      name: "orders",
      component: () => import("@/views/orders/index.vue"),
      meta: { title: "订单中心", titleEn: "Orders", requiresAuth: true }
    },
    {
      path: ":id",
      name: "order-detail",
      component: () => import("@/views/orders/index.vue"),
      meta: { title: "订单详情", titleEn: "Order Detail", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
