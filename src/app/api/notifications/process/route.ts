import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageNotifications, getApiUser } from "@/lib/api-auth";
import { processPendingNotifications } from "@/lib/notification-service";

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageNotifications(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json().catch(() => ({}));
  const limit = Math.min(Math.max(Number(body.limit ?? 20), 1), 100);
  const result = await processPendingNotifications(limit);

  await writeAuditLog({
    adminUserId: user.id,
    module: "notification",
    action: "process",
    targetType: "notificationQueue",
    summary: `执行通知队列：处理 ${result.processed} 条任务`,
    detail: result,
  });

  return NextResponse.json({ data: result });
}
