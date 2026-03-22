import { NextRequest, NextResponse } from "next/server";

import { canViewInvoices, getApiUser } from "@/lib/api-auth";
import { getInvoiceDetailData } from "@/lib/data";

export async function GET(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canViewInvoices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const data = await getInvoiceDetailData(id);

  if (!data) {
    return NextResponse.json({ message: "账单不存在" }, { status: 404 });
  }

  return NextResponse.json({ data });
}
