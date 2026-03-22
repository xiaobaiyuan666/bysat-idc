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

export function buildVncPageUrl(
  payload: Record<string, unknown> | null | undefined,
  options?: {
    ip?: string;
    title?: string;
  }
) {
  if (!isRecord(payload)) {
    return undefined;
  }

  const wsUrl = readString(payload, ["vnc_url_https", "vnc_url", "vnc_url_http", "ws_url"]);
  if (!wsUrl) {
    return undefined;
  }

  const params = new URLSearchParams({
    ws: wsUrl
  });

  const password = readString(payload, ["vnc_pass", "password"]);
  const ip = options?.ip || readString(payload, ["ip", "mainip", "ip_address"]);
  const title = options?.title || readText(payload, ["hostname", "name"]) || ip;

  if (password) {
    params.set("password", password);
  }
  if (ip) {
    params.set("ip", ip);
  }
  if (title) {
    params.set("title", title);
  }

  return `/vnc?${params.toString()}`;
}
