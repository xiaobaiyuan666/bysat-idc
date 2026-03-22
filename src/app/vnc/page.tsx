import { VncConsole } from "@/components/vnc-console";

function readParam(
  params: Record<string, string | string[] | undefined>,
  key: string,
) {
  const value = params[key];
  return typeof value === "string" ? value : undefined;
}

export default async function VncPage({
  searchParams,
}: {
  searchParams: Promise<Record<string, string | string[] | undefined>>;
}) {
  const params = await searchParams;
  const wsUrl = readParam(params, "ws");
  const password = readParam(params, "password");
  const ip = readParam(params, "ip");
  const title = readParam(params, "title");

  return (
    <main className="min-h-screen bg-[radial-gradient(circle_at_top,_rgba(59,130,246,0.15),_transparent_50%),linear-gradient(180deg,#f7fbff_0%,#eef5ff_100%)] px-6 py-8 text-[var(--ink)]">
      <div className="mx-auto max-w-7xl space-y-6">
        <section className="rounded-[32px] border border-[var(--border)] bg-white/82 px-6 py-6 shadow-[0_28px_90px_rgba(15,23,42,0.08)] backdrop-blur">
          <h1 className="text-3xl font-semibold">实例控制台</h1>
          <div className="mt-3 flex flex-wrap gap-6 text-sm text-[var(--muted-ink)]">
            <span>实例: {title ?? "-"}</span>
            <span>目标 IP: {ip ?? "-"}</span>
            <span>连接类型: WebSocket VNC</span>
          </div>
        </section>

        {!wsUrl ? (
          <section className="rounded-[32px] border border-[var(--border)] bg-white/82 px-6 py-8 text-sm text-[var(--muted-ink)] shadow-[0_28px_90px_rgba(15,23,42,0.08)] backdrop-blur">
            缺少 VNC 连接参数，请从后台或客户门户重新点击“获取 VNC”。
          </section>
        ) : (
          <VncConsole wsUrl={wsUrl} password={password} />
        )}
      </div>
    </main>
  );
}
