import { type ButtonHTMLAttributes } from "react";

import { cn } from "@/lib/utils";

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
  tone?: "primary" | "secondary" | "danger" | "ghost";
};

const toneClassMap: Record<NonNullable<ButtonProps["tone"]>, string> = {
  primary:
    "bg-[var(--accent)] text-white shadow-[0_14px_34px_rgba(36,104,242,0.24)] hover:bg-[var(--accent-strong)]",
  secondary:
    "bg-[var(--accent-soft)] text-[var(--accent)] hover:bg-[rgba(36,104,242,0.16)]",
  danger: "bg-[var(--danger)] text-white hover:bg-[#da3a2f]",
  ghost:
    "border border-[var(--border)] bg-[var(--panel-strong)] text-[var(--muted-ink)] hover:bg-[var(--panel-soft)]",
};

export function Button({
  className,
  tone = "primary",
  type = "button",
  ...props
}: ButtonProps) {
  return (
    <button
      type={type}
      className={cn(
        "inline-flex items-center justify-center rounded-2xl px-4 py-2.5 text-sm font-semibold transition duration-200 disabled:cursor-not-allowed disabled:opacity-60",
        toneClassMap[tone],
        className,
      )}
      {...props}
    />
  );
}
