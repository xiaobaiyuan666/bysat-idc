"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireUser } from "@/lib/auth";
import { db } from "@/lib/db";
import { addCycle } from "@/lib/format";
import { orderSchema } from "@/lib/validation";
import { makeCode } from "@/lib/utils";

export async function createOrderAction(formData: FormData) {
  await requireUser();

  const parsed = orderSchema.safeParse({
    customerId: formData.get("customerId"),
    productId: formData.get("productId"),
    serviceName: formData.get("serviceName"),
    cycle: formData.get("cycle"),
    quantity: Number(formData.get("quantity") || 1),
    notes: formData.get("notes") || undefined,
  });

  if (!parsed.success) {
    redirect("/orders");
  }

  const [customer, product] = await Promise.all([
    db.customer.findUnique({ where: { id: parsed.data.customerId } }),
    db.product.findUnique({ where: { id: parsed.data.productId } }),
  ]);

  if (!customer || !product) {
    redirect("/orders");
  }

  const totalAmount = product.price * parsed.data.quantity + product.setupFee;
  const nextDueDate = addCycle(new Date(), parsed.data.cycle);

  await db.$transaction(async (tx) => {
    const order = await tx.order.create({
      data: {
        orderNo: makeCode("ORD"),
        customerId: customer.id,
        status: "PENDING",
        totalAmount,
        source: "后台录单",
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
        orderId: order.id,
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
          source: "manual-order-form",
          productCode: product.code,
        }),
      },
    });

    await tx.orderItem.create({
      data: {
        orderId: order.id,
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
        orderId: order.id,
        serviceId: service.id,
        type: "ORDER",
        status: "ISSUED",
        subtotal: totalAmount,
        totalAmount,
        dueDate: nextDueDate,
        issuedAt: new Date(),
        remark: "后台订单自动生成账单",
      },
    });
  });

  ["/dashboard", "/orders", "/services", "/invoices"].forEach((path) =>
    revalidatePath(path),
  );

  redirect("/orders");
}
