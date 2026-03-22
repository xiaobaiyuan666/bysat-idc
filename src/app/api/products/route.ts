import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { canManageProducts, getApiUser } from "@/lib/api-auth";
import { getProductsPageData } from "@/lib/data";
import { db } from "@/lib/db";
import { productSchema } from "@/lib/validation";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

function toCents(value: unknown) {
  return Math.round(Number(value ?? 0) * 100);
}

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageProducts(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const data = await getProductsPageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
  }

  if (!canManageProducts(user)) {
    return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
  }

  const body = await request.json();
  const parsed = productSchema.safeParse({
    ...body,
    price: toCents(body.price),
    setupFee: toCents(body.setupFee),
    stock: Number(body.stock ?? 0),
  });

  if (!parsed.success) {
    return NextResponse.json(
      { message: "产品参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const product = await db.product.create({
    data: {
      code: parsed.data.code.trim().toUpperCase(),
      name: parsed.data.name,
      category: parsed.data.category,
      status: parsed.data.status,
      billingCycle: parsed.data.billingCycle,
      price: parsed.data.price,
      setupFee: parsed.data.setupFee,
      stock: parsed.data.stock,
      autoProvision: parsed.data.autoProvision,
      providerType: parsed.data.providerType,
      providerProductId: normalizeText(parsed.data.providerProductId),
      regionTemplate: normalizeText(parsed.data.regionTemplate),
      description: normalizeText(parsed.data.description),
    },
  });

  await writeAuditLog({
    adminUserId: user.id,
    module: "product",
    action: "create",
    targetType: "product",
    targetId: product.id,
    summary: `创建产品 ${product.code}`,
  });

  return NextResponse.json({ data: product }, { status: 201 });
}
