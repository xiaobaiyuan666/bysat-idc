import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/login/index.vue"),
    meta: { title: "\u767b\u5f55", titleEn: "Login" }
  },
  {
    path: "/vnc",
    name: "vnc-console",
    component: () => import("@/views/vnc/index.vue"),
    meta: {
      title: "实例控制台",
      titleEn: "Remote Console",
      requiresAuth: true,
      hideLayout: true
    }
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/workbench"
  }
];

export default routes;
