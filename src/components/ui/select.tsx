import { type SelectHTMLAttributes } from "react";

import { cn } from "@/lib/utils";

export function Select({
  className,
  children,
  ...props
}: SelectHTMLAttributes<HTMLSelectElement>) {
  return (
    <select
      className={cn(
        "w-full rounded-2xl border border-[var(--border)] bg-[var(--panel-strong)] px-4 py-3 text-sm text-[var(--ink)] outline-none transition focus:border-[var(--accent)] focus:ring-4 focus:ring-[rgba(36,104,242,0.12)]",
        className,
      )}
      {...props}
    >
      {children}
    </select>
  );
}
