"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Bell,
  Boxes,
  FileText,
  LayoutDashboard,
  LogOut,
  Receipt,
  Server,
  Tickets,
  Wallet,
} from "lucide-react";

import { portalLogoutAction } from "@/actions/portal-auth";
import { Button } from "@/components/ui/button";
import { PORTAL_NAV_ITEMS, PORTAL_ROLE_LABELS } from "@/lib/constants";
import { formatCurrency } from "@/lib/format";
import { cn } from "@/lib/utils";

const iconMap = {
  "/portal": LayoutDashboard,
  "/portal/store": Boxes,
  "/portal/orders": FileText,
  "/portal/services": Server,
  "/portal/invoices": Receipt,
  "/portal/wallet": Wallet,
  "/portal/notifications": Bell,
  "/portal/tickets": Tickets,
} as const;

function getCurrentNavLabel(pathname: string) {
  return (
    PORTAL_NAV_ITEMS.find(
      (item) =>
        item.href === "/portal"
          ? pathname === item.href
          : pathname === item.href || pathname.startsWith(`${item.href}/`),
    )?.label ?? "客户控制台"
  );
}

export function PortalShell({
  user,
  customer,
  children,
}: {
  user: {
    name: string;
    email: string;
    role: string;
  };
  customer: {
    name: string;
    customerNo: string;
    creditBalance: number;
  };
  children: React.ReactNode;
}) {
  const pathname = usePathname();
  const currentNavLabel = getCurrentNavLabel(pathname);

  return (
    <div className="portal-page-shell">
      <div className="mx-auto grid min-h-screen max-w-[1680px] gap-4 px-4 py-4 xl:grid-cols-[268px_minmax(0,1fr)]">
        <aside className="rounded-[28px] bg-[linear-gradient(180deg,var(--sidebar-bg)_0%,var(--sidebar-bg-2)_100%)] p-4 text-white shadow-[0_18px_45px_rgba(18,27,44,0.18)]">
          <div className="flex items-center gap-3 rounded-[22px] px-3 py-2">
            <div className="grid h-11 w-11 place-items-center rounded-2xl bg-[linear-gradient(135deg,#2f80ff,#35c6a7)] text-sm font-extrabold">
              ID
            </div>
            <div>
              <p className="text-sm font-semibold">IDC 云业务系统</p>
              <p className="mt-1 text-xs text-white/55">客户门户与业务控制台</p>
            </div>
          </div>

          <div className="mt-4 rounded-[24px] border border-white/10 bg-white/6 p-5">
            <p className="text-xs uppercase tracking-[0.24em] text-white/45">账户概览</p>
            <p className="mt-3 text-lg font-semibold">{customer.name}</p>
            <p className="mt-1 text-sm text-white/60">{customer.customerNo}</p>
            <div className="mt-4 rounded-2xl bg-white/8 px-4 py-4">
              <p className="text-xs text-white/48">可用余额</p>
              <p className="mt-2 text-2xl font-semibold">
                {formatCurrency(customer.creditBalance)}
              </p>
              <p className="mt-2 text-xs leading-5 text-white/50">
                可用于支付账单、续费服务和处理门店订单。
              </p>
            </div>
          </div>

          <nav className="mt-5 space-y-2">
            {PORTAL_NAV_ITEMS.map((item) => {
              const Icon = iconMap[item.href];
              const active =
                item.href === "/portal"
                  ? pathname === item.href
                  : pathname === item.href || pathname.startsWith(`${item.href}/`);

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm font-medium transition",
                    active
                      ? "bg-white shadow-[0_10px_24px_rgba(18,27,44,0.12)]"
                      : "text-white/74 hover:bg-white/8 hover:text-white",
                  )}
                  style={active ? { color: "var(--ink)" } : undefined}
                >
                  <Icon className="h-4 w-4" />
                  <span>{item.label}</span>
                </Link>
              );
            })}
          </nav>

          <div className="mt-5 rounded-[24px] border border-white/10 bg-white/6 p-5">
            <p className="text-xs uppercase tracking-[0.24em] text-white/45">当前登录</p>
            <p className="mt-3 text-base font-semibold">{user.name}</p>
            <p className="mt-1 text-sm text-white/65">
              {PORTAL_ROLE_LABELS[user.role] ?? user.role}
            </p>
            <p className="mt-1 text-xs text-white/45">{user.email}</p>

            <form action={portalLogoutAction} className="mt-4">
              <Button
                tone="ghost"
                className="w-full justify-center gap-2 border-white/12 bg-white/6 text-white hover:bg-white/10"
              >
                <LogOut className="h-4 w-4" />
                退出门户
              </Button>
            </form>
          </div>
        </aside>

        <main className="min-w-0">
          <header className="rounded-[28px] border border-[var(--border)] bg-white/72 px-6 py-5 shadow-[0_10px_24px_rgba(18,27,44,0.04)] backdrop-blur">
            <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
              <div>
                <p className="text-xs uppercase tracking-[0.24em] text-[var(--muted-ink)]">
                  云产品目录 / 财务结算 / 服务控制 / 工单协同
                </p>
                <h2 className="mt-2 text-[28px] font-bold tracking-[-0.03em] text-[var(--ink)]">
                  {currentNavLabel}
                </h2>
                <p className="mt-2 text-sm leading-7 text-[var(--muted-ink)]">
                  面向客户的统一云业务入口，支持下单、续费、资源查看、财务支付和工单协作。
                </p>
              </div>

              <div className="flex flex-col gap-3 sm:flex-row">
                <div className="rounded-[20px] border border-[var(--border)] bg-white px-4 py-3">
                  <p className="text-xs uppercase tracking-[0.2em] text-[var(--muted-ink)]">
                    客户编号
                  </p>
                  <p className="mt-2 text-sm font-semibold text-[var(--ink)]">
                    {customer.customerNo}
                  </p>
                </div>
                <div className="rounded-[20px] bg-[var(--accent-soft)] px-4 py-3">
                  <p className="text-xs uppercase tracking-[0.2em] text-[var(--accent)]">
                    可用余额
                  </p>
                  <p className="mt-2 text-lg font-semibold text-[var(--ink)]">
                    {formatCurrency(customer.creditBalance)}
                  </p>
                </div>
              </div>
            </div>
          </header>

          <div className="mt-4 space-y-6">{children}</div>
        </main>
      </div>
    </div>
  );
}
