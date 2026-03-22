import crypto from "node:crypto";

import {
  type PaymentGatewayConfig,
  type PaymentMethod,
  type PaymentSignType,
} from "@prisma/client";
import { NextRequest } from "next/server";

import { db } from "@/lib/db";

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

function normalizeHeaderName(
  value: string | null | undefined,
  signType: PaymentSignType,
) {
  const normalized = normalizeText(value)?.toLowerCase();

  if (normalized) {
    return normalized;
  }

  return signType === "HMAC_SHA256" ? "x-payment-signature" : "x-payment-secret";
}

export function getGatewayCallbackHeaderName(config: {
  callbackHeader?: string | null;
  signType: PaymentSignType;
}) {
  return normalizeHeaderName(config.callbackHeader, config.signType);
}

function safeEqual(left: string, right: string) {
  const leftBuffer = Buffer.from(left);
  const rightBuffer = Buffer.from(right);

  if (leftBuffer.length !== rightBuffer.length) {
    return false;
  }

  return crypto.timingSafeEqual(leftBuffer, rightBuffer);
}

function normalizeSignature(signature: string) {
  return signature.trim().replace(/^sha256=/i, "").toLowerCase();
}

export async function getPaymentGatewayConfig(method: PaymentMethod) {
  return db.paymentGatewayConfig.findUnique({
    where: {
      method,
    },
  });
}

export function normalizeGatewayInput(input: {
  method: PaymentMethod;
  name: string;
  merchantId?: string;
  appId?: string;
  apiBaseUrl?: string;
  signType: PaymentSignType;
  callbackSecret: string;
  callbackHeader?: string;
  notifyUrl?: string;
  returnUrl?: string;
  isEnabled: boolean;
  remark?: string;
}) {
  return {
    method: input.method,
    name: input.name.trim(),
    merchantId: normalizeText(input.merchantId),
    appId: normalizeText(input.appId),
    apiBaseUrl: normalizeText(input.apiBaseUrl),
    signType: input.signType,
    callbackSecret: input.callbackSecret.trim(),
    callbackHeader: normalizeHeaderName(input.callbackHeader, input.signType),
    notifyUrl: normalizeText(input.notifyUrl),
    returnUrl: normalizeText(input.returnUrl),
    isEnabled: input.isEnabled,
    remark: normalizeText(input.remark),
  };
}

export function computeGatewaySignature(
  config: Pick<PaymentGatewayConfig, "callbackSecret" | "signType">,
  rawBody: string,
) {
  if (config.signType !== "HMAC_SHA256") {
    return config.callbackSecret;
  }

  return crypto
    .createHmac("sha256", config.callbackSecret)
    .update(rawBody)
    .digest("hex");
}

export function verifyGatewaySignature(input: {
  config: Pick<PaymentGatewayConfig, "callbackSecret" | "callbackHeader" | "signType">;
  request: NextRequest;
  rawBody: string;
  bodySignature?: string;
}) {
  const headerName = normalizeHeaderName(
    input.config.callbackHeader,
    input.config.signType,
  );
  const providedSignature =
    input.request.headers.get(headerName) ??
    input.bodySignature ??
    input.request.headers.get("x-payment-signature") ??
    input.request.headers.get("x-payment-secret");

  if (!providedSignature) {
    return {
      ok: false,
      reason: "MISSING_SIGNATURE",
    } as const;
  }

  const expectedSignature = computeGatewaySignature(input.config, input.rawBody);

  if (input.config.signType === "HMAC_SHA256") {
    const ok = safeEqual(
      normalizeSignature(providedSignature),
      normalizeSignature(expectedSignature),
    );

    return {
      ok,
      reason: ok ? null : "INVALID_SIGNATURE",
    } as const;
  }

  const ok = safeEqual(providedSignature.trim(), expectedSignature.trim());

  return {
    ok,
    reason: ok ? null : "INVALID_SIGNATURE",
  } as const;
}
