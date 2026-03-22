import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import {
  canManageServices,
  canViewServices,
  getApiUser,
} from "@/lib/api-auth";
import { getServiceDetailData } from "@/lib/data";
import { db } from "@/lib/db";
import { serviceUpdateSchema } from "@/lib/validation";

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

  if (!canViewServices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const data = await getServiceDetailData(id);

  if (!data) {
    return NextResponse.json({ message: "服务不存在" }, { status: 404 });
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

  if (!canManageServices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const body = await request.json();
  const parsed = serviceUpdateSchema.safeParse({
    ...body,
    monthlyCost: Math.round(Number(body.monthlyCost ?? 0) * 100),
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "服务参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const service = await db.serviceInstance.update({
    where: {
      id,
    },
    data: {
      name: parsed.data.name,
      hostname: normalizeText(parsed.data.hostname),
      region: normalizeText(parsed.data.region),
      planId: parsed.data.planId || null,
      vpcNetworkId: parsed.data.vpcNetworkId || null,
      status: parsed.data.status,
      monthlyCost: parsed.data.monthlyCost,
      nextDueDate: parsed.data.nextDueDate ? new Date(parsed.data.nextDueDate) : null,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "service",
    action: "update",
    targetType: "service",
    targetId: service.id,
    summary: `更新服务 ${service.serviceNo}`,
  });

  return NextResponse.json({ data: service });
}
