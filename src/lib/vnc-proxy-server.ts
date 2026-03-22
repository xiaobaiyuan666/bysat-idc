import "server-only";

import { createServer } from "node:http";
import type { AddressInfo } from "node:net";
import type { Duplex } from "node:stream";

import { WebSocket, WebSocketServer } from "ws";

type VncProxyState = {
  host: string;
  port: number;
  serverUrl: string;
};

declare global {
  var __IDC_VNC_PROXY_STATE__: Promise<VncProxyState> | undefined;
}

function getConfiguredProxyHost() {
  return process.env.VNC_PROXY_HOST?.trim() || "0.0.0.0";
}

function getConfiguredProxyPort() {
  const parsed = Number(process.env.VNC_PROXY_PORT ?? "3081");
  return Number.isFinite(parsed) && parsed > 0 ? Math.round(parsed) : 3081;
}

function getAllowedVncHosts() {
  const hosts = new Set<string>();
  const baseUrl = process.env.MOFANG_CLOUD_BASE_URL;

  if (baseUrl) {
    try {
      hosts.add(new URL(baseUrl).hostname);
    } catch {
      // Ignore invalid base URL.
    }
  }

  const extraHosts = process.env.VNC_PROXY_ALLOWED_HOSTS;

  if (extraHosts) {
    for (const host of extraHosts.split(",")) {
      const normalized = host.trim();
      if (normalized) {
        hosts.add(normalized);
      }
    }
  }

  return hosts;
}

export function normalizeVncTargetUrl(rawUrl: string) {
  const normalized = new URL(rawUrl);
  const allowedHosts = getAllowedVncHosts();

  if (!["ws:", "wss:"].includes(normalized.protocol)) {
    throw new Error("VNC_PROXY_UNSUPPORTED_PROTOCOL");
  }

  if (allowedHosts.size > 0 && !allowedHosts.has(normalized.hostname)) {
    throw new Error("VNC_PROXY_FORBIDDEN_HOST");
  }

  if (
    normalized.protocol === "ws:" &&
    process.env.MOFANG_CLOUD_BASE_URL?.startsWith("https://")
  ) {
    normalized.protocol = "wss:";
  }

  return normalized.toString();
}

export function buildBrowserVncSocketUrl(baseOrigin: string, targetUrl: string) {
  const origin = new URL(baseOrigin);
  const proxyProtocol = origin.protocol === "https:" ? "wss:" : "ws:";
  const proxyHost = process.env.VNC_PROXY_PUBLIC_HOST?.trim() || origin.hostname;
  const proxyPort = getConfiguredProxyPort();
  const url = new URL(`${proxyProtocol}//${proxyHost}:${proxyPort}/vnc-proxy`);

  url.searchParams.set("target", normalizeVncTargetUrl(targetUrl));
  return url.toString();
}

function rejectUpgrade(socket: Duplex, statusCode = 400, reason = "Bad Request") {
  socket.write(`HTTP/1.1 ${statusCode} ${reason}\r\nConnection: close\r\n\r\n`);
  socket.destroy();
}

function normalizeCloseCode(code?: number) {
  if (code === undefined) {
    return undefined;
  }

  if (code >= 1000 && code <= 1014 && code !== 1004 && code !== 1005 && code !== 1006) {
    return code;
  }

  if (code >= 3000 && code <= 4999) {
    return code;
  }

  return undefined;
}

function wireSockets(browserSocket: WebSocket, targetUrl: string) {
  let downstreamOpen = false;
  const pending: Array<{ data: Parameters<WebSocket["send"]>[0]; isBinary: boolean }> = [];
  const downstream = new WebSocket(targetUrl, {
    rejectUnauthorized: false,
    perMessageDeflate: false,
  });

  const flushPending = () => {
    while (pending.length > 0 && downstream.readyState === WebSocket.OPEN) {
      const next = pending.shift();
      if (!next) {
        continue;
      }
      downstream.send(next.data, { binary: next.isBinary });
    }
  };

  browserSocket.on("message", (data, isBinary) => {
    if (!downstreamOpen || downstream.readyState !== WebSocket.OPEN) {
      pending.push({ data, isBinary });
      return;
    }

    downstream.send(data, { binary: isBinary });
  });

  downstream.on("open", () => {
    downstreamOpen = true;
    flushPending();
  });

  downstream.on("message", (data, isBinary) => {
    if (browserSocket.readyState === WebSocket.OPEN) {
      browserSocket.send(data, { binary: isBinary });
    }
  });

  const closeBrowser = (code?: number, reason?: string) => {
    if (browserSocket.readyState === WebSocket.OPEN) {
      const normalizedCode = normalizeCloseCode(code);
      browserSocket.close(normalizedCode, reason);
    }
  };

  const closeDownstream = (code?: number, reason?: string) => {
    if (downstream.readyState === WebSocket.OPEN || downstream.readyState === WebSocket.CONNECTING) {
      const normalizedCode = normalizeCloseCode(code);
      downstream.close(normalizedCode, reason);
    }
  };

  browserSocket.on("close", (code, reason) => {
    closeDownstream(code, reason.toString());
  });

  browserSocket.on("error", () => {
    closeDownstream(1011, "Browser socket error");
  });

  downstream.on("close", (code, reason) => {
    closeBrowser(code, reason.toString() || "Downstream server closed");
  });

  downstream.on("error", (error) => {
    closeBrowser(1011, error.message || "Failed to connect to downstream server");
  });
}

export async function ensureVncProxyServer() {
  if (!globalThis.__IDC_VNC_PROXY_STATE__) {
    globalThis.__IDC_VNC_PROXY_STATE__ = new Promise<VncProxyState>((resolve, reject) => {
      const host = getConfiguredProxyHost();
      const port = getConfiguredProxyPort();
      const server = createServer((_request, response) => {
        response.statusCode = 426;
        response.end("WebSocket endpoint only");
      });
      const wss = new WebSocketServer({
        noServer: true,
        perMessageDeflate: false,
      });

      server.on("upgrade", (request, socket, head) => {
        try {
          const requestUrl = new URL(
            request.url ?? "/",
            `http://${request.headers.host ?? `127.0.0.1:${port}`}`,
          );

          if (requestUrl.pathname !== "/vnc-proxy") {
            rejectUpgrade(socket, 404, "Not Found");
            return;
          }

          const target = requestUrl.searchParams.get("target");

          if (!target) {
            rejectUpgrade(socket, 400, "Missing target");
            return;
          }

          const normalizedTarget = normalizeVncTargetUrl(target);

          wss.handleUpgrade(request, socket, head, (browserSocket) => {
            wireSockets(browserSocket, normalizedTarget);
          });
        } catch (error) {
          rejectUpgrade(
            socket,
            400,
            error instanceof Error ? error.message : "Invalid target",
          );
        }
      });

      server.once("error", (error: NodeJS.ErrnoException) => {
        if (error.code === "EADDRINUSE") {
          resolve({
            host,
            port,
            serverUrl: `ws://127.0.0.1:${port}/vnc-proxy`,
          });
          return;
        }

        globalThis.__IDC_VNC_PROXY_STATE__ = undefined;
        reject(error);
      });

      server.listen(port, host, () => {
        const address = server.address() as AddressInfo | null;
        const actualPort = address?.port ?? port;

        resolve({
          host,
          port: actualPort,
          serverUrl: `ws://127.0.0.1:${actualPort}/vnc-proxy`,
        });
      });
    });
  }

  return globalThis.__IDC_VNC_PROXY_STATE__;
}
