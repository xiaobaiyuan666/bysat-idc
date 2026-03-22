import { Card } from "@/components/ui/card";
import { StatusChip } from "@/components/ui/status-chip";
import {
  formatDateTime,
  notificationChannelLabel,
  notificationPriorityLabel,
  notificationStatusLabel,
} from "@/lib/format";
import { requirePortalUser } from "@/lib/portal-auth";
import { getPortalNotificationsData } from "@/lib/portal-data";

export default async function PortalNotificationsPage() {
  const user = await requirePortalUser();
  const notifications = await getPortalNotificationsData(user.customerId);

  return (
    <section className="space-y-6">
      <Card>
        <h2 className="text-2xl font-semibold text-[var(--ink)]">通知中心</h2>
        <p className="mt-3 text-sm leading-7 text-[var(--muted-ink)]">
          统一查看订单、账单、退款、工单和魔方云同步状态相关通知。
        </p>
      </Card>

      <Card>
        <div className="grid gap-4 md:grid-cols-4">
          <div className="rounded-[22px] bg-white/70 p-5">
            <p className="text-sm text-[var(--muted-ink)]">通知总数</p>
            <p className="mt-3 text-3xl font-semibold text-[var(--ink)]">{notifications.length}</p>
          </div>
          <div className="rounded-[22px] bg-white/70 p-5">
            <p className="text-sm text-[var(--muted-ink)]">待投递</p>
            <p className="mt-3 text-3xl font-semibold text-[var(--ink)]">
              {notifications.filter((item) => item.status === "PENDING").length}
            </p>
          </div>
          <div className="rounded-[22px] bg-white/70 p-5">
            <p className="text-sm text-[var(--muted-ink)]">已送达</p>
            <p className="mt-3 text-3xl font-semibold text-[var(--ink)]">
              {notifications.filter((item) => item.status === "SENT").length}
            </p>
          </div>
          <div className="rounded-[22px] bg-white/70 p-5">
            <p className="text-sm text-[var(--muted-ink)]">失败</p>
            <p className="mt-3 text-3xl font-semibold text-[var(--ink)]">
              {notifications.filter((item) => item.status === "FAILED").length}
            </p>
          </div>
        </div>
      </Card>

      <Card>
        <h3 className="text-xl font-semibold text-[var(--ink)]">最近通知</h3>
        <div className="mt-5 space-y-4">
          {notifications.length === 0 ? (
            <div className="rounded-[22px] bg-white/70 px-5 py-8 text-sm text-[var(--muted-ink)]">
              暂无通知记录。
            </div>
          ) : null}

          {notifications.map((item) => (
            <div
              key={item.id}
              className="rounded-[24px] border border-[var(--border)] bg-white/70 p-5"
            >
              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <div>
                  <p className="text-lg font-semibold text-[var(--ink)]">
                    {item.subject ?? item.template?.name ?? "系统通知"}
                  </p>
                  <p className="mt-1 text-sm text-[var(--muted-ink)]">
                    {notificationChannelLabel(item.channel)} /{" "}
                    {notificationPriorityLabel(item.priority)} /{" "}
                    {formatDateTime(item.sentAt ?? item.createdAt)}
                  </p>
                </div>
                <StatusChip
                  label={notificationStatusLabel(item.status)}
                  value={item.status}
                />
              </div>

              <p className="mt-4 text-sm leading-7 text-[var(--muted-ink)]">
                {item.content}
              </p>

              {item.errorMessage ? (
                <div className="mt-4 rounded-[18px] bg-[rgba(185,76,55,0.12)] px-4 py-3 text-sm text-[#9d3e2d]">
                  {item.errorMessage}
                </div>
              ) : null}
            </div>
          ))}
        </div>
      </Card>
    </section>
  );
}
