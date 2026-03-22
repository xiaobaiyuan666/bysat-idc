import { NextRequest, NextResponse } from "next/server";

import { canManagePayments, getApiUser } from "@/lib/api-auth";
import { createCsv } from "@/lib/csv";
import { db } from "@/lib/db";
import { formatDateTime } from "@/lib/format";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManagePayments(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const kind = request.nextUrl.searchParams.get("kind") || "payments";
  const today = new Date().toISOString().slice(0, 10);

  if (kind === "refunds") {
    const refunds = await db.refundRecord.findMany({
      include: {
        customer: true,
        payment: true,
        invoice: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    });

    const csv = createCsv(
      refunds.map((refund) => ({
        退款单号: refund.refundNo,
        客户名称: refund.customer.name,
        原收款单号: refund.payment.paymentNo,
        关联账单号: refund.invoice?.invoiceNo ?? "",
        退款方式: refund.method,
        退款状态: refund.status,
        退款金额: (refund.amount / 100).toFixed(2),
        退款原因: refund.reason,
        渠道流水号: refund.transactionNo ?? "",
        处理时间: formatDateTime(refund.processedAt || refund.createdAt),
      })),
    );

    return new NextResponse(csv, {
      status: 200,
      headers: {
        "Content-Type": "text/csv; charset=utf-8",
        "Content-Disposition": `attachment; filename="refunds-${today}.csv"`,
      },
    });
  }

  if (kind === "callbacks") {
    const callbacks = await db.paymentCallbackLog.findMany({
      include: {
        payment: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    });

    const csv = createCsv(
      callbacks.map((log) => ({
        回调时间: formatDateTime(log.createdAt),
        支付方式: log.method ?? "",
        账单号: log.invoiceNo ?? "",
        收款单号: log.paymentNo ?? log.payment?.paymentNo ?? "",
        渠道流水号: log.transactionNo ?? "",
        回调状态: log.callbackStatus ?? "",
        是否处理: log.isHandled ? "是" : "否",
        处理时间: formatDateTime(log.handledAt),
        签名值: log.signature ?? "",
        处理说明: log.message ?? "",
      })),
    );

    return new NextResponse(csv, {
      status: 200,
      headers: {
        "Content-Type": "text/csv; charset=utf-8",
        "Content-Disposition": `attachment; filename="payment-callbacks-${today}.csv"`,
      },
    });
  }

  const payments = await db.payment.findMany({
    include: {
      customer: true,
      invoice: true,
      refunds: true,
    },
    orderBy: {
      createdAt: "desc",
    },
  });

  const csv = createCsv(
    payments.map((payment) => ({
      收款单号: payment.paymentNo,
      客户名称: payment.customer.name,
      账单号: payment.invoice?.invoiceNo ?? "",
      支付方式: payment.method,
      收款状态: payment.status,
      收款金额: (payment.amount / 100).toFixed(2),
      已退款金额: (
        payment.refunds
          .filter((refund) => refund.status === "SUCCESS")
          .reduce((total, refund) => total + refund.amount, 0) / 100
      ).toFixed(2),
      渠道流水号: payment.transactionNo ?? "",
      入账时间: formatDateTime(payment.paidAt || payment.createdAt),
    })),
  );

  return new NextResponse(csv, {
    status: 200,
    headers: {
      "Content-Type": "text/csv; charset=utf-8",
      "Content-Disposition": `attachment; filename="payments-${today}.csv"`,
    },
  });
}
