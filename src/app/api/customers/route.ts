import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, getApiUser } from "@/lib/api-auth";
import { getCustomersPageData } from "@/lib/data";
import { db } from "@/lib/db";
import { makeCode } from "@/lib/utils";
import { customerSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getCustomersPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = customerSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "客户参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const customer = await db.customer.create({
    data: {
      customerNo: makeCode("CUS"),
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
    action: "create",
    targetType: "customer",
    targetId: customer.id,
    summary: `创建客户 ${customer.name}`,
    detail: {
      customerNo: customer.customerNo,
      type: customer.type,
      level: customer.level,
    },
  });

  return NextResponse.json({ data: customer }, { status: 201 });
}
