"use server";

import { revalidatePath } from "next/cache";
import { redirect } from "next/navigation";

import { requireUser } from "@/lib/auth";
import { db } from "@/lib/db";
import { getCloudProvider } from "@/lib/providers";
import {
  buildProviderPayload,
  getNextDueDate,
  syncServiceRecord,
  writeProviderLog,
} from "@/lib/provider-sync";

export async function runServiceAction(formData: FormData) {
  await requireUser();

  const serviceId = String(formData.get("serviceId") ?? "");
  const action = String(formData.get("action") ?? "");

  if (!serviceId || !action) {
    redirect("/services");
  }

  const service = await db.serviceInstance.findUnique({
    where: {
      id: serviceId,
    },
    include: {
      customer: true,
      product: true,
    },
  });

  if (!service) {
    redirect("/services");
  }

  if (service.providerType === "MANUAL") {
    const statusMap = {
      activate: "ACTIVE",
      suspend: "SUSPENDED",
      terminate: "TERMINATED",
      renew: "ACTIVE",
      sync: service.status,
    } as const;

    await db.serviceInstance.update({
      where: {
        id: service.id,
      },
      data: {
        status: statusMap[action as keyof typeof statusMap] ?? service.status,
        expiresAt: action === "terminate" ? new Date() : service.expiresAt,
        nextDueDate:
          action === "renew"
            ? getNextDueDate(service.nextDueDate, service.billingCycle)
            : service.nextDueDate,
        lastSyncAt: new Date(),
      },
    });

    ["/dashboard", "/services"].forEach((path) => revalidatePath(path));
    redirect("/services");
  }

  const provider = getCloudProvider();
  const payload = buildProviderPayload(service);

  const providerResult =
    action === "activate"
      ? await provider.activateService(payload)
      : action === "suspend"
        ? await provider.suspendService(payload)
        : action === "terminate"
          ? await provider.terminateService(payload)
          : action === "renew"
            ? await provider.renewService(payload)
            : await provider.syncService(payload);

  await syncServiceRecord({
    service,
    providerResult,
    fallbackStatus:
      action === "suspend"
        ? "SUSPENDED"
        : action === "terminate"
          ? "TERMINATED"
          : "ACTIVE",
    nextDueDate:
      action === "renew"
        ? getNextDueDate(
            service.nextDueDate,
            service.billingCycle,
            providerResult.nextDueDate,
          )
        : undefined,
  });

  if (action === "terminate") {
    await db.serviceInstance.update({
      where: {
        id: service.id,
      },
      data: {
        expiresAt: new Date(),
      },
    });
  }

  await writeProviderLog({
    action,
    resourceId: providerResult.remoteId ?? service.providerResourceId,
    message: providerResult.message,
    requestBody: providerResult.requestBody,
    responseBody: providerResult.responseBody,
    ok: providerResult.ok,
  });

  ["/dashboard", "/services"].forEach((path) => revalidatePath(path));
  redirect("/services");
}
