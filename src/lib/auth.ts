import bcrypt from "bcryptjs";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";

import { db } from "@/lib/db";
import {
  AUTH_COOKIE_NAME,
  createSessionToken,
  type SessionPayload,
  verifySessionToken,
} from "@/lib/session";

const cookieOptions = {
  httpOnly: true,
  sameSite: "lax" as const,
  secure: process.env.NODE_ENV === "production",
  path: "/",
  maxAge: 60 * 60 * 12,
};

export async function getSession() {
  const cookieStore = await cookies();
  const token = cookieStore.get(AUTH_COOKIE_NAME)?.value;

  if (!token) {
    return null;
  }

  return verifySessionToken(token);
}

export async function getCurrentUser() {
  const session = await getSession();

  if (!session) {
    return null;
  }

  const admin = await db.adminUser.findUnique({
    where: { id: session.sub },
  });

  if (!admin || !admin.isActive) {
    return null;
  }

  return admin;
}

export async function requireUser() {
  const user = await getCurrentUser();

  if (!user) {
    redirect("/login");
  }

  return user;
}

async function writeSessionCookie(payload: SessionPayload) {
  const cookieStore = await cookies();
  const token = await createSessionToken(payload);

  cookieStore.set(AUTH_COOKIE_NAME, token, cookieOptions);
}

export async function loginWithPassword(email: string, password: string) {
  const admin = await db.adminUser.findUnique({
    where: { email },
  });

  if (!admin || !admin.isActive) {
    return false;
  }

  const matched = await bcrypt.compare(password, admin.passwordHash);
  if (!matched) {
    return false;
  }

  await db.adminUser.update({
    where: { id: admin.id },
    data: {
      lastLoginAt: new Date(),
    },
  });

  await writeSessionCookie({
    sub: admin.id,
    name: admin.name,
    email: admin.email,
    role: admin.role,
  });

  return true;
}

export async function logout() {
  const cookieStore = await cookies();
  cookieStore.delete(AUTH_COOKIE_NAME);
}
