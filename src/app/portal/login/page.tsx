import Link from "next/link";

import { PortalLoginForm } from "@/components/auth/portal-login-form";
import { Card } from "@/components/ui/card";
import { getAdminConsoleUrl } from "@/lib/admin-console";

export default async function PortalLoginPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const params = await searchParams;
  const hasError = Boolean(params.error);
  const adminConsoleUrl = getAdminConsoleUrl("/login");

  return (
    <main className="portal-page-shell px-4 py-8">
      <div className="mx-auto grid max-w-6xl gap-6 lg:grid-cols-[1.1fr_0.9fr]">
        <section className="rounded-[28px] bg-[linear-gradient(180deg,var(--sidebar-bg)_0%,var(--sidebar-bg-2)_100%)] p-8 text-white shadow-[0_18px_45px_rgba(18,27,44,0.18)]">
          <div className="flex items-center gap-3">
            <div className="grid h-12 w-12 place-items-center rounded-2xl bg-[linear-gradient(135deg,#2f80ff,#35c6a7)] font-extrabold">
              ID
            </div>
            <div>
              <p className="text-sm font-semibold">IDC 云业务系统</p>
              <p className="mt-1 text-xs text-white/55">客户门户与自助服务平台</p>
            </div>
          </div>

          <h1 className="mt-8 text-4xl font-semibold tracking-[-0.03em]">
            面向商用运营的客户自助云平台
          </h1>
          <p className="mt-5 max-w-2xl text-base leading-8 text-white/72">
            客户可直接完成云产品选购、订单支付、服务实例查看、余额结算、账单处理和工单协同。
            门户端与运营后台共用同一套业务、计费和资源模型。
          </p>

          <div className="mt-8 grid gap-4 md:grid-cols-2">
            <div className="rounded-[24px] border border-white/10 bg-white/6 p-5">
              <p className="text-sm font-semibold">门户能力</p>
              <ul className="mt-3 space-y-2 text-sm leading-6 text-white/68">
                <li>云产品浏览与在线下单</li>
                <li>账单支付与余额流水查询</li>
                <li>实例详情、重装、快照与备份管理</li>
                <li>通知中心与多轮工单协同</li>
              </ul>
            </div>
            <div className="rounded-[24px] border border-white/10 bg-white/6 p-5">
              <p className="text-sm font-semibold">演示账号</p>
              <ul className="mt-3 space-y-2 text-sm leading-6 text-white/68">
                <li>
                  <code>ops@stargalaxy.cn</code> / <code>Portal123!</code>
                </li>
                <li>
                  <code>finance@cloudmatrix.cn</code> / <code>Portal123!</code>
                </li>
              </ul>
            </div>
          </div>

          <div className="mt-8 text-sm text-white/60">
            新版运营后台入口:
            <Link
              href={adminConsoleUrl}
              className="ml-2 font-medium underline underline-offset-4"
            >
              {adminConsoleUrl}
            </Link>
          </div>
        </section>

        <Card className="self-center p-8">
          <div>
            <p className="text-xs uppercase tracking-[0.24em] text-[var(--muted-ink)]">
              Customer Access
            </p>
            <h2 className="mt-3 text-3xl font-semibold tracking-[-0.03em] text-[var(--ink)]">
              登录客户控制台
            </h2>
            <p className="mt-3 text-sm leading-7 text-[var(--muted-ink)]">
              使用客户门户联系人账号登录。门户与后台采用独立会话，适合分角色运营。
            </p>
          </div>

          <div className="mt-8">
            <PortalLoginForm hasError={hasError} />
          </div>
        </Card>
      </div>
    </main>
  );
}
