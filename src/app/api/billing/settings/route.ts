import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageBilling, canViewBilling, getApiUser } from "@/lib/api-auth";
import { getBillingSetting } from "@/lib/billing-engine";
import { db } from "@/lib/db";
import { billingSettingsSchema } from "@/lib/validation";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewBilling(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getBillingSetting();
  return NextResponse.json({ data });
}

export async function PUT(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageBilling(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = billingSettingsSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "财务配置参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  if (parsed.data.autoTerminateDays < parsed.data.autoSuspendDays) {
    return NextResponse.json(
      {
        message: "自动终止天数不能小于自动暂停天数",
      },
      { status: 400 },
    );
  }

  const data = await db.billingSetting.upsert({
    where: {
      id: "default",
    },
    update: parsed.data,
    create: {
      id: "default",
      ...parsed.data,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "billing",
    action: "update-settings",
    targetType: "billingSetting",
    targetId: data.id,
    summary: "更新计费与税务设置",
    detail: parsed.data,
  });

  return NextResponse.json({ data });
}
