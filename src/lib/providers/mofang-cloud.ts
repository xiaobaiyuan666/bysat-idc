import { addCycle } from "@/lib/format";
import {
  type CloudProvider,
  type ProviderAction,
  type ProviderActionPayload,
  type ProviderResponse,
  type ProviderServicePayload,
} from "@/lib/providers/types";

type ActionSpec = {
  method: "GET" | "POST" | "PUT" | "DELETE";
  path: string;
  status: string;
  message: string;
};

type RawRequestResult = {
  ok: boolean;
  status: number;
  text: string;
  parsed?: unknown;
  init?: RequestInit;
};

export type MofangRemoteRecord = Record<string, unknown>;

export type MofangRemoteInventory = {
  detail: ProviderResponse;
  raw?: MofangRemoteRecord;
  ips: MofangRemoteRecord[];
  disks: MofangRemoteRecord[];
  snapshots: MofangRemoteRecord[];
  backups: MofangRemoteRecord[];
  securityGroups: MofangRemoteRecord[];
  vpc: MofangRemoteRecord | null;
};

export type MofangInstanceListResult = {
  ok: boolean;
  message: string;
  items: ProviderResponse[];
  raw?: ProviderResponse;
};

const actionSpecMap: Record<
  Exclude<ProviderAction, "provision" | "renew" | "sync">,
  ActionSpec
> = {
  activate: {
    method: "POST",
    path: "/v1/clouds/:id/on",
    status: "ACTIVE",
    message: "instance activated",
  },
  suspend: {
    method: "POST",
    path: "/v1/clouds/:id/suspend",
    status: "SUSPENDED",
    message: "instance suspended",
  },
  unsuspend: {
    method: "POST",
    path: "/v1/clouds/:id/unsuspend",
    status: "ACTIVE",
    message: "instance unsuspended",
  },
  terminate: {
    method: "DELETE",
    path: "/v1/clouds/:id",
    status: "TERMINATED",
    message: "instance terminated",
  },
  powerOn: {
    method: "POST",
    path: "/v1/clouds/:id/on",
    status: "ACTIVE",
    message: "instance powered on",
  },
  powerOff: {
    method: "POST",
    path: "/v1/clouds/:id/off",
    status: "SUSPENDED",
    message: "instance powered off",
  },
  reboot: {
    method: "POST",
    path: "/v1/clouds/:id/reboot",
    status: "ACTIVE",
    message: "instance rebooted",
  },
  hardReboot: {
    method: "POST",
    path: "/v1/clouds/:id/hard_reboot",
    status: "ACTIVE",
    message: "instance hard rebooted",
  },
  hardPowerOff: {
    method: "POST",
    path: "/v1/clouds/:id/hardoff",
    status: "SUSPENDED",
    message: "instance hard powered off",
  },
  reinstall: {
    method: "PUT",
    path: "/v1/clouds/:id/reinstall",
    status: "PROVISIONING",
    message: "instance reinstall submitted",
  },
  resetPassword: {
    method: "PUT",
    path: "/v1/clouds/:id/password",
    status: "ACTIVE",
    message: "instance password reset submitted",
  },
  getVnc: {
    method: "POST",
    path: "/v1/clouds/:id/vnc",
    status: "ACTIVE",
    message: "instance remote console fetched",
  },
  rescueStart: {
    method: "POST",
    path: "/v1/clouds/:id/rescue",
    status: "PROVISIONING",
    message: "instance rescue started",
  },
  rescueStop: {
    method: "DELETE",
    path: "/v1/clouds/:id/rescue",
    status: "ACTIVE",
    message: "instance rescue stopped",
  },
  lock: {
    method: "POST",
    path: "/admin/v1/mf_cloud/:id/lock",
    status: "ACTIVE",
    message: "instance locked",
  },
  unlock: {
    method: "POST",
    path: "/admin/v1/mf_cloud/:id/unlock",
    status: "ACTIVE",
    message: "instance unlocked",
  },
};

function isRecord(value: unknown): value is MofangRemoteRecord {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value);
}

function safeJsonParse(input: string) {
  try {
    return JSON.parse(input) as unknown;
  } catch {
    return undefined;
  }
}

function pickString(record: MofangRemoteRecord, keys: string[]) {
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

function pickNumber(record: MofangRemoteRecord, keys: string[]) {
  for (const key of keys) {
    const value = record[key];

    if (typeof value === "number" && Number.isFinite(value)) {
      return value;
    }

    if (typeof value === "string" && value.trim()) {
      const parsed = Number(value);

      if (Number.isFinite(parsed)) {
        return parsed;
      }
    }
  }

  return undefined;
}

function pickDate(record: MofangRemoteRecord, keys: string[]) {
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

function unwrapPayload(payload: unknown): unknown {
  if (!isRecord(payload)) {
    return payload;
  }

  for (const key of ["data", "result", "payload"]) {
    if (key in payload) {
      return unwrapPayload(payload[key]);
    }
  }

  return payload;
}

function extractRecords(payload: unknown): MofangRemoteRecord[] {
  const unwrapped = unwrapPayload(payload);

  if (Array.isArray(unwrapped)) {
    return unwrapped.filter(isRecord);
  }

  if (!isRecord(unwrapped)) {
    return [];
  }

  for (const key of [
    "list",
    "items",
    "rows",
    "clouds",
    "servers",
    "instances",
    "records",
  ]) {
    const value = unwrapped[key];

    if (Array.isArray(value)) {
      return value.filter(isRecord);
    }
  }

  return [unwrapped];
}

function extractMeta(payload: unknown) {
  if (!isRecord(payload) || !isRecord(payload.meta)) {
    return undefined;
  }

  return payload.meta;
}

function normalizeProviderStatus(rawStatus?: string) {
  const normalized = rawStatus
    ?.trim()
    .replace(/([a-z])([A-Z])/g, "$1_$2")
    .replace(/[\s-]+/g, "_")
    .toUpperCase();

  if (!normalized) {
    return "ACTIVE";
  }

  if (["ACTIVE", "RUNNING", "ON", "STARTED", "NORMAL"].includes(normalized)) {
    return "ACTIVE";
  }

  if (
    ["SUSPENDED", "SUSPEND", "STOPPED", "OFF", "SHUTOFF", "PAUSED", "DISABLED"].includes(
      normalized,
    )
  ) {
    return "SUSPENDED";
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

  if (["TERMINATED", "DELETED", "DESTROYED", "CANCELLED", "CANCELED"].includes(normalized)) {
    return "TERMINATED";
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

  return normalized;
}

function buildMockIp(seed: string) {
  const total = seed.split("").reduce((sum, char) => sum + char.charCodeAt(0), 0);
  return `10.${(total % 20) + 10}.${(total % 100) + 1}.${(total % 200) + 20}`;
}

function buildMockConsoleUrl(service: ProviderServicePayload) {
  return `https://console.mofang.local/vnc/${service.serviceNo.toLowerCase()}`;
}

function hasProviderStatusField(data?: Record<string, unknown>) {
  if (!data) {
    return false;
  }

  return ["status", "instance_status", "cloud_status", "power_status"].some(
    (key) => typeof data[key] === "string" && String(data[key]).trim().length > 0,
  );
}

function buildMockResponse(
  service: ProviderServicePayload,
  status: string,
  message: string,
  action?: ProviderAction,
  payload?: ProviderActionPayload,
): ProviderResponse {
  return {
    ok: true,
    status,
    remoteId: service.providerResourceId ?? `mf-${service.serviceNo.toLowerCase()}`,
    ipAddress: buildMockIp(service.serviceNo),
    region: service.region ?? "cn-hk-1",
    nextDueDate: addCycle(new Date(), service.billingCycle),
    message,
    consoleUrl: action === "getVnc" ? buildMockConsoleUrl(service) : undefined,
    taskId:
      action && ["reinstall", "reboot", "hardReboot", "rescueStart"].includes(action)
        ? `task-${service.serviceNo.toLowerCase()}`
        : undefined,
    data: payload ? { payload } : undefined,
    requestBody: JSON.stringify({ action, service, payload }, null, 2),
    responseBody: JSON.stringify(
      {
        provider: "MOFANG_CLOUD",
        serviceNo: service.serviceNo,
        status,
        message,
      },
      null,
      2,
    ),
  };
}

export class MofangCloudProvider implements CloudProvider {
  name = "MOFANG_CLOUD";
  private sessionToken?: string;
  private activeBaseUrl?: string;

  private get baseUrl() {
    return process.env.MOFANG_CLOUD_BASE_URL;
  }

  private get useMock() {
    return process.env.MOFANG_CLOUD_USE_MOCK !== "false" || !this.baseUrl;
  }

  isMockMode() {
    return this.useMock;
  }

  private getBaseUrlCandidates() {
    if (!this.baseUrl) {
      return [] as string[];
    }

    const candidates = [this.baseUrl];

    if (this.baseUrl.startsWith("https://")) {
      candidates.push(this.baseUrl.replace(/^https:\/\//i, "http://"));
    }

    return Array.from(new Set(candidates));
  }

  private hasSessionCredentials() {
    return Boolean(
      process.env.MOFANG_CLOUD_USERNAME && process.env.MOFANG_CLOUD_PASSWORD,
    );
  }

  private async getAccessToken(forceRefresh = false) {
    if (this.useMock) {
      return undefined;
    }

    if (!forceRefresh && this.sessionToken) {
      return this.sessionToken;
    }

    if (!this.hasSessionCredentials() || !this.baseUrl) {
      return process.env.MOFANG_CLOUD_API_KEY ?? undefined;
    }

    const body = new URLSearchParams();
    body.set("username", process.env.MOFANG_CLOUD_USERNAME ?? "");
    body.set("password", process.env.MOFANG_CLOUD_PASSWORD ?? "");
    body.set(
      "customfield[google_code]",
      process.env.MOFANG_CLOUD_GOOGLE_CODE ?? "",
    );
    let lastStatus = 0;

    for (const candidate of this.getBaseUrlCandidates()) {
      try {
        const response = await fetch(
          `${candidate}${this.buildAbsolutePath("/v1/login")}`,
          {
            method: "POST",
            headers: {
              "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
              "think-lang": process.env.MOFANG_CLOUD_LANG ?? "zh-cn",
            },
            body: body.toString(),
            cache: "no-store",
          },
        );
        const text = await response.text();
        const parsed = safeJsonParse(text);
        const token =
          (typeof parsed === "string" && parsed.trim()) ||
          (isRecord(parsed)
            ? pickString(parsed, ["token", "access_token"]) ??
              (isRecord(parsed.data)
                ? pickString(parsed.data, ["token", "access_token"])
                : undefined)
            : undefined);

        lastStatus = response.status;
        if (!response.ok || !token) {
          continue;
        }

        this.activeBaseUrl = candidate;
        this.sessionToken = token;
        return token;
      } catch {
        continue;
      }
    }

    throw new Error(
      `failed to authenticate against Mofang Cloud: ${lastStatus || "network"}`,
    );
  }

  private async getHeaders(forceRefresh = false) {
    const sessionToken = await this.getAccessToken(forceRefresh).catch(() => undefined);

    return {
      "Content-Type": "application/json",
      "think-lang": process.env.MOFANG_CLOUD_LANG ?? "zh-cn",
      ...(sessionToken
        ? {
            Authorization: `JWT ${sessionToken}`,
            "access-token": sessionToken,
          }
        : {
            Authorization: `Bearer ${process.env.MOFANG_CLOUD_API_KEY ?? ""}`,
            "x-api-secret": process.env.MOFANG_CLOUD_API_SECRET ?? "",
          }),
    };
  }

  private replacePathParams(path: string, resourceId?: string | null) {
    return path.replace(":id", resourceId ?? "").replace("{id}", resourceId ?? "");
  }

  private buildAbsolutePath(path: string) {
    return path.startsWith("/") ? path : `/${path}`;
  }

  private async buildRequestInit(init?: RequestInit, forceRefresh = false) {
    return {
      ...init,
      headers: {
        ...(await this.getHeaders(forceRefresh)),
        ...(init?.headers ?? {}),
      },
      cache: "no-store" as const,
    };
  }

  private appendFormValue(
    form: URLSearchParams,
    key: string,
    value: unknown,
  ) {
    if (value === undefined || value === null) {
      return;
    }

    if (Array.isArray(value)) {
      value.forEach((item, index) => {
        this.appendFormValue(form, `${key}[${index}]`, item);
      });
      return;
    }

    if (typeof value === "object") {
      Object.entries(value as Record<string, unknown>).forEach(([childKey, childValue]) => {
        this.appendFormValue(form, `${key}[${childKey}]`, childValue);
      });
      return;
    }

    form.append(key, String(value));
  }

  private buildRequestBody(payload?: Record<string, unknown>) {
    if (!payload) {
      return {
        body: undefined,
        headers: undefined,
      };
    }

    if (!this.hasSessionCredentials()) {
      return {
        body: JSON.stringify(payload),
        headers: {
          "Content-Type": "application/json",
        },
      };
    }

    const form = new URLSearchParams();

    Object.entries(payload).forEach(([key, value]) => {
      this.appendFormValue(form, key, value);
    });

    return {
      body: form.toString(),
      headers: {
        "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8",
      },
    };
  }

  private buildResponseFromRecord(
    record: MofangRemoteRecord,
    requestBody?: string,
    responseBody?: string,
  ): ProviderResponse {
    const remoteId = pickString(record, [
      "id",
      "cloudid",
      "cloud_id",
      "instance_id",
      "hostid",
      "server_id",
      "uuid",
    ]);
    const status = normalizeProviderStatus(
      pickString(record, [
        "status",
        "instance_status",
        "cloud_status",
        "power_status",
      ]),
    );
    const ipAddress = pickString(record, [
      "ip",
      "ip_address",
      "main_ip",
      "mainip",
      "public_ip",
      "primary_ip",
    ]);
    const region =
      (isRecord(record.area)
        ? pickString(record.area, ["name", "short_name"])
        : undefined) ??
      pickString(record, [
        "region_name",
        "region",
        "area",
        "zone_name",
        "node_name",
      ]);
    const nextDueDate = pickDate(record, [
      "next_due_date",
      "due_date",
      "renew_time",
      "expire_time",
      "expired_at",
      "end_time",
    ]);
    const consoleUrl = pickString(record, ["url", "vnc_url", "ws_url", "remote_url"]);
    const taskId = pickString(record, ["taskid", "task_id", "job_id"]);
    const message = pickString(record, [
      "msg",
      "message",
      "status_desc",
      "description",
    ]);

    return {
      ok: true,
      status,
      remoteId,
      ipAddress,
      region,
      nextDueDate,
      consoleUrl,
      taskId,
      message,
      requestBody,
      responseBody,
      data: record,
    };
  }

  private normalizeResponse(
    raw: RawRequestResult,
    fallbackMessage?: string,
  ): ProviderResponse {
    const payload = unwrapPayload(raw.parsed);

    if (raw.ok && isRecord(payload)) {
      return {
        ...this.buildResponseFromRecord(
          payload,
          typeof raw.init?.body === "string" ? raw.init.body : undefined,
          raw.text,
        ),
        ok: true,
        message:
          pickString(payload, ["msg", "message", "status_desc"]) ??
          fallbackMessage ??
          "request succeeded",
      };
    }

    const parsedRoot = isRecord(raw.parsed) ? raw.parsed : undefined;

    return {
      ok: raw.ok,
      status: raw.ok ? "ACTIVE" : "FAILED",
      message:
        pickString(parsedRoot ?? {}, ["msg", "message"]) ??
        fallbackMessage ??
        (raw.ok ? "request succeeded" : `request failed: ${raw.status}`),
      requestBody: typeof raw.init?.body === "string" ? raw.init.body : undefined,
      responseBody: raw.text,
    };
  }

  private async requestRaw(
    path: string,
    init?: RequestInit,
    forceRefresh = false,
  ): Promise<RawRequestResult> {
    if (!this.baseUrl) {
      return {
        ok: false,
        status: 0,
        text: "",
        parsed: undefined,
        init,
      };
    }

    try {
      const requestInit = await this.buildRequestInit(init, forceRefresh);
      const candidates =
        this.activeBaseUrl && this.activeBaseUrl.trim()
          ? [this.activeBaseUrl]
          : this.getBaseUrlCandidates();

      for (const candidate of candidates) {
        try {
          const response = await fetch(
            `${candidate}${this.buildAbsolutePath(path)}`,
            requestInit,
          );
          const text = await response.text();
          const parsed = safeJsonParse(text);

          this.activeBaseUrl = candidate;

          if (
            response.status === 401 &&
            this.hasSessionCredentials() &&
            !forceRefresh
          ) {
            this.sessionToken = undefined;
            return this.requestRaw(path, init, true);
          }

          return {
            ok: response.ok,
            status: response.status,
            text,
            parsed,
            init: requestInit,
          };
        } catch {
          continue;
        }
      }
    } catch (error) {
      return {
        ok: false,
        status: 0,
        text: error instanceof Error ? error.message : "network error",
        parsed: undefined,
        init,
      };
    }

    return {
      ok: false,
      status: 0,
      text: "network error",
      parsed: undefined,
      init,
    };
  }

  private async request(path: string, init?: RequestInit): Promise<ProviderResponse> {
    const raw = await this.requestRaw(path, init);

    if (!this.baseUrl) {
      return {
        ok: false,
        status: "FAILED",
        message: "MOFANG_CLOUD_BASE_URL is not configured",
      };
    }

    return this.normalizeResponse(raw);
  }

  private async requestCollection(
    resourceId: string,
    envKey: string,
    fallbackPaths: string[],
  ) {
    if (this.useMock) {
      return [] as MofangRemoteRecord[];
    }

    const candidates = [
      process.env[envKey],
      ...fallbackPaths,
    ].filter((item): item is string => Boolean(item));

    for (const candidate of candidates) {
      const raw = await this.requestRaw(this.replacePathParams(candidate, resourceId), {
        method: "GET",
      });

      if (!raw.ok) {
        continue;
      }

      const records = extractRecords(raw.parsed);

      if (records.length > 0) {
        return records;
      }

      const payload = unwrapPayload(raw.parsed);

      if (isRecord(payload)) {
        return [payload];
      }
    }

    return [] as MofangRemoteRecord[];
  }

  private async requestSingleRecord(
    resourceId: string,
    envKey: string,
    fallbackPaths: string[],
  ) {
    const records = await this.requestCollection(resourceId, envKey, fallbackPaths);
    return records[0] ?? null;
  }

  async fetchInstanceList(): Promise<MofangInstanceListResult> {
    if (this.useMock) {
      return {
        ok: true,
        message: "mock mode is enabled, remote pull sync is skipped",
        items: [],
      };
    }

    const candidates = [
      process.env.MOFANG_CLOUD_LIST_PATH,
      "/v1/clouds?page=1&per_page=100&sort=desc&orderby=id",
    ].filter((item): item is string => Boolean(item));

    for (const candidate of candidates) {
      const items: ProviderResponse[] = [];
      let page = 1;
      let totalPage = 1;
      let lastRaw: RawRequestResult | undefined;

      while (page <= totalPage) {
        const pagePath =
          candidate.includes("page=")
            ? candidate.replace(/page=\d+/i, `page=${page}`)
            : `${candidate}${candidate.includes("?") ? "&" : "?"}page=${page}&per_page=100&sort=desc&orderby=id`;
        const raw = await this.requestRaw(pagePath, { method: "GET" });

        lastRaw = raw;
        if (!raw.ok) {
          items.length = 0;
          break;
        }

        const parsedRoot = isRecord(raw.parsed) ? raw.parsed : undefined;
        const records = extractRecords(raw.parsed);

        items.push(
          ...records
            .map((record) =>
              this.buildResponseFromRecord(
                record,
                JSON.stringify({ path: pagePath }, null, 2),
                raw.text,
              ),
            )
            .filter((item) => Boolean(item.remoteId)),
        );

        const meta = extractMeta(parsedRoot);
        totalPage = Math.max(
          1,
          pickNumber(meta ?? {}, ["total_page"]) ??
            (records.length >= 100 ? page + 1 : page),
        );
        page += 1;
      }

      if (!lastRaw?.ok) {
        continue;
      }

      return {
        ok: true,
        message:
          items.length > 0
            ? `loaded ${items.length} remote instances`
            : "instance list request succeeded but returned no remote instances",
        items,
        raw: this.normalizeResponse(
          lastRaw,
          `instance list request failed: ${candidate}`,
        ),
      };
    }

    return {
      ok: false,
      message: "failed to load instance list from configured MOFANG_CLOUD endpoints",
      items: [],
    };
  }

  async listInstances() {
    const result = await this.fetchInstanceList();
    return result.items;
  }

  async getInstanceDetailById(resourceId: string) {
    if (this.useMock) {
      return buildMockResponse(
        {
          localId: resourceId,
          serviceNo: resourceId,
          name: resourceId,
          providerResourceId: resourceId,
          billingCycle: "MONTHLY",
          customerName: "Mock Customer",
          productName: "Mock Product",
        },
        "ACTIVE",
        "mock instance detail",
        "sync",
      );
    }

    const candidates = [
      process.env.MOFANG_CLOUD_INSTANCE_DETAIL_PATH,
      "/v1/clouds/:id",
    ].filter((item): item is string => Boolean(item));

    for (const candidate of candidates) {
      const raw = await this.requestRaw(this.replacePathParams(candidate, resourceId), {
        method: "GET",
      });

      if (raw.ok) {
        return this.normalizeResponse(raw);
      }
    }

    return {
      ok: false,
      status: "FAILED",
      message: `failed to load remote instance detail: ${resourceId}`,
    };
  }

  async getInstanceInventory(resourceId: string): Promise<MofangRemoteInventory> {
    const detail = await this.getInstanceDetailById(resourceId);
    const raw = isRecord(detail.data) ? detail.data : undefined;

    const [ips, disks, snapshots, backups, securityGroups, vpc] = await Promise.all([
      this.requestCollection(resourceId, "MOFANG_CLOUD_IP_PATH", [
        "/v1/clouds/:id/ip",
      ]),
      this.requestCollection(resourceId, "MOFANG_CLOUD_DISK_PATH", [
        "/v1/clouds/:id/disks",
      ]),
      this.requestCollection(resourceId, "MOFANG_CLOUD_SNAPSHOT_PATH", [
        "/v1/clouds/:id/snapshots",
      ]),
      this.requestCollection(resourceId, "MOFANG_CLOUD_BACKUP_PATH", [
        "/v1/clouds/:id/backups",
      ]),
      this.requestCollection(resourceId, "MOFANG_CLOUD_SECURITY_GROUP_PATH", [
        "/v1/clouds/:id/security_groups",
      ]),
      this.requestSingleRecord(resourceId, "MOFANG_CLOUD_VPC_PATH", [
        "/v1/clouds/:id/vpc_network",
      ]),
    ]);

    return {
      detail,
      raw,
      ips:
        ips.length > 0
          ? ips
          : raw
            ? extractRecords(raw.ips ?? raw.ip_list ?? raw.ips_list ?? raw.ip ?? [])
            : [],
      disks:
        disks.length > 0
          ? disks
          : raw
            ? extractRecords(raw.disks ?? raw.disk_list ?? raw.disk ?? [])
            : [],
      snapshots:
        snapshots.length > 0
          ? snapshots
          : raw
            ? extractRecords(raw.snapshots ?? raw.snapshot_list ?? raw.snapshot ?? [])
            : [],
      backups:
        backups.length > 0
          ? backups
          : raw
            ? extractRecords(raw.backups ?? raw.backup_list ?? raw.backup ?? [])
            : [],
      securityGroups:
        securityGroups.length > 0
          ? securityGroups
          : raw
            ? extractRecords(
                raw.security_groups ??
                  raw.security_group ??
                  raw.firewall_groups ??
                  raw.firewall_group ??
                  [],
              )
            : [],
      vpc:
        vpc ??
        (raw
          ? (extractRecords(raw.vpc_network ?? raw.vpc ?? raw.vpc_info ?? [])[0] ?? null)
          : null),
    };
  }

  async provisionService(service: ProviderServicePayload) {
    if (this.useMock) {
      return buildMockResponse(service, "ACTIVE", "mock provision completed", "provision");
    }

    const path = process.env.MOFANG_CLOUD_INSTANCE_PATH ?? "/v1/clouds";
    const request = this.buildRequestBody({
      remarks: service.name,
      hostname: service.hostname,
      product_id: service.providerProductId,
      ...service.configOptions,
    });

    return this.request(path, {
      method: "POST",
      headers: request.headers,
      body: request.body,
    });
  }

  async activateService(service: ProviderServicePayload) {
    return this.manageServiceAction(service, "activate");
  }

  async suspendService(service: ProviderServicePayload) {
    return this.manageServiceAction(service, "suspend");
  }

  async renewService(service: ProviderServicePayload) {
    if (this.useMock) {
      return buildMockResponse(service, "ACTIVE", "mock renew completed", "renew");
    }

    const path = this.replacePathParams(
      process.env.MOFANG_CLOUD_RENEW_PATH ?? "/v1/clouds/:id/renew",
      service.providerResourceId,
    );
    const request = this.buildRequestBody({
      cycle: service.billingCycle,
    });

    return this.request(path, {
      method: "POST",
      headers: request.headers,
      body: request.body,
    });
  }

  async terminateService(service: ProviderServicePayload) {
    return this.manageServiceAction(service, "terminate");
  }

  async syncService(service: ProviderServicePayload) {
    return this.getInstanceDetailById(service.providerResourceId ?? service.serviceNo);
  }

  async manageServiceAction(
    service: ProviderServicePayload,
    action: ProviderAction,
    payload?: ProviderActionPayload,
  ) {
    if (action === "provision") {
      return this.provisionService(service);
    }

    if (action === "renew") {
      return this.renewService(service);
    }

    if (action === "sync") {
      return this.syncService(service);
    }

    const spec = actionSpecMap[action];

    if (this.useMock) {
      return buildMockResponse(service, spec.status, `mock ${spec.message}`, action, payload);
    }

    const envKey = `MOFANG_CLOUD_${action.toUpperCase()}_PATH`;
    const customPath = process.env[envKey];
    const path = this.replacePathParams(customPath ?? spec.path, service.providerResourceId);
    const request = this.buildRequestBody(payload ? { ...payload } : undefined);

    const response = await this.request(path, {
      method: spec.method,
      headers:
        spec.method === "GET" || spec.method === "DELETE"
          ? undefined
          : request.headers,
      body:
        spec.method === "GET" || spec.method === "DELETE"
          ? undefined
          : request.body,
    });

    if (!response.ok) {
      return response;
    }

    if (!hasProviderStatusField(response.data)) {
      return {
        ...response,
        status: spec.status,
        message:
          response.message && response.message !== "request succeeded"
            ? response.message
            : spec.message,
      };
    }

    return response;
  }
}
