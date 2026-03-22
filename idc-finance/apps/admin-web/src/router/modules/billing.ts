import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/billing",
  redirect: "/billing/invoices",
  meta: {
    title: "财务",
    titleEn: "Finance",
    requiresAuth: true
  },
  children: [
    {
      path: "invoices",
      name: "billing-invoices",
      component: () => import("@/views/billing/invoices/index.vue"),
      meta: { title: "账单管理", titleEn: "Invoices", requiresAuth: true }
    },
    {
      path: "invoices/:id",
      name: "billing-invoice-detail",
      component: () => import("@/views/billing/invoice-detail/index.vue"),
      meta: { title: "账单工作台", titleEn: "Invoice Workbench", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
