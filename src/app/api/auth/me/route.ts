import { NextRequest, NextResponse } from "next/server";

import { getApiUser, serializeAuthUser } from "@/lib/api-auth";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  return NextResponse.json({
    data: serializeAuthUser(user),
  });
}
