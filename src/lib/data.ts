import {
  eachDayOfInterval,
  endOfDay,
  endOfMonth,
  format,
  startOfDay,
  startOfMonth,
  subDays,
  subMonths,
} from "date-fns";

import { getBillingSetting } from "@/lib/billing-engine";
import { db } from "@/lib/db";

function sumMoney(values: number[]) {
  return values.reduce((total, value) => total + value, 0);
}

async function getAuditLogsForTarget(targetType: string, targetId: string, take = 20) {
  return db.auditLog.findMany({
    where: {
      targetType,
      targetId,
    },
    include: {
      adminUser: {
        select: {
          name: true,
          role: true,
        },
      },
    },
    orderBy: {
      createdAt: "desc",
    },
    take,
  });
}

export async function getDashboardData() {
  const now = new Date();
  const monthStart = startOfMonth(now);
  const monthEnd = endOfMonth(now);
  const sixMonthsAgo = startOfMonth(subMonths(now, 5));

  const [
    customers,
    services,
    invoices,
    payments,
    logs,
    billingJobs,
    vpcCount,
    ipCount,
    diskCount,
    snapshotCount,
    backupCount,
    securityGroupCount,
  ] = await Promise.all([
    db.customer.findMany({
      include: {
        invoices: {
          select: {
            totalAmount: true,
            paidAmount: true,
          },
        },
        _count: {
          select: {
            services: true,
          },
        },
      },
    }),
    db.serviceInstance.findMany({
      include: {
        customer: {
          select: {
            name: true,
          },
        },
        product: {
          select: {
            name: true,
          },
        },
      },
      orderBy: {
        nextDueDate: "asc",
      },
    }),
    db.invoice.findMany({
      orderBy: {
        dueDate: "asc",
      },
    }),
    db.payment.findMany({
      where: {
        status: "SUCCESS",
        paidAt: {
          not: null,
          gte: sixMonthsAgo,
        },
      },
      orderBy: {
        paidAt: "asc",
      },
    }),
    db.providerSyncLog.findMany({
      orderBy: {
        syncedAt: "desc",
      },
      take: 6,
    }),
    db.billingJob.findMany({
      include: {
        customer: {
          select: {
            name: true,
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
      take: 6,
    }),
    db.serviceVpcNetwork.count(),
    db.serviceIpAddress.count(),
    db.serviceDisk.count(),
    db.serviceSnapshot.count(),
    db.serviceBackup.count(),
    db.serviceSecurityGroup.count(),
  ]);

  const monthlyRevenue = sumMoney(
    payments
      .filter(
        (payment) =>
          payment.paidAt &&
          payment.paidAt >= monthStart &&
          payment.paidAt <= monthEnd,
      )
      .map((payment) => payment.amount),
  );

  const outstandingAmount = sumMoney(
    invoices
      .filter((invoice) => !["PAID", "VOID"].includes(invoice.status))
      .map((invoice) => Math.max(invoice.totalAmount - invoice.paidAmount, 0)),
  );

  const activeServices = services.filter((service) =>
    ["ACTIVE", "PROVISIONING"].includes(service.status),
  ).length;

  const overdueCustomers = customers.filter(
    (customer) => customer.status === "OVERDUE",
  ).length;

  const totalBalance = sumMoney(customers.map((customer) => customer.creditBalance));

  const openRenewalInvoices = invoices.filter(
    (invoice) =>
      invoice.type === "RENEWAL" &&
      ["ISSUED", "PARTIAL", "OVERDUE"].includes(invoice.status),
  ).length;

  const revenueTrend = Array.from({ length: 6 }).map((_, index) => {
    const target = startOfMonth(subMonths(now, 5 - index));
    const monthKey = format(target, "yyyy-MM");
    const monthPayments = payments.filter(
      (payment) => payment.paidAt && format(payment.paidAt, "yyyy-MM") === monthKey,
    );

    return {
      month: monthKey,
      amount: sumMoney(monthPayments.map((payment) => payment.amount)),
    };
  });

  const topCustomers = customers
    .map((customer) => ({
      id: customer.id,
      name: customer.name,
      creditBalance: customer.creditBalance,
      serviceCount: customer._count.services,
      billedAmount: sumMoney(customer.invoices.map((invoice) => invoice.totalAmount)),
      outstanding: sumMoney(
        customer.invoices.map((invoice) =>
          Math.max(invoice.totalAmount - invoice.paidAmount, 0),
        ),
      ),
      status: customer.status,
    }))
    .sort((a, b) => b.billedAmount - a.billedAmount)
    .slice(0, 5);

  const renewals = services
    .filter((service) => Boolean(service.nextDueDate))
    .slice(0, 8)
    .map((service) => ({
      id: service.id,
      name: service.name,
      customerName: service.customer.name,
      productName: service.product.name,
      status: service.status,
      nextDueDate: service.nextDueDate,
      amount: service.monthlyCost,
    }));

  return {
    metrics: {
      monthlyRevenue,
      outstandingAmount,
      activeServices,
      overdueCustomers,
      totalBalance,
      openRenewalInvoices,
    },
    revenueTrend,
    topCustomers,
    renewals,
    invoices: invoices.slice(0, 8),
    providerLogs: logs,
    billingJobs,
    resourceSummary: {
      vpcs: vpcCount,
      ips: ipCount,
      disks: diskCount,
      snapshots: snapshotCount,
      backups: backupCount,
      securityGroups: securityGroupCount,
    },
  };
}

export async function getReportsOverviewData() {
  const now = new Date();
  const today = endOfDay(now);
  const last14DaysStart = startOfDay(subDays(now, 13));
  const last30DaysStart = startOfDay(subDays(now, 29));
  const monthStart = startOfMonth(now);

  const [payments, refunds, invoices, services, customers] = await Promise.all([
    db.payment.findMany({
      where: {
        status: "SUCCESS",
        paidAt: {
          not: null,
          gte: last30DaysStart,
          lte: today,
        },
      },
      include: {
        invoice: {
          include: {
            service: {
              include: {
                product: {
                  select: {
                    name: true,
                  },
                },
              },
            },
          },
        },
      },
      orderBy: {
        paidAt: "asc",
      },
    }),
    db.refundRecord.findMany({
      where: {
        status: "SUCCESS",
        processedAt: {
          not: null,
          gte: last30DaysStart,
          lte: today,
        },
      },
      orderBy: {
        processedAt: "asc",
      },
    }),
    db.invoice.findMany({
      include: {
        customer: {
          select: {
            id: true,
            name: true,
          },
        },
      },
      orderBy: {
        dueDate: "asc",
      },
    }),
    db.serviceInstance.findMany({
      include: {
        customer: {
          select: {
            name: true,
          },
        },
        product: {
          select: {
            name: true,
          },
        },
      },
      orderBy: {
        nextDueDate: "asc",
      },
    }),
    db.customer.findMany({
      include: {
        invoices: {
          select: {
            status: true,
            totalAmount: true,
            paidAmount: true,
          },
        },
      },
    }),
  ]);

  const monthlyRevenue = sumMoney(
    payments
      .filter((payment) => payment.paidAt && payment.paidAt >= monthStart)
      .map((payment) => payment.amount),
  );
  const monthlyRefunded = sumMoney(
    refunds
      .filter((refund) => refund.processedAt && refund.processedAt >= monthStart)
      .map((refund) => refund.amount),
  );
  const outstandingAmount = sumMoney(
    invoices
      .filter((invoice) => !["PAID", "VOID"].includes(invoice.status))
      .map((invoice) => Math.max(invoice.totalAmount - invoice.paidAmount, 0)),
  );

  const dailyRevenue = eachDayOfInterval({
    start: last14DaysStart,
    end: today,
  }).map((day) => {
    const dayStart = startOfDay(day);
    const dayEnd = endOfDay(day);
    const matchedPayments = payments.filter(
      (payment) => payment.paidAt && payment.paidAt >= dayStart && payment.paidAt <= dayEnd,
    );
    const matchedRefunds = refunds.filter(
      (refund) =>
        refund.processedAt &&
        refund.processedAt >= dayStart &&
        refund.processedAt <= dayEnd,
    );

    return {
      date: format(day, "MM-dd"),
      revenue: sumMoney(matchedPayments.map((payment) => payment.amount)),
      refunds: sumMoney(matchedRefunds.map((refund) => refund.amount)),
    };
  });

  const paymentChannels = Object.values(
    payments.reduce<Record<string, { method: string; amount: number; count: number }>>(
      (acc, payment) => {
        if (!acc[payment.method]) {
          acc[payment.method] = {
            method: payment.method,
            amount: 0,
            count: 0,
          };
        }

        acc[payment.method].amount += payment.amount;
        acc[payment.method].count += 1;
        return acc;
      },
      {},
    ),
  ).sort((a, b) => b.amount - a.amount);

  const invoiceStatusSummary = [
    "DRAFT",
    "ISSUED",
    "PARTIAL",
    "PAID",
    "OVERDUE",
    "VOID",
  ].map((status) => ({
    status,
    count: invoices.filter((invoice) => invoice.status === status).length,
    amount: sumMoney(
      invoices
        .filter((invoice) => invoice.status === status)
        .map((invoice) => invoice.totalAmount),
    ),
  }));

  const productRanking = Object.values(
    payments.reduce<
      Record<string, { productName: string; amount: number; count: number }>
    >((acc, payment) => {
      const productName = payment.invoice?.service?.product?.name || "未关联产品";
      if (!acc[productName]) {
        acc[productName] = {
          productName,
          amount: 0,
          count: 0,
        };
      }
      acc[productName].amount += payment.amount;
      acc[productName].count += 1;
      return acc;
    }, {}),
  )
    .sort((a, b) => b.amount - a.amount)
    .slice(0, 8);

  const regionDistribution = Object.values(
    services.reduce<Record<string, { region: string; count: number; amount: number }>>(
      (acc, service) => {
        const region = service.region || "未分配地域";
        if (!acc[region]) {
          acc[region] = {
            region,
            count: 0,
            amount: 0,
          };
        }
        acc[region].count += 1;
        acc[region].amount += service.monthlyCost;
        return acc;
      },
      {},
    ),
  ).sort((a, b) => b.count - a.count);

  const overdueServices = services
    .filter((service) => ["OVERDUE", "SUSPENDED", "EXPIRED"].includes(service.status))
    .slice(0, 10)
    .map((service) => ({
      id: service.id,
      serviceNo: service.serviceNo,
      name: service.name,
      customerName: service.customer.name,
      productName: service.product.name,
      status: service.status,
      nextDueDate: service.nextDueDate,
      amount: service.monthlyCost,
    }));

  const customerArRanking = customers
    .map((customer) => {
      const openInvoices = customer.invoices.filter((invoice) =>
        ["ISSUED", "PARTIAL", "OVERDUE"].includes(invoice.status),
      );
      return {
        id: customer.id,
        name: customer.name,
        outstandingAmount: sumMoney(
          openInvoices.map((invoice) =>
            Math.max(invoice.totalAmount - invoice.paidAmount, 0),
          ),
        ),
        openInvoiceCount: openInvoices.length,
      };
    })
    .filter((customer) => customer.outstandingAmount > 0)
    .sort((a, b) => b.outstandingAmount - a.outstandingAmount)
    .slice(0, 8);

  return {
    summary: {
      monthlyRevenue,
      monthlyRefunded,
      netRevenue: monthlyRevenue - monthlyRefunded,
      outstandingAmount,
      overdueServiceCount: overdueServices.length,
      activeServiceCount: services.filter((service) => service.status === "ACTIVE").length,
    },
    dailyRevenue,
    paymentChannels,
    invoiceStatusSummary,
    productRanking,
    regionDistribution,
    overdueServices,
    customerArRanking,
  };
}

export async function getCustomersPageData() {
  const customers = await db.customer.findMany({
    include: {
      _count: {
        select: {
          orders: true,
          services: true,
          tickets: true,
        },
      },
      invoices: {
        select: {
          totalAmount: true,
          paidAmount: true,
        },
      },
      creditTransactions: {
        select: {
          amount: true,
        },
      },
    },
    orderBy: {
      createdAt: "desc",
    },
  });

  return customers.map((customer) => ({
    ...customer,
    totalBilled: sumMoney(customer.invoices.map((invoice) => invoice.totalAmount)),
    totalPaid: sumMoney(customer.invoices.map((invoice) => invoice.paidAmount)),
    totalCreditIn: sumMoney(
      customer.creditTransactions
        .filter((transaction) => transaction.amount > 0)
        .map((transaction) => transaction.amount),
    ),
    totalCreditOut: Math.abs(
      sumMoney(
        customer.creditTransactions
          .filter((transaction) => transaction.amount < 0)
          .map((transaction) => transaction.amount),
      ),
    ),
  }));
}

export async function getProductsPageData() {
  return db.product.findMany({
    include: {
      _count: {
        select: {
          services: true,
          items: true,
          plans: true,
        },
      },
    },
    orderBy: [
      {
        status: "asc",
      },
      {
        createdAt: "desc",
      },
    ],
  });
}

export async function getOrdersPageData() {
  const [orders, customers, products] = await Promise.all([
    db.order.findMany({
      include: {
        customer: true,
        items: true,
        services: true,
        invoices: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.customer.findMany({
      orderBy: {
        name: "asc",
      },
    }),
    db.product.findMany({
      where: {
        status: "ACTIVE",
      },
      orderBy: {
        name: "asc",
      },
    }),
  ]);

  return { orders, customers, products };
}

export async function getServicesPageData() {
  return db.serviceInstance.findMany({
    include: {
      customer: true,
      plan: {
        include: {
          region: true,
          zone: true,
          flavor: true,
          image: true,
        },
      },
      product: true,
      order: true,
      vpcNetwork: true,
      _count: {
        select: {
          ipAddresses: true,
          disks: true,
          snapshots: true,
          backups: true,
          securityGroups: true,
        },
      },
    },
    orderBy: [
      {
        status: "asc",
      },
      {
        nextDueDate: "asc",
      },
    ],
  });
}

export async function getInvoicesPageData() {
  const [invoices, customers, services, taxProfiles, billingSetting] = await Promise.all([
    db.invoice.findMany({
      include: {
        customer: true,
        order: true,
        service: true,
        payments: {
          orderBy: {
            createdAt: "desc",
          },
        },
        refunds: {
          orderBy: {
            createdAt: "desc",
          },
        },
      },
      orderBy: [
        {
          createdAt: "desc",
        },
      ],
    }),
    db.customer.findMany({
      where: {
        status: {
          in: ["ACTIVE", "OVERDUE", "SUSPENDED"],
        },
      },
      orderBy: {
        name: "asc",
      },
    }),
    db.serviceInstance.findMany({
      include: {
        customer: {
          select: {
            id: true,
            name: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.invoiceTaxProfile.findMany({
      where: {
        isActive: true,
      },
      orderBy: [
        {
          isDefault: "desc",
        },
        {
          taxRate: "asc",
        },
      ],
    }),
    db.billingSetting.upsert({
      where: {
        id: "default",
      },
      update: {},
      create: {
        id: "default",
      },
    }),
  ]);

  const totalTax = sumMoney(invoices.map((invoice) => invoice.taxAmount));
  const totalIssued = sumMoney(
    invoices
      .filter((invoice) => invoice.status !== "VOID")
      .map((invoice) => invoice.totalAmount),
  );
  const totalOutstanding = sumMoney(
    invoices
      .filter((invoice) => !["PAID", "VOID"].includes(invoice.status))
      .map((invoice) => Math.max(invoice.totalAmount - invoice.paidAmount, 0)),
  );

  return {
    invoices,
    customers,
    services,
    taxProfiles,
    billingSetting,
    summary: {
      invoiceCount: invoices.length,
      draftCount: invoices.filter((invoice) => invoice.status === "DRAFT").length,
      openCount: invoices.filter((invoice) =>
        ["ISSUED", "PARTIAL", "OVERDUE"].includes(invoice.status),
      ).length,
      overdueCount: invoices.filter((invoice) => invoice.status === "OVERDUE").length,
      paidCount: invoices.filter((invoice) => invoice.status === "PAID").length,
      voidCount: invoices.filter((invoice) => invoice.status === "VOID").length,
      totalIssued,
      totalTax,
      totalOutstanding,
    },
  };
}

export async function getCustomerDetailData(customerId: string) {
  const [customer, auditLogs] = await Promise.all([
    db.customer.findUnique({
      where: {
        id: customerId,
      },
      include: {
        _count: {
          select: {
            orders: true,
            services: true,
            tickets: true,
            portalUsers: true,
            contacts: true,
            followUps: true,
          },
        },
        contacts: {
          orderBy: [
            {
              isPrimary: "desc",
            },
            {
              createdAt: "asc",
            },
          ],
        },
        certification: true,
        followUps: {
          orderBy: {
            createdAt: "desc",
          },
          take: 20,
        },
        portalUsers: {
          orderBy: [
            {
              isOwner: "desc",
            },
            {
              createdAt: "asc",
            },
          ],
        },
        orders: {
          include: {
            items: {
              include: {
                product: {
                  select: {
                    name: true,
                    category: true,
                  },
                },
                service: {
                  select: {
                    id: true,
                    serviceNo: true,
                    name: true,
                    status: true,
                  },
                },
              },
            },
            invoices: {
              select: {
                id: true,
                invoiceNo: true,
                status: true,
                totalAmount: true,
                paidAmount: true,
                dueDate: true,
              },
            },
          },
          orderBy: {
            createdAt: "desc",
          },
          take: 12,
        },
        services: {
          include: {
            product: {
              select: {
                name: true,
                category: true,
              },
            },
            plan: {
              include: {
                region: true,
                zone: true,
                flavor: true,
              },
            },
            order: {
              select: {
                id: true,
                orderNo: true,
                status: true,
              },
            },
            _count: {
              select: {
                ipAddresses: true,
                disks: true,
                snapshots: true,
                backups: true,
                securityGroups: true,
              },
            },
          },
          orderBy: [
            {
              status: "asc",
            },
            {
              nextDueDate: "asc",
            },
          ],
          take: 12,
        },
        invoices: {
          include: {
            service: {
              select: {
                id: true,
                serviceNo: true,
                name: true,
                status: true,
              },
            },
            payments: {
              orderBy: {
                createdAt: "desc",
              },
              take: 3,
            },
          },
          orderBy: {
            createdAt: "desc",
          },
          take: 12,
        },
        payments: {
          include: {
            invoice: {
              select: {
                id: true,
                invoiceNo: true,
                status: true,
              },
            },
            refunds: {
              include: {
                processedBy: {
                  select: {
                    name: true,
                    role: true,
                  },
                },
              },
              orderBy: {
                createdAt: "desc",
              },
            },
          },
          orderBy: {
            createdAt: "desc",
          },
          take: 12,
        },
        refunds: {
          include: {
            processedBy: {
              select: {
                name: true,
                role: true,
              },
            },
            payment: {
              select: {
                paymentNo: true,
                method: true,
              },
            },
            invoice: {
              select: {
                id: true,
                invoiceNo: true,
                status: true,
              },
            },
          },
          orderBy: {
            createdAt: "desc",
          },
          take: 12,
        },
        creditTransactions: {
          include: {
            operator: {
              select: {
                name: true,
                role: true,
              },
            },
          },
          orderBy: {
            createdAt: "desc",
          },
          take: 20,
        },
        tickets: {
          include: {
            assignedTo: {
              select: {
                name: true,
                role: true,
              },
            },
          },
          orderBy: {
            updatedAt: "desc",
          },
          take: 12,
        },
      },
    }),
    getAuditLogsForTarget("customer", customerId),
  ]);

  if (!customer) {
    return null;
  }

  const totalBilled = sumMoney(customer.invoices.map((invoice) => invoice.totalAmount));
  const totalPaid = sumMoney(customer.invoices.map((invoice) => invoice.paidAmount));
  const totalOutstanding = sumMoney(
    customer.invoices.map((invoice) =>
      Math.max(invoice.totalAmount - invoice.paidAmount, 0),
    ),
  );
  const totalCreditIn = sumMoney(
    customer.creditTransactions
      .filter((transaction) => transaction.amount > 0)
      .map((transaction) => transaction.amount),
  );
  const totalCreditOut = Math.abs(
    sumMoney(
      customer.creditTransactions
        .filter((transaction) => transaction.amount < 0)
        .map((transaction) => transaction.amount),
    ),
  );
  const activeServices = customer.services.filter((service) =>
    ["ACTIVE", "PROVISIONING"].includes(service.status),
  ).length;
  const overdueInvoices = customer.invoices.filter(
    (invoice) => invoice.status === "OVERDUE",
  ).length;
  const nextFollowUp = customer.followUps
    .filter((followUp) => followUp.nextFollowAt)
    .sort((a, b) => Number(a.nextFollowAt) - Number(b.nextFollowAt))[0] ?? null;

  return {
    customer,
    summary: {
      totalBilled,
      totalPaid,
      totalOutstanding,
      totalCreditIn,
      totalCreditOut,
      activeServices,
      overdueInvoices,
      refundCount: customer.refunds.length,
      paymentCount: customer.payments.length,
      nextFollowUpAt: nextFollowUp?.nextFollowAt ?? null,
    },
    auditLogs,
  };
}

export async function getOrderDetailData(orderId: string) {
  const [order, auditLogs] = await Promise.all([
    db.order.findUnique({
      where: {
        id: orderId,
      },
      include: {
        customer: true,
        items: {
          include: {
            product: true,
            plan: {
              include: {
                region: true,
                zone: true,
                flavor: true,
              },
            },
            service: {
              select: {
                id: true,
                serviceNo: true,
                name: true,
                status: true,
                nextDueDate: true,
              },
            },
          },
        },
        invoices: {
          include: {
            payments: {
              orderBy: {
                createdAt: "desc",
              },
            },
          },
          orderBy: {
            createdAt: "desc",
          },
        },
        payments: {
          include: {
            refunds: {
              orderBy: {
                createdAt: "desc",
              },
            },
          },
          orderBy: {
            createdAt: "desc",
          },
        },
        services: {
          include: {
            product: {
              select: {
                name: true,
                category: true,
              },
            },
            plan: {
              include: {
                region: true,
                zone: true,
                flavor: true,
              },
            },
          },
        },
      },
    }),
    getAuditLogsForTarget("order", orderId),
  ]);

  if (!order) {
    return null;
  }

  const totalInvoiced = sumMoney(order.invoices.map((invoice) => invoice.totalAmount));
  const totalPaid = sumMoney(order.payments.map((payment) => payment.amount));
  const totalRefunded = sumMoney(
    order.payments.flatMap((payment) => payment.refunds.map((refund) => refund.amount)),
  );

  return {
    order,
    summary: {
      itemCount: order.items.length,
      serviceCount: order.services.length,
      invoiceCount: order.invoices.length,
      paymentCount: order.payments.length,
      totalInvoiced,
      totalPaid,
      totalRefunded,
      outstandingAmount: Math.max(order.totalAmount - order.paidAmount, 0),
    },
    auditLogs,
  };
}

export async function getServiceDetailData(serviceId: string) {
  const service = await db.serviceInstance.findUnique({
    where: {
      id: serviceId,
    },
    include: {
      customer: true,
      product: true,
      order: {
        select: {
          id: true,
          orderNo: true,
          status: true,
          totalAmount: true,
          paidAmount: true,
          createdAt: true,
        },
      },
      plan: {
        include: {
          region: true,
          zone: true,
          flavor: true,
          image: true,
        },
      },
      vpcNetwork: true,
      invoices: {
        include: {
          payments: {
            orderBy: {
              createdAt: "desc",
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
      },
      tickets: {
        include: {
          assignedTo: {
            select: {
              name: true,
              role: true,
            },
          },
        },
        orderBy: {
          updatedAt: "desc",
        },
      },
      ipAddresses: {
        orderBy: {
          createdAt: "desc",
        },
      },
      disks: {
        orderBy: {
          createdAt: "desc",
        },
      },
      snapshots: {
        orderBy: {
          createdAt: "desc",
        },
      },
      backups: {
        orderBy: {
          createdAt: "desc",
        },
      },
      securityGroups: {
        include: {
          rules: {
            orderBy: {
              createdAt: "desc",
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
      },
      billingJobs: {
        orderBy: {
          executedAt: "desc",
        },
        take: 12,
      },
    },
  });

  if (!service) {
    return null;
  }

  const [auditLogs, providerLogs] = await Promise.all([
    getAuditLogsForTarget("service", serviceId),
    service.providerResourceId
      ? db.providerSyncLog.findMany({
          where: {
            resourceId: service.providerResourceId,
          },
          orderBy: {
            syncedAt: "desc",
          },
          take: 20,
        })
      : Promise.resolve([]),
  ]);

  const outstandingAmount = sumMoney(
    service.invoices.map((invoice) =>
      Math.max(invoice.totalAmount - invoice.paidAmount, 0),
    ),
  );

  return {
    service,
    summary: {
      invoiceCount: service.invoices.length,
      ticketCount: service.tickets.length,
      ipCount: service.ipAddresses.length,
      diskCount: service.disks.length,
      snapshotCount: service.snapshots.length,
      backupCount: service.backups.length,
      securityGroupCount: service.securityGroups.length,
      outstandingAmount,
    },
    auditLogs,
    providerLogs,
  };
}

export async function getInvoiceDetailData(invoiceId: string) {
  const invoice = await db.invoice.findUnique({
    where: {
      id: invoiceId,
    },
    include: {
      customer: true,
      order: {
        select: {
          id: true,
          orderNo: true,
          status: true,
          totalAmount: true,
          paidAmount: true,
          createdAt: true,
        },
      },
      service: {
        select: {
          id: true,
          serviceNo: true,
          name: true,
          status: true,
          nextDueDate: true,
        },
      },
      payments: {
        include: {
          refunds: {
            orderBy: {
              createdAt: "desc",
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
      },
      refunds: {
        include: {
          processedBy: {
            select: {
              name: true,
              role: true,
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
      },
      billingJobs: {
        orderBy: {
          executedAt: "desc",
        },
        take: 12,
      },
    },
  });

  if (!invoice) {
    return null;
  }

  const paymentIds = invoice.payments.map((payment) => payment.id);
  const [auditLogs, callbackLogs] = await Promise.all([
    getAuditLogsForTarget("invoice", invoiceId),
    paymentIds.length > 0
      ? db.paymentCallbackLog.findMany({
          where: {
            paymentId: {
              in: paymentIds,
            },
          },
          orderBy: {
            createdAt: "desc",
          },
          take: 20,
        })
      : Promise.resolve([]),
  ]);

  const refundedAmount = sumMoney(invoice.refunds.map((refund) => refund.amount));

  return {
    invoice,
    summary: {
      paymentCount: invoice.payments.length,
      refundCount: invoice.refunds.length,
      callbackCount: callbackLogs.length,
      outstandingAmount: Math.max(invoice.totalAmount - invoice.paidAmount, 0),
      refundedAmount,
    },
    auditLogs,
    callbackLogs,
  };
}

export async function getPaymentsPageData() {
  const [payments, invoices, callbackLogs, refunds] = await Promise.all([
    db.payment.findMany({
      include: {
        customer: true,
        invoice: true,
        order: true,
        refunds: {
          orderBy: {
            createdAt: "desc",
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.invoice.findMany({
      where: {
        status: {
          in: ["ISSUED", "PARTIAL", "OVERDUE"],
        },
      },
      include: {
        customer: true,
        order: true,
      },
      orderBy: {
        dueDate: "asc",
      },
    }),
    db.paymentCallbackLog.findMany({
      include: {
        payment: {
          select: {
            id: true,
            paymentNo: true,
            amount: true,
            status: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 20,
    }),
    db.refundRecord.findMany({
      include: {
        customer: true,
        payment: true,
        invoice: true,
        processedBy: true,
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 30,
    }),
  ]);

  const totalSuccess = sumMoney(
    payments
      .filter((payment) => payment.status === "SUCCESS")
      .map((payment) => payment.amount),
  );

  const pendingCallbacks = callbackLogs.filter((log) => !log.isHandled).length;
  const totalRefunded = sumMoney(
    refunds
      .filter((refund) => refund.status === "SUCCESS")
      .map((refund) => refund.amount),
  );
  const pendingRefunds = refunds.filter((refund) => refund.status === "PENDING").length;

  return {
    payments,
    invoices,
    callbackLogs,
    refunds,
    summary: {
      paymentCount: payments.length,
      callbackCount: callbackLogs.length,
      pendingCallbacks,
      totalSuccess,
      refundCount: refunds.length,
      totalRefunded,
      pendingRefunds,
    },
  };
}

export async function getPaymentReconciliationData() {
  const now = new Date();
  const recentStart = startOfDay(subDays(now, 6));

  const [payments, refunds, callbackLogs, invoices, customers] = await Promise.all([
    db.payment.findMany({
      include: {
        customer: {
          select: {
            id: true,
            name: true,
            customerNo: true,
          },
        },
        invoice: {
          select: {
            id: true,
            invoiceNo: true,
            status: true,
            totalAmount: true,
            paidAmount: true,
            dueDate: true,
          },
        },
        refunds: {
          select: {
            id: true,
            amount: true,
            status: true,
            method: true,
            processedAt: true,
            createdAt: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.refundRecord.findMany({
      include: {
        customer: {
          select: {
            id: true,
            name: true,
            customerNo: true,
          },
        },
        invoice: {
          select: {
            id: true,
            invoiceNo: true,
          },
        },
        payment: {
          select: {
            id: true,
            paymentNo: true,
            amount: true,
            method: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.paymentCallbackLog.findMany({
      include: {
        payment: {
          select: {
            id: true,
            paymentNo: true,
            amount: true,
            status: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 100,
    }),
    db.invoice.findMany({
      include: {
        customer: {
          select: {
            id: true,
            name: true,
            customerNo: true,
          },
        },
      },
      orderBy: {
        dueDate: "asc",
      },
    }),
    db.customer.findMany({
      select: {
        id: true,
        name: true,
        customerNo: true,
        creditBalance: true,
        status: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
  ]);

  const effectivePayments = payments.filter((payment) =>
    ["SUCCESS", "REFUNDED"].includes(payment.status),
  );
  const successfulRefunds = refunds.filter((refund) => refund.status === "SUCCESS");

  const grossReceipts = sumMoney(effectivePayments.map((payment) => payment.amount));
  const totalFees = sumMoney(effectivePayments.map((payment) => payment.feeAmount));
  const totalRefunded = sumMoney(successfulRefunds.map((refund) => refund.amount));
  const netReceipts = grossReceipts - totalRefunded - totalFees;

  const openInvoices = invoices
    .filter((invoice) => !["PAID", "VOID", "DRAFT"].includes(invoice.status))
    .map((invoice) => ({
      ...invoice,
      outstandingAmount: Math.max(invoice.totalAmount - invoice.paidAmount, 0),
    }));

  const outstandingAmount = sumMoney(
    openInvoices.map((invoice) => invoice.outstandingAmount),
  );
  const overdueInvoices = openInvoices.filter(
    (invoice) => invoice.status === "OVERDUE" || invoice.dueDate < now,
  );
  const pendingCallbacks = callbackLogs.filter((log) => !log.isHandled);
  const failedCallbacks = callbackLogs.filter(
    (log) => log.callbackStatus === "FAILED",
  );
  const pendingRefunds = refunds.filter((refund) => refund.status === "PENDING");
  const negativeBalanceCustomers = customers
    .filter((customer) => customer.creditBalance < 0)
    .sort((a, b) => a.creditBalance - b.creditBalance)
    .slice(0, 10);

  const issueRows = [
    ...invoices.flatMap((invoice) => {
      const outstanding = Math.max(invoice.totalAmount - invoice.paidAmount, 0);

      if (invoice.paidAmount > invoice.totalAmount) {
        return [
          {
            kind: "INVOICE",
            level: "critical",
            title: "账单实收金额超过应收金额",
            detail: `账单 ${invoice.invoiceNo} 已收金额高于账单总额，请核查收款或退款记录。`,
            relatedNo: invoice.invoiceNo,
            customerName: invoice.customer.name,
            amount: invoice.paidAmount - invoice.totalAmount,
            createdAt: invoice.updatedAt,
          },
        ];
      }

      if (invoice.status === "PAID" && outstanding > 0) {
        return [
          {
            kind: "INVOICE",
            level: "critical",
            title: "账单状态为已支付但仍有未收金额",
            detail: `账单 ${invoice.invoiceNo} 状态为已支付，但仍有 ${(outstanding / 100).toFixed(
              2,
            )} 元未核销。`,
            relatedNo: invoice.invoiceNo,
            customerName: invoice.customer.name,
            amount: outstanding,
            createdAt: invoice.updatedAt,
          },
        ];
      }

      if (!["PAID", "VOID", "DRAFT"].includes(invoice.status) && outstanding === 0) {
        return [
          {
            kind: "INVOICE",
            level: "warning",
            title: "账单金额已结清但状态未关闭",
            detail: `账单 ${invoice.invoiceNo} 已无未收金额，但状态仍为 ${invoice.status}。`,
            relatedNo: invoice.invoiceNo,
            customerName: invoice.customer.name,
            amount: 0,
            createdAt: invoice.updatedAt,
          },
        ];
      }

      return [];
    }),
    ...payments.flatMap((payment) => {
      const refundedAmount = payment.refunds
        .filter((refund) => refund.status === "SUCCESS")
        .reduce((total, refund) => total + refund.amount, 0);

      const issues = [];

      if (
        payment.status !== "FAILED" &&
        payment.method !== "BALANCE" &&
        !payment.transactionNo
      ) {
        issues.push({
          kind: "PAYMENT",
          level: "warning",
          title: "外部渠道收款缺少流水号",
          detail: `收款单 ${payment.paymentNo} 使用 ${payment.method} 入账，但未记录渠道流水号。`,
          relatedNo: payment.paymentNo,
          customerName: payment.customer.name,
          amount: payment.amount,
          createdAt: payment.createdAt,
        });
      }

      if (payment.status !== "FAILED" && !payment.invoiceId) {
        issues.push({
          kind: "PAYMENT",
          level: "warning",
          title: "收款记录未关联账单",
          detail: `收款单 ${payment.paymentNo} 没有关联账单，建议尽快核销。`,
          relatedNo: payment.paymentNo,
          customerName: payment.customer.name,
          amount: payment.amount,
          createdAt: payment.createdAt,
        });
      }

      if (refundedAmount > payment.amount) {
        issues.push({
          kind: "REFUND",
          level: "critical",
          title: "退款金额超过原始收款",
          detail: `收款单 ${payment.paymentNo} 的退款总额已超过原始收款金额。`,
          relatedNo: payment.paymentNo,
          customerName: payment.customer.name,
          amount: refundedAmount - payment.amount,
          createdAt: payment.updatedAt,
        });
      }

      return issues;
    }),
    ...callbackLogs.flatMap((log) => {
      if (!log.isHandled) {
        return [
          {
            kind: "CALLBACK",
            level: "warning",
            title: "存在待处理支付回调",
            detail: log.message || `交易流水 ${log.transactionNo || "-"} 尚未完成入账处理。`,
            relatedNo: log.transactionNo || log.invoiceNo || log.paymentNo || log.id,
            customerName: "-",
            amount: log.payment?.amount ?? 0,
            createdAt: log.createdAt,
          },
        ];
      }

      if (log.callbackStatus === "FAILED") {
        return [
          {
            kind: "CALLBACK",
            level: "critical",
            title: "支付回调处理失败",
            detail: log.message || `交易流水 ${log.transactionNo || "-"} 回调失败。`,
            relatedNo: log.transactionNo || log.invoiceNo || log.paymentNo || log.id,
            customerName: "-",
            amount: log.payment?.amount ?? 0,
            createdAt: log.createdAt,
          },
        ];
      }

      if (log.callbackStatus === "SUCCESS" && !log.paymentId) {
        return [
          {
            kind: "CALLBACK",
            level: "warning",
            title: "成功回调未生成收款单",
            detail: `回调 ${log.transactionNo || log.invoiceNo || log.id} 显示成功，但未绑定收款记录。`,
            relatedNo: log.transactionNo || log.invoiceNo || log.id,
            customerName: "-",
            amount: 0,
            createdAt: log.createdAt,
          },
        ];
      }

      return [];
    }),
    ...negativeBalanceCustomers.map((customer) => ({
      kind: "BALANCE",
      level: "warning" as const,
      title: "客户余额为负数",
      detail: `客户 ${customer.name} 当前余额为负，请核查授信、退款或自动续费扣款。`,
      relatedNo: customer.customerNo,
      customerName: customer.name,
      amount: Math.abs(customer.creditBalance),
      createdAt: now,
    })),
  ]
    .sort(
      (a, b) =>
        new Date(b.createdAt ?? 0).getTime() - new Date(a.createdAt ?? 0).getTime(),
    )
    .slice(0, 80);

  const methodPool = new Set<string>();
  for (const payment of payments) {
    methodPool.add(payment.method);
  }
  for (const refund of refunds) {
    methodPool.add(refund.method);
  }
  for (const log of callbackLogs) {
    if (log.method) {
      methodPool.add(log.method);
    }
  }

  const methodSummary = Array.from(methodPool)
    .map((method) => {
      const methodPayments = effectivePayments.filter((payment) => payment.method === method);
      const methodRefunds = successfulRefunds.filter((refund) => refund.method === method);
      const methodCallbacks = callbackLogs.filter((log) => log.method === method);

      return {
        method,
        paymentCount: methodPayments.length,
        paymentAmount: sumMoney(methodPayments.map((payment) => payment.amount)),
        feeAmount: sumMoney(methodPayments.map((payment) => payment.feeAmount)),
        refundCount: methodRefunds.length,
        refundAmount: sumMoney(methodRefunds.map((refund) => refund.amount)),
        netAmount:
          sumMoney(methodPayments.map((payment) => payment.amount)) -
          sumMoney(methodPayments.map((payment) => payment.feeAmount)) -
          sumMoney(methodRefunds.map((refund) => refund.amount)),
        pendingCallbacks: methodCallbacks.filter((log) => !log.isHandled).length,
        failedCallbacks: methodCallbacks.filter((log) => log.callbackStatus === "FAILED")
          .length,
        lastPaymentAt: methodPayments[0]?.paidAt ?? methodPayments[0]?.createdAt ?? null,
      };
    })
    .sort((a, b) => b.paymentAmount - a.paymentAmount);

  const trend = eachDayOfInterval({
    start: recentStart,
    end: startOfDay(now),
  }).map((day) => {
    const dayStart = startOfDay(day);
    const dayEnd = endOfDay(day);

    const dayPayments = effectivePayments.filter((payment) => {
      const target = payment.paidAt ?? payment.createdAt;
      return target >= dayStart && target <= dayEnd;
    });
    const dayRefunds = successfulRefunds.filter((refund) => {
      const target = refund.processedAt ?? refund.createdAt;
      return target >= dayStart && target <= dayEnd;
    });

    return {
      date: format(day, "MM-dd"),
      paymentAmount: sumMoney(dayPayments.map((payment) => payment.amount)),
      refundAmount: sumMoney(dayRefunds.map((refund) => refund.amount)),
      netAmount:
        sumMoney(dayPayments.map((payment) => payment.amount)) -
        sumMoney(dayRefunds.map((refund) => refund.amount)),
      paymentCount: dayPayments.length,
      refundCount: dayRefunds.length,
    };
  });

  return {
    summary: {
      grossReceipts,
      totalFees,
      totalRefunded,
      netReceipts,
      outstandingAmount,
      openInvoiceCount: openInvoices.length,
      overdueInvoiceCount: overdueInvoices.length,
      pendingCallbackCount: pendingCallbacks.length,
      failedCallbackCount: failedCallbacks.length,
      pendingRefundCount: pendingRefunds.length,
      issueCount: issueRows.length,
      negativeBalanceCustomerCount: negativeBalanceCustomers.length,
    },
    methodSummary,
    issues: issueRows,
    openInvoices: openInvoices.slice(0, 20),
    trend,
    negativeBalanceCustomers,
  };
}

export async function getNotificationsPageData() {
  const [templates, messages, tasks] = await Promise.all([
    db.notificationTemplate.findMany({
      orderBy: [
        {
          isActive: "desc",
        },
        {
          createdAt: "desc",
        },
      ],
    }),
    db.notificationMessage.findMany({
      include: {
        customer: {
          select: {
            id: true,
            name: true,
            email: true,
          },
        },
        template: true,
        createdBy: {
          select: {
            id: true,
            name: true,
            email: true,
            role: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 50,
    }),
    db.asyncTaskJob.findMany({
      where: {
        queueName: "notifications",
      },
      include: {
        notification: {
          select: {
            id: true,
            recipient: true,
            channel: true,
            status: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 50,
    }),
  ]);

  return {
    templates,
    messages,
    tasks,
    summary: {
      templateCount: templates.length,
      totalMessages: messages.length,
      pendingMessages: messages.filter((item) => item.status === "PENDING").length,
      sentMessages: messages.filter((item) => item.status === "SENT").length,
      failedMessages: messages.filter((item) => item.status === "FAILED").length,
      pendingTasks: tasks.filter((item) => item.status === "PENDING").length,
    },
  };
}

export async function getTicketsPageData() {
  const [tickets, customers, services, admins] = await Promise.all([
    db.ticket.findMany({
      include: {
        customer: true,
        service: true,
        assignedTo: true,
        replies: {
          orderBy: {
            createdAt: "asc",
          },
        },
      },
      orderBy: [
        {
          priority: "desc",
        },
        {
          createdAt: "desc",
        },
      ],
    }),
    db.customer.findMany({
      orderBy: {
        name: "asc",
      },
    }),
    db.serviceInstance.findMany({
      include: {
        customer: {
          select: {
            name: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.adminUser.findMany({
      where: {
        isActive: true,
      },
      orderBy: {
        name: "asc",
      },
    }),
  ]);

  return {
    tickets,
    customers,
    services,
    admins,
  };
}

export async function getBillingPageData() {
  const [settings, jobs, dueServices, openInvoices, taxProfiles, paymentGateways] = await Promise.all([
    getBillingSetting(),
    db.billingJob.findMany({
      include: {
        customer: true,
        invoice: true,
        service: true,
      },
      orderBy: {
        executedAt: "desc",
      },
      take: 20,
    }),
    db.serviceInstance.findMany({
      where: {
        nextDueDate: {
          not: null,
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
      take: 12,
    }),
    db.invoice.findMany({
      where: {
        status: {
          in: ["ISSUED", "PARTIAL", "OVERDUE"],
        },
      },
      include: {
        customer: true,
        service: true,
      },
      orderBy: {
        dueDate: "asc",
      },
      take: 20,
    }),
    db.invoiceTaxProfile.findMany({
      orderBy: [
        {
          isDefault: "desc",
        },
        {
          taxRate: "asc",
        },
      ],
    }),
    db.paymentGatewayConfig.findMany({
      orderBy: [
        {
          isEnabled: "desc",
        },
        {
          method: "asc",
        },
      ],
    }),
  ]);

  return {
    settings,
    jobs,
    dueServices,
    openInvoices,
    taxProfiles,
    paymentGateways,
  };
}

export async function getResourcesPageData() {
  const [vpcs, ips, disks, snapshots, backups, securityGroups] = await Promise.all([
    db.serviceVpcNetwork.findMany({
      include: {
        services: {
          select: {
            id: true,
            name: true,
            serviceNo: true,
          },
        },
        securityGroups: {
          select: {
            id: true,
            name: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.serviceIpAddress.findMany({
      include: {
        service: {
          select: {
            id: true,
            name: true,
            serviceNo: true,
            customer: {
              select: {
                name: true,
              },
            },
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.serviceDisk.findMany({
      include: {
        service: {
          select: {
            id: true,
            name: true,
            serviceNo: true,
            customer: {
              select: {
                name: true,
              },
            },
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.serviceSnapshot.findMany({
      include: {
        service: {
          select: {
            id: true,
            name: true,
            serviceNo: true,
          },
        },
        sourceDisk: {
          select: {
            name: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.serviceBackup.findMany({
      include: {
        service: {
          select: {
            id: true,
            name: true,
            serviceNo: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.serviceSecurityGroup.findMany({
      include: {
        service: {
          select: {
            id: true,
            name: true,
            serviceNo: true,
          },
        },
        vpcNetwork: {
          select: {
            id: true,
            name: true,
          },
        },
        rules: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
  ]);

  return {
    vpcs,
    ips,
    disks,
    snapshots,
    backups,
    securityGroups,
  };
}

export async function getAuditPageData() {
  return db.auditLog.findMany({
    include: {
      adminUser: {
        select: {
          id: true,
          name: true,
          email: true,
          role: true,
        },
      },
    },
    orderBy: {
      createdAt: "desc",
    },
    take: 100,
  });
}

export async function getCustomerLedgerData(customerId: string) {
  const customer = await db.customer.findUnique({
    where: {
      id: customerId,
    },
    include: {
      creditTransactions: {
        include: {
          operator: {
            select: {
              id: true,
              name: true,
              email: true,
              role: true,
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
        take: 100,
      },
    },
  });

  return customer;
}

export async function getArchitecturePageData() {
  const [regions, zones, flavors, images, plans, products, customers, portalUsers] =
    await Promise.all([
      db.cloudRegion.findMany({
        include: {
          _count: {
            select: {
              zones: true,
              plans: true,
              images: true,
            },
          },
        },
        orderBy: [
          {
            sortOrder: "asc",
          },
          {
            createdAt: "desc",
          },
        ],
      }),
      db.cloudZone.findMany({
        include: {
          region: true,
          _count: {
            select: {
              plans: true,
            },
          },
        },
        orderBy: [
          {
            sortOrder: "asc",
          },
          {
            createdAt: "desc",
          },
        ],
      }),
      db.cloudFlavor.findMany({
        include: {
          _count: {
            select: {
              plans: true,
            },
          },
        },
        orderBy: [
          {
            sortOrder: "asc",
          },
          {
            createdAt: "desc",
          },
        ],
      }),
      db.cloudImage.findMany({
        include: {
          region: true,
          _count: {
            select: {
              plans: true,
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
      }),
      db.cloudPlan.findMany({
        include: {
          product: true,
          region: true,
          zone: true,
          flavor: true,
          image: true,
          _count: {
            select: {
              services: true,
              orderItems: true,
            },
          },
        },
        orderBy: {
          createdAt: "desc",
        },
      }),
      db.product.findMany({
        where: {
          status: {
            in: ["ACTIVE", "DRAFT"],
          },
        },
        orderBy: {
          name: "asc",
        },
      }),
      db.customer.findMany({
        where: {
          status: {
            in: ["ACTIVE", "OVERDUE", "SUSPENDED"],
          },
        },
        orderBy: {
          name: "asc",
        },
      }),
      db.customerUser.findMany({
        include: {
          customer: true,
        },
        orderBy: {
          createdAt: "desc",
        },
      }),
    ]);

  return {
    regions,
    zones,
    flavors,
    images,
    plans,
    products,
    customers,
    portalUsers,
  };
}

