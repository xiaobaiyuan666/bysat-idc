import bcrypt from "bcryptjs";
import { cookies } from "next/headers";
import { redirect } from "next/navigation";
import { SignJWT, jwtVerify } from "jose";

import { db } from "@/lib/db";

export const PORTAL_COOKIE_NAME = "idc_portal_session";

type PortalSessionPayload = {
  sub: string;
  customerId: string;
  name: string;
  email: string;
  role: string;
};

const cookieOptions = {
  httpOnly: true,
  sameSite: "lax" as const,
  secure: process.env.NODE_ENV === "production",
  path: "/",
  maxAge: 60 * 60 * 12,
};

function getSessionSecret() {
  return new TextEncoder().encode(
    process.env.AUTH_SECRET ?? "replace-this-with-a-long-secret",
  );
}

async function createPortalToken(payload: PortalSessionPayload) {
  return new SignJWT(payload)
    .setProtectedHeader({ alg: "HS256" })
    .setSubject(payload.sub)
    .setIssuedAt()
    .setExpirationTime("12h")
    .sign(getSessionSecret());
}

async function verifyPortalToken(token: string) {
  try {
    const verified = await jwtVerify<PortalSessionPayload>(token, getSessionSecret());
    return verified.payload;
  } catch {
    return null;
  }
}

export async function getPortalSession() {
  const cookieStore = await cookies();
  const token = cookieStore.get(PORTAL_COOKIE_NAME)?.value;

  if (!token) {
    return null;
  }

  return verifyPortalToken(token);
}

export async function getCurrentPortalUser() {
  const session = await getPortalSession();

  if (!session) {
    return null;
  }

  const user = await db.customerUser.findUnique({
    where: {
      id: session.sub,
    },
    include: {
      customer: true,
    },
  });

  if (!user || !user.isActive || user.customer.status === "ARCHIVED") {
    return null;
  }

  return user;
}

export async function requirePortalUser() {
  const user = await getCurrentPortalUser();

  if (!user) {
    redirect("/portal/login");
  }

  return user;
}

export async function loginPortalWithPassword(email: string, password: string) {
  const user = await db.customerUser.findUnique({
    where: {
      email,
    },
  });

  if (!user || !user.isActive) {
    return false;
  }

  const matched = await bcrypt.compare(password, user.passwordHash);

  if (!matched) {
    return false;
  }

  await db.customerUser.update({
    where: {
      id: user.id,
    },
    data: {
      lastLoginAt: new Date(),
    },
  });

  const cookieStore = await cookies();
  const token = await createPortalToken({
    sub: user.id,
    customerId: user.customerId,
    name: user.name,
    email: user.email,
    role: user.role,
  });

  cookieStore.set(PORTAL_COOKIE_NAME, token, cookieOptions);
  return true;
}

export async function logoutPortal() {
  const cookieStore = await cookies();
  cookieStore.delete(PORTAL_COOKIE_NAME);
}
