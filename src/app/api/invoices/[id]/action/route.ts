import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageInvoices, getApiUser } from "@/lib/api-auth";
import { issueInvoice, voidInvoice } from "@/lib/invoice-service";
import { invoiceActionSchema } from "@/lib/validation";

export async function POST(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageInvoices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = invoiceActionSchema.safeParse({
    action: body.action,
    reason: body.reason,
    dueDate: body.dueDate,
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "账单操作参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  try {
    const { id } = await params;

    if (parsed.data.action === "issue") {
      const invoice = await issueInvoice({
        invoiceId: id,
        actorId: user.id,
        dueDate: parsed.data.dueDate || undefined,
      });

      await writeAuditLog({
        adminUserId: user.id,
        module: "invoice",
        action: "issue",
        targetType: "invoice",
        targetId: invoice.id,
        summary: `签发账单：${invoice.invoiceNo}`,
        detail: {
          dueDate: invoice.dueDate,
          totalAmount: invoice.totalAmount,
        },
      });

      return NextResponse.json({ data: invoice });
    }

    const invoice = await voidInvoice({
      invoiceId: id,
      actorId: user.id,
      reason: parsed.data.reason || undefined,
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "invoice",
      action: "void",
      targetType: "invoice",
      targetId: invoice.id,
      summary: `作废账单：${invoice.invoiceNo}`,
      detail: {
        reason: parsed.data.reason || null,
      },
    });

    return NextResponse.json({ data: invoice });
  } catch (error) {
    if (error instanceof Error && error.message === "INVOICE_NOT_FOUND") {
      return NextResponse.json({ message: "账单不存在" }, { status: 404 });
    }

    if (error instanceof Error && error.message === "INVOICE_NOT_DRAFT") {
      return NextResponse.json(
        { message: "只有草稿账单才能执行签发" },
        { status: 400 },
      );
    }

    if (error instanceof Error && error.message === "INVOICE_ALREADY_VOID") {
      return NextResponse.json({ message: "账单已作废" }, { status: 400 });
    }

    if (error instanceof Error && error.message === "INVOICE_HAS_PAYMENT") {
      return NextResponse.json(
        { message: "账单已有收款记录，不允许直接作废" },
        { status: 400 },
      );
    }

    if (error instanceof Error && error.message === "INVOICE_DATE_INVALID") {
      return NextResponse.json({ message: "账单日期无效" }, { status: 400 });
    }

    throw error;
  }
}
