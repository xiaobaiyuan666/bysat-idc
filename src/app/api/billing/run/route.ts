import { NextRequest, NextResponse } from "next/server";

import { canManageBilling, getApiUser } from "@/lib/api-auth";
import { runBillingEngine } from "@/lib/billing-engine";

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageBilling(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await runBillingEngine(user.id);
  return NextResponse.json({ data });
}
