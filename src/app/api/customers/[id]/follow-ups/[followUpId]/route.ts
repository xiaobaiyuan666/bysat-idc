import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { customerFollowUpSchema } from "@/lib/validation";

function parseDate(value?: string | null) {
  if (!value) {
    return undefined;
  }

  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? undefined : date;
}

async function getScopedFollowUp(customerId: string, followUpId: string) {
  return db.customerFollowUp.findFirst({
    where: {
      id: followUpId,
      customerId,
    },
  });
}

export async function PUT(
  request: NextRequest,
  context: { params: Promise<{ id: string; followUpId: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id, followUpId } = await context.params;
  const body = await request.json();
  const parsed = customerFollowUpSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "跟进记录参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const existing = await getScopedFollowUp(id, followUpId);

  if (!existing) {
    return NextResponse.json({ message: "跟进记录不存在" }, { status: 404 });
  }

  const followUp = await db.customerFollowUp.update({
    where: {
      id: followUpId,
    },
    data: {
      type: parsed.data.type,
      title: parsed.data.title,
      content: parsed.data.content,
      nextFollowAt: parseDate(parsed.data.nextFollowAt) ?? null,
      operatorId: user.id,
      operatorName: user.name,
      operatorRole: user.role,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "update-follow-up",
    targetType: "customer",
    targetId: id,
    summary: `更新客户跟进记录 ${followUp.title}`,
    detail: {
      followUpId: followUp.id,
    },
  });

  return NextResponse.json({ data: followUp });
}

export async function DELETE(
  request: NextRequest,
  context: { params: Promise<{ id: string; followUpId: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id, followUpId } = await context.params;
  const existing = await getScopedFollowUp(id, followUpId);

  if (!existing) {
    return NextResponse.json({ message: "跟进记录不存在" }, { status: 404 });
  }

  await db.customerFollowUp.delete({
    where: {
      id: followUpId,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "delete-follow-up",
    targetType: "customer",
    targetId: id,
    summary: `删除客户跟进记录 ${existing.title}`,
    detail: {
      followUpId: existing.id,
    },
  });

  return NextResponse.json({ success: true });
}
