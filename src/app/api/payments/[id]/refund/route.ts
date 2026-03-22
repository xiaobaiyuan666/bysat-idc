import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManagePayments, getApiUser } from "@/lib/api-auth";
import { refundPayment } from "@/lib/payment-service";
import { refundSchema } from "@/lib/validation";

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManagePayments(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = refundSchema.safeParse({
    amount: Math.round(Number(body.amount ?? 0) * 100),
    reason: body.reason,
    method: body.method,
    transactionNo: body.transactionNo,
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "退款参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  try {
    const { id } = await params;
    const refund = await refundPayment({
      paymentId: id,
      amount: parsed.data.amount,
      reason: parsed.data.reason,
      method: parsed.data.method,
      transactionNo: parsed.data.transactionNo,
      operatorId: user.id,
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "payment",
      action: "refund",
      targetType: "refund",
      targetId: refund.id,
      summary: `创建退款：${refund.refundNo}`,
      detail: {
        paymentNo: refund.payment.paymentNo,
        invoiceNo: refund.invoice?.invoiceNo,
        amount: refund.amount,
        reason: refund.reason,
      },
    });

    return NextResponse.json({ data: refund }, { status: 201 });
  } catch (error) {
    if (error instanceof Error && error.message === "PAYMENT_NOT_FOUND") {
      return NextResponse.json({ message: "收款记录不存在" }, { status: 404 });
    }

    if (error instanceof Error && error.message === "PAYMENT_NOT_REFUNDABLE") {
      return NextResponse.json({ message: "当前收款状态不允许退款" }, { status: 400 });
    }

    if (error instanceof Error && error.message === "REFUND_EXCEEDS_AVAILABLE") {
      return NextResponse.json({ message: "退款金额超过可退余额" }, { status: 400 });
    }

    if (error instanceof Error && error.message === "REFUND_AMOUNT_INVALID") {
      return NextResponse.json({ message: "退款金额无效" }, { status: 400 });
    }

    throw error;
  }
}
