import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/services",
  component: () => import("@/layout/components/lay-content/index.vue"),
  meta: { requiresAuth: true },
  children: [
    {
      path: "",
      name: "services",
      component: () => import("@/views/services/index.vue"),
      meta: { title: "服务中心", titleEn: "Services", requiresAuth: true }
    },
    {
      path: ":id",
      name: "service-detail",
      component: () => import("@/views/services/detail.vue"),
      meta: { title: "服务详情", titleEn: "Service Detail", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
