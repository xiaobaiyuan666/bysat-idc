import { db } from "@/lib/db";

function sumMoney(values: number[]) {
  return values.reduce((total, value) => total + value, 0);
}

function hasOpenInvoice(status: string) {
  return ["DRAFT", "ISSUED", "PARTIAL", "OVERDUE"].includes(status);
}

export async function getPortalDashboardData(customerId: string) {
  const [customer, services, invoices, tickets, plans] = await Promise.all([
    db.customer.findUniqueOrThrow({
      where: {
        id: customerId,
      },
    }),
    db.serviceInstance.findMany({
      where: {
        customerId,
      },
      include: {
        product: true,
        plan: {
          include: {
            region: true,
            zone: true,
            flavor: true,
          },
        },
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 6,
    }),
    db.invoice.findMany({
      where: {
        customerId,
      },
      include: {
        service: true,
      },
      orderBy: {
        dueDate: "asc",
      },
      take: 8,
    }),
    db.ticket.findMany({
      where: {
        customerId,
      },
      include: {
        service: true,
        replies: {
          orderBy: {
            createdAt: "desc",
          },
          take: 1,
        },
      },
      orderBy: {
        updatedAt: "desc",
      },
      take: 5,
    }),
    db.cloudPlan.findMany({
      where: {
        isPublic: true,
        isActive: true,
      },
      include: {
        product: true,
        region: true,
        zone: true,
        flavor: true,
      },
      orderBy: [
        {
          salePrice: "asc",
        },
        {
          createdAt: "desc",
        },
      ],
      take: 6,
    }),
  ]);

  const openInvoices = invoices.filter((invoice) =>
    ["ISSUED", "PARTIAL", "OVERDUE"].includes(invoice.status),
  );

  return {
    customer,
    metrics: {
      creditBalance: customer.creditBalance,
      activeServices: services.filter((item) => item.status === "ACTIVE").length,
      openTickets: tickets.filter((item) => item.status !== "CLOSED").length,
      openInvoiceAmount: sumMoney(
        openInvoices.map((invoice) => Math.max(invoice.totalAmount - invoice.paidAmount, 0)),
      ),
    },
    services,
    invoices,
    tickets,
    plans,
  };
}

export async function getPortalCatalogData() {
  return db.cloudPlan.findMany({
    where: {
      isPublic: true,
      isActive: true,
      product: {
        status: "ACTIVE",
      },
    },
    include: {
      product: true,
      region: true,
      zone: true,
      flavor: true,
      image: true,
    },
    orderBy: [
      {
        salePrice: "asc",
      },
      {
        createdAt: "desc",
      },
    ],
  });
}

export async function getPortalOrdersData(customerId: string) {
  return db.order.findMany({
    where: {
      customerId,
    },
    include: {
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
        },
      },
      invoices: true,
    },
    orderBy: {
      createdAt: "desc",
    },
  });
}

export async function getPortalServicesData(customerId: string) {
  return db.serviceInstance.findMany({
    where: {
      customerId,
    },
    include: {
      product: true,
      plan: {
        include: {
          region: true,
          zone: true,
          flavor: true,
          image: true,
        },
      },
      vpcNetwork: true,
      ipAddresses: true,
      disks: true,
      snapshots: true,
      backups: true,
      tickets: {
        select: {
          id: true,
        },
      },
      securityGroups: {
        include: {
          rules: true,
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

export async function getPortalServiceDetail(customerId: string, serviceId: string) {
  const service = await db.serviceInstance.findFirst({
    where: {
      id: serviceId,
      customerId,
    },
    include: {
      customer: true,
      product: true,
      order: true,
      plan: {
        include: {
          region: true,
          zone: true,
          flavor: true,
          image: true,
        },
      },
      vpcNetwork: true,
      ipAddresses: true,
      disks: {
        orderBy: {
          createdAt: "asc",
        },
      },
      snapshots: {
        include: {
          sourceDisk: true,
        },
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
              createdAt: "asc",
            },
          },
        },
      },
      invoices: {
        orderBy: {
          createdAt: "desc",
        },
        take: 8,
      },
      tickets: {
        orderBy: {
          updatedAt: "desc",
        },
        take: 6,
      },
    },
  });

  if (!service) {
    return null;
  }

  const auditTargets = [
    {
      targetType: "service",
      targetId: service.id,
    },
    ...service.disks.map((disk) => ({
      targetType: "serviceDisk",
      targetId: disk.id,
    })),
    ...service.snapshots.map((snapshot) => ({
      targetType: "serviceSnapshot",
      targetId: snapshot.id,
    })),
    ...service.backups.map((backup) => ({
      targetType: "serviceBackup",
      targetId: backup.id,
    })),
    ...service.securityGroups.map((group) => ({
      targetType: "serviceSecurityGroup",
      targetId: group.id,
    })),
    ...service.securityGroups.flatMap((group) =>
      group.rules.map((rule) => ({
        targetType: "serviceSecurityGroupRule",
        targetId: rule.id,
      })),
    ),
  ];

  const providerResourceIds = [
    service.providerResourceId,
    ...service.disks.map((disk) => disk.providerDiskId),
    ...service.snapshots.map((snapshot) => snapshot.providerSnapshotId),
    ...service.backups.map((backup) => backup.providerBackupId),
    ...service.securityGroups.map((group) => group.providerSecurityGroupId),
  ].filter((value): value is string => Boolean(value));

  const [auditLogs, providerLogs, renewalInvoice] = await Promise.all([
    db.auditLog.findMany({
      where: {
        OR: auditTargets,
      },
      orderBy: {
        createdAt: "desc",
      },
      take: 20,
    }),
    providerResourceIds.length > 0
      ? db.providerSyncLog.findMany({
          where: {
            resourceId: {
              in: providerResourceIds,
            },
          },
          orderBy: {
            syncedAt: "desc",
          },
          take: 20,
        })
      : [],
    db.invoice.findFirst({
      where: {
        customerId,
        serviceId: service.id,
        type: "RENEWAL",
        status: {
          in: ["ISSUED", "PARTIAL", "OVERDUE"],
        },
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
  ]);

  return {
    service,
    auditLogs,
    providerLogs,
    renewalInvoice,
  };
}

export async function getPortalInvoicesData(customerId: string) {
  const [customer, invoices, payments, paymentGateways, services] = await Promise.all([
    db.customer.findUniqueOrThrow({
      where: {
        id: customerId,
      },
      select: {
        id: true,
        name: true,
        customerNo: true,
        creditBalance: true,
      },
    }),
    db.invoice.findMany({
      where: {
        customerId,
      },
      include: {
        service: true,
        order: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.payment.findMany({
      where: {
        customerId,
      },
      include: {
        invoice: true,
        order: true,
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
    db.paymentGatewayConfig.findMany({
      where: {
        isEnabled: true,
      },
      orderBy: {
        method: "asc",
      },
    }),
    db.serviceInstance.findMany({
      where: {
        customerId,
        status: {
          in: ["ACTIVE", "SUSPENDED", "OVERDUE", "EXPIRED"],
        },
        nextDueDate: {
          not: null,
        },
      },
      include: {
        product: true,
        order: true,
      },
      orderBy: [
        {
          nextDueDate: "asc",
        },
        {
          createdAt: "desc",
        },
      ],
      take: 20,
    }),
  ]);

  const serviceIdsWithOpenInvoice = new Set(
    invoices
      .filter((invoice) => invoice.serviceId && hasOpenInvoice(invoice.status))
      .map((invoice) => invoice.serviceId as string),
  );
  const renewalPreviews = services
    .filter((service) => !serviceIdsWithOpenInvoice.has(service.id))
    .map((service) => ({
      id: service.id,
      serviceId: service.id,
      serviceNo: service.serviceNo,
      serviceName: service.name,
      productName: service.product.name,
      dueDate: service.nextDueDate,
      billingCycle: service.billingCycle,
      amount: service.monthlyCost,
      status: service.status,
      orderNo: service.order?.orderNo,
    }));

  return {
    customer,
    invoices,
    payments,
    paymentGateways,
    renewalPreviews,
  };
}

export async function getPortalWalletData(customerId: string) {
  const customer = await db.customer.findUniqueOrThrow({
    where: {
      id: customerId,
    },
    include: {
      creditTransactions: {
        orderBy: {
          createdAt: "desc",
        },
        take: 20,
      },
    },
  });

  return customer;
}

export async function getPortalTicketsData(customerId: string) {
  const [tickets, services] = await Promise.all([
    db.ticket.findMany({
      where: {
        customerId,
      },
      include: {
        service: true,
        assignedTo: true,
        replies: {
          orderBy: {
            createdAt: "asc",
          },
        },
      },
      orderBy: {
        updatedAt: "desc",
      },
    }),
    db.serviceInstance.findMany({
      where: {
        customerId,
      },
      orderBy: {
        createdAt: "desc",
      },
    }),
  ]);

  return {
    tickets,
    services,
  };
}

export async function getPortalNotificationsData(customerId: string) {
  return db.notificationMessage.findMany({
    where: {
      customerId,
      status: {
        in: ["PENDING", "SENT", "FAILED"],
      },
    },
    include: {
      template: true,
    },
    orderBy: {
      createdAt: "desc",
    },
    take: 30,
  });
}
