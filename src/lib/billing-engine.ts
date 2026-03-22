import { addDays, differenceInCalendarDays, startOfDay } from "date-fns";
import {
  type BillingSetting,
  type CustomerStatus,
  InvoiceStatus,
  type Prisma,
  type ServiceStatus,
} from "@prisma/client";

import { writeAuditLog } from "@/lib/audit";
import { changeCustomerBalance } from "@/lib/credit";
import { db } from "@/lib/db";
import { getCloudProvider } from "@/lib/providers";
import {
  buildProviderPayload,
  getNextDueDate,
  syncServiceRecord,
  writeProviderLog,
} from "@/lib/provider-sync";
import { makeCode } from "@/lib/utils";

const OPEN_INVOICE_STATUSES: InvoiceStatus[] = [
  "ISSUED",
  "PARTIAL",
  "OVERDUE",
];

type BillingSummary = {
  createdInvoices: number;
  autoRenewed: number;
  markedOverdue: number;
  suspended: number;
  terminated: number;
  jobs: number;
};

const serviceStatusLabelMap: Record<ServiceStatus, string> = {
  PENDING: "待开通",
  PROVISIONING: "开通中",
  ACTIVE: "运行中",
  SUSPENDED: "已暂停",
  OVERDUE: "已逾期",
  TERMINATED: "已终止",
  EXPIRED: "已到期",
  FAILED: "异常",
};

export async function getBillingSetting() {
  return db.billingSetting.upsert({
    where: {
      id: "default",
    },
    update: {},
    create: {
      id: "default",
      invoicePrefix: "INV",
      invoiceIssuerName: "IDC云业务管理系统",
      defaultTaxRate: 13,
    },
  });
}

export async function refreshCustomerStatus(
  customerId: string,
): Promise<CustomerStatus> {
  const [customer, overdueCount] = await Promise.all([
    db.customer.findUnique({
      where: {
        id: customerId,
      },
      select: {
        status: true,
      },
    }),
    db.invoice.count({
      where: {
        customerId,
        status: "OVERDUE",
      },
    }),
  ]);

  if (!customer) {
    return "ACTIVE";
  }

  const nextStatus =
    overdueCount > 0
      ? "OVERDUE"
      : customer.status === "ARCHIVED"
        ? "ARCHIVED"
        : "ACTIVE";

  if (customer.status !== nextStatus) {
    await db.customer.update({
      where: {
        id: customerId,
      },
      data: {
        status: nextStatus,
      },
    });
  }

  return nextStatus;
}

async function createBillingJob(input: {
  jobType: string;
  status: "SUCCESS" | "FAILED" | "SKIPPED";
  customerId?: string | null;
  serviceId?: string | null;
  invoiceId?: string | null;
  message: string;
  payload?: unknown;
}) {
  return db.billingJob.create({
    data: {
      jobType: input.jobType,
      status: input.status,
      customerId: input.customerId ?? undefined,
      serviceId: input.serviceId ?? undefined,
      invoiceId: input.invoiceId ?? undefined,
      message: input.message,
      payload:
        input.payload === undefined
          ? undefined
          : JSON.stringify(input.payload, null, 2),
    },
  });
}

type RenewalServiceModel = Prisma.ServiceInstanceGetPayload<{
  include: {
    customer: true;
    product: true;
  };
}>;

async function createRenewalInvoiceForService(
  service: NonNullable<RenewalServiceModel>,
  input: {
    now: Date;
    actorId?: string | null;
    source: string;
    orderNote: string;
    invoiceRemark: string;
    jobMessagePrefix: string;
    invoicePrefix: string;
  },
) {
  const existingInvoice = await db.invoice.findFirst({
    where: {
      serviceId: service.id,
      type: "RENEWAL",
      status: {
        in: OPEN_INVOICE_STATUSES,
      },
    },
    include: {
      order: true,
    },
  });

  if (existingInvoice) {
    return {
      invoice: existingInvoice,
      order: existingInvoice.order,
      created: false,
    };
  }

  const amount = service.monthlyCost > 0 ? service.monthlyCost : service.product.price;

  if (amount <= 0) {
    return null;
  }

  const created = await db.$transaction(async (tx) => {
    const order = await tx.order.create({
      data: {
        orderNo: makeCode("ORD"),
        customerId: service.customerId,
        status: "PENDING",
        totalAmount: amount,
        currency: "CNY",
        source: input.source,
        orderType: "renew",
        dueDate: service.nextDueDate ?? input.now,
        notes: input.orderNote,
      },
    });

    await tx.orderItem.create({
      data: {
        orderId: order.id,
        productId: service.productId,
        serviceId: service.id,
        title: `${service.product.name} 续费`,
        quantity: 1,
        unitPrice: amount,
        cycle: service.billingCycle,
        totalAmount: amount,
      },
    });

    const invoice = await tx.invoice.create({
      data: {
        invoiceNo: makeCode(input.invoicePrefix || "INV"),
        customerId: service.customerId,
        orderId: order.id,
        serviceId: service.id,
        type: "RENEWAL",
        status: "ISSUED",
        subtotal: amount,
        totalAmount: amount,
        paidAmount: 0,
        dueDate: service.nextDueDate ?? input.now,
        issuedAt: input.now,
        remark: input.invoiceRemark,
      },
    });

    return { order, invoice };
  });

  await createBillingJob({
    jobType: "CREATE_RENEWAL_INVOICE",
    status: "SUCCESS",
    customerId: service.customerId,
    serviceId: service.id,
    invoiceId: created.invoice.id,
    message: `${input.jobMessagePrefix} ${created.invoice.invoiceNo}`,
    payload: {
      serviceNo: service.serviceNo,
      amount,
      dueDate: service.nextDueDate,
      source: input.source,
    },
  });

  await writeAuditLog({
    adminUserId: input.actorId,
    module: "billing",
    action: "create-renewal",
    targetType: "invoice",
    targetId: created.invoice.id,
    summary: `生成续费账单：${service.serviceNo}`,
    detail: {
      orderNo: created.order.orderNo,
      invoiceNo: created.invoice.invoiceNo,
      source: input.source,
    },
  });

  return {
    ...created,
    created: true,
  };
}

export async function createServiceRenewalInvoice(
  serviceId: string,
  input?: {
    actorId?: string | null;
    source?: string;
  },
) {
  const now = new Date();
  const settings = await getBillingSetting();
  const service = await db.serviceInstance.findFirst({
    where: {
      id: serviceId,
      billingCycle: {
        not: "ONETIME",
      },
      status: {
        in: ["ACTIVE", "OVERDUE", "SUSPENDED", "PROVISIONING"],
      },
    },
    include: {
      customer: true,
      product: true,
    },
  });

  if (!service) {
    return null;
  }

  return createRenewalInvoiceForService(service, {
    now,
    actorId: input?.actorId,
    source: input?.source ?? "portal",
    orderNote: `客户在门户为服务 ${service.serviceNo} 创建续费订单`,
    invoiceRemark: `客户在门户为服务 ${service.serviceNo} 创建续费账单`,
    jobMessagePrefix: "已生成续费账单",
    invoicePrefix: settings.invoicePrefix || "INV",
  });
}

async function createRenewalInvoices(
  settings: BillingSetting,
  now: Date,
  actorId: string | null | undefined,
  summary: BillingSummary,
) {
  const cutoff = addDays(startOfDay(now), settings.invoiceLeadDays);

  const services = await db.serviceInstance.findMany({
    where: {
      billingCycle: {
        not: "ONETIME",
      },
      nextDueDate: {
        not: null,
        lte: cutoff,
      },
      status: {
        in: ["ACTIVE", "OVERDUE", "SUSPENDED", "PROVISIONING"],
      },
    },
    include: {
      customer: true,
      product: true,
    },
    orderBy: {
      nextDueDate: "asc",
    },
  });

  for (const service of services) {
    const created = await createRenewalInvoiceForService(service, {
      now,
      actorId,
      source: "billing-engine",
      orderNote: `系统为服务 ${service.serviceNo} 自动生成续费订单`,
      invoiceRemark: `系统自动生成 ${service.serviceNo} 的续费账单`,
      jobMessagePrefix: "已生成续费账单",
      invoicePrefix: settings.invoicePrefix || "INV",
    });

    if (created?.created) {
      summary.createdInvoices += 1;
      summary.jobs += 1;
    }
  }
}

async function settleRenewalByBalance(input: {
  invoiceId: string;
  actorId?: string | null;
  allowNegativeBalance: boolean;
  summary: BillingSummary;
}) {
  const invoice = await db.invoice.findUnique({
    where: {
      id: input.invoiceId,
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
    return;
  }

  const outstanding = Math.max(invoice.totalAmount - invoice.paidAmount, 0);

  if (outstanding <= 0) {
    return;
  }

  await db.$transaction(async (tx) => {
      await changeCustomerBalance({
        customerId: invoice.customerId,
        amount: -outstanding,
        type: "AUTO_RENEW",
        description: `余额自动续费 ${invoice.invoiceNo}`,
        operatorId: input.actorId ?? undefined,
        allowNegative: input.allowNegativeBalance,
        tx,
    });

    await tx.payment.create({
      data: {
        paymentNo: makeCode("PAY"),
        customerId: invoice.customerId,
        invoiceId: invoice.id,
        orderId: invoice.orderId,
        method: "BALANCE",
        status: "SUCCESS",
        amount: outstanding,
        transactionNo: `BAL-${Date.now()}`,
        paidAt: new Date(),
      },
    });

    await tx.invoice.update({
      where: {
        id: invoice.id,
      },
      data: {
        status: "PAID",
        paidAmount: invoice.totalAmount,
        paidAt: new Date(),
      },
    });

    if (invoice.orderId) {
      await tx.order.update({
        where: {
          id: invoice.orderId,
        },
        data: {
          status: "ACTIVE",
          paidAmount: (invoice.order?.paidAmount ?? 0) + outstanding,
          paidAt: new Date(),
        },
      });
    }
  });

  if (invoice.service) {
    if (invoice.service.providerType === "MOFANG_CLOUD") {
      const provider = getCloudProvider();
      const payload = buildProviderPayload(invoice.service);
      const providerResult = await provider.renewService(payload);

      await syncServiceRecord({
        service: invoice.service,
        providerResult,
        fallbackStatus: "ACTIVE",
        nextDueDate: getNextDueDate(
          invoice.service.nextDueDate,
          invoice.service.billingCycle,
          providerResult.nextDueDate,
        ),
      });

      await writeProviderLog({
        action: "renew",
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
          nextDueDate: getNextDueDate(
            invoice.service.nextDueDate,
            invoice.service.billingCycle,
          ),
          lastSyncAt: new Date(),
        },
      });
    }
  }

  await refreshCustomerStatus(invoice.customerId);

  await createBillingJob({
    jobType: "AUTO_RENEW_BALANCE",
    status: "SUCCESS",
    customerId: invoice.customerId,
    serviceId: invoice.serviceId,
    invoiceId: invoice.id,
    message: `账单 ${invoice.invoiceNo} 已通过余额自动结清`,
    payload: {
      outstanding,
    },
  });

  await writeAuditLog({
    adminUserId: input.actorId,
    module: "billing",
    action: "auto-renew-balance",
    targetType: "invoice",
    targetId: invoice.id,
    summary: `余额自动续费账单：${invoice.invoiceNo}`,
    detail: {
      serviceNo: invoice.service?.serviceNo,
      amount: outstanding,
    },
  });

  input.summary.autoRenewed += 1;
  input.summary.jobs += 1;
}

async function markInvoiceOverdue(input: {
  invoiceId: string;
  actorId?: string | null;
  summary: BillingSummary;
}) {
  const invoice = await db.invoice.findUnique({
    where: {
      id: input.invoiceId,
    },
    include: {
      service: true,
    },
  });

  if (!invoice || invoice.status === "OVERDUE" || invoice.status === "PAID") {
    return;
  }

  await db.invoice.update({
    where: {
      id: invoice.id,
    },
    data: {
      status: "OVERDUE",
    },
  });

  if (invoice.serviceId && invoice.service && !["SUSPENDED", "TERMINATED"].includes(invoice.service.status)) {
    await db.serviceInstance.update({
      where: {
        id: invoice.serviceId,
      },
      data: {
        status: "OVERDUE",
        lastSyncAt: new Date(),
      },
    });
  }

  await refreshCustomerStatus(invoice.customerId);

  await createBillingJob({
    jobType: "MARK_OVERDUE",
    status: "SUCCESS",
    customerId: invoice.customerId,
    serviceId: invoice.serviceId,
    invoiceId: invoice.id,
    message: `账单 ${invoice.invoiceNo} 已标记为逾期`,
  });

  await writeAuditLog({
    adminUserId: input.actorId,
    module: "billing",
    action: "mark-overdue",
    targetType: "invoice",
    targetId: invoice.id,
    summary: `标记账单逾期：${invoice.invoiceNo}`,
  });

  input.summary.markedOverdue += 1;
  input.summary.jobs += 1;
}

async function updateServiceLifecycle(input: {
  invoiceId: string;
  targetStatus: ServiceStatus;
  jobType: "AUTO_SUSPEND" | "AUTO_TERMINATE";
  actorId?: string | null;
  summary: BillingSummary;
}) {
  const invoice = await db.invoice.findUnique({
    where: {
      id: input.invoiceId,
    },
    include: {
      service: {
        include: {
          customer: true,
          product: true,
        },
      },
    },
  });

  if (!invoice?.service) {
    return;
  }

  if (
    input.targetStatus === "SUSPENDED" &&
    ["SUSPENDED", "TERMINATED"].includes(invoice.service.status)
  ) {
    return;
  }

  if (input.targetStatus === "TERMINATED" && invoice.service.status === "TERMINATED") {
    return;
  }

  if (invoice.service.providerType === "MOFANG_CLOUD") {
    const provider = getCloudProvider();
    const payload = buildProviderPayload(invoice.service);
    const providerResult =
      input.targetStatus === "TERMINATED"
        ? await provider.terminateService(payload)
        : await provider.suspendService(payload);

    await syncServiceRecord({
      service: invoice.service,
      providerResult,
      fallbackStatus: input.targetStatus,
    });

    await writeProviderLog({
      action: input.targetStatus === "TERMINATED" ? "terminate" : "suspend",
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
        status: input.targetStatus,
        expiresAt: input.targetStatus === "TERMINATED" ? new Date() : invoice.service.expiresAt,
        lastSyncAt: new Date(),
      },
    });
  }

  await createBillingJob({
    jobType: input.jobType,
    status: "SUCCESS",
    customerId: invoice.customerId,
    serviceId: invoice.serviceId,
    invoiceId: invoice.id,
    message: `服务 ${invoice.service.serviceNo} 状态变更为 ${serviceStatusLabelMap[input.targetStatus]}`,
  });

  await writeAuditLog({
    adminUserId: input.actorId,
    module: "service",
    action: input.targetStatus === "TERMINATED" ? "auto-terminate" : "auto-suspend",
    targetType: "service",
    targetId: invoice.service.id,
    summary: `自动处理服务状态：${invoice.service.serviceNo} / ${serviceStatusLabelMap[input.targetStatus]}`,
    detail: {
      invoiceNo: invoice.invoiceNo,
    },
  });

  if (input.targetStatus === "TERMINATED") {
    input.summary.terminated += 1;
  } else {
    input.summary.suspended += 1;
  }

  input.summary.jobs += 1;
}

export async function runBillingEngine(actorId?: string | null) {
  const now = new Date();
  const settings = await getBillingSetting();
  const summary: BillingSummary = {
    createdInvoices: 0,
    autoRenewed: 0,
    markedOverdue: 0,
    suspended: 0,
    terminated: 0,
    jobs: 0,
  };

  await createRenewalInvoices(settings, now, actorId, summary);

  const invoices = await db.invoice.findMany({
    where: {
      status: {
        in: OPEN_INVOICE_STATUSES,
      },
    },
    include: {
      customer: true,
      service: {
        include: {
          customer: true,
          product: true,
        },
      },
    },
    orderBy: {
      dueDate: "asc",
    },
  });

  for (const invoice of invoices) {
    const outstanding = Math.max(invoice.totalAmount - invoice.paidAmount, 0);

    if (outstanding <= 0) {
      continue;
    }

    const canAutoRenew =
      settings.autoRenewByBalance &&
      invoice.type === "RENEWAL" &&
      (invoice.customer.creditBalance >= outstanding || settings.allowNegativeBalance);

    if (canAutoRenew) {
      await settleRenewalByBalance({
        invoiceId: invoice.id,
        actorId,
        allowNegativeBalance: settings.allowNegativeBalance,
        summary,
      });
      continue;
    }

    const overdueDays = differenceInCalendarDays(
      startOfDay(now),
      startOfDay(invoice.dueDate),
    );

    if (overdueDays > 0) {
      await markInvoiceOverdue({
        invoiceId: invoice.id,
        actorId,
        summary,
      });
    }

    if (invoice.serviceId && overdueDays >= settings.autoTerminateDays) {
      await updateServiceLifecycle({
        invoiceId: invoice.id,
        targetStatus: "TERMINATED",
        jobType: "AUTO_TERMINATE",
        actorId,
        summary,
      });
      continue;
    }

    if (invoice.serviceId && overdueDays >= settings.autoSuspendDays) {
      await updateServiceLifecycle({
        invoiceId: invoice.id,
        targetStatus: "SUSPENDED",
        jobType: "AUTO_SUSPEND",
        actorId,
        summary,
      });
    }
  }

  const recentJobs = await db.billingJob.findMany({
    include: {
      customer: {
        select: {
          name: true,
        },
      },
      invoice: {
        select: {
          invoiceNo: true,
        },
      },
      service: {
        select: {
          name: true,
          serviceNo: true,
        },
      },
    },
    orderBy: {
      executedAt: "desc",
    },
    take: 12,
  });

  return {
    settings,
    summary,
    recentJobs,
  };
}
