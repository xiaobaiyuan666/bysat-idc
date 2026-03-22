import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canReplyTickets, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { ticketReplySchema } from "@/lib/validation";

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canReplyTickets(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await params;
  const body = await request.json();
  const parsed = ticketReplySchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "请求参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const ticket = await db.ticket.findUnique({
    where: {
      id,
    },
  });

  if (!ticket) {
    return NextResponse.json({ message: "工单不存在" }, { status: 404 });
  }

  await db.ticketReply.create({
    data: {
      ticketId: ticket.id,
      authorType: "STAFF",
      authorName: user.name,
      content: parsed.data.content,
      isInternal: parsed.data.isInternal,
    },
  });

  const updatedTicket = await db.ticket.update({
    where: {
      id: ticket.id,
    },
    data: {
      status: parsed.data.status ?? ticket.status,
      assignedToId: parsed.data.assignedToId ?? ticket.assignedToId,
      lastReplyAt: new Date(),
    },
    include: {
      customer: true,
      service: true,
      assignedTo: true,
      replies: {
        orderBy: {
          createdAt: "asc",
        },
      },
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "ticket",
    action: "reply",
    targetType: "ticket",
    targetId: ticket.id,
    summary: `回复工单：${ticket.ticketNo}`,
    detail: {
      status: parsed.data.status,
      isInternal: parsed.data.isInternal,
    },
  });

  return NextResponse.json({ data: updatedTicket }, { status: 201 });
}
