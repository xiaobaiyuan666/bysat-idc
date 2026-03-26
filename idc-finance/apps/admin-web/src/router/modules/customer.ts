import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/customer",
  redirect: "/customer/list",
  meta: {
    title: "\u5ba2\u6237",
    titleEn: "Customers",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "customer-list",
      component: () => import("@/views/customer/list/index.vue"),
      meta: { title: "\u5ba2\u6237\u5217\u8868", titleEn: "Customer List", requiresAuth: true }
    },
    {
      path: "detail/:id",
      name: "customer-detail",
      component: () => import("@/views/customer/detail/index.vue"),
      meta: { title: "\u5ba2\u6237\u5de5\u4f5c\u53f0", titleEn: "Customer Workbench", requiresAuth: true },
      props: true
    },
    {
      path: "groups",
      name: "customer-groups",
      component: () => import("@/views/customer/groups/index.vue"),
      meta: { title: "\u5ba2\u6237\u5206\u7ec4", titleEn: "Customer Groups", requiresAuth: true }
    },
    {
      path: "levels",
      name: "customer-levels",
      component: () => import("@/views/customer/levels/index.vue"),
      meta: { title: "\u5ba2\u6237\u7b49\u7ea7", titleEn: "Customer Levels", requiresAuth: true }
    },
    {
      path: "identities",
      name: "customer-identities",
      component: () => import("@/views/customer/identities/index.vue"),
      meta: { title: "\u5b9e\u540d\u8ba4\u8bc1", titleEn: "Identity Review", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
