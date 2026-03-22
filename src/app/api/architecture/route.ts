import bcrypt from "bcryptjs";
import { NextRequest, NextResponse } from "next/server";

import { writeAuditLog } from "@/lib/audit";
import {
  canManageResources,
  canViewResources,
  getApiUser,
} from "@/lib/api-auth";
import { buildCloudPlanConfig } from "@/lib/cloud-plan-config";
import { getArchitecturePageData } from "@/lib/data";
import { db } from "@/lib/db";
import {
  cloudFlavorSchema,
  cloudImageSchema,
  cloudPlanSchema,
  cloudRegionSchema,
  cloudZoneSchema,
  customerUserSchema,
} from "@/lib/validation";

type ArchitectureKind =
  | "region"
  | "zone"
  | "flavor"
  | "image"
  | "plan"
  | "portalUser";

const kindLabels: Record<ArchitectureKind, string> = {
  region: "地域",
  zone: "可用区",
  flavor: "规格",
  image: "镜像",
  plan: "售卖方案",
  portalUser: "门户账号",
};

function unauthorized() {
  return NextResponse.json({ message: "未登录或登录已失效" }, { status: 401 });
}

function forbidden() {
  return NextResponse.json({ message: "没有权限执行该操作" }, { status: 403 });
}

function badRequest(message: string, errors?: unknown) {
  return NextResponse.json({ message, errors }, { status: 400 });
}

function normalizeText(value?: string | null) {
  const normalized = value?.trim();
  return normalized ? normalized : undefined;
}

function normalizeId(value?: string | null) {
  return value && value.length > 0 ? value : undefined;
}

function moneyToCents(value: unknown) {
  return Math.round(Number(value ?? 0) * 100);
}

function normalizeOptionalNumber(value: unknown) {
  if (value === undefined || value === null || value === "") {
    return undefined;
  }

  const normalized = Number(value);
  return Number.isFinite(normalized) ? normalized : undefined;
}

async function parseKind(request: NextRequest) {
  const body = await request.json();
  const kind = String(body.kind ?? "") as ArchitectureKind;

  if (!kindLabels[kind]) {
    return { body, kind: null };
  }

  return { body, kind };
}

export async function GET(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return unauthorized();
  }

  if (!canViewResources(user)) {
    return forbidden();
  }

  const data = await getArchitecturePageData();
  return NextResponse.json({ data });
}

export async function POST(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return unauthorized();
  }

  if (!canManageResources(user)) {
    return forbidden();
  }

  const { body, kind } = await parseKind(request);

  if (!kind) {
    return badRequest("不支持的云架构对象类型");
  }

  try {
    switch (kind) {
      case "region": {
        const parsed = cloudRegionSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("地域参数不正确", parsed.error.flatten());
        }

        const region = await db.cloudRegion.create({
          data: {
            code: parsed.data.code.trim().toLowerCase(),
            name: parsed.data.name,
            location: normalizeText(parsed.data.location),
            providerType: parsed.data.providerType,
            providerRegionId: normalizeText(parsed.data.providerRegionId),
            description: normalizeText(parsed.data.description),
            isActive: parsed.data.isActive,
            sortOrder: parsed.data.sortOrder,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "create",
          targetType: "cloudRegion",
          targetId: region.id,
          summary: `创建地域 ${region.code}`,
        });

        return NextResponse.json({ data: region }, { status: 201 });
      }

      case "zone": {
        const parsed = cloudZoneSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("可用区参数不正确", parsed.error.flatten());
        }

        const zone = await db.cloudZone.create({
          data: {
            regionId: parsed.data.regionId,
            code: parsed.data.code.trim().toLowerCase(),
            name: parsed.data.name,
            providerZoneId: normalizeText(parsed.data.providerZoneId),
            isActive: parsed.data.isActive,
            sortOrder: parsed.data.sortOrder,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "create",
          targetType: "cloudZone",
          targetId: zone.id,
          summary: `创建可用区 ${zone.code}`,
        });

        return NextResponse.json({ data: zone }, { status: 201 });
      }

      case "flavor": {
        const parsed = cloudFlavorSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("规格参数不正确", parsed.error.flatten());
        }

        const flavor = await db.cloudFlavor.create({
          data: {
            code: parsed.data.code.trim().toUpperCase(),
            name: parsed.data.name,
            category: parsed.data.category,
            cpu: parsed.data.cpu,
            memoryGb: parsed.data.memoryGb,
            storageGb: parsed.data.storageGb,
            bandwidthMbps: parsed.data.bandwidthMbps,
            gpuCount: parsed.data.gpuCount,
            description: normalizeText(parsed.data.description),
            isActive: parsed.data.isActive,
            sortOrder: parsed.data.sortOrder,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "create",
          targetType: "cloudFlavor",
          targetId: flavor.id,
          summary: `创建规格 ${flavor.code}`,
        });

        return NextResponse.json({ data: flavor }, { status: 201 });
      }

      case "image": {
        const parsed = cloudImageSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("镜像参数不正确", parsed.error.flatten());
        }

        const image = await db.cloudImage.create({
          data: {
            code: parsed.data.code.trim().toLowerCase(),
            name: parsed.data.name,
            osType: parsed.data.osType,
            version: normalizeText(parsed.data.version),
            architecture: parsed.data.architecture,
            regionId: normalizeId(parsed.data.regionId),
            description: normalizeText(parsed.data.description),
            isPublic: parsed.data.isPublic,
            isActive: parsed.data.isActive,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "create",
          targetType: "cloudImage",
          targetId: image.id,
          summary: `创建镜像 ${image.code}`,
        });

        return NextResponse.json({ data: image }, { status: 201 });
      }

      case "plan": {
        const parsed = cloudPlanSchema.safeParse({
          ...body,
          salePrice: moneyToCents(body.salePrice),
          marketPrice:
            body.marketPrice === undefined || body.marketPrice === ""
              ? undefined
              : moneyToCents(body.marketPrice),
          setupFee: moneyToCents(body.setupFee),
          stock: Number(body.stock ?? 0),
          configCpu: normalizeOptionalNumber(body.configCpu),
          configMemoryGb: normalizeOptionalNumber(body.configMemoryGb),
          systemDiskSize: normalizeOptionalNumber(body.systemDiskSize),
          dataDiskSize: normalizeOptionalNumber(body.dataDiskSize),
          bandwidthMbps: normalizeOptionalNumber(body.bandwidthMbps),
          inboundBandwidthMbps: normalizeOptionalNumber(body.inboundBandwidthMbps),
          flowLimitGb: normalizeOptionalNumber(body.flowLimitGb),
          ipCount: normalizeOptionalNumber(body.ipCount),
          peakDefenceGbps: normalizeOptionalNumber(body.peakDefenceGbps),
        });

        if (!parsed.success) {
          return badRequest("售卖方案参数不正确", parsed.error.flatten());
        }

        const plan = await db.cloudPlan.create({
          data: {
            productId: parsed.data.productId,
            regionId: parsed.data.regionId,
            zoneId: normalizeId(parsed.data.zoneId),
            flavorId: normalizeId(parsed.data.flavorId),
            imageId: normalizeId(parsed.data.imageId),
            code: parsed.data.code.trim().toUpperCase(),
            name: parsed.data.name,
            billingCycle: parsed.data.billingCycle,
            salePrice: parsed.data.salePrice,
            marketPrice: parsed.data.marketPrice,
            setupFee: parsed.data.setupFee,
            stock: parsed.data.stock,
            isPublic: parsed.data.isPublic,
            isActive: parsed.data.isActive,
            configOptions: buildCloudPlanConfig(parsed.data),
            description: normalizeText(parsed.data.description),
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "create",
          targetType: "cloudPlan",
          targetId: plan.id,
          summary: `创建售卖方案 ${plan.code}`,
        });

        return NextResponse.json({ data: plan }, { status: 201 });
      }

      case "portalUser": {
        const parsed = customerUserSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("门户账号参数不正确", parsed.error.flatten());
        }

        if (!parsed.data.password) {
          return badRequest("新增门户账号时必须设置登录密码");
        }

        const portalUser = await db.customerUser.create({
          data: {
            customerId: parsed.data.customerId,
            name: parsed.data.name,
            email: parsed.data.email,
            phone: normalizeText(parsed.data.phone),
            passwordHash: await bcrypt.hash(parsed.data.password, 10),
            role: parsed.data.role,
            isOwner: parsed.data.isOwner,
            isActive: parsed.data.isActive,
          },
          include: {
            customer: true,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "create",
          targetType: "customerUser",
          targetId: portalUser.id,
          summary: `创建门户账号 ${portalUser.email}`,
        });

        return NextResponse.json({ data: portalUser }, { status: 201 });
      }
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : "创建失败";
    return badRequest(message);
  }
}

export async function PUT(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return unauthorized();
  }

  if (!canManageResources(user)) {
    return forbidden();
  }

  const { body, kind } = await parseKind(request);
  const id = String(body.id ?? "");

  if (!kind) {
    return badRequest("不支持的云架构对象类型");
  }

  if (!id) {
    return badRequest("缺少待更新记录 ID");
  }

  try {
    switch (kind) {
      case "region": {
        const parsed = cloudRegionSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("地域参数不正确", parsed.error.flatten());
        }

        const region = await db.cloudRegion.update({
          where: { id },
          data: {
            code: parsed.data.code.trim().toLowerCase(),
            name: parsed.data.name,
            location: normalizeText(parsed.data.location),
            providerType: parsed.data.providerType,
            providerRegionId: normalizeText(parsed.data.providerRegionId),
            description: normalizeText(parsed.data.description),
            isActive: parsed.data.isActive,
            sortOrder: parsed.data.sortOrder,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "update",
          targetType: "cloudRegion",
          targetId: region.id,
          summary: `更新地域 ${region.code}`,
        });

        return NextResponse.json({ data: region });
      }

      case "zone": {
        const parsed = cloudZoneSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("可用区参数不正确", parsed.error.flatten());
        }

        const zone = await db.cloudZone.update({
          where: { id },
          data: {
            regionId: parsed.data.regionId,
            code: parsed.data.code.trim().toLowerCase(),
            name: parsed.data.name,
            providerZoneId: normalizeText(parsed.data.providerZoneId),
            isActive: parsed.data.isActive,
            sortOrder: parsed.data.sortOrder,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "update",
          targetType: "cloudZone",
          targetId: zone.id,
          summary: `更新可用区 ${zone.code}`,
        });

        return NextResponse.json({ data: zone });
      }

      case "flavor": {
        const parsed = cloudFlavorSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("规格参数不正确", parsed.error.flatten());
        }

        const flavor = await db.cloudFlavor.update({
          where: { id },
          data: {
            code: parsed.data.code.trim().toUpperCase(),
            name: parsed.data.name,
            category: parsed.data.category,
            cpu: parsed.data.cpu,
            memoryGb: parsed.data.memoryGb,
            storageGb: parsed.data.storageGb,
            bandwidthMbps: parsed.data.bandwidthMbps,
            gpuCount: parsed.data.gpuCount,
            description: normalizeText(parsed.data.description),
            isActive: parsed.data.isActive,
            sortOrder: parsed.data.sortOrder,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "update",
          targetType: "cloudFlavor",
          targetId: flavor.id,
          summary: `更新规格 ${flavor.code}`,
        });

        return NextResponse.json({ data: flavor });
      }

      case "image": {
        const parsed = cloudImageSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("镜像参数不正确", parsed.error.flatten());
        }

        const image = await db.cloudImage.update({
          where: { id },
          data: {
            code: parsed.data.code.trim().toLowerCase(),
            name: parsed.data.name,
            osType: parsed.data.osType,
            version: normalizeText(parsed.data.version),
            architecture: parsed.data.architecture,
            regionId: normalizeId(parsed.data.regionId),
            description: normalizeText(parsed.data.description),
            isPublic: parsed.data.isPublic,
            isActive: parsed.data.isActive,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "update",
          targetType: "cloudImage",
          targetId: image.id,
          summary: `更新镜像 ${image.code}`,
        });

        return NextResponse.json({ data: image });
      }

      case "plan": {
        const parsed = cloudPlanSchema.safeParse({
          ...body,
          salePrice: moneyToCents(body.salePrice),
          marketPrice:
            body.marketPrice === undefined || body.marketPrice === ""
              ? undefined
              : moneyToCents(body.marketPrice),
          setupFee: moneyToCents(body.setupFee),
          stock: Number(body.stock ?? 0),
          configCpu: normalizeOptionalNumber(body.configCpu),
          configMemoryGb: normalizeOptionalNumber(body.configMemoryGb),
          systemDiskSize: normalizeOptionalNumber(body.systemDiskSize),
          dataDiskSize: normalizeOptionalNumber(body.dataDiskSize),
          bandwidthMbps: normalizeOptionalNumber(body.bandwidthMbps),
          inboundBandwidthMbps: normalizeOptionalNumber(body.inboundBandwidthMbps),
          flowLimitGb: normalizeOptionalNumber(body.flowLimitGb),
          ipCount: normalizeOptionalNumber(body.ipCount),
          peakDefenceGbps: normalizeOptionalNumber(body.peakDefenceGbps),
        });

        if (!parsed.success) {
          return badRequest("售卖方案参数不正确", parsed.error.flatten());
        }

        const plan = await db.cloudPlan.update({
          where: { id },
          data: {
            productId: parsed.data.productId,
            regionId: parsed.data.regionId,
            zoneId: normalizeId(parsed.data.zoneId),
            flavorId: normalizeId(parsed.data.flavorId),
            imageId: normalizeId(parsed.data.imageId),
            code: parsed.data.code.trim().toUpperCase(),
            name: parsed.data.name,
            billingCycle: parsed.data.billingCycle,
            salePrice: parsed.data.salePrice,
            marketPrice: parsed.data.marketPrice,
            setupFee: parsed.data.setupFee,
            stock: parsed.data.stock,
            isPublic: parsed.data.isPublic,
            isActive: parsed.data.isActive,
            configOptions: buildCloudPlanConfig(parsed.data),
            description: normalizeText(parsed.data.description),
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "update",
          targetType: "cloudPlan",
          targetId: plan.id,
          summary: `更新售卖方案 ${plan.code}`,
        });

        return NextResponse.json({ data: plan });
      }

      case "portalUser": {
        const parsed = customerUserSchema.safeParse(body);
        if (!parsed.success) {
          return badRequest("门户账号参数不正确", parsed.error.flatten());
        }

        const portalUser = await db.customerUser.update({
          where: { id },
          data: {
            customerId: parsed.data.customerId,
            name: parsed.data.name,
            email: parsed.data.email,
            phone: normalizeText(parsed.data.phone),
            role: parsed.data.role,
            isOwner: parsed.data.isOwner,
            isActive: parsed.data.isActive,
            passwordHash: parsed.data.password
              ? await bcrypt.hash(parsed.data.password, 10)
              : undefined,
          },
          include: {
            customer: true,
          },
        });

        await writeAuditLog({
          adminUserId: user.id,
          module: "architecture",
          action: "update",
          targetType: "customerUser",
          targetId: portalUser.id,
          summary: `更新门户账号 ${portalUser.email}`,
        });

        return NextResponse.json({ data: portalUser });
      }
    }
  } catch (error) {
    const message = error instanceof Error ? error.message : "更新失败";
    return badRequest(message);
  }
}

export async function DELETE(request: NextRequest) {
  const user = await getApiUser(request);

  if (!user) {
    return unauthorized();
  }

  if (!canManageResources(user)) {
    return forbidden();
  }

  const { body, kind } = await parseKind(request);
  const id = String(body.id ?? "");

  if (!kind) {
    return badRequest("不支持的云架构对象类型");
  }

  if (!id) {
    return badRequest("缺少待删除记录 ID");
  }

  try {
    switch (kind) {
      case "region":
        await db.cloudRegion.delete({ where: { id } });
        break;
      case "zone":
        await db.cloudZone.delete({ where: { id } });
        break;
      case "flavor":
        await db.cloudFlavor.delete({ where: { id } });
        break;
      case "image":
        await db.cloudImage.delete({ where: { id } });
        break;
      case "plan":
        await db.cloudPlan.delete({ where: { id } });
        break;
      case "portalUser":
        await db.customerUser.delete({ where: { id } });
        break;
    }

    await writeAuditLog({
      adminUserId: user.id,
      module: "architecture",
      action: "delete",
      targetType: kind,
      targetId: id,
      summary: `删除${kindLabels[kind]}`,
    });

    return NextResponse.json({ success: true });
  } catch (error) {
    const message = error instanceof Error ? error.message : "删除失败";
    return badRequest(message);
  }
}
