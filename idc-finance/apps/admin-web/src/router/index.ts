import { createRouter, createWebHistory } from "vue-router";
import { useUserStore, usePermissionStore, useMultiTagsStore } from "@/store";
import workbench from "./modules/workbench";
import customer from "./modules/customer";
import catalog from "./modules/catalog";
import order from "./modules/order";
import billing from "./modules/billing";
import service from "./modules/service";
import provider from "./modules/provider";
import report from "./modules/report";
import system from "./modules/system";
import remaining from "./modules/remaining";

const router = createRouter({
  history: createWebHistory(),
  routes: [workbench, customer, catalog, order, billing, service, provider, report, system, ...remaining]
});

router.beforeEach(async to => {
  const userStore = useUserStore();
  const permissionStore = usePermissionStore();

  if (to.path !== "/login" && !userStore.isLoggedIn) {
    return "/login";
  }

  if (to.path === "/login" && userStore.isLoggedIn) {
    return "/workbench";
  }

  if (userStore.isLoggedIn && permissionStore.menus.length === 0) {
    await permissionStore.load();
  }

  return true;
});

router.afterEach(to => {
  if (to.meta?.title && to.path !== "/login") {
    useMultiTagsStore().push({
      title: String(to.meta.title),
      titleEn: typeof to.meta.titleEn === "string" ? String(to.meta.titleEn) : undefined,
      path: to.path
    });
  }
});

export default router;
