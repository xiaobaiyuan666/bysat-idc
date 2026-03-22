import { NextRequest, NextResponse } from "next/server";

import { canViewDashboard, getApiUser } from "@/lib/api-auth";
import { getReportsOverviewData } from "@/lib/data";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewDashboard(user)) {
    return NextResponse.json({ message: "没有权限查看统计报表" }, { status: 403 });
  }

  const data = await getReportsOverviewData();
  return NextResponse.json({ data });
}
