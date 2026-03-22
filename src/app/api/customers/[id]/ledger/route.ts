import { NextRequest, NextResponse } from "next/server";

import { canManageCustomers, getApiUser } from "@/lib/api-auth";
import { getCustomerLedgerData } from "@/lib/data";

export async function GET(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await params;
  const data = await getCustomerLedgerData(id);

  if (!data) {
    return NextResponse.json({ message: "客户不存在" }, { status: 404 });
  }

  return NextResponse.json({ data });
}
