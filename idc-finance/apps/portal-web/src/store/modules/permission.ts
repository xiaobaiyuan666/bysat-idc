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
      { id: 2, title: "服务中心", titleEn: "Services", path: "/services" },
      { id: 3, title: "财务中心", titleEn: "Finance", path: "/invoices" },
      { id: 4, title: "订单中心", titleEn: "Orders", path: "/orders" },
      { id: 5, title: "工单中心", titleEn: "Tickets", path: "/tickets" },
      { id: 6, title: "账户中心", titleEn: "Account", path: "/account" },
      { id: 7, title: "云产品商城", titleEn: "Marketplace", path: "/store" }
    ] as MenuItem[]
  })
});
