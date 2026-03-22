import bcrypt from "bcryptjs";
import {
  BillingCycle,
  OrderStatus,
  ProductCategory,
  ProductStatus,
  ProviderType,
  type Customer,
  type CustomerUser,
  type Product,
  type ServiceInstance,
  type ServiceStatus,
} from "@prisma/client";

import { db } from "@/lib/db";
import { serviceStatusLabel } from "@/lib/format";
import { queueNotification } from "@/lib/notification-service";
import { mapProviderStatus, writeProviderLog } from "@/lib/provider-sync";
import {
  type MofangCloudProvider,
  type MofangRemoteInventory,
} from "@/lib/providers/mofang-cloud";
import { getCloudProvider } from "@/lib/providers";

type PullSyncOptions = {
  limit?: number;
  includeResources?: boolean;
};

type PullSyncItem = {
  remoteId: string;
  serviceId?: string;
  serviceNo?: string;
  orderNo?: string;
  orderStatus?: string;
  customerId?: string;
  customerName?: string;
  customerEmail?: string;
  portalEmail?: string;
  operation: "created" | "updated" | "failed";
  message: string;
};

type PullSyncSummary = {
  mockMode: boolean;
  listOk: boolean;
  listMessage: string;
  remoteCount: number;
  processedCount: number;
  createdServices: number;
  updatedServices: number;
  failedServices: number;
  createdCustomers: number;
  createdProducts: number;
  createdPortalUsers: number;
  updatedPortalUsers: number;
  queuedPortalNotifications: number;
  createdImportOrders: number;
  updatedImportOrders: number;
  createdVpcs: number;
  syncedIps: number;
  syncedDisks: number;
  syncedSnapshots: number;
  syncedBackups: number;
  syncedSecurityGroups: number;
  syncedSecurityRules: number;
};

export type MofangPullSyncResult = {
  summary: PullSyncSummary;
  items: PullSyncItem[];
};

type SyncCounters = Omit<
  PullSyncSummary,
  "mockMode" | "listOk" | "listMessage" | "remoteCount"
>;

type NormalizedIp = {
  address: string;
  version: string;
  type: string;
  bandwidthMbps?: number;
  providerIpId?: string;
  isPrimary: boolean;
  status: string;
};

type NormalizedDisk = {
  name: string;
  type: string;
  sizeGb: number;
  mountPoint?: string;
  providerDiskId?: string;
  isSystem: boolean;
  status: string;
};

type NormalizedSnapshot = {
  name: string;
  providerSnapshotId?: string;
  sizeGb?: number;
  status: string;
  sourceDiskKey?: string;
};

type NormalizedBackup = {
  name: string;
  providerBackupId?: string;
  sizeGb?: number;
  status: string;
  expiresAt?: Date;
};

type NormalizedSecurityRule = {
  direction: string;
  protocol: string;
  portRange: string;
  sourceCidr: string;
  action: string;
  description?: string;
};

type NormalizedSecurityGroup = {
  name: string;
  providerSecurityGroupId?: string;
  status: string;
  rules: NormalizedSecurityRule[];
};

type NormalizedVpc = {
  name: string;
  region?: string;
  cidr: string;
  gateway?: string;
  providerVpcId?: string;
  status: string;
};

function isRecord(value: unknown): value is Record<string, unknown> {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value);
}

function toRecord(value: unknown) {
  return isRecord(value) ? value : undefined;
}

function toRecordArray(value: unknown) {
  if (Array.isArray(value)) {
    return value.filter(isRecord);
  }

  if (isRecord(value)) {
    return [value];
  }

  return [] as Record<string, unknown>[];
}

function pickString(record: Record<string, unknown>, keys: string[]) {
  for (const key of keys) {
    const value = record[key];

    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }

    if (typeof value === "number" && Number.isFinite(value)) {
      return String(value);
    }
  }

  return undefined;
}

function pickNumber(record: Record<string, unknown>, keys: string[]) {
  for (const key of keys) {
    const value = record[key];

    if (typeof value === "number" && Number.isFinite(value)) {
      return value;
    }

    if (typeof value === "string" && value.trim()) {
      const normalized = Number(value.replace(/[^\d.-]/g, ""));

      if (Number.isFinite(normalized)) {
        return normalized;
      }
    }
  }

  return undefined;
}

function pickBoolean(record: Record<string, unknown>, keys: string[]) {
  for (const key of keys) {
    const value = record[key];

    if (typeof value === "boolean") {
      return value;
    }

    if (typeof value === "number") {
      return value !== 0;
    }

    if (typeof value === "string") {
      const normalized = value.trim().toLowerCase();

      if (["1", "true", "yes", "y", "on"].includes(normalized)) {
        return true;
      }

      if (["0", "false", "no", "n", "off"].includes(normalized)) {
        return false;
      }
    }
  }

  return undefined;
}

function pickDate(record: Record<string, unknown>, keys: string[]) {
  for (const key of keys) {
    const value = record[key];

    if (value instanceof Date && !Number.isNaN(value.getTime())) {
      return value;
    }

    if (typeof value === "number" && Number.isFinite(value)) {
      const normalized = value > 1_000_000_000_000 ? value : value * 1000;
      const parsed = new Date(normalized);

      if (!Number.isNaN(parsed.getTime())) {
        return parsed;
      }
    }

    if (typeof value === "string" && value.trim()) {
      const numeric = Number(value);

      if (Number.isFinite(numeric)) {
        const normalized = numeric > 1_000_000_000_000 ? numeric : numeric * 1000;
        const parsed = new Date(normalized);

        if (!Number.isNaN(parsed.getTime())) {
          return parsed;
        }
      }

      const parsed = new Date(value);

      if (!Number.isNaN(parsed.getTime())) {
        return parsed;
      }
    }
  }

  return undefined;
}

function pickNestedRecord(record: Record<string, unknown>, keys: string[]) {
  for (const key of keys) {
    const value = record[key];

    if (isRecord(value)) {
      return value;
    }
  }

  return undefined;
}

function sanitizeToken(value: string, maxLength = 24) {
  const normalized = value
    .normalize("NFKD")
    .replace(/[^\w-]+/g, "-")
    .replace(/_+/g, "-")
    .replace(/-+/g, "-")
    .replace(/^-|-$/g, "")
    .toUpperCase();

  return normalized.slice(0, maxLength) || "SYNC";
}

function normalizeMoneyToCents(value?: number | string | null) {
  if (value === undefined || value === null) {
    return 0;
  }

  const raw =
    typeof value === "string" ? Number(value.replace(/[^\d.-]/g, "")) : value;

  if (!Number.isFinite(raw)) {
    return 0;
  }

  if (
    Math.abs(raw) >= 10000 &&
    ((typeof value === "number" && Number.isInteger(value)) ||
      (typeof value === "string" && !value.includes(".")))
  ) {
    return Math.round(raw);
  }

  return Math.round(raw * 100);
}

function normalizeMemoryGb(value?: number) {
  if (!value || !Number.isFinite(value)) {
    return undefined;
  }

  if (value > 256) {
    return Math.max(1, Math.round(value / 1024));
  }

  return Math.max(1, Math.round(value));
}

function normalizeStorageGb(value?: number) {
  if (!value || !Number.isFinite(value)) {
    return undefined;
  }

  if (value > 16384) {
    return Math.max(1, Math.round(value / 1024));
  }

  return Math.max(1, Math.round(value));
}

function normalizeBillingCycle(
  value?: string | number | null,
  fallback: BillingCycle = BillingCycle.MONTHLY,
) {
  if (value === null || value === undefined) {
    return fallback;
  }

  const normalized = String(value).trim().toLowerCase();

  const map: Record<string, BillingCycle> = {
    monthly: BillingCycle.MONTHLY,
    month: BillingCycle.MONTHLY,
    "1": BillingCycle.MONTHLY,
    quarterly: BillingCycle.QUARTERLY,
    quarter: BillingCycle.QUARTERLY,
    "3": BillingCycle.QUARTERLY,
    semi_annually: BillingCycle.SEMI_ANNUALLY,
    semiannually: BillingCycle.SEMI_ANNUALLY,
    semiannual: BillingCycle.SEMI_ANNUALLY,
    half_year: BillingCycle.SEMI_ANNUALLY,
    "6": BillingCycle.SEMI_ANNUALLY,
    annually: BillingCycle.ANNUALLY,
    annual: BillingCycle.ANNUALLY,
    yearly: BillingCycle.ANNUALLY,
    year: BillingCycle.ANNUALLY,
    "12": BillingCycle.ANNUALLY,
    biennially: BillingCycle.BIENNIALLY,
    biennial: BillingCycle.BIENNIALLY,
    "24": BillingCycle.BIENNIALLY,
    triennially: BillingCycle.TRIENNIALLY,
    triennial: BillingCycle.TRIENNIALLY,
    "36": BillingCycle.TRIENNIALLY,
    onetime: BillingCycle.ONETIME,
    one_time: BillingCycle.ONETIME,
    "one-time": BillingCycle.ONETIME,
    once: BillingCycle.ONETIME,
  };

  return map[normalized] ?? fallback;
}

function getRemoteId(record: Record<string, unknown>) {
  return pickString(record, [
    "id",
    "cloudid",
    "cloud_id",
    "instance_id",
    "hostid",
    "server_id",
    "uuid",
  ]);
}

function getRemoteServiceNo(record: Record<string, unknown>, remoteId: string) {
  const preferred = pickString(record, [
    "service_no",
    "serviceNo",
    "hostid",
    "display_id",
    "server_no",
  ]);

  if (preferred) {
    return preferred;
  }

  return `SRV-MF-${sanitizeToken(remoteId, 16)}`;
}

function getRemoteName(record: Record<string, unknown>, remoteId: string) {
  return (
    pickString(record, [
      "name",
      "remarks",
      "remark",
      "title",
      "instance_name",
      "host_name",
      "hostname",
    ]) ?? `MF Cloud ${remoteId}`
  );
}

function getRemoteRegion(record: Record<string, unknown>) {
  const nestedRegion = pickNestedRecord(record, [
    "region_info",
    "region",
    "area_info",
    "area",
  ]);

  return (
    (nestedRegion
      ? pickString(nestedRegion, ["name", "region_name", "code"])
      : undefined) ??
    pickString(record, ["region_name", "region", "area", "zone_name", "node_name"])
  );
}

function getRemoteStatus(record: Record<string, unknown>) {
  return pickString(record, [
    "status",
    "instance_status",
    "cloud_status",
    "power_status",
    "state",
  ]);
}

function getRemotePrimaryIp(record: Record<string, unknown>) {
  return pickString(record, [
    "ip",
    "ip_address",
    "main_ip",
    "mainip",
    "public_ip",
    "primary_ip",
  ]);
}

function getRemoteCpu(record: Record<string, unknown>) {
  return pickNumber(record, ["cpu", "cpu_core", "cpu_cores", "vcpu", "core"]);
}

function getRemoteMemory(record: Record<string, unknown>) {
  return normalizeMemoryGb(
    pickNumber(record, [
      "memory",
      "memory_gb",
      "memory_mb",
      "mem",
      "ram",
      "ram_mb",
    ]),
  );
}

function getRemoteStorage(record: Record<string, unknown>, disks: NormalizedDisk[]) {
  if (disks.length > 0) {
    return disks.reduce((total, disk) => total + disk.sizeGb, 0);
  }

  return normalizeStorageGb(
    pickNumber(record, [
      "storage",
      "storage_gb",
      "disk_size",
      "system_disk_size",
      "data_disk_size",
    ]),
  );
}

function getRemoteNextDueDate(record: Record<string, unknown>) {
  return pickDate(record, [
    "next_due_date",
    "due_date",
    "renew_time",
    "expire_time",
    "expired_at",
    "end_time",
  ]);
}

function getRemoteActivatedAt(record: Record<string, unknown>) {
  return pickDate(record, [
    "active_time",
    "activated_at",
    "create_time",
    "created_at",
    "start_time",
  ]);
}

function getRemoteExpiresAt(record: Record<string, unknown>) {
  return pickDate(record, ["expire_time", "expired_at", "end_time"]);
}

async function buildUniqueCustomerNo(base: string) {
  let candidate = base;
  let attempt = 1;

  while (await db.customer.findUnique({ where: { customerNo: candidate } })) {
    candidate = `${base}-${attempt}`;
    attempt += 1;
  }

  return candidate;
}

async function buildUniqueEmail(base: string) {
  let candidate = base;
  let attempt = 1;

  while (await db.customer.findUnique({ where: { email: candidate } })) {
    const [local] = base.split("@");
    candidate = `${local}-${attempt}@sync.local`;
    attempt += 1;
  }

  return candidate;
}

async function buildUniquePortalEmail(base: string) {
  let candidate = base;
  let attempt = 1;

  while (await db.customerUser.findUnique({ where: { email: candidate } })) {
    const [local, domain = "sync.local"] = base.split("@");
    candidate = `${local}-${attempt}@${domain}`;
    attempt += 1;
  }

  return candidate;
}

function isSyncPlaceholderEmail(email: string) {
  return email.toLowerCase().endsWith("@sync.local");
}

function getMofangSyncPortalDefaultPassword() {
  return process.env.MOFANG_SYNC_PORTAL_DEFAULT_PASSWORD ?? "Portal123!";
}

async function resolvePortalOwnerEmail(
  desiredEmail: string,
  existingUserId?: string,
) {
  const normalized = desiredEmail.toLowerCase();
  const existingByEmail = await db.customerUser.findUnique({
    where: { email: normalized },
  });

  if (!existingByEmail || existingByEmail.id === existingUserId) {
    return normalized;
  }

  return buildUniquePortalEmail(normalized);
}

async function buildUniqueProductCode(base: string) {
  let candidate = base;
  let attempt = 1;

  while (await db.product.findUnique({ where: { code: candidate } })) {
    candidate = `${base}-${attempt}`;
    attempt += 1;
  }

  return candidate;
}

async function buildUniqueServiceNo(base: string) {
  let candidate = base;
  let attempt = 1;

  while (await db.serviceInstance.findUnique({ where: { serviceNo: candidate } })) {
    candidate = `${base}-${attempt}`;
    attempt += 1;
  }

  return candidate;
}

async function buildUniqueOrderNo(base: string) {
  let candidate = base;
  let attempt = 1;

  while (await db.order.findUnique({ where: { orderNo: candidate } })) {
    candidate = `${base}-${attempt}`;
    attempt += 1;
  }

  return candidate;
}

async function ensureCustomer(
  remote: Record<string, unknown>,
  remoteId: string,
  counters: SyncCounters,
) {
  const customerRecord =
    pickNestedRecord(remote, ["customer", "client", "user", "owner", "member"]) ??
    remote;
  const customerKey =
    pickString(customerRecord, ["customer_no", "client_no", "client_id", "uid", "user_id"]) ??
    remoteId;
  const customerName =
    pickString(customerRecord, [
      "name",
      "company_name",
      "customer_name",
      "username",
      "nick_name",
    ]) ?? `MF Customer ${sanitizeToken(customerKey, 10)}`;
  const rawEmail =
    pickString(customerRecord, ["email", "mail", "user_email"]) ??
    `mf-${sanitizeToken(customerKey, 18).toLowerCase()}@sync.local`;
  const normalizedRawEmail = rawEmail.toLowerCase();
  const expectedCustomerNo = `CUS-MF-${sanitizeToken(customerKey, 12)}`;

  const existing =
    (await db.customer.findFirst({
      where: {
        OR: [{ email: normalizedRawEmail }, { customerNo: expectedCustomerNo }],
      },
    })) ?? undefined;

  if (existing) {
    let nextEmail = existing.email;

    if (
      existing.email !== normalizedRawEmail &&
      isSyncPlaceholderEmail(existing.email) &&
      !isSyncPlaceholderEmail(normalizedRawEmail)
    ) {
      const emailOwner = await db.customer.findUnique({
        where: { email: normalizedRawEmail },
      });

      if (!emailOwner || emailOwner.id === existing.id) {
        nextEmail = normalizedRawEmail;
      }
    }

    return db.customer.update({
      where: { id: existing.id },
      data: {
        name: customerName,
        companyName:
          pickString(customerRecord, ["company_name", "company"]) ?? existing.companyName,
        email: nextEmail,
        phone: pickString(customerRecord, ["phone", "mobile", "tel"]) ?? existing.phone,
        notes: "Synced from Mofang Cloud pull sync",
      },
    });
  }

  const email = await buildUniqueEmail(normalizedRawEmail);
  const customerNo = await buildUniqueCustomerNo(expectedCustomerNo);

  counters.createdCustomers += 1;

  return db.customer.create({
    data: {
      customerNo,
      name: customerName,
      companyName: pickString(customerRecord, ["company_name", "company"]),
      email,
      phone: pickString(customerRecord, ["phone", "mobile", "tel"]),
      type: "COMPANY",
      status: "ACTIVE",
      level: "standard",
      notes: "Synced from Mofang Cloud pull sync",
    },
  });
}

async function ensurePortalOwner(
  customer: Customer,
  remote: Record<string, unknown>,
  counters: SyncCounters,
) {
  const customerRecord =
    pickNestedRecord(remote, ["customer", "client", "user", "owner", "member"]) ??
    remote;
  const existingOwner = await db.customerUser.findFirst({
    where: {
      customerId: customer.id,
      isOwner: true,
    },
    orderBy: {
      createdAt: "asc",
    },
  });

  const desiredName =
    pickString(customerRecord, [
      "name",
      "real_name",
      "customer_name",
      "company_name",
      "username",
    ]) ?? customer.name;
  const desiredPhone =
    pickString(customerRecord, ["phone", "mobile", "tel"]) ??
    customer.phone ??
    existingOwner?.phone ??
    undefined;
  const baseEmail =
    existingOwner?.email &&
    !isSyncPlaceholderEmail(existingOwner.email) &&
    isSyncPlaceholderEmail(customer.email)
      ? existingOwner.email
      : customer.email;
  const desiredEmail = await resolvePortalOwnerEmail(baseEmail, existingOwner?.id);

  if (existingOwner) {
    const updated = await db.customerUser.update({
      where: { id: existingOwner.id },
      data: {
        name: desiredName,
        email: desiredEmail,
        phone: desiredPhone,
        isActive: true,
      },
    });

    if (
      updated.name !== existingOwner.name ||
      updated.email !== existingOwner.email ||
      updated.phone !== existingOwner.phone ||
      !existingOwner.isActive
    ) {
      counters.updatedPortalUsers += 1;
    }

    return updated;
  }

  const passwordHash = await bcrypt.hash(getMofangSyncPortalDefaultPassword(), 10);

  counters.createdPortalUsers += 1;

  return db.customerUser.create({
    data: {
      customerId: customer.id,
      name: desiredName,
      email: desiredEmail,
      phone: desiredPhone,
      passwordHash,
      role: "OWNER",
      isOwner: true,
      isActive: true,
    },
  });
}

function mapImportedOrderStatus(serviceStatus: ServiceStatus): OrderStatus {
  if (serviceStatus === "PENDING") {
    return OrderStatus.PENDING;
  }

  if (serviceStatus === "PROVISIONING") {
    return OrderStatus.PROVISIONING;
  }

  if (serviceStatus === "TERMINATED") {
    return OrderStatus.CANCELLED;
  }

  if (serviceStatus === "FAILED") {
    return OrderStatus.CANCELLED;
  }

  return OrderStatus.ACTIVE;
}

function buildImportedOrderNotes(service: ServiceInstance) {
  const details = [
    `导入自魔方云实例 ${service.serviceNo}`,
    service.providerResourceId ? `远端实例 ID: ${service.providerResourceId}` : undefined,
    `当前状态: ${serviceStatusLabel(service.status)}`,
  ].filter(Boolean);

  return details.join(" / ");
}

async function ensureImportedOrder(input: {
  customer: Customer;
  product: Product;
  service: ServiceInstance;
  counters: SyncCounters;
}) {
  const existingOrder =
    (input.service.orderId
      ? await db.order.findUnique({
          where: { id: input.service.orderId },
        })
      : null) ??
    (await db.orderItem.findFirst({
      where: {
        serviceId: input.service.id,
      },
      include: {
        order: true,
      },
      orderBy: {
        createdAt: "asc",
      },
    }).then((item) => item?.order ?? null));

  const totalAmount = Math.max(input.service.monthlyCost, 0);
  const orderStatus = mapImportedOrderStatus(input.service.status);
  const paidAmount =
    orderStatus === OrderStatus.PENDING || orderStatus === OrderStatus.PROVISIONING
      ? 0
      : totalAmount;
  const dueDate = input.service.nextDueDate ?? input.service.createdAt;
  const paidAt =
    paidAmount > 0 ? input.service.activatedAt ?? input.service.createdAt : undefined;
  const notes = buildImportedOrderNotes(input.service);
  const orderData = {
    customerId: input.customer.id,
    status: orderStatus,
    totalAmount,
    paidAmount,
    source: "mofang-sync",
    orderType: "import",
    dueDate,
    paidAt,
    notes,
  } as const;

  const order =
    existingOrder && existingOrder.source !== "mofang-sync"
      ? existingOrder
      : existingOrder
        ? await db.order.update({
            where: { id: existingOrder.id },
            data: orderData,
          })
        : await db.order.create({
            data: {
              orderNo: await buildUniqueOrderNo(
                `ORD-MF-${sanitizeToken(
                  input.service.providerResourceId ?? input.service.serviceNo,
                  16,
                )}`,
              ),
              ...orderData,
            },
          });

  if (!existingOrder) {
    input.counters.createdImportOrders += 1;
  } else if (existingOrder.source === "mofang-sync") {
    input.counters.updatedImportOrders += 1;
  }

  if (input.service.orderId !== order.id) {
    await db.serviceInstance.update({
      where: { id: input.service.id },
      data: {
        orderId: order.id,
      },
    });
  }

  const itemPayload = {
    orderId: order.id,
    productId: input.product.id,
    serviceId: input.service.id,
    title: `${input.product.name} / ${input.service.name}`,
    quantity: 1,
    unitPrice: totalAmount,
    cycle: input.service.billingCycle,
    totalAmount,
    configSnapshot: JSON.stringify(
      {
        source: "mofang-pull-sync",
        imported: true,
        remoteId: input.service.providerResourceId,
        serviceNo: input.service.serviceNo,
      },
      null,
      2,
    ),
  } as const;
  const existingItem = await db.orderItem.findFirst({
    where: {
      orderId: order.id,
      serviceId: input.service.id,
    },
    orderBy: {
      createdAt: "asc",
    },
  });

  if (existingItem) {
    await db.orderItem.update({
      where: { id: existingItem.id },
      data: itemPayload,
    });
  } else {
    await db.orderItem.create({
      data: itemPayload,
    });
  }

  return order;
}

async function queuePortalSyncNotification(input: {
  customer: Customer;
  portalUser: CustomerUser;
  service: ServiceInstance;
  operation: "created" | "updated";
  previousStatus?: ServiceStatus | null;
  counters: SyncCounters;
}) {
  const shouldNotify =
    input.operation === "created" || input.previousStatus !== input.service.status;

  if (!shouldNotify) {
    return;
  }

  const subject =
    input.operation === "created"
      ? `???????${input.service.name}`
      : `????????${input.service.name}`;
  const content =
    input.operation === "created"
      ? `?? ${input.service.serviceNo} ?????????????????? ${serviceStatusLabel(input.service.status)}?`
      : `?? ${input.service.serviceNo} ?????? ${serviceStatusLabel(input.service.status)}?`;

  await queueNotification({
    customerId: input.customer.id,
    recipient: input.portalUser.email,
    recipientName: input.portalUser.name,
    channel: "SYSTEM",
    module: "service",
    relatedType: "service",
    relatedId: input.service.id,
    subject,
    content,
    payload: {
      source: "mofang-pull-sync",
      serviceNo: input.service.serviceNo,
      status: input.service.status,
    },
  });

  input.counters.queuedPortalNotifications += 1;
}
async function ensureProduct(
  remote: Record<string, unknown>,
  remoteId: string,
  counters: SyncCounters,
) {
  const productRecord =
    pickNestedRecord(remote, ["product", "plan", "package", "goods", "config"]) ??
    remote;
  const providerProductId = pickString(productRecord, [
    "product_id",
    "pid",
    "productid",
    "config_id",
  ]);
  const productName =
    pickString(productRecord, [
      "product_name",
      "name",
      "plan_name",
      "package_name",
      "goods_name",
    ]) ??
    (pickString(remote, ["type", "network_type"])
      ? `MF Cloud ${sanitizeToken(
          pickString(remote, ["type", "network_type"]) ?? "HOST",
          10,
        )}`
      : "MF Cloud Host");
  const billingCycle = normalizeBillingCycle(
    pickString(remote, ["billing_cycle", "cycle", "duration", "pay_type"]),
  );
  const price = normalizeMoneyToCents(
    pickNumber(remote, [
      "renew_amount",
      "renew_price",
      "price",
      "amount",
      "monthly_price",
    ]),
  );

  if (providerProductId) {
    const existingByProviderId = await db.product.findFirst({
      where: {
        providerType: ProviderType.MOFANG_CLOUD,
        providerProductId,
      },
    });

    if (existingByProviderId) {
      return db.product.update({
        where: { id: existingByProviderId.id },
        data: {
          name: productName,
          billingCycle,
          price,
          status: ProductStatus.ACTIVE,
          description: "Synced from Mofang Cloud pull sync",
        },
      });
    }
  }

  const existingByName = await db.product.findFirst({
    where: {
      providerType: ProviderType.MOFANG_CLOUD,
      name: productName,
    },
  });

  if (existingByName) {
    return db.product.update({
      where: { id: existingByName.id },
      data: {
        billingCycle,
        price,
        providerProductId: providerProductId ?? existingByName.providerProductId,
        status: ProductStatus.ACTIVE,
        description: "Synced from Mofang Cloud pull sync",
      },
    });
  }

  const code = await buildUniqueProductCode(
    `MF-${sanitizeToken(providerProductId ?? productName ?? remoteId, 20)}`,
  );

  counters.createdProducts += 1;

  return db.product.create({
    data: {
      code,
      name: productName,
      category: ProductCategory.CLOUD_SERVER,
      status: ProductStatus.ACTIVE,
      billingCycle,
      price,
      stock: 9999,
      autoProvision: true,
      providerType: ProviderType.MOFANG_CLOUD,
      providerProductId: providerProductId ?? undefined,
      description: "Synced from Mofang Cloud pull sync",
    },
  });
}

function mergeRemoteRecord(
  listRecord?: Record<string, unknown>,
  detailRecord?: Record<string, unknown>,
) {
  return {
    ...(listRecord ?? {}),
    ...(detailRecord ?? {}),
    customer:
      pickNestedRecord(detailRecord ?? {}, ["customer", "client", "user", "owner", "member"]) ??
      pickNestedRecord(listRecord ?? {}, ["customer", "client", "user", "owner", "member"]),
    product:
      pickNestedRecord(detailRecord ?? {}, ["product", "plan", "package", "goods", "config"]) ??
      pickNestedRecord(listRecord ?? {}, ["product", "plan", "package", "goods", "config"]),
  };
}

function normalizeIpVersion(record: Record<string, unknown>) {
  const version =
    pickString(record, ["version", "ip_version", "protocol"]) ??
    String(pickNumber(record, ["version", "ip_version"]) ?? "");

  if (version.includes("6")) {
    return "IPv6";
  }

  return "IPv4";
}

function normalizeIps(
  records: Record<string, unknown>[],
  fallbackIpAddress?: string,
) {
  const normalized: NormalizedIp[] = [];
  const seen = new Set<string>();

  for (const [index, record] of records.entries()) {
    const address = pickString(record, [
      "ip",
      "ip_address",
      "address",
      "public_ip",
      "private_ip",
      "value",
    ]);

    if (!address || seen.has(address)) {
      continue;
    }

    seen.add(address);

    const isPrimary =
      pickBoolean(record, ["is_primary", "primary", "main", "default"]) ??
      (fallbackIpAddress ? address === fallbackIpAddress : index === 0);

    normalized.push({
      address,
      version: normalizeIpVersion(record),
      type:
        pickString(record, ["type", "network_type", "purpose"]) ??
        (isPrimary ? "primary" : "secondary"),
      bandwidthMbps: pickNumber(record, [
        "bandwidth",
        "bandwidth_mbps",
        "bw",
        "rate",
      ]),
      providerIpId: pickString(record, ["id", "ip_id", "ipid", "uuid"]),
      isPrimary,
      status: pickString(record, ["status", "state"]) ?? "ACTIVE",
    });
  }

  if (normalized.length === 0 && fallbackIpAddress) {
    normalized.push({
      address: fallbackIpAddress,
      version: "IPv4",
      type: "primary",
      isPrimary: true,
      status: "ACTIVE",
    });
  }

  return normalized;
}

function normalizeDiskSize(record: Record<string, unknown>) {
  return (
    normalizeStorageGb(
      pickNumber(record, [
        "size",
        "size_gb",
        "capacity",
        "disk_size",
        "volume_size",
      ]),
    ) ?? 20
  );
}

function normalizeDisks(
  records: Record<string, unknown>[],
  remote: Record<string, unknown>,
) {
  const normalized: NormalizedDisk[] = [];
  const seen = new Set<string>();

  for (const [index, record] of records.entries()) {
    const providerDiskId = pickString(record, ["id", "disk_id", "volume_id", "uuid"]);
    const mountPoint = pickString(record, ["mount_point", "mount", "path", "device"]);
    const rawType = pickString(record, ["type", "disk_type", "category", "storage_type"]);
    const isSystem =
      pickBoolean(record, ["is_system", "system", "bootable", "boot"]) ??
      Boolean(rawType?.toLowerCase().includes("system")) ??
      ["/", "c:", "c:\\"].includes((mountPoint ?? "").toLowerCase());
    const key = providerDiskId ?? `${mountPoint ?? "disk"}-${index}`;

    if (seen.has(key)) {
      continue;
    }

    seen.add(key);

    normalized.push({
      name:
        pickString(record, ["name", "disk_name", "volume_name"]) ??
        (isSystem ? "System Disk" : `Data Disk ${index + 1}`),
      type: rawType ?? (isSystem ? "system" : "data"),
      sizeGb: normalizeDiskSize(record),
      mountPoint,
      providerDiskId,
      isSystem,
      status: pickString(record, ["status", "state"]) ?? "ATTACHED",
    });
  }

  if (normalized.length === 0) {
    const fallbackSize =
      normalizeStorageGb(
        pickNumber(remote, [
          "system_disk_size",
          "disk_size",
          "storage",
          "storage_gb",
        ]),
      ) ?? 40;

    normalized.push({
      name: "System Disk",
      type: "system",
      sizeGb: fallbackSize,
      mountPoint: "/",
      isSystem: true,
      status: "ATTACHED",
    });
  }

  return normalized;
}

function normalizeSnapshots(records: Record<string, unknown>[]) {
  const normalized: NormalizedSnapshot[] = [];
  const seen = new Set<string>();

  for (const [index, record] of records.entries()) {
    const providerSnapshotId = pickString(record, [
      "id",
      "snapshot_id",
      "snap_id",
      "uuid",
    ]);
    const key = providerSnapshotId ?? `${pickString(record, ["name"]) ?? "snapshot"}-${index}`;

    if (seen.has(key)) {
      continue;
    }

    seen.add(key);

    normalized.push({
      name: pickString(record, ["name", "snapshot_name"]) ?? `Snapshot ${index + 1}`,
      providerSnapshotId,
      sizeGb: normalizeStorageGb(
        pickNumber(record, ["size", "size_gb", "snapshot_size"]),
      ),
      status: pickString(record, ["status", "state"]) ?? "READY",
      sourceDiskKey: pickString(record, [
        "source_disk_id",
        "disk_id",
        "source_disk",
        "disk_name",
        "disk",
      ]),
    });
  }

  return normalized;
}

function normalizeBackups(records: Record<string, unknown>[]) {
  const normalized: NormalizedBackup[] = [];
  const seen = new Set<string>();

  for (const [index, record] of records.entries()) {
    const providerBackupId = pickString(record, [
      "id",
      "backup_id",
      "bak_id",
      "uuid",
    ]);
    const key = providerBackupId ?? `${pickString(record, ["name"]) ?? "backup"}-${index}`;

    if (seen.has(key)) {
      continue;
    }

    seen.add(key);

    normalized.push({
      name: pickString(record, ["name", "backup_name"]) ?? `Backup ${index + 1}`,
      providerBackupId,
      sizeGb: normalizeStorageGb(
        pickNumber(record, ["size", "size_gb", "backup_size"]),
      ),
      status: pickString(record, ["status", "state"]) ?? "READY",
      expiresAt: pickDate(record, ["expires_at", "expire_time", "end_time"]),
    });
  }

  return normalized;
}

function normalizeSecurityRules(
  record: Record<string, unknown>,
  fallbackDirection?: string,
) {
  const sources = [
    ...toRecordArray(record.rules),
    ...toRecordArray(record.rule),
    ...toRecordArray(record.rule_list),
    ...toRecordArray(record.items),
    ...toRecordArray(record.ingress_rules).map((rule) => ({
      ...rule,
      direction: "ingress",
    })),
    ...toRecordArray(record.egress_rules).map((rule) => ({
      ...rule,
      direction: "egress",
    })),
  ];

  return sources.map((rule) => ({
    direction:
      pickString(rule, ["direction", "type"]) ?? fallbackDirection ?? "ingress",
    protocol: pickString(rule, ["protocol", "protocol_name", "ip_protocol"]) ?? "tcp",
    portRange:
      pickString(rule, ["port_range", "ports", "port", "portRange"]) ?? "all",
    sourceCidr:
      pickString(rule, ["source", "cidr", "source_cidr", "address"]) ?? "0.0.0.0/0",
    action: pickString(rule, ["action", "policy"]) ?? "ALLOW",
    description: pickString(rule, ["description", "remark", "name"]),
  }));
}

function normalizeSecurityGroups(records: Record<string, unknown>[]) {
  const normalized: NormalizedSecurityGroup[] = [];
  const seen = new Set<string>();

  for (const [index, record] of records.entries()) {
    const providerSecurityGroupId = pickString(record, [
      "id",
      "security_group_id",
      "group_id",
      "uuid",
    ]);
    const name =
      pickString(record, ["name", "security_group_name", "group_name"]) ??
      `Security Group ${index + 1}`;
    const key = providerSecurityGroupId ?? `${name}-${index}`;

    if (seen.has(key)) {
      continue;
    }

    seen.add(key);

    normalized.push({
      name,
      providerSecurityGroupId,
      status: pickString(record, ["status", "state"]) ?? "ACTIVE",
      rules: normalizeSecurityRules(record),
    });
  }

  return normalized;
}

function normalizeVpc(record?: Record<string, unknown> | null) {
  if (!record) {
    return null;
  }

  const providerVpcId = pickString(record, ["id", "vpc_id", "network_id", "uuid"]);
  const cidr = pickString(record, ["cidr", "cidr_block", "network", "subnet"]);

  if (!cidr && !providerVpcId) {
    return null;
  }

  return {
    name:
      pickString(record, ["name", "vpc_name", "network_name"]) ??
      `VPC ${sanitizeToken(providerVpcId ?? cidr ?? "SYNC", 10)}`,
    region: getRemoteRegion(record),
    cidr: cidr ?? "10.0.0.0/16",
    gateway: pickString(record, ["gateway", "gateway_ip"]),
    providerVpcId,
    status: pickString(record, ["status", "state"]) ?? "ACTIVE",
  } satisfies NormalizedVpc;
}

function hasResourceHints(remote: Record<string, unknown>) {
  const hintKeys = [
    "ips",
    "ip_list",
    "disk",
    "disk_list",
    "snapshot",
    "snapshot_list",
    "backup",
    "backup_list",
    "security_group",
    "security_groups",
    "vpc",
    "vpc_network",
  ];

  return hintKeys.some((key) => key in remote);
}

async function ensureVpc(normalized: NormalizedVpc, counters: SyncCounters) {
  const existing =
    (normalized.providerVpcId
      ? await db.serviceVpcNetwork.findFirst({
          where: { providerVpcId: normalized.providerVpcId },
        })
      : null) ??
    (await db.serviceVpcNetwork.findFirst({
      where: {
        name: normalized.name,
        cidr: normalized.cidr,
        region: normalized.region ?? undefined,
      },
    }));

  if (existing) {
    return db.serviceVpcNetwork.update({
      where: { id: existing.id },
      data: {
        name: normalized.name,
        region: normalized.region,
        cidr: normalized.cidr,
        gateway: normalized.gateway,
        providerVpcId: normalized.providerVpcId ?? existing.providerVpcId,
        status: normalized.status,
      },
    });
  }

  counters.createdVpcs += 1;

  return db.serviceVpcNetwork.create({
    data: {
      name: normalized.name,
      region: normalized.region,
      cidr: normalized.cidr,
      gateway: normalized.gateway,
      providerVpcId: normalized.providerVpcId,
      status: normalized.status,
    },
  });
}

async function syncServiceResources(
  service: ServiceInstance,
  remote: Record<string, unknown>,
  inventory: MofangRemoteInventory,
  counters: SyncCounters,
) {
  if (
    inventory.ips.length === 0 &&
    inventory.disks.length === 0 &&
    inventory.snapshots.length === 0 &&
    inventory.backups.length === 0 &&
    inventory.securityGroups.length === 0 &&
    !inventory.vpc &&
    !hasResourceHints(remote)
  ) {
    return;
  }

  const normalizedIps = normalizeIps(inventory.ips, getRemotePrimaryIp(remote));
  const normalizedDisks = normalizeDisks(inventory.disks, remote);
  const normalizedSnapshots = normalizeSnapshots(inventory.snapshots);
  const normalizedBackups = normalizeBackups(inventory.backups);
  const normalizedSecurityGroups = normalizeSecurityGroups(inventory.securityGroups);
  const normalizedVpc = normalizeVpc(inventory.vpc);
  const vpc = normalizedVpc ? await ensureVpc(normalizedVpc, counters) : null;
  const current = await db.serviceInstance.findUnique({
    where: { id: service.id },
    select: { vpcNetworkId: true },
  });

  if (vpc && current?.vpcNetworkId !== vpc.id) {
    await db.serviceInstance.update({
      where: { id: service.id },
      data: { vpcNetworkId: vpc.id },
    });
  }

  await db.$transaction(async (tx) => {
    await tx.serviceSnapshot.deleteMany({ where: { serviceId: service.id } });
    await tx.serviceBackup.deleteMany({ where: { serviceId: service.id } });
    await tx.serviceIpAddress.deleteMany({ where: { serviceId: service.id } });
    await tx.serviceSecurityGroup.deleteMany({ where: { serviceId: service.id } });
    await tx.serviceDisk.deleteMany({ where: { serviceId: service.id } });

    if (normalizedIps.length > 0) {
      await tx.serviceIpAddress.createMany({
        data: normalizedIps.map((ip) => ({
          serviceId: service.id,
          address: ip.address,
          version: ip.version,
          type: ip.type,
          bandwidthMbps: ip.bandwidthMbps,
          providerIpId: ip.providerIpId,
          isPrimary: ip.isPrimary,
          status: ip.status,
        })),
      });
      counters.syncedIps += normalizedIps.length;
    }

    const diskIdMap = new Map<string, string>();

    for (const disk of normalizedDisks) {
      const createdDisk = await tx.serviceDisk.create({
        data: {
          serviceId: service.id,
          name: disk.name,
          type: disk.type,
          sizeGb: disk.sizeGb,
          mountPoint: disk.mountPoint,
          providerDiskId: disk.providerDiskId,
          isSystem: disk.isSystem,
          status: disk.status,
        },
      });

      counters.syncedDisks += 1;

      if (disk.providerDiskId) {
        diskIdMap.set(disk.providerDiskId, createdDisk.id);
      }

      diskIdMap.set(disk.name, createdDisk.id);
    }

    for (const snapshot of normalizedSnapshots) {
      await tx.serviceSnapshot.create({
        data: {
          serviceId: service.id,
          sourceDiskId: snapshot.sourceDiskKey
            ? (diskIdMap.get(snapshot.sourceDiskKey) ?? null)
            : null,
          name: snapshot.name,
          providerSnapshotId: snapshot.providerSnapshotId,
          sizeGb: snapshot.sizeGb,
          status: snapshot.status,
        },
      });
      counters.syncedSnapshots += 1;
    }

    for (const backup of normalizedBackups) {
      await tx.serviceBackup.create({
        data: {
          serviceId: service.id,
          name: backup.name,
          providerBackupId: backup.providerBackupId,
          sizeGb: backup.sizeGb,
          status: backup.status,
          expiresAt: backup.expiresAt,
        },
      });
      counters.syncedBackups += 1;
    }

    for (const group of normalizedSecurityGroups) {
      const createdGroup = await tx.serviceSecurityGroup.create({
        data: {
          serviceId: service.id,
          vpcNetworkId: vpc?.id,
          name: group.name,
          providerSecurityGroupId: group.providerSecurityGroupId,
          status: group.status,
        },
      });
      counters.syncedSecurityGroups += 1;

      if (group.rules.length > 0) {
        await tx.serviceSecurityGroupRule.createMany({
          data: group.rules.map((rule) => ({
            securityGroupId: createdGroup.id,
            direction: rule.direction,
            protocol: rule.protocol,
            portRange: rule.portRange,
            sourceCidr: rule.sourceCidr,
            action: rule.action,
            description: rule.description,
          })),
        });
        counters.syncedSecurityRules += group.rules.length;
      }
    }
  });
}

async function findExistingService(remoteId: string, serviceNo: string) {
  return (
    (await db.serviceInstance.findFirst({
      where: {
        providerType: ProviderType.MOFANG_CLOUD,
        providerResourceId: remoteId,
      },
    })) ??
    (await db.serviceInstance.findUnique({
      where: {
        serviceNo,
      },
    }))
  );
}

async function upsertService(input: {
  remoteId: string;
  remote: Record<string, unknown>;
  customer: Customer;
  product: Product;
  existing?: ServiceInstance | null;
}) {
  const desiredServiceNo = getRemoteServiceNo(input.remote, input.remoteId);
  const serviceNo =
    input.existing?.serviceNo ??
    (await buildUniqueServiceNo(desiredServiceNo));
  const status = mapProviderStatus(
    getRemoteStatus(input.remote) ?? "ACTIVE",
    (input.existing?.status ?? "PROVISIONING") as ServiceStatus,
  );
  const billingCycle = normalizeBillingCycle(
    pickString(input.remote, ["billing_cycle", "cycle", "duration", "pay_type"]),
    input.existing?.billingCycle ?? BillingCycle.MONTHLY,
  );
  const primaryIp = getRemotePrimaryIp(input.remote);
  const disks = normalizeDisks(
    toRecordArray(
      input.remote.disks ?? input.remote.disk_list ?? input.remote.disk ?? [],
    ),
    input.remote,
  );
  const payload = {
    customerId: input.customer.id,
    productId: input.product.id,
    name: getRemoteName(input.remote, input.remoteId),
    hostname:
      pickString(input.remote, ["hostname", "host_name", "instance_name"]) ??
      input.existing?.hostname ??
      undefined,
    providerType: ProviderType.MOFANG_CLOUD,
    providerResourceId: input.remoteId,
    region: getRemoteRegion(input.remote),
    billingCycle,
    status,
    ipAddress: primaryIp,
    cpuCores: getRemoteCpu(input.remote),
    memoryGb: getRemoteMemory(input.remote),
    storageGb: getRemoteStorage(input.remote, disks),
    monthlyCost: normalizeMoneyToCents(
      pickNumber(input.remote, [
        "renew_amount",
        "renew_price",
        "price",
        "amount",
        "monthly_price",
      ]),
    ),
    nextDueDate: getRemoteNextDueDate(input.remote),
    activatedAt:
      getRemoteActivatedAt(input.remote) ??
      (status === "ACTIVE" ? input.existing?.activatedAt ?? new Date() : undefined),
    expiresAt: getRemoteExpiresAt(input.remote),
    configSnapshot: JSON.stringify(
      {
        source: "mofang-pull-sync",
        pulledAt: new Date().toISOString(),
        remoteId: input.remoteId,
        raw: input.remote,
      },
      null,
      2,
    ),
    lastSyncAt: new Date(),
  };

  if (input.existing) {
    return {
      operation: "updated" as const,
      service: await db.serviceInstance.update({
        where: { id: input.existing.id },
        data: payload,
      }),
    };
  }

  return {
    operation: "created" as const,
    service: await db.serviceInstance.create({
      data: {
        serviceNo,
        ...payload,
      },
    }),
  };
}

export async function pullMofangCloudToLocal(
  options: PullSyncOptions = {},
): Promise<MofangPullSyncResult> {
  const provider = getCloudProvider() as MofangCloudProvider;
  const counters: SyncCounters = {
    processedCount: 0,
    createdServices: 0,
    updatedServices: 0,
    failedServices: 0,
    createdCustomers: 0,
    createdProducts: 0,
    createdPortalUsers: 0,
    updatedPortalUsers: 0,
    queuedPortalNotifications: 0,
    createdImportOrders: 0,
    updatedImportOrders: 0,
    createdVpcs: 0,
    syncedIps: 0,
    syncedDisks: 0,
    syncedSnapshots: 0,
    syncedBackups: 0,
    syncedSecurityGroups: 0,
    syncedSecurityRules: 0,
  };
  const items: PullSyncItem[] = [];
  const includeResources = options.includeResources ?? true;
  const listResult = await provider.fetchInstanceList();
  const remoteCount = listResult.items.length;
  const selectedItems =
    options.limit && options.limit > 0
      ? listResult.items.slice(0, options.limit)
      : listResult.items;

  if (!listResult.ok) {
    await writeProviderLog({
      action: "pull_sync_list",
      ok: false,
      message: listResult.message,
      responseBody: JSON.stringify(listResult, null, 2),
    });

    return {
      summary: {
        mockMode: provider.isMockMode(),
        listOk: false,
        listMessage: listResult.message,
        remoteCount,
        ...counters,
      },
      items,
    };
  }

  await writeProviderLog({
    action: "pull_sync_list",
    ok: true,
    message: listResult.message,
    responseBody: JSON.stringify(
      {
        remoteCount,
        processedCount: selectedItems.length,
      },
      null,
      2,
    ),
  });

  for (const listItem of selectedItems) {
    const listRecord = toRecord(listItem.data);
    const remoteId = listItem.remoteId ?? (listRecord ? getRemoteId(listRecord) : undefined);

    if (!remoteId) {
      counters.failedServices += 1;
      items.push({
        remoteId: "unknown",
        operation: "failed",
        message: "remote instance id is missing",
      });
      continue;
    }

    counters.processedCount += 1;

    try {
      let inventory: MofangRemoteInventory | null = null;
      let detailRecord: Record<string, unknown> | undefined;

      if (includeResources) {
        inventory = await provider.getInstanceInventory(remoteId);
        detailRecord = toRecord(inventory.detail.data) ?? inventory.raw ?? undefined;
      } else {
        const detail = await provider.getInstanceDetailById(remoteId);
        detailRecord = toRecord(detail.data);
      }

      const remote = mergeRemoteRecord(listRecord, detailRecord);
      const serviceNo = getRemoteServiceNo(remote, remoteId);
      const customer = await ensureCustomer(remote, remoteId, counters);
      const product = await ensureProduct(remote, remoteId, counters);
      const existing = await findExistingService(remoteId, serviceNo);
      const previousStatus = existing?.status;
      const result = await upsertService({
        remoteId,
        remote,
        customer,
        product,
        existing,
      });
      let portalUser: CustomerUser | null = null;
      let importedOrderNo: string | undefined;
      let importedOrderStatus: string | undefined;

      try {
        portalUser = await ensurePortalOwner(customer, remote, counters);
      } catch (portalError) {
        await writeProviderLog({
          action: "pull_sync_portal_user",
          resourceId: remoteId,
          ok: false,
          message:
            portalError instanceof Error
              ? portalError.message
              : "failed to sync portal owner",
        });
      }

      try {
        const importedOrder = await ensureImportedOrder({
          customer,
          product,
          service: result.service,
          counters,
        });
        importedOrderNo = importedOrder.orderNo;
        importedOrderStatus = importedOrder.status;
      } catch (orderError) {
        await writeProviderLog({
          action: "pull_sync_import_order",
          resourceId: remoteId,
          ok: false,
          message:
            orderError instanceof Error
              ? orderError.message
              : "failed to sync imported order",
        });
      }

      if (includeResources && inventory) {
        await syncServiceResources(result.service, remote, inventory, counters);
      }

      if (result.operation === "created") {
        counters.createdServices += 1;
      } else {
        counters.updatedServices += 1;
      }

      if (portalUser) {
        try {
          await queuePortalSyncNotification({
            customer,
            portalUser,
            service: result.service,
            operation: result.operation,
            previousStatus,
            counters,
          });
        } catch (notificationError) {
          await writeProviderLog({
            action: "pull_sync_portal_notification",
            resourceId: remoteId,
            ok: false,
            message:
              notificationError instanceof Error
                ? notificationError.message
                : "failed to queue portal sync notification",
          });
        }
      }

      items.push({
        remoteId,
        serviceId: result.service.id,
        serviceNo: result.service.serviceNo,
        orderNo: importedOrderNo,
        orderStatus: importedOrderStatus,
        customerId: customer.id,
        customerName: customer.name,
        customerEmail: customer.email,
        portalEmail: portalUser?.email,
        operation: result.operation,
        message:
          result.operation === "created"
            ? "service created from remote instance"
            : "service updated from remote instance",
      });

      await writeProviderLog({
        action: "pull_sync_service",
        resourceId: remoteId,
        ok: true,
        message: `${result.operation} service ${result.service.serviceNo}`,
        responseBody: JSON.stringify(
          {
            serviceId: result.service.id,
            serviceNo: result.service.serviceNo,
            orderNo: importedOrderNo,
            orderStatus: importedOrderStatus,
            customerId: customer.id,
            customerEmail: customer.email,
            portalEmail: portalUser?.email,
            includeResources,
          },
          null,
          2,
        ),
      });
    } catch (error) {
      counters.failedServices += 1;
      const message = error instanceof Error ? error.message : "unknown sync error";

      items.push({
        remoteId,
        operation: "failed",
        message,
      });

      await writeProviderLog({
        action: "pull_sync_service",
        resourceId: remoteId,
        ok: false,
        message,
      });
    }
  }

  return {
    summary: {
      mockMode: provider.isMockMode(),
      listOk: true,
      listMessage: listResult.message,
      remoteCount,
      ...counters,
    },
    items,
  };
}
