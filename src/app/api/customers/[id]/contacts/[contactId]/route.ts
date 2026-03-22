import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { customerContactSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

async function getScopedContact(customerId: string, contactId: string) {
  return db.customerContact.findFirst({
    where: {
      id: contactId,
      customerId,
    },
  });
}

export async function PUT(
  request: NextRequest,
  context: { params: Promise<{ id: string; contactId: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id, contactId } = await context.params;
  const body = await request.json();
  const parsed = customerContactSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "联系人参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const existing = await getScopedContact(id, contactId);

  if (!existing) {
    return NextResponse.json({ message: "联系人不存在" }, { status: 404 });
  }

  const contact = await db.$transaction(async (tx) => {
    if (parsed.data.isPrimary) {
      await tx.customerContact.updateMany({
        where: {
          customerId: id,
          isPrimary: true,
          NOT: {
            id: contactId,
          },
        },
        data: {
          isPrimary: false,
        },
      });
    }

    return tx.customerContact.update({
      where: {
        id: contactId,
      },
      data: {
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
    action: "update-contact",
    targetType: "customer",
    targetId: id,
    summary: `更新客户联系人 ${contact.name}`,
    detail: {
      contactId: contact.id,
    },
  });

  return NextResponse.json({ data: contact });
}

export async function DELETE(
  request: NextRequest,
  context: { params: Promise<{ id: string; contactId: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageCustomers(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const { id, contactId } = await context.params;
  const existing = await getScopedContact(id, contactId);

  if (!existing) {
    return NextResponse.json({ message: "联系人不存在" }, { status: 404 });
  }

  await db.customerContact.delete({
    where: {
      id: contactId,
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "delete-contact",
    targetType: "customer",
    targetId: id,
    summary: `删除客户联系人 ${existing.name}`,
    detail: {
      contactId: existing.id,
    },
  });

  return NextResponse.json({ success: true });
}
