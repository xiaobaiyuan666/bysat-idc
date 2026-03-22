import { CloudCog } from "lucide-react";

import { portalSignInAction } from "@/actions/portal-auth";
import { Input } from "@/components/ui/input";
import { SubmitButton } from "@/components/ui/submit-button";

export function PortalLoginForm({ hasError }: { hasError: boolean }) {
  return (
    <form action={portalSignInAction} className="space-y-5">
      <div className="inline-flex items-center gap-2 rounded-full bg-[var(--accent-soft)] px-4 py-2 text-sm font-semibold text-[var(--accent)]">
        <CloudCog className="h-4 w-4" />
        客户门户登录
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium text-[var(--ink)]">邮箱账号</label>
        <Input name="email" type="email" placeholder="ops@stargalaxy.cn" required />
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium text-[var(--ink)]">登录密码</label>
        <Input name="password" type="password" placeholder="请输入门户密码" required />
      </div>

      {hasError ? (
        <div className="rounded-2xl bg-[var(--danger-soft)] px-4 py-3 text-sm text-[var(--danger)]">
          登录失败，请检查账号状态或密码是否正确。
        </div>
      ) : null}

      <SubmitButton className="w-full">进入控制台</SubmitButton>
    </form>
  );
}
