import { db } from "@/lib/db";

function serializeDetail(detail: unknown) {
  if (detail === undefined || detail === null) {
    return undefined;
  }

  if (typeof detail === "string") {
    return detail;
  }

  return JSON.stringify(detail, null, 2);
}

export async function writeAuditLog(input: {
  adminUserId?: string | null;
  module: string;
  action: string;
  targetType: string;
  targetId?: string | null;
  summary: string;
  detail?: unknown;
}) {
  return db.auditLog.create({
    data: {
      adminUserId: input.adminUserId ?? undefined,
      module: input.module,
      action: input.action,
      targetType: input.targetType,
      targetId: input.targetId ?? undefined,
      summary: input.summary,
      detail: serializeDetail(input.detail),
    },
  });
}
