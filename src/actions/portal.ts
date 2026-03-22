"use server";

import { type PaymentMethod } from "@prisma/client";
import { revalidatePath } from "next/cache";
import { headers } from "next/headers";
import { redirect } from "next/navigation";

import { writeAuditLog } from "@/lib/audit";
import { createServiceRenewalInvoice } from "@/lib/billing-engine";
import { db } from "@/lib/db";
import {
  computeGatewaySignature,
  getGatewayCallbackHeaderName,
  getPaymentGatewayConfig,
} from "@/lib/finance-config";
import { addCycle } from "@/lib/format";
import { parseCloudPlanConfig } from "@/lib/cloud-plan-config";
import { queueNotification } from "@/lib/notification-service";
import { settleInvoicePayment } from "@/lib/payment-service";
import { requirePortalUser } from "@/lib/portal-auth";
import { type ProviderActionPayload } from "@/lib/providers/types";
import {
  executeResourceAction,
  getResourceOwnershipContext,
  isSupportedResourceAction as isSupportedTypedResourceAction,
  isSupportedResourceType,
  ResourceActionError,
} from "@/lib/resource-operations";
import {
  executeServiceAction,
  isSupportedServiceAction,
} from "@/lib/service-operations";
import {
  portalOrderSchema,
  portalTicketCreateSchema,
  portalTicketReplySchema,
} from "@/lib/validation";
import { makeCode } from "@/lib/utils";
import { buildVncConsoleUrl } from "@/lib/vnc";

function normalizeText(value: FormDataEntryValue | null) {
  if (typeof value !== "string") {
    return undefined;
  }

  const normalized = value.trim();
  return normalized ? normalized : undefined;
}

const portalServiceActionLabels = {
  sync: "同步状态",
  activate: "开通实例",
  suspend: "暂停实例",
  renew: "续费实例",
  terminate: "删除实例",
  powerOn: "开机",
  powerOff: "关机",
  reboot: "重启",
  hardReboot: "强制重启",
  hardPowerOff: "强制关机",
  reinstall: "重装系统",
  resetPassword: "重置密码",
  getVnc: "获取 VNC",
  rescueStart: "进入救援模式",
  rescueStop: "退出救援模式",
  lock: "锁定实例",
  unlock: "解除锁定",
} as const;

const portalResourceActionLabels = {
  disks: {
    createSnapshot: "创建快照",
    createBackup: "创建备份",
    attach: "挂载磁盘",
    detach: "卸载磁盘",
    setBoot: "设为系统盘",
  },
  snapshots: {
    restore: "恢复快照",
    deleteSnapshot: "删除快照",
  },
  backups: {
    restore: "恢复备份",
    expireNow: "归档备份",
    deleteBackup: "删除备份",
  },
  "security-groups": {
    addRule: "新增安全组规则",
    deleteRule: "删除安全组规则",
    deleteGroup: "删除安全组",
  },
} as const;

function getPortalServiceActionLabel(action: string) {
  if (action === "unsuspend") {
    return "解除暂停";
  }

  return portalServiceActionLabels[action as keyof typeof portalServiceActionLabels] ?? action;
}

function getPortalResourceActionLabel(resourceType: string, action: string) {
  const labels =
    portalResourceActionLabels[
      resourceType as keyof typeof portalResourceActionLabels
    ] as Record<string, string> | undefined;

  return labels?.[action] ?? action;
}

function buildServiceResultSearch(input: {
  action: string;
  success: boolean;
  label?: string;
  message?: string;
  taskId?: string;
  consoleUrl?: string;
  invoiceNo?: string;
}) {
  const params = new URLSearchParams({
    action: input.action,
    status: input.success ? "success" : "error",
  });

  if (input.message) {
    params.set("message", input.message);
  }

  if (input.label) {
    params.set("label", input.label);
  }

  if (input.taskId) {
    params.set("taskId", input.taskId);
  }

  if (input.consoleUrl) {
    params.set("consoleUrl", input.consoleUrl);
  }

  if (input.invoiceNo) {
    params.set("invoiceNo", input.invoiceNo);
  }

  return params.toString();
}

export async function portalCreateOrderAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const parsed = portalOrderSchema.safeParse({
    planId: formData.get("planId"),
    quantity: formData.get("quantity"),
    serviceName: formData.get("serviceName"),
    notes: formData.get("notes"),
  });

  if (!parsed.success) {
    redirect("/portal/store?error=order");
  }

  const plan = await db.cloudPlan.findUnique({
    where: {
      id: parsed.data.planId,
    },
    include: {
      product: true,
      region: true,
      zone: true,
      flavor: true,
      image: true,
    },
  });

  if (!plan || !plan.isActive || !plan.isPublic) {
    redirect("/portal/store?error=plan");
  }

  const now = new Date();
  const payableDueDate = new Date(now.getTime() + 1000 * 60 * 60 * 24);
  const nextDueDate = addCycle(now, plan.billingCycle);
  const totalAmount = plan.salePrice * parsed.data.quantity + plan.setupFee;
  const planConfig = parseCloudPlanConfig(plan.configOptions);
  let createdOrderNo = "";
  let createdServiceId = "";
  let createdServiceNo = "";
  let createdInvoiceNo = "";

  await db.$transaction(async (tx) => {
    const order = await tx.order.create({
      data: {
        orderNo: makeCode("ORD"),
        customerId: portalUser.customerId,
        status: "PENDING",
        totalAmount,
        source: "portal",
        orderType: "new",
        dueDate: payableDueDate,
        notes: normalizeText(formData.get("notes")),
      },
    });
    createdOrderNo = order.orderNo;

    const service = await tx.serviceInstance.create({
      data: {
        serviceNo: makeCode("SRV"),
        customerId: portalUser.customerId,
        productId: plan.productId,
        planId: plan.id,
        orderId: order.id,
        name: parsed.data.serviceName,
        hostname: `${plan.code.toLowerCase()}.${portalUser.customer.customerNo.toLowerCase()}.local`,
        providerType: plan.product.providerType,
        region: plan.region.code,
        billingCycle: plan.billingCycle,
        status: "PENDING",
        cpuCores: Number(planConfig.cpu ?? plan.flavor?.cpu ?? 0) || undefined,
        memoryGb: Number(planConfig.memory ?? plan.flavor?.memoryGb ?? 0) || undefined,
        storageGb:
          Number(planConfig.system_disk_size ?? plan.flavor?.storageGb ?? 0) || undefined,
        monthlyCost: plan.salePrice * parsed.data.quantity,
        nextDueDate,
        configSnapshot: JSON.stringify({
          ...planConfig,
          quantity: parsed.data.quantity,
          planCode: plan.code,
          region: plan.region.code,
          zone: plan.zone?.code,
          flavor: plan.flavor?.code,
          image: plan.image?.code,
          source: "portal",
        }),
      },
    });
    createdServiceId = service.id;
    createdServiceNo = service.serviceNo;

    await tx.orderItem.create({
      data: {
        orderId: order.id,
        productId: plan.productId,
        planId: plan.id,
        serviceId: service.id,
        title: `${plan.product.name} / ${plan.name}`,
        quantity: parsed.data.quantity,
        unitPrice: plan.salePrice,
        cycle: plan.billingCycle,
        totalAmount,
      },
    });

    const invoice = await tx.invoice.create({
      data: {
        invoiceNo: makeCode("INV"),
        customerId: portalUser.customerId,
        orderId: order.id,
        serviceId: service.id,
        type: "ORDER",
        status: "ISSUED",
        subtotal: totalAmount,
        totalAmount,
        dueDate: payableDueDate,
        issuedAt: now,
        remark: "门户下单自动生成的待支付账单",
      },
    });
    createdInvoiceNo = invoice.invoiceNo;
  });

  if (createdOrderNo && createdServiceNo) {
    await queueNotification({
      templateCode: "ORDER_CREATED",
      customerId: portalUser.customerId,
      recipient: portalUser.email,
      recipientName: portalUser.name,
      module: "order",
      relatedType: "service",
      relatedId: createdServiceId,
      variables: {
        customer_name: portalUser.customer.name,
        order_no: createdOrderNo,
        service_no: createdServiceNo,
        invoice_no: createdInvoiceNo,
      },
      content: `订单 ${createdOrderNo} 已创建，服务 ${createdServiceNo} 和待支付账单已生成。`,
    });
  }

  revalidatePath("/portal");
  revalidatePath("/portal/store");
  revalidatePath("/portal/orders");
  revalidatePath("/portal/invoices");
  redirect("/portal/orders?created=1");
}

export async function portalPayInvoiceAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const invoiceId = String(formData.get("invoiceId") ?? "");
  const method = String(formData.get("method") ?? "BALANCE");

  if (!invoiceId) {
    redirect("/portal/invoices?error=invoice");
  }

  const invoice = await db.invoice.findFirst({
    where: {
      id: invoiceId,
      customerId: portalUser.customerId,
    },
  });

  if (!invoice) {
    redirect("/portal/invoices?error=invoice");
  }

  const amount = Math.max(invoice.totalAmount - invoice.paidAmount, 0);

  if (amount <= 0) {
    redirect("/portal/invoices?error=paid");
  }

  try {
    if (method === "BALANCE") {
      await settleInvoicePayment({
        invoiceId,
        method: "BALANCE",
        amount,
      });
    } else {
      const gateway = await getPaymentGatewayConfig(method as PaymentMethod);

      if (!gateway || !gateway.isEnabled) {
        redirect("/portal/invoices?error=gateway");
      }

      const requestHeaders = await headers();
      const host =
        requestHeaders.get("x-forwarded-host") ??
        requestHeaders.get("host") ??
        "127.0.0.1:3000";
      const protocol =
        requestHeaders.get("x-forwarded-proto") ??
        (host.includes("localhost") || host.includes("127.0.0.1") ? "http" : "https");
      const callbackUrl = `${protocol}://${host}/api/payments/callback`;
      const transactionNo = `PORTAL-${method}-${Date.now()}`;
      const payload = {
        method,
        invoiceNo: invoice.invoiceNo,
        transactionNo,
        amount: Number((amount / 100).toFixed(2)),
        status: "SUCCESS",
        message: `门户演示支付：${method}`,
      };
      const rawPayload = JSON.stringify(payload);
      const signature = computeGatewaySignature(gateway, rawPayload);
      const signatureHeader = getGatewayCallbackHeaderName(gateway);
      const response = await fetch(callbackUrl, {
        method: "POST",
        headers: {
          "content-type": "application/json",
          [signatureHeader]: signature,
        },
        body: rawPayload,
        cache: "no-store",
      });

      if (!response.ok) {
        redirect(`/portal/invoices?error=payment&method=${method}`);
      }
    }
  } catch (error) {
    if (error instanceof Error && error.message === "INSUFFICIENT_BALANCE") {
      redirect("/portal/invoices?error=balance");
    }

    redirect("/portal/invoices?error=payment");
  }

  revalidatePath("/portal");
  revalidatePath("/portal/invoices");
  revalidatePath("/portal/wallet");
  revalidatePath("/portal/services");
  redirect(`/portal/invoices?paid=1&method=${method}`);
}

export async function portalCreateTicketAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const parsed = portalTicketCreateSchema.safeParse({
    serviceId: formData.get("serviceId"),
    subject: formData.get("subject"),
    priority: formData.get("priority"),
    summary: formData.get("summary"),
  });

  if (!parsed.success) {
    redirect("/portal/tickets?error=create");
  }

  const serviceId = parsed.data.serviceId || undefined;

  if (serviceId) {
    const service = await db.serviceInstance.findFirst({
      where: {
        id: serviceId,
        customerId: portalUser.customerId,
      },
    });

    if (!service) {
      redirect("/portal/tickets?error=service");
    }
  }

  const ticket = await db.ticket.create({
    data: {
      ticketNo: makeCode("TIC"),
      customerId: portalUser.customerId,
      serviceId,
      subject: parsed.data.subject,
      priority: parsed.data.priority,
      status: "OPEN",
      summary: parsed.data.summary,
      lastReplyAt: new Date(),
    },
  });

  await db.ticketReply.create({
    data: {
      ticketId: ticket.id,
      authorType: "CUSTOMER",
      authorName: portalUser.name,
      content: parsed.data.summary,
    },
  });

  await queueNotification({
    templateCode: "TICKET_CREATED",
    customerId: portalUser.customerId,
    recipient: portalUser.email,
    recipientName: portalUser.name,
    module: "ticket",
    relatedType: "ticket",
    relatedId: ticket.id,
    variables: {
      customer_name: portalUser.customer.name,
      ticket_no: ticket.ticketNo,
      subject: ticket.subject,
    },
    content: `工单 ${ticket.ticketNo} 已提交，主题为“${ticket.subject}”。`,
  });

  revalidatePath("/portal");
  revalidatePath("/portal/tickets");
  redirect("/portal/tickets?created=1");
}

export async function portalReplyTicketAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const ticketId = String(formData.get("ticketId") ?? "");
  const parsed = portalTicketReplySchema.safeParse({
    content: formData.get("content"),
  });

  if (!ticketId || !parsed.success) {
    redirect("/portal/tickets?error=reply");
  }

  const ticket = await db.ticket.findFirst({
    where: {
      id: ticketId,
      customerId: portalUser.customerId,
    },
  });

  if (!ticket) {
    redirect("/portal/tickets?error=reply");
  }

  await db.$transaction(async (tx) => {
    await tx.ticketReply.create({
      data: {
        ticketId: ticket.id,
        authorType: "CUSTOMER",
        authorName: portalUser.name,
        content: parsed.data.content,
      },
    });

    await tx.ticket.update({
      where: {
        id: ticket.id,
      },
      data: {
        status: "OPEN",
        lastReplyAt: new Date(),
      },
    });
  });

  await queueNotification({
    customerId: portalUser.customerId,
    recipient: portalUser.email,
    recipientName: portalUser.name,
    channel: "SYSTEM",
    module: "ticket",
    relatedType: "ticket",
    relatedId: ticket.id,
    subject: `工单 ${ticket.ticketNo} 已追加回复`,
    content: `您已在工单 ${ticket.ticketNo} 中追加新的客户回复，客服将继续处理。`,
  });

  revalidatePath("/portal");
  revalidatePath("/portal/tickets");
  redirect("/portal/tickets?replied=1");
}

export async function portalManageServiceAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const serviceId = String(formData.get("serviceId") ?? "");
  const actionValue = String(formData.get("action") ?? "");

  if (!serviceId || !isSupportedServiceAction(actionValue)) {
    redirect("/portal/services?error=action");
  }

  const ownedService = await db.serviceInstance.findFirst({
    where: {
      id: serviceId,
      customerId: portalUser.customerId,
    },
    select: {
      id: true,
      serviceNo: true,
    },
  });

  if (!ownedService) {
    redirect("/portal/services?error=service");
  }

  const payload: ProviderActionPayload = {};
  const password = normalizeText(formData.get("password"));
  const imageId = normalizeText(formData.get("imageId"));

  if (password) {
    payload.password = password;
  }

  if (imageId) {
    payload.imageId = imageId;
  }

  const result = await executeServiceAction(ownedService.id, actionValue, payload);

  if (!result) {
    redirect("/portal/services?error=service");
  }

  const requestHeaders = await headers();
  const host =
    requestHeaders.get("x-forwarded-host") ??
    requestHeaders.get("host") ??
    "127.0.0.1:3000";
  const protocol =
    requestHeaders.get("x-forwarded-proto") ??
    (host.includes("localhost") || host.includes("127.0.0.1") ? "http" : "https");
  const origin = `${protocol}://${host}`;
  const vncConsoleUrl = await buildVncConsoleUrl(origin, result.providerResult.data);

  if (vncConsoleUrl) {
    result.providerResult.consoleUrl = vncConsoleUrl;
  }

  const actionLabel = getPortalServiceActionLabel(actionValue);
  const actionMessage =
    result.providerResult.message || `实例 ${ownedService.serviceNo} 已执行 ${actionLabel}。`;

  await writeAuditLog({
    module: "service",
    action: `portal-${actionValue}`,
    targetType: "service",
    targetId: ownedService.id,
    summary: `门户用户执行服务操作：${ownedService.serviceNo} / ${actionLabel}`,
    detail: {
      portalUserId: portalUser.id,
      portalUserName: portalUser.name,
      portalUserEmail: portalUser.email,
      providerType: result.originalService.providerType,
      remoteId:
        result.providerResult.remoteId ?? result.originalService.providerResourceId,
      taskId: result.providerResult.taskId,
      payload,
    },
  });

  await queueNotification({
    customerId: portalUser.customerId,
    recipient: portalUser.email,
    recipientName: portalUser.name,
    channel: "SYSTEM",
    module: "service",
    relatedType: "service",
    relatedId: ownedService.id,
    subject: `实例操作：${actionLabel}`,
    content: actionMessage,
  });

  revalidatePath("/portal");
  revalidatePath("/portal/services");
  revalidatePath(`/portal/services/${ownedService.id}`);

  redirect(
    `/portal/services/${ownedService.id}?${buildServiceResultSearch({
      action: actionValue,
      label: actionLabel,
      success: result.providerResult.ok,
      message: actionMessage,
      taskId: result.providerResult.taskId,
      consoleUrl: result.providerResult.consoleUrl,
    })}`,
  );
}

export async function portalManageResourceAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const resourceType = String(formData.get("resourceType") ?? "");
  const resourceId = String(formData.get("resourceId") ?? "");
  const action = String(formData.get("action") ?? "");

  if (!resourceId || !isSupportedResourceType(resourceType)) {
    redirect("/portal/services?error=resource");
  }

  if (!isSupportedTypedResourceAction(resourceType, action)) {
    redirect("/portal/services?error=resource-action");
  }

  const ownership = await getResourceOwnershipContext(resourceType, resourceId);

  if (!ownership || ownership.customerId !== portalUser.customerId || !ownership.serviceId) {
    redirect("/portal/services?error=resource-owner");
  }

  const payload = {
    name: normalizeText(formData.get("name")),
    mountPoint: normalizeText(formData.get("mountPoint")),
    expiresAt: normalizeText(formData.get("expiresAt")),
    direction: normalizeText(formData.get("direction")),
    protocol: normalizeText(formData.get("protocol")),
    portRange: normalizeText(formData.get("portRange")),
    sourceCidr: normalizeText(formData.get("sourceCidr")),
    ruleAction: normalizeText(formData.get("ruleAction")),
    description: normalizeText(formData.get("description")),
    ruleId: normalizeText(formData.get("ruleId")),
  } satisfies Record<string, string | undefined>;

  const cleanedPayload = Object.fromEntries(
    Object.entries(payload).filter(([, value]) => value !== undefined),
  );

  try {
    const result = await executeResourceAction(
      resourceType,
      resourceId,
      action,
      cleanedPayload,
    );

    const actionLabel = getPortalResourceActionLabel(resourceType, action);
    const actionMessage =
      result.summary || `资源 ${ownership.serviceNo ?? "-"} 已执行 ${actionLabel}。`;

    await writeAuditLog({
      module: "resource",
      action: `portal-${action}`,
      targetType: result.targetType,
      targetId: result.targetId,
      summary: `门户用户执行资源操作：${ownership.serviceNo ?? "-"} / ${actionLabel}`,
      detail: {
        portalUserId: portalUser.id,
        portalUserName: portalUser.name,
        portalUserEmail: portalUser.email,
        resourceType,
        resourceId,
        payload: cleanedPayload,
      },
    });

    await queueNotification({
      customerId: portalUser.customerId,
      recipient: portalUser.email,
      recipientName: portalUser.name,
      channel: "SYSTEM",
      module: "resource",
      relatedType: "service",
      relatedId: ownership.serviceId,
      subject: `资源操作：${actionLabel}`,
      content: actionMessage,
    });

    revalidatePath("/portal");
    revalidatePath("/portal/services");
    revalidatePath(`/portal/services/${ownership.serviceId}`);

    redirect(
      `/portal/services/${ownership.serviceId}?${buildServiceResultSearch({
        action,
        label: actionLabel,
        success: true,
        message: actionMessage,
      })}`,
    );
  } catch (error) {
    const message =
      error instanceof ResourceActionError
        ? error.message
        : error instanceof Error
          ? error.message
          : "资源操作失败";

    redirect(
      `/portal/services/${ownership.serviceId}?${buildServiceResultSearch({
        action,
        label: getPortalResourceActionLabel(resourceType, action),
        success: false,
        message,
      })}`,
    );
  }
}

export async function portalCreateRenewalInvoiceAction(formData: FormData) {
  const portalUser = await requirePortalUser();
  const serviceId = String(formData.get("serviceId") ?? "");

  if (!serviceId) {
    redirect("/portal/services?error=renew");
  }

  const ownedService = await db.serviceInstance.findFirst({
    where: {
      id: serviceId,
      customerId: portalUser.customerId,
    },
    select: {
      id: true,
      serviceNo: true,
    },
  });

  if (!ownedService) {
    redirect("/portal/services?error=service");
  }

  const created = await createServiceRenewalInvoice(ownedService.id, {
    source: "portal",
  });

  if (!created) {
    redirect(
      `/portal/services/${ownedService.id}?${buildServiceResultSearch({
        action: "renew",
        label: "生成续费账单",
        success: false,
        message: "当前服务不满足续费账单生成条件",
      })}`,
    );
  }

  await writeAuditLog({
    module: "billing",
    action: "portal-renew",
    targetType: "service",
    targetId: ownedService.id,
    summary: `门户用户创建续费账单：${ownedService.serviceNo}`,
    detail: {
      portalUserId: portalUser.id,
      portalUserName: portalUser.name,
      portalUserEmail: portalUser.email,
      invoiceNo: created.invoice.invoiceNo,
      orderNo: created.order?.orderNo,
      created: created.created,
    },
  });

  await queueNotification({
    customerId: portalUser.customerId,
    recipient: portalUser.email,
    recipientName: portalUser.name,
    channel: "SYSTEM",
    module: "billing",
    relatedType: "service",
    relatedId: ownedService.id,
    subject: "续费账单已生成",
    content: created.created
      ? `服务 ${ownedService.serviceNo} 的续费账单 ${created.invoice.invoiceNo} 已生成。`
      : `服务 ${ownedService.serviceNo} 已存在未支付续费账单 ${created.invoice.invoiceNo}。`,
  });

  revalidatePath("/portal");
  revalidatePath("/portal/services");
  revalidatePath(`/portal/services/${ownedService.id}`);
  revalidatePath("/portal/orders");
  revalidatePath("/portal/invoices");

  redirect(
    `/portal/services/${ownedService.id}?${buildServiceResultSearch({
      action: "renew",
      label: "生成续费账单",
      success: true,
      message: created.created ? "续费账单已生成" : "已存在未支付续费账单",
      invoiceNo: created.invoice.invoiceNo,
    })}`,
  );
}
