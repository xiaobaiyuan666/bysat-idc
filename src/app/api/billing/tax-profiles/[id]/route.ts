import { Prisma } from "@prisma/client";
import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageBilling, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { taxProfileSchema } from "@/lib/validation";

export async function PUT(
  request: NextRequest,
  { params }: { params: Promise<{ id: string }> },
) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已过期" }, { status: 401 });
  }

  if (!canManageBilling(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = taxProfileSchema.safeParse({
    code: body.code,
    name: body.name,
    taxRate: Number(body.taxRate ?? 0),
    description: body.description,
    isDefault: Boolean(body.isDefault),
    isActive: body.isActive !== false,
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "税率档案参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const { id } = await params;
  const current = await db.invoiceTaxProfile.findUnique({
    where: {
      id,
    },
  });

  if (!current) {
    return NextResponse.json({ message: "税率档案不存在" }, { status: 404 });
  }

  try {
    const data = await db.$transaction(async (tx) => {
      if (parsed.data.isDefault) {
        await tx.invoiceTaxProfile.updateMany({
          where: {
            id: {
              not: id,
            },
          },
          data: {
            isDefault: false,
          },
        });
      }

      return tx.invoiceTaxProfile.update({
        where: {
          id,
        },
        data: {
          ...parsed.data,
          code: parsed.data.code.trim().toUpperCase(),
        },
      });
    });

    await writeAuditLog({
      adminUserId: user.id,
      module: "billing",
      action: "update-settings",
      targetType: "invoiceTaxProfile",
      targetId: data.id,
      summary: `更新税率档案：${data.name}`,
      detail: data,
    });

    return NextResponse.json({ data });
  } catch (error) {
    if (
      error instanceof Prisma.PrismaClientKnownRequestError &&
      error.code === "P2002"
    ) {
      return NextResponse.json({ message: "税率编码已存在" }, { status: 409 });
    }

    throw error;
  }
}
