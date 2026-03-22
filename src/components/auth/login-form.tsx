import { ShieldCheck } from "lucide-react";

import { signInAction } from "@/actions/auth";
import { SubmitButton } from "@/components/ui/submit-button";
import { Input } from "@/components/ui/input";

export function LoginForm({ hasError }: { hasError: boolean }) {
  return (
    <form action={signInAction} className="space-y-5">
      <div className="inline-flex items-center gap-2 rounded-full bg-[var(--accent-soft)] px-4 py-2 text-sm font-semibold text-[var(--ink)]">
        <ShieldCheck className="h-4 w-4" />
        管理后台登录
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium text-[var(--ink)]">邮箱</label>
        <Input name="email" type="email" placeholder="admin@idc.local" required />
      </div>

      <div className="space-y-2">
        <label className="text-sm font-medium text-[var(--ink)]">密码</label>
        <Input name="password" type="password" placeholder="请输入密码" required />
      </div>

      {hasError ? (
        <div className="rounded-2xl bg-[rgba(185,76,55,0.12)] px-4 py-3 text-sm text-[#9d3e2d]">
          登录失败，请检查邮箱或密码。
        </div>
      ) : null}

      <SubmitButton className="w-full">进入系统</SubmitButton>
    </form>
  );
}
