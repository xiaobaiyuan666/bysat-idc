import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/billing",
  redirect: "/billing/invoices",
  meta: {
    title: "\u8d22\u52a1",
    titleEn: "Finance",
    requiresAuth: true
  },
  children: [
    {
      path: "invoices",
      name: "billing-invoices",
      component: () => import("@/views/billing/invoices/index.vue"),
      meta: { title: "\u8d26\u5355\u7ba1\u7406", titleEn: "Invoices", requiresAuth: true }
    },
    {
      path: "payments",
      name: "billing-payments",
      component: () => import("@/views/billing/payments/index.vue"),
      meta: { title: "\u652f\u4ed8\u8bb0\u5f55", titleEn: "Payments", requiresAuth: true }
    },
    {
      path: "refunds",
      name: "billing-refunds",
      component: () => import("@/views/billing/refunds/index.vue"),
      meta: { title: "\u9000\u6b3e\u8bb0\u5f55", titleEn: "Refunds", requiresAuth: true }
    },
    {
      path: "accounts",
      name: "billing-accounts",
      component: () => import("@/views/billing/accounts/index.vue"),
      meta: { title: "\u8d44\u91d1\u53f0\u8d26", titleEn: "Account Ledger", requiresAuth: true }
    },
    {
      path: "recharges",
      name: "billing-recharges",
      component: () => import("@/views/billing/recharges/index.vue"),
      meta: { title: "\u5145\u503c\u8bb0\u5f55", titleEn: "Recharges", requiresAuth: true }
    },
    {
      path: "adjustments",
      name: "billing-adjustments",
      component: () => import("@/views/billing/adjustments/index.vue"),
      meta: { title: "\u8d44\u91d1\u8c03\u6574", titleEn: "Adjustments", requiresAuth: true }
    },
    {
      path: "invoices/:id",
      name: "billing-invoice-detail",
      component: () => import("@/views/billing/invoice-detail/index.vue"),
      meta: { title: "\u8d26\u5355\u5de5\u4f5c\u53f0", titleEn: "Invoice Workbench", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
