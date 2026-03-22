import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/services",
  redirect: "/services/list",
  meta: {
    title: "业务",
    titleEn: "Business",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "service-list",
      component: () => import("@/views/service/list/index.vue"),
      meta: {
        title: "业务列表",
        titleEn: "Service List",
        requiresAuth: true
      }
    },
    {
      path: "detail/:id",
      name: "service-detail",
      component: () => import("@/views/service/detail/index.vue"),
      meta: {
        title: "服务工作台",
        titleEn: "Service Workbench",
        requiresAuth: true
      }
    }
  ]
} satisfies RouteRecordRaw;
