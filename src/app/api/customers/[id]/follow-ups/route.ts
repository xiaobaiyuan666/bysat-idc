import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, canViewCustomers, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { customerFollowUpSchema } from "@/lib/validation";

function parseDate(value?: string | null) {
  if (!value) {
    return undefined;
  }

  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? undefined : date;
}

export async function GET(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canViewCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const followUps = await db.customerFollowUp.findMany({
    where: {
      customerId: id,
    },
    orderBy: {
      createdAt: "desc",
    },
  });

  return NextResponse.json({ data: followUps });
}

export async function POST(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const body = await request.json();
  const parsed = customerFollowUpSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "跟进记录参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const customer = await db.customer.findUnique({
    where: {
      id,
    },
    select: {
      id: true,
      name: true,
    },
  });

  if (!customer) {
    return NextResponse.json({ message: "客户不存在" }, { status: 404 });
  }

  const followUp = await db.customerFollowUp.create({
    data: {
      customerId: id,
      type: parsed.data.type,
      title: parsed.data.title,
      content: parsed.data.content,
      nextFollowAt: parseDate(parsed.data.nextFollowAt),
      operatorId: user.id,
      operatorName: user.name,
      operatorRole: user.role,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "create-follow-up",
    targetType: "customer",
    targetId: id,
    summary: `新增客户 ${customer.name} 的跟进记录`,
    detail: {
      followUpId: followUp.id,
      type: followUp.type,
      title: followUp.title,
    },
  });

  return NextResponse.json({ data: followUp }, { status: 201 });
}
