import { type TextareaHTMLAttributes } from "react";

import { cn } from "@/lib/utils";

export function Textarea({
  className,
  ...props
}: TextareaHTMLAttributes<HTMLTextAreaElement>) {
  return (
    <textarea
      className={cn(
        "min-h-28 w-full rounded-2xl border border-[var(--border)] bg-[var(--panel-strong)] px-4 py-3 text-sm text-[var(--ink)] outline-none transition placeholder:text-[var(--muted-ink)] focus:border-[var(--accent)] focus:ring-4 focus:ring-[rgba(36,104,242,0.12)]",
        className,
      )}
      {...props}
    />
  );
}
