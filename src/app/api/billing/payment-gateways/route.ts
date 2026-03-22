import { Prisma } from "@prisma/client";
import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageBilling, canViewBilling, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { normalizeGatewayInput } from "@/lib/finance-config";
import { paymentGatewaySchema } from "@/lib/validation";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewBilling(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await db.paymentGatewayConfig.findMany({
    orderBy: [
      {
        isEnabled: "desc",
      },
      {
        method: "asc",
      },
    ],
  });

  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageBilling(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = paymentGatewaySchema.safeParse({
    method: body.method,
    name: body.name,
    merchantId: body.merchantId,
    appId: body.appId,
    apiBaseUrl: body.apiBaseUrl,
    signType: body.signType,
    callbackSecret: body.callbackSecret,
    callbackHeader: body.callbackHeader,
    notifyUrl: body.notifyUrl,
    returnUrl: body.returnUrl,
    isEnabled: body.isEnabled !== false,
    remark: body.remark,
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "支付渠道参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  try {
    const data = await db.paymentGatewayConfig.create({
      data: normalizeGatewayInput(parsed.data),
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "billing",
      action: "create",
      targetType: "paymentGateway",
      targetId: data.id,
      summary: `创建支付渠道：${data.name}`,
      detail: {
        method: data.method,
        signType: data.signType,
        callbackHeader: data.callbackHeader,
        isEnabled: data.isEnabled,
      },
    });

    return NextResponse.json({ data }, { status: 201 });
  } catch (error) {
    if (
      error instanceof Prisma.PrismaClientKnownRequestError &&
      error.code === "P2002"
    ) {
      return NextResponse.json({ message: "该支付方式已存在配置" }, { status: 409 });
    }

    throw error;
  }
}
