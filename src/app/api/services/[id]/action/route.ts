import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageServices, getApiUser } from "@/lib/api-auth";
import { type ProviderActionPayload } from "@/lib/providers/types";
import {
  executeServiceAction,
  isSupportedServiceAction,
  serviceActionLabelMap,
  type SupportedServiceAction,
} from "@/lib/service-operations";
import { buildVncConsoleUrl } from "@/lib/vnc";

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageServices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await params;
  const body = await request.json();
  const payload = (body.payload ?? {}) as ProviderActionPayload;
  const actionValue = String(body.action ?? "");

  if (!isSupportedServiceAction(actionValue)) {
    return NextResponse.json({ message: "不支持的服务操作" }, { status: 400 });
  }

  const action = actionValue as SupportedServiceAction;
  const result = await executeServiceAction(id, action, payload);

  if (!result) {
    return NextResponse.json({ message: "服务不存在" }, { status: 404 });
  }

  const vncConsoleUrl = await buildVncConsoleUrl(
    request.nextUrl.origin,
    result.providerResult.data,
  );

  if (vncConsoleUrl) {
    result.providerResult.consoleUrl = vncConsoleUrl;
  }

  await writeAuditLog({
    adminUserId: user.id,
    module: "service",
    action,
    targetType: "service",
    targetId: result.originalService.id,
    summary: `服务执行操作：${result.originalService.serviceNo} / ${serviceActionLabelMap[action]}`,
    detail: {
      provider: result.originalService.providerType,
      remoteId:
        result.providerResult.remoteId ?? result.originalService.providerResourceId,
      taskId: result.providerResult.taskId,
      payload,
    },
  });

  if (!result.providerResult.ok) {
    return NextResponse.json(
      {
        message: result.providerResult.message ?? "服务操作失败",
        providerResult: result.providerResult,
      },
      { status: 400 },
    );
  }

  return NextResponse.json({
    data: result.data,
    providerResult: result.providerResult,
  });
}
