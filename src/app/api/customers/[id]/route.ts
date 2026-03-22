import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import {
  canManageCustomers,
  canViewCustomers,
  getApiUser,
} from "@/lib/api-auth";
import { getCustomerDetailData } from "@/lib/data";
import { db } from "@/lib/db";
import { customerSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
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
  const data = await getCustomerDetailData(id);

  if (!data) {
    return NextResponse.json({ message: "客户不存在" }, { status: 404 });
  }

  return NextResponse.json({ data });
}

export async function PUT(
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
  const parsed = customerSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "客户参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const customer = await db.customer.update({
    where: {
      id,
    },
    data: {
      name: parsed.data.name,
      companyName: normalizeText(parsed.data.companyName),
      email: parsed.data.email,
      phone: normalizeText(parsed.data.phone),
      contactQQ: normalizeText(parsed.data.contactQQ),
      contactWechat: normalizeText(parsed.data.contactWechat),
      type: parsed.data.type,
      status: parsed.data.status,
      level: parsed.data.level,
      tags: normalizeText(parsed.data.tags),
      notes: normalizeText(parsed.data.notes),
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "update",
    targetType: "customer",
    targetId: customer.id,
    summary: `更新客户 ${customer.name}`,
  });

  return NextResponse.json({ data: customer });
}
