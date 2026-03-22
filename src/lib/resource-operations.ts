import { ProviderSyncStatus, type ProviderType } from "@prisma/client";

import { db } from "@/lib/db";

export const resourceActionLabelMap = {
  disks: {
    createSnapshot: "创建快照",
    createBackup: "创建备份",
    attach: "挂载磁盘",
    detach: "卸载磁盘",
    setBoot: "设为系统盘",
  },
  snapshots: {
    restore: "恢复快照",
    deleteSnapshot: "删除快照",
  },
  backups: {
    restore: "恢复备份",
    expireNow: "归档备份",
    deleteBackup: "删除备份",
  },
  "security-groups": {
    addRule: "新增安全组规则",
    deleteRule: "删除安全组规则",
    deleteGroup: "删除安全组",
  },
} as const;

export type ResourceType = keyof typeof resourceActionLabelMap;
export type SupportedResourceAction = {
  [K in ResourceType]: keyof (typeof resourceActionLabelMap)[K];
}[ResourceType] &
  string;

type OwnershipContext = {
  resourceType: ResourceType;
  id: string;
  name: string;
  serviceId?: string;
  serviceNo?: string;
  customerId?: string;
  providerType?: ProviderType;
  providerResourceId?: string;
};

type ResourceActionResult = {
  data: unknown;
  summary: string;
  targetType: string;
  targetId: string;
  serviceId?: string;
  serviceNo?: string;
  customerId?: string;
};

export class ResourceActionError extends Error {
  constructor(
    message: string,
    public readonly statusCode = 400,
  ) {
    super(message);
    this.name = "ResourceActionError";
  }
}

function normalizeText(value: unknown) {
  if (typeof value !== "string") {
    return undefined;
  }

  const normalized = value.trim();
  return normalized ? normalized : undefined;
}

function daysLater(days: number) {
  return new Date(Date.now() + days * 24 * 60 * 60 * 1000);
}

function singularResourceType(resourceType: ResourceType) {
  if (resourceType === "security-groups") {
    return "security-group";
  }

  return resourceType.slice(0, -1);
}

async function writeResourceProviderLog(input: {
  providerType?: ProviderType;
  resourceType: ResourceType;
  resourceId?: string;
  action: string;
  message: string;
  payload?: Record<string, unknown>;
  responseBody?: unknown;
}) {
  if (!input.providerType) {
    return;
  }

  await db.providerSyncLog.create({
    data: {
      providerType: input.providerType,
      action: input.action,
      resourceType: singularResourceType(input.resourceType),
      resourceId: input.resourceId,
      status: ProviderSyncStatus.SUCCESS,
      message: input.message,
      requestBody: input.payload ? JSON.stringify(input.payload, null, 2) : undefined,
      responseBody:
        input.responseBody !== undefined
          ? JSON.stringify(input.responseBody, null, 2)
          : undefined,
    },
  });
}

async function getDisk(id: string) {
  return db.serviceDisk.findUnique({
    where: { id },
    include: {
      service: {
        select: {
          id: true,
          serviceNo: true,
          customerId: true,
          providerType: true,
          providerResourceId: true,
        },
      },
    },
  });
}

async function getSnapshot(id: string) {
  return db.serviceSnapshot.findUnique({
    where: { id },
    include: {
      service: {
        select: {
          id: true,
          serviceNo: true,
          customerId: true,
          providerType: true,
          providerResourceId: true,
        },
      },
      sourceDisk: {
        select: {
          id: true,
          name: true,
        },
      },
    },
  });
}

async function getBackup(id: string) {
  return db.serviceBackup.findUnique({
    where: { id },
    include: {
      service: {
        select: {
          id: true,
          serviceNo: true,
          customerId: true,
          providerType: true,
          providerResourceId: true,
        },
      },
    },
  });
}

async function getSecurityGroup(id: string) {
  return db.serviceSecurityGroup.findUnique({
    where: { id },
    include: {
      service: {
        select: {
          id: true,
          serviceNo: true,
          customerId: true,
          providerType: true,
          providerResourceId: true,
        },
      },
      rules: {
        orderBy: {
          createdAt: "asc",
        },
      },
      vpcNetwork: {
        select: {
          id: true,
          name: true,
        },
      },
    },
  });
}

export function isSupportedResourceType(value: string): value is ResourceType {
  return Object.prototype.hasOwnProperty.call(resourceActionLabelMap, value);
}

export function isSupportedResourceAction(
  resourceType: string,
  action: string,
): action is SupportedResourceAction {
  return (
    isSupportedResourceType(resourceType) &&
    Object.prototype.hasOwnProperty.call(resourceActionLabelMap[resourceType], action)
  );
}

export function getResourceActionLabel(resourceType: ResourceType, action: string) {
  const labels = resourceActionLabelMap[resourceType] as Record<string, string>;
  return labels[action] ?? action;
}

export async function getResourceOwnershipContext(
  resourceType: ResourceType,
  id: string,
): Promise<OwnershipContext | null> {
  if (resourceType === "disks") {
    const disk = await getDisk(id);

    if (!disk) {
      return null;
    }

    return {
      resourceType,
      id: disk.id,
      name: disk.name,
      serviceId: disk.serviceId,
      serviceNo: disk.service.serviceNo,
      customerId: disk.service.customerId,
      providerType: disk.service.providerType,
      providerResourceId: disk.providerDiskId ?? disk.service.providerResourceId ?? undefined,
    };
  }

  if (resourceType === "snapshots") {
    const snapshot = await getSnapshot(id);

    if (!snapshot) {
      return null;
    }

    return {
      resourceType,
      id: snapshot.id,
      name: snapshot.name,
      serviceId: snapshot.serviceId,
      serviceNo: snapshot.service.serviceNo,
      customerId: snapshot.service.customerId,
      providerType: snapshot.service.providerType,
      providerResourceId:
        snapshot.providerSnapshotId ?? snapshot.service.providerResourceId ?? undefined,
    };
  }

  if (resourceType === "backups") {
    const backup = await getBackup(id);

    if (!backup) {
      return null;
    }

    return {
      resourceType,
      id: backup.id,
      name: backup.name,
      serviceId: backup.serviceId,
      serviceNo: backup.service.serviceNo,
      customerId: backup.service.customerId,
      providerType: backup.service.providerType,
      providerResourceId:
        backup.providerBackupId ?? backup.service.providerResourceId ?? undefined,
    };
  }

  const securityGroup = await getSecurityGroup(id);

  if (!securityGroup) {
    return null;
  }

  return {
    resourceType,
    id: securityGroup.id,
    name: securityGroup.name,
    serviceId: securityGroup.serviceId ?? undefined,
    serviceNo: securityGroup.service?.serviceNo,
    customerId: securityGroup.service?.customerId,
    providerType: securityGroup.service?.providerType,
    providerResourceId:
      securityGroup.providerSecurityGroupId ??
      securityGroup.service?.providerResourceId ??
      undefined,
  };
}

export async function executeResourceAction(
  resourceType: ResourceType,
  id: string,
  action: SupportedResourceAction,
  payload: Record<string, unknown> = {},
): Promise<ResourceActionResult> {
  if (resourceType === "disks") {
    const disk = await getDisk(id);

    if (!disk) {
      throw new ResourceActionError("磁盘不存在", 404);
    }

    if (action === "createSnapshot") {
      const snapshot = await db.serviceSnapshot.create({
        data: {
          serviceId: disk.serviceId,
          sourceDiskId: disk.id,
          name: normalizeText(payload.name) ?? `${disk.name}-快照`,
          providerSnapshotId: `snap-${Date.now()}`,
          sizeGb: disk.sizeGb,
          status: "READY",
        },
      });

      const responseData = await getSnapshot(snapshot.id);
      const summary = `为磁盘 ${disk.name} 创建快照`;

      await writeResourceProviderLog({
        providerType: disk.service.providerType,
        resourceType,
        resourceId: snapshot.providerSnapshotId ?? disk.providerDiskId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: responseData,
      });

      return {
        data: responseData,
        summary,
        targetType: "serviceSnapshot",
        targetId: snapshot.id,
        serviceId: disk.serviceId,
        serviceNo: disk.service.serviceNo,
        customerId: disk.service.customerId,
      };
    }

    if (action === "createBackup") {
      const expiresAt =
        normalizeText(payload.expiresAt) !== undefined
          ? new Date(String(payload.expiresAt))
          : daysLater(30);

      const backup = await db.serviceBackup.create({
        data: {
          serviceId: disk.serviceId,
          name: normalizeText(payload.name) ?? `${disk.name}-备份`,
          providerBackupId: `bak-${Date.now()}`,
          sizeGb: disk.sizeGb,
          status: "READY",
          expiresAt: Number.isNaN(expiresAt.getTime()) ? daysLater(30) : expiresAt,
        },
      });

      const responseData = await getBackup(backup.id);
      const summary = `为磁盘 ${disk.name} 创建备份`;

      await writeResourceProviderLog({
        providerType: disk.service.providerType,
        resourceType,
        resourceId: backup.providerBackupId ?? disk.providerDiskId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: responseData,
      });

      return {
        data: responseData,
        summary,
        targetType: "serviceBackup",
        targetId: backup.id,
        serviceId: disk.serviceId,
        serviceNo: disk.service.serviceNo,
        customerId: disk.service.customerId,
      };
    }

    if (action === "detach") {
      const responseData = await db.serviceDisk.update({
        where: { id },
        data: {
          status: "DETACHED",
          mountPoint: null,
        },
        include: {
          service: {
            select: {
              id: true,
              serviceNo: true,
              customerId: true,
              providerType: true,
              providerResourceId: true,
            },
          },
        },
      });
      const summary = `卸载磁盘 ${disk.name}`;

      await writeResourceProviderLog({
        providerType: disk.service.providerType,
        resourceType,
        resourceId: disk.providerDiskId ?? disk.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: responseData,
      });

      return {
        data: responseData,
        summary,
        targetType: "serviceDisk",
        targetId: disk.id,
        serviceId: disk.serviceId,
        serviceNo: disk.service.serviceNo,
        customerId: disk.service.customerId,
      };
    }

    if (action === "attach") {
      const responseData = await db.serviceDisk.update({
        where: { id },
        data: {
          status: "ATTACHED",
          mountPoint: normalizeText(payload.mountPoint) ?? "/data",
        },
        include: {
          service: {
            select: {
              id: true,
              serviceNo: true,
              customerId: true,
              providerType: true,
              providerResourceId: true,
            },
          },
        },
      });
      const summary = `挂载磁盘 ${disk.name}`;

      await writeResourceProviderLog({
        providerType: disk.service.providerType,
        resourceType,
        resourceId: disk.providerDiskId ?? disk.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: responseData,
      });

      return {
        data: responseData,
        summary,
        targetType: "serviceDisk",
        targetId: disk.id,
        serviceId: disk.serviceId,
        serviceNo: disk.service.serviceNo,
        customerId: disk.service.customerId,
      };
    }

    if (action === "setBoot") {
      await db.serviceDisk.updateMany({
        where: {
          serviceId: disk.serviceId,
        },
        data: {
          isSystem: false,
        },
      });

      const responseData = await db.serviceDisk.update({
        where: { id },
        data: {
          isSystem: true,
          status: "ATTACHED",
        },
        include: {
          service: {
            select: {
              id: true,
              serviceNo: true,
              customerId: true,
              providerType: true,
              providerResourceId: true,
            },
          },
        },
      });
      const summary = `将磁盘 ${disk.name} 设置为系统盘`;

      await writeResourceProviderLog({
        providerType: disk.service.providerType,
        resourceType,
        resourceId: disk.providerDiskId ?? disk.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: responseData,
      });

      return {
        data: responseData,
        summary,
        targetType: "serviceDisk",
        targetId: disk.id,
        serviceId: disk.serviceId,
        serviceNo: disk.service.serviceNo,
        customerId: disk.service.customerId,
      };
    }

    throw new ResourceActionError("不支持的磁盘操作");
  }

  if (resourceType === "snapshots") {
    const snapshot = await getSnapshot(id);

    if (!snapshot) {
      throw new ResourceActionError("快照不存在", 404);
    }

    if (action === "restore") {
      await db.serviceInstance.update({
        where: {
          id: snapshot.serviceId,
        },
        data: {
          lastSyncAt: new Date(),
        },
      });

      const summary = `提交快照 ${snapshot.name} 恢复`;

      await writeResourceProviderLog({
        providerType: snapshot.service.providerType,
        resourceType,
        resourceId:
          snapshot.providerSnapshotId ?? snapshot.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: snapshot,
      });

      return {
        data: snapshot,
        summary,
        targetType: "serviceSnapshot",
        targetId: snapshot.id,
        serviceId: snapshot.serviceId,
        serviceNo: snapshot.service.serviceNo,
        customerId: snapshot.service.customerId,
      };
    }

    if (action === "deleteSnapshot") {
      await db.serviceSnapshot.delete({
        where: { id },
      });
      const summary = `删除快照 ${snapshot.name}`;

      await writeResourceProviderLog({
        providerType: snapshot.service.providerType,
        resourceType,
        resourceId:
          snapshot.providerSnapshotId ?? snapshot.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: { id: snapshot.id },
      });

      return {
        data: { id: snapshot.id },
        summary,
        targetType: "serviceSnapshot",
        targetId: snapshot.id,
        serviceId: snapshot.serviceId,
        serviceNo: snapshot.service.serviceNo,
        customerId: snapshot.service.customerId,
      };
    }

    throw new ResourceActionError("不支持的快照操作");
  }

  if (resourceType === "backups") {
    const backup = await getBackup(id);

    if (!backup) {
      throw new ResourceActionError("备份不存在", 404);
    }

    if (action === "restore") {
      await db.serviceInstance.update({
        where: {
          id: backup.serviceId,
        },
        data: {
          lastSyncAt: new Date(),
        },
      });
      const summary = `提交备份 ${backup.name} 恢复`;

      await writeResourceProviderLog({
        providerType: backup.service.providerType,
        resourceType,
        resourceId:
          backup.providerBackupId ?? backup.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: backup,
      });

      return {
        data: backup,
        summary,
        targetType: "serviceBackup",
        targetId: backup.id,
        serviceId: backup.serviceId,
        serviceNo: backup.service.serviceNo,
        customerId: backup.service.customerId,
      };
    }

    if (action === "expireNow") {
      const responseData = await db.serviceBackup.update({
        where: { id },
        data: {
          status: "ARCHIVED",
          expiresAt: new Date(),
        },
        include: {
          service: {
            select: {
              id: true,
              serviceNo: true,
              customerId: true,
              providerType: true,
              providerResourceId: true,
            },
          },
        },
      });
      const summary = `归档备份 ${backup.name}`;

      await writeResourceProviderLog({
        providerType: backup.service.providerType,
        resourceType,
        resourceId:
          backup.providerBackupId ?? backup.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: responseData,
      });

      return {
        data: responseData,
        summary,
        targetType: "serviceBackup",
        targetId: backup.id,
        serviceId: backup.serviceId,
        serviceNo: backup.service.serviceNo,
        customerId: backup.service.customerId,
      };
    }

    if (action === "deleteBackup") {
      await db.serviceBackup.delete({
        where: { id },
      });
      const summary = `删除备份 ${backup.name}`;

      await writeResourceProviderLog({
        providerType: backup.service.providerType,
        resourceType,
        resourceId:
          backup.providerBackupId ?? backup.service.providerResourceId ?? undefined,
        action,
        message: summary,
        payload,
        responseBody: { id: backup.id },
      });

      return {
        data: { id: backup.id },
        summary,
        targetType: "serviceBackup",
        targetId: backup.id,
        serviceId: backup.serviceId,
        serviceNo: backup.service.serviceNo,
        customerId: backup.service.customerId,
      };
    }

    throw new ResourceActionError("不支持的备份操作");
  }

  const securityGroup = await getSecurityGroup(id);

  if (!securityGroup) {
    throw new ResourceActionError("安全组不存在", 404);
  }

  if (action === "addRule") {
    const direction = normalizeText(payload.direction);
    const protocol = normalizeText(payload.protocol);
    const portRange = normalizeText(payload.portRange);
    const sourceCidr = normalizeText(payload.sourceCidr);

    if (!direction || !protocol || !portRange || !sourceCidr) {
      throw new ResourceActionError("新增安全组规则缺少必要参数");
    }

    const rule = await db.serviceSecurityGroupRule.create({
      data: {
        securityGroupId: id,
        direction,
        protocol,
        portRange,
        sourceCidr,
        action: normalizeText(payload.ruleAction) ?? "ALLOW",
        description: normalizeText(payload.description),
      },
    });

    const responseData = await getSecurityGroup(id);
    const summary = `为安全组 ${securityGroup.name} 新增规则`;

    await writeResourceProviderLog({
      providerType: securityGroup.service?.providerType,
      resourceType,
      resourceId:
        securityGroup.providerSecurityGroupId ??
        securityGroup.service?.providerResourceId ??
        undefined,
      action,
      message: summary,
      payload,
      responseBody: responseData,
    });

    return {
      data: responseData,
      summary,
      targetType: "serviceSecurityGroupRule",
      targetId: rule.id,
      serviceId: securityGroup.serviceId ?? undefined,
      serviceNo: securityGroup.service?.serviceNo,
      customerId: securityGroup.service?.customerId,
    };
  }

  if (action === "deleteRule") {
    const ruleId = normalizeText(payload.ruleId);

    if (!ruleId) {
      throw new ResourceActionError("缺少待删除的规则 ID");
    }

    const rule = await db.serviceSecurityGroupRule.findFirst({
      where: {
        id: ruleId,
        securityGroupId: id,
      },
    });

    if (!rule) {
      throw new ResourceActionError("安全组规则不存在", 404);
    }

    await db.serviceSecurityGroupRule.delete({
      where: { id: ruleId },
    });

    const responseData = await getSecurityGroup(id);
    const summary = `删除安全组 ${securityGroup.name} 的规则`;

    await writeResourceProviderLog({
      providerType: securityGroup.service?.providerType,
      resourceType,
      resourceId:
        securityGroup.providerSecurityGroupId ??
        securityGroup.service?.providerResourceId ??
        undefined,
      action,
      message: summary,
      payload,
      responseBody: responseData,
    });

    return {
      data: responseData,
      summary,
      targetType: "serviceSecurityGroupRule",
      targetId: rule.id,
      serviceId: securityGroup.serviceId ?? undefined,
      serviceNo: securityGroup.service?.serviceNo,
      customerId: securityGroup.service?.customerId,
    };
  }

  if (action === "deleteGroup") {
    await db.serviceSecurityGroup.delete({
      where: { id },
    });
    const summary = `删除安全组 ${securityGroup.name}`;

    await writeResourceProviderLog({
      providerType: securityGroup.service?.providerType,
      resourceType,
      resourceId:
        securityGroup.providerSecurityGroupId ??
        securityGroup.service?.providerResourceId ??
        undefined,
      action,
      message: summary,
      payload,
      responseBody: { id: securityGroup.id },
    });

    return {
      data: { id: securityGroup.id },
      summary,
      targetType: "serviceSecurityGroup",
      targetId: securityGroup.id,
      serviceId: securityGroup.serviceId ?? undefined,
      serviceNo: securityGroup.service?.serviceNo,
      customerId: securityGroup.service?.customerId,
    };
  }

  throw new ResourceActionError("不支持的安全组操作");
}
