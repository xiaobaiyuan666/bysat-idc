import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import {
  canManageInvoices,
  canViewInvoices,
  getApiUser,
} from "@/lib/api-auth";
import { getInvoicesPageData } from "@/lib/data";
import { createManualInvoice } from "@/lib/invoice-service";
import { invoiceCreateSchema } from "@/lib/validation";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewInvoices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getInvoicesPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageInvoices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = invoiceCreateSchema.safeParse({
    customerId: body.customerId,
    serviceId: body.serviceId,
    taxProfileId: body.taxProfileId,
    type: body.type,
    subtotal: Math.round(Number(body.subtotal ?? 0) * 100),
    taxRate: Number(body.taxRate ?? 0),
    dueDate: body.dueDate,
    remark: body.remark,
    issueNow: Boolean(body.issueNow),
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "账单参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  try {
    const invoice = await createManualInvoice({
      ...parsed.data,
      serviceId: parsed.data.serviceId || undefined,
      taxProfileId: parsed.data.taxProfileId || undefined,
      actorId: user.id,
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "invoice",
      action: "create",
      targetType: "invoice",
      targetId: invoice.id,
      summary: `创建账单：${invoice.invoiceNo}`,
      detail: {
        customerId: invoice.customerId,
        serviceId: invoice.serviceId,
        type: invoice.type,
        taxRate: invoice.taxRate,
        taxProfileCode: invoice.taxProfileCode,
        totalAmount: invoice.totalAmount,
        dueDate: invoice.dueDate,
      },
    });

    return NextResponse.json({ data: invoice }, { status: 201 });
  } catch (error) {
    if (error instanceof Error && error.message === "CUSTOMER_NOT_FOUND") {
      return NextResponse.json({ message: "客户不存在" }, { status: 404 });
    }

    if (error instanceof Error && error.message === "SERVICE_CUSTOMER_MISMATCH") {
      return NextResponse.json(
        { message: "所选服务不属于当前客户" },
        { status: 400 },
      );
    }

    if (error instanceof Error && error.message === "INVOICE_TYPE_NOT_ALLOWED") {
      return NextResponse.json(
        { message: "当前仅支持创建手工账单和余额账单" },
        { status: 400 },
      );
    }

    if (error instanceof Error && error.message === "TAX_PROFILE_NOT_FOUND") {
      return NextResponse.json({ message: "税率档案不存在" }, { status: 404 });
    }

    if (error instanceof Error && error.message === "INVOICE_DATE_INVALID") {
      return NextResponse.json({ message: "账单日期无效" }, { status: 400 });
    }

    throw error;
  }
}
