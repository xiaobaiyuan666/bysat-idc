import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManagePayments, getApiUser } from "@/lib/api-auth";
import { getPaymentsPageData } from "@/lib/data";
import { settleInvoicePayment } from "@/lib/payment-service";
import { paymentSchema } from "@/lib/validation";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManagePayments(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getPaymentsPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManagePayments(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = paymentSchema.safeParse({
    invoiceId: body.invoiceId,
    method: body.method,
    amount: Math.round(Number(body.amount ?? 0) * 100),
    transactionNo: body.transactionNo,
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "请求参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  try {
    const data = await settleInvoicePayment({
      invoiceId: parsed.data.invoiceId,
      method: parsed.data.method,
      amount: parsed.data.amount,
      transactionNo: parsed.data.transactionNo,
      operatorId: user.id,
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "payment",
      action: "create",
      targetType: "payment",
      targetId: data.payment.id,
      summary: `创建收款：${data.payment.paymentNo}`,
      detail: {
        invoiceNo: data.invoice.invoiceNo,
        method: data.payment.method,
        amount: data.payment.amount,
      },
    });

    return NextResponse.json({ data: data.payment }, { status: 201 });
  } catch (error) {
    if (error instanceof Error && error.message === "INVOICE_NOT_FOUND") {
      return NextResponse.json({ message: "账单不存在" }, { status: 404 });
    }

    if (error instanceof Error && error.message === "INSUFFICIENT_BALANCE") {
      return NextResponse.json(
        { message: "客户余额不足" },
        { status: 400 },
      );
    }

    if (
      error instanceof Error &&
      error.message === "PAYMENT_EXCEEDS_OUTSTANDING"
    ) {
      return NextResponse.json(
        { message: "收款金额不能超过账单未收金额" },
        { status: 400 },
      );
    }

    throw error;
  }
}
