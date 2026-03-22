import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { db } from "@/lib/db";
import {
  getPaymentGatewayConfig,
  verifyGatewaySignature,
} from "@/lib/finance-config";
import { settleInvoicePayment } from "@/lib/payment-service";
import { paymentCallbackSchema } from "@/lib/validation";

export async function POST(request: NextRequest) {
  const rawBody = await request.text();
  let body: unknown;

  try {
    body = JSON.parse(rawBody);
  } catch {
    return NextResponse.json({ message: "回调数据不是有效的 JSON" }, { status: 400 });
  }

  const parsed = paymentCallbackSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "支付回调参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const payment = parsed.data.paymentNo
    ? await db.payment.findUnique({
        where: {
          paymentNo: parsed.data.paymentNo,
        },
      })
    : parsed.data.transactionNo
      ? await db.payment.findFirst({
          where: {
            transactionNo: parsed.data.transactionNo,
          },
        })
      : null;

  const invoice = parsed.data.invoiceNo
    ? await db.invoice.findUnique({
        where: {
          invoiceNo: parsed.data.invoiceNo,
        },
      })
    : payment?.invoiceId
      ? await db.invoice.findUnique({
          where: {
            id: payment.invoiceId,
          },
        })
      : null;

  const callbackLog = await db.paymentCallbackLog.create({
    data: {
      paymentId: payment?.id,
      method: parsed.data.method,
      invoiceNo: parsed.data.invoiceNo ?? invoice?.invoiceNo,
      paymentNo: parsed.data.paymentNo ?? payment?.paymentNo,
      transactionNo: parsed.data.transactionNo,
      callbackStatus: parsed.data.status,
      signature: parsed.data.signature,
      payload: rawBody || undefined,
      message: parsed.data.message,
    },
  });

  const gateway = await getPaymentGatewayConfig(parsed.data.method);

  if (!gateway) {
    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        isHandled: true,
        handledAt: new Date(),
        message: "未配置对应的支付渠道",
      },
    });

    return NextResponse.json({ message: "支付渠道未配置" }, { status: 403 });
  }

  if (!gateway.isEnabled) {
    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        isHandled: true,
        handledAt: new Date(),
        message: "支付渠道已停用",
      },
    });

    return NextResponse.json({ message: "支付渠道已停用" }, { status: 403 });
  }

  const signatureResult = verifyGatewaySignature({
    config: gateway,
    request,
    rawBody,
    bodySignature: parsed.data.signature,
  });

  if (!signatureResult.ok) {
    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        isHandled: true,
        handledAt: new Date(),
        message:
          signatureResult.reason === "MISSING_SIGNATURE"
            ? "缺少签名信息"
            : "支付回调签名校验失败",
      },
    });

    return NextResponse.json({ message: "支付回调签名无效" }, { status: 401 });
  }

  if (!invoice) {
    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        isHandled: true,
        handledAt: new Date(),
        message: "未找到对应账单，回调已记录",
      },
    });

    return NextResponse.json({ message: "账单不存在" }, { status: 404 });
  }

  if (parsed.data.status !== "SUCCESS") {
    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        isHandled: true,
        handledAt: new Date(),
        message: parsed.data.message ?? `回调状态为 ${parsed.data.status}`,
      },
    });

    return NextResponse.json({
      data: {
        handled: true,
        callbackLogId: callbackLog.id,
        status: parsed.data.status,
      },
    });
  }

  if (parsed.data.transactionNo) {
    const existingPayment = await db.payment.findFirst({
      where: {
        invoiceId: invoice.id,
        transactionNo: parsed.data.transactionNo,
      },
    });

    if (existingPayment) {
      await db.paymentCallbackLog.update({
        where: {
          id: callbackLog.id,
        },
        data: {
          paymentId: existingPayment.id,
          isHandled: true,
          handledAt: new Date(),
          message: "重复回调已忽略",
        },
      });

      return NextResponse.json({
        data: {
          handled: true,
          duplicate: true,
          paymentId: existingPayment.id,
        },
      });
    }
  }

  try {
    const amount =
      parsed.data.amount !== undefined
        ? Math.round(parsed.data.amount * 100)
        : Math.max(invoice.totalAmount - invoice.paidAmount, 0);

    const result = await settleInvoicePayment({
      invoiceId: invoice.id,
      method: parsed.data.method,
      amount,
      transactionNo: parsed.data.transactionNo,
      channelPayload: rawBody || undefined,
    });

    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        paymentId: result.payment.id,
        isHandled: true,
        handledAt: new Date(),
        message: "支付回调已成功入账",
      },
    });

    await writeAuditLog({
      module: "payment",
      action: "callback",
      targetType: "payment",
      targetId: result.payment.id,
      summary: `处理支付回调：${result.invoice.invoiceNo}`,
      detail: {
        method: parsed.data.method,
        transactionNo: parsed.data.transactionNo,
        callbackLogId: callbackLog.id,
        gatewayMethod: gateway.method,
        signType: gateway.signType,
      },
    });

    return NextResponse.json({
      data: {
        handled: true,
        paymentId: result.payment.id,
        callbackLogId: callbackLog.id,
      },
    });
  } catch (error) {
    const callbackErrorMessage =
      error instanceof Error && error.message === "PAYMENT_EXCEEDS_OUTSTANDING"
        ? "回调金额不能超过账单未收金额"
        : error instanceof Error && error.message === "INSUFFICIENT_BALANCE"
          ? "余额结算失败"
          : error instanceof Error
            ? error.message
            : "支付回调处理失败";

    await db.paymentCallbackLog.update({
      where: {
        id: callbackLog.id,
      },
      data: {
        isHandled: true,
        handledAt: new Date(),
        message: callbackErrorMessage,
      },
    });

    if (
      error instanceof Error &&
      error.message === "PAYMENT_EXCEEDS_OUTSTANDING"
    ) {
      return NextResponse.json(
        { message: "回调金额不能超过账单未收金额" },
        { status: 400 },
      );
    }

    if (error instanceof Error && error.message === "INSUFFICIENT_BALANCE") {
      return NextResponse.json(
        { message: "余额结算失败" },
        { status: 400 },
      );
    }

    throw error;
  }
}
