import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/customer",
  redirect: "/customer/list",
  meta: {
    title: "客户",
    titleEn: "Customers",
    requiresAuth: true
  },
  children: [
    {
      path: "list",
      name: "customer-list",
      component: () => import("@/views/customer/list/index.vue"),
      meta: { title: "客户列表", titleEn: "Customer List", requiresAuth: true }
    },
    {
      path: "detail/:id",
      name: "customer-detail",
      component: () => import("@/views/customer/detail/index.vue"),
      meta: { title: "客户工作台", titleEn: "Customer Workbench", requiresAuth: true },
      props: true
    },
    {
      path: "groups",
      name: "customer-groups",
      component: () => import("@/views/customer/groups/index.vue"),
      meta: { title: "客户分组", titleEn: "Customer Groups", requiresAuth: true }
    },
    {
      path: "levels",
      name: "customer-levels",
      component: () => import("@/views/customer/levels/index.vue"),
      meta: { title: "客户等级", titleEn: "Customer Levels", requiresAuth: true }
    },
    {
      path: "identities",
      name: "customer-identities",
      component: () => import("@/views/customer/identities/index.vue"),
      meta: { title: "实名认证", titleEn: "Identity Review", requiresAuth: true }
    }
  ]
} satisfies RouteRecordRaw;
