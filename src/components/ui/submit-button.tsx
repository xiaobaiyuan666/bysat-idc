"use client";

import { useFormStatus } from "react-dom";

import { Button } from "@/components/ui/button";

export function SubmitButton({
  children,
  tone = "primary",
  className,
}: {
  children: React.ReactNode;
  tone?: "primary" | "secondary" | "danger" | "ghost";
  className?: string;
}) {
  const { pending } = useFormStatus();

  return (
    <Button type="submit" tone={tone} className={className} disabled={pending}>
      {pending ? "处理中..." : children}
    </Button>
  );
}
