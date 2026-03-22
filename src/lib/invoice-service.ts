import { InvoiceStatus, InvoiceType } from "@prisma/client";

import { getBillingSetting } from "@/lib/billing-engine";
import { db } from "@/lib/db";
import { queueNotification } from "@/lib/notification-service";
import { makeCode } from "@/lib/utils";

const OPEN_INVOICE_STATUSES: InvoiceStatus[] = ["ISSUED", "PARTIAL", "OVERDUE"];

async function getInvoiceDetail(invoiceId: string) {
  return db.invoice.findUnique({
    where: {
      id: invoiceId,
    },
    include: {
      customer: true,
      order: true,
      service: true,
      payments: true,
      refunds: true,
    },
  });
}

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

function toInvoiceDate(value: string) {
  const parsed = new Date(value);

  if (Number.isNaN(parsed.getTime())) {
    throw new Error("INVOICE_DATE_INVALID");
  }

  return parsed;
}

function isAllowedManualType(type: InvoiceType) {
  return ["MANUAL", "CREDIT"].includes(type);
}

export async function createManualInvoice(input: {
  customerId: string;
  serviceId?: string;
  taxProfileId?: string;
  type: InvoiceType;
  subtotal: number;
  taxRate: number;
  dueDate: string;
  remark?: string;
  issueNow: boolean;
  actorId?: string | null;
}) {
  if (!isAllowedManualType(input.type)) {
    throw new Error("INVOICE_TYPE_NOT_ALLOWED");
  }

  const dueDate = toInvoiceDate(input.dueDate);
  const [customer, service, settings, taxProfile] = await Promise.all([
    db.customer.findUnique({
      where: {
        id: input.customerId,
      },
    }),
    input.serviceId
      ? db.serviceInstance.findUnique({
          where: {
            id: input.serviceId,
          },
          include: {
            customer: true,
          },
        })
      : Promise.resolve(null),
    getBillingSetting(),
    input.taxProfileId
      ? db.invoiceTaxProfile.findFirst({
          where: {
            id: input.taxProfileId,
            isActive: true,
          },
        })
      : Promise.resolve(null),
  ]);

  if (!customer) {
    throw new Error("CUSTOMER_NOT_FOUND");
  }

  if (service && service.customerId !== customer.id) {
    throw new Error("SERVICE_CUSTOMER_MISMATCH");
  }

  if (input.taxProfileId && !taxProfile) {
    throw new Error("TAX_PROFILE_NOT_FOUND");
  }

  const appliedTaxRate = taxProfile?.taxRate ?? input.taxRate ?? settings.defaultTaxRate;
  const taxAmount = Math.round((input.subtotal * appliedTaxRate) / 100);
  const totalAmount = input.subtotal + taxAmount;
  const status: InvoiceStatus = input.issueNow ? "ISSUED" : "DRAFT";

  const invoice = await db.invoice.create({
    data: {
      invoiceNo: makeCode(settings.invoicePrefix || "INV"),
      customerId: customer.id,
      serviceId: service?.id,
      type: input.type,
      status,
      subtotal: input.subtotal,
      taxRate: appliedTaxRate,
      taxAmount,
      taxProfileCode: taxProfile?.code,
      taxProfileName: taxProfile?.name,
      totalAmount,
      dueDate,
      issuedAt: input.issueNow ? new Date() : null,
      remark: normalizeText(input.remark),
    },
    include: {
      customer: true,
      order: true,
      service: true,
      payments: true,
      refunds: true,
    },
  });

  if (status === "ISSUED") {
    await queueNotification({
      customerId: customer.id,
      createdById: input.actorId,
      recipient: customer.email,
      recipientName: customer.name,
      channel: "SYSTEM",
      module: "invoice",
      relatedType: "invoice",
      relatedId: invoice.id,
      subject: `新账单 ${invoice.invoiceNo} 已签发`,
      content: `账单 ${invoice.invoiceNo} 已签发，应付金额 ${(invoice.totalAmount / 100).toFixed(2)} 元，请在 ${invoice.dueDate.toISOString().slice(0, 10)} 前完成支付。`,
    });
  }

  return invoice;
}

export async function issueInvoice(input: {
  invoiceId: string;
  actorId?: string | null;
  dueDate?: string;
}) {
  const invoice = await getInvoiceDetail(input.invoiceId);

  if (!invoice) {
    throw new Error("INVOICE_NOT_FOUND");
  }

  if (invoice.status !== "DRAFT") {
    throw new Error("INVOICE_NOT_DRAFT");
  }

  const nextDueDate = input.dueDate ? toInvoiceDate(input.dueDate) : invoice.dueDate;

  const updated = await db.invoice.update({
    where: {
      id: invoice.id,
    },
    data: {
      status: "ISSUED",
      issuedAt: new Date(),
      dueDate: nextDueDate,
    },
    include: {
      customer: true,
      order: true,
      service: true,
      payments: true,
      refunds: true,
    },
  });

  await queueNotification({
    customerId: updated.customerId,
    createdById: input.actorId,
    recipient: updated.customer.email,
    recipientName: updated.customer.name,
    channel: "SYSTEM",
    module: "invoice",
    relatedType: "invoice",
    relatedId: updated.id,
    subject: `账单 ${updated.invoiceNo} 已签发`,
    content: `账单 ${updated.invoiceNo} 已签发，应付金额 ${(updated.totalAmount / 100).toFixed(2)} 元，请及时支付。`,
  });

  return updated;
}

export async function voidInvoice(input: {
  invoiceId: string;
  actorId?: string | null;
  reason?: string;
}) {
  const invoice = await getInvoiceDetail(input.invoiceId);

  if (!invoice) {
    throw new Error("INVOICE_NOT_FOUND");
  }

  if (invoice.status === "VOID") {
    throw new Error("INVOICE_ALREADY_VOID");
  }

  if (invoice.paidAmount > 0 || invoice.payments.length > 0) {
    throw new Error("INVOICE_HAS_PAYMENT");
  }

  const reason = normalizeText(input.reason);

  const updated = await db.$transaction(async (tx) => {
    const nextRemark = reason
      ? `${invoice.remark ? `${invoice.remark}\n` : ""}作废原因：${reason}`
      : invoice.remark;

    const nextInvoice = await tx.invoice.update({
      where: {
        id: invoice.id,
      },
      data: {
        status: "VOID",
        remark: nextRemark,
      },
      include: {
        customer: true,
        order: true,
        service: true,
        payments: true,
        refunds: true,
      },
    });

    if (invoice.orderId) {
      const openInvoiceCount = await tx.invoice.count({
        where: {
          orderId: invoice.orderId,
          status: {
            in: OPEN_INVOICE_STATUSES,
          },
        },
      });

      if (openInvoiceCount === 0 && (invoice.order?.paidAmount ?? 0) <= 0) {
        await tx.order.update({
          where: {
            id: invoice.orderId,
          },
          data: {
            status: "CANCELLED",
            paidAt: null,
          },
        });
      }
    }

    return nextInvoice;
  });

  await queueNotification({
    customerId: updated.customerId,
    createdById: input.actorId,
    recipient: updated.customer.email,
    recipientName: updated.customer.name,
    channel: "SYSTEM",
    module: "invoice",
    relatedType: "invoice",
    relatedId: updated.id,
    subject: `账单 ${updated.invoiceNo} 已作废`,
    content: reason
      ? `账单 ${updated.invoiceNo} 已作废，原因：${reason}。`
      : `账单 ${updated.invoiceNo} 已作废。`,
  });

  return updated;
}
