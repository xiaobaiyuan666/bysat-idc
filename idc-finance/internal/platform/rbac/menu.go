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
		{ID: 1, ParentID: 0, Title: "工作台", TitleEn: "Workbench", Name: "workbench", Path: "/workbench", Icon: "Monitor", Permission: "workbench:view"},

		{ID: 10, ParentID: 0, Title: "客户", TitleEn: "Customers", Name: "customer", Path: "/customer/list", Icon: "User", Permission: "customer:list"},
		{ID: 11, ParentID: 10, Title: "客户列表", TitleEn: "Customer List", Name: "customer-list", Path: "/customer/list", Icon: "UserFilled", Permission: "customer:list"},
		{ID: 12, ParentID: 10, Title: "实名认证", TitleEn: "Identity Review", Name: "customer-identities", Path: "/customer/identities", Icon: "Postcard", Permission: "customer:identity:view"},
		{ID: 13, ParentID: 10, Title: "客户分组", TitleEn: "Customer Groups", Name: "customer-groups", Path: "/customer/groups", Icon: "Collection", Permission: "customer:update"},
		{ID: 14, ParentID: 10, Title: "客户等级", TitleEn: "Customer Levels", Name: "customer-levels", Path: "/customer/levels", Icon: "Medal", Permission: "customer:update"},

		{ID: 20, ParentID: 0, Title: "业务", TitleEn: "Business", Name: "business", Path: "/orders/list", Icon: "Tickets", Permission: "order:view"},
		{ID: 21, ParentID: 20, Title: "产品订单", TitleEn: "Orders", Name: "order-list", Path: "/orders/list", Icon: "List", Permission: "order:view"},
		{ID: 22, ParentID: 20, Title: "业务列表", TitleEn: "Service List", Name: "service-list", Path: "/services/list", Icon: "DataLine", Permission: "service:view"},

		{ID: 30, ParentID: 0, Title: "财务", TitleEn: "Finance", Name: "finance", Path: "/billing/invoices", Icon: "Coin", Permission: "invoice:view"},
		{ID: 31, ParentID: 30, Title: "账单管理", TitleEn: "Invoices", Name: "billing-invoices", Path: "/billing/invoices", Icon: "DocumentChecked", Permission: "invoice:view"},

		{ID: 40, ParentID: 0, Title: "商品设置", TitleEn: "Catalog", Name: "product-settings", Path: "/catalog/products", Icon: "Goods", Permission: "catalog:product:view"},
		{ID: 41, ParentID: 40, Title: "商品管理", TitleEn: "Products", Name: "catalog-products", Path: "/catalog/products", Icon: "Box", Permission: "catalog:product:view"},

		{ID: 50, ParentID: 0, Title: "资源与商店", TitleEn: "Resources", Name: "resource-store", Path: "/providers/mofang", Icon: "Connection", Permission: "provider:view"},
		{ID: 54, ParentID: 50, Title: "接口账户", TitleEn: "Provider Accounts", Name: "provider-accounts", Path: "/providers/accounts", Icon: "Connection", Permission: "provider:view"},
		{ID: 51, ParentID: 50, Title: "魔方云", TitleEn: "Mofang Cloud", Name: "provider-mofang", Path: "/providers/mofang", Icon: "Connection", Permission: "provider:view"},
		{ID: 52, ParentID: 50, Title: "自动化任务中心", TitleEn: "Automation Tasks", Name: "provider-automation", Path: "/providers/automation", Icon: "Operation", Permission: "automation:view"},
		{ID: 53, ParentID: 50, Title: "自动化策略", TitleEn: "Automation Policies", Name: "provider-automation-settings", Path: "/providers/automation-settings", Icon: "Tools", Permission: "automation:update"},

		{ID: 60, ParentID: 0, Title: "功能", TitleEn: "Reports", Name: "function-center", Path: "/reports/overview", Icon: "PieChart", Permission: "report:view"},
		{ID: 61, ParentID: 60, Title: "统计报表", TitleEn: "Reports", Name: "report-overview", Path: "/reports/overview", Icon: "Histogram", Permission: "report:view"},

		{ID: 70, ParentID: 0, Title: "设置", TitleEn: "Settings", Name: "settings", Path: "/system/admins", Icon: "Setting", Permission: "rbac:role:view"},
		{ID: 71, ParentID: 70, Title: "管理员", TitleEn: "Admins", Name: "system-admins", Path: "/system/admins", Icon: "Avatar", Permission: "rbac:role:view"},
		{ID: 72, ParentID: 70, Title: "角色权限", TitleEn: "Roles & Permissions", Name: "system-roles", Path: "/system/roles", Icon: "Lock", Permission: "rbac:role:view"},
		{ID: 73, ParentID: 70, Title: "菜单管理", TitleEn: "Menus", Name: "system-menus", Path: "/system/menus", Icon: "Menu", Permission: "rbac:menu:view"},
		{ID: 74, ParentID: 70, Title: "审计日志", TitleEn: "Audit Logs", Name: "system-audit", Path: "/system/audit", Icon: "Document", Permission: "audit:view"},
	}
}
