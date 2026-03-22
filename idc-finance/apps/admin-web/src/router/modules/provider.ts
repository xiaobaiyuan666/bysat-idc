import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/providers",
  redirect: "/providers/mofang",
  meta: {
    title: "资源与商店",
    titleEn: "Resources",
    requiresAuth: true
  },
  children: [
    {
      path: "accounts",
      name: "provider-accounts",
      component: () => import("@/views/provider/accounts/index.vue"),
      meta: {
        title: "接口账户",
        titleEn: "Provider Accounts",
        requiresAuth: true
      }
    },
    {
      path: "mofang",
      name: "provider-mofang",
      component: () => import("@/views/provider/mofang/index.vue"),
      meta: {
        title: "魔方云",
        titleEn: "Mofang Cloud",
        requiresAuth: true
      }
    },
    {
      path: "automation",
      name: "provider-automation",
      component: () => import("@/views/provider/automation/index.vue"),
      meta: {
        title: "自动化任务中心",
        titleEn: "Automation Tasks",
        requiresAuth: true
      }
    },
    {
      path: "automation-settings",
      name: "provider-automation-settings",
      component: () => import("@/views/provider/automation-settings/index.vue"),
      meta: {
        title: "自动化策略",
        titleEn: "Automation Policies",
        requiresAuth: true
      }
    }
  ]
} satisfies RouteRecordRaw;
