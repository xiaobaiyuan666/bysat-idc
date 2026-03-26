import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/tickets",
  redirect: "/tickets/list",
  meta: {
    title: "\u5de5\u5355",
    titleEn: "Tickets",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "ticket-list",
      component: () => import("@/views/ticket/list/index.vue"),
      meta: {
        title: "\u5de5\u5355\u4e2d\u5fc3",
        titleEn: "Ticket Center",
        requiresAuth: true
      }
    },
    {
      path: "settings",
      name: "ticket-settings",
      component: () => import("@/views/ticket/settings/index.vue"),
      meta: {
        title: "\u5de5\u5355\u914d\u7f6e",
        titleEn: "Ticket Settings",
        requiresAuth: true
      }
    },
    {
      path: "create",
      name: "ticket-create",
      component: () => import("@/views/ticket/create/index.vue"),
      meta: {
        title: "\u65b0\u5efa\u5de5\u5355",
        titleEn: "Create Ticket",
        requiresAuth: true
      }
    },
    {
      path: ":id",
      name: "ticket-detail",
      component: () => import("@/views/ticket/detail/index.vue"),
      meta: {
        title: "\u5de5\u5355\u5de5\u4f5c\u53f0",
        titleEn: "Ticket Workbench",
        requiresAuth: true
      }
    }
  ]
} satisfies RouteRecordRaw;
