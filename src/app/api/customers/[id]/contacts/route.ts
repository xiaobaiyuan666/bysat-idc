import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, canViewCustomers, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { customerContactSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

export async function GET(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canViewCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const contacts = await db.customerContact.findMany({
    where: {
      customerId: id,
    },
    orderBy: [
      {
        isPrimary: "desc",
      },
      {
        createdAt: "asc",
      },
    ],
  });

  return NextResponse.json({ data: contacts });
}

export async function POST(
  request: NextRequest,
  context: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id } = await context.params;
  const body = await request.json();
  const parsed = customerContactSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "联系人参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const customer = await db.customer.findUnique({
    where: {
      id,
    },
    select: {
      id: true,
      name: true,
    },
  });

  if (!customer) {
    return NextResponse.json({ message: "客户不存在" }, { status: 404 });
  }

  const contact = await db.$transaction(async (tx) => {
    if (parsed.data.isPrimary) {
      await tx.customerContact.updateMany({
        where: {
          customerId: id,
          isPrimary: true,
        },
        data: {
          isPrimary: false,
        },
      });
    }

    return tx.customerContact.create({
      data: {
        customerId: id,
        name: parsed.data.name,
        department: normalizeText(parsed.data.department),
        role: normalizeText(parsed.data.role),
        email: normalizeText(parsed.data.email),
        phone: normalizeText(parsed.data.phone),
        isPrimary: parsed.data.isPrimary,
        notes: normalizeText(parsed.data.notes),
      },
    });
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "create-contact",
    targetType: "customer",
    targetId: id,
    summary: `为客户 ${customer.name} 新增联系人 ${contact.name}`,
    detail: {
      contactId: contact.id,
      isPrimary: contact.isPrimary,
    },
  });

  return NextResponse.json({ data: contact }, { status: 201 });
}
