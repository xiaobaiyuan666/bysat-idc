import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageCustomers, canViewCustomers, getApiUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { customerCertificationSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

function parseDate(value?: string | null) {
  if (!value) {
    return undefined;
  }

  const date = new Date(value);
  return Number.isNaN(date.getTime()) ? undefined : date;
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
  const certification = await db.customerCertification.findUnique({
    where: {
      customerId: id,
    },
  });

  return NextResponse.json({ data: certification });
}

export async function PUT(
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
  const parsed = customerCertificationSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "实名认证参数不正确", errors: parsed.error.flatten() },
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

  const certification = await db.customerCertification.upsert({
    where: {
      customerId: id,
    },
    update: {
      subjectType: parsed.data.subjectType,
      subjectName: parsed.data.subjectName,
      idNumber: normalizeText(parsed.data.idNumber),
      businessLicenseNo: normalizeText(parsed.data.businessLicenseNo),
      status: parsed.data.status,
      submittedAt: parseDate(parsed.data.submittedAt) ?? undefined,
      verifiedAt: parsed.data.status === "VERIFIED" ? new Date() : null,
      reviewNote: normalizeText(parsed.data.reviewNote),
    },
    create: {
      customerId: id,
      subjectType: parsed.data.subjectType,
      subjectName: parsed.data.subjectName,
      idNumber: normalizeText(parsed.data.idNumber),
      businessLicenseNo: normalizeText(parsed.data.businessLicenseNo),
      status: parsed.data.status,
      submittedAt: parseDate(parsed.data.submittedAt) ?? new Date(),
      verifiedAt: parsed.data.status === "VERIFIED" ? new Date() : undefined,
      reviewNote: normalizeText(parsed.data.reviewNote),
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "customer",
    action: "update-certification",
    targetType: "customer",
    targetId: id,
    summary: `更新客户 ${customer.name} 的实名认证资料`,
    detail: {
      status: certification.status,
      subjectType: certification.subjectType,
    },
  });

  return NextResponse.json({ data: certification });
}
