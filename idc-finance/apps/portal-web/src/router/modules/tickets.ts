import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/tickets",
  name: "tickets",
  component: () => import("@/views/tickets/index.vue"),
  meta: { title: "工单中心", titleEn: "Tickets", requiresAuth: true }
} satisfies RouteRecordRaw;
