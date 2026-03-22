import { NextRequest, NextResponse } from "next/server";

import { canViewOrders, getApiUser } from "@/lib/api-auth";
import { getOrderDetailData } from "@/lib/data";

export async function GET(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canViewOrders(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const data = await getOrderDetailData(id);

  if (!data) {
    return NextResponse.json({ message: "订单不存在" }, { status: 404 });
  }

  return NextResponse.json({ data });
}
