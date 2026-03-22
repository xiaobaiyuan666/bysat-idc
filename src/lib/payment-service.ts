import {
  type Payment,
  type PaymentMethod,
  RefundStatus,
} from "@prisma/client";

import { refreshCustomerStatus } from "@/lib/billing-engine";
import { changeCustomerBalance } from "@/lib/credit";
import { db } from "@/lib/db";
import { queueNotification } from "@/lib/notification-service";
import { getCloudProvider } from "@/lib/providers";
import {
  buildProviderPayload,
  syncServiceRecord,
  writeProviderLog,
} from "@/lib/provider-sync";
import { makeCode } from "@/lib/utils";

async function findInvoiceForPayment(invoiceId: string) {
  return db.invoice.findUnique({
    where: {
      id: invoiceId,
    },
    include: {
      customer: true,
      order: true,
      service: {
        include: {
          customer: true,
          product: true,
        },
      },
    },
  });
}

type PaymentInvoice = NonNullable<Awaited<ReturnType<typeof findInvoiceForPayment>>>;

export type SettledInvoicePaymentResult = {
  payment: Payment;
  invoice: PaymentInvoice;
  invoiceFullyPaid: boolean;
};

export async function settleInvoicePayment(input: {
  invoiceId: string;
  method: PaymentMethod;
  amount: number;
  transactionNo?: string;
  channelPayload?: string;
  operatorId?: string | null;
}): Promise<SettledInvoicePaymentResult> {
  const invoice = await findInvoiceForPayment(input.invoiceId);

  if (!invoice) {
    throw new Error("INVOICE_NOT_FOUND");
  }

  const outstanding = Math.max(invoice.totalAmount - invoice.paidAmount, 0);

  if (input.amount > outstanding) {
    throw new Error("PAYMENT_EXCEEDS_OUTSTANDING");
  }

  const invoicePaidAmount = invoice.paidAmount + input.amount;
  const invoiceFullyPaid = invoicePaidAmount >= invoice.totalAmount;

  let payment: Payment | null = null;

  await db.$transaction(async (tx) => {
    if (input.method === "BALANCE") {
      await changeCustomerBalance({
        customerId: invoice.customerId,
        amount: -input.amount,
        type: "CONSUME",
        description: `账单收款 ${invoice.invoiceNo}`,
        operatorId: input.operatorId ?? undefined,
        allowNegative: false,
        tx,
      });
    }

    payment = await tx.payment.create({
      data: {
        paymentNo: makeCode("PAY"),
        customerId: invoice.customerId,
        invoiceId: invoice.id,
        orderId: invoice.orderId,
        method: input.method,
        status: "SUCCESS",
        amount: input.amount,
        transactionNo: input.transactionNo,
        channelPayload: input.channelPayload,
        paidAt: new Date(),
      },
    });

    await tx.invoice.update({
      where: {
        id: invoice.id,
      },
      data: {
        paidAmount: invoicePaidAmount,
        paidAt: invoiceFullyPaid ? new Date() : invoice.paidAt,
        status: invoiceFullyPaid ? "PAID" : "PARTIAL",
      },
    });

    if (invoice.orderId) {
      const orderPayments = await tx.payment.aggregate({
        where: {
          orderId: invoice.orderId,
          status: "SUCCESS",
        },
        _sum: {
          amount: true,
        },
      });

      const orderPaidAmount = orderPayments._sum.amount ?? 0;

      await tx.order.update({
        where: {
          id: invoice.orderId,
        },
        data: {
          paidAmount: orderPaidAmount,
          paidAt: invoiceFullyPaid ? new Date() : invoice.order?.paidAt,
          status:
            orderPaidAmount >= (invoice.order?.totalAmount ?? 0)
              ? "ACTIVE"
              : "PAID",
        },
      });
    }
  });

  if (invoiceFullyPaid && invoice.service) {
    if (invoice.service.providerType === "MOFANG_CLOUD") {
      const provider = getCloudProvider();
      const payload = buildProviderPayload(invoice.service);
      const providerResult =
        invoice.service.providerResourceId || invoice.service.status === "OVERDUE"
          ? await provider.activateService(payload)
          : await provider.provisionService(payload);

      await syncServiceRecord({
        service: invoice.service,
        providerResult,
        fallbackStatus: "ACTIVE",
      });

      await writeProviderLog({
        action:
          invoice.service.providerResourceId || invoice.service.status === "OVERDUE"
            ? "activate"
            : "provision",
        resourceId: providerResult.remoteId ?? invoice.service.providerResourceId,
        message: providerResult.message,
        requestBody: providerResult.requestBody,
        responseBody: providerResult.responseBody,
        ok: providerResult.ok,
      });
    } else {
      await db.serviceInstance.update({
        where: {
          id: invoice.service.id,
        },
        data: {
          status: "ACTIVE",
          activatedAt: new Date(),
          lastSyncAt: new Date(),
        },
      });
    }
  }

  await refreshCustomerStatus(invoice.customerId);

  await queueNotification({
    templateCode: "PAYMENT_SUCCESS",
    customerId: invoice.customerId,
    recipient: invoice.customer.email,
    recipientName: invoice.customer.name,
    module: "payment",
    relatedType: "invoice",
    relatedId: invoice.id,
    variables: {
      customer_name: invoice.customer.name,
      invoice_no: invoice.invoiceNo,
      amount: (input.amount / 100).toFixed(2),
      method: input.method,
    },
    content: `账单 ${invoice.invoiceNo} 已完成收款，到账金额 ${(input.amount / 100).toFixed(2)} 元。`,
  });

  if (!payment) {
    throw new Error("PAYMENT_NOT_CREATED");
  }

  return {
    payment,
    invoice,
    invoiceFullyPaid,
  };
}

async function findPaymentForRefund(paymentId: string) {
  return db.payment.findUnique({
    where: {
      id: paymentId,
    },
    include: {
      customer: true,
      invoice: true,
      order: true,
      refunds: {
        where: {
          status: RefundStatus.SUCCESS,
        },
      },
    },
  });
}

export async function refundPayment(input: {
  paymentId: string;
  amount: number;
  reason: string;
  method?: PaymentMethod;
  transactionNo?: string;
  operatorId?: string | null;
}) {
  const payment = await findPaymentForRefund(input.paymentId);

  if (!payment) {
    throw new Error("PAYMENT_NOT_FOUND");
  }

  if (payment.status === "FAILED") {
    throw new Error("PAYMENT_NOT_REFUNDABLE");
  }

  const refundedAmount = payment.refunds.reduce(
    (total, refund) => total + refund.amount,
    0,
  );
  const refundableAmount = Math.max(payment.amount - refundedAmount, 0);

  if (input.amount > refundableAmount) {
    throw new Error("REFUND_EXCEEDS_AVAILABLE");
  }

  if (input.amount <= 0) {
    throw new Error("REFUND_AMOUNT_INVALID");
  }

  const refundMethod = input.method ?? payment.method;
  let createdRefundId = "";

  await db.$transaction(async (tx) => {
    if (refundMethod === "BALANCE") {
      await changeCustomerBalance({
        customerId: payment.customerId,
        amount: input.amount,
        type: "REFUND",
        description: `收款退款 ${payment.paymentNo}`,
        operatorId: input.operatorId ?? undefined,
        allowNegative: true,
        tx,
      });
    }

    const refund = await tx.refundRecord.create({
      data: {
        refundNo: makeCode("REF"),
        paymentId: payment.id,
        invoiceId: payment.invoiceId,
        customerId: payment.customerId,
        processedById: input.operatorId ?? undefined,
        amount: input.amount,
        method: refundMethod,
        status: RefundStatus.SUCCESS,
        reason: input.reason,
        transactionNo: input.transactionNo,
        payload: JSON.stringify(
          {
            originalPaymentNo: payment.paymentNo,
            originalMethod: payment.method,
          },
          null,
          2,
        ),
        processedAt: new Date(),
      },
    });

    createdRefundId = refund.id;

    if (payment.invoiceId) {
      const nextPaidAmount = Math.max((payment.invoice?.paidAmount ?? 0) - input.amount, 0);
      const now = new Date();
      const nextStatus =
        nextPaidAmount <= 0
          ? payment.invoice?.dueDate && payment.invoice.dueDate < now
            ? "OVERDUE"
            : "ISSUED"
          : nextPaidAmount < (payment.invoice?.totalAmount ?? 0)
            ? "PARTIAL"
            : "PAID";

      await tx.invoice.update({
        where: {
          id: payment.invoiceId,
        },
        data: {
          paidAmount: nextPaidAmount,
          paidAt: nextStatus === "PAID" ? payment.invoice?.paidAt : null,
          status: nextStatus,
        },
      });
    }

    const totalRefundedAfter = refundedAmount + input.amount;
    await tx.payment.update({
      where: {
        id: payment.id,
      },
      data: {
        status: totalRefundedAfter >= payment.amount ? "REFUNDED" : payment.status,
      },
    });

    if (payment.orderId) {
      const orderInvoices = await tx.invoice.findMany({
        where: {
          orderId: payment.orderId,
        },
        select: {
          totalAmount: true,
          paidAmount: true,
        },
      });

      const orderPaidAmount = orderInvoices.reduce(
        (total, invoice) => total + invoice.paidAmount,
        0,
      );
      const orderTotalAmount = orderInvoices.reduce(
        (total, invoice) => total + invoice.totalAmount,
        0,
      );

      await tx.order.update({
        where: {
          id: payment.orderId,
        },
        data: {
          paidAmount: orderPaidAmount,
          paidAt: orderPaidAmount > 0 ? payment.order?.paidAt ?? new Date() : null,
          status:
            orderPaidAmount <= 0
              ? "PENDING"
              : orderPaidAmount >= orderTotalAmount
                ? "ACTIVE"
                : "PAID",
        },
      });
    }
  });

  await refreshCustomerStatus(payment.customerId);

  const refund = await db.refundRecord.findUniqueOrThrow({
    where: {
      id: createdRefundId,
    },
    include: {
      customer: true,
      payment: true,
      invoice: true,
      processedBy: true,
    },
  });

  await queueNotification({
    templateCode: "PAYMENT_REFUND",
    customerId: payment.customerId,
    recipient: payment.customer.email,
    recipientName: payment.customer.name,
    module: "payment",
    relatedType: "refund",
    relatedId: refund.id,
    variables: {
      customer_name: payment.customer.name,
      payment_no: payment.paymentNo,
      refund_no: refund.refundNo,
      amount: (refund.amount / 100).toFixed(2),
      reason: refund.reason,
    },
    content: `收款 ${payment.paymentNo} 已发起退款，退款金额 ${(refund.amount / 100).toFixed(2)} 元。`,
  });

  return refund;
}
