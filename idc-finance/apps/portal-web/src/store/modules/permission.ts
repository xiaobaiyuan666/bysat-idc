import { defineStore } from "pinia";

export interface MenuItem {
  id: number;
  title: string;
  titleEn?: string;
  path: string;
  children?: MenuItem[];
}

export const usePermissionStore = defineStore("permission", {
  state: () => ({
    menus: [
      { id: 1, title: "控制台", titleEn: "Console", path: "/console" },
      { id: 2, title: "云产品商城", titleEn: "Cloud Store", path: "/store" },
      { id: 3, title: "服务中心", titleEn: "Services", path: "/services" },
      { id: 4, title: "订单中心", titleEn: "Orders", path: "/orders" },
      { id: 5, title: "账单中心", titleEn: "Invoices", path: "/invoices" },
      { id: 6, title: "工单中心", titleEn: "Tickets", path: "/tickets" },
      { id: 7, title: "钱包中心", titleEn: "Wallet", path: "/wallet" },
      { id: 8, title: "账户中心", titleEn: "Account", path: "/account" }
    ] as MenuItem[]
  })
});
