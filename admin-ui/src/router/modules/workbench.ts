import { DataAnalysis, TrendCharts } from "@element-plus/icons-vue";
import type { RouteRecordRaw } from "vue-router";

const workbenchRoutes: RouteRecordRaw[] = [
  {
    path: "workbench",
    name: "workbench",
    component: () => import("@/pages/DashboardPage.vue"),
    meta: {
      title: "运营工作台",
      subtitle: "聚合收入、续费、待办、资源库存和魔方云同步结果。",
      icon: DataAnalysis,
      permissions: ["dashboard.view"],
      navGroup: "workbench",
      navOrder: 1,
    },
  },
  {
    path: "reports",
    name: "reports",
    component: () => import("@/pages/ReportsPage.vue"),
    meta: {
      title: "统计报表",
      subtitle: "查看收入趋势、账单结构、渠道表现、区域分布与欠费风险。",
      icon: TrendCharts,
      permissions: ["dashboard.view"],
      navGroup: "workbench",
      navOrder: 2,
    },
  },
];

export default workbenchRoutes;
