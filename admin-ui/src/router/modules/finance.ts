import { Coin, CreditCard, Setting, TrendCharts } from "@element-plus/icons-vue";
import type { RouteRecordRaw } from "vue-router";

const financeRoutes: RouteRecordRaw[] = [
  {
    path: "billing",
    name: "billing",
    component: () => import("@/pages/BillingPage.vue"),
    meta: {
      title: "计费中心",
      subtitle: "维护续费、宽限期、停机、税率和支付渠道配置。",
      icon: Setting,
      permissions: ["billing.view"],
      navGroup: "finance",
      navOrder: 10,
    },
  },
  {
    path: "invoices",
    name: "invoices",
    component: () => import("@/pages/InvoicesPage.vue"),
    meta: {
      title: "账单管理",
      subtitle: "管理草稿、签发、逾期、作废和收款联动的账单生命周期。",
      icon: CreditCard,
      permissions: ["invoices.view"],
      navGroup: "finance",
      navOrder: 20,
    },
  },
  {
    path: "invoices/:id",
    name: "invoice-detail",
    component: () => import("@/pages/InvoiceDetailPage.vue"),
    meta: {
      title: "账单详情",
      subtitle: "查看账单摘要、支付记录、回调日志、退款和财务任务。",
      permissions: ["invoices.view"],
      nav: false,
      activeMenu: "/invoices",
    },
  },
  {
    path: "payments",
    name: "payments",
    component: () => import("@/pages/PaymentsPage.vue"),
    meta: {
      title: "收款管理",
      subtitle: "登记收款、支付回调、退款记录和渠道处理结果。",
      icon: Coin,
      permissions: ["payments.manage"],
      navGroup: "finance",
      navOrder: 30,
    },
  },
  {
    path: "reconciliation",
    name: "reconciliation",
    component: () => import("@/pages/ReconciliationPage.vue"),
    meta: {
      title: "财务对账",
      subtitle: "按渠道汇总收款、退款、异常回调和待处理差异。",
      icon: TrendCharts,
      permissions: ["payments.manage"],
      navGroup: "finance",
      navOrder: 40,
    },
  },
];

export default financeRoutes;
