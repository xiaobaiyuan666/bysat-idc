"use client";

import { useEffect, useRef, useState } from "react";

import type { NoVncModule, NoVncRfb } from "@/types/novnc";

type VncConsoleProps = {
  wsUrl: string;
  password?: string;
};

export function VncConsole({ wsUrl, password }: VncConsoleProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const [status, setStatus] = useState("正在连接控制台...");
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let disposed = false;
    let rfb: NoVncRfb | undefined;

    const connect = async () => {
      if (!containerRef.current) {
        return;
      }

      try {
        const modulePath = "/vendor/novnc/rfb.bundle.mjs";
        const { default: RFB } = (await import(
          /* webpackIgnore: true */ modulePath
        )) as NoVncModule;

        if (disposed || !containerRef.current) {
          return;
        }

        rfb = new RFB(
          containerRef.current,
          wsUrl,
          password ? { credentials: { password } } : undefined,
        );
        rfb.scaleViewport = true;
        rfb.resizeSession = true;
        rfb.clipViewport = false;
        rfb.background = "#0f172a";

        rfb.addEventListener("connect", () => {
          if (!disposed) {
            setStatus("控制台已连接");
            setError(null);
          }
        });

        rfb.addEventListener("disconnect", (event: { detail?: { clean?: boolean } }) => {
          if (!disposed) {
            setStatus(event?.detail?.clean ? "控制台已断开" : "控制台连接中断");
          }
        });

        rfb.addEventListener("credentialsrequired", () => {
          if (password && rfb) {
            rfb.sendCredentials({ password });
            return;
          }

          if (!disposed) {
            setStatus("控制台需要 VNC 密码");
          }
        });
      } catch (connectError) {
        if (!disposed) {
          setError(
            connectError instanceof Error ? connectError.message : "VNC 连接初始化失败",
          );
          setStatus("控制台初始化失败");
        }
      }
    };

    connect();

    return () => {
      disposed = true;

      if (rfb) {
        try {
          rfb.disconnect();
        } catch {
          // Ignore disconnect errors during unmount.
        }
      }
    };
  }, [password, wsUrl]);

  return (
    <div className="space-y-3">
      <div className="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-[var(--border)] bg-white/70 px-4 py-3 text-sm text-[var(--muted-ink)]">
        <span>{status}</span>
        <span className="font-mono text-xs text-[var(--accent-strong)]">{wsUrl}</span>
      </div>
      {error ? (
        <div className="rounded-2xl bg-[rgba(185,76,55,0.12)] px-4 py-3 text-sm text-[#9d3e2d]">
          {error}
        </div>
      ) : null}
      <div
        ref={containerRef}
        className="min-h-[720px] overflow-hidden rounded-[28px] border border-[var(--border)] bg-[#0f172a]"
      />
    </div>
  );
}
