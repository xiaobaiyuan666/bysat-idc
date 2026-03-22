import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/orders",
  redirect: "/orders/list",
  meta: {
    title: "业务",
    titleEn: "Business",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "order-list",
      component: () => import("@/views/order/list/index.vue"),
      meta: { title: "产品订单", titleEn: "Orders", requiresAuth: true }
    },
    {
      path: "detail/:id",
      name: "order-detail",
      component: () => import("@/views/order/detail/index.vue"),
      meta: { title: "订单工作台", titleEn: "Order Workbench", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
