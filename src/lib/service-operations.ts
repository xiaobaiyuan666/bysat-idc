import { type Prisma, type ServiceStatus } from "@prisma/client";

import { db } from "@/lib/db";
import { getCloudProvider } from "@/lib/providers";
import {
  type ProviderAction,
  type ProviderActionPayload,
  type ProviderResponse,
} from "@/lib/providers/types";
import {
  buildProviderPayload,
  getNextDueDate,
  syncServiceRecord,
  writeProviderLog,
} from "@/lib/provider-sync";

export const serviceActionList = [
  "sync",
  "activate",
  "suspend",
  "unsuspend",
  "renew",
  "terminate",
  "powerOn",
  "powerOff",
  "reboot",
  "hardReboot",
  "hardPowerOff",
  "reinstall",
  "resetPassword",
  "getVnc",
  "rescueStart",
  "rescueStop",
  "lock",
  "unlock",
] as const satisfies ProviderAction[];

export type SupportedServiceAction = (typeof serviceActionList)[number];

export const serviceActionLabelMap: Record<SupportedServiceAction, string> = {
  sync: "同步状态",
  activate: "开通实例",
  suspend: "暂停实例",
  unsuspend: "解除暂停",
  renew: "续费实例",
  terminate: "终止实例",
  powerOn: "开机",
  powerOff: "关机",
  reboot: "重启",
  hardReboot: "强制重启",
  hardPowerOff: "强制关机",
  reinstall: "重装系统",
  resetPassword: "重置密码",
  getVnc: "获取 VNC",
  rescueStart: "进入救援模式",
  rescueStop: "退出救援模式",
  lock: "锁定实例",
  unlock: "解除锁定",
};

const serviceActionBaseInclude = {
  customer: true,
  product: true,
  plan: true,
  order: true,
} satisfies Prisma.ServiceInstanceInclude;

const serviceActionResponseInclude = {
  customer: true,
  product: true,
  plan: {
    include: {
      region: true,
      zone: true,
      flavor: true,
      image: true,
    },
  },
  order: true,
  vpcNetwork: true,
  _count: {
    select: {
      ipAddresses: true,
      disks: true,
      snapshots: true,
      backups: true,
      securityGroups: true,
    },
  },
} satisfies Prisma.ServiceInstanceInclude;

export function isSupportedServiceAction(
  value: string,
): value is SupportedServiceAction {
  return serviceActionList.includes(value as SupportedServiceAction);
}

export function isStatusChangingServiceAction(action: SupportedServiceAction) {
  return [
    "activate",
    "suspend",
    "unsuspend",
    "terminate",
    "powerOn",
    "powerOff",
    "renew",
    "reinstall",
    "rescueStart",
    "rescueStop",
    "hardPowerOff",
  ].includes(action);
}

function shouldSyncServiceState(action: SupportedServiceAction) {
  return !["getVnc", "resetPassword", "lock", "unlock"].includes(action);
}

export function fallbackServiceStatus(action: SupportedServiceAction) {
  if (["suspend", "powerOff", "hardPowerOff"].includes(action)) {
    return "SUSPENDED" as const;
  }

  if (action === "terminate") {
    return "TERMINATED" as const;
  }

  if (action === "reinstall") {
    return "PROVISIONING" as const;
  }

  return "ACTIVE" as const;
}

function buildRejectedProviderResult(
  service: {
    status: ServiceStatus;
    providerResourceId?: string | null;
    ipAddress?: string | null;
    region?: string | null;
  },
  message: string,
): ProviderResponse {
  return {
    ok: false,
    status: service.status,
    remoteId: service.providerResourceId ?? undefined,
    ipAddress: service.ipAddress ?? undefined,
    region: service.region ?? undefined,
    message,
  };
}

function getBlockedActionMessage(
  serviceStatus: ServiceStatus,
  action: SupportedServiceAction,
) {
  if (serviceStatus === "TERMINATED" && action !== "sync") {
    return "实例已终止，无法继续执行该操作。";
  }

  if (serviceStatus === "EXPIRED" && !["sync", "renew", "terminate"].includes(action)) {
    return "实例已到期，请先续费后再执行该操作。";
  }

  if (serviceStatus === "FAILED" && !["sync", "reinstall", "terminate"].includes(action)) {
    return "实例当前处于异常状态，请先同步或重装后再继续。";
  }

  if (serviceStatus === "SUSPENDED") {
    if (["sync", "unsuspend", "renew", "terminate"].includes(action)) {
      return undefined;
    }

    if (action === "getVnc") {
      return "实例已暂停，VNC 控制台不可用，请先解除暂停。";
    }

    return "实例已暂停，请先解除暂停后再执行该操作。";
  }

  if (action === "unsuspend") {
    return "实例当前不是暂停状态。";
  }

  if (action === "getVnc" && serviceStatus !== "ACTIVE") {
    return "VNC 仅支持运行中的实例，请先恢复实例后再试。";
  }

  return undefined;
}

export async function getServiceActionContext(serviceId: string) {
  return db.serviceInstance.findUnique({
    where: {
      id: serviceId,
    },
    include: serviceActionBaseInclude,
  });
}

export async function getServiceActionResult(serviceId: string) {
  return db.serviceInstance.findUnique({
    where: {
      id: serviceId,
    },
    include: serviceActionResponseInclude,
  });
}

export async function executeServiceAction(
  serviceId: string,
  action: SupportedServiceAction,
  payload: ProviderActionPayload = {},
) {
  const service = await getServiceActionContext(serviceId);

  if (!service) {
    return null;
  }

  const blockedMessage = getBlockedActionMessage(service.status, action);

  if (blockedMessage) {
    const providerResult = buildRejectedProviderResult(service, blockedMessage);

    await writeProviderLog({
      action,
      resourceId: service.providerResourceId,
      message: blockedMessage,
      ok: false,
    });

    return {
      originalService: service,
      data: await getServiceActionResult(service.id),
      providerResult,
    };
  }

  if (service.providerType === "MANUAL") {
    const updated = await db.serviceInstance.update({
      where: {
        id: service.id,
      },
      data: {
        status: isStatusChangingServiceAction(action)
          ? fallbackServiceStatus(action)
          : service.status,
        expiresAt: action === "terminate" ? new Date() : service.expiresAt,
        nextDueDate:
          action === "renew"
            ? getNextDueDate(service.nextDueDate, service.billingCycle)
            : service.nextDueDate,
        lastSyncAt: new Date(),
      },
    });

    const providerResult: ProviderResponse = {
      ok: true,
      status: updated.status,
      remoteId: updated.providerResourceId ?? undefined,
      ipAddress: updated.ipAddress ?? undefined,
      region: updated.region ?? undefined,
      message: `人工服务已执行“${serviceActionLabelMap[action]}”`,
      consoleUrl:
        action === "getVnc"
          ? `https://manual-console.local/${service.serviceNo.toLowerCase()}`
          : undefined,
    };

    await writeProviderLog({
      action,
      resourceId: service.providerResourceId,
      message: providerResult.message,
      ok: true,
    });

    const refreshed = await getServiceActionResult(service.id);

    return {
      originalService: service,
      data: refreshed ?? updated,
      providerResult,
    };
  }

  const provider = getCloudProvider();
  const providerPayload = buildProviderPayload(service);
  const providerResult = await provider.manageServiceAction(
    providerPayload,
    action,
    payload,
  );

  if (providerResult.ok && shouldSyncServiceState(action)) {
    await syncServiceRecord({
      service,
      providerResult,
      fallbackStatus: isStatusChangingServiceAction(action)
        ? fallbackServiceStatus(action)
        : service.status,
      nextDueDate:
        action === "renew"
          ? getNextDueDate(
              service.nextDueDate,
              service.billingCycle,
              providerResult.nextDueDate,
            )
          : undefined,
    });
  } else {
    await db.serviceInstance.update({
      where: {
        id: service.id,
      },
      data: {
        lastSyncAt: new Date(),
      },
    });
  }

  if (providerResult.ok && action === "terminate") {
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

  const refreshed = await getServiceActionResult(service.id);

  return {
    originalService: service,
    data: refreshed,
    providerResult,
  };
}
