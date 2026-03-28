import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/invoices",
  component: () => import("@/layout/components/route-view/index.vue"),
  meta: { requiresAuth: true },
  children: [
    {
      path: "",
      name: "invoices",
      component: () => import("@/views/invoices/index.vue"),
      meta: { title: "财务中心", titleEn: "Finance", requiresAuth: true }
    },
    {
      path: ":id",
      name: "invoice-detail",
      component: () => import("@/views/invoices/index.vue"),
      meta: { title: "账单详情", titleEn: "Invoice Detail", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
