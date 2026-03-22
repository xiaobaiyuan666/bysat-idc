import { createRouter, createWebHistory } from "vue-router";

import businessRoutes from "@/router/modules/business";
import cloudRoutes from "@/router/modules/cloud";
import financeRoutes from "@/router/modules/finance";
import platformRoutes from "@/router/modules/platform";
import workbenchRoutes from "@/router/modules/workbench";
import PureAdminLayout from "@/layout/index.vue";
import { pinia } from "@/stores/pinia";
import { useAuthStore } from "@/stores/auth";
import { useMultiTagsStore } from "@/stores/modules/multiTags";

function hasRouteAccess(user: { permissions: string[] } | null, permissions?: unknown) {
  if (!permissions || !Array.isArray(permissions) || permissions.length === 0) {
    return true;
  }

  if (!user) {
    return false;
  }

  return permissions.some((permission) => user.permissions.includes(String(permission)));
}

const consoleRoutes = [
  ...workbenchRoutes,
  ...businessRoutes,
  ...financeRoutes,
  ...cloudRoutes,
  ...platformRoutes,
];

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      name: "login",
      component: () => import("@/pages/LoginPage.vue"),
      meta: {
        public: true,
      },
    },
    {
      path: "/",
      component: PureAdminLayout,
      children: [
        {
          path: "",
          redirect: "/workbench",
        },
        {
          path: "dashboard",
          redirect: "/workbench",
          meta: {
            nav: false,
          },
        },
        ...consoleRoutes,
      ],
    },
  ],
});

router.beforeEach(async (to) => {
  const authStore = useAuthStore(pinia);
  const user = await authStore.ensureUser();

  if (to.meta.public) {
    if (user && to.path === "/login") {
      return "/workbench";
    }
    return true;
  }

  if (!user) {
    return "/login";
  }

  if (!hasRouteAccess(user, to.meta.permissions)) {
    return "/workbench";
  }

  return true;
});

router.afterEach((to) => {
  if (to.meta.public || !to.meta.title) {
    return;
  }

  const tagsStore = useMultiTagsStore(pinia);
  tagsStore.addTag({
    path: to.path,
    title: String(to.meta.title),
    activePath: typeof to.meta.activeMenu === "string" ? to.meta.activeMenu : undefined,
    closable: to.path !== "/workbench",
  });
});

export { router };
