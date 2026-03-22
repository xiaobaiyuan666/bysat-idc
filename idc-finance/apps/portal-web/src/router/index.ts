import { createRouter, createWebHistory } from "vue-router";
import { useUserStore, useMultiTagsStore } from "@/store";
import consoleRoute from "./modules/console";
import storeRoute from "./modules/store";
import servicesRoute from "./modules/services";
import ordersRoute from "./modules/orders";
import invoicesRoute from "./modules/invoices";
import ticketsRoute from "./modules/tickets";
import walletRoute from "./modules/wallet";
import accountRoute from "./modules/account";
import remaining from "./modules/remaining";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    consoleRoute,
    storeRoute,
    servicesRoute,
    ordersRoute,
    invoicesRoute,
    ticketsRoute,
    walletRoute,
    accountRoute,
    ...remaining
  ]
});

router.beforeEach(to => {
  const userStore = useUserStore();

  if (to.path !== "/login" && !userStore.isLoggedIn) {
    return "/login";
  }

  if (to.path === "/login" && userStore.isLoggedIn) {
    return "/console";
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
