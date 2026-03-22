import { NextRequest, NextResponse } from "next/server";

import { getApiUser } from "@/lib/api-auth";
import { getDashboardData } from "@/lib/data";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录" }, { status: 401 });
  }

  const data = await getDashboardData();
  return NextResponse.json({ data });
}
