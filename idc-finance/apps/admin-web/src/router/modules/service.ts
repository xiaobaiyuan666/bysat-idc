import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/services",
  redirect: "/services/list",
  meta: {
    title: "\u4e1a\u52a1",
    titleEn: "Business",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "service-list",
      component: () => import("@/views/service/list/index.vue"),
      meta: {
        title: "\u4e1a\u52a1\u5217\u8868",
        titleEn: "Service List",
        requiresAuth: true
      }
    },
    {
      path: "detail/:id",
      name: "service-detail",
      component: () => import("@/views/service/detail/index.vue"),
      meta: {
        title: "\u670d\u52a1\u5de5\u4f5c\u53f0",
        titleEn: "Service Workbench",
        requiresAuth: true
      }
    }
  ]
} satisfies RouteRecordRaw;
