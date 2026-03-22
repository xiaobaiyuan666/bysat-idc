"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import { LayoutDashboard, Receipt, Server, ShoppingCart, Tickets, Users, Wallet, Boxes } from "lucide-react";

import { logoutAction } from "@/actions/auth";
import { Button } from "@/components/ui/button";
import { APP_NAME, NAV_ITEMS, ROLE_LABELS } from "@/lib/constants";
import { cn } from "@/lib/utils";

const iconMap = {
  "/dashboard": LayoutDashboard,
  "/customers": Users,
  "/products": Boxes,
  "/orders": ShoppingCart,
  "/services": Server,
  "/invoices": Receipt,
  "/payments": Wallet,
  "/tickets": Tickets,
} as const;

export function AppShell({
  user,
  children,
}: {
  user: {
    name: string;
    email: string;
    role: string;
  };
  children: React.ReactNode;
}) {
  const pathname = usePathname();

  return (
    <div className="min-h-screen bg-transparent">
      <div className="mx-auto grid min-h-screen max-w-[1600px] gap-6 px-4 py-4 lg:grid-cols-[280px_minmax(0,1fr)]">
        <aside className="rounded-[32px] bg-[linear-gradient(180deg,#132238_0%,#1d334f_100%)] p-6 text-white shadow-[0_32px_80px_rgba(19,34,56,0.24)]">
          <div className="mb-8">
            <p className="text-xs uppercase tracking-[0.32em] text-white/50">IDC Finance</p>
            <h1 className="mt-3 text-2xl font-semibold tracking-tight">{APP_NAME}</h1>
            <p className="mt-2 text-sm text-white/70">
              参考魔方财务的业务链路，预留魔方云资源同步接口。
            </p>
          </div>

          <nav className="space-y-2">
            {NAV_ITEMS.map((item) => {
              const Icon = iconMap[item.href];
              const active = pathname === item.href;

              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={cn(
                    "flex items-center gap-3 rounded-2xl px-4 py-3 text-sm transition",
                    active
                      ? "bg-white text-[var(--ink)]"
                      : "text-white/75 hover:bg-white/10 hover:text-white",
                  )}
                >
                  <Icon className="h-4 w-4" />
                  <span>{item.label}</span>
                </Link>
              );
            })}
          </nav>

          <div className="mt-8 rounded-[24px] border border-white/10 bg-white/8 p-4">
            <p className="text-xs uppercase tracking-[0.24em] text-white/50">当前账号</p>
            <p className="mt-3 text-lg font-semibold">{user.name}</p>
            <p className="text-sm text-white/65">{ROLE_LABELS[user.role] ?? user.role}</p>
            <p className="mt-1 text-xs text-white/50">{user.email}</p>

            <form action={logoutAction} className="mt-4">
              <Button tone="secondary" className="w-full">
                退出登录
              </Button>
            </form>
          </div>
        </aside>

        <main className="min-w-0">
          <header className="mb-6 flex flex-col gap-4 rounded-[32px] border border-[var(--border)] bg-[rgba(255,253,248,0.78)] px-6 py-5 backdrop-blur md:flex-row md:items-center md:justify-between">
            <div>
              <p className="text-xs uppercase tracking-[0.28em] text-[var(--muted-ink)]">
                财务 / 运营 / 资源同步
              </p>
              <h2 className="mt-2 text-2xl font-semibold tracking-tight text-[var(--ink)]">
                中小型 IDC 财务管理后台
              </h2>
            </div>
            <div className="rounded-2xl bg-[var(--accent-soft)] px-4 py-3 text-sm text-[var(--ink)]">
              魔方云适配层默认启用 Mock，可在环境变量中切换真实 API。
            </div>
          </header>

          <div className="space-y-6">{children}</div>
        </main>
      </div>
    </div>
  );
}
