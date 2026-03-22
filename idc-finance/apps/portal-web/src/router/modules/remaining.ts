import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "portal-login",
    component: () => import("@/views/login/index.vue"),
    meta: { title: "登录", titleEn: "Login" }
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/console"
  }
];

export default routes;
