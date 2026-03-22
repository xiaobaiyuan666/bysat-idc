import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageServices, getApiUser } from "@/lib/api-auth";
import { pullMofangCloudToLocal } from "@/lib/mofang-sync";

function normalizeLimit(value: unknown) {
  const limit = Number(value);

  if (!Number.isFinite(limit) || limit <= 0) {
    return 100;
  }

  return Math.min(Math.round(limit), 200);
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageServices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json().catch(() => ({}));
  const limit = normalizeLimit(body.limit);
  const includeResources =
    typeof body.includeResources === "boolean" ? body.includeResources : true;
  const result = await pullMofangCloudToLocal({
    limit,
    includeResources,
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "service",
    action: "provider-sync",
    targetType: "service",
    summary: `执行魔方云拉取同步，处理 ${result.summary.processedCount} 台实例`,
    detail: {
      limit,
      includeResources,
      summary: result.summary,
      items: result.items.slice(0, 20),
    },
  });

  if (!result.summary.listOk) {
    return NextResponse.json(
      {
        message: result.summary.listMessage,
        data: result,
      },
      { status: 502 },
    );
  }

  return NextResponse.json({ data: result });
}
