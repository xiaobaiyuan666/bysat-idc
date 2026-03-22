import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import {
  canManageNotifications,
  canViewNotifications,
  getApiUser,
} from "@/lib/api-auth";
import { getNotificationsPageData } from "@/lib/data";
import { queueNotification } from "@/lib/notification-service";
import { manualNotificationSchema } from "@/lib/validation";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewNotifications(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getNotificationsPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageNotifications(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = manualNotificationSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "通知参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const result = await queueNotification({
    customerId: parsed.data.customerId || undefined,
    createdById: user.id,
    channel: parsed.data.channel,
    priority: parsed.data.priority,
    recipient: parsed.data.recipient,
    recipientName: parsed.data.recipientName || undefined,
    subject: parsed.data.subject || undefined,
    content: parsed.data.content,
    module: parsed.data.module || "notification",
    relatedType: parsed.data.relatedType || undefined,
    relatedId: parsed.data.relatedId || undefined,
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "notification",
    action: "create",
    targetType: "notification",
    targetId: result.message.id,
    summary: `创建通知：${result.message.recipient}`,
    detail: {
      channel: result.message.channel,
      subject: result.message.subject,
      recipient: result.message.recipient,
    },
  });

  return NextResponse.json({ data: result.message }, { status: 201 });
}
