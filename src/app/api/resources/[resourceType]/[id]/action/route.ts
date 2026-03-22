import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageResources, getApiUser } from "@/lib/api-auth";
import {
  executeResourceAction,
  getResourceActionLabel,
  isSupportedResourceAction,
  isSupportedResourceType,
  ResourceActionError,
  type ResourceType,
} from "@/lib/resource-operations";

function unauthorized() {
  return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
}

function forbidden() {
  return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
}

function badRequest(message: string) {
  return NextResponse.json({ message }, { status: 400 });
}

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ resourceType: string; id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return unauthorized();
  }

  if (!canManageResources(user)) {
    return forbidden();
  }

  const { resourceType, id } = await params;
  const body = await request.json();
  const action = String(body.action ?? "");
  const payload = (body.payload ?? {}) as Record<string, unknown>;

  if (!isSupportedResourceType(resourceType)) {
    return badRequest("不支持的资源类型");
  }

  if (!isSupportedResourceAction(resourceType, action)) {
    return badRequest("不支持的资源操作");
  }

  try {
    const result = await executeResourceAction(
      resourceType,
      id,
      action,
      payload,
    );

    await writeAuditLog({
      adminUserId: user.id,
      module: "resource",
      action,
      targetType: result.targetType,
      targetId: result.targetId,
      summary: `${getResourceActionLabel(resourceType as ResourceType, action)}：${result.summary}`,
      detail: {
        resourceType,
        payload,
        serviceId: result.serviceId,
        serviceNo: result.serviceNo,
      },
    });

    return NextResponse.json({
      data: result.data,
      message: result.summary,
    });
  } catch (error) {
    if (error instanceof ResourceActionError) {
      return NextResponse.json(
        { message: error.message },
        { status: error.statusCode },
      );
    }

    const message = error instanceof Error ? error.message : "资源操作失败";
    return badRequest(message);
  }
}
