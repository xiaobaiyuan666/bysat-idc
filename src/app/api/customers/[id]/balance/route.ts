import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, getApiUser } from "@/lib/api-auth";
import { changeCustomerBalance } from "@/lib/credit";
import { db } from "@/lib/db";
import { balanceAdjustmentSchema } from "@/lib/validation";

export async function POST(
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
  const body = await request.json();
  const parsed = balanceAdjustmentSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "请求参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const amount = Math.round(parsed.data.amount * 100);

  if (amount === 0) {
    return NextResponse.json({ message: "调整金额不能为 0" }, { status: 400 });
  }

  try {
    const data = await changeCustomerBalance({
      customerId: id,
      amount,
      type: amount > 0 ? "RECHARGE" : "ADJUSTMENT",
      description: parsed.data.description,
      operatorId: user.id,
      allowNegative: false,
    });

    const customer = await db.customer.findUnique({
      where: {
        id,
      },
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "customer",
      action: "adjust-balance",
      targetType: "customer",
      targetId: id,
      summary: `调整客户余额：${customer?.name ?? id}`,
      detail: {
        amount,
        description: parsed.data.description,
      },
    });

    return NextResponse.json({ data }, { status: 201 });
  } catch (error) {
    if (error instanceof Error && error.message === "INSUFFICIENT_BALANCE") {
      return NextResponse.json(
        { message: "客户余额不足，无法扣减" },
        { status: 400 },
      );
    }

    if (error instanceof Error && error.message === "CUSTOMER_NOT_FOUND") {
      return NextResponse.json({ message: "客户不存在" }, { status: 404 });
    }

    throw error;
  }
}
