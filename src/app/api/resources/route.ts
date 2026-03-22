import { NextRequest, NextResponse } from "next/server";

import { canViewResources, getApiUser } from "@/lib/api-auth";
import { getResourcesPageData } from "@/lib/data";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewResources(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getResourcesPageData();
  return NextResponse.json({ data });
}
