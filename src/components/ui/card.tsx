import { type HTMLAttributes } from "react";

import { cn } from "@/lib/utils";

export function Card({
  className,
  ...props
}: HTMLAttributes<HTMLDivElement>) {
  return (
    <div
      className={cn(
        "rounded-[24px] border border-[var(--border)] bg-[var(--panel)] p-6 shadow-[0_16px_40px_rgba(18,27,44,0.05)]",
        className,
      )}
      {...props}
    />
  );
}
