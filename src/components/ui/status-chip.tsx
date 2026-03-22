import { cn } from "@/lib/utils";

const toneMap: Record<string, string> = {
  ACTIVE: "bg-[var(--success-soft)] text-[var(--success)]",
  PAID: "bg-[var(--success-soft)] text-[var(--success)]",
  SUCCESS: "bg-[var(--success-soft)] text-[var(--success)]",
  SENT: "bg-[var(--success-soft)] text-[var(--success)]",
  OPEN: "bg-[var(--warning-soft)] text-[var(--warning)]",
  ISSUED: "bg-[var(--warning-soft)] text-[var(--warning)]",
  PENDING: "bg-[var(--warning-soft)] text-[var(--warning)]",
  PROVISIONING: "bg-[var(--accent-soft)] text-[var(--accent)]",
  PARTIAL: "bg-[var(--accent-soft)] text-[var(--accent)]",
  PROCESSING: "bg-[var(--accent-soft)] text-[var(--accent)]",
  OVERDUE: "bg-[var(--danger-soft)] text-[var(--danger)]",
  SUSPENDED: "bg-[var(--danger-soft)] text-[var(--danger)]",
  FAILED: "bg-[var(--danger-soft)] text-[var(--danger)]",
  REFUNDED: "bg-[var(--danger-soft)] text-[var(--danger)]",
  TERMINATED: "bg-[rgba(24,35,56,0.08)] text-[var(--ink)]",
  CANCELED: "bg-[rgba(24,35,56,0.08)] text-[var(--ink)]",
  CLOSED: "bg-[rgba(24,35,56,0.08)] text-[var(--ink)]",
  RESOLVED: "bg-[rgba(24,35,56,0.08)] text-[var(--ink)]",
};

export function StatusChip({
  label,
  value,
}: {
  label: string;
  value: string;
}) {
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold",
        toneMap[value] ?? "bg-[rgba(19,34,56,0.08)] text-[var(--ink)]",
      )}
    >
      {label}
    </span>
  );
}
