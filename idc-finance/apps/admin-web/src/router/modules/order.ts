import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/orders",
  redirect: "/orders/list",
  meta: {
    title: "\u4e1a\u52a1",
    titleEn: "Business",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "order-list",
      component: () => import("@/views/order/list/index.vue"),
      meta: { title: "\u4ea7\u54c1\u8ba2\u5355", titleEn: "Orders", requiresAuth: true }
    },
    {
      path: "change-orders",
      name: "service-change-orders",
      component: () => import("@/views/order/change-orders/index.vue"),
      meta: { title: "\u6539\u914d\u5355", titleEn: "Change Orders", requiresAuth: true }
    },
    {
      path: "detail/:id",
      name: "order-detail",
      component: () => import("@/views/order/detail/index.vue"),
      meta: { title: "\u8ba2\u5355\u5de5\u4f5c\u53f0", titleEn: "Order Workbench", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
