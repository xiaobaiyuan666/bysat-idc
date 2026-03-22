import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageTickets, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { getTicketsPageData } from "@/lib/data";
import { ticketSchema } from "@/lib/validation";
import { makeCode } from "@/lib/utils";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageTickets(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getTicketsPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageTickets(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = ticketSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "请求参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const ticket = await db.ticket.create({
    data: {
      ticketNo: makeCode("TIC"),
      customerId: parsed.data.customerId,
      serviceId: parsed.data.serviceId,
      subject: parsed.data.subject,
      priority: parsed.data.priority,
      status: "OPEN",
      summary: parsed.data.summary,
      assignedToId: user.id,
      lastReplyAt: new Date(),
    },
  });

  await db.ticketReply.create({
    data: {
      ticketId: ticket.id,
      authorType: "STAFF",
      authorName: user.name,
      content: parsed.data.summary,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "ticket",
    action: "create",
    targetType: "ticket",
    targetId: ticket.id,
    summary: `创建工单：${ticket.ticketNo}`,
    detail: {
      subject: ticket.subject,
      priority: ticket.priority,
    },
  });

  return NextResponse.json({ data: ticket }, { status: 201 });
}
