import type { RouteRecordRaw } from "vue-router";

const routes: RouteRecordRaw[] = [
  {
    path: "/login",
    name: "login",
    component: () => import("@/views/login/index.vue"),
    meta: { title: "\u767b\u5f55", titleEn: "Login" }
  },
  {
    path: "/:pathMatch(.*)*",
    redirect: "/workbench"
  }
];

export default routes;
