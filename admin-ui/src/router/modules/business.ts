import { Memo, Monitor, Tickets, User } from "@element-plus/icons-vue";
import type { RouteRecordRaw } from "vue-router";

const businessRoutes: RouteRecordRaw[] = [
  {
    path: "customers",
    name: "customers",
    component: () => import("@/pages/CustomersPage.vue"),
    meta: {
      title: "客户列表",
      subtitle: "统一查看客户资料、账户状态、余额、联系人和客户全景。",
      icon: User,
      permissions: ["customers.manage"],
      navGroup: "business",
      navOrder: 10,
    },
  },
  {
    path: "customers/:id",
    name: "customer-detail",
    component: () => import("@/pages/CustomerDetailPage.vue"),
    meta: {
      title: "客户详情",
      subtitle: "查看客户摘要、订单、服务、账单、流水、工单和审计记录。",
      permissions: ["customers.manage"],
      nav: false,
      activeMenu: "/customers",
    },
  },
  {
    path: "orders",
    name: "orders",
    component: () => import("@/pages/OrdersPage.vue"),
    meta: {
      title: "订单中心",
      subtitle: "管理新购、续费、来源渠道、订单金额、账单联动和服务交付。",
      icon: Memo,
      permissions: ["orders.manage"],
      navGroup: "business",
      navOrder: 20,
    },
  },
  {
    path: "orders/:id",
    name: "order-detail",
    component: () => import("@/pages/OrderDetailPage.vue"),
    meta: {
      title: "订单详情",
      subtitle: "查看订单项目、账单、支付、交付服务和订单时间线。",
      permissions: ["orders.manage"],
      nav: false,
      activeMenu: "/orders",
    },
  },
  {
    path: "services",
    name: "services",
    component: () => import("@/pages/ServicesPage.vue"),
    meta: {
      title: "业务列表",
      subtitle: "管理实例生命周期、魔方云动作、资源统计和同步状态。",
      icon: Monitor,
      permissions: ["services.view"],
      navGroup: "business",
      navOrder: 30,
    },
  },
  {
    path: "services/:id",
    name: "service-detail",
    component: () => import("@/pages/ServiceDetailPage.vue"),
    meta: {
      title: "服务详情",
      subtitle: "查看实例摘要、资源清单、账单、工单、同步日志和审计轨迹。",
      permissions: ["services.view"],
      nav: false,
      activeMenu: "/services",
    },
  },
  {
    path: "tickets",
    name: "tickets",
    component: () => import("@/pages/TicketsPage.vue"),
    meta: {
      title: "工单中心",
      subtitle: "处理客户工单、回复记录、内部备注和工单流转。",
      icon: Tickets,
      permissions: ["tickets.manage"],
      navGroup: "business",
      navOrder: 40,
    },
  },
];

export default businessRoutes;
