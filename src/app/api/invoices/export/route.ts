import { NextRequest, NextResponse } from "next/server";

import { canViewInvoices, getApiUser } from "@/lib/api-auth";
import { createCsv } from "@/lib/csv";
import { db } from "@/lib/db";
import { formatDate } from "@/lib/format";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canViewInvoices(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const invoices = await db.invoice.findMany({
    include: {
      customer: true,
      order: true,
      service: true,
    },
    orderBy: {
      createdAt: "desc",
    },
  });

  const csv = createCsv(
    invoices.map((invoice) => ({
      账单编号: invoice.invoiceNo,
      客户名称: invoice.customer.name,
      客户编号: invoice.customer.customerNo,
      账单类型: invoice.type,
      账单状态: invoice.status,
      税率档案: invoice.taxProfileName ?? "",
      税率: `${invoice.taxRate}%`,
      小计金额: (invoice.subtotal / 100).toFixed(2),
      税额: (invoice.taxAmount / 100).toFixed(2),
      账单总额: (invoice.totalAmount / 100).toFixed(2),
      已收金额: (invoice.paidAmount / 100).toFixed(2),
      待收金额: ((invoice.totalAmount - invoice.paidAmount) / 100).toFixed(2),
      到期日: formatDate(invoice.dueDate),
      签发日期: formatDate(invoice.issuedAt),
      订单编号: invoice.order?.orderNo ?? "",
      服务编号: invoice.service?.serviceNo ?? "",
      服务名称: invoice.service?.name ?? "",
      备注: invoice.remark ?? "",
      创建时间: invoice.createdAt.toISOString(),
    })),
  );

  return new NextResponse(csv, {
    status: 200,
    headers: {
      "Content-Type": "text/csv; charset=utf-8",
      "Content-Disposition": `attachment; filename="invoices-${new Date().toISOString().slice(0, 10)}.csv"`,
    },
  });
}
