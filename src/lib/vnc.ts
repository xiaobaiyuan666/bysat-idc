import { buildBrowserVncSocketUrl, ensureVncProxyServer } from "@/lib/vnc-proxy-server";

type VncPayload = Record<string, unknown>;

function isRecord(value: unknown): value is VncPayload {
  return Boolean(value) && typeof value === "object" && !Array.isArray(value);
}

function readString(payload: VncPayload, keys: string[]) {
  for (const key of keys) {
    const value = payload[key];

    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }
  }

  return undefined;
}

function readText(payload: VncPayload, keys: string[]) {
  for (const key of keys) {
    const value = payload[key];

    if (typeof value === "string" && value.trim()) {
      return value.trim();
    }

    if (typeof value === "number" && Number.isFinite(value)) {
      return String(value);
    }
  }

  return undefined;
}

function buildProviderUrl(baseUrl: string, relativePath: string) {
  const normalizedBase = `${baseUrl.replace(/\/+$/, "")}/`;
  return new URL(relativePath.replace(/^\/+/, ""), normalizedBase);
}

function buildMofangNoVncUrl(data: VncPayload) {
  const baseUrl = process.env.MOFANG_CLOUD_BASE_URL?.trim();

  if (!baseUrl) {
    return undefined;
  }

  const token = readString(data, ["token"]);
  const vncPassword = readString(data, ["vnc_token", "vnc_pass", "vncpass"]);
  const pathValue = readText(data, ["path"]);

  if (!token || !vncPassword || !pathValue) {
    return undefined;
  }

  const url = buildProviderUrl(baseUrl, "v1/novnc");
  const ip = readString(data, ["ip", "mainip", "ip_address"]);
  const title = readString(data, ["hostname", "name"]) ?? ip;
  const hostToken = readString(data, ["host_token", "password", "rootpassword", "panel_pass"]);

  url.searchParams.set("token", token);
  url.searchParams.set("vnc_token", vncPassword);
  url.searchParams.set("password", vncPassword);
  url.searchParams.set("path", pathValue);
  url.searchParams.set("lang", process.env.MOFANG_CLOUD_LANG?.trim() || "zh-cn");

  if (hostToken) {
    url.searchParams.set("host_token", hostToken);
  }

  if (ip) {
    url.searchParams.set("ip", ip);
  }

  if (title) {
    url.searchParams.set("title", title);
  }

  return url.toString();
}

export async function buildVncConsoleUrl(
  baseOrigin: string,
  data?: Record<string, unknown>,
) {
  if (!isRecord(data)) {
    return undefined;
  }

  const officialConsoleUrl = buildMofangNoVncUrl(data);

  if (officialConsoleUrl) {
    return officialConsoleUrl;
  }

  const rawWsUrl = readString(data, ["vnc_url_https", "vnc_url", "vnc_url_http", "ws_url"]);

  if (!rawWsUrl) {
    return undefined;
  }

  await ensureVncProxyServer();

  const params = new URLSearchParams({
    ws: buildBrowserVncSocketUrl(baseOrigin, rawWsUrl),
  });

  const password = readString(data, ["vnc_pass", "password"]);
  const ip = readString(data, ["ip"]);
  const title = readString(data, ["hostname", "name"]) ?? ip;

  if (password) {
    params.set("password", password);
  }

  if (ip) {
    params.set("ip", ip);
  }

  if (title) {
    params.set("title", title);
  }

  return `${baseOrigin.replace(/\/$/, "")}/vnc?${params.toString()}`;
}
