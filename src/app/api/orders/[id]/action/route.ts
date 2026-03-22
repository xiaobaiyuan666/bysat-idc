import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageOrders, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { orderActionSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

export async function POST(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageOrders(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const body = await request.json();
  const parsed = orderActionSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "订单动作参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const order = await db.order.findUnique({
    where: {
      id,
    },
    include: {
      customer: {
        select: {
          name: true,
        },
      },
      invoices: true,
      services: true,
      payments: true,
    },
  });

  if (!order) {
    return NextResponse.json({ message: "订单不存在" }, { status: 404 });
  }

  if (parsed.data.action === "approve") {
    const fullyPaid =
      order.paidAmount >= order.totalAmount ||
      order.invoices.some((invoice) => invoice.paidAmount >= invoice.totalAmount);

    if (!fullyPaid) {
      return NextResponse.json(
        { message: "订单尚未完成付款，不能执行核验通过" },
        { status: 400 },
      );
    }

    const now = new Date();

    const updatedOrder = await db.$transaction(async (tx) => {
      await Promise.all(
        order.invoices.map((invoice) =>
          tx.invoice.update({
            where: {
              id: invoice.id,
            },
            data: {
              status:
                invoice.paidAmount >= invoice.totalAmount ? "PAID" : invoice.status,
              paidAt:
                invoice.paidAmount >= invoice.totalAmount
                  ? invoice.paidAt ?? now
                  : invoice.paidAt,
            },
          }),
        ),
      );

      await tx.serviceInstance.updateMany({
        where: {
          orderId: id,
          status: {
            in: ["PENDING", "FAILED"],
          },
        },
        data: {
          status: "PROVISIONING",
        },
      });

      return tx.order.update({
        where: {
          id,
        },
        data: {
          status: order.services.every((service) => service.status === "ACTIVE")
            ? "ACTIVE"
            : "PROVISIONING",
          paidAt: order.paidAt ?? now,
          notes: normalizeText(
            [order.notes, parsed.data.reason && `核验说明：${parsed.data.reason}`]
              .filter(Boolean)
              .join("\n"),
          ),
        },
      });
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "order",
      action: "approve",
      targetType: "order",
      targetId: id,
      summary: `核验通过订单 ${order.orderNo}`,
      detail: {
        customer: order.customer.name,
        reason: normalizeText(parsed.data.reason),
      },
    });

    return NextResponse.json({ data: updatedOrder });
  }

  if (order.paidAmount > 0 || order.payments.length > 0) {
    return NextResponse.json(
      { message: "订单已有收款记录，请先完成退款后再取消" },
      { status: 400 },
    );
  }

  const cancelledOrder = await db.$transaction(async (tx) => {
    await tx.invoice.updateMany({
      where: {
        orderId: id,
      },
      data: {
        status: "VOID",
      },
    });

    await tx.serviceInstance.updateMany({
      where: {
        orderId: id,
        status: {
          in: ["PENDING", "PROVISIONING", "ACTIVE", "FAILED"],
        },
      },
      data: {
        status: "TERMINATED",
      },
    });

    return tx.order.update({
      where: {
        id,
      },
      data: {
        status: "CANCELLED",
        notes: normalizeText(
          [order.notes, parsed.data.reason && `取消原因：${parsed.data.reason}`]
            .filter(Boolean)
            .join("\n"),
        ),
      },
    });
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "order",
    action: "cancel",
    targetType: "order",
    targetId: id,
    summary: `取消订单 ${order.orderNo}`,
    detail: {
      customer: order.customer.name,
      reason: normalizeText(parsed.data.reason),
    },
  });

  return NextResponse.json({ data: cancelledOrder });
}
