import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageOrders, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { getOrdersPageData } from "@/lib/data";
import { addCycle } from "@/lib/format";
import { orderSchema } from "@/lib/validation";
import { makeCode } from "@/lib/utils";

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageOrders(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getOrdersPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageOrders(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = orderSchema.safeParse({
    customerId: body.customerId,
    productId: body.productId,
    serviceName: body.serviceName,
    cycle: body.cycle,
    quantity: Number(body.quantity ?? 1),
    notes: body.notes,
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "请求参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const [customer, product] = await Promise.all([
    db.customer.findUnique({ where: { id: parsed.data.customerId } }),
    db.product.findUnique({ where: { id: parsed.data.productId } }),
  ]);

  if (!customer || !product) {
    return NextResponse.json(
      { message: "客户或产品不存在" },
      { status: 404 },
    );
  }

  const totalAmount = product.price * parsed.data.quantity + product.setupFee;
  const nextDueDate = addCycle(new Date(), parsed.data.cycle);

  const order = await db.$transaction(async (tx) => {
    const createdOrder = await tx.order.create({
      data: {
        orderNo: makeCode("ORD"),
        customerId: customer.id,
        status: "PENDING",
        totalAmount,
        source: "admin",
        orderType: "new",
        dueDate: nextDueDate,
        notes: parsed.data.notes,
      },
    });

    const service = await tx.serviceInstance.create({
      data: {
        serviceNo: makeCode("SRV"),
        customerId: customer.id,
        productId: product.id,
        orderId: createdOrder.id,
        name: parsed.data.serviceName,
        hostname: `${product.code.toLowerCase()}.${customer.name.toLowerCase().replace(/\s+/g, "-")}.local`,
        providerType: product.providerType,
        region: product.regionTemplate,
        billingCycle: parsed.data.cycle,
        status: "PENDING",
        monthlyCost: product.price * parsed.data.quantity,
        nextDueDate,
        configSnapshot: JSON.stringify({
          quantity: parsed.data.quantity,
          source: "api",
          productCode: product.code,
        }),
      },
    });

    await tx.orderItem.create({
      data: {
        orderId: createdOrder.id,
        productId: product.id,
        serviceId: service.id,
        title: product.name,
        quantity: parsed.data.quantity,
        unitPrice: product.price,
        cycle: parsed.data.cycle,
        totalAmount,
      },
    });

    await tx.invoice.create({
      data: {
        invoiceNo: makeCode("INV"),
        customerId: customer.id,
        orderId: createdOrder.id,
        serviceId: service.id,
        type: "ORDER",
        status: "ISSUED",
        subtotal: totalAmount,
        totalAmount,
        dueDate: nextDueDate,
        issuedAt: new Date(),
        remark: "后台订单流程自动生成账单",
      },
    });

    return tx.order.findUnique({
      where: {
        id: createdOrder.id,
      },
      include: {
        customer: true,
        items: true,
        services: true,
        invoices: true,
      },
    });
  });

  if (order) {
    await writeAuditLog({
      adminUserId: user.id,
      module: "order",
      action: "create",
      targetType: "order",
      targetId: order.id,
      summary: `创建订单：${order.orderNo}`,
      detail: {
        customerId: customer.id,
        productId: product.id,
      },
    });
  }

  return NextResponse.json({ data: order }, { status: 201 });
}
