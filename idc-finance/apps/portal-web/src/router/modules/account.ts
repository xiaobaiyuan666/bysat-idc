import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/account",
  alias: ["/account/profile", "/account/contacts", "/account/identity", "/account/security"],
  name: "account",
  component: () => import("@/views/account/index.vue"),
  meta: { title: "账户中心", titleEn: "Account", requiresAuth: true }
} satisfies RouteRecordRaw;
