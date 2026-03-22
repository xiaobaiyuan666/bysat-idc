import {
  type CloudPlan,
  type Customer,
  type Product,
  ProviderSyncStatus,
  type ServiceInstance,
  type ServiceStatus,
} from "@prisma/client";

import { addCycle } from "@/lib/format";
import { db } from "@/lib/db";
import { type ProviderResponse, type ProviderServicePayload } from "@/lib/providers/types";

type ProviderServiceModel = ServiceInstance & {
  customer: Customer;
  product: Product;
  plan?: CloudPlan | null;
};

function parseConfigSnapshot(snapshot?: string | null) {
  if (!snapshot) {
    return undefined;
  }

  try {
    return JSON.parse(snapshot) as Record<string, unknown>;
  } catch {
    return undefined;
  }
}

export function buildProviderPayload(
  service: ProviderServiceModel,
): ProviderServicePayload {
  const configOptions = {
    ...(service.plan?.configOptions ? parseConfigSnapshot(service.plan.configOptions) : {}),
    ...(parseConfigSnapshot(service.configSnapshot) ?? {}),
  };

  return {
    localId: service.id,
    serviceNo: service.serviceNo,
    name: service.name,
    hostname: service.hostname,
    providerResourceId: service.providerResourceId,
    providerProductId: service.product.providerProductId,
    region: service.region,
    billingCycle: service.billingCycle,
    customerName: service.customer.name,
    productName: service.product.name,
    configOptions,
  };
}

export function mapProviderStatus(
  status: string,
  fallback: ServiceStatus,
): ServiceStatus {
  const normalized = status
    .trim()
    .replace(/([a-z])([A-Z])/g, "$1_$2")
    .replace(/[\s-]+/g, "_")
    .toUpperCase();

  if (
    [
      "ACTIVE",
      "RUNNING",
      "ON",
      "STARTED",
      "NORMAL",
      "SUCCESS",
    ].includes(normalized)
  ) {
    return "ACTIVE";
  }

  if (
    [
      "SUSPENDED",
      "SUSPEND",
      "STOPPED",
      "OFF",
      "SHUTOFF",
      "PAUSED",
      "DISABLED",
    ].includes(normalized)
  ) {
    return "SUSPENDED";
  }

  if (["TERMINATED", "DELETED", "DESTROYED", "CANCELLED", "CANCELED"].includes(normalized)) {
    return "TERMINATED";
  }

  if (
    [
      "PROVISIONING",
      "BUILDING",
      "CREATING",
      "INSTALLING",
      "DEPLOYING",
      "REINSTALLING",
      "PENDING",
    ].includes(normalized)
  ) {
    return "PROVISIONING";
  }

  if (normalized === "OVERDUE") {
    return "OVERDUE";
  }

  if (normalized === "EXPIRED") {
    return "EXPIRED";
  }

  if (["FAILED", "ERROR"].includes(normalized)) {
    return "FAILED";
  }

  return fallback;
}

export async function writeProviderLog(input: {
  action: string;
  resourceId?: string | null;
  message?: string;
  requestBody?: string;
  responseBody?: string;
  ok: boolean;
}) {
  await db.providerSyncLog.create({
    data: {
      providerType: "MOFANG_CLOUD",
      action: input.action,
      resourceType: "instance",
      resourceId: input.resourceId ?? undefined,
      message: input.message,
      requestBody: input.requestBody,
      responseBody: input.responseBody,
      status: input.ok ? ProviderSyncStatus.SUCCESS : ProviderSyncStatus.FAILED,
    },
  });
}

export function getNextDueDate(
  currentDate: Date | null | undefined,
  billingCycle: string,
  providerDate?: Date,
) {
  if (providerDate) {
    return providerDate;
  }

  return addCycle(currentDate ?? new Date(), billingCycle);
}

export async function syncServiceRecord(input: {
  service: ProviderServiceModel;
  providerResult: ProviderResponse;
  fallbackStatus: ServiceStatus;
  nextDueDate?: Date | null;
}) {
  const mappedStatus = mapProviderStatus(
    input.providerResult.status,
    input.fallbackStatus,
  );

  await db.serviceInstance.update({
    where: {
      id: input.service.id,
    },
    data: {
      status: mappedStatus,
      providerResourceId:
        input.providerResult.remoteId ?? input.service.providerResourceId,
      ipAddress: input.providerResult.ipAddress ?? input.service.ipAddress,
      region: input.providerResult.region ?? input.service.region,
      nextDueDate:
        input.nextDueDate ??
        input.providerResult.nextDueDate ??
        input.service.nextDueDate,
      activatedAt: mappedStatus === "ACTIVE" ? new Date() : input.service.activatedAt,
      lastSyncAt: new Date(),
    },
  });
}
