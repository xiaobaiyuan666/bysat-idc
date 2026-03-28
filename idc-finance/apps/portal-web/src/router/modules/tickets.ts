import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/tickets",
  component: () => import("@/layout/components/route-view/index.vue"),
  meta: { requiresAuth: true },
  children: [
    {
      path: "",
      name: "tickets",
      component: () => import("@/views/tickets/index.vue"),
      meta: { title: "工单中心", titleEn: "Tickets", requiresAuth: true }
    },
    {
      path: ":id",
      name: "ticket-detail",
      component: () => import("@/views/tickets/index.vue"),
      meta: { title: "工单详情", titleEn: "Ticket Detail", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
