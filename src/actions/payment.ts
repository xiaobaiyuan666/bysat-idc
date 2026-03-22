"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireUser } from "@/lib/auth";
import { db } from "@/lib/db";
import { getCloudProvider } from "@/lib/providers";
import {
  buildProviderPayload,
  syncServiceRecord,
  writeProviderLog,
} from "@/lib/provider-sync";
import { parseMoneyToCent } from "@/lib/format";
import { paymentSchema } from "@/lib/validation";
import { makeCode } from "@/lib/utils";

export async function recordPaymentAction(formData: FormData) {
  await requireUser();

  const parsed = paymentSchema.safeParse({
    invoiceId: formData.get("invoiceId"),
    method: formData.get("method"),
    amount: parseMoneyToCent(formData.get("amount")),
    transactionNo: formData.get("transactionNo") || undefined,
  });

  if (!parsed.success) {
    redirect("/payments");
  }

  const invoice = await db.invoice.findUnique({
    where: {
      id: parsed.data.invoiceId,
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

  if (!invoice) {
    redirect("/payments");
  }

  const payment = await db.payment.create({
    data: {
      paymentNo: makeCode("PAY"),
      customerId: invoice.customerId,
      invoiceId: invoice.id,
      orderId: invoice.orderId,
      method: parsed.data.method,
      status: "SUCCESS",
      amount: parsed.data.amount,
      transactionNo: parsed.data.transactionNo,
      paidAt: new Date(),
    },
  });

  const invoicePaidAmount = invoice.paidAmount + payment.amount;
  const invoiceFullyPaid = invoicePaidAmount >= invoice.totalAmount;

  await db.invoice.update({
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
    const orderPayments = await db.payment.aggregate({
      where: {
        orderId: invoice.orderId,
        status: "SUCCESS",
      },
      _sum: {
        amount: true,
      },
    });

    const orderPaidAmount = orderPayments._sum.amount ?? 0;
    await db.order.update({
      where: {
        id: invoice.orderId,
      },
      data: {
        paidAmount: orderPaidAmount,
        paidAt: invoiceFullyPaid ? new Date() : invoice.order?.paidAt,
        status: orderPaidAmount >= (invoice.order?.totalAmount ?? 0) ? "ACTIVE" : "PAID",
      },
    });
  }

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

  ["/dashboard", "/orders", "/services", "/invoices", "/payments"].forEach((path) =>
    revalidatePath(path),
  );

  redirect("/payments");
}
