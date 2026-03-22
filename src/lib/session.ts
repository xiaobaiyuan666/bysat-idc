import { SignJWT, jwtVerify } from "jose";

export const AUTH_COOKIE_NAME = "idc_finance_session";

export type SessionPayload = {
  sub: string;
  name: string;
  email: string;
  role: string;
};

function getSessionSecret() {
  return new TextEncoder().encode(
    process.env.AUTH_SECRET ?? "replace-this-with-a-long-secret",
  );
}

export async function createSessionToken(payload: SessionPayload) {
  return new SignJWT(payload)
    .setProtectedHeader({ alg: "HS256" })
    .setSubject(payload.sub)
    .setIssuedAt()
    .setExpirationTime("12h")
    .sign(getSessionSecret());
}

export async function verifySessionToken(token: string) {
  try {
    const verified = await jwtVerify<SessionPayload>(token, getSessionSecret());
    return verified.payload;
  } catch {
    return null;
  }
}
