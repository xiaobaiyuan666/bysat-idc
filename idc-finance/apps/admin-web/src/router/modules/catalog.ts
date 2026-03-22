import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/catalog",
  redirect: "/catalog/products",
  meta: {
    title: "商品设置",
    titleEn: "Catalog",
    requiresAuth: true
  },
  children: [
    {
      path: "products",
      name: "catalog-products",
      component: () => import("@/views/catalog/products/index.vue"),
      meta: { title: "商品管理", titleEn: "Products", requiresAuth: true }
    },
    {
      path: "products/:id",
      name: "catalog-product-detail",
      component: () => import("@/views/catalog/products/detail/index.vue"),
      meta: { title: "商品工作台", titleEn: "Product Workbench", requiresAuth: true, hiddenTag: false }
    }
  ]
} satisfies RouteRecordRaw;
