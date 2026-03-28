import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/wallet",
  alias: ["/wallet/transactions", "/wallet/recharges"],
  name: "wallet",
  component: () => import("@/views/wallet/index.vue"),
  meta: { title: "钱包中心", titleEn: "Wallet", requiresAuth: true }
} satisfies RouteRecordRaw;
