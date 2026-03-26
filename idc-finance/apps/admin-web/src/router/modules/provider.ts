import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/providers",
  redirect: "/providers/accounts",
  meta: {
    title: "\u63a5\u53e3\u4e0e\u4e0a\u6e38",
    titleEn: "Providers & Upstream",
    requiresAuth: true
  },
  children: [
    {
      path: "accounts",
      name: "provider-accounts",
      component: () => import("@/views/provider/accounts/index.vue"),
      meta: {
        title: "\u63a5\u53e3\u8d26\u6237",
        titleEn: "Provider Accounts",
        requiresAuth: true
      }
    },
    {
      path: "upstream-sync-history",
      name: "provider-upstream-sync-history",
      component: () => import("@/views/provider/upstream-sync-history/index.vue"),
      meta: {
        title: "上游同步记录",
        titleEn: "Upstream Sync History",
        requiresAuth: true
      }
    },
    {
      path: "resources",
      name: "provider-resources",
      component: () => import("@/views/provider/resources/index.vue"),
      meta: {
        title: "\u6e20\u9053\u8d44\u6e90",
        titleEn: "Channel Resources",
        requiresAuth: true
      }
    },
    {
      path: "automation",
      name: "provider-automation",
      component: () => import("@/views/provider/automation/index.vue"),
      meta: {
        title: "\u81ea\u52a8\u5316\u4efb\u52a1\u4e2d\u5fc3",
        titleEn: "Automation Tasks",
        requiresAuth: true
      }
    },
    {
      path: "automation-settings",
      name: "provider-automation-settings",
      component: () => import("@/views/provider/automation-settings/index.vue"),
      meta: {
        title: "\u81ea\u52a8\u5316\u7b56\u7565",
        titleEn: "Automation Policies",
        requiresAuth: true
      }
    },
    {
      path: "mofang",
      name: "provider-mofang",
      redirect: {
        path: "/providers/resources",
        query: {
          providerType: "MOFANG_CLOUD"
        }
      },
      meta: {
        title: "\u6e20\u9053\u8d44\u6e90",
        titleEn: "Channel Resources",
        requiresAuth: true
      }
    }
  ]
} satisfies RouteRecordRaw;
