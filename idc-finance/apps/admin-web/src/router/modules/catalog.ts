import type { RouteRecordRaw } from "vue-router";

export default {
  path: "/catalog",
  redirect: "/catalog/products",
  meta: {
    title: "\u5546\u54c1\u8bbe\u7f6e",
    titleEn: "Catalog",
    requiresAuth: true
  },
  children: [
    {
      path: "products",
      name: "catalog-products",
      component: () => import("@/views/catalog/products/index.vue"),
      meta: { title: "\u5546\u54c1\u7ba1\u7406", titleEn: "Products", requiresAuth: true }
    },
    {
      path: "products/:id",
      name: "catalog-product-detail",
      component: () => import("@/views/catalog/products/detail/index.vue"),
      meta: { title: "\u5546\u54c1\u5de5\u4f5c\u53f0", titleEn: "Product Workbench", requiresAuth: true, hiddenTag: false }
    }
  ]
} satisfies RouteRecordRaw;
