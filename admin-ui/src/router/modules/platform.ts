import { Bell, Document } from "@element-plus/icons-vue";
import type { RouteRecordRaw } from "vue-router";

const platformRoutes: RouteRecordRaw[] = [
  {
    path: "notifications",
    name: "notifications",
    component: () => import("@/pages/NotificationsPage.vue"),
    meta: {
      title: "通知中心",
      subtitle: "管理通知模板、消息队列、异步任务和发送结果。",
      icon: Bell,
      permissions: ["notifications.view"],
      navGroup: "platform",
      navOrder: 10,
    },
  },
  {
    path: "audits",
    name: "audits",
    component: () => import("@/pages/AuditsPage.vue"),
    meta: {
      title: "审计日志",
      subtitle: "追踪后台操作、财务动作、实例动作和系统留痕。",
      icon: Document,
      permissions: ["audits.view"],
      navGroup: "platform",
      navOrder: 20,
    },
  },
];

export default platformRoutes;
