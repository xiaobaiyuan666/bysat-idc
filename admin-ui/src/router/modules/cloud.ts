import { Box, Operation, Share } from "@element-plus/icons-vue";
import type { RouteRecordRaw } from "vue-router";

const cloudRoutes: RouteRecordRaw[] = [
  {
    path: "products",
    name: "products",
    component: () => import("@/pages/ProductsPage.vue"),
    meta: {
      title: "商品管理",
      subtitle: "管理商品目录、售价、库存、自动开通和销售状态。",
      icon: Box,
      permissions: ["products.manage"],
      navGroup: "cloud",
      navOrder: 10,
    },
  },
  {
    path: "architecture",
    name: "architecture",
    component: () => import("@/pages/ArchitecturePage.vue"),
    meta: {
      title: "云架构中心",
      subtitle: "维护地域、可用区、规格、镜像、售卖方案和门户账号。",
      icon: Operation,
      permissions: ["resources.manage"],
      navGroup: "cloud",
      navOrder: 20,
    },
  },
  {
    path: "resources",
    name: "resources",
    component: () => import("@/pages/ResourcesPage.vue"),
    meta: {
      title: "资源中心",
      subtitle: "管理 VPC、IP、磁盘、快照、备份和安全组资源。",
      icon: Share,
      permissions: ["resources.view"],
      navGroup: "cloud",
      navOrder: 30,
    },
  },
];

export default cloudRoutes;
