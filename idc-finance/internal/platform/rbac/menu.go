package rbac

type Menu struct {
	ID         int64  `json:"id"`
	ParentID   int64  `json:"parentId"`
	Title      string `json:"title"`
	TitleEn    string `json:"titleEn"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	Icon       string `json:"icon"`
	Permission string `json:"permission"`
}

func PhaseOneMenus() []Menu {
	return []Menu{
		{ID: 1, ParentID: 0, Title: "\u5de5\u4f5c\u53f0", TitleEn: "Workbench", Name: "workbench", Path: "/workbench", Icon: "Monitor", Permission: "workbench:view"},

		{ID: 10, ParentID: 0, Title: "\u5ba2\u6237", TitleEn: "Customers", Name: "customer", Path: "/customer/list", Icon: "User", Permission: "customer:list"},
		{ID: 11, ParentID: 10, Title: "\u5ba2\u6237\u5217\u8868", TitleEn: "Customer List", Name: "customer-list", Path: "/customer/list", Icon: "UserFilled", Permission: "customer:list"},
		{ID: 12, ParentID: 10, Title: "\u5b9e\u540d\u8ba4\u8bc1", TitleEn: "Identity Review", Name: "customer-identities", Path: "/customer/identities", Icon: "Postcard", Permission: "customer:identity:view"},
		{ID: 13, ParentID: 10, Title: "\u5ba2\u6237\u5206\u7ec4", TitleEn: "Customer Groups", Name: "customer-groups", Path: "/customer/groups", Icon: "Collection", Permission: "customer:update"},
		{ID: 14, ParentID: 10, Title: "\u5ba2\u6237\u7b49\u7ea7", TitleEn: "Customer Levels", Name: "customer-levels", Path: "/customer/levels", Icon: "Medal", Permission: "customer:update"},

		{ID: 20, ParentID: 0, Title: "\u4e1a\u52a1", TitleEn: "Business", Name: "business", Path: "/orders/list", Icon: "Tickets", Permission: "order:view"},
		{ID: 21, ParentID: 20, Title: "\u4ea7\u54c1\u8ba2\u5355", TitleEn: "Orders", Name: "order-list", Path: "/orders/list", Icon: "List", Permission: "order:view"},
		{ID: 22, ParentID: 20, Title: "\u4e1a\u52a1\u5217\u8868", TitleEn: "Service List", Name: "service-list", Path: "/services/list", Icon: "DataLine", Permission: "service:view"},
		{ID: 24, ParentID: 20, Title: "\u8ba2\u5355\u7533\u8bf7", TitleEn: "Order Requests", Name: "order-requests", Path: "/orders/requests", Icon: "Bell", Permission: "order:view"},
		{ID: 23, ParentID: 20, Title: "\u6539\u914d\u5355", TitleEn: "Change Orders", Name: "service-change-orders", Path: "/orders/change-orders", Icon: "Files", Permission: "order:view"},

		{ID: 30, ParentID: 0, Title: "\u8d22\u52a1", TitleEn: "Finance", Name: "finance", Path: "/billing/invoices", Icon: "Coin", Permission: "invoice:view"},
		{ID: 31, ParentID: 30, Title: "\u8d26\u5355\u7ba1\u7406", TitleEn: "Invoices", Name: "billing-invoices", Path: "/billing/invoices", Icon: "DocumentChecked", Permission: "invoice:view"},
		{ID: 34, ParentID: 30, Title: "\u652f\u4ed8\u8bb0\u5f55", TitleEn: "Payments", Name: "billing-payments", Path: "/billing/payments", Icon: "Wallet", Permission: "invoice:view"},
		{ID: 35, ParentID: 30, Title: "\u9000\u6b3e\u8bb0\u5f55", TitleEn: "Refunds", Name: "billing-refunds", Path: "/billing/refunds", Icon: "RefreshLeft", Permission: "invoice:view"},
		{ID: 36, ParentID: 30, Title: "\u8d44\u91d1\u53f0\u8d26", TitleEn: "Account Ledger", Name: "billing-accounts", Path: "/billing/accounts", Icon: "CreditCard", Permission: "invoice:view"},
		{ID: 37, ParentID: 30, Title: "\u5145\u503c\u8bb0\u5f55", TitleEn: "Recharges", Name: "billing-recharges", Path: "/billing/recharges", Icon: "Wallet", Permission: "invoice:view"},
		{ID: 38, ParentID: 30, Title: "\u8d44\u91d1\u8c03\u6574", TitleEn: "Adjustments", Name: "billing-adjustments", Path: "/billing/adjustments", Icon: "Tools", Permission: "invoice:view"},
		{ID: 32, ParentID: 30, Title: "\u5de5\u5355\u4e2d\u5fc3", TitleEn: "Tickets", Name: "ticket-list", Path: "/tickets/list", Icon: "ChatDotRound", Permission: "ticket:view"},
		{ID: 33, ParentID: 30, Title: "\u5de5\u5355\u914d\u7f6e", TitleEn: "Ticket Settings", Name: "ticket-settings", Path: "/tickets/settings", Icon: "Tools", Permission: "ticket:update"},

		{ID: 40, ParentID: 0, Title: "\u5546\u54c1\u8bbe\u7f6e", TitleEn: "Catalog", Name: "product-settings", Path: "/catalog/products", Icon: "Goods", Permission: "catalog:product:view"},
		{ID: 41, ParentID: 40, Title: "\u5546\u54c1\u7ba1\u7406", TitleEn: "Products", Name: "catalog-products", Path: "/catalog/products", Icon: "Box", Permission: "catalog:product:view"},

		{ID: 50, ParentID: 0, Title: "\u63a5\u53e3\u4e0e\u4e0a\u6e38", TitleEn: "Providers & Upstream", Name: "resource-store", Path: "/providers/accounts", Icon: "Connection", Permission: "provider:view"},
		{ID: 54, ParentID: 50, Title: "\u63a5\u53e3\u8d26\u6237", TitleEn: "Provider Accounts", Name: "provider-accounts", Path: "/providers/accounts", Icon: "Connection", Permission: "provider:view"},
		{ID: 55, ParentID: 50, Title: "\u4e0a\u6e38\u540c\u6b65\u8bb0\u5f55", TitleEn: "Upstream Sync History", Name: "provider-upstream-sync-history", Path: "/providers/upstream-sync-history", Icon: "DocumentCopy", Permission: "provider:view"},
		{ID: 51, ParentID: 50, Title: "\u6e20\u9053\u8d44\u6e90", TitleEn: "Channel Resources", Name: "provider-resources", Path: "/providers/resources", Icon: "DataBoard", Permission: "provider:view"},
		{ID: 52, ParentID: 50, Title: "\u81ea\u52a8\u5316\u4efb\u52a1", TitleEn: "Automation Tasks", Name: "provider-automation", Path: "/providers/automation", Icon: "Operation", Permission: "automation:view"},
		{ID: 53, ParentID: 50, Title: "\u81ea\u52a8\u5316\u7b56\u7565", TitleEn: "Automation Policies", Name: "provider-automation-settings", Path: "/providers/automation-settings", Icon: "Tools", Permission: "automation:update"},

		{ID: 60, ParentID: 0, Title: "\u62a5\u8868\u4e2d\u5fc3", TitleEn: "Reports", Name: "function-center", Path: "/reports/overview", Icon: "PieChart", Permission: "report:view"},
		{ID: 61, ParentID: 60, Title: "\u8fd0\u8425\u62a5\u8868", TitleEn: "Operations Reports", Name: "report-overview", Path: "/reports/overview", Icon: "Histogram", Permission: "report:view"},

		{ID: 70, ParentID: 0, Title: "\u7cfb\u7edf", TitleEn: "Settings", Name: "settings", Path: "/system/admins", Icon: "Setting", Permission: "rbac:role:view"},
		{ID: 71, ParentID: 70, Title: "\u7ba1\u7406\u5458", TitleEn: "Admins", Name: "system-admins", Path: "/system/admins", Icon: "Avatar", Permission: "rbac:role:view"},
		{ID: 72, ParentID: 70, Title: "\u89d2\u8272\u6743\u9650", TitleEn: "Roles & Permissions", Name: "system-roles", Path: "/system/roles", Icon: "Lock", Permission: "rbac:role:view"},
		{ID: 73, ParentID: 70, Title: "\u83dc\u5355\u7ba1\u7406", TitleEn: "Menus", Name: "system-menus", Path: "/system/menus", Icon: "Menu", Permission: "rbac:menu:view"},
		{ID: 74, ParentID: 70, Title: "\u5ba1\u8ba1\u65e5\u5fd7", TitleEn: "Audit Logs", Name: "system-audit", Path: "/system/audit", Icon: "Document", Permission: "audit:view"},
	}
}
