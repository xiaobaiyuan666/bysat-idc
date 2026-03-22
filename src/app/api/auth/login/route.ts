import bcrypt from "bcryptjs";
import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import { serializeAuthUser } from "@/lib/api-auth";
import { db } from "@/lib/db";
import { loginSchema } from "@/lib/validation";
import { AUTH_COOKIE_NAME, createSessionToken } from "@/lib/session";

const cookieOptions = {
  httpOnly: true,
  sameSite: "lax" as const,
  secure: process.env.NODE_ENV === "production",
  path: "/",
  maxAge: 60 * 60 * 12,
};

export async function POST(request: NextRequest) {
  const body = await request.json();
  const parsed = loginSchema.safeParse(body);

  if (!parsed.success) {
    return NextResponse.json(
      { message: "请求参数不正确", errors: parsed.error.flatten() },
      { status: 400 },
    );
  }

  const admin = await db.adminUser.findUnique({
    where: {
      email: parsed.data.email,
    },
  });

  if (!admin || !admin.isActive) {
    return NextResponse.json({ message: "账号不存在或已被停用" }, { status: 401 });
  }

  const matched = await bcrypt.compare(parsed.data.password, admin.passwordHash);
  if (!matched) {
    return NextResponse.json({ message: "邮箱或密码错误" }, { status: 401 });
  }

  await db.adminUser.update({
    where: {
      id: admin.id,
    },
    data: {
      lastLoginAt: new Date(),
    },
  });

  const token = await createSessionToken({
    sub: admin.id,
    name: admin.name,
    email: admin.email,
    role: admin.role,
  });

  const response = NextResponse.json({
    data: serializeAuthUser(admin),
  });

  response.cookies.set(AUTH_COOKIE_NAME, token, cookieOptions);

  await writeAuditLog({
    adminUserId: admin.id,
    module: "auth",
    action: "login",
    targetType: "adminUser",
    targetId: admin.id,
    summary: `管理员登录：${admin.email}`,
  });

  return response;
}
