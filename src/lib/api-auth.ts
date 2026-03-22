import { type AdminRole, type AdminUser } from "@prisma/client";
import { NextRequest } from "next/server";

import { db } from "@/lib/db";
import { AUTH_COOKIE_NAME, verifySessionToken } from "@/lib/session";

const ROLE_PERMISSIONS = {
  SUPER_ADMIN: [
    "dashboard.view",
    "customers.view",
    "customers.manage",
    "products.view",
    "products.manage",
    "orders.view",
    "orders.manage",
    "services.view",
    "services.manage",
    "billing.view",
    "billing.manage",
    "resources.view",
    "resources.manage",
    "invoices.view",
    "invoices.manage",
    "payments.view",
    "payments.manage",
    "payments.callback",
    "notifications.view",
    "notifications.manage",
    "audits.view",
    "tickets.view",
    "tickets.manage",
    "tickets.reply",
  ],
  FINANCE: [
    "dashboard.view",
    "customers.view",
    "customers.manage",
    "products.view",
    "products.manage",
    "orders.view",
    "orders.manage",
    "services.view",
    "services.manage",
    "billing.view",
    "billing.manage",
    "invoices.view",
    "invoices.manage",
    "payments.view",
    "payments.manage",
    "payments.callback",
    "notifications.view",
    "notifications.manage",
    "audits.view",
    "tickets.view",
    "tickets.manage",
    "tickets.reply",
  ],
  SUPPORT: [
    "dashboard.view",
    "services.view",
    "tickets.view",
    "notifications.view",
    "tickets.manage",
    "tickets.reply",
  ],
  OPERATIONS: [
    "dashboard.view",
    "products.view",
    "orders.view",
    "orders.manage",
    "services.view",
    "services.manage",
    "billing.view",
    "billing.manage",
    "resources.view",
    "resources.manage",
    "invoices.view",
    "notifications.view",
    "notifications.manage",
    "audits.view",
    "tickets.view",
    "tickets.manage",
    "tickets.reply",
  ],
} as const satisfies Record<AdminRole, readonly string[]>;

export type Permission = (typeof ROLE_PERMISSIONS)[AdminRole][number];

export async function getApiUser(request: NextRequest) {
  const token = request.cookies.get(AUTH_COOKIE_NAME)?.value;

  if (!token) {
    return null;
  }

  const session = await verifySessionToken(token);

  if (!session) {
    return null;
  }

  const user = await db.adminUser.findUnique({
    where: {
      id: session.sub,
    },
  });

  if (!user || !user.isActive) {
    return null;
  }

  return user;
}

export function getPermissionsForRole(role: AdminRole) {
  return [...ROLE_PERMISSIONS[role]];
}

export function hasAnyRole(user: Pick<AdminUser, "role">, roles: AdminRole[]) {
  return roles.includes(user.role);
}

export function hasPermission(
  user: Pick<AdminUser, "role">,
  permission: Permission,
) {
  return (ROLE_PERMISSIONS[user.role] as readonly string[]).includes(permission);
}

export function hasAnyPermission(
  user: Pick<AdminUser, "role">,
  permissions: Permission[],
) {
  return permissions.some((permission) => hasPermission(user, permission));
}

export function canManageBilling(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "billing.manage");
}

export function canViewBilling(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "billing.view");
}

export function canViewDashboard(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "dashboard.view");
}

export function canViewAudits(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "audits.view");
}

export function canManageCustomers(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "customers.manage");
}

export function canViewCustomers(user: Pick<AdminUser, "role">) {
  return hasAnyPermission(user, ["customers.view", "customers.manage"]);
}

export function canManagePayments(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "payments.manage");
}

export function canViewNotifications(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "notifications.view");
}

export function canManageNotifications(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "notifications.manage");
}

export function canManageProducts(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "products.manage");
}

export function canManageOrders(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "orders.manage");
}

export function canViewOrders(user: Pick<AdminUser, "role">) {
  return hasAnyPermission(user, ["orders.view", "orders.manage"]);
}

export function canViewInvoices(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "invoices.view");
}

export function canManageInvoices(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "invoices.manage");
}

export function canManageServices(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "services.manage");
}

export function canViewServices(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "services.view");
}

export function canViewResources(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "resources.view");
}

export function canManageResources(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "resources.manage");
}

export function canManageTickets(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "tickets.manage");
}

export function canReplyTickets(user: Pick<AdminUser, "role">) {
  return hasPermission(user, "tickets.reply");
}

export function serializeAuthUser(user: Pick<AdminUser, "id" | "name" | "email" | "role">) {
  return {
    id: user.id,
    name: user.name,
    email: user.email,
    role: user.role,
    permissions: getPermissionsForRole(user.role),
  };
}
