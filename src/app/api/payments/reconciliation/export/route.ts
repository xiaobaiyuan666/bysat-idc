import { NextRequest, NextResponse } from "next/server";

import { canManagePayments, getApiUser } from "@/lib/api-auth";
import { createCsv } from "@/lib/csv";
import { getPaymentReconciliationData } from "@/lib/data";
import { formatDateTime } from "@/lib/format";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManagePayments(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getPaymentReconciliationData();
  const today = new Date().toISOString().slice(0, 10);

  const rows =
    data.issues.length > 0
      ? data.issues.map((issue) => ({
          异常级别: issue.level,
          异常类型: issue.kind,
          异常标题: issue.title,
          相关单号: issue.relatedNo,
          客户名称: issue.customerName,
          涉及金额: (issue.amount / 100).toFixed(2),
          发生时间: formatDateTime(issue.createdAt),
          异常详情: issue.detail,
        }))
      : [
          {
            异常级别: "info",
            异常类型: "SYSTEM",
            异常标题: "当前无对账异常",
            相关单号: "",
            客户名称: "",
            涉及金额: "0.00",
            发生时间: formatDateTime(new Date()),
            异常详情: "导出时系统未检测到待处理的对账异常。",
          },
        ];

  const csv = createCsv(rows);

  return new NextResponse(csv, {
    status: 200,
    headers: {
      "Content-Type": "text/csv; charset=utf-8",
      "Content-Disposition": `attachment; filename="payment-reconciliation-${today}.csv"`,
    },
  });
}
